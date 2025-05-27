package config

import (
	"os"
	"testing"
)

func TestLoadEnvs_WithEnvDocker_SkipsGodotenv(t *testing.T) {
	os.Setenv("ENV", "docker")
	defer os.Unsetenv("ENV")

	os.Setenv("PORT", "8080")
	defer os.Unsetenv("PORT")

	envs, _ := LoadEnvs("../../../.env")

	if envs.Port != "8080" {
		t.Errorf("expected Port to be '8080', got '%s'", envs.Port)
	}
}
