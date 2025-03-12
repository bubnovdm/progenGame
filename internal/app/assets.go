package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"log"
)

func loadImage(path string) *ebiten.Image {
	img, _, err := ebitenutil.NewImageFromFile(path)
	if err != nil {
		log.Fatal(err)
	}
	return img
}

func loadAssets(g *Game) {
	// Инициализация текстур
	g.textures[BackgroundSymbol] = loadImage("assets/textures/empty.png")
	g.textures[PathSymbol] = loadImage("assets/textures/floor.png")
	g.textures[StartSymbol] = loadImage("assets/textures/start.png")
	g.textures[ExitSymbol] = loadImage("assets/textures/exit.png")
	g.textures[WallSymbol] = loadImage("assets/textures/wall.png")
	g.enemyImage = loadImage("assets/textures/enemy.png")

	// Загрузка больших изображений классов (для CharacterSheet)
	g.characterImages[WarriorClass] = loadImage("assets/textures/warrior.png")
	g.characterImages[MageClass] = loadImage("assets/textures/mage.png")
	g.characterImages[ArcherClass] = loadImage("assets/textures/archer.png")

	// Загрузка мини-иконок для карты
	g.classImages[WarriorClass] = loadImage("assets/textures/warrior_mini.png")
	g.classImages[MageClass] = loadImage("assets/textures/mage_mini.png")
	g.classImages[ArcherClass] = loadImage("assets/textures/archer_mini.png")

	// Установим начальное изображение игрока для карты
	g.playerImage = g.classImages[WarriorClass]

	// Загрузка фонового изображения
	g.backgroundImage = loadImage("assets/textures/backmenu.png")
}
