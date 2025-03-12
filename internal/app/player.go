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

type Player struct {
	ID           string      // 16 байт, выравнивание 8
	Name         string      // 16 байт, выравнивание 8
	Class        PlayerClass // 16 байт, выравнивание 8
	Inventory    []Item      // 24 байта, выравнивание 8
	Skills       []Skill     // 24 байта, выравнивание 8
	Experience   uint32      // 4 байта, выравнивание 4
	HP           uint16      // 2 байта, выравнивание 2
	Mana         uint16      // 2 байта, выравнивание 2
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
	player := Player{
		ID:         "player1",
		Name:       "Hero",
		Class:      class,
		Inventory:  []Item{},
		Skills:     []Skill{},
		Experience: 0,
		X:          0,
		Y:          0,
		Level:      1,
	}

	// Инициализация характеристик в зависимости от класса
	switch class {
	case WarriorClass:
		player.HP = 120
		player.Mana = 30
		player.Strength = 15
		player.Agility = 8
		player.Intelligence = 5
		player.PhDefense = 10
		player.MgDefense = 5
		player.Skills = append(player.Skills, Skill{Name: "Heavy Strike", Level: 1, Cooldown: 5.0})
		player.Inventory = append(player.Inventory,
			Item{ID: 1, Name: "Sword", Price: 10.0},
			Item{ID: 2, Name: "Shield", Price: 8.0},
		)

	case MageClass:
		player.HP = 80
		player.Mana = 70
		player.Strength = 5
		player.Agility = 8
		player.Intelligence = 15
		player.PhDefense = 3
		player.MgDefense = 10
		player.Skills = append(player.Skills, Skill{Name: "Fireball", Level: 1, Cooldown: 3.0})
		player.Inventory = append(player.Inventory,
			Item{ID: 3, Name: "Staff", Price: 12.0},
			Item{ID: 4, Name: "Robe", Price: 6.0},
		)

	case ArcherClass:
		player.HP = 100
		player.Mana = 40
		player.Strength = 8
		player.Agility = 15
		player.Intelligence = 7
		player.PhDefense = 5
		player.MgDefense = 5
		player.Skills = append(player.Skills, Skill{Name: "Rapid Shot", Level: 1, Cooldown: 4.0})
		player.Inventory = append(player.Inventory,
			Item{ID: 5, Name: "Bow", Price: 9.0},
			Item{ID: 6, Name: "Quiver", Price: 5.0},
		)
	}

	return player
}
