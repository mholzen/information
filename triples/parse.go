package triples

import (
	"bufio"
	"bytes"
	"encoding/json"
	"log"
	"os"
	"strings"
)

// parse examples.json into a Go data structure

func Parse(filename string) (*Triples, error) {
	// read the file
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var buffer bytes.Buffer

	// remove comments
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		// remove comments from line
		line = strings.Split(line, "//")[0]

		buffer.WriteString(line)
	}
	return parseString(buffer)
}

func parseString(buffer bytes.Buffer) (*Triples, error) {
	decoder := json.NewDecoder(&buffer)
	var data interface{}
	err := decoder.Decode(&data)
	if err != nil {
		return nil, err
	}

	res := NewTriples()

	switch typedData := data.(type) {
	case map[string]interface{}:
		_, err := res.NewTriplesFromMap(typedData)
		if err != nil {
			return nil, err
		}

	case []interface{}:
		_, err := res.NewTriplesFromSlice(typedData)
		if err != nil {
			return nil, err
		}

	default:
		log.Printf("unknown")
	}

	printAllTriples(res)

	return res, nil
}

// func printKeys(triples *Triples) {
// 	keys := maps.Keys(triples.Nodes)
// 	sort.Strings(keys)

// 	log.Printf("keys (%d):", len(keys))
// 	// for _, key := range keys {
// 	// 	log.Printf("%s", key)
// 	// }
// }

func printAllTriples(triples *Triples) {
	log.Printf("triples (%d):", len(triples.TripleSet))
	for triple := range triples.TripleSet {
		log.Printf("%s", triple)
	}
}
