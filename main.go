package main

import (
	"fmt"
	"image/color"
	"io/ioutil"
	"log"
	"math/rand"
	"time"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/text"
	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"
)

const (
	WIDTH    = 840
	HEIGHT   = 640
	TILESIZE = 100
)

var (
	RED    = color.RGBA{255, 0, 0, 255}
	GREEN  = color.RGBA{0, 255, 0, 255}
	ORANGE = color.RGBA{222, 136, 24, 255}
	YELLOW = color.RGBA{237, 233, 12, 255}
	GRAY   = color.RGBA{108, 122, 137, 255}
	PURPLE = color.RGBA{227, 11, 202, 255}
	CYAN   = color.RGBA{25, 209, 194, 255}
)

func CheckPointCollision(mx int, my int, rect Kotak) bool {
	if mx >= int(rect.x) && // right of the left edge AND
		mx <= int(rect.x+float64(rect.width)) && // left of the right edge AND
		my >= int(rect.y) && // below the top AND
		my <= int(rect.y+float64(rect.height)) { // above the bottom
		return true
	}
	return false
}

type Kotak struct {
	width, height int
	x, y          float64
	open, clear   bool
	color         color.RGBA
}

func (k *Kotak) Draw(screen *ebiten.Image) {
	if k.open {
		ebitenutil.DrawRect(screen, k.x, k.y, float64(k.width), float64(k.height), k.color)
	} else {
		ebitenutil.DrawRect(screen, k.x, k.y, float64(k.width), float64(k.height), GRAY)
	}
}

type Game struct {
	bg            *ebiten.Image
	kotaks        []Kotak
	mousePosition struct {
		x int
		y int
	}
	gameOver      bool
	previousKotak *Kotak
	nextKotak     *Kotak
	turn          int
}

var (
	count     = 0
	f         font.Face
	creditF   font.Face
	gameOverF font.Face
	overCount = 0
)

func (g *Game) Update() error {
	g.mousePosition.x, g.mousePosition.y = ebiten.CursorPosition()

	for i := 0; i < len(g.kotaks); i++ {
		if ebiten.IsMouseButtonPressed(ebiten.MouseButtonLeft) && CheckPointCollision(g.mousePosition.x, g.mousePosition.y, g.kotaks[i]) {
			if g.previousKotak == nil {
				g.kotaks[i].open = true
				g.previousKotak = &g.kotaks[i]

			} else if g.nextKotak == nil && g.previousKotak != &g.kotaks[i] {
				g.kotaks[i].open = true
				g.nextKotak = &g.kotaks[i]

			}
		}
	}

	if g.previousKotak != nil && g.nextKotak != nil {
		if g.previousKotak.color == g.nextKotak.color {
			count++
			if count == 50 {
				g.previousKotak.clear = true
				g.nextKotak.clear = true
				count = 0
			}
			if count == 0 {
				g.previousKotak = nil
				g.nextKotak = nil
				g.turn++
				overCount++
			}

		} else {
			count++
			if count == 50 {
				g.previousKotak.open = false
				g.nextKotak.open = false
				count = 0
			}
			if count == 0 {
				g.previousKotak = nil
				g.nextKotak = nil
				g.turn++
			}
		}
	}

	if overCount == 6 {
		g.gameOver = true
	}

	if ebiten.IsKeyPressed(ebiten.KeyR) && g.gameOver {
		overCount = 0
		g.turn = 0
		count = 0
		g.gameOver = false
		g.previousKotak = nil
		g.nextKotak = nil
		for i := 0; i < len(g.kotaks); i++ {
			g.kotaks[i].clear = false
			g.kotaks[i].open = false
		}
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(g.kotaks), func(i, j int) { g.kotaks[i].color, g.kotaks[j].color = g.kotaks[j].color, g.kotaks[i].color })
	}

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	bp := &ebiten.DrawImageOptions{}
	bp.GeoM.Scale(1.5, 1.5)
	screen.DrawImage(g.bg, bp)

	for _, k := range g.kotaks {
		if !k.clear {
			k.Draw(screen)
		}
	}

	if g.gameOver {
		text.Draw(screen, "Game Over", gameOverF, WIDTH/2-110, HEIGHT/3, RED)
		text.Draw(screen, "Tekan \"R\" untuk restart!", f, WIDTH/2-150, HEIGHT/2, color.White)
	}

	text.Draw(screen, "Remember Me", f, WIDTH/2-80, 35, color.White)
	text.Draw(screen, fmt.Sprintf("Turn: %v", g.turn), f, 376, 535, color.White)
	text.Draw(screen, "created by aji mustofa @pepega90", creditF, 10, 620, color.RGBA{202, 222, 24, 255})

	// ebitenutil.DebugPrint(screen, fmt.Sprintf("g.gameOver: %v", g.gameOver))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("count: %v\n", count))
	// ebitenutil.DebugPrint(screen, fmt.Sprintf("Mouse x: %v\nMouse y: %v", g.mousePosition.x, g.mousePosition.y))
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (screenWidth, screenHeight int) {
	return WIDTH, HEIGHT
}

func main() {
	ebiten.SetWindowSize(WIDTH, HEIGHT)
	ebiten.SetWindowTitle("Remember Me")

	// load font
	fontBytes, err := ioutil.ReadFile("./assets/Pixellari.ttf")
	if err != nil {
		log.Fatalf("readfile: %v", err)
		return
	}

	tt, err := opentype.Parse(fontBytes) // <= custom font
	// tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf) // <= wasm font
	if err != nil {
		log.Fatalf("opentype load font: %v", err)
		return
	}

	f, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    30,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	creditF, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    20,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	gameOverF, _ = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    50,
		DPI:     72,
		Hinting: font.HintingFull,
	})

	// load assets
	bgImg, _, _ := ebitenutil.NewImageFromFile("./assets/bg.png")

	g := new(Game)

	g.turn = 0
	g.bg = bgImg
	g.gameOver = false

	// load kotak
	baris := 3
	kolom := 4

	w := []color.RGBA{RED, ORANGE, GREEN, YELLOW}

	for y := 0; y < baris; y++ {
		wIndex := 0
		for x := 0; x < kolom; x++ {
			if y == 0 {
				g.kotaks = append(g.kotaks, Kotak{
					x:      150 + float64(x*150),
					y:      60 + float64(y*150),
					width:  TILESIZE,
					height: TILESIZE,
					color:  w[wIndex],
					open:   false,
					clear:  false,
				})
				wIndex++
			}
			if y == 1 {
				g.kotaks = append(g.kotaks, Kotak{
					x:      150 + float64(x*150),
					y:      60 + float64(y*150),
					width:  TILESIZE,
					height: TILESIZE,
					color:  w[wIndex],
					open:   false,
					clear:  false,
				})
				wIndex++
			}
			if y == 2 {
				if x == 0 || x == 1 {
					g.kotaks = append(g.kotaks, Kotak{
						x:      150 + float64(x*150),
						y:      60 + float64(y*150),
						width:  TILESIZE,
						height: TILESIZE,
						color:  PURPLE,
						open:   false,
						clear:  false,
					})
				}

				if x == 2 || x == 3 {
					g.kotaks = append(g.kotaks, Kotak{
						x:      150 + float64(x*150),
						y:      60 + float64(y*150),
						width:  TILESIZE,
						height: TILESIZE,
						color:  CYAN,
						open:   false,
						clear:  false,
					})
				}
			}

		}
	}

	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(g.kotaks), func(i, j int) { g.kotaks[i].color, g.kotaks[j].color = g.kotaks[j].color, g.kotaks[i].color })

	if err := ebiten.RunGame(g); err != nil {
		log.Fatal(err)
	}
}
