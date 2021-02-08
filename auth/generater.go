package auth

import "time"

type Generater interface {
	// Token creates a new token
	Token(string) (Token, error)
	// Code creates a new code
	Code() (string, error)
}

type generater struct{}

func (g *generater) Token(id string) (Token, error) {
	return Token{ID: id, Iat: time.Now().Unix(), Ein: 1}, nil
}

func (g *generater) Code() (string, error) {
	return "123456", nil
}
