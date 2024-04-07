package main

import (
	"image/color"
	"log"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"gocv.io/x/gocv"
)

func main() {
	r := gin.Default()

	r.LoadHTMLGlob("templates/*")

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	r.POST("/upload", func(c *gin.Context) {
		file, err := c.FormFile("file")
		if err != nil {
			log.Println("Error getting uploaded file:", err)
			c.String(400, "Bad request")
			return
		}

		filename := filepath.Join("upload", file.Filename)
		if err := c.SaveUploadedFile(file, filename); err != nil {
			log.Println("Error saving uploaded file:", err)
			c.String(500, "Internal server error")
			return
		}

		log.Println("Processing uploaded video:", file.Filename)
		if err := processVideo(filename); err != nil {
			log.Println("Error processing video:", err)
			c.String(500, "Error processing video")
			return
		}
		log.Println("Video processing completed")

		newPath := filepath.Join("verified", file.Filename)
		if err := os.Rename(filename, newPath); err != nil {
			log.Println("Error moving file:", err)
			c.String(500, "Error moving file")
			return
		}

		c.String(200, "Video uploaded and verified successfully")
	})

	r.Run(":8080")
}

func processVideo(videoPath string) error {
	webcam, err := gocv.VideoCaptureFile(videoPath)
	if err != nil {
		return err
	}
	defer webcam.Close()

	window := gocv.NewWindow("Face Detection")
	defer window.Close()

	classifier := gocv.NewCascadeClassifier()
	classifier.Load("haarcascade_frontalface_default.xml")

	for {
		frame := gocv.NewMat()
		if webcam.Read(&frame) {
			if frame.Empty() {
				continue
			}

			faces := classifier.DetectMultiScale(frame)
			for _, face := range faces {
				gocv.Rectangle(&frame, face, color.RGBA{0, 255, 0, 0}, 2)
			}

			window.IMShow(frame)
			if window.WaitKey(1) >= 0 {
				break
			}
		} else {
			break
		}
	}

	return nil
}
