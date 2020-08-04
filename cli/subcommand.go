package cli

type SubCommand struct {
	Name    string
	Default bool
	Run     func(app *App) error
}
