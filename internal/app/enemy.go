package app

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
