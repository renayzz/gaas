package main

import (
	"encoding/xml"
	"fmt"
	"os"
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

	// Parse SVG file
	decoder := xml.NewDecoder(file)
	var elements []*SVGElement
	var currentElement *SVGElement

	for {
		// Read tokens from the XML document
		token, err := decoder.Token()
		if err != nil {
			break
		}

		// Process the token
		switch t := token.(type) {
		case xml.StartElement:
			// Create a new SVG element
			currentElement = &SVGElement{
				Type:  t.Name.Local,
				Attrs: map[string]string{},
			}

			// Store the element attributes
			for _, attr := range t.Attr {
				currentElement.Attrs[attr.Name.Local] = attr.Value
			}

			// Check if this element has a link
			if link, ok := currentElement.Attrs["xlink:href"]; ok {
				// Find the linked element
				for _, linkedElement := range elements {
					if linkedElement.Attrs["id"] == link {
						currentElement.LinkedTo = linkedElement
						currentElement.LinkStyle = currentElement.Attrs["style"]
						break
					}
				}
			}

			// Add the element to the slice
			elements = append(elements, currentElement)
		}
	}

	// Print the results
	for _, element := range elements {
		fmt.Printf("%s element found\n", element.Type)

		if element.LinkedTo != nil {
			fmt.Printf("  Linked to %s element with style %s\n", element.LinkedTo.Type, element.LinkStyle)
		}
	}
}
