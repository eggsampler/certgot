package cli

import "fmt"

func SetConfigValue(name string) func(f *Flag, ctx *Context) error {
	return func(f *Flag, ctx *Context) error {
		cfg := ctx.App.Configs.Get(name)
		if cfg == nil {
			return fmt.Errorf("no config %q for flag %q", name, f.Name)
		}
		return cfg.set(f.valuesRaw, ConfigSource{
			Source: SourceFlag,
			Extra:  f.Name,
		})
	}
}

func GetConfigDefault(name string) func(*Context) (string, error) {
	return func(ctx *Context) (string, error) {
		cfg := ctx.App.Configs.Get(name)
		if cfg == nil {
			return "", fmt.Errorf("invalid config: %s", name)
		}
		return cfg.HelpDefault, nil
	}
}
