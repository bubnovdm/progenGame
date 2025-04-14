package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"golang.org/x/image/font"
	"log"
	"os"
	"sync"
	"unsafe"
)

type GameState uint8

const (
	Menu GameState = iota
	CharacterSheet
	Dungeon
	InGameMenu
	CombatState
)

type Game struct {
	// Состояние игры
	State               GameState // Текущее состояние игры (Menu, Dungeon и т.д.)
	HasSave             bool      // Есть ли сохранение
	FloorSelectorOpen   bool      // Открыт ли выпадающий список этажей
	moveDelay           uint8     // Задержка движения (кадры)
	CurrentFloor        int       // Текущий этаж
	MaxFloor            int       // Максимальный достигнутый этаж
	SelectedFloor       int       // Выбранный этаж в выпадающем списке
	AutoAttackCooldown  float64   // Кулдаун автоатаки
	EnemyAttackCooldown float64   // Кулдаун атаки врагов

	// Данные игрока
	Player             Player                 // Игрок
	selectedClassIndex int                    // Индекс текущего выбранного класса
	classes            []PlayerClass          // Список доступных классов
	ClassConfig        map[string]ClassConfig // Конфигурации классов
	AbilityCooldowns   map[string]float64     // Кулдауны способностей
	AvailableBuffs     []Buff                 // Бафы, которые можно выбрать

	// Данные врагов
	Enemies      []Enemy // Список врагов на карте
	CurrentEnemy *Enemy  // Текущий враг в бою

	// Данные карты
	GameMap GameMap // Игровая карта

	// Ресурсы (изображения и текстуры)
	textures              map[rune]*ebiten.Image        // Текстуры для карты (трава, стены и т.д.)
	classImages           map[PlayerClass]*ebiten.Image // Мини-иконки для карты
	characterImages       map[PlayerClass]*ebiten.Image // Большие изображения для CharacterSheet и боя
	enemyImage            *ebiten.Image                 // Изображение врага на карте
	enemyLargeImages      map[string]*ebiten.Image      // Изображения врагов в бою
	backgroundImage       *ebiten.Image                 // Фоновое изображение для меню
	combatBackgroundImage *ebiten.Image                 // Фоновое изображение для боя

	// UI
	CombatLog   []string               // Лог боя
	DamageStats map[string]*DamageStat // Типа рекаунт :D

	// Шрифт
	Font font.Face // Поле для хранения шрифта
}

func (g *Game) Update() error {
	// Рассчитываем dt на основе реального времени
	tps := ebiten.ActualTPS()
	if tps == 0 {
		tps = 144.0 // Значение по умолчанию, если TPS ещё не определён
	}
	dt := 1.0 / tps

	// Обновляем кулдауны
	g.updateCooldowns(dt)

	// Обработка эффектов на текущем враге
	if g.CurrentEnemy != nil {
		// Обрабатываем все эффекты
		for i := 0; i < len(g.CurrentEnemy.ActiveEffects); i++ {
			effect := g.CurrentEnemy.ActiveEffects[i]
			logMessage := effect.Update(dt, g.CurrentEnemy, g)
			if logMessage != "" {
				g.CombatLog = append(g.CombatLog, logMessage)
			}

			// Проверяем, не умер ли враг от эффекта
			if g.CurrentEnemy.HP <= 0 {
				g.HandleEnemyDeath()
				break // Прерываем цикл, так как враг мёртв
			}

			// Удаляем завершённые эффекты
			if effect.IsFinished() {
				g.CurrentEnemy.ActiveEffects = append(g.CurrentEnemy.ActiveEffects[:i], g.CurrentEnemy.ActiveEffects[i+1:]...)
				i-- // Уменьшаем индекс, так как слайс сдвинулся
			}
		}
	}

	// Контратака врага
	if g.CurrentEnemy != nil && g.State == CombatState {
		g.EnemyAttackCooldown -= dt
		if g.EnemyAttackCooldown <= 0 {
			enemyDamage := int(g.CurrentEnemy.Strength)
			var defense int
			defense = int(g.Player.PhDefense) // Предполагаем, что враги наносят физический урон
			effectiveDamage := int(float64(enemyDamage) * (100.0 / (100.0 + float64(defense))))
			if effectiveDamage < minimalEnemyDamage {
				effectiveDamage = minimalEnemyDamage // Минимальный урон 3
			}
			oldHP := g.Player.HP
			g.Player.HP -= uint16(effectiveDamage)
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s counterattacks for %d damage. Player HP: %d", g.CurrentEnemy.Name, effectiveDamage, g.Player.HP))
			fmt.Printf("Player HP changed from %d to %d\n", oldHP, g.Player.HP)
			g.EnemyAttackCooldown = 2.0
			if g.Player.HP <= 0 {
				g.CombatLog = append(g.CombatLog, "Player defeated! Game Over.")
				g.State = Menu
				g.Player = NewPlayer(WarriorClass, g)
				g.CurrentEnemy = nil
			}
		}
	}

	switch g.State {
	case Menu:
		g.updateMenu()

	case CharacterSheet:
		g.updateCharacterSheet()

	case InGameMenu:
		g.updateInGameMenu()

	case CombatState:
		g.updateCombat(dt)

	case Dungeon:
		g.updateDungeon()
	}
	return nil
}

// Нужен для имплиментации интерфейса ebiten
func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case Menu:
		g.drawMenu(screen)
	case CharacterSheet:
		g.drawCharacterSheet(screen)
	case InGameMenu:
		g.drawInGameMenu(screen)
		g.drawDamageStats(screen)
	case Dungeon:
		g.drawDungeon(screen)
		g.drawDamageStats(screen)
	case CombatState:
		g.drawCombat(screen)
		g.drawDamageStats(screen)
	}
}

// Нужен для имплиментации интерфейса ebiten
func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func Start() {
	// Загружаем конфигурации
	var wg sync.WaitGroup
	wg.Add(3)
	go func() {
		defer wg.Done()
		if err := LoadEnemyConfigs("assets/enemies/enemies.json"); err != nil {
			log.Fatalf("Failed to load enemy configs: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := LoadClassConfigs("assets/classes/classes.json"); err != nil {
			log.Fatalf("Failed to load class configs: %v", err)
		}
	}()

	go func() {
		defer wg.Done()
		if err := LoadAbilityConfigs("assets/abilities/abilities.json"); err != nil {
			log.Fatalf("Failed to load ability configs: %v", err)
		}

	}()
	wg.Wait()

	// Надо бы как-то доработать
	/*
		if err := LoadFont("assets/fonts/PressStart2P-Regular.ttf", 24); err != nil {
			log.Fatalf("Failed to load font: %v", err)
		}
	*/

	game := &Game{
		textures:         make(map[rune]*ebiten.Image),
		CurrentFloor:     1,
		MaxFloor:         1,
		SelectedFloor:    1,
		State:            Menu,
		classes:          []PlayerClass{WarriorClass, MageClass, ArcherClass},
		classImages:      make(map[PlayerClass]*ebiten.Image),
		characterImages:  make(map[PlayerClass]*ebiten.Image),
		AbilityCooldowns: make(map[string]float64),
		CombatLog:        []string{},
		enemyLargeImages: make(map[string]*ebiten.Image),
		ClassConfig:      ToMap(),
		AvailableBuffs:   make([]Buff, 0), // Инициализируем пустой слайс
		DamageStats:      map[string]*DamageStat{},
	}

	// !!! Заглушка для первого отображения кд способностей, мб есть вариант получше, надо подумать
	game.AbilityCooldowns = map[string]float64{
		"1": 0,
		"2": 0,
		"3": 0,
		"4": 0,
	}

	// Выводим размер структуры Game
	fmt.Printf("Size of Game struct: %d bytes\n", unsafe.Sizeof(*game))
	fmt.Printf(" |-Size of GameMap struct: %d bytes\n", unsafe.Sizeof(game.GameMap))
	fmt.Printf(" |-Size of Enemies struct: %d bytes\n", unsafe.Sizeof(game.Enemies))
	fmt.Printf(" |-Size of Player struct: %d bytes\n", unsafe.Sizeof(game.Player))
	fmt.Printf("   |-Size of Player struct (Inventory): %d bytes\n", unsafe.Sizeof(game.Player.Inventory))
	fmt.Printf("   |-Size of Player struct (DamageType): %d bytes\n", unsafe.Sizeof(game.Player.DamageType))
	fmt.Printf(" |-Size of class images: %d bytes\n", unsafe.Sizeof(game.classImages))
	fmt.Printf(" |-Size of Class config: %d bytes\n", unsafe.Sizeof(game.ClassConfig))

	// Проверяем наличие сохранений
	_, err := os.Stat("save.json")
	game.HasSave = err == nil // Добавим поле HasSave в Game

	// Устанавливаем игрока после создания game
	game.Player = NewPlayer(WarriorClass, game)

	loadAssets(game)
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Map Generator")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
