package config

import (
	"fmt"
	"testing"
)

func TestGetConfig(t *testing.T) {
	config := GetConfig()
	fmt.Print(config.XMLName)
}

func TestOspfconfig(t *testing.T) {
	config := GetConfig()
	config.OspfdConfig()
}
