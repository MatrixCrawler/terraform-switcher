package param_parsing

import (
	"github.com/hashicorp/terraform-config-inspect/tfconfig"
	"github.com/warrensbox/terraform-switcher/lib"
)

const versionTfFileName = "version.tf"

func GetVersionFromVersionsTF(params Params) Params {
	logger.Infof("Reading version from terraform module at %q", params.ChDirPath)
	module, err := tfconfig.LoadModule(params.ChDirPath)
	if err != nil {
		logger.Fatalf("Could not load terraform module at %q", params.ChDirPath)
	}
	tfconstraint := module.RequiredCore[0]
	version, err2 := lib.GetSemver(tfconstraint, params.MirrorURL)
	if err2 != nil {
		logger.Fatalf("No version found matching %q", tfconstraint)
	}
	params.Version = version
	return params
}

func isTerraformModule(params Params) bool {
	module, err := tfconfig.LoadModule(params.ChDirPath)
	return err != nil && len(module.RequiredCore) > 0
}
