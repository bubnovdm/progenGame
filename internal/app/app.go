package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
	"math/rand"
	"time"
)

type GameState uint8

const (
	Menu GameState = iota
	CharacterSheet
	Dungeon
	InGameMenu
	CombatState
)

type DotEffect struct {
	DamagePerTick int
	Duration      float64
	TickInterval  float64
	TimeRemaining float64
}

type RapidShotEffect struct {
	DamagePerHit  int
	HitsRemaining int
	HitInterval   float64
	TimeUntilNext float64
}

type Game struct {
	GameMap               GameMap                       // Игровая карта
	Player                Player                        // Игрок
	Enemies               []Enemy                       // Слайс врагов
	moveDelay             int                           // Ограничение кадров на шаг
	textures              map[rune]*ebiten.Image        // Мапа для текстур
	playerImage           *ebiten.Image                 // Изображение игрока
	enemyImage            *ebiten.Image                 // Изображение врага на карте
	enemyLargeImages      map[string]*ebiten.Image      // Мапа изображений врагов в бою
	Level                 int                           // Уровень карты
	State                 GameState                     // Статус экрана
	selectedClassIndex    int                           // Индекс текущего выбранного класса
	classes               []PlayerClass                 // Список доступных классов
	classImages           map[PlayerClass]*ebiten.Image // Мини-иконки для карты
	characterImages       map[PlayerClass]*ebiten.Image // Большие изображения для CharacterSheet
	backgroundImage       *ebiten.Image                 // Фоновое изображение
	CurrentEnemy          *Enemy                        // Враг, с которым идёт бой
	AutoAttackCooldown    float64                       // Таймер для автоатаки
	AbilityCooldowns      map[string]float64            // Таймеры для способностей
	CombatLog             []string                      // Добавляем поле для лога боя
	combatBackgroundImage *ebiten.Image                 // Новое поле для фона боя
	ActiveDotEffect       *DotEffect                    // Эффект урона со временем
	ActiveRapidShot       *RapidShotEffect              // Эффект последовательных ударов
	EnemyAttackCooldown   float64                       // Кд атак врагов
}

func (g *Game) Update() error {
	dt := 1.0 / 60.0 // Предполагаем 60 FPS

	// Обновляем кулдауны
	if g.AutoAttackCooldown > 0 {
		g.AutoAttackCooldown -= dt
		if g.AutoAttackCooldown < 0 {
			g.AutoAttackCooldown = 0
		}
	}

	for key := range g.AbilityCooldowns {
		if g.AbilityCooldowns[key] > 0 {
			g.AbilityCooldowns[key] -= dt
			if g.AbilityCooldowns[key] < 0 {
				g.AbilityCooldowns[key] = 0
			}
		}
	}

	// Обработка DoT-эффекта
	if g.ActiveDotEffect != nil && g.CurrentEnemy != nil {
		g.ActiveDotEffect.TimeRemaining -= dt
		if g.ActiveDotEffect.TimeRemaining <= 0 {
			g.ActiveDotEffect = nil
		} else {
			g.ActiveDotEffect.TickInterval -= dt
			if g.ActiveDotEffect.TickInterval <= 0 {
				g.CurrentEnemy.HP -= g.ActiveDotEffect.DamagePerTick
				g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s takes %d burn damage. Enemy HP: %d", g.CurrentEnemy.Name, g.ActiveDotEffect.DamagePerTick, g.CurrentEnemy.HP))
				g.ActiveDotEffect.TickInterval = 1.0 // Сбрасываем интервал
				if g.CurrentEnemy.HP <= 0 {
					g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
					enemyID := g.CurrentEnemy.ID
					g.CurrentEnemy = nil
					g.ActiveDotEffect = nil
					g.ActiveRapidShot = nil
					g.State = Dungeon
					g.Enemies = removeEnemy(g.Enemies, enemyID)
				}
			}
		}
	}

	// Обработка Rapid Shot
	if g.ActiveRapidShot != nil && g.CurrentEnemy != nil {
		g.ActiveRapidShot.TimeUntilNext -= dt
		if g.ActiveRapidShot.TimeUntilNext <= 0 && g.ActiveRapidShot.HitsRemaining > 0 {
			g.CurrentEnemy.HP -= g.ActiveRapidShot.DamagePerHit
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("Rapid Shot hits %s for %d damage. Enemy HP: %d", g.CurrentEnemy.Name, g.ActiveRapidShot.DamagePerHit, g.CurrentEnemy.HP))
			g.ActiveRapidShot.HitsRemaining--
			g.ActiveRapidShot.TimeUntilNext = g.ActiveRapidShot.HitInterval
			if g.ActiveRapidShot.HitsRemaining <= 0 {
				g.ActiveRapidShot = nil
			}
			if g.CurrentEnemy.HP <= 0 {
				g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
				enemyID := g.CurrentEnemy.ID
				g.CurrentEnemy = nil
				g.ActiveDotEffect = nil
				g.ActiveRapidShot = nil
				g.State = Dungeon
				g.Enemies = removeEnemy(g.Enemies, enemyID)
			}
		}
	}

	// Контратака врага
	if g.CurrentEnemy != nil && g.State == CombatState {
		g.EnemyAttackCooldown -= dt
		if g.EnemyAttackCooldown <= 0 {
			enemyDamage := int(g.CurrentEnemy.Strength)
			var defense int
			defense = int(g.Player.PhDefense) * 2
			effectiveDamage := enemyDamage - defense
			if effectiveDamage < 0 {
				effectiveDamage = 0
			}
			oldHP := g.Player.HP // Сохраняем старое значение HP для отладки
			g.Player.HP -= uint16(effectiveDamage)
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s counterattacks for %d damage. Player HP: %d", g.CurrentEnemy.Name, effectiveDamage, g.Player.HP))
			fmt.Printf("Player HP changed from %d to %d\n", oldHP, g.Player.HP) // Отладка
			g.EnemyAttackCooldown = 2.0
			if g.Player.HP <= 0 {
				g.CombatLog = append(g.CombatLog, "Player defeated! Game Over.")
				g.State = Menu
				g.Player = NewPlayer(WarriorClass)
				g.CurrentEnemy = nil
				g.ActiveDotEffect = nil
				g.ActiveRapidShot = nil
			}
		}
	}

	switch g.State {
	case Menu:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for _, button := range g.getMenuButtons() {
				if mx >= button.X && mx <= button.X+button.Width &&
					my >= button.Y && my <= button.Y+button.Height {
					button.Action(g)
				}
			}
		}

	case CharacterSheet:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for _, button := range g.getCharacterSheetButtons() {
				if mx >= button.X && mx <= button.X+button.Width &&
					my >= button.Y && my <= button.Y+button.Height {
					button.Action(g)
				}
			}
		}

	case InGameMenu:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for _, button := range g.getInGameMenuButtons() {
				if mx >= button.X && mx <= button.X+button.Width &&
					my >= button.Y && my <= button.Y+button.Height {
					button.Action(g)
				}
			}
		}

	case CombatState:
		if g.CurrentEnemy == nil {
			g.State = Dungeon
			return nil
		}

		if g.AbilityCooldowns["BasicAttack"] > 0 {
			g.AbilityCooldowns["BasicAttack"] -= 1.0 / 60.0
		}
		if g.AutoAttackCooldown > 0 {
			g.AutoAttackCooldown -= 1.0 / 60.0
		}
		if g.AutoAttackCooldown <= 0 && g.CurrentEnemy != nil { // Добавляем проверку
			g.autoAttack()
			g.AutoAttackCooldown = 2.0
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = Dungeon
			g.CurrentEnemy = nil
			return nil
		}
		if inpututil.IsKeyJustPressed(ebiten.Key1) {
			g.useAbility("1")
		}

	case Dungeon:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = InGameMenu
			return nil
		}

		if g.moveDelay > 0 {
			g.moveDelay--
			return nil
		}

		if g.GameMap.Floor[g.Player.Y][g.Player.X] == ExitSymbol {
			g.Level++
			var mapType MapType
			if g.Level%2 == 0 {
				mapType = OpenMap
			} else {
				mapType = DungeonMap
			}
			g.GameMap = GenerateMap(mapType)
			g.moveToStartPosition()
			g.spawnEnemies()
			g.moveDelay = 10
			return nil
		}

		newX, newY := g.Player.X, g.Player.Y

		if ebiten.IsKeyPressed(ebiten.KeyW) && g.Player.Y > 0 {
			newY--
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && g.Player.Y < MapSize-1 {
			newY++
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) && g.Player.X > 0 {
			newX--
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) && g.Player.X < MapSize-1 {
			newX++
		}

		if newX != g.Player.X || newY != g.Player.Y {
			if g.IsWalkable(newX, newY) {
				g.Player.X = newX
				g.Player.Y = newY
				g.moveDelay = 10
			}
			// Проверка столкновения с врагом
			for _, enemy := range g.Enemies {
				if g.isAdjacent(g.Player.X, g.Player.Y, enemy.X, enemy.Y) {
					g.CurrentEnemy = &enemy
					g.State = CombatState
					g.AutoAttackCooldown = 2.0 // Таймер автоатаки в секундах
					return nil
				}
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case Menu:
		g.drawMenu(screen)
	case CharacterSheet:
		g.drawCharacterSheet(screen)
	case InGameMenu:
		g.drawInGameMenu(screen)
	case Dungeon:
		g.drawDungeon(screen)
	case CombatState:
		g.drawCombat(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func Start() {
	rand.Seed(time.Now().UnixNano())

	// Загружаем конфигурации
	if err := LoadEnemyConfigs("assets/enemies/enemies.json"); err != nil {
		log.Fatalf("Failed to load enemy configs: %v", err)
	}
	if err := LoadClassConfigs("assets/classes/classes.json"); err != nil {
		log.Fatalf("Failed to load class configs: %v", err)
	}
	if err := LoadAbilityConfigs("assets/abilities/abilities.json"); err != nil {
		log.Fatalf("Failed to load ability configs: %v", err)
	}

	game := &Game{
		Player:           NewPlayer(WarriorClass),
		textures:         make(map[rune]*ebiten.Image),
		Level:            1,
		State:            Menu,
		classes:          []PlayerClass{WarriorClass, MageClass, ArcherClass},
		classImages:      make(map[PlayerClass]*ebiten.Image),
		characterImages:  make(map[PlayerClass]*ebiten.Image),
		AbilityCooldowns: make(map[string]float64),
		CombatLog:        []string{},
	}

	loadAssets(game)
	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Map Generator")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

// Базовые методы
func (g *Game) isAdjacent(x1, y1, x2, y2 int) bool {
	return (x1 == x2 && (y1 == y2-1 || y1 == y2+1)) || // Вертикально
		(y1 == y2 && (x1 == x2-1 || x1 == x2+1)) // Горизонтально
}
