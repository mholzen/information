package set

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

type Statement = List[StringShape]
type Statements []Statement

func NewStatement(args ...StringShape) Statement {
	res := make(List[StringShape], 0)
	res = append(res, args...)
	return res
}

func parseString(content string) (Statements, error) {
	// scanner := bufio.NewScanner(strings.NewReader(content))
	return nil, nil
}

func parse(filename string) (Statements, error) {
	// Open the file for reading
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	// Create a scanner to read the file line by line
	scanner := bufio.NewScanner(file)

	// Create a regex pattern to match quoted words
	// Loop through each line of the file
	// Use the regex pattern to split the line into individual words,
	// including those enclosed in double quotes
	// Remove any quotes around words that were enclosed in double quotes
	// Check for any errors that may have occurred during scanning
	return newFunction(scanner)
}

func newFunction(scanner *bufio.Scanner) (Statements, error) {
	pattern := regexp.MustCompile(`"([^"]+)"|\S+|\t|    `)

	contexts := make(Statements, 0)
	contexts = append(contexts, NewStatement())

	indentCount := 0
	context := &contexts[indentCount]

	statements := make(Statements, 0)

	for scanner.Scan() {

		line := scanner.Text()
		words := pattern.FindAllString(line, -1)

		for i, word := range words {
			if strings.HasPrefix(word, `"`) && strings.HasSuffix(word, `"`) {
				words[i] = strings.Trim(word, `"`)
			}
		}
		tokens := getTokens(words)
		tokens = append(tokens, Token{TokenType: EOL})

		for _, token := range tokens {
			log.Printf("parsing %+v, indentCount: %d, contexts: %+v", token, indentCount, contexts)
			switch token.TokenType {
			case INDENT:
				indentCount++

			case WORD:
				contexts = contexts[0:indentCount]
				if indentCount >= len(contexts) {
					contexts = append(contexts, NewStatement())
				}
				context = &contexts[indentCount]
				if len(*context) >= 3 {
					return nil, fmt.Errorf("more than 3 words on line")
				}
				*context = append(*context, StringShape(token.Word))

			case EOL:
				statement := NewStatement()
				for _, context := range contexts {
					statement = append(statement, context...)
				}
				if len(statement) == 3 {
					statements = append(statements, statement)
					log.Printf("added statement %+v", statement)
					contexts = contexts[:len(contexts)-1]
				}
				indentCount = 0

			default:
				return nil, fmt.Errorf("unknown token '%v'", token)
			}
		}
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
	return statements, nil
}

type TokenType int

const (
	WORD TokenType = iota
	INDENT
	EOL
)

func getTokenType(word string) TokenType {
	if word == "    " || word == `\t` {
		return INDENT
	} else {
		return WORD
	}
}

type Token struct {
	TokenType TokenType
	Word      string
}

func getToken(word string) Token {
	return Token{
		TokenType: getTokenType(word),
		Word:      word,
	}
}

func getTokens(words []string) []Token {
	res := make([]Token, 0)
	for _, word := range words {
		res = append(res, getToken(word))
	}
	return res
}
