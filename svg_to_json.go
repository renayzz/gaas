package main

import (
	"encoding/json"
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
)

type Text struct {
	X        int    `xml:"x,attr"`
	Y        int    `xml:"y,attr"`
	FontSize int    `xml:"font-size,attr"`
	Fill     string `xml:"fill,attr"`
	Content  string `xml:",chardata"`
}

type Arrow struct {
	Start string `xml:"d,attr"`
	End   string `xml:"marker-end,attr"`
}

type Form struct {
	ID     string `json:"id"`
	X      int    `json:"x"`
	Y      int    `json:"y"`
	Width  int    `json:"width"`
	Height int    `json:"height"`
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

	var texts []TextJson
	var arrows []ArrowJson

	svgObj := svg.New(os.Stdout)
	svgObj.Start(500, 500)

	err = svgObj.Parse(svgFile)
	if err != nil {
		fmt.Println("Error parsing SVG:", err)
		return
	}

	for i, elem := range svgObj.Elements() {
		switch elem := elem.(type) {
		case svg.Text:
			form := Form{
				ID:     fmt.Sprintf("form%d", i),
				X:      elem.X,
				Y:      elem.Y,
				Width:  elem.FontSize,
				Height: elem.FontSize,
			}
			textJson := TextJson{
				ID:      fmt.Sprintf("text%d", i),
				Content: elem.Content,
				Form:    form,
			}
			texts = append(texts, textJson)
		case svg.Path:
			start := elem.D
			end := elem.MarkerEnd
			arrowJson := ArrowJson{
				Start: start,
				End:   end,
			}
			arrows = append(arrows, arrowJson)
		}
	}

	output := Output{
		Texts:  texts,
		Arrows: arrows,
	}
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
