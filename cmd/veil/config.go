package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type Config struct {
	Template  string
	Directory string
}

func readFileAsString(filePath string) (string, error) {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return "", nil // File doesn't exist, return empty string
	} else if err != nil {
		return "", err // Other error occurred
	}

	// File exists, read its contents
	data, err := ioutil.ReadFile(filePath)
	if err != nil {
		return "", err
	}

	return string(data), nil
}

func lookupConfig() *Config {

	c := &Config{
		Directory: ".",
	}

	if s := os.Getenv("VEIL_CONFIG"); s != "" {
		c.ParseConfig(s)
	} else if s := os.Getenv("VEIL_CONFIG_FILE"); s != "" {
		c.ParseConfigFile(s)
	} else {
		s, err := readFileAsString("~/.veil")
		if err != nil && s != "" {
			c.ParseConfig(s)
		}
	}

	return c
}

func (c *Config) ParseConfigFile(f string) {
	s, err := readFileAsString(f)
	if err != nil && s != "" {
		c.ParseConfig(s)
	}
}

func (c *Config) ParseConfig(s string) {
	cl := flag.NewFlagSet("", flag.ExitOnError)
	template := cl.String("t", "gokit", "Template variable")
	dir := cl.String("d", ".", "Directory")
	cl.Parse(strings.Fields(s))
	c.Set(template, dir)
}

func (c *Config) Set(template *string, dir *string) {
	if template != nil {
		c.Template = *template
	}
	if dir != nil {
		c.Directory = *dir
	}
}

func (c *Config) GetTemplateString() string {
	// Load templates
	var templateStr []byte

	switch c.Template {
	case "gokit":
		templateStr = gokit_service
	case "rpc":
		templateStr = rpc_service
	default:
		panic(fmt.Sprintf("Unsupported template: %s", c.Template))
	}
	return string(templateStr)
}
