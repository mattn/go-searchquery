package searchquery

import (
	"fmt"
	"strings"
	"unicode"
)

// TokenType represents the type of a token
type TokenType int

const (
	TokenTerm TokenType = iota // Regular term or phrase (e.g., "hello" or "hello world")
	TokenAND
	TokenOR
	TokenLParen
	TokenRParen
	TokenEOF
)

// Token represents a parsed token
type Token struct {
	Type  TokenType
	Value string // Used only for Term tokens (including phrases)
}

// TermNode represents a term or phrase node in the AST
type TermNode struct {
	Phrase string // Full phrase if it contains spaces
}

func (TermNode) isNode() {}

// AndNode and OrNode represent boolean operator nodes
type AndNode struct {
	Left, Right Node
}

func (AndNode) isNode() {}

type OrNode struct {
	Left, Right Node
}

func (OrNode) isNode() {}

// Node is the interface for AST nodes
type Node interface {
	isNode()
}

// Parser parses the query string
type Parser struct {
	input   string
	pos     int
	tokens  []Token
	current int
}

// NewParser creates a new parser instance
func NewParser(query string) *Parser {
	p := &Parser{input: strings.TrimSpace(query)}
	p.tokenize()
	return p
}

// tokenize splits the query into tokens
func (p *Parser) tokenize() {
	input := p.input
	i := 0
	for i < len(input) {
		if unicode.IsSpace(rune(input[i])) {
			i++
			continue
		}

		if input[i] == '(' {
			p.tokens = append(p.tokens, Token{Type: TokenLParen})
			i++
			continue
		}
		if input[i] == ')' {
			p.tokens = append(p.tokens, Token{Type: TokenRParen})
			i++
			continue
		}

		// Match AND or OR (case-sensitive, uppercase only)
		if strings.HasPrefix(input[i:], "AND") && (i+3 == len(input) || unicode.IsSpace(rune(input[i+3])) || input[i+3] == '(' || input[i+3] == ')') {
			p.tokens = append(p.tokens, Token{Type: TokenAND})
			i += 3
			continue
		}
		if strings.HasPrefix(input[i:], "OR") && (i+2 == len(input) || unicode.IsSpace(rune(input[i+2])) || input[i+2] == '(' || input[i+2] == ')') {
			p.tokens = append(p.tokens, Token{Type: TokenOR})
			i += 2
			continue
		}

		// Quoted phrase or unquoted term
		if input[i] == '"' {
			i++ // skip opening quote
			start := i
			for i < len(input) && input[i] != '"' {
				i++
			}
			if i >= len(input) {
				// Unclosed quote: treat remaining text as term (simple fallback)
				p.tokens = append(p.tokens, Token{Type: TokenTerm, Value: input[start:]})
				i = len(input)
			} else {
				p.tokens = append(p.tokens, Token{Type: TokenTerm, Value: input[start:i]})
				i++ // skip closing quote
			}
			continue
		}

		// Regular term (key:value extensions are treated as terms and ignored)
		start := i
		for i < len(input) && !unicode.IsSpace(rune(input[i])) && input[i] != '(' && input[i] != ')' {
			i++
		}
		term := input[start:i]
		if term != "" {
			p.tokens = append(p.tokens, Token{Type: TokenTerm, Value: term})
		}
	}
	p.tokens = append(p.tokens, Token{Type: TokenEOF})
}

// Parse builds the AST using a shunting-yard-like approach
func (p *Parser) Parse() (Node, error) {
	var stack []Node
	var opStack []Token

	for {
		token := p.nextToken()
		if token.Type == TokenEOF {
			break
		}
		if token.Type == TokenTerm {
			stack = append(stack, TermNode{Phrase: token.Value})
		} else if token.Type == TokenLParen {
			opStack = append(opStack, token)
		} else if token.Type == TokenRParen {
			for len(opStack) > 0 && opStack[len(opStack)-1].Type != TokenLParen {
				var err error
				stack, err = p.applyOp(stack, opStack)
				if err != nil {
					return nil, err
				}
			}
			if len(opStack) == 0 {
				return nil, fmt.Errorf("mismatched parentheses")
			}
			opStack = opStack[:len(opStack)-1] // pop LParen
		} else if token.Type == TokenAND || token.Type == TokenOR {
			for len(opStack) > 0 {
				top := opStack[len(opStack)-1]
				if top.Type == TokenLParen {
					break
				}
				// Apply higher or equal precedence operators
				if (token.Type == TokenOR && top.Type == TokenAND) ||
					(token.Type == TokenAND && top.Type == TokenAND) {
					var err error
					stack, err = p.applyOp(stack, opStack)
					if err != nil {
						return nil, err
					}
				} else {
					break
				}
			}
			opStack = append(opStack, token)
		}
	}

	// Apply remaining operators
	for len(opStack) > 0 {
		if opStack[len(opStack)-1].Type == TokenLParen {
			return nil, fmt.Errorf("mismatched parentheses")
		}
		var err error
		stack, err = p.applyOp(stack, opStack)
		if err != nil {
			return nil, err
		}
	}

	// If no explicit operators, apply implicit AND (required minimum behavior)
	if len(stack) > 1 {
		for len(stack) > 1 {
			left := stack[0]
			right := stack[1]
			stack = append([]Node{AndNode{Left: left, Right: right}}, stack[2:]...)
		}
	}

	if len(stack) != 1 {
		return nil, fmt.Errorf("invalid expression")
	}
	return stack[0], nil
}

// applyOp pops an operator and applies it to the top two nodes on the stack
func (p *Parser) applyOp(stack []Node, opStack []Token) ([]Node, error) {
	if len(stack) < 2 {
		return nil, fmt.Errorf("invalid expression")
	}
	op := opStack[len(opStack)-1]
	opStack = opStack[:len(opStack)-1]
	right := stack[len(stack)-1]
	left := stack[len(stack)-2]
	stack = stack[:len(stack)-2]

	var node Node
	if op.Type == TokenAND {
		node = AndNode{Left: left, Right: right}
	} else {
		node = OrNode{Left: left, Right: right}
	}
	stack = append(stack, node)
	return stack, nil
}

// nextToken returns the next token in the stream
func (p *Parser) nextToken() Token {
	if p.current >= len(p.tokens) {
		return Token{Type: TokenEOF}
	}
	token := p.tokens[p.current]
	p.current++
	return token
}

// Match determines whether the event content matches the query
func Match(content, query string) (bool, error) {
	if query == "" {
		return true, nil
	}
	parser := NewParser(query)
	ast, err := parser.Parse()
	if err != nil {
		return false, err
	}
	contentLower := strings.ToLower(content)
	return eval(ast, contentLower), nil
}

// eval recursively evaluates the AST against the lowercase content
func eval(node Node, contentLower string) bool {
	switch n := node.(type) {
	case TermNode:
		termLower := strings.ToLower(n.Phrase)
		if strings.Contains(n.Phrase, " ") || strings.HasPrefix(n.Phrase, "\"") {
			// Exact phrase match (contiguous sequence)
			phraseLower := strings.ToLower(n.Phrase)
			return strings.Contains(contentLower, phraseLower)
		}
		// Single term: substring match
		return strings.Contains(contentLower, termLower)
	case AndNode:
		return eval(n.Left, contentLower) && eval(n.Right, contentLower)
	case OrNode:
		return eval(n.Left, contentLower) || eval(n.Right, contentLower)
	default:
		return false
	}
}
