package main

import (
	"runtime/debug"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResolveVersion_KeepsTheVersionInjectedAtBuildTime(t *testing.T) {
	info := &debug.BuildInfo{Main: debug.Module{Version: "v9.9.9"}}

	assert.Equal(t, "0.6.0", resolveVersion("0.6.0", info, true))
}

func TestResolveVersion_UsesTheModuleVersion_WhenInstalledWithGoInstall(t *testing.T) {
	info := &debug.BuildInfo{Main: debug.Module{Version: "v0.6.0"}}

	assert.Equal(t, "0.6.0", resolveVersion("dev", info, true))
}

func TestResolveVersion_StaysDev_WhenBuiltFromASourceCheckout(t *testing.T) {
	for _, info := range []*debug.BuildInfo{{Main: debug.Module{Version: "(devel)"}}, {}, nil} {
		assert.Equal(t, "dev", resolveVersion("dev", info, info != nil))
	}
}
