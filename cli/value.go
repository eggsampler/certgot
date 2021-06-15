package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// DefaultValue represents a type used for a Flag.DefaultValue
// These can be used to represent simple types like a string, or boolean
// But also more complex types, for example ones that can prompt user for input
type DefaultValue interface {
	// Get returns the value set
	Get() interface{}

	// GetUsageDefault shows a default value when printing help
	GetUsageDefault() string
}

type SimpleValue struct {
	// Value is the value returned from the DefaultValue.Get function
	Value interface{}

	// UsageDefault is the value returned from the DefaultValue.GetUsageDefault function
	UsageDefault string
}

func (sv SimpleValue) Get() interface{} {
	return sv.Value
}

func (sv SimpleValue) GetUsageDefault() string {
	return sv.UsageDefault
}

type AskValue struct {
	Query  string
	Cancel string
	value  string
}

func (av *AskValue) Get() interface{} {
	if av.value == "" {
		fmt.Printf("%s (Enter 'c' to cancel): ", av.Query)
		av.value = readLine(av.Cancel)
	}
	return av.value
}

func (av AskValue) GetUsageDefault() string {
	return "Ask"
}

type ListValueOption struct {
	Option string
	Value  interface{}
}

type ListValue struct {
	Query   string
	Cancel  string
	Options []ListValueOption
	value   interface{}
}

func (lv *ListValue) Get() interface{} {
	if lv.value != nil {
		return lv.value
	}
	fmt.Println(lv.Query)
	fmt.Println(strings.Repeat("- ", 40))
	for k, v := range lv.Options {
		fmt.Printf("%d: %s (%v)\n", k+1, v.Option, v.Value)
	}
	fmt.Println(strings.Repeat("- ", 40))
	fmt.Printf("Select the appropriate number [1-%d] then [enter] (press 'c' to cancel): ", len(lv.Options))
	num, _ := strconv.Atoi(readLine(lv.Cancel))
	if num <= 0 || num > len(lv.Options) {
		fmt.Println("Invalid value")
		os.Exit(1)
	}
	lv.value = lv.Options[num-1].Value
	return lv.value
}

func (lv ListValue) GetUsageDefault() string {
	return "Ask"
}

func readLine(cancel string) string {
	scanner := bufio.NewScanner(os.Stdin)
	success := scanner.Scan()
	if err := scanner.Err(); err != nil {
		fmt.Printf("Error: %v\n", err)
		os.Exit(1)
		return ""
	}
	input := scanner.Text()
	if !success || strings.ToLower(input) == "c" {
		if cancel != "" {
			fmt.Println(cancel)
		}
		os.Exit(1)
		return ""
	}
	return input
}
