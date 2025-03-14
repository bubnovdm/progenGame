package app

import "log"

// player.go

type PlayerClass int

const (
	WarriorClass PlayerClass = 0
	MageClass    PlayerClass = 1
	ArcherClass  PlayerClass = 2
)

func (c PlayerClass) String() string {
	switch c {
	case WarriorClass:
		return "Warrior"
	case MageClass:
		return "Mage"
	case ArcherClass:
		return "Archer"
	default:
		return "Unknown"
	}
}

type Item struct {
	ID    int
	Name  string
	Price float64
}

type Skill struct {
	Name     string
	Level    int
	Cooldown float64
}

// DamageType определяет тип урона
type DamageType string

const (
	PhysicalDamage DamageType = "physical"
	MagicalDamage  DamageType = "magical"
)

// MainStat определяет основную характеристику для урона
type MainStat string

const (
	StrengthStat     MainStat = "strength"
	AgilityStat      MainStat = "agility"
	IntelligenceStat MainStat = "intelligence"
)

type Player struct {
	ID           string // 16 байт, выравнивание 8
	Name         string // 16 байт, выравнивание 8
	MainStat     MainStat
	DamageType   DamageType
	Class        PlayerClass // 16 байт, выравнивание 8
	Inventory    []Item      // 24 байта, выравнивание 8
	Skills       []Skill     // 24 байта, выравнивание 8
	Experience   uint64      // 4 байта, выравнивание 4
	HP           uint16      // 2 байта, выравнивание 2
	MaxHP        uint16      // 2 байта, выравнивание 2
	X            int         // 8 байт, выравнивание 8
	Y            int         // 8 байт, выравнивание 8
	Level        uint8       // 1 байт, выравнивание 1
	Strength     uint8       // 1 байт, выравнивание 1
	Agility      uint8       // 1 байт, выравнивание 1
	Intelligence uint8       // 1 байт, выравнивание 1
	PhDefense    uint8       // 1 байт, выравнивание 1
	MgDefense    uint8       // 1 байт, выравнивание 1

}

func NewPlayer(class PlayerClass) Player {
	config := GetClassConfigForType(class.String())
	if config == nil {
		log.Printf("Warning: No config found for class %s, using default Warrior", class.String())
		return Player{
			X:            0,
			Y:            0,
			HP:           50,
			MaxHP:        50,
			Strength:     10,
			Agility:      5,
			Intelligence: 5,
			PhDefense:    5,
			MgDefense:    5,
			Class:        class,
			Inventory:    []Item{},
			Skills:       []Skill{},
			MainStat:     StrengthStat,
			DamageType:   PhysicalDamage,
		}
	}

	return Player{
		X:            0,
		Y:            0,
		HP:           uint16(config.BaseStats.MaxHP),       // Приведение к uint16
		MaxHP:        uint16(config.BaseStats.MaxHP),       // Приведение к uint16
		Strength:     uint8(config.BaseStats.Strength),     // Приведение к uint8
		Agility:      uint8(config.BaseStats.Agility),      // Приведение к uint8
		Intelligence: uint8(config.BaseStats.Intelligence), // Приведение к uint8
		PhDefense:    uint8(config.BaseStats.PhDefense),    // Приведение к uint8
		MgDefense:    uint8(config.BaseStats.MgDefense),    // Приведение к uint8
		Class:        class,
		Inventory:    []Item{},
		Skills:       []Skill{},
		MainStat:     config.MainStat,
		DamageType:   config.DamageType,
	}
}

// AddExperience добавляет опыт и проверяет повышение уровня
func (p *Player) AddExperience(exp uint64) {
	p.Experience += exp
	const expPerLevel = 100 // Опыт для следующего уровня
	for p.Experience >= expPerLevel {
		p.Experience -= expPerLevel
		p.LevelUp()
	}
}

// LevelUp повышает уровень и улучшает характеристики
func (p *Player) LevelUp() {
	p.Level++
	switch p.Class {
	case WarriorClass:
		p.MaxHP += 10
		p.HP += 10
		p.Strength += 2
		p.Agility += 1
		p.Intelligence += 1
		p.PhDefense += 3
		p.MgDefense += 1
	case MageClass:
		p.MaxHP += 5
		p.HP += 5
		p.Strength += 1
		p.Agility += 1
		p.Intelligence += 3
		p.PhDefense += 1
		p.MgDefense += 3
	case ArcherClass:
		p.MaxHP += 7
		p.HP += 7
		p.Strength += 1
		p.Agility += 3
		p.Intelligence += 1
		p.PhDefense += 2
		p.MgDefense += 2
	}
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
}
