package config

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestGetKeptnGoUtilsConfig(t *testing.T) {
	got := GetKeptnGoUtilsConfig()
	require.NotEmpty(t, got.ShKeptnSpecVersion)
}
