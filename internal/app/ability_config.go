package app

import (
	"encoding/json"
	"fmt"
	"os"
)

// AbilityConfig описывает параметры способности
type AbilityConfig struct {
	Key            string  `json:"key"`
	Name           string  `json:"name"`
	Multiplier     float64 `json:"multiplier"`
	IgnoreDefense  bool    `json:"ignore_defense,omitempty"`
	DotMultiplier  float64 `json:"dot_multiplier,omitempty"`
	DotDuration    float64 `json:"dot_duration,omitempty"`
	HitCount       int     `json:"hit_count,omitempty"`
	HitInterval    float64 `json:"hit_interval,omitempty"`
	HealPercentage float64 `json:"heal_percentage,omitempty"`
	Cooldown       float64 `json:"cooldown"`
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
