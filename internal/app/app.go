package app

import (
	"github.com/bubnovdm/progenGame/internal/layers"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
	"time"
)

type Game struct {
	GameMap    layers.Map
	startImage *ebiten.Image
	exitImage  *ebiten.Image
	pathImage  *ebiten.Image
	emptyImage *ebiten.Image
}

func (g *Game) Update() error {
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Проходим по карте и отрисовываем её
	for y, row := range g.GameMap {
		for x, cell := range row {
			var img *ebiten.Image

			switch cell {
			case layers.StartSymbol:
				img = g.startImage
			case layers.ExitSymbol:
				img = g.exitImage
			case layers.PathSymbol:
				img = g.pathImage
			default:
				img = g.emptyImage
			}

			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*25), float64(y*25))
			screen.DrawImage(img, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000 // Подгонка под размер карты (40x40 точек по 25 пикселей)
}

func Start() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		GameMap: layers.GenerateMap(),
	}

	var err error
	game.startImage, _, err = ebitenutil.NewImageFromFile("assets/textures/start.png")
	if err != nil {
		log.Fatal(err)
	}
	game.exitImage, _, err = ebitenutil.NewImageFromFile("assets/textures/exit.png")
	if err != nil {
		log.Fatal(err)
	}
	game.pathImage, _, err = ebitenutil.NewImageFromFile("assets/textures/floor.png")
	if err != nil {
		log.Fatal(err)
	}
	game.emptyImage, _, err = ebitenutil.NewImageFromFile("assets/textures/empty.png")
	if err != nil {
		log.Fatal(err)
	}

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Map Generator")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
