package pet_maker

import (
	"bytes"
	"fmt"
	_ "golang.org/x/image/webp" // Register WebP format
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"net/http"
	"os"

	xdraw "golang.org/x/image/draw"
)

func loadImageFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status code %d", resp.StatusCode)
	}

	img, _, err := image.Decode(resp.Body)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func loadPNG(path string) (image.Image, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	img, err := png.Decode(f)
	if err != nil {
		return nil, err
	}
	return img, nil
}

func MakePetImage(url string, speed float64, width int, height int) *bytes.Reader {
	const frames = 5
	var gifFrames []*image.Paletted
	var delays []int

	baseImg, err := loadImageFromURL(url)
	if err != nil {
		panic(err)
	}

	// support transparency
	customPalette := make([]color.Color, len(palette.Plan9))
	copy(customPalette, palette.Plan9)
	customPalette[0] = color.RGBA{0, 0, 0, 0}
	bounceOffset := []int{0, 3, 6, 4, 2} // based on a 128x28 grid

	for i := range frames {
		overlayPath := fmt.Sprintf("./pet_maker/pet_images/%d.png", i)
		overlayImg, err := loadPNG(overlayPath)
		if err != nil {
			panic(err)
		}

		rgba := image.NewRGBA(image.Rect(0, 0, width, height))
		draw.Draw(rgba, rgba.Bounds(), &image.Uniform{C: image.Transparent}, image.Point{}, draw.Src)

		left := int(0.23*float64(width)) - bounceOffset[i]
		top := int(0.18*float64(height)) + bounceOffset[i]
		right := int(0.98*float64(width)) + bounceOffset[i]
		bottom := int(0.94 * float64(height))
		baseTargetRect := image.Rect(left, top, right, bottom)
		baseSrcRect := baseImg.Bounds()
		xdraw.CatmullRom.Scale(rgba, baseTargetRect, baseImg, baseSrcRect, draw.Over, nil)

		overlayTargetRect := image.Rect(0, 0, width, height)
		overlaySrcRect := overlayImg.Bounds()
		xdraw.CatmullRom.Scale(rgba, overlayTargetRect, overlayImg, overlaySrcRect, draw.Over, nil)

		paletted := image.NewPaletted(rgba.Bounds(), customPalette)
		draw.FloydSteinberg.Draw(paletted, rgba.Bounds(), rgba, image.Point{})

		gifFrames = append(gifFrames, paletted)
		delays = append(delays, int(5*speed))
	}

	outGif := &gif.GIF{
		Image:    gifFrames,
		Delay:    delays,
		Disposal: make([]byte, len(gifFrames)),
	}

	for i := range outGif.Disposal {
		outGif.Disposal[i] = gif.DisposalBackground
	}

	var bufMem bytes.Buffer
	err = gif.EncodeAll(&bufMem, outGif)
	if err != nil {
		panic(err)
	}
	return bytes.NewReader(bufMem.Bytes())
}
