package cli

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

	// lastFlagExpectingValue represents the previously parsed flag, if it takes a value but no value has been set on it
	// this is used to hold a flag if it takes a value, so the next argument can be applied as the value to that flag
	var lastFlagExpectingValue *Flag

	// parse all of the arguments, excluding
	for _, arg := range argsToParse[1:] {

		// first check if it's a flag
		if strings.HasPrefix(arg, "-") {

			// if the previous flag expected a value, but we're now parsing a flag
			if lastFlagExpectingValue != nil && lastFlagExpectingValue.RequiresValue {
				lastValue := lastFlagExpectingValue.valuesInfo[len(lastFlagExpectingValue.valuesInfo)-1]
				return fmt.Errorf("flag %q (%s) requires value, but none provided before next flag %q",
					lastValue.RawFlag, lastValue.FlagName, arg)
			}

			// unset the last flag expecting a value, as we're past that flag and onto the current flag
			lastFlagExpectingValue = nil

			// grab the flag parts using the regexes
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

			// add the flag to the context, using FlagList.Put to make sure it's not duplicated
			ctx.Flags.Put(currentFlag)

			hasValue := false
			var value string

			// if the argument is a flag with an inline value, set the flag's value now
			if flagMatch[2] == "=" {

				// if the flag is not explicitly set to take a value, return an error
				if !currentFlag.TakesValue {
					return fmt.Errorf("flag doesn't take a value: %s", flagName)
				}

				// the flag has a value
				hasValue = true
				value = flagMatch[3]

				// add the current value to the list of all values
				currentFlag.valuesRaw = append(currentFlag.valuesRaw, flagMatch[3])

			} else if currentFlag.TakesValue {
				// set the variable so we can put an arg on it
				lastFlagExpectingValue = currentFlag
			}

			// and the information about the value
			currentFlag.valuesInfo = append(currentFlag.valuesInfo, FlagValue{
				FlagName: flagName,
				RawFlag:  flagMatch[1],
				HasValue: hasValue,
				Value:    value,
			})

		} else if lastFlagExpectingValue != nil {
			// if the last argument was a flag that is expecting a value

			// set the value
			lastFlagExpectingValue.valuesRaw = append(lastFlagExpectingValue.valuesRaw, arg)

			// and include the value in the flag value information list
			endIdx := len(lastFlagExpectingValue.valuesInfo) - 1
			lastValue := lastFlagExpectingValue.valuesInfo[endIdx]
			lastValue.HasValue = true
			lastValue.Value = arg
			lastFlagExpectingValue.valuesInfo[endIdx] = lastValue

			// and reset the last flag to nil so we don't keep setting the value
			lastFlagExpectingValue = nil

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

	// if the previous flag expected a value, but there's no more arguments left
	if lastFlagExpectingValue != nil && lastFlagExpectingValue.RequiresValue {
		lastValue := lastFlagExpectingValue.valuesInfo[len(lastFlagExpectingValue.valuesInfo)-1]
		return fmt.Errorf("flag %q (%s) requires value, but more arguments were provided",
			lastValue.RawFlag, lastValue.FlagName)
	}

	return nil
}

var (
	regFlagLong  = regexp.MustCompile(`^--([[:alnum:]]+(?:-[[:alnum:]]+)*)(?:(=)(.+))?$`)
	regFlagShort = regexp.MustCompile("^-(a+|b+|c+|d+|e+|f+|g+|h+|i+|j+|k+|l+|m+|n+|o+|p+|q+|r+|s+|t+|u+|v+|w+|x+|y+|z+)(?:(=)(.+))?$")
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
