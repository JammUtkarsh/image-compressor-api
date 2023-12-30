package imageservice

import (
	"io"
	"net/http"

	"github.com/h2non/bimg"
)

type Image struct {
	Name string
	Body []byte
	Err  error
}

// This function downloadis an image from a given URL and stores it in a buffer.
func (img *Image) ImgURLToBuffer(url string) {
	res, err := http.Get(url)
	if err != nil {
		return
	}
	defer res.Body.Close()

	img.Body, img.Err = io.ReadAll(res.Body)
	if img.Err != nil {
		return
	}
}

// The function takes an image buffer, applies image processing operations such as conversion and
// quality adjustment, and saves the processed image to a file.
func (img *Image) ImgCompress(imgQuality int, imgType bimg.ImageType) {
	img.Body, img.Err = bimg.NewImage(img.Body).Convert(imgType)
	if img.Err != nil {
		return
	}

	img.Body, img.Err = bimg.NewImage(img.Body).Process(bimg.Options{Quality: imgQuality})
	if img.Err != nil {
		return
	}
}
