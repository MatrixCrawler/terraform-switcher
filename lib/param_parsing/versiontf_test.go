package param_parsing

import (
	"fmt"
	"github.com/hashicorp/go-version"
	"github.com/warrensbox/terraform-switcher/lib"
	"strings"
	"testing"
)

func TestGetVersionFromVersionsTF_matches_version(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/integration-tests/test_versiontf"
	params, _ = GetVersionFromVersionsTF(params)
	v1, _ := version.NewVersion("1.0.5")
	actualVersion, _ := version.NewVersion(params.Version)
	if !actualVersion.Equal(v1) {
		t.Error("Determined version is not 1.0.5")
	}
}

func TestGetVersionFromVersionsTF_non_matching_constraints(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_non_matching_constraints"
	params, err := GetVersionFromVersionsTF(params)
	if err == nil {
		t.Error("Expected error but got nil")
	} else {
		expected := "exact constraint ("
		expected2 := ") cannot be combined with other conditions"
		if !strings.Contains(fmt.Sprint(err), expected) || !strings.Contains(fmt.Sprint(err), expected2) {
			t.Errorf("Expected error string containing %q and %q. Got %q", expected, expected2, err)
		}
	}
}

func TestGetVersionFromVersionsTF_erroneous_file(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_error"
	params, err := GetVersionFromVersionsTF(params)
	if err == nil {
		t.Error("Expected error got nil")
	} else {
		expected := "error parsing constraint: Malformed constraint: ~527> 1.0.0"
		if fmt.Sprint(err) != expected {
			t.Errorf("Expected error %q, got %q", expected, err)
		}
	}
}

func TestGetVersionFromVersionsTF_non_existent_constraint(t *testing.T) {
	logger = lib.InitLogger("DEBUG")
	var params Params
	params = initParams(params)
	params.ChDirPath = "../../test-data/skip-integration-tests/test_versiontf_non_existent"
	params, err := GetVersionFromVersionsTF(params)
	if err == nil {
		t.Error("Expected error got nil")
	} else {
		expected := "did not find version matching constraint: > 99999.0.0"
		if fmt.Sprint(err) != expected {
			t.Errorf("Expected error %q, got %q", expected, err)
		}
	}
}
