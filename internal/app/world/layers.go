package world

import (
	"github.com/bubnovdm/progenGame/internal/app/game"
	"math/rand"
)

const (
	MapSize          = 40
	PathLength       = 100
	StartSymbol      = 'S'
	ExitSymbol       = 'X'
	PathSymbol       = '1'
	EmptySymbol      = '0'
	BackgroundSymbol = 'G'
)

func GenerateMap() game.GameMap {
	var m game.GameMap

	// Инициализируем все слои пустыми значениями
	for i := 0; i < MapSize; i++ {
		for j := 0; j < MapSize; j++ {
			m.Background[i][j] = BackgroundSymbol
			m.Floor[i][j] = EmptySymbol
			m.Objects[i][j] = EmptySymbol
		}
	}

	// Рандомные начальная и конечная точки
	startX, startY := rand.Intn(MapSize), rand.Intn(MapSize)
	endX, endY := rand.Intn(MapSize), rand.Intn(MapSize)
	for startX == endX && startY == endY {
		endX, endY = rand.Intn(MapSize), rand.Intn(MapSize)
	}

	// Устанавливаем начальную и конечную точки на слое Floor
	m.Floor[startX][startY] = StartSymbol
	m.Floor[endX][endY] = ExitSymbol

	// Генерируем путь на слое Floor
	GeneratePath(&m, startX, startY, endX, endY)

	return m
}

func GeneratePath(m *game.GameMap, startX, startY, endX, endY int) {
	currentX, currentY := startX, startY
	pathLength := 0
	correctionInterval := 15 // Корректируем направление каждые 5 шагов

	for pathLength < PathLength || (currentX != endX || currentY != endY) {
		if currentX != startX || currentY != startY {
			m.Floor[currentX][currentY] = PathSymbol
		}

		if pathLength >= PathLength && currentX == endX && currentY == endY {
			break
		}

		// Корректируем направление каждые correctionInterval шагов
		if pathLength%correctionInterval == 0 {
			if currentX < endX && currentX < MapSize-1 {
				currentX++
			} else if currentX > endX && currentX > 0 {
				currentX--
			} else if currentY < endY && currentY < MapSize-1 {
				currentY++
			} else if currentY > endY && currentY > 0 {
				currentY--
			}
		} else {
			// Случайное блуждание
			switch rand.Intn(4) {
			case 0: // вверх
				if currentY > 0 {
					currentY--
				}
			case 1: // вниз
				if currentY < MapSize-1 {
					currentY++
				}
			case 2: // влево
				if currentX > 0 {
					currentX--
				}
			case 3: // вправо
				if currentX < MapSize-1 {
					currentX++
				}
			}
		}
		pathLength++
	}
	m.Floor[endX][endY] = ExitSymbol // Устанавливаем выход
}
