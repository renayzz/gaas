package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/tdewolff/minify/v2"
	"github.com/tdewolff/minify/v2/svg"
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
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println(err)
		return
	}

	m := minify.New()
	m.AddFunc("image/svg+xml", svg.Minify)

	minifiedBytes, err := m.Bytes("image/svg+xml", bytes)
	if err != nil {
		fmt.Println(err)
		return
	}

	var elements []*SVGElement
	var currentElement *SVGElement

	lines := strings.Split(string(minifiedBytes), "\n")

	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}

		if strings.HasPrefix(line, "<") && strings.HasSuffix(line, ">") {
			tag := strings.TrimPrefix(line, "<")
			tag = strings.TrimSuffix(tag, ">")

			if strings.HasPrefix(tag, "/") {
				if currentElement != nil && strings.TrimPrefix(tag, "/") == currentElement.Type {
					currentElement = currentElement.LinkedTo
				}
			} else {
				parts := strings.Split(tag, " ")
				element := &SVGElement{
					Type:  parts[0],
					Attrs: make(map[string]string),
				}

				for _, part := range parts[1:] {
					kv := strings.Split(part, "=")
					if len(kv) == 2 {
						attrName := kv[0]
						attrValue := strings.Trim(kv[1], "\"")
						element.Attrs[attrName] = attrValue
					}
				}

				if currentElement != nil {
					element.LinkedTo = currentElement
					currentElement = element
				} else {
					currentElement = element
					elements = append(elements, element)
				}
			}
		} else {
			if currentElement != nil {
				currentElement.Text += line
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
			item["link_style"] = element.LinkedTo.LinkStyle
			item["arrow_type"] = element.LinkedTo.ArrowType
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
