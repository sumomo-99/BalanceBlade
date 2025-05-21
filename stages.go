package main

type Stage struct {
	ShapeWidth  int
	ShapeHeight int
	BarVertical bool
	BarSpeed    float64
}

var Stages = []Stage{
	{ShapeWidth: 200, ShapeHeight: 150, BarVertical: true, BarSpeed: 2.0},   // Stage 1
	{ShapeWidth: 150, ShapeHeight: 200, BarVertical: false, BarSpeed: 2.5},  // Stage 2
	{ShapeWidth: 250, ShapeHeight: 100, BarVertical: true, BarSpeed: -3.0},  // Stage 3
	{ShapeWidth: 180, ShapeHeight: 180, BarVertical: false, BarSpeed: -2.0},  // Stage 4
	{ShapeWidth: 220, ShapeHeight: 130, BarVertical: true, BarSpeed: 3.5},   // Stage 5
}
