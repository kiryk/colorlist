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
		<tr><td style="background-color: {{.RGB}};width: 10pt;">{{.Freq}}</td></tr>
		{{end}}
		</table>
	</body>
</html>
`))

type byFreq []ColorStats
type ColorStats struct {
	Color color.Color
	Freq int
}

func (cs byFreq) Len() int           { return len(cs) }
func (cs byFreq) Less(i, j int) bool { return cs[i].Freq < cs[j].Freq }
func (cs byFreq) Swap(i, j int)      { cs[i], cs[j] = cs[j], cs[i] }

func (c *ColorStats) RGB() string {
	r, g, b, _ := c.Color.RGBA()
	return fmt.Sprintf("#%02x%02x%02x", r>>8, g>>8, b>>8)
}

func PrintStats(w io.Writer, stats []ColorStats) {
	tmpl.Execute(w, stats)
}

func ReadStats(r io.Reader) ([]ColorStats, error) {
	img, err := png.Decode(r)
	if err != nil {
		return []ColorStats{}, err
	}

	counts := make(map[color.Color]int)
	for y := img.Bounds().Min.Y; y < img.Bounds().Max.Y; y++ {
		for x := img.Bounds().Min.X; x < img.Bounds().Max.X; x++ {
			counts[img.At(x, y)]++
		}
	}

	stats := []ColorStats{}
	for color, freq := range counts {
		stats = append(stats, ColorStats{color, freq})
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
