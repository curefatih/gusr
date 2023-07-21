package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

type GitUser struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	GPGKey string `json:"gpgKey"`
}

var configFile = "git-users.json"

func main() {
	if err := createConfigFileIfNotExist(); err != nil {
		log.Fatal(err)
	}

	rootCmd := &cobra.Command{
		Use:   "gusr",
		Short: "A tool for managing Git users",
	}

	rootCmd.AddCommand(&cobra.Command{
		Use:   "list",
		Short: "List all Git users",
		Run:   list,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "set",
		Short: "Set the Git user globally",
		Run:   set,
	})

	rootCmd.AddCommand(&cobra.Command{
		Use:   "add",
		Short: "Add a Git user",
		Run:   add,
	})

	if err := rootCmd.Execute(); err != nil {
		log.Fatal(err)
	}
}

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

func createConfigFileIfNotExist() error {
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
	err = saveUsers(users)
	if err != nil {
		return err
	}

	return nil
}

func set(cmd *cobra.Command, args []string) {
	var isGlobal bool
	cmd.Flags().BoolVarP(&isGlobal, "global", "g", false, "set the Git user globally")

	users, err := loadUsers()
	if err != nil {
		log.Fatal(err)
	}

	prompt := promptui.Select{
		Label: "Select a Git user",
		Items: users,
		Templates: &promptui.SelectTemplates{
			Label:    "{{ .Name }}",
			Selected: "{{ .Name }}",
		},
	}

	index, _, err := prompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	user := users[index]

	if isGlobal {
		cmd := exec.Command("git", "config", "--global", "user.name", user.Name)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		cmd = exec.Command("git", "config", "--global", "user.email", user.Email)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if user.GPGKey != "" {
			cmd = exec.Command("git", "config", "--global", "user.signingkey", user.GPGKey)
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("Git user %s set globally\n", user.Name)
	} else {
		cmd := exec.Command("git", "config", "user.name", user.Name)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		cmd = exec.Command("git", "config", "user.email", user.Email)
		if err := cmd.Run(); err != nil {
			log.Fatal(err)
		}

		if user.GPGKey != "" {
			cmd = exec.Command("git", "config", "user.signingkey", user.GPGKey)
			if err := cmd.Run(); err != nil {
				log.Fatal(err)
			}
		}

		fmt.Printf("Git user %s set locally\n", user.Name)
	}
}

func add(cmd *cobra.Command, args []string) {
	namePrompt := promptui.Prompt{
		Label: "Git user name",
	}

	name, err := namePrompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	emailPrompt := promptui.Prompt{
		Label: "Git user email",
	}

	email, err := emailPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	gpgKeyPrompt := promptui.Prompt{
		Label: "GPG Key (optional)",
	}

	gpgKey, err := gpgKeyPrompt.Run()
	if err != nil {
		log.Fatal(err)
	}

	newUser := GitUser{
		Name:   name,
		Email:  email,
		GPGKey: gpgKey,
	}

	users, err := loadUsers()
	if err != nil {
		log.Fatal(err)
	}

	users = append(users, newUser)

	if err := saveUsers(users); err != nil {
		log.Fatal(err)
	}

	fmt.Printf("Git user %s <%s> added\n", name, email)
}

func loadUsers() ([]GitUser, error) {
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

func saveUsers(users []GitUser) error {
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

func list(cmd *cobra.Command, args []string) {
	users, err := loadUsers()
	if err != nil {
		log.Fatal(err)
	}

	if len(users) == 0 {
		fmt.Println("There is no user saved.")
		return
	}

	fmt.Println("Git users:")
	for _, user := range users {
		fmt.Printf("- %s <%s>\n", user.Name, user.Email)
		if user.GPGKey != "" {
			fmt.Printf("  GPG key: %s\n", user.GPGKey)
		}
	}
}

func runGitCommand(args []string) error {
	cmd := exec.Command("git", args...)
	cmd.Stderr = os.Stderr
	return cmd.Run()
}
