package app

import (
	"fmt"
	"github.com/bubnovdm/progenGame/internal/utils"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"image/color"
	"math"
	"strconv"
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Floor: %d", g.CurrentFloor), 450, 15)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Enemies remaining: %d", len(g.Enemies)), 415, 25)
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
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CritChance: %g%%", g.Player.GetTotalCritChance()), 60, 320)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("CritDamage: %v%%", (g.Player.BaseCritDamage+g.Player.CritDamageBonus)*100), 60, 350)

	// Бафы (справа от характеристик)
	//buffsBg := ebiten.NewImage(300, 150) // Фон для бафов, можно будет попробовать добавить прокрутку, если бафов будет много
	//buffsBg.Fill(color.RGBA{R: 50, G: 50, B: 50, A: 255})
	//geom := ebiten.GeoM{}
	//geom.Translate(700, 110) // Позиция справа, на уровне характеристик
	//screen.DrawImage(buffsBg, &ebiten.DrawImageOptions{GeoM: geom})

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

				// Заполнение полоски DoT
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

	ability3Config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), "3")
	ability3Text := "3: Skill 3"
	if ability3Config != nil {
		ability3Text = fmt.Sprintf("3: %s", ability3Config.Name)
	}

	ability4Config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), "4")
	ability4Text := "4: Ult"
	if ability4Config != nil {
		ability4Text = fmt.Sprintf("4: %s", ability4Config.Name)
	}

	// Подложка под названия способностей для отслеживания кд
	for i := range g.AbilityCooldowns {
		index, _ := strconv.Atoi(i)
		y := float64(600 + 30*(index-1))
		var bgColor color.Color
		if g.AbilityCooldowns[i] > 0 {
			bgColor = color.RGBA{255, 0, 0, 128} // Красный фон, если на кд
		} else {
			bgColor = color.RGBA{0, 255, 0, 128} // Зеленый фон, если готово
		}
		ebitenutil.DrawRect(screen, 50, y, 150, 20, bgColor)
	}

	ebitenutil.DebugPrintAt(screen, ability1Text, 50, 600)
	ebitenutil.DebugPrintAt(screen, ability2Text, 50, 630)
	ebitenutil.DebugPrintAt(screen, ability3Text, 50, 660)
	ebitenutil.DebugPrintAt(screen, ability4Text, 50, 690)

	// Лог боя
	yOffset := 720
	for i, log := range g.CombatLog {
		ebitenutil.DebugPrintAt(screen, log, 50, yOffset+20*int(i))
	}
	if len(g.CombatLog) >= 6 {
		g.CombatLog = g.CombatLog[1:]
	}
}
