package bow

import "testing"

func TestL1Scorer(t *testing.T) {
	var s L1Scorer

	score := s.Score(BowVector{
		wordIDs:  []WordID{1, 2, 3, 4, 5, 6},
		wordVals: []WordValue{0.2, 0.3, 0.5, 0.1, 0.8, 0.9},
	}, BowVector{
		wordIDs:  []WordID{1, 3, 5, 6, 7, 8, 12},
		wordVals: []WordValue{0.7, 0.1, 0.2, 0.6, 0.1, 0.01},
	})
	// we should do abs(a,b) - abs(a) - abs(b) for each matching word ID then do -1 * score / 2
	// matching word IDs are 1,3,5,6
	// A       => 0.2, 0.5, 0.8, 0.9
	// B       => 0.7, 0.1, 0.2, 0.6
	// abs(a,b)=> 0.5, 0.4, 0.6, 0.3
	// -abs(a) => 0.3, -0.1 -0.2 -0.6
	// -abs(b) => -0.4 -0.2 -0.4 -1.2
	// score = -2.2
	// -score/2 = 1.1
	if score != 1.1 {
		t.Fatalf("got score %v want 1.1", score)
	}
}
