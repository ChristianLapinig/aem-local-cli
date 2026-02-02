package models

type Config struct {
	EnvsPath     string        `json:"envsPath"`
	Environments []Environment `json:"environments"`
}

type Environment struct {
	Name string `json:"name"`
	Path string `json:"path"`
}
