package param_parsing

import (
	"fmt"
	"github.com/pborman/getopt"
	"github.com/warrensbox/terraform-switcher/lib"
	"os"
)

type Params struct {
	ChDirPath        string
	CustomBinaryPath string
	DefaultVersion   string
	HelpFlag         bool
	LatestFlag       bool
	LatestPre        string
	LatestStable     string
	ListAllFlag      bool
	MirrorURL        string
	ShowLatestFlag   bool
	ShowLatestPre    string
	ShowLatestStable string
	Version          string
	VersionFlag      bool
}

var logger = lib.InitLogger()

func GetParameters() Params {
	var params Params
	params = initParams(params)

	getopt.StringVarLong(&params.ChDirPath, "chdir", 'c', "Switch to a different working directory before executing the given command. Ex: tfswitch --chdir terraform_project will run tfswitch in the terraform_project directory")
	getopt.StringVarLong(&params.CustomBinaryPath, "bin", 'b', "Custom binary path. Ex: tfswitch -b "+lib.ConvertExecutableExt("/Users/username/bin/terraform"))
	getopt.StringVarLong(&params.DefaultVersion, "default", 'd', "Default to this version in case no other versions could be detected. Ex: tfswitch --default 1.2.4")
	getopt.BoolVarLong(&params.HelpFlag, "help", 'h', "Displays help message")
	getopt.BoolVarLong(&params.LatestFlag, "latest", 'u', "Get latest stable version")
	getopt.StringVarLong(&params.LatestPre, "latest-pre", 'p', "Latest pre-release implicit version. Ex: tfswitch --latest-pre 0.13 downloads 0.13.0-rc1 (latest)")
	getopt.StringVarLong(&params.LatestStable, "latest-stable", 's', "Latest implicit version based on a constraint. Ex: tfswitch --latest-stable 0.13.0 downloads 0.13.7 and 0.13 downloads 0.15.5 (latest)")
	getopt.BoolVarLong(&params.ListAllFlag, "list-all", 'l', "List all versions of terraform - including beta and rc")
	getopt.StringVarLong(&params.MirrorURL, "mirror", 'm', "Install from a remote API other than the default. Default: "+lib.DefaultMirror)
	getopt.BoolVarLong(&params.ShowLatestFlag, "show-latest", 'U', "Show latest stable version")
	getopt.StringVarLong(&params.ShowLatestPre, "show-latest-pre", 'P', "Show latest pre-release implicit version. Ex: tfswitch --show-latest-pre 0.13 prints 0.13.0-rc1 (latest)")
	getopt.StringVarLong(&params.ShowLatestStable, "show-latest-stable", 'S', "Show latest implicit version. Ex: tfswitch --show-latest-stable 0.13 prints 0.13.7 (latest)")
	getopt.BoolVarLong(&params.VersionFlag, "version", 'v', "Displays the version of tfswitch")

	// Parse the command line parameters to fetch stuff like chdir
	getopt.Parse()

	// Read configuration files
	if tomlFileExists(params) {
		params = getParamsTOML(params)
	} else if tfSwitchFileExists(params) {
		params = GetParamsFromTfSwitch(params)
	} else if terraformVersionFileExists(params) {
		params = GetParamsFromTerraformVersion(params)
	} else if isTerraformModule(params) {
		params = GetVersionFromVersionsTF(params)
	} else if terraGruntFileExists(params) {
		params = GetVersionFromTerragrunt(params)
	} else {
		params = GetParamsFromEnvironment(params)
	}

	// Parse again to overwrite anything that might by defined on the cli AND in any config file (CLI always wins)
	getopt.Parse()
	args := getopt.Args()
	if len(args) == 1 {
		/* version provided on command line as arg */
		params.Version = args[0]
	}
	return params
}

func initParams(params Params) Params {
	params.ChDirPath = lib.GetCurrentDirectory()
	params.CustomBinaryPath = lib.ConvertExecutableExt(lib.GetDefaultBin())
	params.DefaultVersion = lib.DefaultLatest
	params.HelpFlag = false
	params.LatestFlag = false
	params.LatestPre = lib.DefaultLatest
	params.LatestStable = lib.DefaultLatest
	params.ListAllFlag = false
	params.MirrorURL = lib.DefaultMirror
	params.ShowLatestFlag = false
	params.ShowLatestPre = lib.DefaultLatest
	params.ShowLatestStable = lib.DefaultLatest
	params.Version = lib.DefaultLatest
	params.VersionFlag = false
	return params
}

func UsageMessage() {
	fmt.Print("\n\n")
	getopt.PrintUsage(os.Stderr)
	fmt.Println("Supply the terraform version as an argument, or choose from a menu")
}
