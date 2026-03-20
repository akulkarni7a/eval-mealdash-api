package routing

import (
	"math"
	"time"
)

// RestaurantPenalty tracks slow-restaurant detection scores.
// When avg driver_pickup_duration exceeds 6 minutes over a 3-day window,
// the penalty accumulates at 0.1/day. Decays 5%/day.
// When penalty > 0.5, restaurant is deprioritized in routing.
type RestaurantPenalty struct {
	RestaurantID string
	Score        float64
	LastUpdated  time.Time
}

const (
	pickupThreshold  = 6 * time.Minute
	penaltyIncrement = 0.1
	decayRate        = 0.05
	deprioritizeAt   = 0.5
	lookbackDays     = 3
)

func (p *RestaurantPenalty) Update(avgPickupDuration time.Duration) {
	// Apply daily decay
	daysSinceUpdate := time.Since(p.LastUpdated).Hours() / 24
	p.Score *= math.Pow(1-decayRate, daysSinceUpdate)

	// Check if restaurant is slow
	if avgPickupDuration > pickupThreshold {
		p.Score += penaltyIncrement
	}

	p.LastUpdated = time.Now()
}

func (p *RestaurantPenalty) IsDeprioritized() bool {
	return p.Score > deprioritizeAt
}
