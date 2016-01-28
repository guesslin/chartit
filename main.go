package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"github.com/ajstarks/svgo"
	"io"
	"log"
	"math"
	"os"
	"sort"
	"strconv"
	"strings"
)

// Chart is the simplest element only make up by label and value.
type Chart struct {
	Label string
	Value int
}

// Charts is slice of Chart
type Charts []Chart

// byValue is alias for sort.Sort
type byValue Charts

// Sum all element value in Charts
func (c Charts) Sum() (s int) {
	for _, col := range c {
		s += col.Value
	}
	return
}

// Percentage return giving label percent in Charts
func (c Charts) Percentage(label string) float64 {
	var numerator int
	for _, col := range c {
		if strings.Compare(col.Label, label) == 0 {
			numerator = col.Value
			break
		}
	}
	return float64(numerator) / float64(c.Sum())
}

func (b byValue) Len() int {
	return len(b)
}

func (b byValue) Swap(i, j int) {
	b[i], b[j] = b[j], b[i]
}
func (b byValue) Less(i, j int) bool {
	return b[i].Value < b[j].Value
}

// readCSV read from filename input by cmd flag, convert CSV into Charts struct and error if any file operand error happen.
func readCSV(filename string) (c Charts, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	r := csv.NewReader(file)
	records, err := r.ReadAll()
	if err != nil {
		return nil, err
	}
	for _, item := range records {
		if v, err := strconv.Atoi(item[1]); err == nil {
			col := Chart{item[0], v}
			c = append(c, col)
		}
	}
	return
}

func DrawBar(c Charts, width, height int, w io.Writer) {
}

// DrawPie draws pie chart from data by SVG PATH
func DrawPie(c Charts, width, height int, w io.Writer) {
	var angle float64
	clockwise := 1
	color := []string{"red", "blue", "black", "green"}
	r := int(math.Min(float64(width), float64(height)) * 0.7 / 2)
	// start Canvas with w io.Writer
	canvas := svg.New(w)
	canvas.Start(width, height)
	for i, col := range c {
		large := 0
		half := float64(360)*c.Percentage(col.Label)/2 + angle
		start := degreeToRadian(angle)
		angle += float64(360) * c.Percentage(col.Label)
		end := degreeToRadian(angle)
		if (end - start) >= math.Pi {
			large = 1
		}
		// Draw Sector from center of circle
		// then line to start angle point,
		// then Arc to end angle point with r raidus and clockwise direction
		// then line back to center of circle
		canvas.Path(fmt.Sprintf("M%d,%d L%d,%d A%d,%d 0 %d,%d %d,%d L%d,%d",
			(width/2),  // circle x
			(height/2), // circle y
			(width/2)+(int(math.Sin(start)*float64(r))),  // start angle x
			(height/2)-(int(math.Cos(start)*float64(r))), // start angle y
			r,     // r of circle
			r,     // r of circle
			large, // if over than 180 degree
			clockwise,
			(width/2)+(int(math.Sin(end)*float64(r))),  // end angle x
			(height/2)-(int(math.Cos(end)*float64(r))), // end angle y
			(width/2),   // circle x
			(height/2)), // circle y
			fmt.Sprintf("fill:%s;stroke:%s", color[i%4], color[i%4]))

		canvas.Text((width/2)+(int(math.Sin(degreeToRadian(half))*float64(r)*1.2)), (height/2)-(int(math.Cos(degreeToRadian(half))*float64(r)*1.2)), col.Label)
	}
	canvas.End()
}

// degreeToRadian convert degree into radian
func degreeToRadian(angle float64) (radian float64) {
	radian = math.Pi * angle / 180
	return
}

func main() {
	var csvFile string
	var chartFile string
	var width, height int
	var pie, bar bool
	flag.StringVar(&csvFile, "csv", "input.csv", "CSV filename")
	flag.StringVar(&chartFile, "output", "output", "OUTPUT filename without file extend")
	flag.IntVar(&width, "width", 1000, "OUTPUT file width")
	flag.IntVar(&height, "height", 800, "OUTPUT file height")
	flag.BoolVar(&pie, "pie", false, "Generate Pie Chart or not")
	flag.BoolVar(&bar, "bar", false, "Generate Bar Chart or not")
	flag.Parse()
	c, err := readCSV(csvFile)
	if err != nil {
		log.Fatal("Read csv file with some problem!")
	}
	sort.Sort(sort.Reverse(byValue(c)))
	if pie {
		out, err := os.Create(fmt.Sprintf("%s-pie.csv", chartFile))
		if err != nil {
			log.Fatal("Create file with some problem!")
		}
		defer out.Close()
		DrawPie(c, width, height, out)
	}
	if bar {
		out, err := os.Create(fmt.Sprintf("%s-bar.csv", chartFile))
		if err != nil {
			log.Fatal("Create file with some problem!")
		}
		defer out.Close()
		DrawBar(c, width, height, out)
	}
}
