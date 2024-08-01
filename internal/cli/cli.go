package cli

import (
	"fmt"
	"os"

	"github.com/logeshkannan96/dbcli/internal/config"
	"github.com/logeshkannan96/dbcli/internal/database"
	"github.com/logeshkannan96/dbcli/internal/shell"
	"github.com/spf13/cobra"
)

var (
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	configName string
)

var rootCmd = &cobra.Command{
	Use:   "dbcli",
	Short: "Just another db client",
	Long:  `A simple db client with environment switch.`,
}

var connectDBCmd = &cobra.Command{
	Use:   "connect",
	Short: "Connect to the MySQL database and start an interactive shell",
	Run: func(cmd *cobra.Command, args []string) {
		var cfg config.DatabaseConfig
		var err error

		if configName != "" {
			cfg, err = config.LoadConfig(configName)
			if err != nil {
				fmt.Printf("Failed to load config: %v\n", err)
				return
			}
		} else {
			cfg = config.DatabaseConfig{
				Host:     dbHost,
				Port:     dbPort,
				User:     dbUser,
				Password: dbPassword,
				DBName:   dbName,
			}
		}

		err = database.Connect(cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)
		if err != nil {
			fmt.Printf("Failed to connect to database: %v\n", err)
			return
		}
		fmt.Println("Successfully connected to the database")

		shell.StartShell()

		database.Close()
	},
}

var saveConfigCmd = &cobra.Command{
	Use:   "save-config [config-name]",
	Short: "Save the current database configuration",
	Args:  cobra.ExactArgs(1),
	Run: func(cmd *cobra.Command, args []string) {
		cfg := config.DatabaseConfig{
			Name:     args[0],
			Host:     dbHost,
			Port:     dbPort,
			User:     dbUser,
			Password: dbPassword,
			DBName:   dbName,
		}
		if err := config.SaveConfig(cfg); err != nil {
			fmt.Printf("Failed to save config: %v\n", err)
			return
		}
		fmt.Printf("Configuration '%s' saved successfully\n", args[0])
	},
}

var listConfigsCmd = &cobra.Command{
	Use:   "list-configs",
	Short: "List all saved database configurations",
	Run: func(cmd *cobra.Command, args []string) {
		configs, err := config.ListConfigs()
		if err != nil {
			fmt.Printf("Failed to list configs: %v\n", err)
			return
		}
		fmt.Println("Saved configurations:")
		for _, cfg := range configs {
			fmt.Println(cfg)
		}
	},
}

func init() {
	rootCmd.AddCommand(connectDBCmd, saveConfigCmd, listConfigsCmd)

	connectDBCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	connectDBCmd.Flags().IntVar(&dbPort, "port", 3306, "Database port")
	connectDBCmd.Flags().StringVar(&dbUser, "user", "root", "Database user")
	connectDBCmd.Flags().StringVar(&dbPassword, "password", "", "Database password")
	connectDBCmd.Flags().StringVar(&dbName, "dbname", "", "Database name")
	connectDBCmd.Flags().StringVar(&configName, "config", "", "Name of the saved configuration to use")

	saveConfigCmd.Flags().StringVar(&dbHost, "host", "localhost", "Database host")
	saveConfigCmd.Flags().IntVar(&dbPort, "port", 3306, "Database port")
	saveConfigCmd.Flags().StringVar(&dbUser, "user", "root", "Database user")
	saveConfigCmd.Flags().StringVar(&dbPassword, "password", "", "Database password")
	saveConfigCmd.Flags().StringVar(&dbName, "dbname", "", "Database name")
}

func Run() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
