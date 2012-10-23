package main

import (
	"fmt"
	"os"
	"log"
	"image"
	"flag"
	"strings"
	"path/filepath"
)

import _ "image/png"


func Pic(dx int, dy int) [][][]uint64 {
	pic := make([][][]uint64, dx) /* type declaration */
	for i := range pic {
		pic[i] = make([][]uint64, dy) /* again the type? */
		for j := range pic[i] {
			pic[i][j] = []uint64{0, 0, 0}
		}
	}
	return pic
	//return pic
}

func main() {
	flag.Parse()
	dirname := flag.Arg(0) + string(filepath.Separator)
	fmt.Println(dirname)

	d, err := os.Open(dirname)
	if err != nil {
		log.Fatal(err)
	}

	fi, err := d.Readdir(-1)
	if err != nil {
		log.Fatal(err)
	}

	avgdata := [][][]uint64{}
	picinit := false;

	n := 0

	for _, fi := range fi {
		fname := fi.Name()
		if !fi.IsDir() && fname[0] != '.' && strings.HasSuffix(fname, ".png") {
			fmt.Println("Loading", fname)

			file, err := os.Open(dirname + fname)
			if err != nil {
				log.Fatal(err)
			}

			m, _, err := image.Decode(file)
			if err != nil {
				log.Fatal(err)
			}
			bounds := m.Bounds()

			if !picinit {
				avgdata = Pic(bounds.Max.X + 2, bounds.Max.Y + 2)
				picinit = true;
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

			n++
			if(n > 3) {
				break;
			}
		}
	}

	img := image.NewRGBA(image.Rect(0, 0, len(avgdata), len(avgdata[0])));
	_ = img

}