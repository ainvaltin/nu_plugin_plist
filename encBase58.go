package main

import (
	"context"
	"fmt"
	"io"
	"strings"

	"github.com/mr-tron/base58"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

func encodeBase58() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:        "encode base58",
			Category:    "Formats",
			Desc:        `Encode data as base58.`,
			SearchTerms: []string{"base58"},
			InputOutputTypes: []nu.InOutTypes{
				{In: types.String(), Out: types.String()},
				{In: types.Binary(), Out: types.String()},
			},
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
			{Description: `Encode data as base58`, Example: `'some data' | encode base58 -a flickr`, Result: &nu.Value{Value: `2tdRS31QBvroH`}},
		},
		OnRun: encodeBase58Handler,
	}
}

func encodeBase58Handler(ctx context.Context, call *nu.ExecCommand) error {
	alphabet := base58alphabet(call)

	var buf []byte
	switch in := call.Input.(type) {
	case nil:
		return nil
	case nu.Value:
		switch data := in.Value.(type) {
		case []byte:
			buf = data
		case string:
			buf = []byte(data)
		default:
			return fmt.Errorf("unsupported Value type %T", data)
		}
	case io.ReadCloser:
		var err error
		if buf, err = io.ReadAll(in); err != nil {
			return fmt.Errorf("reading input stream: %w", err)
		}
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}

	return call.ReturnValue(ctx, nu.Value{Value: base58.FastBase58EncodingAlphabet(buf, alphabet)})
}

func base58alphabet(call *nu.ExecCommand) *base58.Alphabet {
	flag, _ := call.FlagValue("alphabet")
	switch s := strings.ToLower(flag.Value.(string)); s {
	case "btc", "":
		return base58.BTCAlphabet
	case "flickr":
		return base58.FlickrAlphabet
	default:
		return base58.NewAlphabet(s)
	}
}
