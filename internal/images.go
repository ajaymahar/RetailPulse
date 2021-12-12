package internal

import (
	"fmt"
	"image/jpeg"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"
)

// GetFileName from an image
func GetFileName(imageURL string) (string, error) {

	iURL, err := url.Parse(imageURL)
	if err != nil {
		return "", fmt.Errorf("invalid url format: %w", err)
	}
	path := iURL.Path
	seg := strings.Split(path, "/")
	name := seg[len(seg)-1]
	return name, nil

}

//DownloadImage will download and save it
func DownloadImage(imageURL, iName string) error {
	res, err := http.Get(imageURL)
	if err != nil {
		return fmt.Errorf("downloading failed: %w ", err)
	}
	// cleanup
	defer res.Body.Close()

	// saving image locally
	if err := saveImage(iName, res.Body); err != nil {
		return fmt.Errorf("downloading failed: %w", err)
	}
	return nil
}

//Save image locally
func saveImage(fileName string, stream io.Reader) error {
	f, err := os.Create("images/" + fileName)
	if err != nil {
		return fmt.Errorf("can't save image: %w", err)
	}

	// clean up
	defer f.Close()

	_, err = io.Copy(f, stream)
	if err != nil {
		return fmt.Errorf("can't copy image bytes: %w", err)
	}
	return nil
}

// GetDimmensions will return the perimeter 2* [Height+Width] of each image
func GetDimmensions(imgName string) (int, error) {

	imgf, err := os.Open("images/" + imgName)
	if err != nil {
		return 0, fmt.Errorf("can't open image: %w", err)
	}
	defer imgf.Close()

	img, err := jpeg.Decode(imgf)
	if err != nil {
		return 0, fmt.Errorf("can't decode <%v>: %w", imgName, err)
	}

	perimeter := 2 * (img.Bounds().Max.X + img.Bounds().Max.Y)
	// fmt.Printf("Width: %v Height: %v\n", img.Bounds().Max.X, img.Bounds().Max.Y)

	// random sleep for 0-4 seconds
	// random sleep time of 0.1 to 0.4 secs (this is to imitate GPU processing)
	rand.Seed(time.Now().UnixNano())
	r := rand.Intn(4)
	time.Sleep(time.Duration(r) * time.Second)
	return perimeter, nil
}
