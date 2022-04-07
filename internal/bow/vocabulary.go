package bow

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// Read a vocabulary file that is compatible with DBoW2 (only a subset of functionality is supported)
func NewVocabularyFromFile(filename string) ([]string, error) {
	// Example contents:
	// 10 6  0 0
	// 0 0 252 188 188 242 169 109 85 143 187 191 164 25 222 255 72 27 129 215 237 16 58 111 219 51 219 211 85 127 192 112 134 34  0
	// 0 0 93 125 221 103 180 14 111 184 112 234 255 76 215 115 153 115 22 196 124 110 233 240 249 46 237 239 101 20 104 243 66 33  0
	// 0 0 58 185 58 250 93 221 82 239 143 13 252 9 46 221 102 16 200 187 215 80 78 43 250 245 251 221 0 123 83 14 238 202  0
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	scanner := bufio.NewScanner(file)

	// The first line contains the following metadata, grab it.
	var branchingFactor *int // aka 'k'
	var depthLevels *int     // aka 'L'
	var scoringType *int
	var weightingType *int

	if scanner.Scan() {
		firstLine := scanner.Text()
		fields := strings.Fields(firstLine)
		if len(fields) != 4 {
			return nil, fmt.Errorf("invalid vocab file, expected %d fields on first line for metadata", len(fields))
		}
		entries := []struct {
			val **int
			str string
			min int
			max int
		}{
			{
				val: &branchingFactor, min: 0, max: 20, str: fields[0],
			},
			{
				val: &depthLevels, min: 1, max: 10, str: fields[1],
			},
			{
				// enums from DBoW2: we only implement L1_NORM
				// L1_NORM=0, L2_NORM=1, CHI_SQUARE=2, KL=3, BHATTACHARYYA=4, DOT_PRODUCT=5
				val: &scoringType, min: 0, max: 0, str: fields[2],
			},
			{
				// enums from DBoW2: we only implement TF_IDF
				// TF_IDF=0, TF=1, IDF=2, BINARY=3
				val: &weightingType, min: 0, max: 0, str: fields[3],
			},
		}
		for _, entry := range entries {
			parsedInt, err := strconv.ParseInt(entry.str, 10, 64)
			if err != nil {
				return nil, fmt.Errorf("unable to parse metadata field value '%s' - expected integer: %s", entry.str, err)
			}
			pint := int(parsedInt)
			if pint > entry.max || pint < entry.min {
				return nil, fmt.Errorf("metadata field unknown, value '%d'", pint)
			}
			*entry.val = &pint
		}
	}

	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println(line)
		// TODO parse line
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return nil, nil
}
