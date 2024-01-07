package main

import (
	"fmt"
	"image/color"

	"github.com/jung-kurt/gofpdf"
)

type PDFOption func(*gofpdf.Fpdf)

func FillColor(c color.RGBA) PDFOption {
	return func(pdf *gofpdf.Fpdf) {
		r, g, b := rgb(c)
		pdf.SetFillColor(r, g, b)
	}
}

func rgb(c color.RGBA) (int, int, int) {
	alpha := float64(c.A) / 255.0
	alphaWhite := int(255 * (1.0 - alpha))
	r := int(float64(c.R)*alpha) + alphaWhite
	g := int(float64(c.G)*alpha) + alphaWhite
	b := int(float64(c.B)*alpha) + alphaWhite
	return r, g, b
}

type PDF struct {
	fpdf *gofpdf.Fpdf
	x, y float64
}

func (p *PDF) Move(xDelta, yDelta float64) {
	p.x, p.y = p.x+xDelta, p.y+yDelta
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) MoveAbs(x, y float64) {
	p.x, p.y = x, y
	p.fpdf.MoveTo(p.x, p.y)
}

func (p *PDF) Text(text string) {
	p.fpdf.Text(p.x, p.y, text)
}

func (p *PDF) Polygon(pts []gofpdf.PointType, opts ...PDFOption) {
	for _, opt := range opts {
		opt(p.fpdf)
	}
	p.fpdf.Polygon(pts, "F")
}

func main() {
	pdf := gofpdf.New("P", "mm", "A4", "")
	pdf.AddPage()
	pdf.SetFont("Arial", "B", 16)
	pdf.Cell(40, 10, "ドキュメントベースでいじってみる")
	err := pdf.OutputFileAndClose("hello.pdf")
	if err != nil {
		panic(err)
	}
}

func drawGrid(pdf *gofpdf.Fpdf) {
	w, h := pdf.GetPageSize()
	pdf.SetFont("courier", "", 12)
	pdf.SetTextColor(80, 80, 80)
	pdf.SetDrawColor(200, 200, 200)
	for x := 0.0; x < w; x = x + (w / 20.0) {
		pdf.SetTextColor(200, 200, 200)
		pdf.Line(x, 0, x, h)
		_, lineHt := pdf.GetFontSize()
		pdf.Text(x, lineHt, fmt.Sprintf("%d", int(x)))
	}
	for y := 0.0; y < h; y = y + (w / 20.0) {
		pdf.SetTextColor(80, 80, 80)
		pdf.Line(0, y, w, y)
		pdf.Text(0, y, fmt.Sprintf("%d", int(y)))
	}
}
