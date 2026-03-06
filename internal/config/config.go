package config

type Config struct {
	Name string
}

func Load(filename string) Config {
	return Config{
		Name: "My Michel Site",
	}
}
