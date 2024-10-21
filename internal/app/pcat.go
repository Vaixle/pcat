package app

import (
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"pcat/internal"
	"strings"
)

type Application interface {
	Run(context.Context)
}

type Pcat struct {
	paths      []string
	textBuffer []string
	fileReader internal.FileReader
	eg         *errgroup.Group
}

func NewPcat(eg *errgroup.Group, flagParser internal.FlagParser, paths []string) Application {
	flags, info := flagParser.ParseFlags()
	if info {
		return nil
	}

	var fileReader internal.FileReader
	textBuffer := make([]string, len(paths))

	if len(flags) != 0 {
		flagHandler := internal.NewPcatFlagHandler(flags...)
		fileReader = internal.NewFileByLineReader(flagHandler, textBuffer)
	} else {
		fileReader = internal.NewFileAllReader(textBuffer)
	}

	return &Pcat{
		paths:      paths,
		textBuffer: textBuffer,
		fileReader: fileReader,
		eg:         eg,
	}
}

func (p *Pcat) Run(ctx context.Context) {
	for i, path := range p.paths {
		p.eg.Go(
			func() error {
				return p.fileReader.ReadFile(ctx, path, i)
			})
	}
	if err := p.eg.Wait(); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println(strings.Join(p.textBuffer, "\n"))
	}
}
