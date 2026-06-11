package main

import (
	"syscall/js"
)

func main() {
	// Block forever after wiring up the game — keeps Go runtime alive in WASM
	setupAndStart()
	select {}
}

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

	// game state
	paddleWidth := 100.0
	paddleHeight := 10.0
	paddleX := (width - paddleWidth) / 2

	ballX := width / 2
	ballY := height / 2
	ballDX := 2.5
	ballDY := -2.5
	ballR := 8.0

	rows := 5
	cols := 8
	brickW := 75.0
	brickH := 20.0
	brickP := 10.0
	offsetTop := 30.0
	offsetLeft := 30.0

	bricks := make([][]int, rows)
	for i := 0; i < rows; i++ {
		bricks[i] = make([]int, cols)
		for j := 0; j < cols; j++ {
			bricks[i][j] = 1
		}
	}

	left := false
	right := false

	keyDown := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ev := args[0]
		k := ev.Get("key").String()
		if k == "ArrowLeft" || k == "Left" {
			left = true
		}
		if k == "ArrowRight" || k == "Right" {
			right = true
		}
		return nil
	})
	keyUp := js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		ev := args[0]
		k := ev.Get("key").String()
		if k == "ArrowLeft" || k == "Left" {
			left = false
		}
		if k == "ArrowRight" || k == "Right" {
			right = false
		}
		return nil
	})
	js.Global().Get("document").Call("addEventListener", "keydown", keyDown)
	js.Global().Get("document").Call("addEventListener", "keyup", keyUp)

	var raf js.Func
	raf = js.FuncOf(func(this js.Value, args []js.Value) interface{} {
		// update paddle
		if left {
			paddleX -= 6
			if paddleX < 0 {
				paddleX = 0
			}
		}
		if right {
			paddleX += 6
			if paddleX+paddleWidth > width {
				paddleX = width - paddleWidth
			}
		}

		// update ball
		ballX += ballDX
		ballY += ballDY

		if ballX+ballR > width || ballX-ballR < 0 {
			ballDX = -ballDX
		}
		if ballY-ballR < 0 {
			ballDY = -ballDY
		} else if ballY+ballR > height {
			// bottom
			if ballX > paddleX && ballX < paddleX+paddleWidth {
				ballDY = -ballDY
				diff := ballX - (paddleX + paddleWidth/2)
				ballDX = diff / (paddleWidth / 2) * 4
			} else {
				// reset
				ballX = width / 2
				ballY = height / 2
				ballDX = 2.5
				ballDY = -2.5
			}
		}

		// brick collisions
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				if bricks[i][j] == 1 {
					bx := offsetLeft + float64(j)*(brickW+brickP)
					by := offsetTop + float64(i)*(brickH+brickP)
					if ballX > bx && ballX < bx+brickW && ballY > by && ballY < by+brickH {
						ballDY = -ballDY
						bricks[i][j] = 0
					}
				}
			}
		}

		// draw
		ctx.Call("clearRect", 0, 0, width, height)

		// ball
		ctx.Call("beginPath")
		ctx.Call("arc", ballX, ballY, ballR, 0, 2*3.14159)
		ctx.Set("fillStyle", "#0095DD")
		ctx.Call("fill")
		ctx.Call("closePath")

		// paddle
		ctx.Call("beginPath")
		ctx.Call("rect", paddleX, height-paddleHeight-10, paddleWidth, paddleHeight)
		ctx.Set("fillStyle", "#0095DD")
		ctx.Call("fill")
		ctx.Call("closePath")

		// bricks
		for i := 0; i < rows; i++ {
			for j := 0; j < cols; j++ {
				if bricks[i][j] == 1 {
					bx := offsetLeft + float64(j)*(brickW+brickP)
					by := offsetTop + float64(i)*(brickH+brickP)
					ctx.Call("beginPath")
					ctx.Call("rect", bx, by, brickW, brickH)
					ctx.Set("fillStyle", "#FF5733")
					ctx.Call("fill")
					ctx.Call("closePath")
				}
			}
		}

		js.Global().Call("requestAnimationFrame", raf)
		return nil
	})

	js.Global().Call("requestAnimationFrame", raf)
}
