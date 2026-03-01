package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/ainvaltin/nu-plugin"
)

func main() {
	p, err := nu.New(
		[]*nu.Command{
			toPlist(),
			fromPlist(),
			encodeBase85(),
			decodeBase85(),
			encodeBase58(),
			decodeBase58(),
		},
		"0.0.1",
		debugCfg(),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "failed to create plugin", err)
		return
	}
	if err := p.Run(quitSignalContext()); err != nil && !errors.Is(err, nu.ErrGoodbye) {
		fmt.Fprintln(os.Stderr, "plugin exited with error", err)
	}
}

func quitSignalContext() context.Context {
	ctx, cancel := context.WithCancelCause(context.Background())

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)
		defer signal.Stop(sigChan)
		sig := <-sigChan
		cancel(fmt.Errorf("got quit signal: %s", sig))
	}()

	return ctx
}

func debugCfg() *nu.Config {
	fIn, err := os.Create("/Users/ain/go/src/nu_plugin_plist/input.log")
	if err != nil {
		panic(err)
	}
	fOut, err := os.Create("/Users/ain/go/src/nu_plugin_plist/output.log")
	if err != nil {
		panic(err)
	}
	fLog, err := os.Create("/Users/ain/go/src/nu_plugin_plist/log.txt")
	if err != nil {
		panic(err)
	}
	return &nu.Config{
		Logger:   slog.New(slog.NewTextHandler(fLog, &slog.HandlerOptions{Level: slog.LevelDebug})),
		SniffIn:  fIn,
		SniffOut: fOut,
	}
}
