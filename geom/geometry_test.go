package geometry

import (
	"testing"
)

func TestIsInsidePolygon(t *testing.T) {
	poly := []Point{Point{0, 0}, Point{10, 0}, Point{10, 10.011234443211}, Point{0, 10}}
	tesPoint := Point{10, 10.011234443211}
	if IsInsidePolygon(poly, tesPoint) != true {
		t.Errorf("ERROR")
	}
}
