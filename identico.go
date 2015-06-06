package identico

import (
	"image"
	"image/color"
	"image/draw"
	"io/ioutil"

	"code.google.com/p/freetype-go/freetype"
	"code.google.com/p/freetype-go/freetype/truetype"
)

type Font struct {
	font    *truetype.Font
	size    float64
	spacing float64
	ctx     *freetype.Context
}

func Classic(mask image.Image, bg, fg color.Color, font Font, letter rune) image.Image {
	bounds := mask.Bounds()
	w, h := bounds.Max.X, bounds.Max.Y
	bgimg := FillBackground(w, h, bg)
	fgimg := ReplaceMask(mask, fg)

	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, bgimg, image.ZP, draw.Src)
	draw.Draw(dst, bounds, fgimg, image.ZP, draw.Over)

	text := image.NewNRGBA(bounds)
	font.ctx.SetClip(bounds)
	font.ctx.SetSrc(image.White)
	font.ctx.SetDst(text)

	offX, offY := offsets(font, letter)
	center := image.Rect((w/2 - offX/2), (h/2 - offY/2), w, h)

	drawLetter(letter, font)
	draw.Draw(dst, center, text, image.ZP, draw.Over)

	return dst
}

func NewFont(path string, size float64) (Font, error) {
	font := Font{}
	fontBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return font, err
	}
	parsed, err := freetype.ParseFont(fontBytes)
	if err != nil {
		return font, err
	}
	ctx := freetype.NewContext()
	font.font = parsed
	font.size = size
	ctx.SetFont(parsed)
	ctx.SetFontSize(size)
	ctx.SetDPI(72)
	ctx.SetHinting(freetype.FullHinting)
	font.ctx = ctx
	return font, nil
}

func FillBackground(width, height int, col color.Color) image.Image {
	img := image.NewNRGBA(image.Rect(0, 0, width, height))
	draw.Draw(img, img.Bounds(), &image.Uniform{col}, image.ZP, draw.Src)
	return img
}

func ReplaceMask(mask image.Image, col color.Color) image.Image {
	bounds := mask.Bounds()
	dst := image.NewNRGBA(bounds)
	w, h := bounds.Max.X, bounds.Max.Y
	r, g, b, _ := col.RGBA()

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pixel := mask.At(x, y)
			_, _, _, alpha := pixel.RGBA()
			if alpha != 0 {
				rgba := color.NRGBA{shift(r), shift(g), shift(b), shift(alpha)}
				dst.Set(x, y, rgba)
			} else {
				dst.Set(x, y, pixel)
			}
		}
	}
	return dst
}

func shift(v uint32) uint8 {
	return uint8(v >> 8)
}

func drawLetter(letter rune, font Font) error {
	pt := freetype.Pt(0, int(font.size))
	_, err := font.ctx.DrawString(string(letter), pt)
	return err
}

func offsets(f Font, l rune) (int, int) {
	font := f.font
	scale := f.size / float64(font.FUnitsPerEm())
	index := font.Index(l)
	width := int(font.HMetric(font.FUnitsPerEm(), index).AdvanceWidth)
	height := int(font.VMetric(font.FUnitsPerEm(), index).AdvanceHeight)
	return int(float64(width) * scale), int(float64(height) * scale)
}
