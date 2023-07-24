package cmd

import (
	"fmt"
	"log"
	"os/exec"

	"github.com/curefatih/gusr/pkg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var isGlobal bool

var setCmd = &cobra.Command{
	Use:   "set",
	Short: "Set the Git user globally",
	Run: func(cmd *cobra.Command, args []string) {
		users, err := pkg.LoadUsers()
		if err != nil {
			log.Fatal(err)
		}

		prompt := promptui.Select{
			Label: fmt.Sprintf("Select a Git user to set %s", func() string {
				if isGlobal {
					return "globally"
				}
				return "locally"
			}()),
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
	},
}

func init() {
	rootCmd.AddCommand(setCmd)
	setCmd.Flags().BoolVarP(&isGlobal, "global", "g", false, "--global or -g")
}
