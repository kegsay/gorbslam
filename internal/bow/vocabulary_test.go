package bow

import (
	"bytes"
	_ "embed"
	"testing"
)

//go:embed ORBShortvoc.txt
var VocabFile []byte

func TestVocabulary(t *testing.T) {
	voc, err := NewVocabularyFromReader(bytes.NewReader(VocabFile))
	if err != nil {
		t.Fatalf("NewVocabularyFromReader: %s", err)
	}
	if voc.BranchingFactor != 10 {
		t.Errorf("Branching factor wrong, got %v want 10", voc.BranchingFactor)
	}
	if voc.DepthLevels != 6 {
		t.Errorf("Depth levels wrong, got %v want 6", voc.DepthLevels)
	}
	if len(voc.Nodes) != 303 {
		t.Errorf("num Nodes wrong, got %v want 303", len(voc.Nodes))
	}
}
