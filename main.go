package main

import (
	"fmt"
	"log"

	"github.com/hajimehoshi/ebiten/v2"
)

const (
	screenWidth  = 640
	screenHeight = 480
	tileSize     = 5
)

type Game struct {
	Snake    *Snake
	Food     *Food
	GameOver bool
	Tick     int
	Speed    int
}

func (g *Game) String() string {
	return fmt.Sprintf("Snake: %v\nFood: %v\nGameState: %v\n", g.Snake, g.Food, g.GameOver)
}

func main() {
	game := &Game{
		Snake:    NewSnake(),
		Food:     &Food{},
		GameOver: false,
		Tick:     0,
		Speed:    1,
	}
	fmt.Println(game.Snake)

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Goo Game")
	if err := ebiten.RunGame(game); err != nil {
		log.Fatal(err)
	}
}
