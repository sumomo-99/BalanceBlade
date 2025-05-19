package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
)

const (
	screenWidth  = 640
	screenHeight = 480
	barSpeed     = 2
	successMargin = 0.10 // 10% margin for success
)

type Game struct {
	shape       Shape
	bar         Bar
	score       int
	rand        *rand.Rand
	gameState   GameState
	clickHandled bool // Flag to prevent multiple clicks in one frame
}

type GameState int

const (
	GameStatePlaying GameState = iota
	GameStateResult
)

type Shape struct {
	x      int
	y      int
	width  int
	height int
	area   float64
}

type Bar struct {
	position int
	vertical bool
}

func (g *Game) Init() {
	g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize shape (rectangle)
	shapeWidth := 100 + g.rand.Intn(150) // Random width between 100 and 250
	shapeHeight := 100 + g.rand.Intn(150) // Random height between 100 and 250
	g.shape = Shape{
		x:      screenWidth/2 - shapeWidth/2,
		y:      screenHeight/2 - shapeHeight/2,
		width:  shapeWidth,
		height: shapeHeight,
		area:   float64(shapeWidth * shapeHeight),
	}

	// Determine bar orientation (vertical or horizontal) randomly
	vertical := g.rand.Float64() < 0.5
	g.bar = Bar{
		vertical: vertical,
	}

	// Initialize bar position
	if vertical {
		g.bar.position = g.shape.x
	} else {
		g.bar.position = g.shape.y
	}

	g.score = 0
	g.gameState = GameStatePlaying
	g.clickHandled = false
}

func (g *Game) Update() error {
	switch g.gameState {
	case GameStatePlaying:
		// Move the bar
		if g.bar.vertical {
			g.bar.position += barSpeed
			if g.bar.position > g.shape.x+g.shape.width || g.bar.position < g.shape.x {
				barSpeed *= -1
			}
		} else {
			g.bar.position += barSpeed
			if g.bar.position > g.shape.y+g.shape.height || g.bar.position < g.shape.y {
				barSpeed *= -1
			}
		}

		// Handle user input (click)
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) && !g.clickHandled {
			g.clickHandled = true
			// Calculate areas and determine success
			area1, area2 := g.calculateAreas()
			diff := math.Abs(area1 - area2)
			relativeDiff := diff / g.shape.area

			if relativeDiff <= successMargin {
				// Success!
				scoreIncrement := int(math.Round((1 - relativeDiff/successMargin) * 100)) // Higher score for smaller difference
				g.score += scoreIncrement
				fmt.Printf("Success! Score increment: %d\n", scoreIncrement)
			} else {
				// Failure
				fmt.Println("Failure!")
			}
			g.gameState = GameStateResult
		}
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.clickHandled = false
		}

	case GameStateResult:
		// Wait for click to restart
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.Init() // Restart the game
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the shape
	ebitenutil.DrawRect(screen, float64(g.shape.x), float64(g.shape.y), float64(g.shape.width), float64(g.shape.height), color.White)

	// Draw the bar
	barColor := color.RGBA{255, 0, 0, 255} // Red
	if g.bar.vertical {
		ebitenutil.DrawRect(screen, float64(g.bar.position), float64(g.shape.y), 2, float64(g.shape.height), barColor)
	} else {
		ebitenutil.DrawRect(screen, float64(g.shape.x), float64(g.bar.position), float64(g.shape.width), 2, barColor)
	}

	// Display score and game state
	ebitenutil.DebugPrint(screen, fmt.Sprintf("Score: %d\n", g.score))

	switch g.gameState {
	case GameStatePlaying:
		ebitenutil.DebugPrint(screen, "State: Playing\n")
	case GameStateResult:
		ebitenutil.DebugPrint(screen, "State: Result\nClick to restart")
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) calculateAreas() (float64, float64) {
	if g.bar.vertical {
		width1 := g.bar.position - g.shape.x
		width2 := g.shape.x + g.shape.width - g.bar.position
		area1 := float64(width1 * g.shape.height)
		area2 := float64(width2 * g.shape.height)
		return area1, area2
	} else {
		height1 := g.bar.position - g.shape.y
		height2 := g.shape.y + g.shape.height - g.bar.position
		area1 := float64(height1 * g.shape.width)
		area2 := float64(height2 * g.shape.width)
		return area1, area2
	}
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Half Slice Game")

	game := &Game{}
	game.Init()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
