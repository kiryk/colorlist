package main

import (
	"fmt"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"sort"
)

type ColorFreq struct {
	C color.Color
	F int
}

type ByColorFreq []ColorFreq

func (cfs ByColorFreq) Len() int {
	return len(cfs)
}

func (cfs ByColorFreq) Swap(i, j int) {
	cfs[i], cfs[j] = cfs[j], cfs[i]
}

func (cfs ByColorFreq) Less(i, j int) bool {
	/*ir, ig, ib, _ := cfs[i].C.RGBA()
	jr, jg, jb, _ := cfs[j].C.RGBA()*/
	if cfs[i].F > cfs[j].F {
		return true
	} /*else if cfs[i].F == cfs[j].F && ir > jr {
		return true
	} else if ir == jr && ig > jg {
		return true
	} else if ig == jg && ib > jb {
		return true
	}*/
	return false
}

func ReadStats(r io.Reader) (map[color.Color]int, error) {
	stats := make(map[color.Color]int)
	img, err := png.Decode(r)

	if err != nil {
		return stats, err
	}

	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			stats[img.At(x, y)]++
		}
	}
	return stats, nil
}

func GetTopN(stats map[color.Color]int, n int) []ColorFreq {
	top := []ColorFreq{}

	for c, f := range stats {
		top = append(top, ColorFreq{c, f})
	}
	sort.Sort(ByColorFreq(top))
	return top
}

func PrintStats(w io.Writer, cf []ColorFreq) {
	fmt.Fprintf(w, "<!DOCTYPE HTML>\n<html>"+
		"<head><style>td {min-width: 30pt;font-family: monospace;}</style></head>"+
		"<body><table><tr>\n")
	for i := range cf {
		if (i+1)%10 == 0 {
			fmt.Fprintf(w, "</tr><tr>\n")
		}
		r, g, b, _ := cf[i].C.RGBA()
		r /= 0x100
		g /= 0x100
		b /= 0x100
		fmt.Fprintf(w,
			"\t<td style=\"background-color: rgb(%d,%d,%d);width: 10pt;\">%d</td>\n",
			r, g, b, cf[i].F)
	}
	fmt.Fprintf(w, "</tr></table></body></html>\n")
}

func main() {
	var err error
	var r io.Reader
	var w io.Writer

	if len(os.Args) == 2 {
		if r, err = os.Open(os.Args[1]); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
		if w, err = os.Create(os.Args[1] + ".html"); err != nil {
			log.Fatal(err)
			os.Exit(1)
		}
	} else {
		r, w = os.Stdin, os.Stdout
	}
	stats, err := ReadStats(r)
	if err != nil {
		log.Fatal(err)
		os.Exit(1)
	}
	PrintStats(w, GetTopN(stats, len(stats)))
}
