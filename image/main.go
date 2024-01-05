package main

import (
	"log"
	"net/http"

	svg "github.com/ajstarks/svgo"
)

func main() {
	http.Handle("/circle", http.HandlerFunc(circle))
	err := http.ListenAndServe(":2003", nil)
	if err != nil {
		log.Fatal("ListenAndServe:", err)
	}

	//http://localhost:2003/circle 
}

func circle(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "image/svg+xml")
	canvas := svg.New(w)
	data := []struct {
		Month string
		Usage int
	}{
		{"Jan", 171},
		{"Feb", 180},
		{"Mar", 100},
		{"Apr", 87},
		{"May", 65},
		{"Jun", 40},
		{"Jul", 32},
		{"Aug", 55},
		{"Sep", 222},
		{"Oct", 0},
		{"Nov", 0},
		{"Dec", 0},
	}
	width := len(data)*60 + 10
	height := 300
	threshold := 160
	max := 0

	// グラフはなれている 数値いじるだけ
	for _, item := range data {
		if item.Usage > max {
			max = item.Usage
		}
	}
	canvas.Start(width, height)
	canvas.Circle(0, 0, width, "fill:white")
	for i, val := range data {
		percent := val.Usage * (height - 50) / max
		canvas.Rect(i*60+10, (height-50)-percent, 50, percent, "fill:rgb(77,20,232)")
		canvas.Text(i*60+35, height-24, val.Month, "font-size:14pt;fill:rgb(150, 150, 150);text-anchor:middle")
	}
	threshPercent := threshold * (height - 50) / max
	canvas.Line(0, height-threshPercent, width, height-threshPercent, "stroke: rgb(255,100,100); opacity: 0.8; stroke-width: 2px")
	canvas.Rect(0, 0, width, height-threshPercent, "fill:rgb(255, 100, 100); opacity: 0.1")
	canvas.Line(0, height-50, width, height-50, "stroke: rgb(150, 150, 150); stroke-width:2")

	canvas.End()
}
