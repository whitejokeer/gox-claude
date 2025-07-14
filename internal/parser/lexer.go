package parser

import (
	"fmt"
	"regexp"
	"strings"
)

// TokenType represents the type of a token
type TokenType int

const (
	TokenSectionStart TokenType = iota
	TokenSectionEnd
	TokenContent
	TokenEOF
)

// Token represents a lexical token
type Token struct {
	Type       TokenType
	Value      string
	Attributes map[string]string
	Line       int
	Column     int
}

// Lexer tokenizes .gox file content
type Lexer struct {
	content  []byte
	filename string
	pos      int
	line     int
	column   int
}

// NewLexer creates a new lexer
func NewLexer(content []byte, filename string) *Lexer {
	return &Lexer{
		content:  content,
		filename: filename,
		pos:      0,
		line:     1,
		column:   1,
	}
}

// Tokenize processes the content and returns tokens
func (l *Lexer) Tokenize() ([]Token, error) {
	var tokens []Token
	inSection := false
	
	for l.pos < len(l.content) {
		// Skip whitespace at the beginning
		if l.isWhitespace(l.current()) {
			l.skipWhitespace()
			continue
		}
		
		// Check for section tags
		if l.current() == '<' {
			// Look ahead to see if this is a .gox section tag
			remaining := string(l.content[l.pos:])
			
			if l.isOpeningSectionTag(remaining) && !inSection {
				// Opening section tag
				token, err := l.parseTag()
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, token)
				inSection = true
			} else if l.isClosingSectionTag(remaining) && inSection {
				// Closing section tag
				token, err := l.parseTag()
				if err != nil {
					return nil, err
				}
				tokens = append(tokens, token)
				inSection = false
			} else {
				// Regular content (HTML or other)
				token := l.parseContent()
				if strings.TrimSpace(token.Value) != "" {
					tokens = append(tokens, token)
				}
			}
		} else {
			// Parse content
			token := l.parseContent()
			if strings.TrimSpace(token.Value) != "" {
				tokens = append(tokens, token)
			}
		}
	}
	
	return tokens, nil
}

// current returns the current character
func (l *Lexer) current() byte {
	if l.pos >= len(l.content) {
		return 0
	}
	return l.content[l.pos]
}

// peek returns the character at offset from current position
func (l *Lexer) peek(offset int) byte {
	pos := l.pos + offset
	if pos >= len(l.content) {
		return 0
	}
	return l.content[pos]
}

// advance moves to the next character
func (l *Lexer) advance() {
	if l.pos < len(l.content) {
		if l.content[l.pos] == '\n' {
			l.line++
			l.column = 1
		} else {
			l.column++
		}
		l.pos++
	}
}

// isWhitespace checks if a character is whitespace
func (l *Lexer) isWhitespace(ch byte) bool {
	return ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r'
}

// skipWhitespace skips whitespace characters
func (l *Lexer) skipWhitespace() {
	for l.pos < len(l.content) && l.isWhitespace(l.current()) {
		l.advance()
	}
}

// parseTag parses an opening or closing tag
func (l *Lexer) parseTag() (Token, error) {
	startLine := l.line
	startColumn := l.column
	
	if l.current() != '<' {
		return Token{}, &ParseError{
			File:    l.filename,
			Line:    l.line,
			Column:  l.column,
			Message: "expected '<'",
		}
	}
	
	l.advance() // consume '<'
	
	isClosing := false
	if l.current() == '/' {
		isClosing = true
		l.advance() // consume '/'
	}
	
	// Parse tag name
	tagName := l.parseIdentifier()
	if tagName == "" {
		return Token{}, &ParseError{
			File:    l.filename,
			Line:    l.line,
			Column:  l.column,
			Message: "expected tag name",
		}
	}
	
	// Validate tag name
	if !l.isValidSectionName(tagName) {
		return Token{}, &ParseError{
			File:    l.filename,
			Line:    l.line,
			Column:  l.column,
			Message: fmt.Sprintf("invalid section name: %s", tagName),
		}
	}
	
	var attributes map[string]string
	
	if !isClosing {
		// Parse attributes for opening tags
		attributes = l.parseAttributes()
	}
	
	// Skip whitespace before '>'
	l.skipWhitespace()
	
	if l.current() != '>' {
		return Token{}, &ParseError{
			File:    l.filename,
			Line:    l.line,
			Column:  l.column,
			Message: "expected '>'",
		}
	}
	
	l.advance() // consume '>'
	
	tokenType := TokenSectionStart
	if isClosing {
		tokenType = TokenSectionEnd
	}
	
	return Token{
		Type:       tokenType,
		Value:      tagName,
		Attributes: attributes,
		Line:       startLine,
		Column:     startColumn,
	}, nil
}

// parseIdentifier parses an identifier (tag name)
func (l *Lexer) parseIdentifier() string {
	start := l.pos
	
	for l.pos < len(l.content) {
		ch := l.current()
		if (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') || ch == '-' || ch == '_' {
			l.advance()
		} else {
			break
		}
	}
	
	return string(l.content[start:l.pos])
}

// parseAttributes parses tag attributes
func (l *Lexer) parseAttributes() map[string]string {
	attributes := make(map[string]string)
	
	for l.pos < len(l.content) {
		l.skipWhitespace()
		
		if l.current() == '>' {
			break
		}
		
		// Parse attribute name
		name := l.parseIdentifier()
		if name == "" {
			break
		}
		
		l.skipWhitespace()
		
		value := ""
		if l.current() == '=' {
			l.advance() // consume '='
			l.skipWhitespace()
			value = l.parseAttributeValue()
		}
		
		attributes[name] = value
	}
	
	return attributes
}

// parseAttributeValue parses an attribute value
func (l *Lexer) parseAttributeValue() string {
	if l.current() == '"' {
		return l.parseQuotedString('"')
	} else if l.current() == '\'' {
		return l.parseQuotedString('\'')
	} else {
		// Unquoted value
		start := l.pos
		for l.pos < len(l.content) {
			ch := l.current()
			if l.isWhitespace(ch) || ch == '>' {
				break
			}
			l.advance()
		}
		return string(l.content[start:l.pos])
	}
}

// parseQuotedString parses a quoted string
func (l *Lexer) parseQuotedString(quote byte) string {
	if l.current() != quote {
		return ""
	}
	
	l.advance() // consume opening quote
	start := l.pos
	
	for l.pos < len(l.content) {
		if l.current() == quote {
			value := string(l.content[start:l.pos])
			l.advance() // consume closing quote
			return value
		}
		l.advance()
	}
	
	// Unclosed string - return what we have
	return string(l.content[start:l.pos])
}

// parseContent parses content between tags
func (l *Lexer) parseContent() Token {
	startLine := l.line
	startColumn := l.column
	start := l.pos
	
	for l.pos < len(l.content) {
		if l.current() == '<' {
			// Check if this is actually a .gox section tag
			if l.peek(1) == '/' {
				// Check if it's a closing tag for a valid section
				remaining := string(l.content[l.pos:])
				if l.isClosingSectionTag(remaining) {
					break
				}
			} else if l.isLetter(l.peek(1)) {
				// Check if it's an opening tag for a valid section
				remaining := string(l.content[l.pos:])
				if l.isOpeningSectionTag(remaining) {
					break
				}
			}
		}
		l.advance()
	}
	
	content := string(l.content[start:l.pos])
	
	return Token{
		Type:   TokenContent,
		Value:  content,
		Line:   startLine,
		Column: startColumn,
	}
}

// isOpeningSectionTag checks if the content starting at current position is a .gox section opening tag
func (l *Lexer) isOpeningSectionTag(content string) bool {
	if len(content) < 3 || content[0] != '<' {
		return false
	}
	
	// Extract tag name
	tagEnd := 1
	for tagEnd < len(content) && content[tagEnd] != '>' && content[tagEnd] != ' ' {
		tagEnd++
	}
	
	if tagEnd >= len(content) {
		return false
	}
	
	tagName := content[1:tagEnd]
	return l.isValidSectionName(tagName)
}

// isClosingSectionTag checks if the content starting at current position is a .gox section closing tag
func (l *Lexer) isClosingSectionTag(content string) bool {
	if len(content) < 4 || content[0] != '<' || content[1] != '/' {
		return false
	}
	
	// Extract tag name
	tagEnd := 2
	for tagEnd < len(content) && content[tagEnd] != '>' {
		tagEnd++
	}
	
	if tagEnd >= len(content) {
		return false
	}
	
	tagName := content[2:tagEnd]
	return l.isValidSectionName(tagName)
}

// isLetter checks if a character is a letter
func (l *Lexer) isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z')
}

// isValidSectionName checks if a section name is valid
func (l *Lexer) isValidSectionName(name string) bool {
	validSections := []string{"template", "go", "style"}
	for _, valid := range validSections {
		if name == valid {
			return true
		}
	}
	return false
}

// Regular expressions for validation
var (
	// Identifier regex: letters, numbers, underscores, hyphens
	identifierRegex = regexp.MustCompile(`^[a-zA-Z][a-zA-Z0-9_-]*$`)
)