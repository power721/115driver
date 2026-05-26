package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/power721/115driver/cli/internal/auth"
	"github.com/power721/115driver/cli/internal/output"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Show current configuration",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := resolveConfigPath()

		v, err := readConfig(path)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		// Resolve profile: --profile flag > env > config default_profile > "main"
		currentProfile := profile
		if currentProfile == "" {
			if envProfile := os.Getenv(auth.EnvProfile); envProfile != "" {
				currentProfile = envProfile
			} else {
				currentProfile = v.GetString("default_profile")
				if currentProfile == "" {
					currentProfile = auth.DefaultProfile
				}
			}
		}

		if jsonOutput {
			printer.PrintSuccess(map[string]interface{}{
				"config_path": path,
				"profile":     currentProfile,
			})
		} else {
			fmt.Printf("Config file: %s\n", path)
			fmt.Printf("Profile:     %s\n", currentProfile)
		}
		return nil
	},
}

var configListCmd = &cobra.Command{
	Use:   "list",
	Short: "List all configured profiles",
	Args:  cobra.NoArgs,
	RunE: func(cmd *cobra.Command, args []string) error {
		path := resolveConfigPath()
		v, err := readConfig(path)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		defaultProfile := v.GetString("default_profile")
		if defaultProfile == "" {
			defaultProfile = auth.DefaultProfile
		}

		profiles := v.GetStringMap("profiles")
		if len(profiles) == 0 {
			if jsonOutput {
				printer.PrintSuccess(map[string]interface{}{
					"profiles": []map[string]interface{}{},
				})
			} else {
				fmt.Println("No profiles configured. Run '115driver login' to create one.")
			}
			return nil
		}

		// Sort profile names for consistent output
		names := make([]string, 0, len(profiles))
		for name := range profiles {
			names = append(names, name)
		}
		sort.Strings(names)

		type profileEntry struct {
			Name    string `json:"name"`
			Default bool   `json:"default"`
			HasAuth bool   `json:"has_auth"`
		}

		entries := make([]profileEntry, 0, len(names))
		for _, name := range names {
			cookie := v.GetString("profiles." + name + ".cookie")
			entries = append(entries, profileEntry{
				Name:    name,
				Default: name == defaultProfile,
				HasAuth: cookie != "",
			})
		}

		if jsonOutput {
			printer.PrintSuccess(map[string]interface{}{
				"default_profile": defaultProfile,
				"profiles":        entries,
			})
		} else {
			fmt.Printf("Profiles (default: %s):\n\n", defaultProfile)
			for _, e := range entries {
				marker := "  "
				if e.Default {
					marker = "* "
				}
				authStatus := "no auth"
				if e.HasAuth {
					authStatus = "authenticated"
				}
				fmt.Printf("%s%s (%s)\n", marker, e.Name, authStatus)
			}
		}
		return nil
	},
}

var configUseCmd = &cobra.Command{
	Use:   "use <profile>",
	Short: "Set the default profile",
	Args:  cobra.ExactArgs(1),
	RunE: func(cmd *cobra.Command, args []string) error {
		targetProfile := args[0]
		path := resolveConfigPath()

		v, err := readConfig(path)
		if err != nil {
			return &exitError{code: output.ExitError, msg: err.Error()}
		}

		// Verify profile exists
		cookie := v.GetString("profiles." + targetProfile + ".cookie")
		if cookie == "" {
			return &exitError{code: output.ExitNotFound, msg: fmt.Sprintf("Profile '%s' not found. Run '115driver login --profile %s' to create it.", targetProfile, targetProfile)}
		}

		v.Set("default_profile", targetProfile)

		if err := writeConfig(v, path); err != nil {
			return &exitError{code: output.ExitError, msg: fmt.Sprintf("Failed to save config: %v", err)}
		}

		printer.PrintSuccess(map[string]interface{}{
			"default_profile": targetProfile,
		})
		if !jsonOutput {
			fmt.Printf("Default profile set to: %s\n", targetProfile)
		}
		return nil
	},
}

func resolveConfigPath() string {
	if configPath != "" {
		return configPath
	}
	if envPath := os.Getenv(auth.EnvConfig); envPath != "" {
		return envPath
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, auth.DefaultConfigDir, auth.DefaultConfigFile)
}

func readConfig(path string) (*viper.Viper, error) {
	v := viper.New()
	v.SetConfigFile(path)
	if err := v.ReadInConfig(); err != nil {
		if os.IsNotExist(err) {
			return v, nil // empty config is fine
		}
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return v, nil
		}
		return nil, fmt.Errorf("failed to read config: %w", err)
	}
	return v, nil
}

func writeConfig(v *viper.Viper, path string) error {
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return fmt.Errorf("create config dir: %w", err)
	}
	if err := v.WriteConfig(); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func init() {
	configCmd.AddCommand(configListCmd)
	configCmd.AddCommand(configUseCmd)
	rootCmd.AddCommand(configCmd)
}
