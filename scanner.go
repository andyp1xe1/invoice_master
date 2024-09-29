package main

import (
	"github.com/fumiama/go-docx"
	"github.com/gen2brain/go-fitz"
	"github.com/otiai10/gosseract/v2"

	"fmt"
	"image/png"
	"os"
	"path/filepath"
	"strings"
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

func (s *Scanner) pdfTxt(path string) (string, error) {
	doc, err := fitz.New(path)
	if err != nil {
		return "", fmt.Errorf("failed to open PDF: %w", err)
	}
	defer doc.Close()

	var content string

	for n := 0; n < doc.NumPage(); n++ {
		pageText, err := doc.Text(n)
		if err != nil {
			return "", fmt.Errorf("failed to extract text from page %d: %w", n, err)
		}

		content += pageText
	}
	if content != "" {
		return content, nil
	}

	for n := 0; n < doc.NumPage(); n++ {

		img, err := doc.Image(n)
		if err != nil {
			return "", fmt.Errorf("could not render page %d to image: %v", n, err)
		}

		imgFileName := fmt.Sprintf("page_%d.png", n+1)
		imgFile, err := os.Create(imgFileName)
		if err != nil {
			return "", fmt.Errorf("could not create image file %s: %v", imgFileName, err)
		}
		defer imgFile.Close()

		if err := png.Encode(imgFile, img); err != nil {
			return "", fmt.Errorf("could not encode image to PNG: %v", err)
		}

		fmt.Printf("Saved %s\n", imgFileName)

		text, err := s.imgOcr(imgFileName)
		if err != nil {
			return "", fmt.Errorf("OCR failed for %s: %v", imgFileName, err)
		}

		fmt.Printf("Text from page %d:\n%s\n", n+1, text)
		content += text + "\n"
	}

	return content, nil
}

func (s *Scanner) docTxt(path string) (string, error) {
	readFile, err := os.Open(path)
	if err != nil {
		return "", err
	}
	defer readFile.Close()

	fileinfo, err := readFile.Stat()
	if err != nil {
		return "", err
	}
	size := fileinfo.Size()

	doc, err := docx.Parse(readFile, size)
	if err != nil {
		return "", err
	}

	var textOutput string
	for _, it := range doc.Document.Body.Items {
		switch it.(type) {
		case *docx.Paragraph:
			textOutput += fmt.Sprintf("%v", it) + "\n---\n" // Add marker for new paragraph
		case *docx.Table:
			textOutput += fmt.Sprintf("%v", it) + "\n"
		}
	}
	return textOutput, nil
}

func (s *Scanner) extSwitch(filePath string) (string, error) {
	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".pdf":
		return s.pdfSwitch(filePath)
	case ".docx":
		println("test")
		return s.docTxt(filePath)
	case ".png":
		fallthrough
	case ".jpg", ".jpeg":
		return s.imgOcr(filePath)
	default:
		return "", nil
	}
}

func (s *Scanner) pdfSwitch(path string) (string, error) {
	var err error
	var txt string
	if txt, err = s.pdfTxt(path); txt == "" {
		txt, err = s.imgOcr(path)
	}
	if err != nil {
		return "", err
	}
	return txt, err
}
