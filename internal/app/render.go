package app

import (
	"fmt"
	"github.com/bubnovdm/progenGame/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
)

func (g *Game) drawCharacterSheet(screen *ebiten.Image) {
	screen.Fill(color.Black)

	// Отрисовка фонового изображения
	if g.backgroundImage != nil {
		geom := ebiten.GeoM{}
		geom.Scale(1000.0/float64(g.backgroundImage.Bounds().Dx()), 1000.0/float64(g.backgroundImage.Bounds().Dy()))
		screen.DrawImage(g.backgroundImage, &ebiten.DrawImageOptions{GeoM: geom})
	}

	// Фон для заголовка
	titleBg := ebiten.NewImage(200, 50)
	titleBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom := ebiten.GeoM{}
	geom.Translate(400, 30)
	screen.DrawImage(titleBg, &ebiten.DrawImageOptions{GeoM: geom})
	ebitenutil.DebugPrintAt(screen, "Select character", 450, 50)

	// Фон для характеристик
	statsBg := ebiten.NewImage(300, 250)
	statsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom = ebiten.GeoM{}
	geom.Translate(50, 100)
	screen.DrawImage(statsBg, &ebiten.DrawImageOptions{GeoM: geom})

	// Характеристики (слева)
	ebitenutil.DebugPrintAt(screen, "Stats", 60, 110)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Class: %s", g.Player.Class), 60, 140)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Str: %d", g.Player.Strength), 60, 170)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Agi: %d", g.Player.Agility), 60, 200)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Int: %d", g.Player.Intelligence), 60, 230)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("pDef: %d", g.Player.PhDefense), 60, 260)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mDef: %d", g.Player.MgDefense), 60, 290)

	// Фон для талантов
	talentsBg := ebiten.NewImage(300, 150)
	talentsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom = ebiten.GeoM{}
	geom.Translate(50, 350)
	screen.DrawImage(talentsBg, &ebiten.DrawImageOptions{GeoM: geom})

	// Таланты (слева, ниже характеристик)
	ebitenutil.DebugPrintAt(screen, "Talents", 60, 360)
	for i, skill := range g.Player.Skills {
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s (Lv. %d, CD: %.1f)", skill.Name, skill.Level, skill.Cooldown), 60, 390+i*30)
	}

	// Фон для инвентаря
	inventoryBg := ebiten.NewImage(200, 250)
	inventoryBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom = ebiten.GeoM{}
	geom.Translate(700, 100)
	screen.DrawImage(inventoryBg, &ebiten.DrawImageOptions{GeoM: geom})

	// Инвентарь (справа)
	ebitenutil.DebugPrintAt(screen, "Inventory", 710, 110)
	for i := 0; i < 3; i++ {
		for j := 0; j < 3; j++ {
			slot := ebiten.NewImage(50, 50)
			slot.Fill(color.RGBA{R: 100, G: 100, B: 100, A: 255})
			geom = ebiten.GeoM{}
			geom.Translate(float64(710+j*60), float64(130+i*60))
			screen.DrawImage(slot, &ebiten.DrawImageOptions{GeoM: geom})
			if i*3+j < len(g.Player.Inventory) {
				item := g.Player.Inventory[i*3+j]
				ebitenutil.DebugPrintAt(screen, item.Name, 710+j*60+5, 130+i*60+20)
			}
		}
	}

	// Изображение персонажа в центре (используем characterImages)
	if img, ok := g.characterImages[g.Player.Class]; ok {
		geom = ebiten.GeoM{}
		geom.Translate(400, 300) // Центр изображения (500 - 100, 400 - 100)
		//geom.Scale(0.75, 0.75)   // Масштабирование до 150x150 (если исходное 200x200)
		op := &ebiten.DrawImageOptions{GeoM: geom}
		screen.DrawImage(img, op)
	}

	// Кнопки переключения класса и "Play"
	for _, button := range g.getCharacterSheetButtons() {
		buttonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
		mx, my := ebiten.CursorPosition()
		if mx >= button.X && mx <= button.X+button.Width &&
			my >= button.Y && my <= button.Y+button.Height {
			buttonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		}
		buttonImage := ebiten.NewImage(button.Width, button.Height)
		buttonImage.Fill(buttonColor)
		geom = ebiten.GeoM{}
		geom.Translate(float64(button.X), float64(button.Y))
		screen.DrawImage(buttonImage, &ebiten.DrawImageOptions{GeoM: geom})
		ebitenutil.DebugPrintAt(screen, button.Label, button.X+10, button.Y+15)
	}
}

func (g *Game) drawDungeon(screen *ebiten.Image) {
	screen.Fill(color.Black)

	visibleRadius := 7
	playerX := float64(g.Player.X * 25)
	playerY := float64(g.Player.Y * 25)

	minX := utils.Max(0, int(playerX/25)-visibleRadius) * 25
	maxX := utils.Min(MapSize-1, int(playerX/25)+visibleRadius) * 25
	minY := utils.Max(0, int(playerY/25)-visibleRadius) * 25
	maxY := utils.Min(MapSize-1, int(playerY/25)+visibleRadius) * 25

	// Отрисовка фона с туманом войны
	for y := minY; y <= maxY; y += 25 {
		for x := minX; x <= maxX; x += 25 {
			cellX := int(x / 25)
			cellY := int(y / 25)
			dx := float64(x) - playerX
			dy := float64(y) - playerY
			distance := math.Sqrt(dx*dx+dy*dy) / 25

			if distance <= float64(visibleRadius) {
				cell := g.GameMap.Background[cellY][cellX]
				if img, ok := g.textures[cell]; ok {
					op := &ebiten.DrawImageOptions{}
					alpha := 1.0 - (distance / float64(visibleRadius))
					if alpha < 0.3 {
						alpha = 0.3
					}
					op.ColorScale.SetA(float32(alpha))
					geom := ebiten.GeoM{}
					geom.Translate(float64(x), float64(y))
					op.GeoM = geom
					screen.DrawImage(img, op)
				}
			}
		}
	}

	// Отрисовка пола без затухания внутри радиуса
	for y := minY; y <= maxY; y += 25 {
		for x := minX; x <= maxX; x += 25 {
			cellX := int(x / 25)
			cellY := int(y / 25)
			dx := float64(x) - playerX
			dy := float64(y) - playerY
			distance := math.Sqrt(dx*dx+dy*dy) / 25

			if distance <= float64(visibleRadius) {
				cell := g.GameMap.Floor[cellY][cellX]
				if cell == EmptySymbol {
					continue
				}
				if img, ok := g.textures[cell]; ok {
					op := &ebiten.DrawImageOptions{}
					op.ColorScale.SetA(1.0)
					geom := ebiten.GeoM{}
					geom.Translate(float64(x), float64(y))
					op.GeoM = geom
					screen.DrawImage(img, op)
				}
			}
		}
	}

	// Отрисовка объектов (стены) без затухания внутри радиуса
	for y := minY; y <= maxY; y += 25 {
		for x := minX; x <= maxX; x += 25 {
			cellX := int(x / 25)
			cellY := int(y / 25)
			dx := float64(x) - playerX
			dy := float64(y) - playerY
			distance := math.Sqrt(dx*dx+dy*dy) / 25

			if distance <= float64(visibleRadius) {
				cell := g.GameMap.Objects[cellY][cellX]
				if cell == EmptySymbol {
					continue
				}
				if img, ok := g.textures[cell]; ok {
					op := &ebiten.DrawImageOptions{}
					op.ColorScale.SetA(1.0)
					geom := ebiten.GeoM{}
					geom.Translate(float64(x), float64(y))
					op.GeoM = geom
					screen.DrawImage(img, op)
				}
			}
		}
	}

	// Отрисовка врагов
	for _, enemy := range g.Enemies {
		dx := float64(enemy.X*25) - playerX
		dy := float64(enemy.Y*25) - playerY
		distance := math.Sqrt(dx*dx+dy*dy) / 25

		if distance <= float64(visibleRadius) {
			op := &ebiten.DrawImageOptions{}
			geom := ebiten.GeoM{}
			geom.Translate(float64(enemy.X*25), float64(enemy.Y*25))
			op.GeoM = geom
			screen.DrawImage(g.enemyImage, op)
		}
	}

	// Отрисовка игрока
	if g.playerImage != nil {
		op := &ebiten.DrawImageOptions{}
		geom := ebiten.GeoM{}
		geom.Translate(playerX, playerY)
		op.GeoM = geom
		screen.DrawImage(g.playerImage, op)
	}

	// Отрисовка полосок HP и опыта
	const (
		barWidth  = 200
		barHeight = 20
		barX      = 10
		hpY       = 10
		expY      = 35 // Полоска опыта ниже HP
	)

	// Полоска HP
	hpRatio := float64(g.Player.HP) / float64(g.Player.MaxHP) // Используем g.Player.MaxHP
	hpBarWidth := int(float64(barWidth) * hpRatio)
	if hpBarWidth < 1 {
		hpBarWidth = 1 // Минимальная ширина 1 пиксель
	}

	hpBackground := ebiten.NewImage(barWidth, barHeight)
	hpBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom := ebiten.GeoM{}
	geom.Translate(float64(barX), float64(hpY))
	screen.DrawImage(hpBackground, &ebiten.DrawImageOptions{GeoM: geom})

	hpFill := ebiten.NewImage(hpBarWidth, barHeight)
	hpFill.Fill(color.RGBA{R: 40, G: 170, B: 40, A: 255})
	geom = ebiten.GeoM{}
	geom.Translate(float64(barX), float64(hpY))
	screen.DrawImage(hpFill, &ebiten.DrawImageOptions{GeoM: geom})

	// Полоска опыта
	const expPerLevel = 100 // Опыт для следующего уровня
	expRatio := float64(g.Player.Experience) / float64(expPerLevel)
	expBarWidth := int(float64(barWidth) * expRatio)
	if expBarWidth < 1 {
		expBarWidth = 1 // Минимальная ширина 1 пиксель
	}

	expBackground := ebiten.NewImage(barWidth, barHeight)
	expBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom = ebiten.GeoM{}
	geom.Translate(float64(barX), float64(expY))
	screen.DrawImage(expBackground, &ebiten.DrawImageOptions{GeoM: geom})

	expFill := ebiten.NewImage(expBarWidth, barHeight)
	expFill.Fill(color.RGBA{R: 0, G: 255, B: 215, A: 255}) // Голубой цвет для опыта
	geom = ebiten.GeoM{}
	geom.Translate(float64(barX), float64(expY))
	screen.DrawImage(expFill, &ebiten.DrawImageOptions{GeoM: geom})

	// Текст для HP и опыта
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("HP: %d/%d", g.Player.HP, g.Player.MaxHP), barX+5, hpY+3) // Используем g.Player.MaxHP
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Level: %d (Exp: %d/%d)", g.Player.Level, g.Player.Experience, expPerLevel), barX+5, expY+3)

	// Отображение уровня по центру (имеется в виду этаж, а не уровень персонажа)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Floor: %d", g.Level), 500, 15)
}

func (g *Game) drawMenu(screen *ebiten.Image) {
	// Отрисовка фонового изображения
	if g.backgroundImage != nil {
		geom := ebiten.GeoM{}
		geom.Scale(1000.0/float64(g.backgroundImage.Bounds().Dx()), 1000.0/float64(g.backgroundImage.Bounds().Dy()))
		screen.DrawImage(g.backgroundImage, &ebiten.DrawImageOptions{GeoM: geom})
	} else {
		screen.Fill(color.Black)
	}

	// Отрисовка кнопок
	for _, button := range g.getMenuButtons() {
		buttonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
		mx, my := ebiten.CursorPosition()
		if mx >= button.X && mx <= button.X+button.Width &&
			my >= button.Y && my <= button.Y+button.Height {
			buttonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		}
		buttonImage := ebiten.NewImage(button.Width, button.Height)
		buttonImage.Fill(buttonColor)
		geom := ebiten.GeoM{}
		geom.Translate(float64(button.X), float64(button.Y))
		screen.DrawImage(buttonImage, &ebiten.DrawImageOptions{GeoM: geom})
		ebitenutil.DebugPrintAt(screen, button.Label, button.X+50, button.Y+15)
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

func (g *Game) drawCombat(screen *ebiten.Image) {
	// Отрисовка фона
	if g.combatBackgroundImage != nil {
		screen.DrawImage(g.combatBackgroundImage, &ebiten.DrawImageOptions{})
	} else {
		// Запасной вариант, если фон не загрузился
		overlay := ebiten.NewImage(1000, 1000)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 150})
		screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
	}

	// Отрисовка персонажа (слева от центра)
	if img, ok := g.characterImages[g.Player.Class]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(250, 400) // X=500-250 (центр - 250)
		screen.DrawImage(img, op)
	}

	// Отрисовка врага (справа от центра)
	if g.CurrentEnemy != nil {
		// Выбор изображения врага по имени
		enemyImage, ok := g.enemyLargeImages[g.CurrentEnemy.Name]
		if !ok || enemyImage == nil {
			// Запасной вариант: используем изображение гоблина
			enemyImage, ok = g.enemyLargeImages["Goblin"]
			if !ok || enemyImage == nil {
				// Если даже гоблина нет, пропускаем отрисовку
				enemyImage = nil
			}
		}
		if enemyImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(550, 400) // X=500+50 (центр + 50)
			screen.DrawImage(enemyImage, op)
		}
	}

	// Полоски HP и опыта
	const barWidth, barHeight = 200, 20 // Фиксированная длина 200 пикселей

	// Полоска HP персонажа
	hpRatioPlayer := float64(g.Player.HP) / float64(g.Player.MaxHP)
	hpBarWidthPlayer := int(float64(barWidth) * hpRatioPlayer) // Заполнение в процентах
	if hpBarWidthPlayer < 1 {
		hpBarWidthPlayer = 1 // Минимальная ширина 1 пиксель
	}
	hpBackgroundPlayer := ebiten.NewImage(barWidth, barHeight)
	hpBackgroundPlayer.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geomPlayer := ebiten.GeoM{}
	geomPlayer.Translate(250, 350) // Над персонажем
	screen.DrawImage(hpBackgroundPlayer, &ebiten.DrawImageOptions{GeoM: geomPlayer})
	hpFillPlayer := ebiten.NewImage(hpBarWidthPlayer, barHeight)
	hpFillPlayer.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255}) // Зелёный для игрока
	geomPlayer = ebiten.GeoM{}
	geomPlayer.Translate(250, 350)
	screen.DrawImage(hpFillPlayer, &ebiten.DrawImageOptions{GeoM: geomPlayer})
	// Текст HP персонажа
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%d/%d", g.Player.HP, g.Player.MaxHP), 250, 330)

	// Полоска опыта персонажа
	const expPerLevel = 100 // Опыт для следующего уровня
	expRatio := float64(g.Player.Experience) / float64(expPerLevel)
	expBarWidth := int(float64(barWidth) * expRatio)
	if expBarWidth < 1 {
		expBarWidth = 1 // Минимальная ширина 1 пиксель
	}
	expBackground := ebiten.NewImage(barWidth, barHeight)
	expBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geomExp := ebiten.GeoM{}
	geomExp.Translate(250, 380) // Под полоской HP
	screen.DrawImage(expBackground, &ebiten.DrawImageOptions{GeoM: geomExp})
	expFill := ebiten.NewImage(expBarWidth, barHeight)
	expFill.Fill(color.RGBA{R: 0, G: 255, B: 215, A: 255}) // Желтый цвет для опыта
	geomExp = ebiten.GeoM{}
	geomExp.Translate(250, 380)
	screen.DrawImage(expFill, &ebiten.DrawImageOptions{GeoM: geomExp})
	// Текст уровня и опыта
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Level: %d (Exp: %d/%d)", g.Player.Level, g.Player.Experience, expPerLevel), 250, 383)

	// Полоска HP врага
	if g.CurrentEnemy != nil {
		hpRatio := float64(g.CurrentEnemy.HP) / float64(g.CurrentEnemy.MaxHP)
		hpBarWidth := int(float64(barWidth) * hpRatio) // Заполнение в процентах
		if hpBarWidth < 1 {
			hpBarWidth = 1 // Минимальная ширина 1 пиксель
		}
		hpBackground := ebiten.NewImage(barWidth, barHeight)
		hpBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom := ebiten.GeoM{}
		geom.Translate(550, 350) // Над врагом
		screen.DrawImage(hpBackground, &ebiten.DrawImageOptions{GeoM: geom})
		hpFill := ebiten.NewImage(hpBarWidth, barHeight)
		hpFill.Fill(color.RGBA{R: 200, G: 0, B: 0, A: 255}) // Красный для врага
		geom = ebiten.GeoM{}
		geom.Translate(550, 350)
		screen.DrawImage(hpFill, &ebiten.DrawImageOptions{GeoM: geom})
		// Текст HP и имя врага
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s\n%d/%d", g.CurrentEnemy.Name, g.CurrentEnemy.HP, g.CurrentEnemy.MaxHP), 550, 310)
	}

	// Кнопки способностей (пока просто текст)
	ability1Config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), "1")
	ability1Text := "1: Basic Attack" // Значение по умолчанию
	if ability1Config != nil {
		ability1Text = fmt.Sprintf("1: %s", ability1Config.Name)
	}

	ability2Config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), "2")
	ability2Text := "2: Skill"
	if ability2Config != nil {
		ability2Text = fmt.Sprintf("2: %s", ability2Config.Name)
	}

	// Подложка под названия способностей для отслеживания кд
	if g.AbilityCooldowns["1"] > 0 {
		ebitenutil.DrawRect(screen, 50, 600, 150, 20, color.RGBA{255, 0, 0, 128}) // Красный фон
	} else {
		ebitenutil.DrawRect(screen, 50, 600, 150, 20, color.RGBA{0, 255, 0, 128}) // Зеленый фон
	}

	if g.AbilityCooldowns["2"] > 0 {
		ebitenutil.DrawRect(screen, 50, 630, 150, 20, color.RGBA{255, 0, 0, 128}) // Красный фон
	} else {
		ebitenutil.DrawRect(screen, 50, 630, 150, 20, color.RGBA{0, 255, 0, 128}) // Зеленый фон
	}

	if g.AbilityCooldowns["3"] > 0 {
		ebitenutil.DrawRect(screen, 50, 660, 150, 20, color.RGBA{255, 0, 0, 128}) // Красный фон
	} else {
		ebitenutil.DrawRect(screen, 50, 660, 150, 20, color.RGBA{0, 255, 0, 128}) // Зеленый фон
	}

	ebitenutil.DebugPrintAt(screen, ability1Text, 50, 600)
	ebitenutil.DebugPrintAt(screen, ability2Text, 50, 630)
	ebitenutil.DebugPrintAt(screen, "3: Ult", 50, 660) // Пока оставим как есть

	// Лог боя
	yOffset := 700
	for i, log := range g.CombatLog {
		ebitenutil.DebugPrintAt(screen, log, 50, yOffset+20*int(i))
	}
	if len(g.CombatLog) > 5 {
		g.CombatLog = g.CombatLog[len(g.CombatLog)-5:]
	}
}

// Удаление врага по ID
func removeEnemy(enemies []Enemy, id string) []Enemy {
	for i, enemy := range enemies {
		if enemy.ID == id {
			return append(enemies[:i], enemies[i+1:]...)
		}
	}
	return enemies
}
