package utility

import (
	"strings"
	"github.com/kellydunn/golang-geo"
)

type PokeUtility struct {
}

func (u *PokeUtility) LocateValueInKeyWithMapIntString(str string, mapstr map[int]string) (int, error) {
	if len(str) == 0 {
		return -1, nil
	} else {
		for k, v := range mapstr {
			if strings.EqualFold(v, str) {
				return k, nil
			}
		}
		return -1, nil
	}
}
func (u *PokeUtility) CalculateTwoPointsDistanceInMiles(p1lan,p1lng,p2lan,p2lng float64) float64 {
	p := geo.NewPoint(34.117671, -118.073250)
	p2 := geo.NewPoint(34.114826, -118.075295)
	// find the great circle distance between them in km
	dist := p.GreatCircleDistance(p2)
	// change the km to miles
	miles:= dist/0.621371
	return miles
}