package orb

import (
	"gocv.io/x/gocv"
)

func MatchDescriptors(descA, descB gocv.Mat) []gocv.DMatch {
	// NORM_HAMMING should be used with ORB, BRISK and BRIEF,
	// NORM_HAMMING2 should be used with ORB when WTA_K==3 or 4
	normType := gocv.NormHamming
	if WTAK == 3 || WTAK == 4 {
		normType = gocv.NormHamming2
	}
	// we're using the ratio test so don't enable cross-checking
	crossCheck := false
	// brute force matcher takes each feature descriptor in the query matrix and computes the distance
	// for each descriptor in training, returning the closest K matches.
	matcher := gocv.NewBFMatcherWithParams(normType, crossCheck)
	// we want to return the top 2 matches for each then apply the Lowe's ratio test.
	// Lowe's test checks that the two distances are sufficiently different.
	// Background: we can't just return the closest match and expect things to work well, because
	// many points will be present in many images and therefore not be particularly useful for
	// _discriminating_ descriptors (effectively these points are background noise and should be ignored).
	// You could downweight these points and upweight novel points (TF-IDF style) but that's slow.
	// Instead, Lowe's ratio test is an ingenious way to do this by selecting 2 possible candidates
	// for the match. The assumption is that there is exactly 1 match in the other image, and the
	// other match is wrong/noise. The match with the smallest distance is the "good" match, and so
	// therefore the 2nd match is "bad" and just noise. If this is true, we can use the "bad" match
	// as a control and check that the "good" match is MUCH closer than the "bad" match. If it isn't
	// then this implies the "good" match isn't good at all, it's no better than the "bad" match and
	// hence should be ignored.
	matches := matcher.KnnMatch(descA, descB, 2) // look for A in B

	// apply ratio test
	goodMatches := make([]gocv.DMatch, 0, len(matches))
	for _, kMatches := range matches {
		if len(kMatches) != 2 {
			continue // possible to only have 1 match if we only have 1 key point
		}
		if kMatches[0].Distance < kMatches[1].Distance*0.75 {
			goodMatches = append(goodMatches, kMatches[0])
		}
	}
	return goodMatches
}
