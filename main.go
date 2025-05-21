package main

import (
	"fmt"
	"image/color"
	"log"
	"math"
	"math/rand"
	"os"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text/v2"
)

const (
	screenWidth   = 640
	screenHeight  = 480
	successMargin = 0.10 // 10% margin for success
	fontSize      = 24
)

var barSpeed = 2.0 // Make it a float to avoid type issues later

type Game struct {
	shape        Shape
	bar          Bar
	stageIndex   int
	score        int
	rand         *rand.Rand
	gameState    GameState
	clickHandled bool // Flag to prevent multiple clicks in one frame
	lives        int
	fontFace     *text.GoTextFace
	areaRatio    float64 // To store the calculated area ratio
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
	kind   ShapeKind
}

type ShapeKind int

const (
	Rectangle ShapeKind = iota
	Circle
	Triangle
)

type Bar struct {
	position int
	vertical bool
}

func (g *Game) Init() {
	g.rand = rand.New(rand.NewSource(time.Now().UnixNano()))

	// Initialize shape (rectangle)
	shapeWidth := 100 + g.rand.Intn(150)  // Random width between 100 and 250
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
	g.lives = 3 // Initial number of lives
	g.stageIndex = 0
	g.InitLevel()
}

func (g *Game) InitLevel() {
	// Use stage settings
	stage := Stages[g.stageIndex%len(Stages)] // Cycle through stages
	shapeWidth := stage.ShapeWidth
	shapeHeight := stage.ShapeHeight

	// Randomly choose shape kind
	shapeKind := ShapeKind(g.rand.Intn(3)) // 3 is the number of shape kinds

	g.shape = Shape{
		x:      screenWidth/2 - shapeWidth/2,
		y:      screenHeight/2 - shapeHeight/2,
		width:  shapeWidth,
		height: shapeHeight,
		area:   float64(shapeWidth * shapeHeight),
		kind:   shapeKind,
	}

	g.bar = Bar{
		vertical: stage.BarVertical,
	}

	// Initialize bar position
	if g.bar.vertical {
		g.bar.position = g.shape.x
	} else {
		g.bar.position = g.shape.y
	}

	barSpeed = stage.BarSpeed
}

func (g *Game) Update() error {
	switch g.gameState {
	case GameStatePlaying:
		// Move the bar
		if g.bar.vertical {
			g.bar.position += int(barSpeed)
			if g.bar.position > g.shape.x+g.shape.width || g.bar.position < g.shape.x {
				barSpeed = -barSpeed
			}
		} else {
			g.bar.position += int(barSpeed)
			if g.bar.position > g.shape.y+g.shape.height || g.bar.position < g.shape.y {
				barSpeed = -barSpeed
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
				g.lives-- // Decrement lives on failure
				if g.lives <= 0 {
					g.gameState = GameStateResult // Game over if no lives left
				} else {
					g.gameState = GameStateResult
				}
			}
			// Calculate area ratio
			if area1+area2 != 0 {
				g.areaRatio = math.Min(area1, area2) / math.Max(area1, area2)
			} else {
				g.areaRatio = 0 // Avoid division by zero
			}
			g.gameState = GameStateResult // Transition to result state
		}
		if !inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			g.clickHandled = false
		}

	case GameStateResult:
		// Wait for click to proceed
		if inpututil.IsMouseButtonJustPressed(ebiten.MouseButtonLeft) {
			if g.lives <= 0 {
				g.Init() // Restart the game only if game over
			} else {
				// Automatically proceed to the next stage if lives remain
				g.stageIndex++
				g.InitLevel()
				g.gameState = GameStatePlaying
			}
		}
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	// Draw the shape
	switch g.shape.kind {
	case Rectangle:
		ebitenutil.DrawRect(screen, float64(g.shape.x), float64(g.shape.y), float64(g.shape.width), float64(g.shape.height), color.White)
	case Circle:
		ebitenutil.DrawCircle(screen, float64(g.shape.x+g.shape.width/2), float64(g.shape.y+g.shape.height/2), float64(g.shape.width/2), color.White)
	case Triangle:
		// Define the vertices of the triangle
		x1 := float64(g.shape.x + g.shape.width/2)
		y1 := float64(g.shape.y)
		x2 := float64(g.shape.x)
		y2 := float64(g.shape.y + g.shape.height)
		x3 := float64(g.shape.x + g.shape.width)
		y3 := float64(g.shape.y + g.shape.height)

		// Draw the triangle using ebitenutil.DrawLine
		ebitenutil.DrawLine(screen, x1, y1, x2, y2, color.White)
		ebitenutil.DrawLine(screen, x2, y2, x3, y3, color.White)
		ebitenutil.DrawLine(screen, x3, y3, x1, y1, color.White)
	}

	// Draw the bar
	barColor := color.RGBA{255, 0, 0, 255}
	if g.bar.vertical {
		ebitenutil.DrawRect(screen, float64(g.bar.position), float64(g.shape.y), 2, float64(g.shape.height), barColor)
	} else {
		ebitenutil.DrawRect(screen, float64(g.shape.x), float64(g.bar.position), float64(g.shape.width), 2, barColor)
	}

	// Display score and game state
	stateText := fmt.Sprintf("Score: %d\nLives: %d\n", g.score, g.lives)
	textop := &text.DrawOptions{}
	textop.LineSpacing = fontSize * 1.5

	switch g.gameState {
	case GameStatePlaying:
		stateText += "State: Playing\n"
	case GameStateResult:
		if g.lives <= 0 {
			stateText += "State: Game Over\nClick to restart"
		} else {
			stateText += fmt.Sprintf("Area Ratio: %.2f\nClick to continue", g.areaRatio)
		}
	}

	text.Draw(screen, stateText, g.fontFace, textop)
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) calculateAreas() (float64, float64) {
	switch g.shape.kind {
	case Rectangle:
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
	case Circle:
		// Approximate circle area division (more complex calculation needed for accuracy)
		radius := float64(g.shape.width / 2) // Assuming width is diameter
		circleArea := math.Pi * radius * radius
		if g.bar.vertical {
			//Approximation, needs proper calculation
			x := float64(g.bar.position)
			area1 := calculateCircularSegmentArea(radius, x - float64(g.shape.x) - radius)
			area2 := circleArea - area1
			return area1, area2
		} else {
			y := float64(g.bar.position)
			area1 := calculateCircularSegmentArea(radius, y - float64(g.shape.y) - radius)
			area2 := circleArea - area1
			return area1, area2
		}

	case Triangle:
		// Approximation for triangle (basic split) - needs proper calculation based on bar orientation
		if g.bar.vertical {
			width1 := g.bar.position - g.shape.x
			width2 := g.shape.x + g.shape.width - g.bar.position
			area1 := 0.5 * float64(width1 * g.shape.height) // Approximation
			area2 := 0.5 * float64(width2 * g.shape.height) // Approximation
			return area1, area2
		} else {
			height1 := g.bar.position - g.shape.y
			height2 := g.shape.y + g.shape.height - g.bar.position
			area1 := 0.5 * float64(height1 * g.shape.width) // Approximation
			area2 := 0.5 * float64(height2 * g.shape.width) // Approximation
			return area1, area2
		}
	default:
		return 0, 0 // Default case
	}
}

// Calculate the area of a circular segment
func calculateCircularSegmentArea(radius, height float64) float64 {
	if height > radius || height < -radius {
		return 0 // Height out of bounds
	}
	theta := 2 * math.Acos((radius - height) / radius)
	area := (radius*radius)/2 * (theta - math.Sin(theta))
	return area
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Balance Blade")

	f, err := os.Open("font/Roboto-Medium.ttf")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	src, err := text.NewGoTextFaceSource(f)
	if err != nil {
		log.Fatal(err)
	}

	game := &Game{
		fontFace: &text.GoTextFace{Source: src, Size: fontSize},
	}
	game.Init()

	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
