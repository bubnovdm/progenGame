package app

import (
	//"github.com/hajimehoshi/ebiten/v2"
	"os"
)

type Button struct {
	X, Y, Width, Height int
	Label               string
	Action              func(*Game)
}

func (g *Game) getMenuButtons() []Button {
	return []Button{
		{
			X:      400,
			Y:      300,
			Width:  200,
			Height: 50,
			Label:  "New Game",
			Action: func(g *Game) {
				g.Level = 1
				g.selectedClassIndex = 0                              // По умолчанию первый класс (Воин)
				g.Player = NewPlayer(g.classes[g.selectedClassIndex]) // Сбрасываем характеристики
				g.Player.X = 0                                        // Устанавливаем начальные координаты
				g.Player.Y = 0
				g.State = CharacterSheet
			},
		},
		{
			X:      400,
			Y:      400,
			Width:  200,
			Height: 50,
			Label:  "Continue",
			Action: func(g *Game) {},
		},
		{
			X:      400,
			Y:      500,
			Width:  200,
			Height: 50,
			Label:  "Exit",
			Action: func(g *Game) {
				os.Exit(0)
			},
		},
	}
}

func (g *Game) getInGameMenuButtons() []Button {
	return []Button{
		{
			X:      400,
			Y:      400,
			Width:  200,
			Height: 50,
			Label:  "Cancel",
			Action: func(g *Game) {
				g.State = Dungeon
			},
		},
		{
			X:      400,
			Y:      500,
			Width:  200,
			Height: 50,
			Label:  "Exit",
			Action: func(g *Game) {
				os.Exit(0)
			},
		},
	}
}

func (g *Game) getCharacterSheetButtons() []Button {
	currentX, currentY := g.Player.X, g.Player.Y // Сохраняем текущие координаты
	return []Button{
		{
			X:      400, // Кнопка "влево"
			Y:      600,
			Width:  60,
			Height: 60,
			Label:  "<",
			Action: func(g *Game) {
				g.selectedClassIndex = (g.selectedClassIndex - 1 + len(g.classes)) % len(g.classes)
				tempPlayer := NewPlayer(g.classes[g.selectedClassIndex])
				g.Player = tempPlayer // Полное обновление игрока, включая Inventory
				g.Player.X = currentX // Восстанавливаем координаты
				g.Player.Y = currentY
				g.playerImage = g.classImages[g.Player.Class] // Обновляем изображение игрока
			},
		},
		{
			X:      540, // Кнопка "вправо"
			Y:      600,
			Width:  60,
			Height: 60,
			Label:  ">",
			Action: func(g *Game) {
				g.selectedClassIndex = (g.selectedClassIndex + 1) % len(g.classes)
				tempPlayer := NewPlayer(g.classes[g.selectedClassIndex])
				g.Player = tempPlayer // Полное обновление игрока, включая Inventory
				g.Player.X = currentX // Восстанавливаем координаты
				g.Player.Y = currentY
				g.playerImage = g.classImages[g.Player.Class] // Обновляем изображение игрока
			},
		},
		{
			X:      450, // Кнопка "Play"
			Y:      700,
			Width:  100,
			Height: 60,
			Label:  "Play",
			Action: func(g *Game) {
				var mapType MapType
				if g.Level%2 == 0 {
					mapType = OpenMap
				} else {
					mapType = DungeonMap
				}
				g.GameMap = GenerateMap(mapType)
				g.moveToStartPosition() // Устанавливает X, Y для Dungeon
				g.spawnEnemies()
				g.State = Dungeon
			},
		},
	}
}
