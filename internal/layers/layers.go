package layers

import (
	"math/rand"
)

const (
	MapSize     = 40
	PathLength  = 50
	StartSymbol = 'S'
	ExitSymbol  = 'X'
	PathSymbol  = '1'
	EmptySymbol = '0'
)

type Map [MapSize][MapSize]rune

func GenerateMap() Map {
	var m Map

	// Инициализируем карту пустыми значениями
	for i := range m {
		for j := range m[i] {
			m[i][j] = EmptySymbol
		}
	}

	// Рандомно определим начальную и конечную точки
	startX, startY := rand.Intn(MapSize), rand.Intn(MapSize)
	endX, endY := rand.Intn(MapSize), rand.Intn(MapSize)

	// Убедимся, что начальная и конечная точки не совпадают
	for startX == endX && startY == endY {
		endX, endY = rand.Intn(MapSize), rand.Intn(MapSize)
	}

	// Устанавливаем начальную и конечную точки на карте
	m[startX][startY] = StartSymbol
	m[endX][endY] = ExitSymbol

	// Генерируем путь между точками
	GeneratePath(&m, startX, startY, endX, endY)

	return m
}

func GeneratePath(m *Map, startX, startY, endX, endY int) {
	currentX, currentY := startX, startY
	pathLength := 0

	for pathLength < PathLength || (currentX != endX || currentY != endY) {
		// Убедимся, что мы не перезаписываем начальную точку
		if currentX != startX || currentY != startY {
			m[currentX][currentY] = PathSymbol
		}

		if pathLength >= PathLength && currentX == endX && currentY == endY {
			break
		}

		// Определяем направление движения
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
		pathLength++
	}
}
