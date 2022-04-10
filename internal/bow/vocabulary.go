package bow

import (
	"bufio"
	"fmt"
	"io"
	"math"
	"os"
	"strconv"
	"strings"
)

type Node struct {
	ID         int
	Parent     int
	Children   []int
	Descriptor []uint8

	// only if the node is a word
	WordID WordID
	Weight WordValue
}

type Vocabulary struct {
	BranchingFactor int
	DepthLevels     int
	// words for this vocabulary (leaves)
	Words []*Node
	// tree nodes
	Nodes []Node
}

// Read a vocabulary file that is compatible with DBoW2 (only a subset of functionality is supported)
func NewVocabularyFromFile(filename string) (*Vocabulary, error) {
	file, err := os.Open(filename)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return NewVocabularyFromReader(file)
}

func NewVocabularyFromReader(reader io.Reader) (*Vocabulary, error) {
	// Example contents:
	// 10 6  0 0
	// 0 0 252 188 188 242 169 109 85 143 187 191 164 25 222 255 72 27 129 215 237 16 58 111 219 51 219 211 85 127 192 112 134 34  0
	// 0 0 93 125 221 103 180 14 111 184 112 234 255 76 215 115 153 115 22 196 124 110 233 240 249 46 237 239 101 20 104 243 66 33  0
	// 0 0 58 185 58 250 93 221 82 239 143 13 252 9 46 221 102 16 200 187 215 80 78 43 250 245 251 221 0 123 83 14 238 202  0
	scanner := bufio.NewScanner(reader)

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
	} else {
		return nil, fmt.Errorf("no lines read")
	}
	v := &Vocabulary{
		BranchingFactor: *branchingFactor,
		DepthLevels:     *depthLevels,
	}
	// subsequent lines are nodes, and not all nodes are words for the BoW vectors.
	expectedNodes := (math.Pow(float64(v.BranchingFactor), 1+float64(v.DepthLevels)) - 1) / (float64(v.BranchingFactor) - 1)
	expectedWords := math.Pow(float64(v.BranchingFactor), 1+float64(v.DepthLevels)) - 1
	nodes := make([]Node, 0, int(expectedNodes))
	words := make([]*Node, 0, int(expectedWords))
	// start with a root node
	var root Node
	nodes = append(nodes, root)

	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		n := Node{
			ID: len(nodes),
		}
		parentID, err := strconv.Atoi(fields[0])
		isLeaf := fields[1] != "0"
		if err != nil {
			return nil, fmt.Errorf("parsing parent ID failed: line %v -> %s", line, err)
		}
		parent := nodes[parentID]
		parent.Children = append(parent.Children, n.ID)
		nodes[parentID] = parent
		n.Parent = parentID

		// fields [2]->[n-1] are depth levels
		// TODO: is this right? We skip a bunch of fields here, 32 entries for ORB vectors, but depth is 6.
		descriptor := make([]uint8, *depthLevels)
		for i := 0; i < len(descriptor); i++ {
			desc, err := strconv.ParseUint(fields[2+i], 10, 8)
			if err != nil {
				return nil, fmt.Errorf("parsing depth %v at index %v failed: line %v -> %s", i, i+2, line, err)
			}
			descriptor[i] = uint8(desc)
		}
		n.Descriptor = descriptor

		// [n] is the weight
		weight, err := strconv.ParseFloat(fields[len(fields)-1], 64)
		if err != nil {
			return nil, fmt.Errorf("parsing weight failed: line %v -> %s", line, err)
		}
		n.Weight = WordValue(weight)

		if isLeaf {
			wordID := len(words)
			n.WordID = WordID(wordID)
			words = append(words, &n)
		} else if len(n.Children) == 0 {
			n.Children = make([]int, 0, *branchingFactor)
		}
		nodes = append(nodes, n)
	}
	v.Nodes = nodes
	v.Words = words

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return v, nil
}
