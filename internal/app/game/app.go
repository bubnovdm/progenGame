package game

//app.go

import (
	"github.com/bubnovdm/progenGame/internal/app/world"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
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

type PlayerPosition struct {
	X, Y int // Позиция игрока на карте
}

type Game struct {
	GameMap        world.GameMap
	PlayerPosition PlayerPosition
	moveDelay      int
	textures       map[rune]*ebiten.Image // Мапа для текстур
	playerImage    *ebiten.Image          // Отдельное поле для игрока
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
	screen.Fill(color.Black)

	visibleRadius := 7
	playerX := g.PlayerPosition.X
	playerY := g.PlayerPosition.Y

	// Ограничим область отрисовки
	minX := max(0, playerX-visibleRadius)
	maxX := min(world.MapSize-1, playerX+visibleRadius)
	minY := max(0, playerY-visibleRadius)
	maxY := min(world.MapSize-1, playerY+visibleRadius)

	// Отрисовка фона
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			cell := g.GameMap.Background[y][x]
			if img, ok := g.textures[cell]; ok {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*25), float64(y*25))
				screen.DrawImage(img, op)
			}
		}
	}

	// Отрисовка пола
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			cell := g.GameMap.Floor[y][x]
			if img, ok := g.textures[cell]; ok {
				op := &ebiten.DrawImageOptions{}
				op.GeoM.Translate(float64(x*25), float64(y*25))
				screen.DrawImage(img, op)
			}
		}
	}

	// Отрисовка игрока
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.PlayerPosition.X*25), float64(g.PlayerPosition.Y*25))
	screen.DrawImage(g.playerImage, op)
}

// Вспомогательные функции
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

// Вспомогательная функция для вычисления абсолютного значения
func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000 // Подгонка под размер карты (40x40 точек по 25 пикселей)
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func Start() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		GameMap:  world.GenerateMap(),
		textures: make(map[rune]*ebiten.Image), // Инициализация мапы
	}

	// Загрузка текстур в мапу
	game.textures[world.BackgroundSymbol] = loadImage("assets/textures/empty.png")
	game.textures[world.PathSymbol] = loadImage("assets/textures/floor.png")
	game.textures[world.StartSymbol] = loadImage("assets/textures/start.png")
	game.textures[world.ExitSymbol] = loadImage("assets/textures/exit.png")
	game.playerImage = loadImage("assets/textures/player.png")

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
