// DecodeImage decodifica uma imagem genérica a partir de bytes
package imageutil

import (
	"bytes"
	"image"

	"github.com/disintegration/imaging"
)

func DecodeImage(data []byte) (image.Image, error) {
	return imaging.Decode(bytes.NewReader(data))
}

type ResizeOptions struct {
	Width   int
	Height  int
	Format  imaging.Format // imaging.JPEG, imaging.PNG, etc.
	Quality int            // 1-100 (JPEG)
}

// ResizeImage redimensiona uma imagem mantendo alta qualidade

// ResizeImage redimensiona uma imagem mantendo alta qualidade
func ResizeImage(input []byte, opts ResizeOptions) ([]byte, error) {
	img, err := imaging.Decode(bytes.NewReader(input))
	if err != nil {
		return nil, err
	}
	resized := imaging.Resize(img, opts.Width, opts.Height, imaging.Lanczos)
	var buf bytes.Buffer
	switch opts.Format {
	case imaging.JPEG:
		err = imaging.Encode(&buf, resized, imaging.JPEG, imaging.JPEGQuality(opts.Quality))
	case imaging.PNG:
		err = imaging.Encode(&buf, resized, imaging.PNG)
	case imaging.GIF:
		err = imaging.Encode(&buf, resized, imaging.GIF)
	default:
		err = imaging.Encode(&buf, resized, imaging.PNG)
	}
	if err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
