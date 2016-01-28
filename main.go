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

// DrawPie draws pie chart from data by SVG PATH
func DrawPie(c Charts, width, height int, w io.Writer) (err error) {
	var angle float64
	canvas := svg.New(w)
	color := []string{"red", "blue", "black", "green"}
	r := width * 3 / 15
	canvas.Start(width, height)
	canvas.Circle(width/2, height/2, r, "fill:none;stroke:red")
	for i, col := range c {
		large := 0
		half := float64(360)*c.Percentage(col.Label)/2 + angle
		start := degreeToRadian(angle)
		angle += float64(360) * c.Percentage(col.Label)
		end := degreeToRadian(angle)
		if (end - start) >= math.Pi {
			large = 1
		}
		canvas.Path(fmt.Sprintf("M%d,%d L%d,%d A%d,%d 0 %d,1 %d,%d L%d,%d",
			(width/2),  // circle x
			(height/2), // circle y
			(width/2)+(int(math.Sin(start)*float64(r))),  // start angle x
			(height/2)-(int(math.Cos(start)*float64(r))), // start angle y
			r,     // r of circle
			r,     // r of circle
			large, // if over than 180 degree
			(width/2)+(int(math.Sin(end)*float64(r))),  // end angle x
			(height/2)-(int(math.Cos(end)*float64(r))), // end angle y
			(width/2),   // circle x
			(height/2)), // circle y
			fmt.Sprintf("fill:%s;stroke:%s", color[i%4], color[i%4]))
		canvas.Text((width/2)+(int(math.Sin(degreeToRadian(half))*float64(r+50))), (height/2)-(int(math.Cos(degreeToRadian(half))*float64(r+50))), col.Label)
	}
	canvas.End()
	return
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
	flag.StringVar(&csvFile, "csv", "input.csv", "CSV filename")
	flag.StringVar(&chartFile, "output", "output.svg", "OUTPUT filename")
	flag.IntVar(&width, "width", 1000, "OUTPUT file width")
	flag.IntVar(&height, "height", 800, "OUTPUT file height")
	flag.Parse()
	c, err := readCSV(csvFile)
	if err != nil {
		log.Fatal("Read csv file with some problem!")
	}
	out, err := os.Create(chartFile)
	if err != nil {
		log.Fatal("Create file with some problem!")
	}
	if err = DrawPie(c, width, height, out); err != nil {
		log.Fatal("Write canvas file with some problem!")
	}
}
