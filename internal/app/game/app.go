package game

//app.go

import (
	"github.com/bubnovdm/progenGame/internal/app/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
	"math/rand"
	"time"
)

type GameState int

const (
	Menu           GameState = iota // 0
	CharacterSheet                  // 1
	Dungeon                         // 2
)

type Layer [world.MapSize][world.MapSize]rune

type PlayerPosition struct {
	X, Y int // Позиция игрока на карте
}

type GameMap struct {
	Background Layer // Фон (например, трава или декоративные элементы)
	Floor      Layer // Пол (пути, по которым можно ходить)
	Objects    Layer // Объекты окружения (стены, сундуки, двери и т.д.)
}

type Game struct {
	GameMap        GameMap
	PlayerPosition PlayerPosition
	moveDelay      int
	startImage     *ebiten.Image
	exitImage      *ebiten.Image
	pathImage      *ebiten.Image
	emptyImage     *ebiten.Image
	playerImage    *ebiten.Image
}

func (g *Game) Update() error {
	if g.moveDelay > 0 {
		g.moveDelay--
		return nil
	}

	newX, newY := g.PlayerPosition.X, g.PlayerPosition.Y

	if ebiten.IsKeyPressed(ebiten.KeyW) && g.PlayerPosition.Y > 0 {
		newY--
	}
	if ebiten.IsKeyPressed(ebiten.KeyS) && g.PlayerPosition.Y < world.MapSize-1 {
		newY++
	}
	if ebiten.IsKeyPressed(ebiten.KeyA) && g.PlayerPosition.X > 0 {
		newX--
	}
	if ebiten.IsKeyPressed(ebiten.KeyD) && g.PlayerPosition.X < world.MapSize-1 {
		newX++
	}

	// Проверяем, можно ли двигаться на новую позицию
	if newX != g.PlayerPosition.X || newY != g.PlayerPosition.Y {
		if g.IsWalkable(newX, newY) {
			g.PlayerPosition.X = newX
			g.PlayerPosition.Y = newY
			g.moveDelay = 10
		}
	}
	return nil
}

func (g *Game) IsWalkable(x, y int) bool {
	return g.GameMap.Floor[y][x] == world.PathSymbol ||
		g.GameMap.Floor[y][x] == world.StartSymbol ||
		g.GameMap.Floor[y][x] == world.ExitSymbol
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Отрисовка фона
	for y, row := range g.GameMap.Background {
		for x, cell := range row {
			var img *ebiten.Image
			if cell == 'G' {
				img = g.emptyImage
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*25), float64(y*25))
			screen.DrawImage(img, op)
		}
	}

	// Отрисовка пола
	for y, row := range g.GameMap.Floor {
		for x, cell := range row {
			var img *ebiten.Image
			switch cell {
			case world.StartSymbol:
				img = g.startImage
			case world.ExitSymbol:
				img = g.exitImage
			case world.PathSymbol:
				img = g.pathImage
			default:
				continue
			}
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(float64(x*25), float64(y*25))
			screen.DrawImage(img, op)
		}
	}

	// Отрисовка игрока
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.PlayerPosition.X*25), float64(g.PlayerPosition.Y*25))
	screen.DrawImage(g.playerImage, op)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000 // Подгонка под размер карты (40x40 точек по 25 пикселей)
}

func Start() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		GameMap: world.GenerateMap(),
	}

	// Загрузка текстур
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
	game.playerImage, _, err = ebitenutil.NewImageFromFile("assets/textures/player.png")
	if err != nil {
		log.Fatal(err)
	}

	// Устанавливаем игрока на начальную позицию
	for y := 0; y < world.MapSize; y++ {
		for x := 0; x < world.MapSize; x++ {
			if game.GameMap.Floor[y][x] == world.StartSymbol {
				game.PlayerPosition.X = x
				game.PlayerPosition.Y = y
				break
			}
		}
	}

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Map Generator")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
