package cmd

import (
	"fmt"

	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/db"
	"github.com/jniltinho/go-nagiosql/internal/services/nagimport"
	"github.com/spf13/cobra"
)

var importOverwrite bool
var importConfigID uint8

var importCmd = &cobra.Command{
	Use:   "import <file.cfg>",
	Short: "Import a Nagios .cfg file into the database",
	Args:  cobra.ExactArgs(1),
	RunE:  runImport,
}

func init() {
	importCmd.Flags().BoolVar(&importOverwrite, "overwrite", false, "overwrite existing objects with the same name")
	importCmd.Flags().Uint8Var(&importConfigID, "config-id", 0, "domain/config id to assign (default 0 = common)")
	rootCmd.AddCommand(importCmd)
}

func runImport(cmd *cobra.Command, args []string) error {
	filePath := args[0]

	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	database, err := db.Open(cfg)
	if err != nil {
		return err
	}

	objects, err := nagimport.ParseFile(filePath)
	if err != nil {
		return fmt.Errorf("parsing %s: %w", filePath, err)
	}

	inserted, updated, skipped, errs := 0, 0, 0, 0
	for _, obj := range objects {
		ok, wasNew, importErr := nagimport.ImportObject(database, obj, importConfigID, importOverwrite)
		if importErr != nil {
			name := obj.Fields["host_name"] + obj.Fields["service_description"] + obj.Fields["command_name"]
			fmt.Printf("[ERROR] %s %q: %v\n", obj.Type, name, importErr)
			errs++
			continue
		}
		switch {
		case !ok:
			skipped++
		case wasNew:
			inserted++
		default:
			updated++
		}
	}

	fmt.Printf("import %s: inserted=%d updated=%d skipped=%d errors=%d\n",
		filePath, inserted, updated, skipped, errs)
	if errs > 0 {
		return fmt.Errorf("%d object(s) failed to import", errs)
	}
	return nil
}
