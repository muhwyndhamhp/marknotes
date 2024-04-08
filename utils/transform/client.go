package transform

import (
	"fmt"
	"io"
	"mime/multipart"
	"os"

	"github.com/h2non/bimg"
	"github.com/muhwyndhamhp/marknotes/utils/errs"
)

type Client struct{}

func (c *Client) ResizeWeb(f *multipart.FileHeader) (*os.File, error) {
	reader, err := f.Open()
	if err != nil {
		return nil, errs.Wrap(err)
	}

	// read all bytes
	buffer, err := io.ReadAll(reader)
	if err != nil {
		return nil, errs.Wrap(err)
	}

	newImage, err := bimg.NewImage(buffer).Resize(800, 600)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	fmt.Println(string(newImage))

	return nil, nil
}
