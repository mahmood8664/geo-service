package config

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestLoadConfigs(t *testing.T) {
	loadConfigs("../resources/config.yml")
	assert.Greater(t, len(C.MongoDB.URL), 0)
	assert.Greater(t, len(C.HttpPort), 0)
}
