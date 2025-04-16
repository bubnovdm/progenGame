package app

import (
	"encoding/json"
	"fmt"
	"os"
)

// ClassStats описывает базовые характеристики класса
type ClassStats struct {
	MaxHP          uint16  `json:"max_hp"`
	Strength       uint16  `json:"strength"`
	Agility        uint16  `json:"agility"`
	Intelligence   uint16  `json:"intelligence"`
	PhDefense      uint16  `json:"ph_defense"`
	MgDefense      uint16  `json:"mg_defense"`
	BaseCritChance float64 `json:"crit_chance"`
	BaseCritDamage float64 `json:"crit_damage"`
}

type LevelUpStats struct {
	MaxHP        uint16 `json:"max_hp"`
	Strength     uint16 `json:"strength"`
	Agility      uint16 `json:"agility"`
	Intelligence uint16 `json:"intelligence"`
	PhDefense    uint16 `json:"ph_defense"`
	MgDefense    uint16 `json:"mg_defense"`
}

// ClassConfig описывает конфигурацию класса
type ClassConfig struct {
	Type         string       `json:"type"`
	BaseStats    ClassStats   `json:"base_stats"`
	MainStat     MainStat     `json:"main_stat"`
	DamageType   DamageType   `json:"damage_type"`
	LevelUpStats LevelUpStats `json:"level_up_stats"`
}

// ClassConfigs хранит все конфигурации классов
var classConfigs []ClassConfig

// LoadClassConfigs загружает конфигурации классов из JSON-файла
func LoadClassConfigs(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read class config file: %v", err)
	}

	err = json.Unmarshal(data, &classConfigs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal class config: %v", err)
	}

	return nil
}

// ToMap преобразует слайс ClassConfigs в мапу для удобного доступа по имени класса
func ToMap() map[string]ClassConfig {
	configMap := make(map[string]ClassConfig)
	for _, cfg := range classConfigs {
		configMap[cfg.Type] = cfg
	}
	return configMap
}
