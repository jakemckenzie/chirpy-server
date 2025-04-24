package services

import (
	"strings"
	"unicode"
)

type TextService struct{}

func NewTextService() *TextService {
	return &TextService{}
}

func (ts *TextService) CleanProfanity(text string) string {
	tokens := strings.Fields(text)
	for i, token := range tokens {
		if isAlphabetic(token) && profaneWords[strings.ToLower(token)] {
			tokens[i] = "****"
		}
	}
	return strings.Join(tokens, " ")
}

func isAlphabetic(s string) bool {
	for _, r := range s {
		if !unicode.IsLetter(r) {
			return false
		}
	}
	return true
}

var profaneWords = map[string]bool{
	"kerfuffle": true,
	"sharbert":  true,
	"fornax":    true,
}
