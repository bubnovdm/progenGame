package app

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

// EnemyStats описывает базовые характеристики врага
type EnemyStats struct {
	Name         string `json:"name"`
	HP           int    `json:"hp"`
	HPPerLevel   int    `json:"hp_per_level"`
	Strength     int    `json:"strength"`
	Agility      int    `json:"agility"`
	Intelligence int    `json:"intelligence"`
	PhDefense    int    `json:"ph_defense"`
	MgDefense    int    `json:"mg_defense"`
}

// EnemyConfig описывает конфигурацию врага
type EnemyConfig struct {
	Type      string     `json:"type"`
	Levels    []int      `json:"levels"`
	BaseStats EnemyStats `json:"base_stats"`
}

// EnemyConfigs хранит все конфигурации врагов
var enemyConfigs []EnemyConfig

// LoadEnemyConfigs загружает конфигурации врагов из JSON-файла
func LoadEnemyConfigs(filePath string) error {
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("failed to read enemy config file: %v", err)
	}

	err = json.Unmarshal(data, &enemyConfigs)
	if err != nil {
		return fmt.Errorf("failed to unmarshal enemy config: %v", err)
	}

	return nil
}

// GetEnemyConfigForLevel возвращает конфигурацию врага для заданного уровня
func GetEnemyConfigForLevel(level int) *EnemyConfig {
	for _, config := range enemyConfigs {
		for _, lvl := range config.Levels {
			if lvl == level {
				return &config
			}
		}
	}
	// Если уровень не найден, возвращаем первого врага (например, гоблина)
	if len(enemyConfigs) > 0 {
		return &enemyConfigs[0]
	}
	return nil
}
