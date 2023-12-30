package services

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/h2non/bimg"
)

type Image struct {
	Name string
	Body []byte
}

// This function downloadis an image from a given URL and stores it in a buffer.
func (img *Image) ImgURLToBuffer(url string) error{
	res, err := http.Get(url)
	if err != nil {
		return err
	}
	defer res.Body.Close()
	if res.StatusCode != 200 {
		return errors.New("received non 200 response code")
	}
	if img.Body, err = io.ReadAll(res.Body); err != nil {
		return err
	}
	img.Name = filepath.Base(url)
	return nil
}

// The function takes a byte array and compresses it to the given quality and image type.
func (img *Image) ImgCompress(imgQuality int, imgType bimg.ImageType) (err error) {
	if img.Body, err = bimg.NewImage(img.Body).Convert(imgType); err != nil {
		return err
	}

	if img.Body, err = bimg.NewImage(img.Body).Process(bimg.Options{Quality: imgQuality}); err != nil {
		return err
	}
	return nil
}

// This function take a byte array and saves it to a file in the given directory.
func (img *Image) SaveFile(directory string) (err error) {
	fPath := filepath.Join(directory, img.Name)
	if _, err = os.Stat(directory); os.IsNotExist(err) {
		if err = os.MkdirAll(directory, 0755); err != nil {
			return err
		}
	}
	
	obj, err := os.Create(fPath)
	if err != nil {
		return err
	}
	defer obj.Close()
	
	if _, err = obj.Write(img.Body); err != nil {
		return err
	}
	return nil
}
