package config

import (
	"fmt"
	"os"

	"github.com/spf13/viper"
)

const configDir = "/.config/kttn/"
const configFile = "config"
const configFileType = "yaml"

func InitConfig() {
  // Find home directory.
  home, err := os.UserHomeDir()
  if err != nil {
    fmt.Fprintln(os.Stderr, "Could not get user's home directory.") 
  }

  viper.SetDefault("word_set", "en_500")
  viper.SetDefault("number_of_words", 50)
  viper.SetDefault("max_wmp", 0)
  viper.SetDefault("max_cmp", 0)
  
  // Search config in home directory with name ". toolbox" (without extension).
  viper.AddConfigPath (home + configDir)
  viper.SetConfigType (configFileType)
  viper.SetConfigName (configFile)

  if err := viper.ReadInConfig(); err == nil {
    fmt.Fprintln(os.Stdout, "Using config file:", viper.ConfigFileUsed())
  } else if _, ok := err.(viper.ConfigFileNotFoundError); ok {
    fmt.Fprintln(os.Stderr, "Config file not found, will create one with defaults.") 
    createConfigDirIfNeeded(home)
    confError := viper.WriteConfigAs(home + configDir + configFile + "." + configFileType)
    if confError != nil {
      fmt.Printf("Error writing config %v", confError)
    }
	} else {
    fmt.Fprintln(os.Stderr, "Error reading existing config file. Will just use defaults.") 
  }
}

func createConfigDirIfNeeded(homeDir string) {
    os.MkdirAll(homeDir + configDir, 0755)
}
