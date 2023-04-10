package geometry

import (
	"fmt"
	"math"

	geo "github.com/kellydunn/golang-geo"
)

const accuracy float64 = 0.00001

type Point struct {
	X float64
	Y float64
}
type DDDPoint struct {
	X float64
	Y float64
	Z float64
}
type LineEquation struct {
	slope, a, b, c float64
}
type BoundingBox struct {
	BottomLeft Point
	TopRight   Point
}
type Coverage struct {
	points                []Point
	bigsidePoints         [2]int
	bigLinesEquationSlice []LineEquation
	SlopeInf              bool
	FLightCoeff           float64
}

func Round(a float64) float64 {
	return math.Round(a*10) / 10
}
func RoundP(a Point) Point {
	return Point{math.Round(a.X*100) / 100, math.Round(a.Y*100) / 100}
}
func AlmostEqual(a, b float64) (res bool) {
	m := math.Abs(math.Abs(a) - math.Abs(b))
	res = m < accuracy
	return
}

func NewCoverage(cpoints []Point, coeff float64) *Coverage {
	return &Coverage{points: cpoints, FLightCoeff: coeff}
}
func Reverse(input []int) []int {
	inputLen := len(input)
	output := make([]int, inputLen)

	for i, n := range input {
		j := inputLen - i - 1

		output[j] = n
	}

	return output
}
func dist(a, b Point) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}
func FarthestPoint(points []Point, le LineEquation) Point {
	var maxlen float64 = 0
	var num int
	if le.slope == math.Inf(1) {
		for i, point := range points {
			dist := math.Abs(-le.c - point.X)
			if dist > maxlen {
				maxlen = dist
				num = i
			}
		}
	} else {
		for i, point := range points {
			dist := lineToPointDist(point, le)
			if dist > maxlen {
				maxlen = dist
				num = i
			}
		}

	}
	return points[num]
}
func lineToPointDist(point Point, le LineEquation) float64 {
	return math.Abs(le.a*point.X+le.b*point.Y+le.c) / (math.Sqrt(le.a*le.a + le.b*le.b))
}
func Intersection(linea, lineb LineEquation) (p Point, parallels bool) {

	if linea.slope == lineb.slope {
		parallels = true
		return
	}
	parallels = false

	p.X = (linea.b*lineb.c - lineb.b*linea.c) / (linea.a*lineb.b - lineb.a*linea.b)
	p.Y = (lineb.a*linea.c - linea.a*lineb.c) / (linea.a*lineb.b - lineb.a*linea.b)
	// if linea.a == 0 && lineb.b == 0 {
	// 	p.X = lineb.c
	// 	p.Y = linea.c
	// } else if lineb.a == 0 && linea.b == 0 {
	// 	p.X = linea.c
	// 	p.Y = lineb.c

	// } else if linea.a == 0 {
	// 	p.Y = linea.c
	// 	p.X = (-p.Y*lineb.b - lineb.c) / lineb.a

	// } else if lineb.a == 0 {
	// 	p.Y = lineb.c
	// 	p.X = (-p.Y*linea.b - linea.c) / linea.a

	// } else if linea.b == 0 {
	// 	p.X = linea.c
	// 	p.Y = (-p.X*linea.b - lineb.c) / lineb.a

	// } else if lineb.b == 0 {
	// 	p.X = lineb.c
	// 	p.Y = (-p.X*lineb.b - linea.c) / linea.a

	// } else {
	// 	p.Y = (linea.a*lineb.c - lineb.a*linea.c) / (lineb.a*linea.b - linea.a*lineb.b)
	// 	p.X = (-p.Y*linea.b - linea.c) / linea.a

	// }
	return
}
func CalcC(point Point, le LineEquation) float64 {
	return -le.a*point.X - le.b*point.Y
}
func CalcSlopeAxByC(fpoint, spoint Point) (le LineEquation) {
	if IsSlopeInf(fpoint, spoint) {
		le.slope = math.Inf(1)
		le.a = 1
		le.b = 0
		le.c = -fpoint.X
	} else {
		le.slope = GetSlope(fpoint, spoint)
		le.a = le.slope
		le.b = -1.0
		le.c = -le.a*spoint.X + spoint.Y
	}
	return
}
func IsSlopeInf(a, b Point) bool {
	return (b.X - a.X) == 0
}
func GetSlope(a, b Point) float64 {
	return (b.Y - a.Y) / (b.X - a.X)
}
func Orientation(p, q, r Point) int {
	val := (((q.Y - p.Y) *
		(r.X - q.X)) -
		((q.X - p.X) *
			(r.Y - q.Y)))
	if Round(val) == 0.0 {
		return 0
	}
	if val > 0.0 {
		return 1 //Collinear
	} else {
		return 2 // Clock or counterclock
	}
}
func OnSegment(p, q, r Point) bool {
	if (q.X <= math.Max(p.X, r.X)) &&
		(q.X >= math.Min(p.X, r.X)) &&
		(q.Y <= math.Max(p.Y, r.Y)) &&
		(q.Y >= math.Min(p.Y, r.Y)) {
		return true
	}

	return false
}
func IsInsidePolygon(poly []Point, point Point) bool {

	var p geo.Polygon
	l := len(poly)
	for i, polpoint := range poly {
		if AlmostEqual(Round(polpoint.X), Round(point.X)) && AlmostEqual(Round(polpoint.Y), Round(point.Y)) {
			return true
		}
		if Orientation(poly[i], point, poly[(i+1)%l]) == 0 {
			return OnSegment(RoundP(poly[i]), RoundP(point), RoundP(poly[(i+1)%l]))
		}

		p.Add(geo.NewPoint(Round(polpoint.X), Round(polpoint.Y)))
	}
	gpoint := geo.NewPoint(Round(point.X), Round(point.Y))

	return p.Contains(gpoint)
}
func (c *Coverage) CalcMaxLenghtNums() {
	var maxLen float64 = 0
	lenght := len(c.points)
	for k, _ := range c.points {
		if d := dist(c.points[k], c.points[(k+1)%lenght]); d > maxLen {
			maxLen = d
			c.bigsidePoints[0] = k
			c.bigsidePoints[1] = (k + 1) % lenght
		}

	}
}
func (cov *Coverage) CreateBigLinesSlice() {
	cov.CalcMaxLenghtNums()
	le := CalcSlopeAxByC(cov.points[cov.bigsidePoints[0]], cov.points[cov.bigsidePoints[1]])
	var farthestlen float64
	var farthestnum int

	if le.slope == math.Inf(1) {
		cov.SlopeInf = true
	}

	for k, point := range cov.points {

		if k != cov.bigsidePoints[0] && k != cov.bigsidePoints[1] {

			dist := lineToPointDist(point, le)
			if dist > farthestlen {
				farthestlen = dist
				farthestnum = k
			}
		}

	}
	stopC := -le.a*cov.points[farthestnum].X - le.b*cov.points[farthestnum].Y

	toAdd := math.Sqrt(le.a*le.a+le.b*le.b) * cov.FLightCoeff

	if cov.SlopeInf {

		to := cov.points[farthestnum].X
		from := cov.points[cov.bigsidePoints[0]].X
		if to > from {
			from += toAdd
			for to > from {
				cov.bigLinesEquationSlice = append(cov.bigLinesEquationSlice, LineEquation{le.slope, le.a, le.b, from})
				from += toAdd
			}
		} else {
			from -= toAdd
			for to < from {
				cov.bigLinesEquationSlice = append(cov.bigLinesEquationSlice, LineEquation{le.slope, le.a, le.b, from})
				from -= toAdd
			}
		}

	} else {

		if stopC > le.c {
			curC := le.c + toAdd
			for curC < stopC {
				cov.bigLinesEquationSlice = append(cov.bigLinesEquationSlice, LineEquation{le.slope, le.a, le.b, curC})
				curC += toAdd
			}
		} else {
			curC := le.c - toAdd
			for curC > stopC {
				cov.bigLinesEquationSlice = append(cov.bigLinesEquationSlice, LineEquation{le.slope, le.a, le.b, curC})
				curC -= toAdd
			}
		}

	}
}
func (cov *Coverage) CreateCoverageEquations() (eqslice []LineEquation) {

	for k, point := range cov.points {

		if IsSlopeInf(cov.points[k], cov.points[(k+1)%len(cov.points)]) {
			line := CalcSlopeAxByC(cov.points[k], cov.points[(k+1)%len(cov.points)])
			farp := FarthestPoint(cov.points, line)
			if farp.X > cov.points[k].X {
				eqslice = append(eqslice, LineEquation{-1, 1, 0, point.X + cov.FLightCoeff})
			} else {
				eqslice = append(eqslice, LineEquation{-1, 1, 0, point.X - cov.FLightCoeff})
			}
		} else {
			line := CalcSlopeAxByC(cov.points[k], cov.points[(k+1)%len(cov.points)])
			farp := FarthestPoint(cov.points, line)
			stopC := CalcC(farp, line)
			toAdd := math.Sqrt(line.a*line.a + line.b*line.b)
			if stopC > line.c {
				line.c += toAdd * cov.FLightCoeff
				eqslice = append(eqslice, line)

			} else {
				line.c -= toAdd * cov.FLightCoeff
				eqslice = append(eqslice, line)
			}
		}
	}
	return
}

func CreateInsideCoors(eqslice []LineEquation) (insideCoors []Point) {
	for i := range eqslice {
		inter, p := Intersection(eqslice[i], eqslice[(i+1)%len(eqslice)])
		if !p {
			insideCoors = append(insideCoors, inter)
		}
	}
	//insideCoors = append(insideCoors, insideCoors[0])
	return
}
func (cov *Coverage) PreparePointsSlice(insideCoors []Point, eqslice []LineEquation) (finalTrajectory []Point) {
	lastVisited := []int{-1, -1}
	l := len(eqslice)
	for _, line := range cov.bigLinesEquationSlice {
		curPos := []int{}
		var tpoint []Point

		for x := range eqslice {
			IntersectionPoint, parallels := Intersection(line, eqslice[x])
			if parallels {
				continue
			} else {
				if IsInsidePolygon(insideCoors, IntersectionPoint) {
					curPos = append(curPos, x)
					tpoint = append(tpoint, IntersectionPoint)
				}
			}
		}

		if len(curPos) < 2 {
			break
		} else if lastVisited[0] == -1 {
			finalTrajectory = append(finalTrajectory, tpoint[0])
			finalTrajectory = append(finalTrajectory, tpoint[1])
			lastVisited[0] = curPos[0]
			lastVisited[1] = curPos[1]

		} else {

			if curPos[0] == lastVisited[0] && curPos[1] == lastVisited[1] {

				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited = Reverse(lastVisited)

			} else if curPos[0] == lastVisited[1] && curPos[1] == lastVisited[0] {

				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited = Reverse(lastVisited)

			} else if lastVisited[0] == curPos[0] {

				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if lastVisited[1] == curPos[1] {

				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if lastVisited[0] == curPos[1] {

				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited[0] = curPos[0]
				lastVisited[1] = curPos[1]

			} else if lastVisited[1] == curPos[0] {

				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited[0] = curPos[0]
				lastVisited[1] = curPos[1]

			} else if lastVisited[0] == l-1 && curPos[0] == 0 {

				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]
			} else if lastVisited[1] == l-1 && curPos[0] == 0 {

				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited[0] = curPos[0]
				lastVisited[1] = curPos[1]
			} else if curPos[1] == l-1 && lastVisited[1] == 0 {

				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if curPos[1]-lastVisited[0] == 1 && curPos[0]-lastVisited[1] == -1 {
				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited[0] = curPos[0]
				lastVisited[1] = curPos[1]

			} else if lastVisited[0] == 0 && curPos[1] == l-1 {

				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if curPos[0]-lastVisited[0] == -1 && curPos[1]-lastVisited[1] == 1 {
				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if curPos[0]-lastVisited[0] == 1 && curPos[1]-lastVisited[1] == -1 {
				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if curPos[0]-lastVisited[0] == 1 && curPos[1]-lastVisited[1] == 1 {
				finalTrajectory = append(finalTrajectory, tpoint[1])
				finalTrajectory = append(finalTrajectory, tpoint[0])
				lastVisited[0] = curPos[1]
				lastVisited[1] = curPos[0]

			} else if curPos[0]-lastVisited[1] == 1 && lastVisited[0]-curPos[1] == 1 {
				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited[0] = curPos[0]
				lastVisited[1] = curPos[1]

			} else if curPos[1]-lastVisited[0] == 1 && curPos[0]-lastVisited[1] == 1 {
				finalTrajectory = append(finalTrajectory, tpoint[0])
				finalTrajectory = append(finalTrajectory, tpoint[1])
				lastVisited[0] = curPos[0]
				lastVisited[1] = curPos[1]

			} else {
				fmt.Errorf("bug")

			}

		}
	}
	return
}
