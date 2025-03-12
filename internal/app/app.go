package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
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
	GameMap            GameMap
	Player             Player
	Enemies            []Enemy
	moveDelay          int
	textures           map[rune]*ebiten.Image
	playerImage        *ebiten.Image
	enemyImage         *ebiten.Image
	Level              int
	State              GameState
	selectedClassIndex int                           // Индекс текущего выбранного класса
	classes            []PlayerClass                 // Список доступных классов
	classImages        map[PlayerClass]*ebiten.Image // Мини-иконки для карты
	characterImages    map[PlayerClass]*ebiten.Image // Большие изображения для CharacterSheet
	backgroundImage    *ebiten.Image                 // Фоновое изображение
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

	case CharacterSheet:
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			mx, my := ebiten.CursorPosition()
			for _, button := range g.getCharacterSheetButtons() {
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
			g.spawnEnemies()
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

func (g *Game) Draw(screen *ebiten.Image) {
	switch g.State {
	case Menu:
		g.drawMenu(screen)

	case CharacterSheet:
		g.drawCharacterSheet(screen)

	case InGameMenu:
		g.drawInGameMenu(screen)

	case Dungeon:
		g.drawDungeon(screen)
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return 1000, 1000
}

func Start() {
	rand.Seed(time.Now().UnixNano())

	game := &Game{
		Player:          NewPlayer(WarriorClass), // По умолчанию Воин
		textures:        make(map[rune]*ebiten.Image),
		Level:           1,
		State:           Menu,
		classes:         []PlayerClass{WarriorClass, MageClass, ArcherClass},
		classImages:     make(map[PlayerClass]*ebiten.Image),
		characterImages: make(map[PlayerClass]*ebiten.Image),
	}

	// Загрузка ресурсов
	loadAssets(game)

	ebiten.SetWindowSize(1000, 1000)
	ebiten.SetWindowTitle("Map Generator")
	if err := ebiten.RunGame(game); err != nil {
		panic(err)
	}
}
