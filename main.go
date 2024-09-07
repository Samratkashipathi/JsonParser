package main

import (
	"fmt"
	"strconv"
)

// https://datatracker.ietf.org/doc/html/rfc8259#page-5
const (
	BeginObject    = '{'
	BeginArray     = '['
	EndObject      = '}'
	EndArray       = ']'
	NameSeparator  = ':'
	ValueSeparator = ','
)

type JSON interface{}

type Parser struct {
	input string
	pos   int
}

type ParseError struct {
	msg string
	pos int
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("Parse error at position %d: %s", e.pos, e.msg)
}

func NewParser(input string) *Parser {
	return &Parser{input, 0}
}

func (p *Parser) Parse() (JSON, error) {
	if len(p.input) <= 0 {
		fmt.Println("Empty String")
		return nil, nil
	}

	p.skipWhiteSpace()
	if p.pos >= len(p.input) {
		fmt.Println("Empty string")
		return nil, nil
	}

	value, err := p.parseValue()

	if err != nil {
		fmt.Printf("Error while parsing the value. Pos %d\n", p.pos)
		return nil, err
	}

	p.skipWhiteSpace()
	if p.pos < len(p.input) {
		fmt.Println("Trailing character at the end")
		return nil, err
	}

	fmt.Println(value)
	return nil, nil
}

func (p *Parser) parseValue() (JSON, error) {
	p.skipWhiteSpace()

	cur := p.input[p.pos]

	switch cur {
	case BeginObject:
		return p.parseObject()
	case '"':
		return p.parseString()
	case BeginArray:
		return p.parseArray()
	case 'f':
		return p.parseLiteral("false")
	case 't':
		return p.parseLiteral("true")
	case 'n':
		return p.parseLiteral("null")
	case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
		return p.parseNumber()
	default:
		return nil, nil
	}
}

func (p *Parser) parseObject() (JSON, error) {
	obj := make(map[string]JSON)
	p.pos++

	for {
		p.skipWhiteSpace()

		if p.pos >= len(p.input) {
			return nil, &ParseError{msg: "unexpected end of input", pos: p.pos}
		}

		if p.input[p.pos] == EndObject {
			p.pos++
			return obj, nil
		}

		key, err := p.parseString()
		if err != nil {
			return nil, err
		}

		p.skipWhiteSpace()

		if p.input[p.pos] != NameSeparator {
			return nil, &ParseError{msg: "expected : after key", pos: p.pos}
		}
		p.pos++

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		obj[key] = value

		if p.input[p.pos] == ValueSeparator {
			p.pos++
		}
	}
}

func (p *Parser) parseString() (string, error) {
	p.pos++
	start := p.pos

	for p.input[p.pos] != '"' {
		p.pos++
	}

	p.pos++

	return p.input[start : p.pos-1], nil
}

func (p *Parser) parseArray() ([]interface{}, error) {
	arr := make([]interface{}, 0)
	p.pos++

	for {
		p.skipWhiteSpace()

		value, err := p.parseValue()
		if err != nil {
			return nil, err
		}

		arr = append(arr, value)

		p.skipWhiteSpace()

		if p.input[p.pos] == EndArray {
			p.pos++
			return arr, nil
		}

		if p.input[p.pos] != ValueSeparator {
			return nil, &ParseError{msg: "Expected , in array value", pos: p.pos}
		}

		p.pos++
	}

}

func (p *Parser) parseLiteral(literal string) (interface{}, error) {
	start := p.pos
	p.pos++

	for {
		current := p.input[p.pos]
		switch current {
		case ValueSeparator, EndArray, EndObject:
			foundLiteral := p.input[start:p.pos]

			if foundLiteral != literal {
				return nil, &ParseError{msg: fmt.Sprintf("Expected %q, got %q", literal, foundLiteral), pos: p.pos}
			}

			switch literal {
			case "true":
				return true, nil
			case "false":
				return false, nil
			case "null":
				return nil, nil
			}

			return nil, &ParseError{msg: fmt.Sprintf("Expected %q, got %q", literal, foundLiteral), pos: p.pos}
		default:
			p.pos++
		}
	}
}

func (p *Parser) parseNumber() (int, error) {
	start := p.pos
	for {
		switch p.input[p.pos] {
		case 48, 49, 50, 51, 52, 53, 54, 55, 56, 57:
			p.pos++
		case ValueSeparator, EndArray, EndObject:
			return strconv.Atoi(p.input[start:p.pos])
		default:
			return 0, &ParseError{msg: fmt.Sprintf("Expected digit, got %q", p.input[p.pos]), pos: p.pos}
		}
	}
}

func (p *Parser) skipWhiteSpace() {
	for p.pos < len(p.input) && p.input[p.pos] == ' ' {
		p.pos = p.pos + 1
	}
}

func main() {
	s := `{"a": {"b": false}, "c": true, "d" : [1, null]}`
	p := NewParser(s)
	p.Parse()
}
