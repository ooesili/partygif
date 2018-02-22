package main

import (
	"errors"
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"io"
	"io/ioutil"
	"os"

	colorful "github.com/lucasb-eyer/go-colorful"
)

var config struct {
	inFile     string
	outputFile string
	cycles     int
	repeats    int
	frameRate  int
}

func init() {
	flag.StringVar(&config.inFile, "in", "", "input GIF file")
	flag.StringVar(&config.outputFile, "out", "", "output GIF file")
	flag.IntVar(&config.cycles, "cycles", 1, "number of color cycles during the GIF")
	flag.IntVar(&config.repeats, "repeats", 0, "number of times to repeat GIF before color shifting")
	flag.IntVar(&config.frameRate, "framerate", 10, "frame rate in 100ths of seconds for static GIFs")
}

func main() {
	if err := mainErr(); err != nil {
		fmt.Fprintf(os.Stderr, "partygif: error: %s\n", err)
		os.Exit(1)
	}
}

func mainErr() error {
	flag.Parse()

	if len(flag.Args()) != 0 {
		return errors.New("does not take any non-flag arguments")
	}

	inputFile, err := openInFile()
	if err != nil {
		return fmt.Errorf("opening input file: %s", err)
	}
	defer func() { _ = inputFile.Close() }()

	outputFile, err := openOutputFile()
	if err != nil {
		return fmt.Errorf("opening output file: %s", err)
	}
	defer func() { _ = outputFile.Close() }()

	return partyGIF(inputFile, outputFile)
}

func partyGIF(inputFile io.Reader, outputFile io.Writer) error {
	img, err := gif.DecodeAll(inputFile)
	if err != nil {
		return fmt.Errorf("decoding input file: %s", err)
	}

	addRepeats(img)
	colorShift(img)

	if err := gif.EncodeAll(outputFile, img); err != nil {
		return fmt.Errorf("encoding gif to output file: %s", err)
	}

	return nil
}

func addRepeats(img *gif.GIF) {
	if len(img.Image) == 1 {
		img.Delay = []int{config.frameRate}
	}

	originalLength := len(img.Image)

	for i := 0; i < config.repeats-1; i++ {
		for j := 0; j < originalLength; j++ {
			img.Image = append(img.Image, copyFrame(img.Image[j]))
			img.Delay = append(img.Delay, img.Delay[j])
			img.Disposal = append(img.Disposal, img.Disposal[j])
		}
	}

	// clean restart of original GIF
	for repeatIndex := 0; repeatIndex < config.repeats; repeatIndex++ {
		repeatBegin := repeatIndex * originalLength
		lastFrameInRepeat := repeatBegin + originalLength - 1
		img.Disposal[lastFrameInRepeat] = gif.DisposalBackground
	}
}

func copyFrame(frame *image.Paletted) *image.Paletted {
	copyPix := make([]uint8, len(frame.Pix))
	copy(copyPix, frame.Pix)

	copyPalette := make(color.Palette, len(frame.Palette))
	copy(copyPalette, frame.Palette)

	return &image.Paletted{
		Pix:     copyPix,
		Stride:  frame.Stride,
		Rect:    frame.Rect,
		Palette: copyPalette,
	}
}

func colorShift(img *gif.GIF) {
	hueStep := 360 / float64(len(img.Image)) * float64(config.cycles)

	for frameIndex, frame := range img.Image {
		hueShift := hueStep * float64(frameIndex)

		for i := range frame.Palette {
			frame.Palette[i] = shiftHue(hueShift, frame.Palette[i])
		}
	}
}

func shiftHue(shift float64, col color.Color) color.Color {
	_, _, _, alpha := col.RGBA()
	if alpha == 0 {
		return col
	}

	hue, chroma, lum := colorful.MakeColor(col).Hcl()
	hue += shift

	return alphaOverride{
		color: colorful.Hcl(hue, chroma, lum).Clamped(),
		alpha: alpha,
	}
}

type alphaOverride struct {
	color color.Color
	alpha uint32
}

func (c alphaOverride) RGBA() (uint32, uint32, uint32, uint32) {
	r, g, b, _ := c.color.RGBA()
	a := c.alpha
	return r, g, b, a
}

func openInFile() (io.ReadCloser, error) {
	if config.inFile == "" {
		return ioutil.NopCloser(os.Stdin), nil
	}

	file, err := os.Open(config.inFile)
	if err != nil {
		return nil, fmt.Errorf("opening file: %s: %s", config.inFile, err)
	}

	return file, nil
}

func openOutputFile() (io.WriteCloser, error) {
	if config.outputFile == "" {
		return struct {
			io.Writer
			io.Closer
		}{
			Writer: os.Stdout,
			Closer: ioutil.NopCloser(nil),
		}, nil
	}

	flags := os.O_WRONLY | os.O_CREATE | os.O_TRUNC
	file, err := os.OpenFile(config.outputFile, flags, 0644)
	if err != nil {
		return nil, fmt.Errorf("opening file: %s: %s", config.outputFile, err)
	}

	return file, nil
}
