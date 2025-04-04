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

	// Бафы (под характеристиками, вместо талантов)
	talentsBg := ebiten.NewImage(300, 150)
	talentsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom = ebiten.GeoM{}
	geom.Translate(50, 350)
	screen.DrawImage(talentsBg, &ebiten.DrawImageOptions{GeoM: geom})

	// Бафы вместо талантов
	ebitenutil.DebugPrintAt(screen, "Buffs:", 60, 360) // Заголовок "Buffs" вместо "Talents"
	for i, buff := range g.AvailableBuffs {
		ebitenutil.DebugPrintAt(screen, buff.Name(), 60, 390+i*30) // Каждый баф с отступом 30 пикселей
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

	// Выпадающий список для выбора этажа (под инвентарем)
	const (
		floorButtonWidth  = 200
		floorButtonHeight = 30
		floorButtonX      = 700
		floorButtonY      = 360 // Под инвентарем (100 + 250 + 10 отступ)
	)

	// Кнопка для открытия/закрытия выпадающего списка
	floorButton := ebiten.NewImage(floorButtonWidth, floorButtonHeight)
	floorButtonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
	mx, my := ebiten.CursorPosition()
	if mx >= floorButtonX && mx <= floorButtonX+floorButtonWidth &&
		my >= floorButtonY && my <= floorButtonY+floorButtonHeight {
		floorButtonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
	}
	floorButton.Fill(floorButtonColor)
	geom = ebiten.GeoM{}
	geom.Translate(float64(floorButtonX), float64(floorButtonY))
	screen.DrawImage(floorButton, &ebiten.DrawImageOptions{GeoM: geom})
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Floor: %d", g.SelectedFloor), floorButtonX+10, floorButtonY+5)

	// Отрисовка выпадающего списка, если он открыт
	if g.FloorSelectorOpen {
		for i := 1; i <= g.MaxFloor; i++ {
			option := ebiten.NewImage(floorButtonWidth, floorButtonHeight)
			optionColor := color.RGBA{R: 80, G: 80, B: 80, A: 255}
			if mx >= floorButtonX && mx <= floorButtonX+floorButtonWidth &&
				my >= floorButtonY+floorButtonHeight*i && my <= floorButtonY+floorButtonHeight*(i+1) {
				optionColor = color.RGBA{R: 120, G: 120, B: 120, A: 255}
			}
			option.Fill(optionColor)
			geom = ebiten.GeoM{}
			geom.Translate(float64(floorButtonX), float64(floorButtonY+floorButtonHeight*i))
			screen.DrawImage(option, &ebiten.DrawImageOptions{GeoM: geom})
			ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Floor %d", i), floorButtonX+10, floorButtonY+floorButtonHeight*i+5)
		}
	}

	// Изображение персонажа в центре
	if img, ok := g.characterImages[g.Player.Class]; ok {
		geom = ebiten.GeoM{}
		geom.Translate(400, 300)
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
				if img, ok := g.textures[rune(cell)]; ok {
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
				if img, ok := g.textures[rune(cell)]; ok {
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
				if img, ok := g.textures[rune(cell)]; ok {
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
	if img, ok := g.classImages[g.Player.Class]; ok {
		op := &ebiten.DrawImageOptions{}
		geom := ebiten.GeoM{}
		geom.Translate(playerX, playerY)
		op.GeoM = geom
		screen.DrawImage(img, op)
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Floor: %d", g.CurrentFloor), 500, 15)
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
		buttonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255} // Серый цвет по умолчанию
		mx, my := ebiten.CursorPosition()

		// Если кнопка не отключена, проверяем наведение курсора
		if !button.Disabled {
			if mx >= button.X && mx <= button.X+button.Width &&
				my >= button.Y && my <= button.Y+button.Height {
				buttonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255} // Светлее при наведении
			}
		} else {
			buttonColor = color.RGBA{R: 50, G: 50, B: 50, A: 255} // Темнее, если кнопка отключена
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

	// Характеристики (слева)
	ebitenutil.DebugPrintAt(screen, "Stats", 60, 110)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Class: %s", g.Player.Class), 60, 140)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Str: %d", g.Player.Strength), 60, 170)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Agi: %d", g.Player.Agility), 60, 200)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Int: %d", g.Player.Intelligence), 60, 230)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("pDef: %d", g.Player.PhDefense), 60, 260)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("mDef: %d", g.Player.MgDefense), 60, 290)

	// Бафы (справа от характеристик)
	buffsBg := ebiten.NewImage(300, 150) // Фон для бафов, можно будет попробовать добавить прокрутку, если бафов будет много
	buffsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geom := ebiten.GeoM{}
	geom.Translate(700, 110) // Позиция справа, на уровне характеристик
	screen.DrawImage(buffsBg, &ebiten.DrawImageOptions{GeoM: geom})

	ebitenutil.DebugPrintAt(screen, "Buffs:", 710, 120) // Заголовок для бафов
	for i, buff := range g.AvailableBuffs {
		ebitenutil.DebugPrintAt(screen, buff.Name(), 710, 150+i*30) // Каждый баф с отступом 30 пикселей
	}

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
	expFill.Fill(color.RGBA{R: 0, G: 255, B: 215, A: 255}) // Голубой цвет для опыта
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

		// Полоски для DoT-эффектов
		const dotBarHeight = 10        // Высота полоски DoT (меньше, чем у HP)
		yOffset := 350 + barHeight + 5 // Начальная позиция Y для первой полоски DoT (сразу под HP с отступом 5 пикселей)

		for _, effect := range g.CurrentEnemy.ActiveEffects {
			// Проверяем, является ли эффект DoT (по имени или типу)
			if dotEffect, ok := effect.(*DotEffect); ok {
				// Рассчитываем прогресс эффекта
				dotRatio := dotEffect.TimeRemaining / dotEffect.Duration
				if dotRatio < 0 {
					dotRatio = 0
				}
				dotBarWidth := int(float64(barWidth) * dotRatio)
				if dotBarWidth < 1 {
					dotBarWidth = 1 // Минимальная ширина 1 пиксель
				}

				// Фон полоски DoT
				dotBackground := ebiten.NewImage(barWidth, dotBarHeight)
				dotBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
				geom := ebiten.GeoM{}
				geom.Translate(550, float64(yOffset))
				screen.DrawImage(dotBackground, &ebiten.DrawImageOptions{GeoM: geom})

				// Заполнение полоски DoT (например, фиолетовый цвет для DoT)
				dotFill := ebiten.NewImage(dotBarWidth, dotBarHeight)
				dotColor := color.RGBA{R: 128, G: 0, B: 128, A: 255} // Фиолетовый по умолчанию
				switch dotEffect.Name {
				case "Poison":
					dotColor = color.RGBA{R: 0, G: 128, B: 0, A: 255} // Зеленый для яда
				case "Ignite":
					dotColor = color.RGBA{R: 255, G: 69, B: 0, A: 255} // Оранжевый для горения
				case "Bleed":
					dotColor = color.RGBA{R: 200, G: 0, B: 0, A: 255} // Красный для кровотечения
				}
				dotFill.Fill(dotColor)
				geom = ebiten.GeoM{}
				geom.Translate(550, float64(yOffset))
				screen.DrawImage(dotFill, &ebiten.DrawImageOptions{GeoM: geom})

				// Текст с названием эффекта и оставшимся временем
				ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: %.1fs", dotEffect.Name, dotEffect.TimeRemaining), 550, yOffset+2)

				// Сдвигаем Y для следующей полоски
				yOffset += dotBarHeight + 5 // Отступ между полосками
			}
		}
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

//Попытался сделать кастомный шрифт - пока не вышло
/*
	func (g *Game) drawCharacterSheet(screen *ebiten.Image) {
		screen.Fill(color.Black)

		// Отрисовка фонового изображения
		if g.backgroundImage != nil {
			geom := ebiten.GeoM{}
			geom.Scale(1000.0/float64(g.backgroundImage.Bounds().Dx()), 1000.0/float64(g.backgroundImage.Bounds().Dy()))
			screen.DrawImage(g.backgroundImage, &ebiten.DrawImageOptions{GeoM: geom})
		}

		// Цвет текста и тени
		textColor := color.White
		shadowColor := color.RGBA{R: 0, G: 0, B: 0, A: 128} // Полупрозрачная чёрная тень

		// Фон для заголовка
		titleBg := ebiten.NewImage(200, 50)
		titleBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom := ebiten.GeoM{}
		geom.Translate(400, 30)
		screen.DrawImage(titleBg, &ebiten.DrawImageOptions{GeoM: geom})
		text.Draw(screen, "Select character", g.Font, 450+2, 50+2, shadowColor) // Тень
		text.Draw(screen, "Select character", g.Font, 450, 50, textColor)       // Основной текст

		// Фон для характеристик
		statsBg := ebiten.NewImage(300, 250)
		statsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom = ebiten.GeoM{}
		geom.Translate(50, 100)
		screen.DrawImage(statsBg, &ebiten.DrawImageOptions{GeoM: geom})

		// Характеристики (слева)
		text.Draw(screen, "Stats", g.Font, 60+2, 110+2, shadowColor)
		text.Draw(screen, "Stats", g.Font, 60, 110, textColor)
		text.Draw(screen, fmt.Sprintf("Class: %s", g.Player.Class), g.Font, 60+2, 140+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Class: %s", g.Player.Class), g.Font, 60, 140, textColor)
		text.Draw(screen, fmt.Sprintf("Str: %d", g.Player.Strength), g.Font, 60+2, 170+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Str: %d", g.Player.Strength), g.Font, 60, 170, textColor)
		text.Draw(screen, fmt.Sprintf("Agi: %d", g.Player.Agility), g.Font, 60+2, 200+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Agi: %d", g.Player.Agility), g.Font, 60, 200, textColor)
		text.Draw(screen, fmt.Sprintf("Int: %d", g.Player.Intelligence), g.Font, 60+2, 230+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Int: %d", g.Player.Intelligence), g.Font, 60, 230, textColor)
		text.Draw(screen, fmt.Sprintf("pDef: %d", g.Player.PhDefense), g.Font, 60+2, 260+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("pDef: %d", g.Player.PhDefense), g.Font, 60, 260, textColor)
		text.Draw(screen, fmt.Sprintf("mDef: %d", g.Player.MgDefense), g.Font, 60+2, 290+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("mDef: %d", g.Player.MgDefense), g.Font, 60, 290, textColor)

		// Бафы (под характеристиками, вместо талантов)
		talentsBg := ebiten.NewImage(300, 150)
		talentsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom = ebiten.GeoM{}
		geom.Translate(50, 350)
		screen.DrawImage(talentsBg, &ebiten.DrawImageOptions{GeoM: geom})

		// Бафы вместо талантов
		text.Draw(screen, "Buffs", g.Font, 60+2, 360+2, shadowColor)
		text.Draw(screen, "Buffs", g.Font, 60, 360, textColor)
		for i, buff := range g.AvailableBuffs {
			yPos := 390 + i*30
			text.Draw(screen, buff.Name(), g.Font, 60+2, yPos+2, shadowColor)
			text.Draw(screen, buff.Name(), g.Font, 60, yPos, textColor)
		}

		// Фон для инвентаря
		inventoryBg := ebiten.NewImage(200, 250)
		inventoryBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom = ebiten.GeoM{}
		geom.Translate(700, 100)
		screen.DrawImage(inventoryBg, &ebiten.DrawImageOptions{GeoM: geom})

		// Инвентарь (справа)
		text.Draw(screen, "Inventory", g.Font, 710+2, 110+2, shadowColor)
		text.Draw(screen, "Inventory", g.Font, 710, 110, textColor)
		for i := 0; i < 3; i++ {
			for j := 0; j < 3; j++ {
				slot := ebiten.NewImage(50, 50)
				slot.Fill(color.RGBA{R: 100, G: 100, B: 100, A: 255})
				geom = ebiten.GeoM{}
				geom.Translate(float64(710+j*60), float64(130+i*60))
				screen.DrawImage(slot, &ebiten.DrawImageOptions{GeoM: geom})
				if i*3+j < len(g.Player.Inventory) {
					item := g.Player.Inventory[i*3+j]
					xPos, yPos := 710+j*60+5, 130+i*60+20
					text.Draw(screen, item.Name, g.Font, xPos+2, yPos+2, shadowColor)
					text.Draw(screen, item.Name, g.Font, xPos, yPos, textColor)
				}
			}
		}

		// Выпадающий список для выбора этажа (под инвентарем)
		const (
			floorButtonWidth  = 200
			floorButtonHeight = 30
			floorButtonX      = 700
			floorButtonY      = 360
		)

		// Кнопка для открытия/закрытия выпадающего списка
		floorButton := ebiten.NewImage(floorButtonWidth, floorButtonHeight)
		floorButtonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
		mx, my := ebiten.CursorPosition()
		if mx >= floorButtonX && mx <= floorButtonX+floorButtonWidth &&
			my >= floorButtonY && my <= floorButtonY+floorButtonHeight {
			floorButtonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
		}
		floorButton.Fill(floorButtonColor)
		geom = ebiten.GeoM{}
		geom.Translate(float64(floorButtonX), float64(floorButtonY))
		screen.DrawImage(floorButton, &ebiten.DrawImageOptions{GeoM: geom})
		floorText := fmt.Sprintf("Floor: %d", g.SelectedFloor)
		text.Draw(screen, floorText, g.Font, floorButtonX+10+2, floorButtonY+5+2, shadowColor)
		text.Draw(screen, floorText, g.Font, floorButtonX+10, floorButtonY+5, textColor)

		// Отрисовка выпадающего списка, если он открыт
		if g.FloorSelectorOpen {
			for i := 1; i <= g.MaxFloor; i++ {
				option := ebiten.NewImage(floorButtonWidth, floorButtonHeight)
				optionColor := color.RGBA{R: 80, G: 80, B: 80, A: 255}
				if mx >= floorButtonX && mx <= floorButtonX+floorButtonWidth &&
					my >= floorButtonY+floorButtonHeight*i && my <= floorButtonY+floorButtonHeight*(i+1) {
					optionColor = color.RGBA{R: 120, G: 120, B: 120, A: 255}
				}
				option.Fill(optionColor)
				geom = ebiten.GeoM{}
				geom.Translate(float64(floorButtonX), float64(floorButtonY+floorButtonHeight*i))
				screen.DrawImage(option, &ebiten.DrawImageOptions{GeoM: geom})
				floorOptionText := fmt.Sprintf("Floor %d", i)
				xPos, yPos := floorButtonX+10, floorButtonY+floorButtonHeight*i+5
				text.Draw(screen, floorOptionText, g.Font, xPos+2, yPos+2, shadowColor)
				text.Draw(screen, floorOptionText, g.Font, xPos, yPos, textColor)
			}
		}

		// Изображение персонажа в центре
		if img, ok := g.characterImages[g.Player.Class]; ok {
			geom = ebiten.GeoM{}
			geom.Translate(400, 300)
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
			xPos, yPos := button.X+10, button.Y+15
			text.Draw(screen, button.Label, g.Font, xPos+2, yPos+2, shadowColor)
			text.Draw(screen, button.Label, g.Font, xPos, yPos, textColor)
		}
	}
*/
/*
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
					if img, ok := g.textures[rune(cell)]; ok {
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
					if img, ok := g.textures[rune(cell)]; ok {
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
					if img, ok := g.textures[rune(cell)]; ok {
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
		if img, ok := g.classImages[g.Player.Class]; ok {
			op := &ebiten.DrawImageOptions{}
			geom := ebiten.GeoM{}
			geom.Translate(playerX, playerY)
			op.GeoM = geom
			screen.DrawImage(img, op)
		}

		// Отрисовка полосок HP и опыта
		const (
			barWidth  = 200
			barHeight = 20
			barX      = 10
			hpY       = 10
			expY      = 35
		)

		// Цвет текста и тени
		textColor := color.White
		shadowColor := color.RGBA{R: 0, G: 0, B: 0, A: 128}

		// Полоска HP
		hpRatio := float64(g.Player.HP) / float64(g.Player.MaxHP)
		hpBarWidth := int(float64(barWidth) * hpRatio)
		if hpBarWidth < 1 {
			hpBarWidth = 1
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
		const expPerLevel = 100
		expRatio := float64(g.Player.Experience) / float64(expPerLevel)
		expBarWidth := int(float64(barWidth) * expRatio)
		if expBarWidth < 1 {
			expBarWidth = 1
		}

		expBackground := ebiten.NewImage(barWidth, barHeight)
		expBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom = ebiten.GeoM{}
		geom.Translate(float64(barX), float64(expY))
		screen.DrawImage(expBackground, &ebiten.DrawImageOptions{GeoM: geom})

		expFill := ebiten.NewImage(expBarWidth, barHeight)
		expFill.Fill(color.RGBA{R: 0, G: 255, B: 215, A: 255})
		geom = ebiten.GeoM{}
		geom.Translate(float64(barX), float64(expY))
		screen.DrawImage(expFill, &ebiten.DrawImageOptions{GeoM: geom})

		// Текст для HP и опыта
		hpText := fmt.Sprintf("HP: %d/%d", g.Player.HP, g.Player.MaxHP)
		expText := fmt.Sprintf("Level: %d (Exp: %d/%d)", g.Player.Level, g.Player.Experience, expPerLevel)
		text.Draw(screen, hpText, g.Font, barX+5+2, hpY+3+2, shadowColor)
		text.Draw(screen, hpText, g.Font, barX+5, hpY+3, textColor)
		text.Draw(screen, expText, g.Font, barX+5+2, expY+3+2, shadowColor)
		text.Draw(screen, expText, g.Font, barX+5, expY+3, textColor)

		// Отображение уровня по центру (имеется в виду этаж, а не уровень персонажа)
		floorText := fmt.Sprintf("Floor: %d", g.CurrentFloor)
		text.Draw(screen, floorText, g.Font, 500+2, 15+2, shadowColor)
		text.Draw(screen, floorText, g.Font, 500, 15, textColor)
	}
*/
/*
	func (g *Game) drawMenu(screen *ebiten.Image) {
		// Отрисовка фонового изображения
		if g.backgroundImage != nil {
			geom := ebiten.GeoM{}
			geom.Scale(1000.0/float64(g.backgroundImage.Bounds().Dx()), 1000.0/float64(g.backgroundImage.Bounds().Dy()))
			screen.DrawImage(g.backgroundImage, &ebiten.DrawImageOptions{GeoM: geom})
		} else {
			screen.Fill(color.Black)
		}

		// Цвет текста и тени
		textColor := color.White
		shadowColor := color.RGBA{R: 0, G: 0, B: 0, A: 128}

		// Отрисовка кнопок
		for _, button := range g.getMenuButtons() {
			buttonColor := color.RGBA{R: 100, G: 100, B: 100, A: 255}
			mx, my := ebiten.CursorPosition()

			if !button.Disabled {
				if mx >= button.X && mx <= button.X+button.Width &&
					my >= button.Y && my <= button.Y+button.Height {
					buttonColor = color.RGBA{R: 150, G: 150, B: 150, A: 255}
				}
			} else {
				buttonColor = color.RGBA{R: 50, G: 50, B: 50, A: 255}
			}

			buttonImage := ebiten.NewImage(button.Width, button.Height)
			buttonImage.Fill(buttonColor)
			geom := ebiten.GeoM{}
			geom.Translate(float64(button.X), float64(button.Y))
			screen.DrawImage(buttonImage, &ebiten.DrawImageOptions{GeoM: geom})
			xPos, yPos := button.X+50, button.Y+15
			text.Draw(screen, button.Label, g.Font, xPos+2, yPos+2, shadowColor)
			text.Draw(screen, button.Label, g.Font, xPos, yPos, textColor)
		}
	}
*/
/*
	func (g *Game) drawInGameMenu(screen *ebiten.Image) {
		// Рисуем карту на фоне
		g.drawDungeon(screen)

		// Полупрозрачный фон для меню
		overlay := ebiten.NewImage(1000, 1000)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 200})
		op := &ebiten.DrawImageOptions{}
		screen.DrawImage(overlay, op)

		// Цвет текста и тени
		textColor := color.White
		shadowColor := color.RGBA{R: 0, G: 0, B: 0, A: 128}

		// Заголовок
		text.Draw(screen, "Pause Menu", g.Font, 450+2, 300+2, shadowColor)
		text.Draw(screen, "Pause Menu", g.Font, 450, 300, textColor)

		// Характеристики (слева)
		text.Draw(screen, "Stats", g.Font, 60+2, 110+2, shadowColor)
		text.Draw(screen, "Stats", g.Font, 60, 110, textColor)
		text.Draw(screen, fmt.Sprintf("Class: %s", g.Player.Class), g.Font, 60+2, 140+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Class: %s", g.Player.Class), g.Font, 60, 140, textColor)
		text.Draw(screen, fmt.Sprintf("Str: %d", g.Player.Strength), g.Font, 60+2, 170+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Str: %d", g.Player.Strength), g.Font, 60, 170, textColor)
		text.Draw(screen, fmt.Sprintf("Agi: %d", g.Player.Agility), g.Font, 60+2, 200+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Agi: %d", g.Player.Agility), g.Font, 60, 200, textColor)
		text.Draw(screen, fmt.Sprintf("Int: %d", g.Player.Intelligence), g.Font, 60+2, 230+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("Int: %d", g.Player.Intelligence), g.Font, 60, 230, textColor)
		text.Draw(screen, fmt.Sprintf("pDef: %d", g.Player.PhDefense), g.Font, 60+2, 260+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("pDef: %d", g.Player.PhDefense), g.Font, 60, 260, textColor)
		text.Draw(screen, fmt.Sprintf("mDef: %d", g.Player.MgDefense), g.Font, 60+2, 290+2, shadowColor)
		text.Draw(screen, fmt.Sprintf("mDef: %d", g.Player.MgDefense), g.Font, 60, 290, textColor)

		// Бафы (справа от характеристик)
		buffsBg := ebiten.NewImage(300, 150)
		buffsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom := ebiten.GeoM{}
		geom.Translate(700, 110)
		screen.DrawImage(buffsBg, &ebiten.DrawImageOptions{GeoM: geom})

		text.Draw(screen, "Buffs", g.Font, 710+2, 120+2, shadowColor)
		text.Draw(screen, "Buffs", g.Font, 710, 120, textColor)
		for i, buff := range g.AvailableBuffs {
			yPos := 150 + i*30
			text.Draw(screen, buff.Name(), g.Font, 710+2, yPos+2, shadowColor)
			text.Draw(screen, buff.Name(), g.Font, 710, yPos, textColor)
		}

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
			xPos, yPos := button.X+20, button.Y+15
			text.Draw(screen, button.Label, g.Font, xPos+2, yPos+2, shadowColor)
			text.Draw(screen, button.Label, g.Font, xPos, yPos, textColor)
		}
	}
*/
/*
func (g *Game) drawCombat(screen *ebiten.Image) {
	// Отрисовка фона
	if g.combatBackgroundImage != nil {
		screen.DrawImage(g.combatBackgroundImage, &ebiten.DrawImageOptions{})
	} else {
		overlay := ebiten.NewImage(1000, 1000)
		overlay.Fill(color.RGBA{R: 0, G: 0, B: 0, A: 150})
		screen.DrawImage(overlay, &ebiten.DrawImageOptions{})
	}

	// Отрисовка персонажа (слева от центра)
	if img, ok := g.characterImages[g.Player.Class]; ok {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(250, 400)
		screen.DrawImage(img, op)
	}

	// Отрисовка врага (справа от центра)
	if g.CurrentEnemy != nil {
		enemyImage, ok := g.enemyLargeImages[g.CurrentEnemy.Name]
		if !ok || enemyImage == nil {
			enemyImage, ok = g.enemyLargeImages["Goblin"]
			if !ok || enemyImage == nil {
				enemyImage = nil
			}
		}
		if enemyImage != nil {
			op := &ebiten.DrawImageOptions{}
			op.GeoM.Translate(550, 400)
			screen.DrawImage(enemyImage, op)
		}
	}

	// Полоски HP и опыта
	const barWidth, barHeight = 200, 20

	// Цвет текста и тени
	textColor := color.White
	shadowColor := color.RGBA{R: 0, G: 0, B: 0, A: 128}

	// Полоска HP персонажа
	hpRatioPlayer := 0.0
	if g.Player.MaxHP > 0 { // Проверка на деление на ноль
		hpRatioPlayer = float64(g.Player.HP) / float64(g.Player.MaxHP)
	}
	hpBarWidthPlayer := int(float64(barWidth) * hpRatioPlayer)
	if hpBarWidthPlayer < 1 {
		hpBarWidthPlayer = 1
	}
	hpBackgroundPlayer := ebiten.NewImage(barWidth, barHeight)
	hpBackgroundPlayer.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geomPlayer := ebiten.GeoM{}
	geomPlayer.Translate(250, 350)
	screen.DrawImage(hpBackgroundPlayer, &ebiten.DrawImageOptions{GeoM: geomPlayer})
	hpFillPlayer := ebiten.NewImage(hpBarWidthPlayer, barHeight)
	hpFillPlayer.Fill(color.RGBA{R: 0, G: 200, B: 0, A: 255})
	geomPlayer = ebiten.GeoM{}
	geomPlayer.Translate(250, 350)
	screen.DrawImage(hpFillPlayer, &ebiten.DrawImageOptions{GeoM: geomPlayer})
	hpText := fmt.Sprintf("%d/%d", g.Player.HP, g.Player.MaxHP)
	text.Draw(screen, hpText, g.Font, 250+2, 330+2, shadowColor)
	text.Draw(screen, hpText, g.Font, 250, 330, textColor)

	// Полоска опыта персонажа
	const expPerLevel = 100
	expRatio := 0.0
	if expPerLevel > 0 { // Проверка на деление на ноль
		expRatio = float64(g.Player.Experience) / float64(expPerLevel)
	}
	expBarWidth := int(float64(barWidth) * expRatio)
	if expBarWidth < 1 {
		expBarWidth = 1
	}
	expBackground := ebiten.NewImage(barWidth, barHeight)
	expBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	geomExp := ebiten.GeoM{}
	geomExp.Translate(250, 380)
	screen.DrawImage(expBackground, &ebiten.DrawImageOptions{GeoM: geomExp})
	expFill := ebiten.NewImage(expBarWidth, barHeight)
	expFill.Fill(color.RGBA{R: 0, G: 255, B: 215, A: 255})
	geomExp = ebiten.GeoM{}
	geomExp.Translate(250, 380)
	screen.DrawImage(expFill, &ebiten.DrawImageOptions{GeoM: geomExp})
	expText := fmt.Sprintf("Level: %d (Exp: %d/%d)", g.Player.Level, g.Player.Experience, expPerLevel)
	text.Draw(screen, expText, g.Font, 250+2, 383+2, shadowColor)
	text.Draw(screen, expText, g.Font, 250, 383, textColor)

	// Полоска HP врага
	if g.CurrentEnemy != nil {
		hpRatio := 0.0
		if g.CurrentEnemy.MaxHP > 0 { // Проверка на деление на ноль
			hpRatio = float64(g.CurrentEnemy.HP) / float64(g.CurrentEnemy.MaxHP)
		}
		hpBarWidth := int(float64(barWidth) * hpRatio)
		if hpBarWidth < 1 {
			hpBarWidth = 1
		}
		hpBackground := ebiten.NewImage(barWidth, barHeight)
		hpBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
		geom := ebiten.GeoM{}
		geom.Translate(550, 350)
		screen.DrawImage(hpBackground, &ebiten.DrawImageOptions{GeoM: geom})
		hpFill := ebiten.NewImage(hpBarWidth, barHeight)
		hpFill.Fill(color.RGBA{R: 200, G: 0, B: 0, A: 255})
		geom = ebiten.GeoM{}
		geom.Translate(550, 350)
		screen.DrawImage(hpFill, &ebiten.DrawImageOptions{GeoM: geom})
		// Текст HP и имя врага
		enemyText := fmt.Sprintf("%s\n%d/%d", g.CurrentEnemy.Name, g.CurrentEnemy.HP, g.CurrentEnemy.MaxHP)
		text.Draw(screen, enemyText, g.Font, 550+2, 310+2, shadowColor)
		text.Draw(screen, enemyText, g.Font, 550, 310, textColor)

		// Полоски для DoT-эффектов
		const dotBarHeight = 10
		yOffset := 350 + barHeight + 5

		for _, effect := range g.CurrentEnemy.ActiveEffects {
			if dotEffect, ok := effect.(*DotEffect); ok {
				dotRatio := 0.0
				if dotEffect.Duration > 0 { // Проверка на деление на ноль
					dotRatio = dotEffect.TimeRemaining / dotEffect.Duration
				}
				if dotRatio < 0 {
					dotRatio = 0
				}
				dotBarWidth := int(float64(barWidth) * dotRatio)
				if dotBarWidth < 1 {
					dotBarWidth = 1
				}

				dotBackground := ebiten.NewImage(barWidth, dotBarHeight)
				dotBackground.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
				geom := ebiten.GeoM{}
				geom.Translate(550, float64(yOffset))
				screen.DrawImage(dotBackground, &ebiten.DrawImageOptions{GeoM: geom})

				dotFill := ebiten.NewImage(dotBarWidth, dotBarHeight)
				dotColor := color.RGBA{R: 128, G: 0, B: 128, A: 255}
				switch dotEffect.Name {
				case "Poison":
					dotColor = color.RGBA{R: 0, G: 128, B: 0, A: 255}
				case "Ignite":
					dotColor = color.RGBA{R: 255, G: 69, B: 0, A: 255}
				case "Bleed":
					dotColor = color.RGBA{R: 200, G: 0, B: 0, A: 255}
				}
				dotFill.Fill(dotColor)
				geom = ebiten.GeoM{}
				geom.Translate(550, float64(yOffset))
				screen.DrawImage(dotFill, &ebiten.DrawImageOptions{GeoM: geom})

				dotText := fmt.Sprintf("%s: %.1fs", dotEffect.Name, dotEffect.TimeRemaining)
				text.Draw(screen, dotText, g.Font, 550+2, yOffset+2+2, shadowColor)
				text.Draw(screen, dotText, g.Font, 550, yOffset+2, textColor)

				yOffset += dotBarHeight + 5
			}
		}
	}

	// Кнопки способностей
	ability1Config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), "1")
	ability1Text := "1: Basic Attack"
	if ability1Config != nil {
		ability1Text = fmt.Sprintf("1: %s", ability1Config.Name)
	}

	ability2Config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), "2")
	ability2Text := "2: Skill"
	if ability2Config != nil {
		ability2Text = fmt.Sprintf("2: %s", ability2Config.Name)
	}

	if g.AbilityCooldowns["1"] > 0 {
		ebitenutil.DrawRect(screen, 50, 600, 150, 20, color.RGBA{255, 0, 0, 128})
	} else {
		ebitenutil.DrawRect(screen, 50, 600, 150, 20, color.RGBA{0, 255, 0, 128})
	}

	if g.AbilityCooldowns["2"] > 0 {
		ebitenutil.DrawRect(screen, 50, 630, 150, 20, color.RGBA{255, 0, 0, 128})
	} else {
		ebitenutil.DrawRect(screen, 50, 630, 150, 20, color.RGBA{0, 255, 0, 128})
	}

	if g.AbilityCooldowns["3"] > 0 {
		ebitenutil.DrawRect(screen, 50, 660, 150, 20, color.RGBA{255, 0, 0, 128})
	} else {
		ebitenutil.DrawRect(screen, 50, 660, 150, 20, color.RGBA{0, 255, 0, 128})
	}

	text.Draw(screen, ability1Text, g.Font, 50+2, 600+2, shadowColor)
	text.Draw(screen, ability1Text, g.Font, 50, 600, textColor)
	text.Draw(screen, ability2Text, g.Font, 50+2, 630+2, shadowColor)
	text.Draw(screen, ability2Text, g.Font, 50, 630, textColor)
	text.Draw(screen, "3: Ult", g.Font, 50+2, 660+2, shadowColor)
	text.Draw(screen, "3: Ult", g.Font, 50, 660, textColor)

	// Лог боя
	yOffset := 700
	for i, log := range g.CombatLog {
		yPos := yOffset + 20*int(i)
		text.Draw(screen, log, g.Font, 50+2, yPos+2, shadowColor)
		text.Draw(screen, log, g.Font, 50, yPos, textColor)
	}
	if len(g.CombatLog) > 5 {
		g.CombatLog = g.CombatLog[len(g.CombatLog)-5:]
	}
}
*/
