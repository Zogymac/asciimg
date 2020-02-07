package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/jpeg"
	//"math/bits"

	// Side-effect import.
	// Сайд-эффект — добавление декодера PNG в пакет image.
	_ "image/png"
	"os"
	// Внешняя зависимость.
	"golang.org/x/image/draw"
	"github.com/buger/goterm"
)
var (
	output = flag.String("o", "", "file to write")
	noscale = flag.Bool("noscale", false,"scaling")
	width = flag.Int("w", 200, "width")
	height = flag.Int("h", 40, "height")
	terminalSiz = flag.Bool("trm", false, "Terminal size")
)
func scale(img image.Image, w int, h int) image.Image {
	dstImg := image.NewRGBA(image.Rect(0, 0, w, h))
	draw.NearestNeighbor.Scale(dstImg, dstImg.Bounds(), img, img.Bounds(), draw.Over, nil)
	return dstImg
}

func decodeImageFile(imgName string) (image.Image, error) {
	imgFile, err := os.Open(imgName)
	if err != nil {
		return nil, err
	}

	img, _, err := image.Decode(imgFile)

	return img, err
}

func processPixel(c color.Color) rune {
	gc := color.GrayModel.Convert(c)
	r, _, _, _ := gc.RGBA()
	r = r >> 8
	symbols := []rune("@80GCLft1i;:,. ")
	return symbols[r*uint32(len(symbols)-1)/256]
}

func convertToAscii(img image.Image) [][]rune {
	textImg := make([][]rune, img.Bounds().Dy())
	for i := range textImg {
		textImg[i] = make([]rune, img.Bounds().Dx())
	}

	for i := range textImg {
		for j := range textImg[i] {
			textImg[i][j] = processPixel(img.At(j, i))
		}
	}
	return textImg
}

func checkError(e error) {
	if e != nil {
		fmt.Println(e)
		os.Exit(1)
	}
}
func main() {
	flag.Parse()

	if flag.NArg() == 0 {
		fmt.Println("Usage: asciimg <image.jpg>")
		os.Exit(0)
	}
	imgName := flag.Arg(0)

	img, err := decodeImageFile(imgName)
	if err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
	o := os.Stdout
	if *output != "" {
		f, err := os.Create(*output)
		defer f.Close()
		checkError(err)
		o = f
	}
	if *terminalSiz == true {
		*width = goterm.Width()
		*height = goterm.Height()
	}
	if *output == "" && *noscale == false {
		img = scale(img, *width, *height)
	}
	if *noscale == true {
		img = scale(img, *width, *height)
	}
	textImg := convertToAscii(img)
	for i := range textImg {
		for j := range textImg[i] {
			fmt.Fprintf(o,"%c", textImg[i][j])
			//Printf("%c", textImg[i][j])
		}
		fmt.Fprint(o, "\n")
	}
}
