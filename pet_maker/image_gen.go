package pet_maker

import (
	"image"
	"image/color"
	"image/color/palette"
	"image/draw"
	"image/gif"
	"image/png"
	"os"
	"net/http"
	"fmt"

	xdraw "golang.org/x/image/draw"
)

func loadPNGFromURL(url string) (image.Image, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch image: status code %d", resp.StatusCode)
	}

	img, err := png.Decode(resp.Body)
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

func main() {
	const frames = 5
	var gifFrames []*image.Paletted
	var delays []int

	baseImg, err := loadPNGFromURL("https://cdn.discordapp.com/avatars/143090142360371200/a73420b217a77a77b17fb42fa7ecfbcc.png?size=4096")
	if err != nil {
		panic(err)
	}

	// support transparency
	customPalette := make([]color.Color, len(palette.Plan9))
	copy(customPalette, palette.Plan9)
	customPalette[0] = color.RGBA{0, 0, 0, 0}

	bounceOffset := []int {0, 70, 140, 90, 50}

	for i := range frames {
		overlayPath := fmt.Sprintf("%d.png", i)
		overlayImg, err := loadPNG(overlayPath)
		if err != nil {
			panic(err)
		}

		rgba := image.NewRGBA(image.Rect(0, 0, 1024, 1024))
		draw.Draw(rgba, rgba.Bounds(), &image.Uniform{C: image.Transparent}, image.Point{}, draw.Src)

		baseTargetRect := image.Rect(200, 200 + bounceOffset[i], 1024, 1024)
		baseSrcRect := baseImg.Bounds()
		xdraw.CatmullRom.Scale(rgba, baseTargetRect, baseImg, baseSrcRect, draw.Over, nil)

		overlayTargetRect := image.Rect(0, 0, 1024, 1024)
		overlaySrcRect := overlayImg.Bounds()
		xdraw.CatmullRom.Scale(rgba, overlayTargetRect, overlayImg, overlaySrcRect, draw.Over, nil)

		paletted := image.NewPaletted(rgba.Bounds(), customPalette)
		draw.FloydSteinberg.Draw(paletted, rgba.Bounds(), rgba, image.Point{})

		gifFrames = append(gifFrames, paletted)
		delays = append(delays, 5)
	}

	outGif := &gif.GIF{
		Image: gifFrames,
		Delay: delays,
		Disposal: make([]byte, len(gifFrames)),
	}

	for i := range outGif.Disposal {
		outGif.Disposal[i] = gif.DisposalBackground
	}

	outFile, err := os.Create("output.gif")
	if err != nil {
		panic(err)
	}
	defer outFile.Close()

	err = gif.EncodeAll(outFile, outGif)
	if err != nil {
		panic(err)
	}

	fmt.Println("GIF created as output.gif")
}
