package app

import (
	"github.com/bubnovdm/progenGame/internal/utils"
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
