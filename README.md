# gorbslam
Monocular ORB-SLAM in Go. Based on the [paper](https://arxiv.org/abs/1610.06475):
```
ORB-SLAM2: an Open-Source SLAM System for Monocular, Stereo and RGB-D Cameras
Raul Mur-Artal, Juan D. Tardos
arXiv:1610.06475 [cs.RO]
```

Currently you can run the following binaries:
 - `orbcheck`: run the ORB feature extractor on an image or on a video capture device.
 - `parsevocab`: parse an ORBvoc.txt file to create an ORB vocabulary.


### Requirements
 - OpenCV 4+ : `brew install opencv` or equivalent on Linux.

### Status
 - [x] ORB feature extration: https://www.youtube.com/watch?v=4AvTMVD9ig0 (describes SIFT but same principles apply)
 - [ ] Lens distortion correction / camera intrinsics: https://www.youtube.com/watch?v=26nV4oDLiqc
 - [x] ORB Vocabulary from file
 - [ ] Bag of Visual Words to store features for relocalisation (e.g during tracking failure): https://www.youtube.com/watch?v=a4cFONdc6nc 
 - [ ] Bundle adjustment / reprojection error: https://www.youtube.com/watch?v=lmj2Jk5tl60
 - [ ] Tracker goroutine to localise the camera with every frame (reprojection error and motion-only bundle adjustment): https://www.youtube.com/watch?v=0I30M6yTklo&t=191s
 - [ ] Local mapping goroutine to manage the local map and do local bundle adjustment: https://youtu.be/0I30M6yTklo?t=281
 - [ ] Loop closing goroutine using pose-graph optimisation (then full bundle adjustment). 

Long-term:
 - [ ] Stereo ORB-SLAM: fixes scale issues with Monocular and is generally more accurate.
