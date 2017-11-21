package snake

import (
	_ "fmt"
	"time"

	"github.com/nsf/termbox-go"
)

const (
	Speed           = 25 * time.Millisecond
	GrowAmount      = 10
	FoodCount       = 5
	TextColor       = termbox.ColorGreen
	BackgroundColor = termbox.ColorDefault
	SnakeColor      = termbox.ColorGreen
	FoodColor       = termbox.ColorGreen
	Wide            = 50
	High            = 50
)

type Direction int

const (
	Up Direction = iota
	Down
	Left
	Right
)

type Coord struct {
	X, Y int
}

type Snake struct {
	Direction Direction
	Body      []Coord
	Grow      int // Amount left to grow
}

type Context struct {
	Quit       bool
	Wide, High int
	Snake      map[string]Snake
}

func (c *Coord) Draw(color termbox.Attribute) {
	termbox.SetCell(c.X, c.Y, ' ', BackgroundColor, color)
}

func NewSnake(x, y int) Snake {
	snake := Snake{
		Direction: Up,
		Grow:      GrowAmount,
	}
	snake.Push(Coord{x, y})
	return snake
}

func (s *Snake) Draw() {
	for _, c := range s.Body {
		c.Draw(SnakeColor)
	}
}

func (s *Snake) Push(c Coord) {
	s.Body = append(s.Body, c)
}

func (s *Snake) Pop() {
	s.Body = s.Body[1:]
}

func (s *Snake) Head() Coord {
	return s.Body[len(s.Body)-1]
}

func NewContext() *Context {
	return &Context{
		Snake: make(map[string]Snake),
		Wide:  Wide,
		High:  High,
	}
}

func (ctx *Context) Grow(s *Snake) {
	w, h := ctx.Wide, ctx.High
	head := s.Head()
	c := Coord{head.X, head.Y}
	switch s.Direction {
	case Up:
		c.Y--
		if c.Y < 0 {
			c.Y = h - 1
		}
	case Down:
		c.Y++
		if c.Y >= h {
			c.Y = 0
		}
	case Left:
		c.X--
		if c.X < 0 {
			c.X = w - 1
		}
	case Right:
		c.X++
		if c.X >= w {
			c.X = 0
		}
	}
	s.Push(c)
}

func (ctx *Context) Move(s *Snake) {
	ctx.Grow(s)
	if s.Grow <= 0 {
		s.Pop()
	} else {
		s.Grow--
	}
}

func (ctx *Context) Draw() {
	termbox.Clear(BackgroundColor, BackgroundColor)
	for _, snake := range ctx.Snake {
		snake.Draw()
	}
	termbox.Flush()
}

func (ctx *Context) Update() {
	for client, snake := range ctx.Snake {
		ctx.Move(&snake)
		ctx.Snake[client] = snake
	}
}

func (ctx *Context) AddSnake(client string, w, h int) {
	ctx.Snake[client] = NewSnake(w, h)
}

func (ctx *Context) HandleKey(client string, key termbox.Key) {
	switch key {
	case termbox.KeyArrowUp:
		if ctx.Snake[client].Direction != Down {
			snake := ctx.Snake[client]
			snake.Direction = Up
			ctx.Snake[client] = snake
		}
	case termbox.KeyArrowDown:
		if ctx.Snake[client].Direction != Up {
			snake := ctx.Snake[client]
			snake.Direction = Down
			ctx.Snake[client] = snake
		}
	case termbox.KeyArrowLeft:
		if ctx.Snake[client].Direction != Right {
			snake := ctx.Snake[client]
			snake.Direction = Left
			ctx.Snake[client] = snake
		}
	case termbox.KeyArrowRight:
		if ctx.Snake[client].Direction != Left {
			snake := ctx.Snake[client]
			snake.Direction = Right
			ctx.Snake[client] = snake
		}
	}
}
