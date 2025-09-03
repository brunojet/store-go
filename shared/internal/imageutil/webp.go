package imageutil

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
)

func ResizeImageToPNG(img image.Image, width, height int) ([]byte, error) {
	resized := imaging.Resize(img, width, height, imaging.Lanczos)
	var buf bytes.Buffer
	err := imaging.Encode(&buf, resized, imaging.PNG)
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
