package main

import (
	"image/color"
	"log"
	"math/rand"
	"time"

	"golang.org/x/image/font"
	"golang.org/x/image/font/opentype"

	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/hajimehoshi/ebiten/v2/examples/resources/fonts"
	"github.com/hajimehoshi/ebiten/v2/text"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

const sampleText = `  The quick brown fox jumps
over the lazy dog.`

var (
	mplusNormalFont font.Face
	mplusBigFont    font.Face
)

func init() {
	tt, err := opentype.Parse(fonts.MPlus1pRegular_ttf)
	if err != nil {
		log.Fatal(err)
	}

	const dpi = 72
	mplusNormalFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    24,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
	mplusBigFont, err = opentype.NewFace(tt, &opentype.FaceOptions{
		Size:    32,
		DPI:     dpi,
		Hinting: font.HintingFull,
	})
	if err != nil {
		log.Fatal(err)
	}
}

func init() {
	rand.Seed(time.Now().UnixNano())
}

type Game struct {
	counter        int
	kanjiText      []rune
	kanjiTextColor color.RGBA
	glyphs         []text.Glyph
}

func (g *Game) Update() error {
	// Initialize the glyphs for special (colorful) rendering.
	if len(g.glyphs) == 0 {
		g.glyphs = text.AppendGlyphs(g.glyphs, mplusNormalFont, sampleText)
	}
	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	gray := color.RGBA{0x80, 0x80, 0x80, 0xff}

	{
		const x, y = 160, 240
		const lineHeight = 80
		b := text.BoundString(text.FaceWithLineHeight(mplusBigFont, lineHeight), sampleText)
		ebitenutil.DrawRect(screen, float64(b.Min.X+x), float64(b.Min.Y+y), float64(b.Dx()), float64(b.Dy()), gray)
		text.Draw(screen, sampleText, text.FaceWithLineHeight(mplusBigFont, lineHeight), x, y, color.White)
	}
	{
		const x, y = 240, 400
		op := &ebiten.DrawImageOptions{}
		// g.glyphs is initialized by text.AppendGlyphs.
		// You can customize how to render each glyph.
		// In this example, multiple colors are used to render glyphs.
		for i, gl := range g.glyphs {
			op.GeoM.Reset()
			op.GeoM.Translate(x, y)
			op.GeoM.Translate(gl.X, gl.Y)
			op.ColorM.Reset()
			r := 1.0
			if i%3 == 0 {
				r = 0.5
			}
			g := 1.0
			if i%3 == 1 {
				g = 0.5
			}
			b := 1.0
			if i%3 == 2 {
				b = 0.5
			}
			op.ColorM.Scale(r, g, b, 1)
			screen.DrawImage(gl.Image, op)
		}
	}
}

func (g *Game) Layout(outsideWidth, outsideHeight int) (int, int) {
	return screenWidth, screenHeight
}

func main() {
	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Game Engine Demo")
	if err := ebiten.RunGame(&Game{}); err != nil {
		log.Fatal(err)
	}
}
