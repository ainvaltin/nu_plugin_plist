package main

import (
	"bytes"
	"context"
	"fmt"
	"io"

	"github.com/mr-tron/base58"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

func decodeBase58() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:             "decode base58",
			Category:         "Formats",
			Desc:             `Decode base58 encoded data.`,
			SearchTerms:      []string{"base58"},
			InputOutputTypes: []nu.InOutTypes{{In: types.String(), Out: types.Binary()}},
			Named: []nu.Flag{{
				Long:    "alphabet",
				Short:   'a',
				Shape:   syntaxshape.String(),
				Default: &nu.Value{Value: "btc"},
				Desc:    "Alphabet to use, must be 58 characters long. There is two shorthand values:\n\t - BTC: use the bitcoin base58 alphabet;\n\t - flickr: use the Flickr base58 alphabet;",
			}},
			AllowMissingExamples: true,
		},
		Examples: []nu.Example{
			{Description: `Decode base58 data`, Example: `'2tdRS31QBvroH' | decode base58 -a flickr`, Result: &nu.Value{Value: `some data`}},
		},
		OnRun: decodeBase58Handler,
	}
}

func decodeBase58Handler(ctx context.Context, call *nu.ExecCommand) error {
	alphabet := base58alphabet(call)

	var buf string
	switch in := call.Input.(type) {
	case nil:
		return nil
	case nu.Value:
		switch data := in.Value.(type) {
		case []byte:
			buf = string(data)
		case string:
			buf = data
		default:
			return fmt.Errorf("unsupported Value type %T", data)
		}
	case io.ReadCloser:
		out := &bytes.Buffer{}
		if _, err := io.Copy(out, in); err != nil {
			return err
		}
		buf = out.String()
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}

	bin, err := base58.DecodeAlphabet(buf, alphabet)
	if err != nil {
		return fmt.Errorf("decoding input as base58: %w", err)
	}
	return call.ReturnValue(ctx, nu.Value{Value: bin})
}
