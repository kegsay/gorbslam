package main

import (
	"flag"
	"fmt"
	"image/color"
	"log"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"time"

	"net/http"
	_ "net/http/pprof"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"github.com/kegsay/gorbslam/internal/orb"
	"gocv.io/x/gocv"
)

// Checks that the ORB features are being extracted from images correctly.

var (
	flagImg                = flag.String("image", "", "image to plot ORB features onto")
	flagVideoCaptureDevice = flag.Int("cam", 0, "video capture device to read from")
)

func parseFrame(w fyne.Window, vc *gocv.VideoCapture, windowCanvas fyne.Canvas, imgCanvas *canvas.Image) {
	if !vc.IsOpened() {
		log.Println("waiting for video capture device...")
		time.Sleep(100 * time.Millisecond)
		return
	}
	frame := gocv.NewMat()
	defer frame.Close()
	if !vc.Read(&frame) {
		log.Println("failed to read frame")
		return
	}
	w.Resize(fyne.NewSize(float32(frame.Cols()), float32(frame.Rows())))

	//frame := gocv.IMRead(*flagImg, gocv.IMReadAnyColor)
	if frame.Empty() {
		log.Printf("cannot load image at %s", *flagImg)
		os.Exit(1)
	}

	keypoints, _ := orb.Features(frame, false)
	log.Printf("Detected %d keypoints", len(keypoints))
	if len(keypoints) == 0 {
		return
	}
	outputMat := gocv.NewMatWithSize(frame.Rows(), frame.Cols(), frame.Type())
	defer outputMat.Close()
	gocv.DrawKeyPoints(frame, keypoints, &outputMat, color.RGBA{G: 255}, gocv.DrawDefault)
	outputImg, err := outputMat.ToImage()
	if err != nil {
		log.Printf("failed to convert output to image: %s", err)
		os.Exit(1)
	}
	imgCanvas.Image = outputImg
	imgCanvas.Refresh()
}

func main() {
	flag.Parse()
	// pprof
	go func() {
		if err := http.ListenAndServe(":6060", nil); err != nil {
			panic(err)
		}
	}()
	// run it on an image
	if flagImg != nil && *flagImg != "" {
		img := gocv.IMRead(*flagImg, gocv.IMReadAnyColor)
		if img.Empty() {
			log.Printf("cannot load image at %s", *flagImg)
			os.Exit(1)
		}
		keypoints, _ := orb.Features(img, false)
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
		frameCount := 0

		go func() {
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
			var canvasImg canvas.Image
			canvasImg.ScaleMode = canvas.ImageScaleFastest
			windowCanvas.SetContent(&canvasImg)
			for {
				select {
				case <-c:
					return
				default:
				}
				parseFrame(w, vc, windowCanvas, &canvasImg)
				frameCount++
				if frameCount > 60 {
					frameCount = 0
					runtime.GC()
					var stats runtime.MemStats
					runtime.ReadMemStats(&stats)
					fmt.Printf("Memory: %v MB\n", (float64(stats.Alloc)/1024.0)/1024.0)
					debug.FreeOSMemory()
				}
			}
		}()

		w.Resize(fyne.NewSize(640, 480))
		w.ShowAndRun()
	}
}
