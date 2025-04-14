package app

import "fmt"

type Effect interface {
	Update(dt float64, target *Enemy, g *Game) (logMessage string) // Обновляет эффект и возвращает сообщение для лога
	IsFinished() bool                                              // Возвращает true, если эффект завершён
}

// Это структура для доток. Ожог от фаербола мага, кровотечение от вара, мб в будущем что-то ещё
type DotEffect struct {
	Name          string  // Имя эффекта
	DamagePerTick int     // Урон за тик
	Duration      float64 // Общая длительность эффекта
	TickInterval  float64 // Интервал между тиками
	TimeRemaining float64 // Оставшееся время действия эффекта
	TickTimer     float64 // Таймер до следующего тика
}

func (e *DotEffect) Update(dt float64, target *Enemy, g *Game) (logMessage string) {
	e.TimeRemaining -= dt
	if e.TimeRemaining <= 0 {
		return ""
	}

	e.TickTimer -= dt
	if e.TickTimer <= 0 {
		target.HP -= e.DamagePerTick
		e.TickTimer = e.TickInterval // Сбрасываем таймер тика
		g.UpdateDamageStat(e.Name, e.DamagePerTick)
		return fmt.Sprintf("%s takes %d %s damage. Enemy HP: %d", target.Name, e.DamagePerTick, e.Name, target.HP)
	}
	return ""
}

func (e *DotEffect) IsFinished() bool {
	return e.TimeRemaining <= 0
}

// Это структура для потокового урона, типа быстрой стрельбы ханта. Мб что-то ещё появится в будущем из идей.
type RapidShotEffect struct {
	Name          string  // Имя эффекта
	DamagePerHit  int     // Урон за удар
	HitsRemaining int     // Оставшееся количество ударов
	HitInterval   float64 // Интервал между ударами
	TimeUntilNext float64 // Время до следующего удара
}

func (e *RapidShotEffect) Update(dt float64, target *Enemy, g *Game) (logMessage string) {
	e.TimeUntilNext -= dt
	if e.TimeUntilNext <= 0 && e.HitsRemaining > 0 {
		target.HP -= e.DamagePerHit
		e.HitsRemaining--
		e.TimeUntilNext = e.HitInterval
		g.UpdateDamageStat(e.Name, e.DamagePerHit)
		return fmt.Sprintf("%s hits %s for %d damage. Enemy HP: %d", e.Name, target.Name, e.DamagePerHit, target.HP)
	}
	return ""
}

func (e *RapidShotEffect) IsFinished() bool {
	return e.HitsRemaining <= 0
}
