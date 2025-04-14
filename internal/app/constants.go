package app

// Сюда, по-хорошему, нужно вынести все константы, что есть в остальных файлах

// Константы для классов
type PlayerClass uint8

const (
	WarriorClass PlayerClass = 0
	MageClass    PlayerClass = 1
	ArcherClass  PlayerClass = 2
)

// DamageType определяет тип урона
type DamageType string

const (
	PhysicalDamage DamageType = "physical"
	MagicalDamage  DamageType = "magical"
)

const (
	minimalAADamage    = 3
	minimalSpellDamage = 3
	minimalDoTDamage   = 1
	minimalEnemyDamage = 3
)

// MainStat определяет основную характеристику для урона
type MainStat string

const (
	StrengthStat     MainStat = "strength"
	AgilityStat      MainStat = "agility"
	IntelligenceStat MainStat = "intelligence"
)

// Константы для карты

const (
	MapSize          = 40
	PathLength       = 100
	StartSymbol      = 'S'
	ExitSymbol       = 'X'
	PathSymbol       = '1'
	EmptySymbol      = '0'
	BackgroundSymbol = 'G'
	WallSymbol       = 'W'
)

type MapType int

const (
	OpenMap MapType = iota
	DungeonMap
)

const (
	maxRooms    = 10
	minRoomSize = 4
	maxRoomSize = 8
)

// Размер выпадающего списка в CharacterSheet
const (
	floorButtonWidth  = 200
	floorButtonHeight = 30
	floorButtonX      = 700
	floorButtonY      = 360
)
