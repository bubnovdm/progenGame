package app

//player.go

type Item struct {
	ID    int
	Name  string
	Value float64
}

type Skill struct {
	Name     string
	Level    int
	Cooldown float64
}

type Player struct {
	ID           string  // 16 байт, выравнивание 8
	Name         string  // 16 байт, выравнивание 8
	Inventory    []Item  // 24 байта, выравнивание 8
	Skills       []Skill // 24 байта, выравнивание 8
	Experience   uint32  // 4 байта, выравнивание 4
	HP           uint16  // 2 байта, выравнивание 2
	Mana         uint16  // 2 байта, выравнивание 2
	X            int     // 8 байт, выравнивание 8
	Y            int     // 8 байт, выравнивание 8
	Level        uint8   // 1 байт, выравнивание 1
	Strength     uint8   // 1 байт, выравнивание 1
	Agility      uint8   // 1 байт, выравнивание 1
	Intelligence uint8   // 1 байт, выравнивание 1
	PhDefense    uint8   // 1 байт, выравнивание 1
	MgDefense    uint8   // 1 байт, выравнивание 1
}
