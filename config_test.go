package apgo_test

import (
    "testing"
    "github.com/igkuz/apgo"
)

func TestNewConfig(t *testing.T) {
  expectedDbName := "sgap_test"
  expectedDbUser := "root"
  expectedDbPass := "123"
  expectedDbDialect := "mysql"
  config := apgo.NewConfig()

  if config.DB["name"] != expectedDbName {
    t.Error("Wrong DB name was configured")
  }

  if config.DB["user"] != expectedDbUser {
    t.Error("Wrong DB user was configured")
  }

  if config.DB["password"] != expectedDbPass {
    t.Error("Wrong DB password was configured")
  }
  if config.DB["dialect"] != expectedDbDialect {
    t.Error("Wrong DB dialect was configured")
  }
}
