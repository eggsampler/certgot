package cli2

import (
	"fmt"
	"regexp"
	"strings"
)

// parseArguments takes a list of arguments (ie, os.Args) and parses them into a Context,
// given a list of valid flags and commands.
// It is expected that the first element in argsToParse will be the binary name, and will be skipped
func parseArguments(argsToParse []string, ctx *Context, validFlags FlagList, validCommands CommandList) error {

	// if there are no flags or commands, don't bother parsing
	if len(argsToParse) <= 1 {
		return nil
	}

	// flg represents the current or previously parsed flag, if no value has been set on it
	// this is used to hold a flag if it takes a value, so the next argument can be applied as the value
	var lastFlag *Flag

	// parse all of the argumentes, excluding
	for _, arg := range argsToParse[1:] {

		// first check if it's a flag
		if strings.HasPrefix(arg, "-") {

			// grab the flag using the regexes
			flagMatch := extractFlag(arg)
			if len(flagMatch) < 2 {
				return fmt.Errorf("invalid flag: %s", arg)
			}

			// for both regexes, the flag name is the first match group
			flagName := flagMatch[1]

			// if the flag is a repeated short flag the name isn't the repeated characters,
			// so take the name as the first character and take a count of them
			var shortRepeatCount int
			if !strings.HasPrefix(arg, "--") {
				shortRepeatCount = len(flagName)
				flagName = string(flagName[0])
			}

			// grab the flag from the list, throw an error if it doesn't exist
			currentFlag := validFlags.Get(flagName)
			if currentFlag == nil {
				return fmt.Errorf("unknown flag: %s (%s)", flagName, flagMatch[1])
			}

			// check if the flag allows repeated short flags
			if shortRepeatCount > 1 && !currentFlag.AllowShortRepeat {
				return fmt.Errorf("flag doesn't allow repeated short flags: %s", flagName)
			}

			// check if the flag allows multiple
			existingFlag := ctx.Flags.Get(flagName)
			if existingFlag != nil && !currentFlag.AllowMultiple {
				return fmt.Errorf("flag doesn't allow multiples: %s", flagName)
			}

			// add the flag to the context
			ctx.Flags = append(ctx.Flags, currentFlag)

			// add the original flag name
			currentFlag.flags = append(currentFlag.flags, flagMatch[1])

			// if the argument is i flag with an inline value, set the flag's value now
			if strings.Contains(arg, "=") {

				// if the flag is not explicitly set to take a value, return an error
				if !currentFlag.TakesValue {
					return fmt.Errorf("flag doesn't take a value: %s", flagName)
				}
				currentFlag.values = append(currentFlag.values, flagMatch[2])

			} else if currentFlag.TakesValue {
				// if the flag takes a value, set the last flag a as the current flag
				// so the next argument can be set as the flag's value
				lastFlag = currentFlag
			}

		} else if lastFlag != nil {
			// if the last argument was a flag that takes a value,
			// set the current argument as that flag's value
			lastFlag.values = append(lastFlag.values, arg)

			// and reset the last flag to nil so we don't keep setting the value
			lastFlag = nil

		} else if ctx.Command != nil {
			// anything after a command has been found is added to the extra arguments
			ctx.ExtraArguments = append(ctx.ExtraArguments, arg)

		} else {
			// if no command has been found, attempt to locate it
			ctx.Command = validCommands.Get(arg)

			// and if there is no valid command with that name, return an error
			if ctx.Command == nil {
				return fmt.Errorf("invalid command: %s", arg)
			}
		}
	}

	// if the last flag is still set and it takes a value (it should, this is checked before lastFlag is set)
	// show an error because there are no arguments left to set the value
	if lastFlag != nil && lastFlag.TakesValue {
		return fmt.Errorf("flag requires value: %s", lastFlag.Name)
	}

	return nil
}

var (
	regFlagLong  = regexp.MustCompile(`^--([[:alnum:]]+(?:-[[:alnum:]]+)*)(?:=(.+))?$`)
	regFlagShort = regexp.MustCompile("^-(a+|b+|c+|d+|e+|f+|g+|h+|i+|j+|k+|l+|m+|n+|o+|p+|q+|r+|s+|t+|u+|v+|w+|x+|y+|z+)(?:=(.+))?$")
)

// extractFlag takes a given string which is already determined to start with a single dash character '-'
// and parses it using either of the regular expressions for a long flag (prefixed with two dashes),
// or a short flag (prefixed with a single dash)
func extractFlag(s string) []string {
	if strings.HasPrefix(s, "--") {
		return regFlagLong.FindStringSubmatch(s)
	}
	return regFlagShort.FindStringSubmatch(s)
}
