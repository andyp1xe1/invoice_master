package main

import (
	"github.com/otiai10/gosseract/v2"
)

type Scanner struct {
	tess *gosseract.Client
}

func NewScanner() *Scanner {
	sc := &Scanner{}
	sc.NewTess()
	return sc
}

func (s *Scanner) NewTess() {
	s.tess = gosseract.NewClient()
	s.tess.SetLanguage("ron", "rus", "eng")
}

func (s *Scanner) imgOcr(path string) (string, error) {
	s.tess.SetImage(path)
	if text, err := s.tess.Text(); err != nil {
		return "", err
	} else {
		return text, nil
	}
}
