package app

import (
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
	"log"
	"os"
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

	// Загружаем большие изображения врагов
	g.enemyLargeImages = make(map[string]*ebiten.Image)
	g.enemyLargeImages["Goblin"] = loadImage("assets/textures/goblin_large.png")
	g.enemyLargeImages["Golem"] = loadImage("assets/textures/golem_large.png")
	g.enemyLargeImages["Wolf"] = loadImage("assets/textures/wolf_large.png")

	// Загрузка больших изображений классов
	g.characterImages[WarriorClass] = loadImage("assets/textures/warrior.png")
	g.characterImages[MageClass] = loadImage("assets/textures/mage.png")
	g.characterImages[ArcherClass] = loadImage("assets/textures/archer.png")

	// Загрузка мини-иконок для карты
	g.classImages[WarriorClass] = loadImage("assets/textures/warrior_mini.png")
	g.classImages[MageClass] = loadImage("assets/textures/mage_mini.png")
	g.classImages[ArcherClass] = loadImage("assets/textures/archer_mini.png")

	// Загрузка фоновых изображений
	g.backgroundImage = loadImage("assets/textures/backmenu.png")
	g.combatBackgroundImage = loadImage("assets/textures/backcombat.png")
}

// Новая функция для загрузки шрифта
func LoadFont(path string, size float64) font.Face {
	// Читаем файл шрифта
	data, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("Failed to read font file %s: %v", path, err)
	}
	log.Printf("Successfully read font file %s, size: %d bytes", path, len(data))

	// Проверяем, не пустой ли файл
	if len(data) == 0 {
		log.Fatalf("Font file %s is empty", path)
	}

	// Парсим TTF-шрифт
	ttfFont, err := opentype.Parse(data)
	if err != nil {
		log.Fatalf("Failed to parse font %s: %v", path, err)
	}
	log.Printf("Successfully parsed font %s", path)

	// Создаём шрифт нужного размера
	face, err := opentype.NewFace(ttfFont, &opentype.FaceOptions{
		Size:    size,
		DPI:     72,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatalf("Failed to create font face for %s: %v", path, err)
	}
	log.Printf("Successfully created font face for %s", path)

	return face

}
