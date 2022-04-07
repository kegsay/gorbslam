package bow

import "math"

type Scorer interface {
	// Computes the score between two vectors. Vectors must be sorted and normalised if needed.
	Score(a, b []BowVector) float64
}

type L1Scorer struct{}

func (s *L1Scorer) Score(a, b BowVector) float64 {
	var score float64
	var i, j int
	for i < len(a.wordIDs) && j < len(b.wordIDs) {
		aVal := float64(a.wordVals[i])
		bVal := float64(b.wordVals[j])
		if a.wordIDs[i] == b.wordIDs[j] {
			score += math.Abs(aVal-bVal) - math.Abs(aVal) - math.Abs(bVal)
			i++
			j++
		} else if a.wordIDs[i] < b.wordIDs[j] {
			// move A forward
			i = a.lowerBound(b.wordIDs[j])
		} else {
			// move B forward
			j = b.lowerBound(a.wordIDs[i])
		}
	}
	score = (-1 * score) / 2
	return score
}
