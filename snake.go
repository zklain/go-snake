package main

import (
	"github.com/pkg/profile"
	"time"
	"image/color"
	"log"
	"fmt"
	"math/rand"

	"github.com/hajimehoshi/ebiten/ebitenutil"
	"github.com/hajimehoshi/ebiten/inpututil"
	"github.com/hajimehoshi/ebiten"
	// "github.com/hajimehoshi/ebiten/text"
)

type Coordinates struct {
	x int
	y int
}

type Snake struct {
	x int
	y int
	body []Coordinates
	length int
	direction int
}

type Food struct {
	coordinates Coordinates
	eaten bool
}

// TODO: rename to Game
type World struct {
	width int
	height int
	snake *Snake
	score int
	play bool
	food *Food
}

func init() {
	rand.Seed(time.Now().UnixNano())

	// pixels := make([]uint8, 10 * 10)
	// for i := range pixels { 
	// 	pixels[i] = snakeColor
	// }

	// snakeBody, _ = ebiten.NewImageFromImage(&image.Alpha{
	// 	Pix:    pixels,
	// 	Stride: 10,
	// 	Rect:   image.Rect(0, 0, 10, 10),
	// }, ebiten.FilterDefault)

	snakeBody, _ = ebiten.NewImage(10, 10, ebiten.FilterDefault)
	snakeBody.Fill(color.RGBA{36, 180, 129, 255})

	foodImage, _ = ebiten.NewImage(10, 10, ebiten.FilterDefault)
	foodImage.Fill(color.RGBA{224, 36, 39, 0xff})

	canvasImage, _ = ebiten.NewImage(screenWidth, screenHeight, ebiten.FilterDefault)
	canvasImage.Fill(color.RGBA{36, 224, 127, 255})
}

func GenerateWorld(width, height int) *World {
	// create snake
	s := &Snake{
		x: width / 2,
		y: height / 2,
		length: snakeInitLen,
		direction: 0, 
		body: make([]Coordinates, snakeInitLen),
	}

	for i := range s.body {
		sc := &Coordinates{s.x, s.y + ((i + 1) * 10)}
		s.body[i] = *sc
	}
	
	// TODO: place food
	// TODO: find random till not one of snake coords
	fc := &Coordinates{rand.Intn(width / 10) * 10, rand.Intn(height / 10) * 10}
	f := &Food {
		coordinates: *fc,
		eaten: false,
	}

	w := &World{
		width: width,
		height: height,
		snake: s,
		score: 0,
		play: true,
		food: f,
	}

	return w
}

func (w* World) DrawSnake(image* ebiten.Image) {
	
	// TODO: mege with draw food
	for _, coords := range w.snake.body {
		op := &ebiten.DrawImageOptions{}
		op.GeoM.Translate(float64(coords.x), float64(coords.y))
		image.DrawImage(snakeBody, op)
	}
}

func (w* World) DrawFood(image* ebiten.Image) {
	f := w.food
	op := &ebiten.DrawImageOptions{}
	op.GeoM.Translate(float64(f.coordinates.x), float64(f.coordinates.y))
	image.DrawImage(foodImage, op)
}

func (w* World) MoveSnake(direction int) {
	s := w.snake
	f := w.food
	head := s.body[0]
	tail := s.body[s.length - 1]
	origTail := tail
	switch direction {
	// move up
	case 0:
		if s.direction == 1 {
			return
		}
		if head.y <= 0 {
			tail.x = head.x
			tail.y = screenHeight
		} else {
			tail.x = head.x
			tail.y = head.y - 10
		}
	// DOWN
	case 1:
		if s.direction == 0 {
			return
		}
		if head.y == w.height {
			tail.y = 0
			tail.x = head.x
		} else {
			tail.y = head.y + 10
			tail.x = head.x
		}

	// LEFT
	case 2:
		if s.direction == 3 {
			return
		}
		if head.x == 0 {
			tail.x = w.width
			tail.y = head.y
		} else {
			tail.x = head.x - moveBy
			tail.y = head.y
		}
	// RIGHT
	case 3:
		if s.direction == 2 {
			return
		}
		if head.x == w.width {
			tail.x = 0
			tail.y = head.y
		} else {
			tail.x = head.x + moveBy
			tail.y = head.y
		}
	}
	if s.hasBiten() {
		return
	}
	s.direction = direction
	s.updateBody(tail)
	if head.x == f.coordinates.x && head.y == f.coordinates.y {
		w.food.eaten = true
		// TODO: dont remove in updateBody
		s.eat(origTail)
		w.score++
	}


}

func (w* World) PlaceFood() {
	f := w.food
	f.coordinates.x = rand.Intn(w.width / 10) * 10
	f.coordinates.y = rand.Intn(w.height / 10) * 10
	f.eaten = false
}

func (s* Snake) updateBody(newHead Coordinates) {
	newBody := make([]Coordinates, 0, s.length)
	newBody = append(newBody, newHead)
	for _, v := range s.body[:len(s.body) - 1] {
		newBody = append(newBody, v)
	}
	s.body = newBody
}

func (s* Snake) eat(tail Coordinates) {
	s.body = append(s.body, tail)
	s.length++
}

func (s* Snake) hasBiten() bool {
	head := s.body[0]
	for _, v := range s.body[1:] {
		if head.x == v.x && head.y == v.y {
			return true
		}
	}
	return false
}

const (
	snakeInitLen = 5
	screenWidth = 640
	screenHeight = 480
	moveBy = 10
)

var (
	// pixels = make([]byte, screenWidth * screenHeight * 4)
	// TODO: Enum
	keys = []ebiten.Key{
		ebiten.KeyW, // up
		ebiten.KeyS, // down
		ebiten.KeyA, // left
		ebiten.KeyD, // right
	}
	canvasImage *ebiten.Image
	snakeBody *ebiten.Image
	foodImage *ebiten.Image
	world = GenerateWorld(screenWidth, screenHeight)
)

func update(screen* ebiten.Image) error {
	for i, key := range keys {
		if inpututil.IsKeyJustPressed(key) {
			world.MoveSnake(i)
		} else {
			world.MoveSnake(world.snake.direction)
		}
	}

	if ebiten.IsDrawingSkipped() {
		return nil
	}

	world.DrawSnake(screen)
	if world.food.eaten {
		world.PlaceFood()
	}
	world.DrawFood(screen)


	// print game info
	ebitenutil.DebugPrint(screen, fmt.Sprintf("TPS: %0.2f", ebiten.CurrentTPS()))
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("X: %d Y: %d", world.snake.body[0].x, world.snake.body[0].y), 0, 20)
	ebitenutil.DebugPrintAt(screen, fmt.Sprintf("Score: %d", world.score), 0, 40)

	return nil
}


func main() {
	defer profile.Start(profile.MemProfile).Stop()
	ebiten.SetMaxTPS(10)
	if err := ebiten.Run(update, screenWidth, screenHeight, 2, "Snake"); err != nil {
		log.Fatal(err)
	}
}

// TODO: git
// TODO: Play pause
// TODO: restart
// TODO: safe food placement
// TODO: start screen
// todo: speed (TPS?)