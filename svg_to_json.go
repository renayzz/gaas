package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	svg "github.com/ajstarks/svgo"
)

type SVGElement struct {
	Type      string
	Attrs     map[string]string
	Text      string
	LinkedTo  *SVGElement
	LinkStyle string
	ArrowType string
}

func main() {
	// Check if SVG file path is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide an SVG file path")
		return
	}

	// Read SVG file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	data, err := os.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	// Parse SVG using svgo
	var elements []*SVGElement
	canvas := svg.New(&elements)
	canvas.SetParser(svg.XMLParse)
	canvas.SetIndent("", "  ")
	canvas.Parse(data)

	// Find linked elements
	for _, element := range elements {
		// Check if this element has a link
		if link, ok := element.Attrs["xlink:href"]; ok {
			// Find the linked element
			for _, linkedElement := range elements {
				if linkedElement.Attrs["id"] == link {
					element.LinkedTo = linkedElement
					element.LinkStyle = element.Attrs["style"]
					element.ArrowType = element.Attrs["marker-end"]
					break
				}
			}
		}
	}

	// Create a JSON array with the results
	var result []map[string]string

	for _, element := range elements {
		item := map[string]string{
			"type": element.Type,
			"text": element.Text,
		}

		if element.LinkedTo != nil {
			item["linked_to"] = element.LinkedTo.Type
			item["link_style"] = element.LinkStyle
			item["arrow_type"] = element.ArrowType
		}

		result = append(result, item)
	}

	// Output the JSON array to standard output
	jsonData, err := json.Marshal(result)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(jsonData))
}
