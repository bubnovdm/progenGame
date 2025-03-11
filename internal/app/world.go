package app

import (
	"github.com/bubnovdm/progenGame/internal/utils"
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
	WallSymbol       = 'W'
)

type MapType int

const (
	OpenMap MapType = iota
	DungeonMap
)

type Layer [MapSize][MapSize]rune

type GameMap struct {
	Background Layer
	Floor      Layer
	Objects    Layer
	Type       MapType
}

func GenerateMap(mapType MapType) GameMap {
	var m GameMap
	m.Type = mapType

	for i := 0; i < MapSize; i++ {
		for j := 0; j < MapSize; j++ {
			m.Background[i][j] = BackgroundSymbol
			m.Floor[i][j] = EmptySymbol
			m.Objects[i][j] = EmptySymbol
		}
	}

	if mapType == OpenMap {
		startX, startY := rand.Intn(MapSize), rand.Intn(MapSize)
		endX, endY := rand.Intn(MapSize), rand.Intn(MapSize)
		for startX == endX && startY == endY {
			endX, endY = rand.Intn(MapSize), rand.Intn(MapSize)
		}
		m.Floor[startX][startY] = StartSymbol
		m.Floor[endX][endY] = ExitSymbol
		GeneratePath(&m, startX, startY, endX, endY)
	} else {
		GenerateRoom(&m)
	}

	return m
}

func GeneratePath(m *GameMap, startX, startY, endX, endY int) {
	currentX, currentY := startX, startY
	pathLength := 0
	correctionInterval := 5

	for pathLength < PathLength || (currentX != endX || currentY != endY) {
		if currentX != startX || currentY != startY {
			m.Floor[currentX][currentY] = PathSymbol
		}

		if pathLength >= PathLength && currentX == endX && currentY == endY {
			break
		}

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
	m.Floor[endX][endY] = ExitSymbol
}

func GenerateRoom(m *GameMap) {
	const (
		maxRooms    = 10
		minRoomSize = 4
		maxRoomSize = 8
	)

	type Room struct {
		X, Y, W, H int
	}

	rooms := []Room{}

	// Шаг 1: Генерируем комнаты
	for i := 0; i < maxRooms; i++ {
		w := rand.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		h := rand.Intn(maxRoomSize-minRoomSize+1) + minRoomSize
		x := rand.Intn(MapSize-w-1) + 1
		y := rand.Intn(MapSize-h-1) + 1

		newRoom := Room{x, y, w, h}

		overlap := false
		for _, room := range rooms {
			if newRoom.X < room.X+room.W+1 && newRoom.X+newRoom.W+1 > room.X &&
				newRoom.Y < room.Y+room.H+1 && newRoom.Y+newRoom.H+1 > room.Y {
				overlap = true
				break
			}
		}

		if !overlap {
			// Добавляем комнату на слой Floor
			for ry := newRoom.Y; ry < newRoom.Y+newRoom.H; ry++ {
				for rx := newRoom.X; rx < newRoom.X+newRoom.W; rx++ {
					m.Floor[ry][rx] = PathSymbol
				}
			}
			rooms = append(rooms, newRoom)
		}
	}

	// Шаг 2: Соединяем комнаты коридорами
	for i := 0; i < len(rooms)-1; i++ {
		room1 := rooms[i]
		room2 := rooms[i+1]

		x1 := room1.X + room1.W/2
		y1 := room1.Y + room1.H/2
		x2 := room2.X + room2.W/2
		y2 := room2.Y + room2.H/2

		// Горизонтальный коридор
		for x := utils.Min(x1, x2); x <= utils.Max(x1, x2); x++ {
			m.Floor[y1][x] = PathSymbol
		}

		// Вертикальный коридор
		for y := utils.Min(y1, y2); y <= utils.Max(y1, y2); y++ {
			m.Floor[y][x2] = PathSymbol
		}
	}

	// Шаг 3: Устанавливаем начальную и конечную точки
	if len(rooms) > 0 {
		startRoom := rooms[0]
		m.Floor[startRoom.Y+startRoom.H/2][startRoom.X+startRoom.W/2] = StartSymbol

		endRoom := rooms[len(rooms)-1]
		m.Floor[endRoom.Y+endRoom.H/2][endRoom.X+endRoom.W/2] = ExitSymbol
	}

	// Шаг 4: Добавляем стены после создания всех путей
	for _, room := range rooms {
		// Добавляем стены вокруг комнаты
		for ry := room.Y - 1; ry <= room.Y+room.H; ry++ {
			for rx := room.X - 1; rx <= room.X+room.W; rx++ {
				if ry < 0 || ry >= MapSize || rx < 0 || rx >= MapSize {
					continue
				}
				// Добавляем стену только если клетка не путь, не старт и не выход
				if (ry == room.Y-1 || ry == room.Y+room.H || rx == room.X-1 || rx == room.X+room.W) &&
					m.Floor[ry][rx] != PathSymbol && m.Floor[ry][rx] != StartSymbol && m.Floor[ry][rx] != ExitSymbol {
					m.Objects[ry][rx] = WallSymbol
				}
			}
		}
	}

	// Добавляем стены вдоль коридоров
	for y := 0; y < MapSize; y++ {
		for x := 0; x < MapSize; x++ {
			if m.Floor[y][x] == PathSymbol {
				// Проверяем, является ли клетка частью коридора (рядом нет комнаты)
				isCorridor := true
				for _, room := range rooms {
					if x >= room.X && x < room.X+room.W && y >= room.Y && y < room.Y+room.H {
						isCorridor = false
						break
					}
				}
				if isCorridor {
					// Добавляем стены по бокам коридора
					if y-1 >= 0 && m.Floor[y-1][x] != PathSymbol && m.Floor[y-1][x] != StartSymbol && m.Floor[y-1][x] != ExitSymbol {
						m.Objects[y-1][x] = WallSymbol
					}
					if y+1 < MapSize && m.Floor[y+1][x] != PathSymbol && m.Floor[y+1][x] != StartSymbol && m.Floor[y+1][x] != ExitSymbol {
						m.Objects[y+1][x] = WallSymbol
					}
					if x-1 >= 0 && m.Floor[y][x-1] != PathSymbol && m.Floor[y][x-1] != StartSymbol && m.Floor[y][x-1] != ExitSymbol {
						m.Objects[y][x-1] = WallSymbol
					}
					if x+1 < MapSize && m.Floor[y][x+1] != PathSymbol && m.Floor[y][x+1] != StartSymbol && m.Floor[y][x+1] != ExitSymbol {
						m.Objects[y][x+1] = WallSymbol
					}
				}
			}
		}
	}
}
