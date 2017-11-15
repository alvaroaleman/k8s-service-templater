package main

type AppConfig struct {
	Template            string `yaml:"template"`
	TemplateDestination string `yaml:"template_destination"`
	Command             string `yaml:"command"`
}
