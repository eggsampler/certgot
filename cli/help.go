package cli

import (
	"fmt"
	"strings"
)

type HelpTopic struct {
	Topics          []string
	Name            string
	Description     string
	LongDescription string
	ShowTopic       string
	ShowFunc        func(*App, string) bool
}

func ShowAlways(*App, string) bool       { return true }
func ShowAnyTopic(_ *App, s string) bool { return s != "" }
func ShowNoTopic(_ *App, s string) bool  { return s == "" }
func ShowNotSubcommand(app *App, topic string) bool {
	for _, sc := range app.subCommands {
		if sc.Name == topic {
			return false
		}
	}
	return true
}

func DefaultHelpPrinter(app *App) {
	fmt.Println(strings.Repeat("- ", 40))
	defer func() {
		fmt.Println(strings.Repeat("- ", 40))
	}()

	specifiedTopic := ""

	argHelp := app.GetArgument("help")
	if argHelp != nil {
		specifiedTopic = argHelp.StringOrDefault()
	}

	// if there is no help topic, iterate all help topics
	// print help topic if showfunc evaluates to true
	if specifiedTopic == "" {
		foundTopic := false
		for _, helpTopic := range app.helpTopics {
			if helpTopic.ShowFunc != nil && helpTopic.ShowFunc(app, specifiedTopic) {
				foundTopic = true
				printHelpTopic(helpTopic)
			}
		}
		if !foundTopic {
			fmt.Printf("No help for specified topic: %s\n", specifiedTopic)
		}
		return
	}

	// check if topic is a subcommand and print usage + description
	if sc := app.subCommands[specifiedTopic]; sc != nil && len(sc.HelpTopic.Name) > 0 {
		fmt.Println()
		fmt.Println("  " + sc.HelpTopic.Name)
		fmt.Println()

		if len(sc.HelpTopic.Description) > 0 {
			fmt.Println(sc.HelpTopic.Description)
		}

		// print any non-specific help topics for the specified topic
		for _, helpTopic := range app.helpTopics {
			if helpTopic.ShowFunc != nil && helpTopic.ShowFunc(app, specifiedTopic) {
				printHelpTopic(helpTopic)
			}
		}
	}

	// then grab the helptopic and print that
	ht, ok := getHelpTopic(app.helpTopics, specifiedTopic)
	if ok {
		printHelpTopic(ht)
	}
}

func printHelpTopic(topic HelpTopic) {
	fmt.Println()
	if topic.Name != "" {
		fmt.Println(topic.Name + ":")
	}
	if topic.Description != "" {
		fmt.Println(topic.Description)
		fmt.Println()
	}
	if topic.LongDescription != "" {
		fmt.Println(topic.LongDescription)
		fmt.Println()
	}
}

func getHelpTopic(topics []HelpTopic, s string) (HelpTopic, bool) {
	for _, v := range topics {
		for _, t := range v.Topics {
			if t == s {
				return v, true
			}
		}
	}
	return HelpTopic{}, false
}
