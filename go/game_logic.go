package main

// Ball models the moving ball in the game.
type Ball struct {
    X  float64
    Y  float64
    DX float64
    DY float64
    R  float64
}

// Paddle models the player's paddle.
type Paddle struct {
    X      float64
    Y      float64
    Width  float64
    Height float64
    Speed  float64
}

// Brick models a breakable brick.
type Brick struct {
    Active bool
    X      float64
    Y      float64
    Width  float64
    Height float64
}

// Update moves the ball and resolves wall, paddle, and brick collisions.
func (b *Ball) Update(width, height float64, paddle *Paddle, bricks [][]*Brick) {
    b.X += b.DX
    b.Y += b.DY

    // Wall collisions
    if b.X+b.R > width || b.X-b.R < 0 {
        b.DX = -b.DX
    }
    if b.Y-b.R < 0 {
        b.DY = -b.DY
    } else if b.Y+b.R > paddle.Y && b.Y-b.R < paddle.Y+paddle.Height {
        if b.X > paddle.X && b.X < paddle.X+paddle.Width {
            b.DY = -b.DY
            diff := b.X - (paddle.X + paddle.Width/2)
            b.DX = diff / (paddle.Width / 2) * 4
        }
    } else if b.Y-b.R > height {
        b.Reset(width, height)
    }

    // Brick collisions
    for i := range bricks {
        for j := range bricks[i] {
            brick := bricks[i][j]
            if brick.Active && b.CollidesWith(brick) {
                b.DY = -b.DY
                brick.Active = false
            }
        }
    }
}

// CollidesWith returns true if the ball overlaps the brick.
func (b *Ball) CollidesWith(brick *Brick) bool {
    return (b.X-b.R) < brick.X+brick.Width &&
        (b.X+b.R) > brick.X &&
        (b.Y-b.R) < brick.Y+brick.Height &&
        (b.Y+b.R) > brick.Y
}

// Reset puts the ball back in the center with the starting velocity.
func (b *Ball) Reset(width, height float64) {
    b.X = width / 2
    b.Y = height / 2
    b.DX = 2.5
    b.DY = -2.5
}

// Update moves the paddle left or right while clamping to the canvas.
func (p *Paddle) Update(leftPressed, rightPressed bool, width float64) {
    if leftPressed {
        p.X -= p.Speed
        if p.X < 0 {
            p.X = 0
        }
    }
    if rightPressed {
        p.X += p.Speed
        if p.X+p.Width > width {
            p.X = width - p.Width
        }
    }
}
