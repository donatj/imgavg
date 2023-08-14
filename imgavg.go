package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/gif"
	"image/jpeg"
	"image/png"
	"log"
	"os"
)

func pictable(dx int, dy int) [][][]uint64 {
	pic := make([][][]uint64, dx) /* type declaration */
	for i := range pic {
		pic[i] = make([][]uint64, dy) /* again the type? */
		for j := range pic[i] {
			pic[i][j] = []uint64{0, 0, 0}
		}
	}
	return pic
}

func avgImageFromPictable(avgdata [][][]uint64, n int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, len(avgdata), len(avgdata[0])))

	o := uint64(n)

	for x := 0; x < len(avgdata); x++ {
		for y := 0; y < len(avgdata[0]); y++ {
			mycolor := color.RGBA{uint8(avgdata[x][y][0] / o), uint8(avgdata[x][y][1] / o), uint8(avgdata[x][y][2] / o), 255}
			img.Set(x, y, mycolor)
		}
	}

	return img
}

func init() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of %s [dir] [outputfile]:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse()

	if len(flag.Args()) < 1 {
		flag.Usage()
		os.Exit(2)
	}

	image.RegisterFormat("png", "png", png.Decode, png.DecodeConfig)
	image.RegisterFormat("jpeg", "jpeg", jpeg.Decode, jpeg.DecodeConfig)
	image.RegisterFormat("gif", "gif", gif.Decode, gif.DecodeConfig)
}

func main() {
	outputfile := ""
	if flag.NArg() > 1 {
		outputfile = flag.Arg(flag.NArg() - 1)
	} else {
		outputfile = "output.png"
	}

	inputfiles := flag.Args()
	inputfiles = inputfiles[:len(inputfiles)-1]

	fileList, err := getFiles(inputfiles)
	if err != nil {
		log.Fatal(err)
	}

	fileList = filterFiles(fileList,
		[]string{".png", ".jpg", ".jpeg", ".gif"})

	// Lets create this before hand just in case so the user doesn't get screwed
	f, err := os.OpenFile(outputfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}

	avgdata := [][][]uint64{}
	picinit := false

	n := 0

	for _, fname := range fileList {
		n++
		log.Println("Loading", fname)

		file, err := os.Open(fname)
		if err != nil {
			log.Fatal(fname, " ", err)
		}

		m, _, err := image.Decode(file)
		if err != nil {
			log.Fatal(fname, " ", err)
		}
		bounds := m.Bounds()

		if !picinit {
			avgdata = pictable(bounds.Max.X+2, bounds.Max.Y+2)
			picinit = true
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				avgdata[x][y][0] += uint64((float32(r) / 65535) * 255)
				avgdata[x][y][1] += uint64((float32(g) / 65535) * 255)
				avgdata[x][y][2] += uint64((float32(b) / 65535) * 255)
			}
		}

		file.Close()
	}

	img := avgImageFromPictable(avgdata, n)

	if err = png.Encode(f, img); err != nil {
		log.Fatal(err)
	}

}
