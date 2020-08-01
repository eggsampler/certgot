package parser

type Result struct {
	Success   bool
	String    string
	Remaining Input
}

type Parser func(Input) Result

func ZeroOrMore(parser Parser) Parser {
	return func(input Input) Result {
		var result = Result{
			Success:   true,
			Remaining: input,
		}
		for result.Remaining != nil {
			var oneMoreResult = parser(result.Remaining)
			if !oneMoreResult.Success {
				return result
			}
			result.String += oneMoreResult.String
			result.Remaining = oneMoreResult.Remaining
		}
		return result
	}
}

func OneOrMore(parser Parser) Parser {
	return func(input Input) Result {
		var result = ZeroOrMore(parser)(input)
		if len(result.String) > 0 {
			return result
		}
		return Result{Success: false, Remaining: input}
	}
}

func Or(parsers ...Parser) Parser {
	return func(input Input) Result {
		for _, parser := range parsers {
			result := parser(input)
			if result.Success {
				return result
			}
		}
		return Result{Success: false, Remaining: input}
	}
}

func And(parsers ...Parser) Parser {
	return func(input Input) Result {
		finalResult := Result{
			Success:   true,
			Remaining: input,
		}
		for _, currentParser := range parsers {
			currentResult := currentParser(finalResult.Remaining)
			if !currentResult.Success {
				return Result{Success: false, Remaining: input}
			}
			finalResult.String += currentResult.String
			finalResult.Remaining = currentResult.Remaining
		}
		return finalResult
	}
}

func Not(parser Parser) Parser {
	return func(input Input) Result {
		result := parser(input)
		return Result{
			Success:   !result.Success,
			Remaining: input,
		}
	}
}
