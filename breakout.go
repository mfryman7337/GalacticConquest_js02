package main

import (
    "fmt"
//    "math"
    "math/rand"
    "time"
	"github.com/faiface/pixel"
    "github.com/faiface/pixel/imdraw"
	"github.com/faiface/pixel/pixelgl"
	"golang.org/x/image/colornames"
    "golang.org/x/image/font/basicfont"
	"github.com/faiface/pixel/text"
)

const WINDOW_HEIGHT  = 1300     // Width of entire PNG image
const WINDOW_WIDTH   = 1300     // Height of entire PNG image
const GAME_WIDTH     = 1100
const GAME_CENTER    = WINDOW_WIDTH / 2
const EDGE_WIDTH     = 20
const LEFT_EDGE      = (WINDOW_WIDTH - GAME_WIDTH) / 2
const RIGHT_EDGE     = WINDOW_WIDTH - LEFT_EDGE
const TOP_EDGE       = WINDOW_HEIGHT - 60
const BOTTOM_EDGE    = 20

//var imd  imdraw
type block_t struct {
    Xloc                 float64
    Yloc                 float64
    Rect                 pixel.Rect
    Width                float64
    Score                float64
    Visible              bool
//    Color                RGBA
}

type paddle_t struct {
    Xloc                 float64
    Yloc                 float64
    Width                float64
    Rect                 pixel.Rect
//    Color                RGBA
}

type ball_t struct {
    Xloc                 float64
    Yloc                 float64
    Xspeed               float64
    Yspeed               float64
//    Color                RGBA
}


var Block   [100]   block_t
var Paddle          paddle_t
var Ball            ball_t
var EdgeL, EdgeT, EdgeR   pixel.Rect
var TotalScore      float64
var TotalBalls      float64
var ColorWheel      pixel.RGBA
var R, G, B         float64
var C, CC           int32


func Start_Color() {   
    if CC >= 400 {
        CC = 0
    } else {
        CC = CC + 1
    }
    C = CC
}

func Next_Color() pixel.RGBA {
    var  row, col float64
    
    C = C + 1
    if C >= 400 {
        C = 0
    }
    row = float64(C / 100)
    col = float64(C % 100)
    
    switch(int(row)) {
        case 0:
            R = col * 0.01
            B = (100.0-col) * 0.01
            break
            
        case 1:
            G = col * 0.01
            break

        case 2:
            R = (100.0-col) * 0.01
            break
            
        case 3:
            B = col * 0.01
            break
    }

    ColorWheel = pixel.RGB(R, G, B)
    
    return(ColorWheel)
}


func Game_Init() {
    var X, Y int
    var i, columns, width int
    
    TotalScore = 0.0
    TotalBalls = 5.0
    
    Paddle.Xloc  = 650.0
    Paddle.Yloc  = 50.0
    Paddle.Width = 80.0
    
    Ball.Xloc = GAME_CENTER
    Ball.Yloc = 100
    Ball.Xspeed = 5
    Ball.Yspeed = 10
    
    EdgeL.Min = pixel.V(LEFT_EDGE-EDGE_WIDTH, BOTTOM_EDGE)
    EdgeL.Max = pixel.V(LEFT_EDGE, TOP_EDGE+EDGE_WIDTH)
    EdgeT.Min = pixel.V(LEFT_EDGE-EDGE_WIDTH, TOP_EDGE)
    EdgeT.Max = pixel.V(RIGHT_EDGE, TOP_EDGE+EDGE_WIDTH)
    EdgeR.Min = pixel.V(RIGHT_EDGE, BOTTOM_EDGE)
    EdgeR.Max = pixel.V(RIGHT_EDGE+EDGE_WIDTH, TOP_EDGE+EDGE_WIDTH)
    
    columns = 10
    width = GAME_WIDTH / columns
    for i=0; i<60; i++ {
        X = i % columns
        Y = i / columns
        Block[i].Xloc = float64(GAME_CENTER - ((columns-1) * width)/2 + (X * width))
        Block[i].Yloc = float64((TOP_EDGE - 150) - (Y * 50))
        Block[i].Width = float64(GAME_WIDTH / columns) * 0.6
        Block[i].Rect.Min = pixel.V(Block[i].Xloc-Block[i].Width/2, Block[i].Yloc-5.0)
        Block[i].Rect.Max = pixel.V(Block[i].Xloc+Block[i].Width/2, Block[i].Yloc+5.0)
        Block[i].Score = float64(8 - Y) * 100.0
        Block[i].Visible = true
    }
}


func run() {
    var loc, ball pixel.Vec
    var i int
//    var test pixel.RGBA
    
	cfg := pixelgl.WindowConfig{
		Title:  "Breakout!",
		Bounds: pixel.R(0, 0, WINDOW_WIDTH, WINDOW_HEIGHT),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}


    imd := imdraw.New(nil)
    imd.Color = pixel.RGB(0, 0, 80)
    imd.EndShape = imdraw.RoundEndShape

	for !win.Closed() {
        // Check for the Left Arrow
        if win.JustPressed(pixelgl.KeyLeft) {
            if Paddle.Xloc > 100 {
                Paddle.Xloc -= 20
            }
        }
        // Check for the Right Arrow
        if win.JustPressed(pixelgl.KeyRight) {
            if Paddle.Xloc < 1200 {
                Paddle.Xloc += 20
            }
        }
        
        // Heck, just follow the mouse pointer if it's in the window
        loc = win.MousePosition()
        
        // Check for a wall-bounce
        ball = pixel.V(Ball.Xloc, Ball.Yloc)
        if EdgeL.Contains(ball) || EdgeR.Contains(ball) {
            Ball.Xspeed = -Ball.Xspeed
        }
        if EdgeT.Contains(ball) {
            Ball.Yspeed = -Ball.Yspeed
        }
        
        // Check if we've hit any blocks
        for i=0; i<100; i++ {
            if Block[i].Rect.Contains(ball) {
                if Block[i].Visible {
                    Block[i].Visible = false
                    TotalScore += Block[i].Score
                    Ball.Yspeed = -Ball.Yspeed
                    fmt.Printf("TotalScore %8.0f\n", TotalScore)
                }
            }
        }
//        fmt.Printf("%6.4f, %6.4f, %6.4f\n", test.R, test.G, test.B)
        
        // Check if the ball hits the paddle, and create an exit angle
        if Paddle.Rect.Contains(ball) {
            Ball.Yspeed = -Ball.Yspeed
            Ball.Xspeed = (Ball.Xloc - Paddle.Xloc) / 8.0
        }
        
        // Check if the ball was missed
        if (Ball.Yloc < (Paddle.Yloc - 10.0)) {
            TotalBalls = TotalBalls - 1.0
            Ball.Xloc = GAME_CENTER
            Ball.Yloc = 100
            if TotalBalls > 0 {
                Ball.Xspeed = 5
                Ball.Yspeed = 10
            } else {
                Ball.Xspeed = 0
                Ball.Yspeed = 0
            }
        }
        
        ///////////////////////////////////////
        // OK, we're done checking thngs...
        ///////////////////////////////////////
        
        // Redraw the entire playing field
        imd.Clear()
        
        // Draw the outer edges of the game
        imd.Color = pixel.RGB(0.0, 0.0, 0.4)
        imd.EndShape = imdraw.RoundEndShape
        imd.Push(EdgeL.Min, EdgeL.Max)
        imd.Rectangle(0.0)
        imd.Push(EdgeT.Min, EdgeT.Max)
        imd.Rectangle(0.0)
        imd.Push(EdgeR.Min, EdgeR.Max)
        imd.Rectangle(0.0)
        
        // Draw the individual blocks still showing
        Start_Color()
        imd.EndShape = imdraw.RoundEndShape     //SharpEndShape
        for i=0; i<100; i++ {
            if Block[i].Visible {
//                imd.Color = pixel.RGB((Block[i].Yloc-500) / 1000.0, 0.0, 0.0)
                imd.Color = Next_Color()
                imd.Push(pixel.V(Block[i].Xloc-(Block[i].Width/2.0), Block[i].Yloc), pixel.V(Block[i].Xloc+(Block[i].Width/2.0), Block[i].Yloc))
                imd.Line(20)
            }
        }

        // Draw the paddle in it's current location
        if (loc.X > (LEFT_EDGE+Paddle.Width/2)) && (loc.X < (RIGHT_EDGE-Paddle.Width/2)) {
            Paddle.Xloc = loc.X
        }
        Paddle.Rect.Min = pixel.V(Paddle.Xloc-(Paddle.Width/2), Paddle.Yloc-5.0)
        Paddle.Rect.Max = pixel.V(Paddle.Xloc+(Paddle.Width/2), Paddle.Yloc+5.0)
        imd.Color = pixel.RGB(0.0, 0.4, 0.0)
        imd.EndShape = imdraw.RoundEndShape
        imd.Push(pixel.V(Paddle.Xloc-(Paddle.Width/2), Paddle.Yloc), pixel.V(Paddle.Xloc+(Paddle.Width/2), Paddle.Yloc))
        imd.Line(10)
        
        // Draw the ball in it's new location
        Ball.Xloc += Ball.Xspeed
        Ball.Yloc += Ball.Yspeed
        imd.Color = pixel.RGB(0.8, 0.0, 0.0)
        imd.Push(pixel.V(Ball.Xloc, Ball.Yloc))
        imd.Circle(5, 10)        
            
        basicAtlas := text.NewAtlas(basicfont.Face7x13, text.ASCII)
        basicTxt := text.New(pixel.V(GAME_CENTER-200, WINDOW_HEIGHT-30), basicAtlas)
        basicTxt.Color = colornames.Blue
        fmt.Fprintln(basicTxt, "Your Score -", TotalScore)

        gameTxt := text.New(pixel.V(GAME_CENTER-200, 600), basicAtlas)
        if TotalBalls == 0 {
            gameTxt.Color = colornames.Red
            fmt.Fprintln(gameTxt, "G A M E  O V E R")
            fmt.Fprintln(gameTxt, " Score -", TotalScore)


            if win.JustPressed(pixelgl.MouseButtonLeft) {
                Game_Init()
            }
		}
        
        win.Clear(colornames.Whitesmoke)
        imd.Draw(win)
        basicTxt.Draw(win, pixel.IM.Scaled(basicTxt.Orig, 3))
        gameTxt.Draw(win, pixel.IM.Scaled(gameTxt.Orig, 4))
		win.Update()
	}
}

func main() {
    fmt.Printf("\r\n\n\nBreakout started\n")
    Game_Init()
    now := time.Now()
    rand.Seed(int64(now.Nanosecond()))

	pixelgl.Run(run)
}


/*

package main

import (
    "gg"
    "fmt"
    "math"
    "os"
)

const HEIGHT  = 2048      // Width of entire PNG image
const WIDTH   = 2048      // Height of entire PNG image


func Draw_Pattern(dia1 float64, step1 float64, dia2 float64, step2 float64) {

    var r, angle, radian, graphCenterX, graphCenterY float64
    var X1, Y1, X2, Y2 float64
    var imageFname string
    
	dc := gg.NewContext(WIDTH, HEIGHT)
//    graphWidth  = WIDTH
//    graphHeight = HEIGHT
    graphCenterX = WIDTH / 2.0
    graphCenterY = HEIGHT / 2.0 
    
    dc.SetRGB(240.0, 240.0, 240.0)
    dc.Clear()
    dc.SetLineWidth(1.0)


    r = 300.0
    dc.SetRGB(80.0, 80.0, 80.0)
    for angle=0.0; angle<360.0; angle=angle+1.0 {
        radian = angle * (180.0 / math.Pi)
        X1 = (r * dia1) * math.Cos(radian * step1)
        Y1 = (r * dia1) * math.Sin(radian * step1)
        X2 = (r * dia2) * math.Cos(radian * step2)
        Y2 = (r * dia2) * math.Sin(radian * step2)
//        fmt.Printf("%6.3f %6.3f %6.3f %6.3f\n", X1, Y1, X2, Y2)
        dc.DrawLine(graphCenterX+X1, graphCenterY+Y1, graphCenterX+X2, graphCenterY+Y2)
        dc.Stroke()
    }
    
//    sMsg = fmt.Sprintf("Parameters %d  %d", x, y)
//    dc.DrawString(sMsg, HEIGHT-10, (WIDTH / 2) - 100)
//
    imageFname = fmt.Sprintf("lissajous.png")
    os.Remove(imageFname)
    dc.SavePNG(imageFname)
}



func main ( ) {
//    var command, param string
    
    Draw_Pattern(2.0, 7.0, 3.0, 3.0);
    
    
//    fmt.Scanln(&command, &param)
}





func run() {
	cfg := pixelgl.WindowConfig{
		Title:  "Pixel Rocks!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync:  true,
	}
	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}

    imd := imdraw.New(nil)

    imd.Color = pixel.RGB(0, 0, 0)
	imd.Push(pixel.V(200, 100))
	imd.Color = pixel.RGB(0.5, 0.5, 0.5)
	imd.Push(pixel.V(800, 100))
	imd.Color = pixel.RGB(1, 1, 1)
	imd.Push(pixel.V(500, 700))
	imd.Polygon(0)
    
	for !win.Closed() {
		win.Clear(colornames.Aliceblue)
        imd.Draw(win)
		win.Update()
	}
}



*/
