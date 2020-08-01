package parser

type Input interface {
	CurrentRune() rune
	RemainingInput() Input
}

func NewStringInput(text string) Input {
	return &StringInput{[]rune(text), 0}
}

type StringInput struct {
	text  []rune
	index int
}

func (i StringInput) String() string {
	return string(i.text[i.index:])
}

func (i StringInput) RemainingInput() Input {
	if i.index >= len(i.text)-1 {
		return nil
	}
	return StringInput{i.text, i.index + 1}
}

func (i StringInput) CurrentRune() rune {
	if i.index >= len(i.text) {
		return 0
	}
	return i.text[i.index]
}
