package parser

import "unicode"

type ExpectFunc func(r rune) bool

func Expect(f ExpectFunc) Parser {
	return func(input Input) Result {
		if f(input.CurrentRune()) {
			return Result{
				Success:   true,
				String:    string(input.CurrentRune()),
				Remaining: input.RemainingInput(),
			}
		}
		return Result{Remaining: input}
	}
}

func ExpectRune(c rune) Parser {
	return Expect(func(r rune) bool {
		return c == r
	})
}

func Letter() Parser {
	return Expect(unicode.IsLetter)
}

func Number() Parser {
	return Expect(unicode.IsNumber)
}

func Space() Parser {
	return ExpectRune(' ')
}

func EOF() Parser {
	return func(input Input) Result {
		if input == nil || input.Finished() {
			return Result{
				Success: true,
			}
		}
		return Result{Remaining: input}
	}
}

func EOL() Parser {
	return Or(
		And(ExpectRune('\r'), ExpectRune('\n')),
		ExpectRune('\n'),
		EOF())
}
