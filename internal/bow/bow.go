package bow

type WordID uint64
type WordValue float64

type BowVector struct {
	// Sorted word IDs, allows binary search lookups
	wordIDs []WordID
	// indexes match word IDs
	wordVals []WordValue
}

func (v BowVector) AddWeight(id WordID, val WordValue) {
	existingVal, _ := v.find(id)
	val += existingVal
	v.set(id, val)
}

func (v BowVector) AddWeightIfNotExist(id WordID, val WordValue) {
	_, exists := v.find(id)
	if exists {
		return
	}
	v.set(id, val)
}

// L1-Normalizes the values in the vector
func (v BowVector) Normalise() {
	// TODO
}

func (v BowVector) set(id WordID, val WordValue) {
	// TODO
}

// Return the index of the next item that is >= than the provided word ID. Returns -1 if there is no such key.
func (v BowVector) lowerBound(w WordID) (i int) {
	low := 0
	hi := len(v.wordIDs) - 1
	var mid int
	for low <= hi && hi < len(v.wordIDs) {
		mid = int((low + hi) / 2)
		if v.wordIDs[mid] < w { // value was smaller than the one we want, search higher
			low = mid + 1
		} else if v.wordIDs[mid] > w { // value was higher than the one we want, search lower
			hi = mid - 1
		} else { // we have the value
			return mid
		}
	}
	return mid
}

func (v BowVector) find(id WordID) (val WordValue, ok bool) {
	low := 0
	hi := len(v.wordIDs) - 1
	for low <= hi && hi < len(v.wordIDs) {
		mid := int((low + hi) / 2)
		if v.wordIDs[mid] < id { // value was smaller than the one we want, search higher
			low = mid + 1
		} else if v.wordIDs[mid] > id { // value was higher than the one we want, search lower
			hi = mid - 1
		} else { // we have the value
			return v.wordVals[mid], true
		}
	}
	return 0, false
}
