package app

import (
	"fmt"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"sort"
)

type DamageStat struct {
	AbilityName string
	TotalDamage int
	CountUses   int
}

func (g *Game) UpdateDamageStat(abilityName string, damage int) {
	if stat, exists := g.DamageStats[abilityName]; exists {
		stat.TotalDamage += damage
		stat.CountUses++
	} else {
		g.DamageStats[abilityName] = &DamageStat{
			AbilityName: abilityName,
			TotalDamage: damage,
			CountUses:   1,
		}
	}
}

// Сортировка списка (но надо чуть доработать, при одинаковых значениях строки бесконтрольно прыгают туда-сюда
func (g *Game) GetSortedDamageStats() []DamageStat {
	stats := make([]DamageStat, 0, len(g.DamageStats))
	for _, stat := range g.DamageStats {
		stats = append(stats, *stat)
	}
	sort.Slice(stats, func(i, j int) bool {
		return stats[i].TotalDamage > stats[j].TotalDamage
	})
	return stats
}

func (g *Game) drawDamageStats(screen *ebiten.Image) {
	if g.State != Dungeon && g.State != CombatState {
		return
	}

	stats := g.GetSortedDamageStats()
	x, y := 650, 20
	ebitenutil.DebugPrintAt(screen, "Damage Stats:", x, y)
	y += 20
	for i, stat := range stats {
		if i >= 5 { // Ограничиваем количество отображаемых строк (например, топ-5)
			break
		}
		ebitenutil.DebugPrintAt(screen, fmt.Sprintf("%s: TotalDmg: %d, AvgDmg: %d, (Count: %d)", stat.AbilityName, stat.TotalDamage, stat.TotalDamage/stat.CountUses, stat.CountUses), x, y)
		y += 20
	}
}
