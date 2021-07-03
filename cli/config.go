package cli

import (
	"strconv"
	"strings"
)

type ConfigList []*Config

func (cl ConfigList) Get(name string) *Config {
	for _, cfg := range cl {
		if cfg == nil {
			continue
		}
		if strings.EqualFold(cfg.Name, name) {
			return cfg
		}
	}
	return nil
}

type Config struct {
	Name        string
	Default     []string
	HelpDefault string
	OnSet       func(*Config, []string) error

	value []string
	isSet bool
}

func (c Config) IsSet() bool {
	return c.isSet
}

func (c *Config) set(v []string) error {
	if c.OnSet != nil {
		err := c.OnSet(c, v)
		if err != nil {
			return err
		}
	}

	c.value = v
	c.isSet = true

	return nil
}

func (c Config) Bool() bool {
	return c.isSet
}

func (c Config) Int() int {
	if len(c.value) > 0 {
		i, _ := strconv.Atoi(c.value[0])
		return i
	} else if len(c.Default) > 0 {
		i, _ := strconv.Atoi(c.Default[0])
		return i
	}
	return 0
}

func (c Config) String() string {
	if len(c.value) > 0 {
		return c.value[0]
	} else if len(c.Default) > 0 {
		return c.Default[0]
	}
	return ""
}

func (c Config) StringSlice() []string {
	if len(c.value) > 0 {
		return c.value
	} else if len(c.Default) > 0 {
		return c.Default
	}
	return c.value
}
