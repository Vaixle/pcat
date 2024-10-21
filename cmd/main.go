package main

import (
	"bufio"
	"context"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"os"
	"os/signal"
	"pcat/internal"
	"pcat/internal/app"
	"runtime"
	"syscall"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	eg, ctx := errgroup.WithContext(ctx)
	eg.SetLimit(runtime.NumCPU())

	args := os.Args
	var paths []string
	if len(os.Args) == 1 || args[1] == "-" {

		//skip '-' in args
		if len(os.Args) > 1 && os.Args[1] == "-" {
			os.Args = append(os.Args[:1], os.Args[2:]...)
		}
		paths = readPathsFromStdIn()
	} else {
		paths = args[1:]
	}

	flagParser := internal.NewPcatFlagParser()
	pCat := app.NewPcat(eg, flagParser, paths)
	if pCat == nil {
		return
	}
	pCat.Run(ctx)
}

func readPathsFromStdIn() []string {
	var paths []string
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Printf("Supply values for the following parameters:\nPath[0]: ")
	for i := 1; scanner.Scan(); i++ {
		text := scanner.Text()
		if text == "" {
			break
		}

		fmt.Printf("Path[%d]: ", i)
		paths = append(paths, text)
	}

	if len(paths) == 0 {
		log.Fatalln("path is empty")
	}

	return paths
}
