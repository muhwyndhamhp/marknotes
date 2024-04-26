package imageprocessing

import (
	"bytes"
	"image/jpeg"
	"log"
	"mime/multipart"

	"github.com/disintegration/imaging"
	"github.com/kolesa-team/go-webp/encoder"
	"github.com/kolesa-team/go-webp/webp"
	"github.com/muhwyndhamhp/marknotes/utils/tern"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (cl *Client) ResizeImage(file *multipart.FileHeader, size int) (*bytes.Reader, int, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, -1, err
	}

	defer reader.Close()

	src, err := imaging.Decode(reader)
	if err != nil {
		return nil, -1, err
	}

	dstImage800 := imaging.Resize(src, tern.Int(size, 1200), 0, imaging.Lanczos)

	var b []byte
	w := bytes.NewBuffer(b)
	err = imaging.Encode(w, dstImage800, imaging.JPEG)
	if err != nil {
		return nil, -1, err
	}

	r := bytes.NewReader(w.Bytes())

	options, err := encoder.NewLossyEncoderOptions(encoder.PresetDefault, 75)
	if err != nil {
		log.Fatalln(err)
	}

	img, err := jpeg.Decode(r)
	if err != nil {
		return nil, -1, err
	}

	var o []byte
	wo := bytes.NewBuffer(o)
	err = webp.Encode(wo, img, options)
	if err != nil {
		return nil, -1, err
	}

	r = bytes.NewReader(wo.Bytes())

	return r, int(r.Size()), nil
}
