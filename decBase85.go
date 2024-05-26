package main

import (
	"bytes"
	"context"
	"encoding/ascii85"
	"fmt"
	"io"

	"github.com/ainvaltin/nu-plugin"
)

func decodeBase85() *nu.Command {
	return &nu.Command{
		Sig: nu.PluginSignature{
			Name:     "decode base85",
			Category: "Hash",
			Usage:    `Decode Ascii85, also called Base85.`,
			UsageEx: `Implements the ascii85 data encoding as used in the btoa tool and Adobe's PostScript and PDF document formats.` +
				"\n\nIf the decoded data is binary (not printable) then add 'into binary' to the pipeline ie '... decode base85 | into binary'",
			SearchTerms:      []string{"Ascii85", "Base85"},
			InputOutputTypes: [][]string{{"String", "Binary"}},
			/*Named: []nu.Flag{
				{Long: "trim-markers", Short: "m", Desc: "trim <~ and ~> symbols from the beginning and end of the data (if present)"},
			},*/
			AllowMissingExamples: true,
		},
		Exm: nu.Examples{
			{Description: `Decode base85 data`, Example: `'F)Po,+Cno&@/' | decode base85`, Result: &nu.Value{Value: `some data`}},
		},
		OnRun: decodeBase85Handler,
	}
}

func decodeBase85Handler(ctx context.Context, exec *nu.ExecCommand) error {
	out, err := exec.ReturnRawStream(ctx)
	if err != nil {
		return fmt.Errorf("creating response stream: %w", err)
	}
	defer out.Close()

	switch in := exec.Input.(type) {
	case nu.Empty:
		return nil
	case nu.Value:
		var buf []byte
		switch data := in.Value.(type) {
		case []byte:
			buf = data
		case string:
			buf = []byte(data)
		default:
			return fmt.Errorf("unsupported Value type %T", data)
		}
		_, err = io.Copy(out, ascii85.NewDecoder(bytes.NewReader(buf)))
		return err
	case io.ReadCloser:
		_, err = io.Copy(out, ascii85.NewDecoder(in))
		return err
	default:
		return fmt.Errorf("unsupported input type %T", exec.Input)
	}
}
