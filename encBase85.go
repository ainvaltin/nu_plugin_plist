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

func encodeBase85() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:             "encode base85",
			Category:         "Formats",
			Desc:             `Encode binary as Ascii85, also called Base85.`,
			Description:      `Implements the ascii85 data encoding as used in the btoa tool and Adobe's PostScript and PDF document formats.`,
			SearchTerms:      []string{"Ascii85", "Base85"},
			InputOutputTypes: []nu.InOutTypes{{In: types.Binary(), Out: types.String()}, {In: types.String(), Out: types.String()}},
			/*Named: []nu.Flag{
				{Long: "add-markers", Short: "m", Desc: "wrap output in <~ and ~> symbols"},
			},*/
			AllowMissingExamples: true,
		},
		Examples: []nu.Example{
			{Description: `Convert binary to base85`, Example: `0x[0102030405] | encode base85`, Result: &nu.Value{Value: `!<N?+"T`}},
			{Description: `Convert an string to base85`, Example: `'Some Data' | encode base85`, Result: &nu.Value{Value: `;f?Ma+@KX[@/`}},
		},
		OnRun: encodeBase85Handler,
	}
}

func encodeBase85Handler(ctx context.Context, call *nu.ExecCommand) error {
	out, err := call.ReturnRawStream(ctx, nu.StringStream())
	if err != nil {
		return fmt.Errorf("creating response stream: %w", err)
	}
	defer out.Close()

	enc := ascii85.NewEncoder(out)
	defer enc.Close()

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
		_, err = io.Copy(enc, bytes.NewBuffer(buf))
		return err
	case io.ReadCloser:
		_, err = io.Copy(enc, in)
		return err
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}
