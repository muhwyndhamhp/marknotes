package fileman

import (
	"fmt"
	"os"
)

func CheckDir(path string) error {
	if _, err := os.Stat(path); os.IsNotExist(err) {
		err := os.Mkdir(path, 0o755)
		if err != nil {
			fmt.Println(err)
		}
	}
	return nil
}

func DeleteFile(filepath string) error {
	if _, err := os.Stat(filepath); os.IsNotExist(err) {
		return nil
	}

	err := os.Remove(filepath)
	if err != nil {
		fmt.Println(err)
		return err
	}
	return nil
}

func DeletAllFiles(dir string) error {
	files, err := os.ReadDir(dir)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return nil
	}

	for _, file := range files {
		if file.IsDir() {
			continue
		}

		err := os.Remove(fmt.Sprintf("%s/%s", dir, file.Name()))
		if err != nil {
			fmt.Println("Error deleting file:", err)
			return nil
		}
	}

	return nil
}
