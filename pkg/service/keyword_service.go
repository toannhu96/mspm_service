package service

import (
	"bufio"
	"context"
	"github.com/toannhu96/mspm_service/pkg/ds/search"
	"os"
	"strings"
)

type keywordService struct {
	model *search.M
}

func NewKeywordService() (KeywordService, error) {
	model := search.NewModel("")
	file, err := os.Open("/resources/patterns.txt") // change this pattern for local development
	if err != nil {
		return nil, err
	}
	words := bufio.NewReader(file)
	model.Build(words)
	return &keywordService{model: model}, nil
}

type KeywordService interface {
	MultiTermMatch(ctx context.Context, input string) (output search.Output, err error)
}

func (svc *keywordService) MultiTermMatch(ctx context.Context, input string) (output search.Output, err error) {
	document := strings.NewReader(input)
	return svc.model.MultiTermMatch(document)
}
