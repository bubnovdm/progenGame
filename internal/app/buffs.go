package app

import "math/rand"

// Buff описывает пассивный постоянный баф
type Buff interface {
	Apply(player *Player) // Применяет баф к игроку
	Name() string         // Возвращает имя бафа (для отображения в UI)
}

// MainStatBuff увеличивает мейнстат игрока (Strength, Agility или Intelligence в зависимости от класса)
type MainStatBuff struct {
	amount uint16
}

func (b *MainStatBuff) Apply(player *Player) {
	switch player.Class {
	case WarriorClass:
		player.Strength += b.amount
	case ArcherClass:
		player.Agility += b.amount
	case MageClass:
		player.Intelligence += b.amount
	}
}

func (b *MainStatBuff) Name() string {
	return "Main Stat Boost"
}

// HealthBuff увеличивает максимальное здоровье игрока
type HealthBuff struct {
	amount uint16
}

func (b *HealthBuff) Apply(player *Player) {
	player.MaxHP += b.amount
	// Если текущее здоровье меньше нового максимума, увеличиваем его
	if player.HP > player.MaxHP {
		player.HP = player.MaxHP
	}
}

func (b *HealthBuff) Name() string {
	return "Health Boost"
}

// AttackSpeedBuff увеличивает скорость автоатак (уменьшает кулдаун автоатак на 5%)
type AttackSpeedBuff struct {
	speedIncrease float64 // Процент увеличения скорости (например, 0.05 для 5%)
}

func (b *AttackSpeedBuff) Apply(player *Player) {
	// Уменьшаем кулдаун автоатак (увеличиваем скорость)
	// Если текущий кулдаун = 2.0, то после бафа он станет 2.0 * (1 - 0.05) = 1.9
	player.AutoAttackCooldownMultiplier *= (1 - b.speedIncrease)
	// Устанавливаем минимальный кулдаун, чтобы избежать 0
	if player.AutoAttackCooldownMultiplier < 0.1 {
		player.AutoAttackCooldownMultiplier = 0.1
	}
}

func (b *AttackSpeedBuff) Name() string {
	return "Attack Speed Boost"
}

// GetRandomBuff возвращает случайный баф
func GetRandomBuff() Buff {
	buffs := []Buff{
		&MainStatBuff{amount: 2},              // +2 к мейнстату
		&HealthBuff{amount: 10},               // +10 к хп
		&AttackSpeedBuff{speedIncrease: 0.05}, // 5% увеличение скорости
	}
	return buffs[rand.Intn(len(buffs))]
}

// BuffData используется для сериализации/десериализации бафов
type BuffData struct {
	Name          string  `json:"name"`
	Amount        uint16  `json:"amount"`
	SpeedIncrease float64 `json:"speed_increase"`
}

// ToBuffData конвертирует Buff в BuffData
func ToBuffData(buff Buff) BuffData {
	data := BuffData{Name: buff.Name()}
	switch b := buff.(type) {
	case *MainStatBuff:
		data.Amount = b.amount
	case *HealthBuff:
		data.Amount = b.amount
	case *AttackSpeedBuff:
		data.SpeedIncrease = b.speedIncrease
	}
	return data
}

// FromBuffData создает Buff из BuffData
func FromBuffData(data BuffData) Buff {
	switch data.Name {
	case "Main Stat Boost":
		return &MainStatBuff{amount: data.Amount}
	case "Health Boost":
		return &HealthBuff{amount: data.Amount}
	case "Attack Speed Boost":
		return &AttackSpeedBuff{speedIncrease: data.SpeedIncrease}
	default:
		return nil // или дефолтный баф
	}
}

func (g *Game) ApplyBuffs() {
	for _, buff := range g.AvailableBuffs {
		buff.Apply(&g.Player)
	}
}
