package main

import "testing"

func TestBallResetResetsToCenter(t *testing.T) {
	b := &Ball{X: 10, Y: 20, DX: 5, DY: -3, R: 8}
	b.Reset(800, 600)

	if b.X != 400 || b.Y != 300 {
		t.Fatalf("Reset() expected ball at center (400, 300), got (%v, %v)", b.X, b.Y)
	}
	if b.DX != 2.5 || b.DY != -2.5 {
		t.Fatalf("Reset() expected velocity (2.5, -2.5), got (%v, %v)", b.DX, b.DY)
	}
}

func TestPaddleUpdateClampsLeftBoundary(t *testing.T) {
	p := &Paddle{X: 5, Y: 580, Width: 100, Height: 10, Speed: 6}
	p.Update(true, false, 800)

	if p.X < 0 || p.X != 0 {
		t.Fatalf("Update() should clamp paddle X to 0, got %v", p.X)
	}
}

func TestCollidesWithBrickReturnsTrueOnOverlap(t *testing.T) {
	b := &Ball{X: 50, Y: 50, R: 5}
	brick := &Brick{X: 40, Y: 40, Width: 30, Height: 30}

	if !b.CollidesWith(brick) {
		t.Fatal("CollidesWith() expected overlap")
	}
}
