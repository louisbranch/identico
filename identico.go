package identico

import (
	"image"
	"image/color"
)

func ReplaceMask(src image.Image, col color.Color) image.Image {
	bounds := src.Bounds()
	dst := image.NewRGBA(bounds)
	w, h := bounds.Max.X, bounds.Max.Y

	for x := 0; x < w; x++ {
		for y := 0; y < h; y++ {
			rgba := col.(color.NRGBA)
			pixel := src.At(x, y).(color.NRGBA)
			if pixel.A != 0 {
				rgba.A = pixel.A
				dst.Set(x, y, rgba)
			} else {
				dst.Set(x, y, pixel)
			}
		}
	}
	return dst
}
