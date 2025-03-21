package app

import (
	"encoding/json"
	"fmt"
	"os"
)

// AbilityConfig описывает параметры способности
type AbilityConfig struct {
	Class          string  `json:"class"`
	Key            string  `json:"key"`
	Name           string  `json:"name"`
	Multiplier     float64 `json:"multiplier"`
	Cooldown       float64 `json:"cooldown"`
	DotDuration    float64 `json:"dot_duration"`
	DotMultiplier  float64 `json:"dot_multiplier"`
	DotName        string  `json:"dot_name"`
	HitCount       int     `json:"hit_count"`
	HitInterval    float64 `json:"hit_interval"`
	IgnoreDefense  bool    `json:"ignore_defense"`
	HealPercentage float64 `json:"heal_percentage"`
}

// ClassAbilityConfig описывает способности для класса
type ClassAbilityConfig struct {
	Class     string          `json:"class"`
	Abilities []AbilityConfig `json:"abilities"`
}

var abilityConfigs []ClassAbilityConfig

// LoadAbilityConfigs загружает конфигурации способностей из JSON-файла
func LoadAbilityConfigs(filePath string) error {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("Failed to read ability config file: %v", err)
	}

	err = json.Unmarshal(data, &abilityConfigs)
	if err != nil {
		return fmt.Errorf("Failed to unmarshal ability config: %v", err)
	}

	return nil
}

// GetAbilityConfigForClassAndKey возвращает конфигурацию способности для класса по ключу
func GetAbilityConfigForClassAndKey(class string, key string) *AbilityConfig {
	for _, classConfig := range abilityConfigs {
		if classConfig.Class == class {
			for _, ability := range classConfig.Abilities {
				if ability.Key == key {
					return &ability
				}
			}
		}
	}
	return nil
}
