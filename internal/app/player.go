package app

// player.go

type PlayerClass string

const (
	WarriorClass PlayerClass = "Warrior" // Воин
	MageClass    PlayerClass = "Mage"    // Маг
	ArcherClass  PlayerClass = "Archer"  // Лучник
)

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
	Experience   uint32      // 4 байта, выравнивание 4
	HP           uint16      // 2 байта, выравнивание 2
	Mana         uint16      // 2 байта, выравнивание 2
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
	var maxHP int
	var strength, agility, intelligence int
	var mainStat MainStat
	var damageType DamageType

	switch class {
	case WarriorClass:
		maxHP = 120
		strength = 10
		agility = 5
		intelligence = 5
		mainStat = StrengthStat
		damageType = PhysicalDamage
	case MageClass:
		maxHP = 80
		strength = 5
		agility = 5
		intelligence = 10
		mainStat = IntelligenceStat
		damageType = MagicalDamage
	case ArcherClass:
		maxHP = 100
		strength = 5
		agility = 10
		intelligence = 5
		mainStat = AgilityStat
		damageType = PhysicalDamage
	default:
		maxHP = 40
		strength = 5
		agility = 5
		intelligence = 5
		mainStat = StrengthStat
		damageType = PhysicalDamage
	}

	return Player{
		X:            0,
		Y:            0,
		HP:           uint16(maxHP),
		MaxHP:        uint16(maxHP),
		Mana:         30,
		Strength:     uint8(strength),
		Agility:      uint8(agility),
		Intelligence: uint8(intelligence),
		PhDefense:    5,
		MgDefense:    5,
		Class:        class,
		Inventory:    []Item{},
		Skills:       []Skill{},
		MainStat:     mainStat,
		DamageType:   damageType,
	}
}
