package main

import (
	"fmt"
	"image"
	"image/color"
	"log"
	"math"
	"os"

	"gioui.org/app"
	"gioui.org/io/system"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/unit"
)

type Cell struct {
	volume, velX, velY, density float64
	pos                         image.Point
}

var cellGrid []*Cell

func main() {
	go func() {
		w := app.NewWindow(
			app.Title("GF Simulator"),
			app.Size(unit.Dp(640), unit.Dp(640)),
		)
		if err := draw(w); err != nil {
			log.Fatal(err)
		}

		os.Exit(0)
	}()
	app.Main()
}

func drawCell(gtx layout.Context, ops *op.Ops, i int) {
	if cellGrid[i].volume > 0 {
		defer clip.Rect{Min: image.Pt(cellGrid[i].pos.X, cellGrid[i].pos.Y), Max: image.Pt(cellGrid[i].pos.X+gtx.Constraints.Max.X/200, cellGrid[i].pos.Y+gtx.Constraints.Max.Y/200)}.Push(ops).Pop()
		paint.ColorOp{Color: color.NRGBA{B: 0x80, A: uint8(math.Round(0xFF * cellGrid[i].volume))}}.Add(ops)
		paint.PaintOp{}.Add(ops)
	}
}

func draw(w *app.Window) error {
	//var for operations from UI
	var ops op.Ops

	for e := range w.Events() {
		switch e := e.(type) {

		// this is sent when the application should re-render.
		case system.FrameEvent:
			gtx := layout.NewContext(&ops, e)
			// Let's try out the flexbox layout concept:
			// layout.Flex{
			// 	// Vertical alignment, from top to bottom
			// 	Axis: layout.Vertical,
			// 	// Empty space is left at the start, i.e. at the top
			// 	Spacing: layout.SpaceStart,
			// }.Layout(gtx,
			// 	layout.Rigid(
			// 		func(gtx layout.Context) layout.Dimensions {
			// 			circle := clip.Ellipse{
			// 				Min: image.Pt(80, 0),
			// 				Max: image.Pt(320, 240),
			// 			}.Op(gtx.Ops)
			// 			color := color.NRGBA{B: 200, A: 255}
			// 			paint.FillShape(gtx.Ops, color, circle)
			// 			d := image.Point{Y: 400}
			// 			return layout.Dimensions{Size: d}
			// 		},
			// 	),
			// )
			cellGrid = append(cellGrid, new(Cell))
			cellGrid[0].volume = 1
			cellGrid[0].pos = image.Pt(3, 3)
			fmt.Print(gtx.Constraints.Max)
			for i := range cellGrid {
				drawCell(gtx, gtx.Ops, i)
			}
			e.Frame(gtx.Ops)

		case system.DestroyEvent:
			return e.Err
		}

	}
	return nil
}
