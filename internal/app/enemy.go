package app

//enemy.go

type Ability struct {
	Name     string
	Level    int
	Cooldown float64
}

type Enemy struct {
	ID        string    // Уникальный идентификатор врага
	Name      string    // Имя врага
	Level     uint8     // Уровень
	Strength  uint8     // Сила
	Defense   uint8     // Защита
	HP        uint16    // Здоровье
	Abilities []Ability // Умения
}

type Goblin struct {
	Enemy
	SpecialDrop string // Особая вещь, выпадающая с Гоблина
}

type Dragon struct {
	Enemy
	FireBreathDamage uint16 // Урон от огненного дыхания
}
