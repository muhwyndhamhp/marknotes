package storage

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/muhwyndhamhp/marknotes/config"
)

func ServeFile(filename string) (string, error) {
	filePath := fmt.Sprintf("%s/%s", config.Get(config.STORE_VOL_PATH), filename)

	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return "", err
	}

	return filePath, nil
}

func WriteFile(f *multipart.FileHeader, prefix string) (string, error) {
	file, err := f.Open()
	if err != nil {
		return "", err
	}

	fname := strings.ReplaceAll(f.Filename, " ", "_")
	name := fmt.Sprintf("%s-%s", prefix, AppendTimestamp(fname))

	dst, err := os.Create(fmt.Sprintf("%s/%s", config.Get(config.STORE_VOL_PATH), name))
	if err != nil {
		return "", err
	}

	_, err = io.Copy(dst, file)
	if err != nil {
		return "", err
	}

	file.Close()
	dst.Close()

	baseURL := strings.Split(config.Get(config.OAUTH_URL), "/callback")[0]

	return fmt.Sprintf("%s/posts/media/%s", baseURL, name), nil
}

func AppendTimestamp(fileName string) string {
	extension := filepath.Ext(fileName)
	name := fileName[0 : len(fileName)-len(extension)]
	fileName = fmt.Sprintf("%s-%s%s", name, time.Now().Format("20060102150405"), extension)
	return fileName
}

func IsValidFileType(fileHeader *multipart.FileHeader) (string, bool) {
	allowedImageTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
		// Add more image types as needed
	}

	allowedVideoTypes := map[string]bool{
		"video/mp4":  true,
		"video/avi":  true,
		"video/mpeg": true,
		"video/mov":  true,
		// Add more video types as needed
	}

	file, err := fileHeader.Open()
	if err != nil {
		return "", false
	}
	defer file.Close()

	buffer := make([]byte, 512) // Read the first 512 bytes to detect file type
	_, err = file.Read(buffer)
	if err != nil {
		return "", false
	}

	contentType := http.DetectContentType(buffer)

	// Check if the content type is allowed for either images or videos
	return contentType, allowedImageTypes[contentType] || allowedVideoTypes[contentType]
}
