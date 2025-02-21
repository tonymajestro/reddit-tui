package utils

import (
	"os"
	"path/filepath"
)

const (
	appName          = "reddittui"
	defaultConfigDir = ".config"
	defaultStateDir  = ".local/state"
	defaultCacheDir  = ".cache"
	logFileName      = "reddittui.log"
)

func GetConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, defaultConfigDir, appName), nil
}

func GetStateDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, defaultStateDir, appName), nil
}

func GetCacheDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(homeDir, defaultCacheDir, appName), nil
}

func OpenLogFile() (*os.File, error) {
	stateDir, err := GetStateDir()
	if err != nil {
		return nil, err
	}

	err = os.Mkdir(stateDir, 0750)
	if err != nil && !os.IsExist(err) {
		return nil, err
	}

	logPath := filepath.Join(stateDir, logFileName)
	return os.OpenFile(logPath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
}

func FileExists(path string) bool {
	_, err := os.Open(path)
	return os.IsNotExist(err)
}
