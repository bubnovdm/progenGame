package app

import (
	//"github.com/hajimehoshi/ebiten/v2"
	"log"
	"os"
)

type Button struct {
	X, Y, Width, Height int
	Label               string
	Action              func(*Game)
	Disabled            bool
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
				g.CurrentFloor = 1
				g.selectedClassIndex = 0
				g.Player = NewPlayer(g.classes[g.selectedClassIndex], g)
				g.Player.X = 0
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
			Action: func(g *Game) {
				if g.HasSave {
					if err := g.LoadGame(); err != nil {
						log.Printf("Failed to load game: %v", err)
					} else {
						g.State = CharacterSheet
					}
				}
			},
			Disabled: !g.HasSave,
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
	return []Button{
		{
			X:      400,
			Y:      600,
			Width:  50,
			Height: 50,
			Label:  "<",
			Action: func(g *Game) {
				g.selectedClassIndex = (g.selectedClassIndex - 1 + len(g.classes)) % len(g.classes)
				g.Player = NewPlayer(g.classes[g.selectedClassIndex], g)
			},
		},
		{
			X:      550,
			Y:      600,
			Width:  50,
			Height: 50,
			Label:  ">",
			Action: func(g *Game) {
				g.selectedClassIndex = (g.selectedClassIndex + 1) % len(g.classes)
				g.Player = NewPlayer(g.classes[g.selectedClassIndex], g)
			},
		},
		{
			X:      400,
			Y:      700,
			Width:  200,
			Height: 50,
			Label:  "Play",
			Action: func(g *Game) {
				g.CurrentFloor = g.SelectedFloor
				var mapType MapType
				if g.CurrentFloor%2 == 0 {
					mapType = OpenMap
				} else {
					mapType = DungeonMap
				}
				g.GameMap = GenerateMap(mapType)
				g.moveToStartPosition()
				g.spawnEnemies()
				g.State = Dungeon
			},
		},
	}
}
