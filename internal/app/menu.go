package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
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
				g.Player = Player{
					ID:           "player1",
					Name:         "Hero",
					HP:           100,
					Mana:         50,
					Level:        1,
					Strength:     10,
					Agility:      10,
					Intelligence: 10,
				}
				var mapType MapType
				if g.Level%2 == 0 {
					mapType = OpenMap
				} else {
					mapType = DungeonMap
				}
				g.GameMap = GenerateMap(mapType)
				g.moveToStartPosition()
				g.State = Dungeon
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

func (g *Game) drawMenu(screen *ebiten.Image) {
	screen.Fill(color.Black)
	ebitenutil.DebugPrint(screen, "Main Menu")

	for _, button := range g.getMenuButtons() {
		buttonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
		mx, my := ebiten.CursorPosition()
		if mx >= button.X && mx <= button.X+button.Width &&
			my >= button.Y && my <= button.Y+button.Height {
			buttonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		}
		buttonImage := ebiten.NewImage(button.Width, button.Height)
		buttonImage.Fill(buttonColor)
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(button.X), float64(button.Y))
		screen.DrawImage(buttonImage, op)
		ebitenutil.DebugPrintAt(screen, button.Label, button.X+20, button.Y+15)
	}
}

func (g *Game) drawInGameMenu(screen *ebiten.Image) {
	// Рисуем карту на фоне
	g.drawDungeon(screen)

	// Полупрозрачный фон для меню
	overlay := ebiten.NewImage(1000, 1000)
	overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 200})
	op := &ebiten.DrawImageOptions{}
	screen.DrawImage(overlay, op)

	// Заголовок
	ebitenutil.DebugPrintAt(screen, "Pause Menu", 450, 300)

	// Отрисовка кнопок
	for _, button := range g.getInGameMenuButtons() {
		buttonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
		mx, my := ebiten.CursorPosition()
		if mx >= button.X && mx <= button.X+button.Width &&
			my >= button.Y && my <= button.Y+button.Height {
			buttonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		}
		buttonImage := ebiten.NewImage(button.Width, button.Height)
		buttonImage.Fill(buttonColor)
		op = &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(button.X), float64(button.Y))
		screen.DrawImage(buttonImage, op)
		ebitenutil.DebugPrintAt(screen, button.Label, button.X+20, button.Y+15)
	}
}
