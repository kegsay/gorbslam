package main

import (
	"flag"
	"image/color"
	"log"
	"os"
	"os/signal"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/kegsay/gorbslam/internal/orb"
	"gocv.io/x/gocv"
)

// Matches ORB keypoints between 2 images

var (
	flagImg1               = flag.String("img1", "", "First image")
	flagImg2               = flag.String("img2", "", "Second image")
	flagVideoCaptureDevice = flag.Int("cam", 0, "video capture device to read from")
)

func matchWebcam(cap int) {
	a := app.New()
	w := a.NewWindow("ORB Matches [video]")
	windowCanvas := w.Canvas()
	go func() {
		resized := false
		vc, err := gocv.VideoCaptureDevice(cap)
		if err != nil {
			log.Printf("failed to load video capture device: %s", err)
			os.Exit(1)
		}
		c := make(chan os.Signal, 1)
		signal.Notify(c, os.Interrupt)
		go func() {
			for range c {
				//vc.Close()
				w.Close()
				close(c)
			}
		}()
		var prevFrame *gocv.Mat
		var prevKPs []gocv.KeyPoint
		var prevDescriptors gocv.Mat
		for {
			select {
			case <-c:
				return
			default:
			}
			if !vc.IsOpened() {
				log.Println("waiting for video capture device...")
				time.Sleep(100 * time.Millisecond)
				continue
			}
			frame := gocv.NewMat()
			if !vc.Read(&frame) {
				log.Println("failed to read frame")
			}
			if !resized {
				w.Resize(fyne.NewSize(float32(frame.Cols()*2), float32(frame.Rows())))
				resized = true
			}
			if frame.Empty() {
				log.Printf("cannot load image")
				os.Exit(1)
			}

			kpFrame, dFrame := orb.Features(frame, false)
			if len(kpFrame) == 0 {
				continue
			}
			if prevFrame == nil {
				prevFrame = &frame
				prevDescriptors = dFrame
				prevKPs = kpFrame
				continue
			}
			// compare frame and prevFrame
			matches := orb.MatchDescriptors(dFrame, prevDescriptors)
			log.Printf("found %d matches", len(matches))
			if len(matches) > 0 {
				outputMat := gocv.NewMatWithSize(frame.Rows()*2, frame.Cols()*2, frame.Type())
				gocv.DrawMatches(frame, kpFrame, *prevFrame, prevKPs, matches, &outputMat, color.RGBA{G: 255}, color.RGBA{B: 255}, nil, gocv.NotDrawSinglePoints)
				outputImg, err := outputMat.ToImage()
				if err != nil {
					log.Printf("cannot draw matches: %s", err)
					continue
				}
				outputImage := canvas.NewImageFromImage(outputImg)
				outputImage.ScaleMode = canvas.ImageScaleFastest
				windowCanvas.SetContent(outputImage)
			}

			prevFrame = &frame
			prevDescriptors = dFrame
			prevKPs = kpFrame
		}
	}()

	w.Resize(fyne.NewSize(640, 480))
	w.ShowAndRun()
}

func main() {
	flag.Parse()
	if *flagImg1 == "" && *flagImg2 == "" && flagVideoCaptureDevice != nil {
		matchWebcam(*flagVideoCaptureDevice)
		return
	}

	if *flagImg1 == "" && *flagImg2 == "" {
		flag.Usage()
		os.Exit(1)
	}
	img1 := gocv.IMRead(*flagImg1, gocv.IMReadAnyColor)
	if img1.Empty() {
		log.Printf("cannot load image at %s", *flagImg1)
		os.Exit(1)
	}
	img2 := gocv.IMRead(*flagImg2, gocv.IMReadAnyColor)
	if img2.Empty() {
		log.Printf("cannot load image at %s", *flagImg2)
		os.Exit(1)
	}
	keypoints1, descriptors1 := orb.Features(img1, false)
	log.Printf("image 1: detected %d keypoints (%d x %d)", len(keypoints1), img1.Rows(), img1.Cols())
	keypoints2, descriptors2 := orb.Features(img2, false)
	log.Printf("image 2: detected %d keypoints (%d x %d)", len(keypoints2), img2.Rows(), img2.Cols())
	matches := orb.MatchDescriptors(descriptors1, descriptors2)
	log.Printf("found %d matches", len(matches))
	outputMat := gocv.NewMatWithSize(img1.Rows()+img2.Rows(), img1.Cols()+img2.Cols(), img1.Type())
	gocv.DrawMatches(img1, keypoints1, img2, keypoints2, matches, &outputMat, color.RGBA{G: 255}, color.RGBA{B: 255}, nil, gocv.NotDrawSinglePoints)
	outputImg, err := outputMat.ToImage()
	if err != nil {
		log.Printf("failed to convert output to image: %s", err)
		os.Exit(1)
	}
	a := app.New()
	w := a.NewWindow("ORB Matches")
	outputImage := canvas.NewImageFromImage(outputImg)
	w.SetContent(outputImage)
	w.Resize(fyne.NewSize(float32(img1.Cols()), float32(img1.Rows())))
	w.ShowAndRun()
}
