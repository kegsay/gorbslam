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
	"github.com/kegsay/gorbslam/internal"
	"gocv.io/x/gocv"
)

// Checks that the ORB features are being extracted from images correctly.

var (
	flagImg                = flag.String("image", "", "image to plot ORB features onto")
	flagVideoCaptureDevice = flag.Int("cam", 0, "video capture device to read from")
)

func main() {
	flag.Parse()
	// run it on an image
	if flagImg != nil && *flagImg != "" {
		img := gocv.IMRead(*flagImg, gocv.IMReadAnyColor)
		if img.Empty() {
			log.Printf("cannot load image at %s", *flagImg)
			os.Exit(1)
		}
		keypoints, _ := internal.OrbFeatures(img, false)
		log.Printf("Detected %d keypoints", len(keypoints))
		outputMat := gocv.NewMatWithSize(img.Rows(), img.Cols(), img.Type())
		gocv.DrawKeyPoints(img, keypoints, &outputMat, color.RGBA{G: 255}, gocv.DrawDefault)
		outputImg, err := outputMat.ToImage()
		if err != nil {
			log.Printf("failed to convert output to image: %s", err)
			os.Exit(1)
		}
		a := app.New()
		w := a.NewWindow("ORB Features")

		outputImage := canvas.NewImageFromImage(outputImg)
		w.SetContent(outputImage)
		w.Resize(fyne.NewSize(640, 480))
		w.ShowAndRun()
	}
	// run it on a webcam
	if flagVideoCaptureDevice != nil {
		a := app.New()
		w := a.NewWindow("ORB Features [video]")
		windowCanvas := w.Canvas()

		go func() {
			resized := false
			vc, err := gocv.VideoCaptureDevice(*flagVideoCaptureDevice)
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
					w.Resize(fyne.NewSize(float32(frame.Cols()), float32(frame.Rows())))
					resized = true
				}

				//frame := gocv.IMRead(*flagImg, gocv.IMReadAnyColor)
				if frame.Empty() {
					log.Printf("cannot load image at %s", *flagImg)
					os.Exit(1)
				}

				keypoints, _ := internal.OrbFeatures(frame, false)
				log.Printf("Detected %d keypoints", len(keypoints))
				if len(keypoints) == 0 {
					continue
				}
				outputMat := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), frame.Type())
				gocv.DrawKeyPoints(frame, keypoints, &outputMat, color.RGBA{G: 255}, gocv.DrawDefault)
				outputImg, err := outputMat.ToImage()
				if err != nil {
					log.Printf("failed to convert output to image: %s", err)
					os.Exit(1)
				}
				outputImage := canvas.NewImageFromImage(outputImg)
				windowCanvas.SetContent(outputImage)
			}
		}()

		w.Resize(fyne.NewSize(640, 480))
		w.ShowAndRun()
	}
}
