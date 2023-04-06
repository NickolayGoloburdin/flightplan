package geometry

import "math"

type Point struct {
	x float64
	y float64
}
type Line struct {
	point1 Point
	point2 Point
}
type Coverage struct {
	points                  []Point
	firstPoint, SecondPoint int
	a, b, startC, stopC     float64
}

func NewCoverage(cpoints []Point) *Coverage {
	return &Coverage{points: cpoints}
}
func dist(a, b Point) float64 {
	return math.Sqrt(math.Pow(b.x-a.x, 2) + math.Pow(b.y-a.y, 2))
}
func lineToPointDist(point Point, a, b, c float64) float64 {
	return a*point.x + b*point.y + c/(math.Sqrt(a*a+b*b))
}
func (c *Coverage) CalcMaxLenghtAndNum(points []Point) {
	var maxLen float64 = 0
	var corners [2]int
	lenght := len(points)
	for k, _ := range points {
		if d := dist(points[k], points[(k+1)%lenght]); d > maxLen {
			maxLen = d
			corners[0] = k
			corners[1] = (k + 1) % lenght
		}

	}
	c.firstPoint, c.SecondPoint = corners
}
func (c *Coverage) CalcAxByCBase() {
	if IsSlopeInf(c.points[c.firstPoint], c.points[c.SecondPoint]) {
		c.a = 1
		c.b = 0
		c.startC = -c.points[c.firstPoint].x
	} else {
		c.a = GetSlope(c.points[c.firstPoint], c.points[c.SecondPoint])
		c.b = -1.0
		c.startC = -c.a*c.points[c.SecondPoint].x + c.points[c.SecondPoint].y
	}
}
func IsSlopeInf(a, b Point) bool {
	return (b.x - a.x) == 0
}
func GetSlope(a, b Point) float64 {
	return (b.y - a.y) / (b.x - a.x)
}
func findFarthestPoint(points []Point, corners [2]int) (float64, int) {
	var farthestpoint int
	var maxDist float64
	for i := range points {
		if i != corners[0] && i != corners[1] {
			dist := lineToPointDist()
		}
	}
}
