package config_test

import (
    "testing"
    "../config"
)

func TestNew(t *testing.T) {
  expectedDbName := "sgap_test"
  config := config.NewConfig()
  if config.DB["name"] != expectedDbName {
    t.Error("Wrong DB name was configured")
  }
}
