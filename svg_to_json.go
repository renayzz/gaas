package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"os"
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

	// Parse SVG file
	var svg struct {
		XMLName xml.Name `xml:"svg"`
		Elements []struct {
			XMLName xml.Name
			Text    string `xml:",chardata"`
			Attrs   []struct {
				Name  xml.Name
				Value string
			} `xml:",any"`
			LinkedTo  string `xml:"http://www.w3.org/1999/xlink href,attr"`
			LinkStyle string `xml:"style"`
			ArrowType string `xml:"marker-end,attr"`
		} `xml:",any"`
	}

	decoder := xml.NewDecoder(file)
	if err := decoder.Decode(&svg); err != nil {
		fmt.Println(err)
		return
	}

	// Convert parsed SVG elements to our own SVGElement type
	var elements []*SVGElement
	for _, elem := range svg.Elements {
		element := &SVGElement{
			Type:  elem.XMLName.Local,
			Attrs: make(map[string]string),
			Text:  elem.Text,
			LinkStyle: elem.LinkStyle,
			ArrowType: elem.ArrowType,
		}

		for _, attr := range elem.Attrs {
			element.Attrs[attr.Name.Local] = attr.Value
		}

		elements = append(elements, element)
	}

	// Find linked elements
	for _, element := range elements {
		if element.LinkedTo != "" {
			for _, linkedElement := range elements {
				if linkedElement.Attrs["id"] == element.LinkedTo {
					element.LinkedTo = linkedElement
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
