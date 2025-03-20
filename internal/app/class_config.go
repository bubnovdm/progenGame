package app

import (
	"encoding/json"
	"fmt"
	"os"
)

// ClassStats описывает базовые характеристики класса
type ClassStats struct {
	MaxHP        int `json:"max_hp"`
	Strength     int `json:"strength"`
	Agility      int `json:"agility"`
	Intelligence int `json:"intelligence"`
	PhDefense    int `json:"ph_defense"`
	MgDefense    int `json:"mg_defense"`
}

type LevelUpStats struct {
	MaxHP        int `json:"max_hp"`
	Strength     int `json:"strength"`
	Agility      int `json:"agility"`
	Intelligence int `json:"intelligence"`
	PhDefense    int `json:"ph_defense"`
	MgDefense    int `json:"mg_defense"`
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
