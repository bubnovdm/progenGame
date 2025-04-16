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

		// Проверяем, сработал ли критический удар
		isCrit := g.Player.RollCrit()
		effectiveDamage := int(float64(damage) * (100.0 / (100.0 + float64(defense))))
		if isCrit {
			effectiveDamage = int(float64(effectiveDamage) * (g.Player.BaseCritDamage + g.Player.CritDamageBonus))
		}
		if effectiveDamage < minimalAADamage {
			effectiveDamage = minimalAADamage // Минимальный урон 3
		}

		g.CurrentEnemy.HP -= effectiveDamage
		g.CombatLog = append(g.CombatLog, fmt.Sprintf("Autoattack hits %s for %d %s damage. Enemy HP: %d", g.CurrentEnemy.Name, effectiveDamage, g.Player.DamageType, g.CurrentEnemy.HP))
		fmt.Printf("Autoattack hits %s for %d %s damage. Enemy HP: %d\n", g.CurrentEnemy.Name, effectiveDamage, g.Player.DamageType, g.CurrentEnemy.HP) // Отладка
		// Обновляем статистику
		g.UpdateDamageStat("Auto Attack", effectiveDamage)
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
		fmt.Printf("Player Strength: %d, Multiplier: %.2f\n", g.Player.Strength, config.Multiplier) // Отладка
		damage = int(float64(g.Player.Strength) * config.Multiplier)
	case AgilityStat:
		fmt.Printf("Player Agility: %d, Multiplier: %.2f\n", g.Player.Agility, config.Multiplier) // Отладка
		damage = int(float64(g.Player.Agility) * config.Multiplier)
	case IntelligenceStat:
		fmt.Printf("Player Intelligence: %d, Multiplier: %.2f\n", g.Player.Intelligence, config.Multiplier) // Отладка
		damage = int(float64(g.Player.Intelligence) * config.Multiplier)
	default:
		damage = int(10 * config.Multiplier)
		fmt.Printf("Default damage: 10, Multiplier: %.2f\n", config.Multiplier) // Отладка
	}

	// Кидаем крть
	isCrit := g.Player.RollCrit()

	// Применяем защиту (если не игнорируется)
	effectiveDamage := damage
	if !config.IgnoreDefense {
		var defense int
		switch g.Player.DamageType {
		case PhysicalDamage:
			defense = g.CurrentEnemy.PhDefense
		case MagicalDamage:
			defense = g.CurrentEnemy.MgDefense
		default:
			defense = 0
		}
		effectiveDamage = int(float64(damage) * (100.0 / (100.0 + float64(defense))))
	}
	if isCrit {
		effectiveDamage = int(float64(effectiveDamage) * (g.Player.BaseCritDamage + g.Player.CritDamageBonus))
	}
	if effectiveDamage < minimalSpellDamage {
		effectiveDamage = minimalSpellDamage // Минимальный урон 3
	}

	// Применяем мгновенный урон
	g.CurrentEnemy.HP -= effectiveDamage
	g.AbilityCooldowns[ability] = config.Cooldown
	g.CombatLog = append(g.CombatLog, fmt.Sprintf("Used %s for %d %s damage. Enemy HP: %d", config.Name, effectiveDamage, g.Player.DamageType, g.CurrentEnemy.HP))
	fmt.Printf("Used %s for %d %s damage. Enemy HP: %d\n", config.Name, effectiveDamage, g.Player.DamageType, g.CurrentEnemy.HP) // Отладка

	g.UpdateDamageStat(config.Name, effectiveDamage)

	// Проверяем дополнительные эффекты
	if config.DotDuration > 0 {
		//Кидаем крит
		isCrit = g.Player.RollCrit()

		var dotDamage int
		switch g.Player.MainStat {
		case StrengthStat:
			dotDamage = int(float64(g.Player.Strength) * config.DotMultiplier)
		case AgilityStat:
			dotDamage = int(float64(g.Player.Agility) * config.DotMultiplier)
		case IntelligenceStat:
			dotDamage = int(float64(g.Player.Intelligence) * config.DotMultiplier)
		default:
			dotDamage = int(10 * config.DotMultiplier)
		}

		if !config.IgnoreDefense {
			var defense int
			switch g.Player.DamageType {
			case PhysicalDamage:
				defense = g.CurrentEnemy.PhDefense
			case MagicalDamage:
				defense = g.CurrentEnemy.MgDefense
			default:
				defense = 0
			}
			dotDamage = int(float64(dotDamage) * (100.0 / (100.0 + float64(defense))) / config.DotDuration)
		}
		if isCrit {
			dotDamage = int(float64(dotDamage) * (g.Player.BaseCritDamage + g.Player.CritDamageBonus))
		}
		if dotDamage < minimalDoTDamage {
			dotDamage = minimalDoTDamage
		}

		// Определяем имя эффекта в зависимости от класса
		dotName := config.DotName
		if dotName == "" {
			dotName = "DOT" // Значение по умолчанию
		}

		dotEffect := &DotEffect{
			Name:          dotName,
			DamagePerTick: dotDamage,
			Duration:      config.DotDuration,
			TickInterval:  1.0,
			TimeRemaining: config.DotDuration,
			TickTimer:     1.0,
		}
		if g.CurrentEnemy.ApplyEffect(dotEffect) {
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s is affected by %s for %d damage per second!", g.CurrentEnemy.Name, dotName, dotDamage))
			fmt.Printf("%s is affected by %s for %d damage per second!\n", g.CurrentEnemy.Name, dotName, dotDamage) // Отладка
		} else {
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("Cannot apply %s effect to %s: too many effects!", dotName, g.CurrentEnemy.Name))
		}
	}

	if config.HitCount > 0 {
		rapidDamage := effectiveDamage
		if !config.IgnoreDefense {
			//Кидаем крит
			isCrit = g.Player.RollCrit()
			var defense int
			switch g.Player.DamageType {
			case PhysicalDamage:
				defense = int(g.CurrentEnemy.PhDefense)
			case MagicalDamage:
				defense = int(g.CurrentEnemy.MgDefense)
			default:
				defense = 0
			}
			rapidDamage = int(float64(rapidDamage) * (100.0 / (100.0 + float64(defense))))
		}
		if isCrit {
			rapidDamage = int(float64(rapidDamage) * (g.Player.BaseCritDamage + g.Player.CritDamageBonus))
		}
		if rapidDamage < 1 {
			rapidDamage = 1
		}
		rapidShot := &RapidShotEffect{
			Name:          "Rapid Shot",
			DamagePerHit:  rapidDamage,
			HitsRemaining: config.HitCount - 1,
			HitInterval:   config.HitInterval,
			TimeUntilNext: config.HitInterval,
		}
		if g.CurrentEnemy.ApplyEffect(rapidShot) {
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("%s triggers Rapid Shot on %s!", config.Name, g.CurrentEnemy.Name))
		} else {
			g.CombatLog = append(g.CombatLog, fmt.Sprintf("Cannot apply Rapid Shot to %s: too many effects!", g.CurrentEnemy.Name))
		}
	}
}
