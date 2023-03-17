package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"os"
)

type Svg struct {
	XMLName xml.Name `xml:"svg"`
	Texts   []Text   `xml:"text"`
	Arrows  []Arrow  `xml:"path"`
}

type Text struct {
	XMLName  xml.Name `xml:"text"`
	X        string   `xml:"x,attr"`
	Y        string   `xml:"y,attr"`
	FontSize string   `xml:"font-size,attr"`
	Fill     string   `xml:"fill,attr"`
	Content  string   `xml:",chardata"`
}

type Arrow struct {
	XMLName xml.Name `xml:"path"`
	Start   string   `xml:"d,attr"`
	End     string   `xml:"marker-end,attr"`
}

type Form struct {
	ID     string `json:"id"`
	X      string `json:"x"`
	Y      string `json:"y"`
	Width  string `json:"width"`
	Height string `json:"height"`
}

type TextJson struct {
	ID      string `json:"id"`
	Content string `json:"content"`
	Form    Form   `json:"form"`
}

type ArrowJson struct {
	Start string `json:"start"`
	End   string `json:"end"`
}

type Output struct {
	Texts  []TextJson  `json:"texts"`
	Arrows []ArrowJson `json:"arrows"`
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run svg_analyzer.go <input_file.svg>")
		return
	}

	inputFilePath := os.Args[1]
	svgFile, err := os.Open(inputFilePath)
	if err != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	defer svgFile.Close()

	svgBytes, err := ioutil.ReadAll(svgFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return
	}

	var svg Svg
	err = xml.Unmarshal(svgBytes, &svg)
	if err != nil {
		fmt.Println("Error parsing SVG:", err)
		return
	}

	var texts []TextJson
	for i, text := range svg.Texts {
		form := Form{ID: fmt.Sprintf("form%d", i), X: text.X, Y: text.Y, Width: text.FontSize, Height: text.FontSize}
		textJson := TextJson{ID: fmt.Sprintf("text%d", i), Content: text.Content, Form: form}
		texts = append(texts, textJson)
	}

	var arrows []ArrowJson
	for i, arrow := range svg.Arrows {
		start := arrow.Start
		end := arrow.End
		arrowJson := ArrowJson{Start: start, End: end}
		arrows = append(arrows, arrowJson)
	}

	output := Output{Texts: texts, Arrows: arrows}
	jsonBytes, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Println("Error generating JSON:", err)
		return
	}

	outputFilePath := inputFilePath + ".json"
	err = ioutil.WriteFile(outputFilePath, jsonBytes, 0644)
	if err != nil {
		fmt.Println("Error writing file:", err)
		return
	}

	fmt.Printf("JSON output written to %s\n", outputFilePath)
}
