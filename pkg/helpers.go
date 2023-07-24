package pkg

import (
	"encoding/json"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

var configFile = "git-users.json"

func configDir() (string, error) {
	switch runtime.GOOS {
	case "windows":
		return filepath.Join(os.Getenv("APPDATA"), "GitUser"), nil
	case "darwin", "linux":
		dir, err := os.UserHomeDir()
		if err != nil {
			return "", err
		}
		return filepath.Join(dir, ".git-user"), nil
	default:
		return "", errors.New("unsupported operating system")
	}
}

func getConfigFilePath() (string, error) {
	configPath, err := configDir()
	if err != nil {
		return "", err
	}

	return filepath.Join(configPath, configFile), nil
}

func CreateConfigFileIfNotExist() error {
	configPath, err := configDir()
	if err != nil {
		return err
	}

	// Create the directory if it doesn't exist.
	err = os.MkdirAll(configPath, os.ModePerm)
	if err != nil {
		return err
	}

	configFilePath := filepath.Join(configPath, configFile)

	// If the config file already exists, we don't need to do anything.
	if _, err := os.Stat(configFilePath); err == nil {
		return nil
	}

	// Create the config file.
	users := []GitUser{}
	err = SaveUsers(users)
	if err != nil {
		return err
	}

	return nil
}

func LoadUsers() ([]GitUser, error) {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return nil, err
	}
	file, err := os.Open(configFilePath)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var users []GitUser
	err = json.NewDecoder(file).Decode(&users)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func SaveUsers(users []GitUser) error {
	configFilePath, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(configFilePath)
	if err != nil {
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	if err := encoder.Encode(users); err != nil {
		return err
	}

	return nil
}

func runGitCommand(args []string) error {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
