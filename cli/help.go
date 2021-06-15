package cli

import (
	"fmt"
	"os"
	"strings"

	"github.com/eggsampler/certgot/log"
	"golang.org/x/term"
)

var termWidth = 80

func init() {
	w, _, err := term.GetSize(int(os.Stdout.Fd()))
	if err == nil {
		termWidth = w
	}
}

type HelpTopic struct {
	Topic           string
	Name            string
	Usage           string
	Description     string
	LongDescription string
	ShowTopic       string
	ShowFunc        func(*App, string) bool
}

func ShowAlways(*App, string) bool       { return true }
func ShowAnyTopic(_ *App, s string) bool { return s != "" }
func ShowNoTopic(_ *App, s string) bool  { return s == "" }
func ShowNotSubcommand(app *App, topic string) bool {
	_, ok := app.subCommandMap[topic]
	return !ok
}

func DefaultHelpPrinter(app *App, specifiedTopic string) {
	specifiedTopic = strings.ToLower(specifiedTopic)
	if specifiedTopic == "all" {
		specifiedTopic = ""
	}

	// check topic exists
	sc := app.subCommandMap[specifiedTopic]
	ht, foundHt := getHelpTopic(app.helpTopics, specifiedTopic)

	if specifiedTopic != "" && sc == nil && !foundHt {
		fmt.Printf("Unknown topic/command: %q\n", specifiedTopic)
		allTopics := []string{"all"}
		for _, t := range app.helpTopics {
			allTopics = append(allTopics, t.Topic)
		}
		fmt.Printf("Valid topics: %s\n", strings.Join(allTopics, ", "))
		var allCommands []string
		for k, _ := range app.subCommandMap {
			allCommands = append(allCommands, k)
		}
		fmt.Printf("Valid commands: %s\n", strings.Join(allCommands, ", "))
		return
	}

	// check if topic is a subcommand and print usage + description
	if sc != nil {
		fmt.Println("usage:")
		fmt.Println()

		if len(sc.Usage.LongUsage) > 0 {
			fmt.Println("  " + app.Name + " " + sc.Name + " " + sc.Usage.LongUsage)
		} else {
			fmt.Println("  " + app.Name + " " + sc.Name + " [options] ...")
		}
		fmt.Println()

		if len(sc.Usage.UsageDescription) > 0 {
			fmt.Println(sc.Usage.UsageDescription)
			fmt.Println()
		}
	}

	// then and print the helptopic if found
	if foundHt {
		printHelpTopic(app, ht)
	}

	// print any non-specific help topics for the specified topic (if present)
	for _, helpTopic := range app.helpTopics {
		if helpTopic.ShowFunc != nil && helpTopic.ShowFunc(app, specifiedTopic) {
			printHelpTopic(app, helpTopic)
		}
	}

	if sc != nil {
		printHelpSubCommand(app, sc)
	}
}

func printHelpSubCommand(app *App, sc *SubCommand) {
	fmt.Println(sc.Name + ":")
	if sc.Usage.ArgumentDescription != "" {
		fmt.Println(log.Wrap(sc.Usage.ArgumentDescription, termWidth, "  "))
		fmt.Println()
	}
	for _, argName := range sc.Usage.Flags {
		arg := app.Flag(argName)
		if arg == nil {
			// TODO: handle this more gracefully ?
			panic(fmt.Sprintf("help subcommand %q has no argument %q", sc.Name, argName))
		}
		printFlagHelp(arg)
	}
	fmt.Println()
}

func printHelpTopic(app *App, topic HelpTopic) {
	if topic.Name != "" {
		fmt.Println(topic.Name + ":")
	}
	if topic.Description != "" {
		fmt.Println(log.Wrap(topic.Description, termWidth, "  "))
		fmt.Println()
	}
	if topic.LongDescription != "" {
		fmt.Println(log.Wrap(topic.LongDescription, termWidth, ""))
		//fmt.Println()
	}
	for _, cmd := range app.subCommandList {
		if contains(cmd.HelpTopics, topic.Topic) {
			cmdName := cmd.Name
			if cmd.Default {
				cmdName = "(default) " + cmdName
			}
			printHelpLine(cmdName, cmd.Usage.UsageDescription)
		}
	}

	for _, arg := range app.flagsList {
		if contains(arg.HelpTopics, topic.Topic) {
			printFlagHelp(arg)
		}
	}
	fmt.Println()
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

func printFlagHelp(f *Flag) {
	argList := []string{
		strings.TrimSpace(flagDashes(f.Name) + f.Name + " " + f.Usage.ArgName),
	}
	for _, n := range f.AltNames {
		s := "-"
		if len(n) > 1 {
			s += "-"
		}
		s += n + " " + f.Usage.ArgName
		argList = append(argList, strings.TrimSpace(s))
	}
	args := strings.Join(argList, ", ")
	if strings.HasPrefix(args, "--") {
		// nothing
	} else {
		args = " " + args
	}

	// desc includes the argument description and any default value, if set
	desc := f.Usage.Description
	if f.DefaultValue != nil && f.DefaultValue.GetUsageDefault() != "" {
		// TODO: %s ? stringer something something
		desc += fmt.Sprintf(" (default: %v)", f.DefaultValue.GetUsageDefault())
	}

	printHelpLine(args, desc)
}

func getHelpTopic(topics []HelpTopic, s string) (HelpTopic, bool) {
	for _, v := range topics {
		if v.Topic == s {
			return v, true
		}
	}
	return HelpTopic{}, false
}

func contains(s []string, c string) bool {
	for _, v := range s {
		if v == c {
			return true
		}
	}
	return false
}
