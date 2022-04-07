package bow

import (
	"reflect"
	"testing"
)

func TestBow(t *testing.T) {
	var vec BowVector

	vec.AddWeight(30, 0.6)
	vec.AddWeight(10, 0.4)
	vec.AddWeight(20, 0.8)
	vec.AddWeight(10, 0.1) // = 0.5

	vec.AddWeightIfNotExist(100, 0.5)
	vec.AddWeightIfNotExist(100, 0.2)
	if !reflect.DeepEqual(vec.wordIDs, []int{10, 20, 30, 100}) {
		t.Fatalf("wordIDs not stored correctly: %v", vec.wordIDs)
	}
	if !reflect.DeepEqual(vec.wordVals, []float64{0.5, 0.8, 0.6, 0.5}) {
		t.Fatalf("wordVals not stored correctly: %v", vec.wordVals)
	}

}
