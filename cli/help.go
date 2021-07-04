package cli

import (
	"fmt"
	"strings"

	"github.com/eggsampler/certgot/log"
)

type HelpCategories []*HelpCategory

func (hc HelpCategories) Get(s string) *HelpCategory {
	for _, c := range hc {
		if strings.EqualFold(c.Name, s) {
			return c
		}
	}
	return nil
}

type HelpCategory struct {
	// Category is the name given to this category, and this is what is used when looking up help
	Category string

	// Name is a longer name for the category, and is the first thing shown when printing help for the category
	Name string

	Description string

	Usage string

	UsageDescription string

	// ShowFunc is used to determine whether this category should be shown to the user
	ShowFunc func(ctx *Context, category string) bool
}

func ShowAlways(*Context, string) bool              { return true }
func ShowAnyCategory(_ *Context, s string) bool     { return len(s) > 0 }
func ShowNoCategory(_ *Context, s string) bool      { return len(s) == 0 }
func ShowNoCommand(ctx *Context, topic string) bool { return ctx.App.Commands.Get(topic) == nil }

func DefaultHelpPrinter(ctx *Context, requestedCategory string) {
	requestedCategory = strings.ToLower(requestedCategory)
	if requestedCategory == "all" {
		requestedCategory = ""
	}

	// check category exists
	cmd := ctx.App.Commands.Get(requestedCategory)
	category := ctx.App.Help.Get(requestedCategory)

	if len(requestedCategory) > 0 && cmd == nil && category == nil {
		fmt.Printf("Unknown topic/command: %q\n", requestedCategory)
		allTopics := []string{"all"}
		for _, t := range ctx.App.Help {
			allTopics = append(allTopics, t.Category)
		}
		fmt.Printf("Valid topics: %s\n", strings.Join(allTopics, ", "))
		var allCommands []string
		for _, c := range ctx.App.Commands {
			allCommands = append(allCommands, c.Name)
		}
		fmt.Printf("Valid commands: %s\n", strings.Join(allCommands, ", "))
		return
	}

	if requestedCategory == "" {
		fmt.Println("showing help for 'all'")
	} else {
		fmt.Printf("showing help for '%s'\n", requestedCategory)
	}
	fmt.Println()

	// check if topic is a command and print usage + description
	if cmd != nil {
		fmt.Println("usage:")
		fmt.Println()

		if len(cmd.Usage) > 0 {
			fmt.Println("  " + ctx.App.Name + " " + cmd.Name + " " + cmd.Usage)
		} else {
			fmt.Println("  " + ctx.App.Name + " " + cmd.Name + " [options] ...")
		}
		fmt.Println()

		if len(cmd.UsageDescription) > 0 {
			fmt.Println(cmd.UsageDescription)
			fmt.Println()
		}
	}

	// then print the category if found
	if category != nil {
		printHelpCategory(ctx, category)
	}

	// print any non-specific help topics for the specified topic (if present)
	for _, cat := range ctx.App.Help {
		if cat.ShowFunc != nil && cat.ShowFunc(ctx, requestedCategory) {
			printHelpCategory(ctx, cat)
		}
	}

	if cmd != nil {
		printHelpCommand(ctx, cmd)
	}
}

func printHelpCategory(ctx *Context, category *HelpCategory) {
	if len(category.Name) > 0 {
		fmt.Println(category.Name + ":")
	}

	if len(category.Usage) > 0 {
		fmt.Println(log.Wrap(category.Usage, termWidth, "  "))
		fmt.Println()
		if len(category.UsageDescription) > 0 {
			fmt.Println(log.Wrap(category.UsageDescription, termWidth, ""))
		}
	}

	if len(category.Description) > 0 {
		fmt.Println(log.Wrap(category.Description, termWidth, "  "))
	}

	for _, cmd := range ctx.App.Commands {
		if contains(cmd.HelpCategories, category.Category) {
			cmdName := cmd.Name
			if cmd.Default {
				cmdName = "(default) " + cmdName
			}
			printHelpLine(cmdName, cmd.UsageDescription)
		}
	}

	for _, flg := range ctx.App.Flags {
		if contains(flg.HelpCategories, category.Category) {
			printFlagHelp(ctx, flg)
		}
	}
	fmt.Println()
}

func printHelpCommand(ctx *Context, cmd *Command) {
	fmt.Println(cmd.Name + ":")
	if len(cmd.ArgumentDescription) > 0 {
		fmt.Println(log.Wrap(cmd.ArgumentDescription, termWidth, "  "))
		fmt.Println()
	}
	for _, flagName := range cmd.HelpFlags {
		arg := ctx.App.Flags.Get(flagName)
		if arg == nil {
			// TODO: handle this more gracefully ?
			panic(fmt.Sprintf("flag %q does not exist in app for command %q", flagName, cmd.Name))
		}
		printFlagHelp(ctx, arg)
	}
	fmt.Println()
}

func flagDashes(s string) string {
	if len(s) == 1 {
		return "-"
	}
	return "--"
}

func printFlagHelp(ctx *Context, f *Flag) {
	argList := []string{
		strings.TrimSpace(flagDashes(f.Name) + f.Name + " " + f.HelpValueName),
	}
	for _, n := range f.AltNames {
		s := "-"
		if len(n) > 1 {
			s += "-"
		}
		s += n + " " + f.HelpValueName
		argList = append(argList, strings.TrimSpace(s))
	}
	args := strings.Join(argList, ", ")
	if strings.HasPrefix(args, "--") {
		// nothing
	} else {
		args = " " + args
	}

	// desc includes the argument description and any default value, if set
	desc := f.HelpDescription
	if f.HelpDefault != nil {
		defaultValueName, err := f.HelpDefault(ctx)
		if err != nil {
			// TODO: return this error?
			log.WithField("flag", f.Name).WithError(err).Error("fetching help default string")
		}
		if len(defaultValueName) > 0 {
			desc += fmt.Sprintf(" (default: %s)", defaultValueName)
		}
	}

	printHelpLine(args, desc)
}

func printHelpLine(flagOrCmdName, desc string) {
	flagOrCmdName = "  " + flagOrCmdName
	descPrefix := strings.Repeat(" ", 20)

	// if the length of the flag/command is greater than 20, put the description on the next line
	if len(flagOrCmdName) > 20 {
		// print the argument/cmd (ie, --hello THING)
		fmt.Println(flagOrCmdName)

		// print the description for the argument flag
		if len(desc) > termWidth {
			lines := log.WrapSlice(desc, termWidth, descPrefix)
			for _, line := range lines {
				fmt.Println(line)
			}
		} else {
			fmt.Println(descPrefix + desc)
		}

		return
	}

	combinedLine := flagOrCmdName + strings.Repeat(" ", 20-len(flagOrCmdName)) + desc
	if len(combinedLine) > termWidth {
		lines := log.WrapSlice(combinedLine, termWidth, "")
		fmt.Println(lines[0])
		for _, line := range lines[1:] {
			fmt.Println(descPrefix + line)
		}
	} else {
		fmt.Println(combinedLine)
	}
}
