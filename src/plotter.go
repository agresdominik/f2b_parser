package main

import (
	"fmt"

	"gonum.org/v1/plot"
	"gonum.org/v1/plot/plotter"
	"gonum.org/v1/plot/plotutil"
	"gonum.org/v1/plot/vg"
)


func barChart(action string, variable string, counter map[string]float64) {

	labels := []string{}

	plottingCount := plotter.Values{}
	for key, value := range counter {
		if value > 100 {
			labels = append(labels, key)
			plottingCount = append(plottingCount, value)
		}
	}

	p := plot.New()

	p.Title.Text = fmt.Sprintf("%v by %v", action, variable)
	p.Y.Label.Text = fmt.Sprintf("%v", action)

	w := vg.Points(15)

	barsA, err := plotter.NewBarChart(plottingCount, w)
	if err != nil {
		panic(err)
	}
	//barsA.LineStyle.Width = vg.Length(0)
	barsA.Color = plotutil.Color(0)
	barsA.Offset = 0

	p.Add(barsA)
	p.NominalX(labels...)

	fileName := fmt.Sprintf("%v_%v_barchart.png", action, variable)

	if err := p.Save(12*vg.Inch, 6*vg.Inch, fileName); err != nil {
		panic(err)
	}
}
