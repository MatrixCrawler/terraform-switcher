package param_parsing

import (
	"fmt"
	"github.com/spf13/viper"
	"github.com/warrensbox/terraform-switcher/lib"
	"log"
)

const tfSwitchTOMLFileName = ".tfswitch.toml"

// getParamsTOML parses everything in the toml file, return required version and bin path
func getParamsTOML(params Params) Params {
	tomlPath := params.ChDirPath + "/" + tfSwitchTOMLFileName
	if tomlFileExists(params) {
		fmt.Printf("Reading configuration from %s\n", tomlPath)
		configfileName := lib.GetFileName(tfSwitchTOMLFileName)
		viper.SetConfigType("toml")
		viper.SetConfigName(configfileName)
		viper.AddConfigPath(params.ChDirPath)

		errs := viper.ReadInConfig() // Find and read the config file
		if errs != nil {
			log.Fatalf("Unable to read %s provided\n", tomlPath)
		}

		params.Version = viper.GetString("version") //attempt to get the version if it's provided in the toml
		params.CustomBinaryPath = viper.GetString("bin")
	} else {
		fmt.Println("No configuration file at " + tomlPath)
	}
	return params
}

func tomlFileExists(params Params) bool {
	tomlPath := params.ChDirPath + "/" + tfSwitchTOMLFileName
	return lib.CheckFileExist(tomlPath)
}
