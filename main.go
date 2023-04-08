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

func convertWGStoCart(points []g.DDDPoint) []g.DDDPoint {
	newpoints := make([]g.DDDPoint, len(points))
	for i, el := range points {
		x, y := mercator.LatLonToMeters(el.X, el.Y)
		newpoints[i] = g.DDDPoint{x, y, el.Z}

	}
	return newpoints
}
func convertCarttoWGS(points []g.DDDPoint) []g.DDDPoint {
	newpoints := make([]g.DDDPoint, len(points))
	for _, i := range points {
		x, y := mercator.MetersToLatLon(i.X, i.Y)
		newpoints = append(newpoints, g.DDDPoint{x, y, i.Z})

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
	pts := []g.Point{g.Point{22.782258064516125, 29.54545454545455}, g.Point{47.983870967741936, 68.77705627705629}, g.Point{82.86290322580645, 50.10822510822511}, g.Point{65.52419354838709, 9.523809523809526}}
	cv := g.NewCoverage(pts, 1)
	cv.CreateBigLinesSlice()
	eq := cv.CreateCoverageEquations()
	ic := g.CreateInsideCoors(eq)
	points := cv.PreparePointsSlice(ic, eq)

	err := plt(PointtoXY(points))
	if err != nil {
		fmt.Sprintf(err.Error())
	}
}
