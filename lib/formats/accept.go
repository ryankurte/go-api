package formats

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
)

// Regex for matching / grouping accept headers
var acceptRegex = regexp.MustCompile(`((?P<name>[a-z\-]+\/[a-z\-]+)(; q=(?P<weight>[0-1].[0-9+])){0,1}[, ]*)`)
var groupNames = acceptRegex.SubexpNames()

// Accept data type for internal use
type acceptType struct {
	name   string
	weight float32
}

type acceptTypes []acceptType

func (slice acceptTypes) Len() int {
	return len(slice)
}

func (slice acceptTypes) Less(i, j int) bool {
	return slice[i].weight < slice[j].weight
}

func (slice acceptTypes) Swap(i, j int) {
	slice[i], slice[j] = slice[j], slice[i]
}

// Parses an HTTP accept header to return an ordered list of content types
func ParseAcceptHeader(accept string) ([]string, error) {
	result := make([]string, 0)

	var types acceptTypes

	// Check for match
	if res := acceptRegex.MatchString(accept); res == false {
		return result, fmt.Errorf("Invalid accept header")
	}

	// Locate all accept types
	for _, match := range acceptRegex.FindAllStringSubmatch(accept, -1) {
		t := match[2]
		w, err := strconv.ParseFloat(match[4], 32)
		if err != nil {
			w = 1.0
		}

		types = append(types, acceptType{t, float32(w)})
	}

	// Sort data by weight (assumes equal ordering will be maintained)
	sort.Sort(sort.Reverse(types))

	// Pack into ordered result array
	for _, t := range types {
		result = append(result, t.name)
	}

	return result, nil
}
