package imageprocessing

import (
	"bytes"
	"mime/multipart"

	"github.com/disintegration/imaging"
)

type Client struct{}

func NewClient() *Client {
	return &Client{}
}

func (cl *Client) ResizeImage(file *multipart.FileHeader) (*bytes.Reader, int, error) {
	reader, err := file.Open()
	if err != nil {
		return nil, -1, err
	}

	defer reader.Close()

	src, err := imaging.Decode(reader)
	if err != nil {
		return nil, -1, err
	}

	dstImage800 := imaging.Resize(src, 1200, 0, imaging.Lanczos)

	var b []byte
	w := bytes.NewBuffer(b)
	err = imaging.Encode(w, dstImage800, imaging.JPEG)
	if err != nil {
		return nil, -1, err
	}

	r := bytes.NewReader(w.Bytes())

	return r, int(r.Size()), nil
}
