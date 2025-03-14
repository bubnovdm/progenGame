package app

import "fmt"

func (g *Game) autoAttack() {
	if g.CurrentEnemy != nil {
		var damage int
		switch g.Player.MainStat {
		case StrengthStat:
			damage = int(g.Player.Strength)
		case AgilityStat:
			damage = int(g.Player.Agility)
		case IntelligenceStat:
			damage = int(g.Player.Intelligence)
		default:
			damage = 5
		}

		var defense int
		switch g.Player.DamageType {
		case PhysicalDamage:
			defense = int(g.CurrentEnemy.PhDefense) * 2
		case MagicalDamage:
			defense = int(g.CurrentEnemy.MgDefense) * 2
		default:
			defense = 0
		}

		effectiveDamage := damage - defense
		if effectiveDamage < 0 {
			effectiveDamage = 0
		}

		g.CurrentEnemy.HP -= effectiveDamage
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("Autoattack hits %s for %d %s damage. Enemy HP: %d", g.CurrentEnemy.Name, effectiveDamage, g.Player.DamageType, g.CurrentEnemy.HP))
		if g.CurrentEnemy.HP <= 0 {
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
			enemyID := g.CurrentEnemy.ID
			g.CurrentEnemy = nil
			g.ActiveDotEffect = nil
			g.ActiveRapidShot = nil
			g.State = Dungeon
			g.Enemies = removeEnemy(g.Enemies, enemyID)
		}
	}
}

func (g *Game) useAbility(ability string) {
	if g.AbilityCooldowns[ability] > 0 {
		g.CombatLog = append(g.CombatLog, "Ability on cooldown!")
		return
	}

	// Загружаем конфигурацию способности
	config := GetAbilityConfigForClassAndKey(g.Player.Class.String(), ability)
	if config == nil {
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("Ability %s not found for class %s", ability, g.Player.Class.String()))
		return
	}

	// Обработка лечения
	if config.HealPercentage > 0 {
		healAmount := uint16(float64(g.Player.MaxHP) * config.HealPercentage)
		oldHP := g.Player.HP
		g.Player.HP += healAmount
		if g.Player.HP > g.Player.MaxHP {
			g.Player.HP = g.Player.MaxHP
		}
		g.AbilityCooldowns[ability] = config.Cooldown
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("Used %s, healed from %d to %d HP. Player HP: %d/%d", config.Name, oldHP, g.Player.HP, g.Player.HP, g.Player.MaxHP))
		return // Завершаем выполнение после лечения
	}

	// Проверяем наличие врага для атакующих способностей
	if g.CurrentEnemy == nil {
		g.CombatLog = append(g.CombatLog, "No enemy to target!")
		return
	}

	// Рассчитываем базовый урон
	var damage int
	switch g.Player.MainStat {
	case StrengthStat:
		damage = int(float64(g.Player.Strength) * config.Multiplier)
	case AgilityStat:
		damage = int(float64(g.Player.Agility) * config.Multiplier)
	case IntelligenceStat:
		damage = int(float64(g.Player.Intelligence) * config.Multiplier)
	default:
		damage = int(5 * config.Multiplier)
	}

	// Применяем защиту (если не игнорируется)
	effectiveDamage := damage
	if !config.IgnoreDefense {
		var defense int
		switch g.Player.DamageType {
		case PhysicalDamage:
			defense = int(g.CurrentEnemy.PhDefense) * 2
		case MagicalDamage:
			defense = int(g.CurrentEnemy.MgDefense) * 2
		default:
			defense = 0
		}
		effectiveDamage = damage - defense
		if effectiveDamage < 0 {
			effectiveDamage = 0
		}
	}

	// Применяем мгновенный урон
	g.CurrentEnemy.HP -= effectiveDamage
	g.AbilityCooldowns[ability] = config.Cooldown
	g.CombatLog = append(g.CombatLog, fmt.Sprintf("Used %s for %d %s damage. Enemy HP: %d", config.Name, effectiveDamage, g.Player.DamageType, g.CurrentEnemy.HP))

	// Проверяем дополнительные эффекты
	if config.DotDuration > 0 {
		dotDamage := int(float64(g.Player.Intelligence) * config.DotMultiplier)
		g.ActiveDotEffect = &DotEffect{
			DamagePerTick: dotDamage,
			Duration:      config.DotDuration,
			TickInterval:  1.0,
			TimeRemaining: config.DotDuration,
		}
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s is burning for %d damage per second!", g.CurrentEnemy.Name, dotDamage))
	}

	if config.HitCount > 0 {
		g.ActiveRapidShot = &RapidShotEffect{
			DamagePerHit:  effectiveDamage,
			HitsRemaining: config.HitCount - 1,
			HitInterval:   config.HitInterval,
			TimeUntilNext: config.HitInterval,
		}
	}

	if g.CurrentEnemy.HP <= 0 {
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s defeated!", g.CurrentEnemy.Name))
		enemyID := g.CurrentEnemy.ID
		g.CurrentEnemy = nil
		g.ActiveDotEffect = nil
		g.ActiveRapidShot = nil
		g.State = Dungeon
		g.Enemies = removeEnemy(g.Enemies, enemyID)
	}
}
