package internal

import (
	"bufio"
	"context"
	"os"
	"strings"
)

type FileReader interface {
	ReadFile(context.Context, string, int) error
}

type FileByLineReader struct {
	flagHandler FlagHandler
	textBuffer  []string
}

func NewFileByLineReader(flagHandler FlagHandler, textBuffer []string) FileReader {
	return &FileByLineReader{
		flagHandler: flagHandler,
		textBuffer:  textBuffer,
	}
}

func (r *FileByLineReader) ReadFile(ctx context.Context, path string, i int) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	strBuilder := strings.Builder{}

	for scanner.Scan() {
		text := scanner.Text()

		mText := &ModifiableText{
			text: text,
		}
		r.flagHandler.Execute(mText)

		if mText.skip {
			continue
		}

		strBuilder.WriteString(mText.text + "\n")
	}
	r.textBuffer[i] = strBuilder.String()
	return nil
}

type FileAllReader struct {
	textBuffer []string
}

func NewFileAllReader(textBuffer []string) FileReader {
	return &FileAllReader{
		textBuffer: textBuffer,
	}
}

func (r *FileAllReader) ReadFile(ctx context.Context, path string, i int) error {
	select {
	case <-ctx.Done():
		return nil
	default:
	}

	text, err := os.ReadFile(path)
	if err != nil {
		return err
	}

	r.textBuffer[i] = string(text)
	return nil
}
