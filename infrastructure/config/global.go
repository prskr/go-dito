package config

var current *AppConfig

func Current() *AppConfig {
	return current
}

func SetCurrent(c *AppConfig) {
	current = c
}
