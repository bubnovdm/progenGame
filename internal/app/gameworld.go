package app

import (
	"fmt"
	"github.com/bubnovdm/progenGame/internal/utils"
	"log"
	"math/rand"
)

func (g *Game) IsWalkable(x, y int) bool {
	return (g.GameMap.Floor[y][x] == PathSymbol ||
		g.GameMap.Floor[y][x] == StartSymbol ||
		g.GameMap.Floor[y][x] == ExitSymbol) &&
		g.GameMap.Objects[y][x] != WallSymbol
}

func (g *Game) moveToStartPosition() {
	for y := 0; y < MapSize; y++ {
		for x := 0; x < MapSize; x++ {
			if g.GameMap.Floor[y][x] == StartSymbol {
				g.Player.X = x
				g.Player.Y = y
				return
			}
		}
	}
	g.Player.X = 0
	g.Player.Y = 0
}

func (g *Game) countFloorTiles() int {
	count := 0
	for y := 0; y < MapSize; y++ {
		for x := 0; x < MapSize; x++ {
			if g.GameMap.Floor[y][x] == PathSymbol {
				count++
			}
		}
	}
	return count
}

func (g *Game) spawnEnemies() {
	g.Enemies = nil // Очищаем старых врагов

	floorTiles := g.countFloorTiles()
	// 1 враг на каждые 25 клеток пола, минимум 1 враг
	enemyCount := utils.Max(1, floorTiles/25)

	startX, startY := -1, -1
	exitX, exitY := -1, -1

	// Находим стартовую и конечную позиции
	for y := 0; y < MapSize; y++ {
		for x := 0; x < MapSize; x++ {
			if g.GameMap.Floor[y][x] == StartSymbol {
				startX, startY = x, y
			}
			if g.GameMap.Floor[y][x] == ExitSymbol {
				exitX, exitY = x, y
			}
		}
	}

	// Размещаем врагов
	for i := 0; i < enemyCount; i++ {
		for {
			x := rand.Intn(MapSize)
			y := rand.Intn(MapSize)

			// Проверяем, что клетка — пол, не старт, не выход и не занята другим врагом
			if g.GameMap.Floor[y][x] == PathSymbol &&
				!(x == startX && y == startY) &&
				!(x == exitX && y == exitY) &&
				!g.isEnemyAt(x, y) {
				enemy := NewEnemy(x, y, g.Level)
				g.Enemies = append(g.Enemies, enemy)
				break
			}
		}
	}
}

func (g *Game) isEnemyAt(x, y int) bool {
	for _, enemy := range g.Enemies {
		if enemy.X == x && enemy.Y == y {
			return true
		}
	}
	return false
}

// NewEnemy создаёт нового врага на основе уровня.
// Если на заданном уровне доступно несколько типов врагов, выбирает случайный.
// Если конфигурация для уровня не найдена, использует запасной вариант.
func NewEnemy(x, y int, level int) Enemy {
	// Находим все конфигурации врагов для текущего уровня
	var possibleConfigs []*EnemyConfig
	for _, config := range enemyConfigs {
		for _, lvl := range config.Levels {
			if lvl == level {
				possibleConfigs = append(possibleConfigs, &config)
				break
			}
		}
	}

	// Если конфигурации не найдены, используем запасной вариант
	if len(possibleConfigs) == 0 {
		log.Printf("Warning: No enemy config found for level %d, using default Goblin", level)
		maxHP := 30 + (level * 5)
		return Enemy{
			ID:           fmt.Sprintf("enemy_%d", rand.Int()),
			Name:         "Goblin",
			X:            x,
			Y:            y,
			HP:           maxHP,
			MaxHP:        maxHP,
			Strength:     5,
			Agility:      3,
			Intelligence: 2,
			PhDefense:    2,
			MgDefense:    2,
		}
	}

	// Выбираем случайного врага из доступных для уровня
	config := possibleConfigs[rand.Intn(len(possibleConfigs))]

	// Вычисляем HP с учётом уровня
	maxHP := config.BaseStats.HP + (config.BaseStats.HPPerLevel * (level - 1))
	return Enemy{
		ID:           fmt.Sprintf("enemy_%d", rand.Int()),
		Name:         config.BaseStats.Name,
		X:            x,
		Y:            y,
		HP:           maxHP,
		MaxHP:        maxHP,
		Strength:     config.BaseStats.Strength,
		Agility:      config.BaseStats.Agility,
		Intelligence: config.BaseStats.Intelligence,
		PhDefense:    config.BaseStats.PhDefense,
		MgDefense:    config.BaseStats.MgDefense,
	}
}

// Тригеррит бой с врагом
func (g *Game) isAdjacent(x1, y1, x2, y2 int) bool {
	dx := x1 - x2
	dy := y1 - y2
	return utils.Abs(dx) <= 1 && utils.Abs(dy) <= 1
}
