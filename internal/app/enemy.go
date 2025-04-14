package app

import "fmt"

type Enemy struct {
	X, Y          int
	ID            string
	Name          string
	HP            int
	MaxHP         int
	Strength      int
	Agility       int
	Intelligence  int
	PhDefense     int
	MgDefense     int
	ActiveEffects []Effect
}

const MaxEffects = 4 // Максимальное количество эффектов на враге

func (e *Enemy) ApplyEffect(effect Effect) bool {
	if len(e.ActiveEffects) >= MaxEffects {
		return false // Не можем добавить эффект, если уже достигнут лимит
	}
	e.ActiveEffects = append(e.ActiveEffects, effect)
	return true
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

// Обработчик смерти врага
func (g *Game) HandleEnemyDeath() {
	if g.CurrentEnemy == nil {
		return
	}
	g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
	levelUpMsg := g.Player.AddExperience(20, g)
	if levelUpMsg != "" {
		g.CombatLog = append(g.CombatLog, levelUpMsg)
	}
	fmt.Printf("Enemies before removal: %d\n", len(g.Enemies))
	g.Enemies = removeEnemy(g.Enemies, g.CurrentEnemy.ID)
	fmt.Printf("Enemies after removal: %d\n", len(g.Enemies))
	g.CurrentEnemy = nil
	g.State = Dungeon
	if len(g.Enemies) == 0 {
		newBuff := GetRandomBuff()
		g.AvailableBuffs = append(g.AvailableBuffs, newBuff)
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("Received buff: %s", newBuff.Name()))
		newBuff.Apply(&g.Player)
	} else {
		fmt.Printf("Buff not awarded, enemies remaining: %d\n", len(g.Enemies))
	}
}
