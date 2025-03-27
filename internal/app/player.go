package app

import "fmt"

// player.go

type PlayerClass uint8

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
	ID                           string      // 16 байт
	Name                         string      // 16 байт
	X                            int         // 8 байт
	Y                            int         // 8 байт
	Inventory                    []Item      // 24 байта
	Skills                       []Skill     // 24 байта
	Class                        PlayerClass // 16 байт
	MainStat                     MainStat    // 16 байт
	DamageType                   DamageType  // 16 байт
	AutoAttackCooldownMultiplier float64     // 8 байт
	CritChance                   float64     // 8 байт - шанс крита в %
	CritDamage                   float64     // 8 байт - множитель крита (коэффициент)
	HP                           uint16      // 2 байта
	MaxHP                        uint16      // 2 байта
	Strength                     uint16      // 2 байта
	Agility                      uint16      // 2 байта
	Intelligence                 uint16      // 2 байта
	PhDefense                    uint16      // 2 байта
	MgDefense                    uint16      // 2 байта
	Experience                   uint8       // 1 байт
	Level                        uint8       // 1 байт

}

func NewPlayer(class PlayerClass, g *Game) Player {
	// Получаем конфигурацию класса из g.ClassConfig
	var classConfig ClassConfig
	if g != nil {
		var ok bool
		classConfig, ok = g.ClassConfig[class.String()]
		if !ok {
			// Если класс не найден, используем Warrior по умолчанию
			classConfig = g.ClassConfig["Warrior"]
		}
	} else {
		// Если g == nil, используем значения по умолчанию
		classConfig = ClassConfig{
			Type: class.String(),
			BaseStats: ClassStats{
				MaxHP:        120,
				Strength:     10,
				Agility:      5,
				Intelligence: 5,
				PhDefense:    7,
				MgDefense:    3,
			},
			MainStat:   StrengthStat,
			DamageType: PhysicalDamage,
		}
	}

	baseStats := classConfig.BaseStats
	return Player{
		Class:                        class,
		Level:                        1,
		Experience:                   0,
		HP:                           baseStats.MaxHP,
		MaxHP:                        baseStats.MaxHP,
		Strength:                     baseStats.Strength,
		Agility:                      baseStats.Agility,
		Intelligence:                 baseStats.Intelligence,
		PhDefense:                    baseStats.PhDefense,
		MgDefense:                    baseStats.MgDefense,
		MainStat:                     classConfig.MainStat,
		DamageType:                   classConfig.DamageType,
		AutoAttackCooldownMultiplier: 1.0,
		X:                            1,
		Y:                            1,
		Inventory:                    []Item{},
		Skills:                       []Skill{},
	}
}

// AddExperience добавляет опыт и проверяет повышение уровня
func (p *Player) AddExperience(exp uint8, g *Game) string {
	p.Experience += exp
	// Проверка на повышение уровня
	if p.Experience >= 100 {
		p.Experience -= 100
		p.LevelUp(g)
		return fmt.Sprintf("CurrentFloor Up! You are now level %d", p.Level)
	}
	return ""
}

// LevelUp повышает уровень и улучшает характеристики
func (p *Player) LevelUp(g *Game) {
	p.Level++

	// Получаем конфигурацию класса
	classConfig, ok := g.ClassConfig[p.Class.String()]
	if !ok {
		// Если класс не найден, используем значения по умолчанию (например, Warrior)
		classConfig = g.ClassConfig["Warrior"]
	}

	// Применяем прирост характеристик из level_up_stats
	levelUpStats := classConfig.LevelUpStats
	p.MaxHP += uint16(levelUpStats.MaxHP)
	p.HP += uint16(levelUpStats.MaxHP)
	p.Strength += levelUpStats.Strength
	p.Agility += levelUpStats.Agility
	p.Intelligence += levelUpStats.Intelligence
	p.PhDefense += levelUpStats.PhDefense
	p.MgDefense += levelUpStats.MgDefense

	// Убедимся, что HP не превышает MaxHP
	if p.HP > p.MaxHP {
		p.HP = p.MaxHP
	}
}
