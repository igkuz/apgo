package config

type APConfig struct {
  DB    map[string]string
}

func NewConfig() (*APConfig) {
  db := make(map[string]string)
  db["name"] = "sgap_test"
  return &APConfig{ DB: db }
}
