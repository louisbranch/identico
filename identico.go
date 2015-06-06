package identico

import (
	"image"
	"image/color"
	"image/draw"
)

func Classic(mask image.Image, bg, fg color.Color) image.Image {
	bounds := mask.Bounds()
	bgimg := FillBackground(bounds, bg)
	fgimg := ReplaceMask(mask, fg)

	dst := image.NewNRGBA(bounds)
	draw.Draw(dst, bounds, bgimg, image.ZP, draw.Src)
	draw.Draw(dst, bounds, fgimg, image.ZP, draw.Over)
	return dst
}

func FillBackground(bounds image.Rectangle, col color.Color) image.Image {
	img := image.NewNRGBA(bounds)
	draw.Draw(img, bounds, &image.Uniform{col}, image.ZP, draw.Src)
	return img
}

func ReplaceMask(mask image.Image, col color.Color) image.Image {
	bounds := mask.Bounds()
	dst := image.NewNRGBA(bounds)
	w, h := bounds.Max.X, bounds.Max.Y
	r, g, b, _ := col.RGBA()
	rgba := color.NRGBA{
		uint8(r >> 8),
		uint8(g >> 8),
		uint8(b >> 8),
		0,
	}

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			pixel := mask.At(x, y)
			_, _, _, alpha := pixel.RGBA()
			if alpha != 0 {
				rgba.A = uint8(alpha >> 8)
				dst.Set(x, y, rgba)
			} else {
				dst.Set(x, y, pixel)
			}
		}
	}
	return dst
}
