//go:build js && wasm

package main

import (
	"syscall/js"
)

type Game struct {
	Canvas       js.Value
	Ctx          js.Value
	Width        float64
	Height       float64
	Ball         *Ball
	Paddle       *Paddle
	Bricks       [][]*Brick
	Rows         int
	Cols         int
	LeftPressed  bool
	RightPressed bool
}

func main() {
	// Block forever after wiring up the game — keeps Go runtime alive in WASM
	setupAndStart()
	select {}
}

// TODO
// create start game button and reset button
// clearing blocks doesn't give a "win" screen.
// ball should speed up after x number of bounces
// brick layers should be different colors
//

func setupAndStart() {
	doc := js.Global().Get("document")
	canvas := doc.Call("getElementById", "breakoutCanvas")
	if canvas.IsNull() || canvas.IsUndefined() {
		// nothing to do if canvas not present
		return
	}

	width := 800.0
	height := 600.0
	canvas.Set("width", width)
	canvas.Set("height", height)
	ctx := canvas.Call("getContext", "2d")

	game := &Game{
		Canvas: canvas,
		Ctx:    ctx,
		Width:  width,
		Height: height,
		Ball:   &Ball{X: width / 2, Y: height / 2, DX: 2.5, DY: -2.5, R: 8.0},
		Paddle: &Paddle{X: (width - 100.0) / 2, Y: height - 10.0 - 10.0, Width: 100.0, Height: 10.0, Speed: 6.0},
		Rows:   5,
		Cols:   8,
	}

	game.initBricks()
	game.setupKeyListeners()
	game.startGameLoop()
}

func (g *Game) initBricks() {
	brickW := 75.0
	brickH := 20.0
	brickP := 10.0
	offsetTop := 30.0
	offsetLeft := 30.0

	g.Bricks = make([][]*Brick, g.Rows)
	for i := 0; i < g.Rows; i++ {
		g.Bricks[i] = make([]*Brick, g.Cols)
		for j := 0; j < g.Cols; j++ {
			bx := offsetLeft + float64(j)*(brickW+brickP)
			by := offsetTop + float64(i)*(brickH+brickP)
			g.Bricks[i][j] = &Brick{
				Active: true,
				X:      bx,
				Y:      by,
				Width:  brickW,
				Height: brickH,
			}
		}
	}
}

func (g *Game) setupKeyListeners() {
	keyDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ev := args[0]
		k := ev.Get("key").String()
		if k == "ArrowLeft" || k == "Left" {
			g.LeftPressed = true
		}
		if k == "ArrowRight" || k == "Right" {
			g.RightPressed = true
		}
		return nil
	})
	keyUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ev := args[0]
		k := ev.Get("key").String()
		if k == "ArrowLeft" || k == "Left" {
			g.LeftPressed = false
		}
		if k == "ArrowRight" || k == "Right" {
			g.RightPressed = false
		}
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "keydown", keyDown)
	js.Global().Get("document").Call("addEventListener", "keyup", keyUp)
}

func (g *Game) startGameLoop() {
	var raf js.Func
	raf = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		g.update()
		g.draw()
		js.Global().Call("requestAnimationFrame", raf)
		return nil
	})
	js.Global().Call("requestAnimationFrame", raf)
}

func (g *Game) update() {
	g.Paddle.Update(g.LeftPressed, g.RightPressed, g.Width)
	g.Ball.Update(g.Width, g.Height, g.Paddle, g.Bricks)
}

func (g *Game) draw() {
	g.Ctx.Call("clearRect", 0, 0, g.Width, g.Height)
	g.Ball.Draw(g.Ctx)
	g.Paddle.Draw(g.Ctx, g.Height)
	g.drawBricks()
}

func (g *Game) drawBricks() {
	for i := 0; i < g.Rows; i++ {
		for j := 0; j < g.Cols; j++ {
			if g.Bricks[i][j].Active {
				g.Bricks[i][j].Draw(g.Ctx)
			}
		}
	}
}

func (b *Ball) Draw(ctx js.Value) {
	ctx.Call("beginPath")
	ctx.Call("arc", b.X, b.Y, b.R, 0, 2*3.14159)
	ctx.Set("fillStyle", "#0095DD")
	ctx.Call("fill")
	ctx.Call("closePath")
}

func (p *Paddle) Draw(ctx js.Value, height float64) {
	ctx.Call("beginPath")
	ctx.Call("rect", p.X, p.Y, p.Width, p.Height)
	ctx.Set("fillStyle", "#0095DD")
	ctx.Call("fill")
	ctx.Call("closePath")
}

func (b *Brick) Draw(ctx js.Value) {
	if b.Active {
		ctx.Call("beginPath")
		ctx.Call("rect", b.X, b.Y, b.Width, b.Height)
		ctx.Set("fillStyle", "#FF5733")
		ctx.Call("fill")
		ctx.Call("closePath")
	}
}
