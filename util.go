package sbctl

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"path/filepath"
)

func ChecksumFile(file string) string {
	hasher := sha256.New()
	s, err := os.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	hasher.Write(s)

	return hex.EncodeToString(hasher.Sum(nil))
}

func ReadOrCreateFile(filePath string) ([]byte, error) {
	// Try to access or create the file itself
	f, err := os.ReadFile(filePath)
	if err != nil {
		// Errors will mainly happen due to permissions or non-existing file
		if os.IsNotExist(err) {
			// First, guarantee the directory's existence
			// os.MkdirAll simply returns nil if the directory already exists
			fileDir := filepath.Dir(filePath)
			if err = os.MkdirAll(fileDir, os.ModePerm); err != nil {
				return nil, err
			}

			file, err := os.Create(filePath)
			if err != nil {
				return nil, err
			}
			file.Close()

			// Create zero-length f, which is equivalent to what would be read from empty file
			f = make([]byte, 0)
		} else {
			if os.IsPermission(err) {
				return nil, err
			}
			return nil, err
		}
	}

	return f, nil
}

func IsImmutable(file string) (bool, error) {
	f, err := os.Open(file)
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	attr, err := GetAttr(f)
	if err != nil {
		return false, err
	}
	if (attr & FS_IMMUTABLE_FL) != 0 {
		return false, nil
	}
	return true, nil
}
