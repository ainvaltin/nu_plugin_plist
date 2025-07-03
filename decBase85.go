package main

import (
	"bytes"
	"context"
	"encoding/ascii85"
	"fmt"
	"io"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/types"
)

func decodeBase85() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:     "decode base85",
			Category: "Formats",
			Desc:     `Decode Ascii85, also called Base85.`,
			Description: `Implements the ascii85 data encoding as used in the btoa tool and Adobe's PostScript and PDF document formats.` +
				"\n\nIf the decoded data is binary (not printable) then add 'into binary' to the pipeline ie '... decode base85 | into binary'",
			SearchTerms:      []string{"Ascii85", "Base85"},
			InputOutputTypes: []nu.InOutTypes{{In: types.String(), Out: types.Binary()}},
			/*Named: []nu.Flag{
				{Long: "trim-markers", Short: "m", Desc: "trim <~ and ~> symbols from the beginning and end of the data (if present)"},
			},*/
			AllowMissingExamples: true,
		},
		Examples: []nu.Example{
			{Description: `Decode base85 data`, Example: `'F)Po,+Cno&@/' | decode base85`, Result: &nu.Value{Value: `some data`}},
		},
		OnRun: decodeBase85Handler,
	}
}

func decodeBase85Handler(ctx context.Context, call *nu.ExecCommand) error {
	out, err := call.ReturnRawStream(ctx)
	if err != nil {
		return fmt.Errorf("creating response stream: %w", err)
	}
	defer out.Close()

	switch in := call.Input.(type) {
	case nil:
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
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}
