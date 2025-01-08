package utils

import (
	"os"
	"path/filepath"
	"repeat-what-shit/internal/consts"
)

func GetAppDirPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, consts.AppDataDirName), nil
}

func CreateAppDirIfNotExists() error {
	appDir, err := GetAppDirPath()
	if err != nil {
		return err
	}

	if _, err := os.Stat(appDir); os.IsNotExist(err) {
		return os.MkdirAll(appDir, 0755)
	}

	return nil
}
