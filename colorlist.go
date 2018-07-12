package main

import (
	"flag"
	"fmt"
	"html/template"
	"image/color"
	"image/png"
	"io"
	"log"
	"os"
	"sort"
)

var (
	input   = flag.String("f", "", "input file (default stdin)")
	reverse = flag.Bool("r", false, "reverse sorting")
)

var tmpl = template.Must(template.New("output").Parse(`
<!DOCTYPE html>
<html>
	<head><style>td {min-width: 30pt; font-family: monospace;}</style></head>
	<body>
		<table>
		{{range .}}
			<tr><td style="background-color: {{.RGB}};width: 10pt;">{{.F}}</td></tr>
		{{end}}
		</table>
	</body>
</html>
`))

type byFreq []ColorStats
type ColorStats struct {
	C color.Color
	F int
}

func (cfs byFreq) Len() int           { return len(cfs) }
func (cfs byFreq) Less(i, j int) bool { return cfs[i].F < cfs[j].F }
func (cfs byFreq) Swap(i, j int)      { cfs[i], cfs[j] = cfs[j], cfs[i] }

func (c *ColorStats) RGB() string {
	r, g, b, _ := c.C.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r/0x100, g/0x100, b/0x100)
}

func PrintStats(w io.Writer, stats []ColorStats) {
	tmpl.Execute(w, stats)
}

func ReadStats(r io.Reader) ([]ColorStats, error) {
	img, err := png.Decode(r)
	if err != nil {
		return stats, err
	}

	counts := make(map[color.Color]int)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			counts[img.At(x, y)]++
		}
	}

	stats := []ColorStats{}
	for c, f := range counts {
		stats = append(stats, ColorStats{c, f})
	}
	return stats, nil
}

func main() {
	var err error

	flag.Parse()

	r, w := os.Stdin, os.Stdout
	if *input != "" {
		if r, err = os.Open(*input); err != nil {
			log.Fatal(err)
		}
		defer r.Close()
		if w, err = os.Create(*input+".html"); err != nil {
			log.Fatal(err)
		}
		defer w.Close()
	}

	stats, err := ReadStats(r)
	if err != nil {
		log.Fatal(err)
	}

	if *reverse {
		sort.Sort(sort.Reverse(byFreq(stats)))
	} else {
		sort.Sort(byFreq(stats))
	}
	PrintStats(w, stats)
}
