package internal

import (
	"gocv.io/x/gocv"
)

// ORB feature extractor variables TODO: configurable
const (
	// The maximum number of features to retain.
	nFeatures        = 1000 // opencv default is 500; we like features so let's get more
	nInitialFeatures = 5000 // we need more features on the first few frames to speed up bootstrapping process
	// Pyramid decimation ratio, greater than 1. scaleFactor==2 means the classical pyramid, where
	// each next level has 4x less pixels than the previous, but such a big scale factor will degrade
	// feature matching scores dramatically. On the other hand, too close to 1 scale factor will mean
	// that to cover certain scale range you will need more pyramid levels and so the speed will suffer.
	scaleFactor = 1.2 // opencv default.
	// The number of pyramid levels. The smallest level will have linear size equal to input_image_linear_size/pow(scaleFactor, nlevels - firstLevel).
	nLevels = 8 // opencv default.
	// FAST threshold
	fastThreshold = 20 // opencv default. TODO: if no corners detected, use lower value?
	// The number of points that produce each element of the oriented BRIEF descriptor.
	// The default value 2 means the BRIEF where we take a random point pair and compare their brightnesses, so we get 0/1 response.
	WTAK = 2 // opencv default
	// The level of pyramid to put source image to.
	firstLevel = 0 // opencv default
	// This is size of the border where the features are not detected. It should roughly match the patchSize parameter.
	edgeThreshold = 31 // opencv default
	// size of the patch used by the oriented BRIEF descriptor. Of course, on smaller pyramid layers the perceived image area covered by a feature will be larger.
	patchSize = 31 // opencv default
	// The default HARRIS_SCORE means that Harris algorithm is used to rank features (the score is
	// written to KeyPoint::score and is used to retain best nfeatures features); FAST_SCORE is
	// alternative value of the parameter that produces slightly less stable keypoints, but it is a little faster to compute.
	scoreType = gocv.ORBScoreTypeHarris // opencv default
)

// Compute the feature vectors for the input image.
// Learning: https://www.youtube.com/watch?v=4AvTMVD9ig0 for more information (it describes SIFT but works similarly for ORB too)
// - Detect keypoints in the image
// - Compute descriptors around each keypoint
func OrbFeatures(img gocv.Mat, initial bool) ([]gocv.KeyPoint, gocv.Mat) {
	features := nFeatures
	if initial {
		features = nInitialFeatures
	}
	orb := gocv.NewORBWithParams(
		features, scaleFactor, nLevels, edgeThreshold, firstLevel, WTAK, scoreType, patchSize, fastThreshold,
	)
	// Masks specifying where to look for keypoints
	mask := gocv.NewMatWithSizeFromScalar(gocv.NewScalar(1, 1, 1, 1), img.Rows(), img.Cols(), gocv.MatTypeCV8U)
	keypoints, descriptors := orb.DetectAndCompute(img, mask)
	// TODO: undistort
	// TODO: assign features to grid?
	return keypoints, descriptors
}

// Undistort the key points to account for lens distortion
// Learning: https://www.youtube.com/watch?v=26nV4oDLiqc
func UndistortKeyPoints(keypoints []gocv.KeyPoint) {
	// gocv.UndistortPoints(src, dst, cameraMatrix, distCoeffs, nil, nil)
}
