package cmd

import (
	"fmt"
	"log"

	"github.com/jniltinho/go-nagiosql/internal/config"
	"github.com/jniltinho/go-nagiosql/internal/db"
	"github.com/jniltinho/go-nagiosql/internal/models"
	"github.com/jniltinho/go-nagiosql/internal/services/nagconfig"
	"github.com/jniltinho/go-nagiosql/internal/services/nagios"
	"github.com/spf13/cobra"
	"gorm.io/gorm"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Manage Nagios configuration files",
}

var configWriteCmd = &cobra.Command{
	Use:   "write [host|service|all]",
	Short: "Write Nagios config files from the database",
	Args:  cobra.MaximumNArgs(1),
	RunE:  runConfigWrite,
}

var configVerifyCmd = &cobra.Command{
	Use:   "verify",
	Short: "Run nagios -v to verify the configuration",
	RunE:  runConfigVerify,
}

var configRestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Verify config then trigger a graceful Nagios reload via reload.trigger",
	RunE:  runConfigRestart,
}

func init() {
	configCmd.AddCommand(configWriteCmd, configVerifyCmd, configRestartCmd)
	rootCmd.AddCommand(configCmd)
}

func loadConfigTarget(database *gorm.DB) (models.Configtarget, error) {
	var ct models.Configtarget
	err := database.Where("domain_id = ?", 0).First(&ct).Error
	return ct, err
}

func runConfigWrite(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	database, err := db.Open(cfg)
	if err != nil {
		return err
	}
	ct, err := loadConfigTarget(database)
	if err != nil {
		return fmt.Errorf("loading configtarget: %w", err)
	}

	gen := nagconfig.New(database, ct.HostPath, ct.ServicePath, ct.BackupPath)

	objectType := "all"
	if len(args) > 0 {
		objectType = args[0]
	}

	switch objectType {
	case "host", "hosts":
		n, err := gen.WriteAllHosts()
		if err != nil {
			return fmt.Errorf("write hosts: %w", err)
		}
		fmt.Printf("wrote %d host config file(s)\n", n)
	case "service", "services":
		n, err := gen.WriteAllServices()
		if err != nil {
			return fmt.Errorf("write services: %w", err)
		}
		fmt.Printf("wrote %d service config file(s)\n", n)
	default:
		n, err := gen.WriteAll()
		if err != nil {
			return fmt.Errorf("write all: %w", err)
		}
		fmt.Printf("wrote %d config file(s)\n", n)
	}
	return nil
}

func runConfigVerify(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	result := nagios.Verify(cfg.Nagios.Binary, cfg.Nagios.ConfigFile)
	fmt.Println(result.Output)
	if !result.Valid {
		return fmt.Errorf("config verification failed")
	}
	fmt.Println("config verify: OK")
	return nil
}

func runConfigRestart(cmd *cobra.Command, args []string) error {
	cfg, err := config.Load(cfgFile)
	if err != nil {
		return fmt.Errorf("loading config: %w", err)
	}
	result, err := nagios.Restart(cfg.Nagios.Binary, cfg.Nagios.ConfigFile, cfg.Nagios.ReloadTrigger)
	fmt.Println(result.Output)
	if err != nil {
		return fmt.Errorf("restart: %w", err)
	}
	if !result.Valid {
		return fmt.Errorf("config verification failed; reload not triggered")
	}
	log.Printf("reload trigger written to %s", cfg.Nagios.ReloadTrigger)
	fmt.Println("config restart: OK")
	return nil
}
