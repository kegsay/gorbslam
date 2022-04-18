package internal

import "gocv.io/x/gocv"

// Undistort the key points to account for lens distortion
// Learning: https://www.youtube.com/watch?v=26nV4oDLiqc
func UndistortKeyPoints(keypoints []gocv.KeyPoint) {
	// gocv.UndistortPoints(src, dst, cameraMatrix, distCoeffs, nil, nil)
}
