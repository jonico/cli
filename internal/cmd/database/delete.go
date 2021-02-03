package database

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/planetscale/cli/internal/cmdutil"
	"github.com/planetscale/cli/internal/config"

	"github.com/planetscale/planetscale-go/planetscale"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/terminal"
	"github.com/spf13/cobra"
)

// DeleteCmd is the Cobra command for deleting a database for an authenticated
// user.
func DeleteCmd(cfg *config.Config) *cobra.Command {
	var force bool

	cmd := &cobra.Command{
		Use:     "delete <database_name>",
		Short:   "Delete a database instance",
		Aliases: []string{"rm"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := cfg.NewClientFromConfig()
			if err != nil {
				return err
			}

			if len(args) == 0 {
				return errors.New("<database_name> is missing")
			}

			name := args[0]

			if !force {
				userInput := ""
				confirmationMessage := fmt.Sprintf("%s %s %s", cmdutil.Bold("Please type"), cmdutil.BoldBlue(name), cmdutil.Bold("to confirm:"))

				prompt := &survey.Input{
					Message: confirmationMessage,
				}

				err := survey.AskOne(prompt, &userInput)
				if err != nil {
					if err == terminal.InterruptErr {
						os.Exit(0)
					} else {
						return err
					}
				}

				if userInput != name {
					return errors.New("Incorrect database name entered, skipping database deletion...")
				}
			}

			end := cmdutil.PrintProgress(fmt.Sprintf("Deleting database %s...", cmdutil.BoldBlue(name)))
			defer end()

			err = client.Databases.Delete(ctx, &planetscale.DeleteDatabaseRequest{
				Organization: cfg.Organization,
				Database:     name,
			})
			if err != nil {
				return err
			}

			end()
			fmt.Printf("Database %s was successfully deleted!\n", cmdutil.BoldBlue(name))

			return nil
		},
	}

	cmd.Flags().BoolVarP(&force, "force", "f", false, "Delete a databse without confirmation")
	return cmd
}