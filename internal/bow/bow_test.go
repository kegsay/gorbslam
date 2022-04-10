package bow

import (
	"math"
	"reflect"
	"testing"
)

func TestBow(t *testing.T) {
	var vec BowVector

	vec.AddWeight(30, 0.6)
	if !reflect.DeepEqual(vec.wordIDs, []WordID{30}) {
		t.Fatalf("wordID not stored correctly: %v", vec.wordIDs)
	}
	if !reflect.DeepEqual(vec.wordVals, []WordValue{0.6}) {
		t.Fatalf("wordVal not stored correctly: %v", vec.wordVals)
	}
	vec.AddWeight(10, 0.4)
	if !reflect.DeepEqual(vec.wordIDs, []WordID{10, 30}) {
		t.Fatalf("wordID not stored correctly: %v", vec.wordIDs)
	}
	if !reflect.DeepEqual(vec.wordVals, []WordValue{0.4, 0.6}) {
		t.Fatalf("wordVal not stored correctly: %v", vec.wordVals)
	}
	vec.AddWeight(20, 0.8)
	if !reflect.DeepEqual(vec.wordIDs, []WordID{10, 20, 30}) {
		t.Fatalf("wordID not stored correctly: %v", vec.wordIDs)
	}
	if !reflect.DeepEqual(vec.wordVals, []WordValue{0.4, 0.8, 0.6}) {
		t.Fatalf("wordVal not stored correctly: %v", vec.wordVals)
	}
	vec.AddWeight(10, 0.1) // = 0.5
	if !reflect.DeepEqual(vec.wordIDs, []WordID{10, 20, 30}) {
		t.Fatalf("wordIDs not stored correctly: %v", vec.wordIDs)
	}

	vec.AddWeightIfNotExist(100, 0.5)
	if !reflect.DeepEqual(vec.wordIDs, []WordID{10, 20, 30, 100}) {
		t.Fatalf("wordIDs not stored correctly: %v", vec.wordIDs)
	}
	vec.AddWeightIfNotExist(100, 0.2)
	if !reflect.DeepEqual(vec.wordIDs, []WordID{10, 20, 30, 100}) {
		t.Fatalf("wordIDs not stored correctly: %v", vec.wordIDs)
	}
	if !reflect.DeepEqual(vec.wordVals, []WordValue{0.5, 0.8, 0.6, 0.5}) {
		t.Fatalf("wordVals not stored correctly: %v", vec.wordVals)
	}

	// 0.5 + 0.8 + 0.6 + 0.5 = 2.4
	// Divide each entry by 2.4
	vec.Normalise()
	wantNorm := []WordValue{0.5 / 2.4, 0.8 / 2.4, 0.6 / 2.4, 0.5 / 2.4}
	for i := range wantNorm {
		if math.Abs(float64(vec.wordVals[i]-wantNorm[i])) > 0.0001 {
			t.Fatalf("index %d wordVal not normalised correctly: %v != %v", i, vec.wordVals[i], wantNorm[i])
		}
	}

}
