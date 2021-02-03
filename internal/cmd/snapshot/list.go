package snapshot

import (
	"context"
	"fmt"

	"github.com/planetscale/cli/internal/cmdutil"
	"github.com/planetscale/cli/internal/config"
	"github.com/planetscale/cli/internal/printer"
	"github.com/planetscale/planetscale-go/planetscale"
	"github.com/spf13/cobra"
)

// ListCmd makes a command for listing all snapshots for a database branch.
func ListCmd(cfg *config.Config) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "list <database> <branch>",
		Aliases: []string{"ls"},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := context.Background()
			client, err := cfg.NewClientFromConfig()
			if err != nil {
				return err
			}

			if len(args) != 2 {
				return cmd.Usage()
			}

			database, branch := args[0], args[1]

			end := cmdutil.PrintProgress(fmt.Sprintf("Fetching schema snapshots for %s in %s...", cmdutil.BoldBlue(branch), cmdutil.BoldBlue(database)))
			defer end()

			snapshots, err := client.SchemaSnapshots.List(ctx, &planetscale.ListSchemaSnapshotsRequest{
				Organization: cfg.Organization,
				Database:     database,
				Branch:       branch,
			})
			if err != nil {
				return err
			}
			end()

			isJSON, err := cmd.Flags().GetBool("json")
			if err != nil {
				return err
			}

			if len(snapshots) == 0 && !isJSON {
				fmt.Printf("No schema snapshots exist for %s in %s.\n", cmdutil.BoldBlue(branch), cmdutil.BoldBlue(database))
				return nil
			}

			err = printer.PrintOutput(isJSON, printer.NewSchemaSnapshotSlicePrinter(snapshots))
			if err != nil {
				return err
			}

			return nil
		},
	}

	return cmd
}