package app

import (
	"fmt"
	"github.com/bubnovdm/progenGame/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"
)

type GameState int

const (
	Menu GameState = iota
	CharacterSheet
	Dungeon
	InGameMenu
)

type Game struct {
	GameMap     GameMap
	Player      Player
	moveDelay   int
	textures    map[rune]*ebiten.Image
	playerImage *ebiten.Image
	Level       int
	State       GameState
}

func (g *Game) Update() error {
	switch g.State {
	case Menu:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for _, button := range g.getMenuButtons() {
				if mx >= button.X && mx <= button.X+button.Width &&
					my >= button.Y && my <= button.Y+button.Height {
					button.Action(g)
				}
			}
		}

	case InGameMenu:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for _, button := range g.getInGameMenuButtons() {
				if mx >= button.X && mx <= button.X+button.Width &&
					my >= button.Y && my <= button.Y+button.Height {
					button.Action(g)
				}
			}
		}

	case Dungeon:
		if inpututil.IsKeyJustPressed(ebiten.KeyEscape) {
			g.State = InGameMenu
			return nil
		}

		if g.moveDelay > 0 {
			g.moveDelay--
			return nil
		}

		if g.GameMap.Floor[g.Player.Y][g.Player.X] == ExitSymbol {
			g.Level++
			var mapType MapType
			if g.Level%2 == 0 {
				mapType = OpenMap
			} else {
				mapType = DungeonMap
			}
			g.GameMap = GenerateMap(mapType)
			g.moveToStartPosition()
			g.moveDelay = 10
			return nil
		}

		newX, newY := g.Player.X, g.Player.Y

		if ebiten.IsKeyPressed(ebiten.KeyW) && g.Player.Y > 0 {
			newY--
		}
		if ebiten.IsKeyPressed(ebiten.KeyS) && g.Player.Y < MapSize-1 {
			newY++
		}
		if ebiten.IsKeyPressed(ebiten.KeyA) && g.Player.X > 0 {
			newX--
		}
		if ebiten.IsKeyPressed(ebiten.KeyD) && g.Player.X < MapSize-1 {
			newX++
		}

		if newX != g.Player.X || newY != g.Player.Y {
			if g.IsWalkable(newX, newY) {
				g.Player.X = newX
				g.Player.Y = newY
				g.moveDelay = 10
			}
		}
	}
	return nil
}

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

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case Menu:
		g.drawMenu(screen)

	case InGameMenu:
		g.drawInGameMenu(screen)

	case Dungeon:
		g.drawDungeon(screen)
	}
}

func (g *Game) drawDungeon(screen *ebiten.Image) {
	screen.Fill(color.Black)

	visibleRadius := 7
	playerX := g.Player.X
	playerY := g.Player.Y

	minX := utils.Max(0, playerX-visibleRadius)
	maxX := utils.Min(MapSize-1, playerX+visibleRadius)
	minY := utils.Max(0, playerY-visibleRadius)
	maxY := utils.Min(MapSize-1, playerY+visibleRadius)

	// Отрисовка фона с туманом войны
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := float64(x - playerX)
			dy := float64(y - playerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= float64(visibleRadius) {
				cell := g.GameMap.Background[y][x]
				if img, ok := g.textures[cell]; ok {
					op := &ebiten.DrawImageOptions{}
					alpha := 1.0 - (distance / float64(visibleRadius))
					if alpha < 0.3 {
						alpha = 0.3
					}
					op.ColorScale.SetA(float32(alpha))
					op.GeoM.Translate(float64(x*25), float64(y*25))
					screen.DrawImage(img, op)
				}
			}
		}
	}

	// Отрисовка пола без затухания внутри радиуса
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := float64(x - playerX)
			dy := float64(y - playerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= float64(visibleRadius) {
				cell := g.GameMap.Floor[y][x]
				if cell == EmptySymbol {
					continue
				}
				if img, ok := g.textures[cell]; ok {
					op := &ebiten.DrawImageOptions{}
					op.ColorScale.SetA(1.0)
					op.GeoM.Translate(float64(x*25), float64(y*25))
					screen.DrawImage(img, op)
				} else {
					log.Printf("No texture for Floor cell '%c' at (%d, %d)", cell, x, y)
				}
			}
		}
	}

	// Отрисовка объектов (стены) без затухания внутри радиуса
	for y := minY; y <= maxY; y++ {
		for x := minX; x <= maxX; x++ {
			dx := float64(x - playerX)
			dy := float64(y - playerY)
			distance := math.Sqrt(dx*dx + dy*dy)

			if distance <= float64(visibleRadius) {
				cell := g.GameMap.Objects[y][x]
				if cell == EmptySymbol {
					continue
				}
				if img, ok := g.textures[cell]; ok {
					op := &ebiten.DrawImageOptions{}
					op.ColorScale.SetA(1.0)
					op.GeoM.Translate(float64(x*25), float64(y*25))
					screen.DrawImage(img, op)
				} else {
					log.Printf("No texture for Objects cell '%c' at (%d, %d)", cell, x, y)
				}
			}
		}
	}

	// Отрисовка игрока
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(g.Player.X*25), float64(g.Player.Y*25))
	screen.DrawImage(g.playerImage, op)

	// Отображение уровня в верхнем левом углу
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Level: %d", g.Level))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func Start() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		Player: Player{
			ID:           "player1",
			Name:         "Hero",
			HP:           100,
			Mana:         50,
			Level:        1,
			Strength:     10,
			Agility:      10,
			Intelligence: 10,
		},
		textures: make(map[rune]*ebiten.Image),
		Level:    1,
		State:    Menu,
	}

	// Инициализация текстур
	game.textures[BackgroundSymbol] = loadImage("assets/textures/empty.png")
	game.textures[PathSymbol] = loadImage("assets/textures/floor.png")
	game.textures[StartSymbol] = loadImage("assets/textures/start.png")
	game.textures[ExitSymbol] = loadImage("assets/textures/exit.png")
	game.textures[WallSymbol] = loadImage("assets/textures/wall.png")
	game.playerImage = loadImage("assets/textures/player.png")

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Map Generator")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}
