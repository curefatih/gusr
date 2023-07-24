package cmd

import (
	"fmt"
	"log"

	"github.com/curefatih/gusr/pkg"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List all Git users",
	Run: func(cmd *cobra.Command, args []string) {
		users, err := pkg.LoadUsers()
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
	},
}

func init() {
	rootCmd.AddCommand(listCmd)
}
