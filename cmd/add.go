package cmd

import (
	"fmt"
	"log"

	"github.com/curefatih/gusr/pkg"
	"github.com/manifoldco/promptui"
	"github.com/spf13/cobra"
)

var addCmd = &cobra.Command{
	Use:   "add",
	Short: "Add a Git user",
	Run: func(cmd *cobra.Command, args []string) {
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

		newUser := pkg.GitUser{
			Name:   name,
			Email:  email,
			GPGKey: gpgKey,
		}

		users, err := pkg.LoadUsers()
		if err != nil {
			log.Fatal(err)
		}

		users = append(users, newUser)

		if err := pkg.SaveUsers(users); err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Git user %s <%s> added\n", name, email)
	},
}

func init() {
	rootCmd.AddCommand(addCmd)
}
