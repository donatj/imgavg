package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	_ "image/gif"
	_ "image/jpeg"
	"image/png"
	"log"
	"os"
	"runtime"
	"sync"

	"github.com/schollz/progressbar/v3"
)

type pictable struct {
	data   []uint64 // R,G,B,R,G,B...
	dx, dy int
}

func newPictable(dx, dy int) pictable {
	return pictable{
		data: make([]uint64, dx*dy*3),
		dx:   dx,
		dy:   dy,
	}
}

func (p pictable) add(x, y int, r, g, b uint32) {
	i := (y*p.dx + x) * 3
	p.data[i] += uint64(r >> 8)
	p.data[i+1] += uint64(g >> 8)
	p.data[i+2] += uint64(b >> 8)
}

func avgImageFromPictable(p pictable, n int) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, p.dx, p.dy))
	div := uint64(n)

	for y := 0; y < p.dy; y++ {
		for x := 0; x < p.dx; x++ {
			i := (y*p.dx + x) * 3
			img.Set(x, y, color.RGBA{
				R: uint8(p.data[i] / div),
				G: uint8(p.data[i+1] / div),
				B: uint8(p.data[i+2] / div),
				A: 255,
			})
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
}

func main() {
	outputfile := "output.png"
	if flag.NArg() > 1 {
		outputfile = flag.Arg(flag.NArg() - 1)
	}

	inputfiles := flag.Args()[:len(flag.Args())-1]

	fileList, err := getFiles(inputfiles)
	if err != nil {
		log.Fatal(err)
	}

	fileList = filterFiles(fileList, []string{".png", ".jpg", ".jpeg", ".gif"})

	f, err := os.OpenFile(outputfile, os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()

	var (
		avgdata pictable
		n       int
	)

	images := make(chan image.Image, len(fileList))
	files := make(chan string, len(fileList))

	var wg sync.WaitGroup
	for i := 0; i < runtime.NumCPU()*2; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for fname := range files {
				file, err := os.Open(fname)
				if err != nil {
					log.Fatal(fname, " ", err)
				}
				m, _, err := image.Decode(file)
				file.Close()
				if err != nil {
					log.Fatal(fname, " ", err)
				}
				images <- m
			}
		}()
	}

	for _, fname := range fileList {
		files <- fname
	}
	close(files)

	go func() {
		wg.Wait()
		close(images)
	}()

	bar := progressbar.Default(int64(len(fileList)))
	for m := range images {
		n++
		bounds := m.Bounds()
		_ = bar.Add(1)

		if avgdata.data == nil {
			avgdata = newPictable(bounds.Max.X+2, bounds.Max.Y+2)
		}

		for y := bounds.Min.Y; y < bounds.Max.Y; y++ {
			for x := bounds.Min.X; x < bounds.Max.X; x++ {
				r, g, b, _ := m.At(x, y).RGBA()
				avgdata.add(x, y, r, g, b)
			}
		}
	}

	if err := png.Encode(f, avgImageFromPictable(avgdata, n)); err != nil {
		log.Fatal(err)
	}
}
