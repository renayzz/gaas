package main

import (
	"fmt"
	"os"

	svg "github.com/ajstarks/svgo"
)

type SVGElement struct {
	Type      string
	Attrs     map[string]string
	LinkedTo  *SVGElement
	LinkStyle string
}

func main() {
	// Check if SVG file path is provided
	if len(os.Args) < 2 {
		fmt.Println("Please provide an SVG file path")
		return
	}

	// Open SVG file
	file, err := os.Open(os.Args[1])
	if err != nil {
		fmt.Println(err)
		return
	}
	defer file.Close()

	// Create new SVG decoder
	decoder := svg.NewDecoder(file)

	// Parse SVG image
	image, err := decoder.Decode()
	if err != nil {
		fmt.Println(err)
		return
	}

	// Create a slice of SVG elements
	elements := []*SVGElement{}

	// Iterate over SVG elements
	for _, element := range image.Children {
		svgElement := &SVGElement{
			Type:  element.XMLName.Local,
			Attrs: element.Attr,
		}

		// Check if this element has a link
		if link, ok := element.Attr["xlink:href"]; ok {
			// Find the linked element
			for _, linkedElement := range elements {
				if linkedElement.Attrs["id"] == link {
					svgElement.LinkedTo = linkedElement
					linkedElement.LinkStyle = element.Attr["style"]
					break
				}
			}
		}

		// Add the element to the slice
		elements = append(elements, svgElement)
	}

	// Print the results
	for _, element := range elements {
		fmt.Printf("%s element found\n", element.Type)

		if element.LinkedTo != nil {
			fmt.Printf("  Linked to %s element with style %s\n", element.LinkedTo.Type, element.LinkStyle)
		}
	}
}