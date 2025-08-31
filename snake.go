package main

import (
	"fmt"
	"image/color"
	"math/rand/v2"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/inpututil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"github.com/hajimehoshi/ebiten/v2/vector"
	"golang.org/x/image/font/basicfont"
)

type Point struct {
	X int
	Y int
}

type Snake struct {
	Body        []Point
	Head        Point
	Direction   Point
	GrowCounter int
}

type Food struct {
	ID       int
	Location Point
}

var (
	Up    = Point{X: 0, Y: -1}
	Down  = Point{X: 0, Y: 1}
	Left  = Point{X: -1, Y: 0}
	Right = Point{X: 1, Y: 0}
)

func NewSnake() *Snake {
	return &Snake{
		Body: []Point{
			{X: 10, Y: 10},
		},
		Direction:   Point{X: 1, Y: 0},
		GrowCounter: 0,
	}
}

func (s *Snake) Move() {
	if len(s.Body) == 0 {
		return
	}

	head := s.Body[0]
	newHead := Point{
		X: head.X + s.Direction.X,
		Y: head.Y + s.Direction.Y,
	}

	if s.IsValidMove(newHead) {
		s.Head = newHead
		s.Body = append([]Point{newHead}, s.Body...)

		if s.GrowCounter > 0 {
			s.GrowCounter--
		} else {
			s.Body = s.Body[:len(s.Body)-1]
		}
	}
}

var (
	font         = basicfont.Face7x13
	edgePadding  = 20
	foodTileSize = 1
	foodColor    = color.RGBA{0, 255, 0, 128}
)

func (g *Game) Draw(screen *ebiten.Image) {
	if g.GameOver {
		screen.Fill(color.RGBA{255, 0, 0, 1})
		textWithOffset := screenWidth/2 - 100
		text.Draw(screen, "You Lost! Press \"R\" to restart.", font, textWithOffset, edgePadding, color.Black)
	}

	if !g.GameOver {
		screen.Fill(color.RGBA{255, 255, 1, 1})
		textWithOffset := screenWidth/2 - 40
		text.Draw(screen, "Game Started!", font, textWithOffset, edgePadding, color.Black)
	}

	if g.Food.ID != 0 {
		fX := float32(g.Food.Location.X * tileSize)
		fY := float32(g.Food.Location.Y * tileSize)
		// vector.DrawFilledCircle(screen, fX, fY, float32(tileSize*foodTileSize), foodColor, false)
		vector.DrawFilledRect(screen, fX, fY, float32(tileSize*foodTileSize), float32(tileSize*foodTileSize), foodColor, false)
	}

	for _, s := range g.Snake.Body {
		vector.DrawFilledRect(screen, float32(s.X*tileSize), float32(s.Y*tileSize), tileSize, tileSize, color.Black, false)
	}
}

func (g *Game) Layout(outsideWidth int, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func (g *Game) Update() error {
	head := g.Snake.Body[0]
	if head.X >= screenWidth/tileSize-1 || head.Y >= screenHeight/tileSize-1 || head.X < 0+1 {
		g.GameOver = true
	}

	if g.GameOver && inpututil.IsKeyJustPressed(ebiten.KeyR) {
		g.Restart()
	}

	if !g.GameOver {
		g.Snake.Move()

		if g.Food.ID == 0 && g.Food.Location.X == 0 && g.Food.Location.Y == 0 {
			g.NewFood()
		}

		if inpututil.IsKeyJustPressed(ebiten.KeyArrowUp) {
			g.Snake.Direction = Up
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowDown) {
			g.Snake.Direction = Down
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowLeft) {
			g.Snake.Direction = Left
		}
		if inpututil.IsKeyJustPressed(ebiten.KeyArrowRight) {
			g.Snake.Direction = Right
		}

		// Eat food
		if g.AteFood() {
			fmt.Println("ate!!")
			g.Snake.GrowCounter += 1
			g.NewFood()
		}

	}

	return nil
}

func (g *Game) Restart() {
	g.Snake = NewSnake()
	g.GameOver = false
}

/*
Need to check if snake is going backwards onto its body.
Haven't validated whether this works.
*/
func (s *Snake) IsValidMove(next Point) bool {
	for _, v := range s.Body {
		if v.X == next.X && v.Y == next.Y {
			return false
		}
	}
	return true
}

func (g *Game) NewFood() {
	locX := rand.IntN(screenWidth / tileSize)
	locY := rand.IntN(screenHeight / tileSize)

	fmt.Printf("Food Location: (X: %v, Y: %v )", locX, locY)

	var newID int
	if g.Food.ID > 0 {
		newID = g.Food.ID + 1
	} else {
		newID = 1
	}

	g.Food.ID = newID
	g.Food.Location = Point{X: locX, Y: locY}
}

func (g *Game) AteFood() bool {
	if g.Snake.Head.X == g.Food.Location.X && g.Snake.Head.Y == g.Food.Location.Y {
		return true
	}

	return false
}
