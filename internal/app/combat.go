package app

import "fmt"

func (g *Game) autoAttack() {
	if g.CurrentEnemy != nil {
		// Определяем базовый урон на основе главной характеристики
		var damage int
		switch g.Player.MainStat {
		case StrengthStat:
			damage = int(g.Player.Strength)
		case AgilityStat:
			damage = int(g.Player.Agility)
		case IntelligenceStat:
			damage = int(g.Player.Intelligence)
		default:
			damage = 5 // Запасной вариант
		}

		// Определяем защиту врага в зависимости от типа урона
		var defense int
		switch g.Player.DamageType {
		case PhysicalDamage:
			defense = int(g.CurrentEnemy.PhDefense) * 2
		case MagicalDamage:
			defense = int(g.CurrentEnemy.MgDefense) * 2
		default:
			defense = 0
		}

		// Рассчитываем итоговый урон
		effectiveDamage := damage - defense
		if effectiveDamage < 0 {
			effectiveDamage = 0
		}

		g.CurrentEnemy.HP -= effectiveDamage // Изменяем только HP
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("Autoattack hits %s for %d damage. Enemy HP: %d", g.CurrentEnemy.Name, effectiveDamage, g.CurrentEnemy.HP))
		if g.CurrentEnemy.HP <= 0 {
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
			enemyID := g.CurrentEnemy.ID // Сохраняем ID перед обнулением
			g.CurrentEnemy = nil
			g.State = Dungeon
			g.Enemies = removeEnemy(g.Enemies, enemyID)
		}
	}
}

func (g *Game) useAbility(ability string) {
	switch ability {
	case "1":
		if g.AbilityCooldowns["BasicAttack"] <= 0 {
			if g.CurrentEnemy != nil {
				damage := int(g.Player.Strength) * 3
				effectiveDamage := damage - (int(g.CurrentEnemy.PhDefense) * 2)
				if effectiveDamage < 0 {
					effectiveDamage = 0
				}
				g.CurrentEnemy.HP -= effectiveDamage // Изменяем только HP
				g.AbilityCooldowns["BasicAttack"] = 3.0
				g.CombatLog = append(g.CombatLog, fmt.Sprintf("Used Basic Attack for %d damage. Enemy HP: %d", effectiveDamage, g.CurrentEnemy.HP))
				if g.CurrentEnemy.HP <= 0 {
					g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
					enemyID := g.CurrentEnemy.ID // Сохраняем ID перед обнулением
					g.CurrentEnemy = nil
					g.State = Dungeon
					g.Enemies = removeEnemy(g.Enemies, enemyID)
				}
			}
		} else {
			g.CombatLog = append(g.CombatLog, "Basic Attack on cooldown!")
		}
	}
}
