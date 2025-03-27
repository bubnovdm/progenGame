package app

import (
	"encoding/json"
	"fmt"
	"os"
)

type GameSave struct {
	Player         Player      `json:"player"`          // Характеристики персонажа
	MaxFloor       int         `json:"max_floor"`       // Максимальный этаж, до которого дошел игрок
	CurrentFloor   int         `json:"current_floor"`   // Текущий этаж
	SelectedClass  PlayerClass `json:"selected_class"`  // Выбранный класс
	AvailableBuffs []BuffData  `json:"available_buffs"` // Список бафов игрока
}

func (g *Game) SaveGame() error {
	buffData := make([]BuffData, len(g.AvailableBuffs))
	for i, buff := range g.AvailableBuffs {
		buffData[i] = ToBuffData(buff)
	}

	saveData := GameSave{
		Player:         g.Player,
		CurrentFloor:   g.CurrentFloor,
		MaxFloor:       g.MaxFloor,
		SelectedClass:  g.Player.Class,
		AvailableBuffs: buffData,
	}

	data, err := json.MarshalIndent(saveData, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal save data: %v", err)
	}

	err = os.WriteFile("save.json", data, 0644)
	if err != nil {
		return fmt.Errorf("failed to write save file: %v", err)
	}

	fmt.Println("Game saved successfully")
	return nil
}

// LoadGame загружает прогресс из файла
func (g *Game) LoadGame() error {
	data, err := os.ReadFile("save.json")
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to read save file: %v", err)
	}

	var saveData GameSave
	err = json.Unmarshal(data, &saveData)
	if err != nil {
		return fmt.Errorf("failed to unmarshal save data: %v", err)
	}

	g.Player = saveData.Player
	g.CurrentFloor = saveData.CurrentFloor
	g.MaxFloor = saveData.MaxFloor
	g.SelectedFloor = saveData.CurrentFloor
	g.Player.Class = saveData.SelectedClass

	g.AvailableBuffs = make([]Buff, len(saveData.AvailableBuffs))
	for i, buffData := range saveData.AvailableBuffs {
		g.AvailableBuffs[i] = FromBuffData(buffData)
	}

	var mapType MapType
	if g.CurrentFloor%2 == 0 {
		mapType = OpenMap
	} else {
		mapType = DungeonMap
	}
	g.GameMap = GenerateMap(mapType)
	g.moveToStartPosition()
	g.spawnEnemies()

	// Предполагается, что ApplyBuffs есть в другом файле
	g.ApplyBuffs()

	fmt.Println("Game loaded successfully")
	return nil
}
