package main

import (
	"fmt"
	"image/color"
	"math"
	"os"

	g "github.com/Nickolaygoloburdin/flightplanner/geom"
	"github.com/davvo/mercator"
	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/vg/draw"
)

func convertWGStoCart(points []g.Point) []g.Point {
	newpoints := make([]g.Point, len(points))
	for i, el := range points {
		x, y := mercator.LatLonToMeters(el.X, el.Y)
		newpoints[i] = g.Point{x, y}

	}
	return newpoints
}
func convertCarttoWGS(points []g.Point) []g.Point {
	newpoints := make([]g.Point, len(points))
	for _, i := range points {
		x, y := mercator.MetersToLatLon(i.X, i.Y)
		newpoints = append(newpoints, g.Point{x, y})

	}
	return newpoints
}
func PointtoXY(points []g.Point) plotter.XYs {
	pts := make(plotter.XYs, len(points))
	for i := range pts {
		pts[i].X = points[i].X
		pts[i].Y = points[i].Y
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

func Dist(a, b g.Point) float64 {
	return math.Sqrt(math.Pow(b.X-a.X, 2) + math.Pow(b.Y-a.Y, 2))
}

func main() {
	//pts := PointtoXY(convertWGStoCart(points))
	points := []g.Point{g.Point{39.11909179137994, 45.763133084558405}, g.Point{39.12513548767524, 45.7632710483264}, g.Point{39.125447468361614, 45.75795265224616}, g.Point{39.1221870569868, 45.75846919338599},
		g.Point{39.12036648632977, 45.75908922423628}, g.Point{39.11973611839241, 45.75962453390596}, g.Point{39.119626628876716, 45.76015995663019}, g.Point{39.11975016674302, 45.761943564761964},
		g.Point{39.119900557578234, 45.762676699792195}, g.Point{39.120091791379934, 45.763133084558405}}
	pts := convertWGStoCart(points)
	cv := g.NewCoverage(pts, 10)
	cv.CreateBigLinesSlice()
	eq := cv.CreateCoverageEquations()
	ic := g.CreateInsideCoors(eq)
	finpoints := cv.PreparePointsSlice(ic, eq)

	err := plt(PointtoXY((finpoints)))
	if err != nil {
		fmt.Sprintf(err.Error())
	}
	fmt.Printf("Points :", convertCarttoWGS(finpoints))
}
