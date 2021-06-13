package cli

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/eggsampler/certgot/log"
)

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

func DefaultHelpPrinter(app *App) {
	specifiedTopic := ""

	argHelp := app.GetArgument("help")
	if argHelp != nil {
		specifiedTopic = argHelp.StringOrDefault()
	}

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
}

func printHelpTopic(app *App, topic HelpTopic) {
	if topic.Name != "" {
		fmt.Println(topic.Name + ":")
	}
	if topic.Description != "" {
		fmt.Println(log.Wrap(topic.Description, 80, "  "))
		fmt.Println()
	}
	if topic.LongDescription != "" {
		fmt.Println(log.Wrap(topic.LongDescription, 80, ""))
		//fmt.Println()
	}
	w := tabwriter.NewWriter(os.Stdout, 0, 4, 0, '\t', 0)
	for _, cmd := range app.subCommandList {
		if contains(cmd.HelpTopics, topic.Topic) {
			name := cmd.Name
			if cmd.Default {
				name = "(default) " + name
			}
			_, _ = fmt.Fprintf(w, "    %s\t%s\n", name, cmd.Usage.UsageDescription)
		}
	}
	_ = w.Flush()

	for _, arg := range app.argsList {
		if contains(arg.HelpTopics, topic.Topic) {
			argList := []string{
				strings.TrimSpace("--" + arg.Name + " " + arg.Usage.ArgName),
			}
			for _, n := range arg.AltNames {
				s := "-"
				if len(n) > 1 {
					s += "-"
				}
				s += n + " " + arg.Usage.ArgName
				argList = append(argList, strings.TrimSpace(s))
			}
			args := strings.Join(argList, ", ")
			padding := 1
			if strings.HasPrefix(args, "--") {
				padding = 2
			}

			argsExtra := ""
			if arg.DefaultValue != nil {
				argsExtra = fmt.Sprintf(" (default: %v)", arg.DefaultValue.Get())
			}

			if len(args) > 18 {
				_, _ = fmt.Printf("%s%s\t\n\t\t%s%s\n", strings.Repeat(" ", padding), args, arg.Usage.Description, argsExtra)
				continue
			}

			_, _ = fmt.Fprintf(w, "%s%s\t%s%s\n", strings.Repeat(" ", padding), args, arg.Usage.Description, argsExtra)
		}
	}
	_ = w.Flush()
	fmt.Println()
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
