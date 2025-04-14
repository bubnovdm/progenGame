package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"log"
)

/*
Перенёс из app.go из Update() сюда всякие разные блоки кода, чтобы app.go не выглядел таким монструозным
*/

// Обрабатываем кд автоатак и способностей
func (g *Game) updateCooldowns(dt float64) {
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
}

func (g *Game) updateMenu() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, button := range g.getMenuButtons() {
			if button.Disabled {
				continue // Пропускаем отключенные кнопки
			}
			if mx >= button.X && mx <= button.X+button.Width &&
				my >= button.Y && my <= button.Y+button.Height {
				button.Action(g)
				break
			}
		}
	}
}

func (g *Game) updateCharacterSheet() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()

		// Обработка кнопок
		for _, button := range g.getCharacterSheetButtons() {
			if mx >= button.X && mx <= button.X+button.Width &&
				my >= button.Y && my <= button.Y+button.Height {
				button.Action(g)
				break
			}
		}

		// Нажатие на кнопку "Floor"
		if mx >= floorButtonX && mx <= floorButtonX+floorButtonWidth &&
			my >= floorButtonY && my <= floorButtonY+floorButtonHeight {
			g.FloorSelectorOpen = !g.FloorSelectorOpen
		}

		// Выбор этажа из выпадающего списка
		if g.FloorSelectorOpen {
			for i := 1; i <= g.MaxFloor; i++ {
				optionY := floorButtonY + floorButtonHeight*i
				if mx >= floorButtonX && mx <= floorButtonX+floorButtonWidth &&
					my >= optionY && my <= optionY+floorButtonHeight {
					g.SelectedFloor = i
					g.FloorSelectorOpen = false
					break
				}
			}
		}
	}
}

func (g *Game) updateInGameMenu() {
	if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
		mx, my := ebiten.CursorPosition()
		for _, button := range g.getInGameMenuButtons() {
			if mx >= button.X && mx <= button.X+button.Width &&
				my >= button.Y && my <= button.Y+button.Height {
				button.Action(g)
			}
		}
	}
}

func (g *Game) updateCombat(dt float64) error {
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
	if g.AutoAttackCooldown <= 0 && g.CurrentEnemy != nil {
		g.autoAttack()
		// Проверяем, не умер ли враг после автоатаки
		if g.CurrentEnemy.HP <= 0 {
			g.HandleEnemyDeath()
		} else {
			// Учитываем множитель скорости автоатак, если враг еще жив
			g.AutoAttackCooldown = 2.0 * g.Player.AutoAttackCooldownMultiplier
		}
	}

	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		// Удаляем врага из списка, чтобы он не мешал
		if g.CurrentEnemy != nil {
			g.Enemies = removeEnemy(g.Enemies, g.CurrentEnemy.ID)
			fmt.Printf("Enemies after escape: %d\n", len(g.Enemies))
			if len(g.Enemies) == 0 {
				newBuff := GetRandomBuff()
				g.AvailableBuffs = append(g.AvailableBuffs, newBuff)
				g.CombatLog = append(g.CombatLog, fmt.Sprintf("Received buff: %s", newBuff.Name()))
				newBuff.Apply(&g.Player)
			}
		}
		g.State = Dungeon
		g.CurrentEnemy = nil
		return nil
	}

	abilityKeys := map[ebiten.Key]string{
		ebiten.Key1: "1",
		ebiten.Key2: "2",
		ebiten.Key3: "3",
		ebiten.Key4: "4",
	}

	for key, ability := range abilityKeys {
		if inpututil.IsKeyJustPressed(key) {
			g.useAbility(ability)
			if g.CurrentEnemy != nil && g.CurrentEnemy.HP <= 0 {
				g.HandleEnemyDeath()
			}
		}
	}
	return nil
}

func (g *Game) updateDungeon() error {
	if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
		g.State = InGameMenu
		return nil
	}

	if g.moveDelay > 0 {
		g.moveDelay--
		return nil
	}

	if g.GameMap.Floor[g.Player.Y][g.Player.X] == ExitSymbol {
		g.CurrentFloor++
		if g.CurrentFloor > g.MaxFloor {
			g.MaxFloor = g.CurrentFloor // Обновляем максимальный этаж
		}
		var mapType MapType
		if g.CurrentFloor%2 == 0 {
			mapType = OpenMap
		} else {
			mapType = DungeonMap
		}
		g.GameMap = GenerateMap(mapType)
		g.moveToStartPosition()
		g.spawnEnemies()
		// Отладка: сколько врагов заспавнилось на новом этаже
		fmt.Printf("Spawned %d enemies on floor %d\n", len(g.Enemies), g.CurrentFloor)
		g.moveDelay = 10

		// Сохраняем игру
		go func() {
			if err := g.SaveGame(); err != nil {
				log.Printf("Save error: %v", err)
			}
		}()

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

	// Переходим в режим боя, если сталкиваемся с врагом
	if newX != g.Player.X || newY != g.Player.Y {
		if g.IsWalkable(newX, newY) {
			g.Player.X = newX
			g.Player.Y = newY
			g.moveDelay = 10
		}
		for _, enemy := range g.Enemies {
			if g.isAdjacent(g.Player.X, g.Player.Y, enemy.X, enemy.Y) {
				g.CurrentEnemy = &enemy
				g.State = CombatState
				g.AutoAttackCooldown = 2.0 * g.Player.AutoAttackCooldownMultiplier
				return nil
			}
		}
	}
	return nil
}
