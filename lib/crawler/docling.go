package crawler

import (
	"os"
	"os/exec"
	"path/filepath"

	"github.com/google/uuid"
)

func PdfToMarkdown(tempDir string, pdf []byte) (string, error) {
	inFileName := uuid.New().String()
	inFile := filepath.Join(tempDir, inFileName+".pdf")
	defer os.Remove(inFile)
	outDir := filepath.Join(tempDir, uuid.New().String())
	outFile := filepath.Join(outDir, inFileName+".md")
	defer os.RemoveAll(outDir)
	f, err := os.Create(inFile)
	if err != nil {
		return "", err
	}
	_, err = f.Write(pdf)
	if err != nil {
		return "", err
	}
	f.Close()

	cmd := exec.Command("docling", "--from", "pdf", "--to", "md", "--image-export-mode", "placeholder", "--output", outDir, inFile)
	_, err = cmd.Output()
	if err != nil {
		return "", err
	}

	markdown, err := os.ReadFile(outFile)
	return string(markdown), err
}
