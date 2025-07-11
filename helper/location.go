package helper

import (
	"math"
)

// HaversineDistance calculates the distance between two points on Earth using the Haversine formula.
// lat1, lon1 are the coordinates of the first point (in degrees).
// lat2, lon2 are the coordinates of the second point (in degrees).
// Returns the distance in meters.
func HaversineDistance(lat1, lon1, lat2, lon2 float64) float64 {
	const R = 6371e3 // Earth's radius in meters

	phi1 := lat1 * math.Pi / 180
	phi2 := lat2 * math.Pi / 180
	deltaPhi := (lat2 - lat1) * math.Pi / 180
	deltaLambda := (lon2 - lon1) * math.Pi / 180

	a := math.Sin(deltaPhi/2)*math.Sin(deltaPhi/2) + 
		math.Cos(phi1)*math.Cos(phi2)* 
		math.Sin(deltaLambda/2)*math.Sin(deltaLambda/2)
	c := 2 * math.Atan2(math.Sqrt(a), math.Sqrt(1-a))

	distance := R * c
	return distance
}
