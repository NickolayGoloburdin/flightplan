package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	"github.com/davvo/mercator"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

type Point struct {
	x float64
	y float64
	z float64
}

func convertWGStoCart(points []Point) []Point {
	newpoints := make([]Point, len(points))
	for i, el := range points {
		x, y := mercator.LatLonToMeters(el.x, el.y)
		newpoints[i] = Point{x, y, el.z}

	}
	return newpoints
}
func convertCarttoWGS(points []Point) []Point {
	newpoints := make([]Point, len(points))
	for _, i := range points {
		x, y := mercator.MetersToLatLon(i.x, i.y)
		newpoints = append(newpoints, Point{x, y, i.z})

	}
	return newpoints
}
func PointtoXY(points []Point) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i := range pts {
		pts[i].X = points[i].x
		pts[i].Y = points[i].y
	}
	return pts
}
func plt(pts plotter.XYs) error {
	p := plot.New()
	p.Title.Text = "Title"
	p.X.Label.Text = "X [mm]"
	p.Y.Label.Text = "Y [A.U.]"
	p.X.Label.Position = draw.PosRight
	p.Y.Label.Position = draw.PosTop
	s, err := plotter.NewScatter(pts)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	s.GlyphStyle.Shape = draw.CrossGlyph{}
	s.Color = color.RGBA{R: 255, A: 255}
	p.Add(s)
	wt, err := p.WriterTo(512, 512, "png")
	if err != nil {
		return fmt.Errorf(err.Error())

	}
	f, err := os.Create("out.png")
	defer f.Close()
	_, err = wt.WriteTo(f)
	if err != nil {
		return fmt.Errorf(err.Error())
	}
	return nil
}

func Dist(a, b Point) float64 {
	return math.Sqrt(math.Pow(b.x-a.x, 2) + math.Pow(b.y-a.y, 2))
}

func main() {
	points := []Point{Point{45.91043204152349, 50.50752024855144, 0}, Point{45.90929808523507, 50.50737081218062, 0}, Point{45.90721698424876, 50.50676072184594, 0},
		Point{45.90506986198657, 50.5057890410317, 0}, Point{45.90767218680866, 50.50375778662693, 0}, Point{45.91184648100185, 50.50651414134762, 0}, Point{45.91043204152349, 50.50752024855144, 0}}

	pts := PointtoXY(convertWGStoCart(points))
	err := plt(pts)
	if err != nil {
		fmt.Sprintf(err.Error())
	}
}
