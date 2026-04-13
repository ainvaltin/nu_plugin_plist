package main

import (
	"context"
	"fmt"
	"io"
	"math/rand/v2"
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
				{In: types.OneOf(types.Binary(), types.String()), Out: types.String()},
			},
			Named: []nu.Flag{
				base58AlphabetFlag(),
			},
			AllowMissingExamples: true,
		},
		Examples: []nu.Example{
			{Description: `Encode data as base58`, Example: `'some data' | encode base58 -a flickr`, Result: &nu.Value{Value: `2tdRS31QBvroH`}},
		},
		OnRun: encodeBase58Handler,
	}
}

func encodeBase58Handler(ctx context.Context, call *nu.ExecCommand) error {
	alphabet, err := base58alphabet(call)
	if err != nil {
		return err
	}

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

func base58alphabet(call *nu.ExecCommand) (*base58.Alphabet, error) {
	flag, _ := call.FlagValue("alphabet")
	switch s := strings.ToLower(flag.Value.(string)); s {
	case "btc", "":
		return base58.BTCAlphabet, nil
	case "flickr":
		return base58.FlickrAlphabet, nil
	case "xrp":
		return XRPLedgerAlphabet, nil
	default:
		if n := len(s); n != 58 {
			return nil, fmt.Errorf("alphabet must have 58 characters, got %d", n)
		}
		return base58.NewAlphabet(s), nil
	}
}

var XRPLedgerAlphabet = base58.NewAlphabet("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")

func base58AlphabetFlag() nu.Flag {
	return nu.Flag{
		Long:    "alphabet",
		Short:   'a',
		Shape:   syntaxshape.String(),
		Default: &nu.Value{Value: "btc"},
		Desc:    "Alphabet to use, must be 58 characters long. There is three shorthand values:\n\t - flickr: use the Flickr base58 alphabet;\n\t - XRP: use the XRP Ledger alphabet;\n\t - BTC: use the bitcoin base58 alphabet;",
		Completions: nu.DynamicCompletion(func() []nu.DynamicSuggestion {
			ra := []byte("rpshnaf39wBUDNEGHJKLM4PQRST7VWXYZ2bcdeCg65jkm8oFqi1tuvAxyz")
			rand.Shuffle(len(ra), func(i, j int) { ra[i], ra[j] = ra[j], ra[i] })
			return []nu.DynamicSuggestion{
				{Value: "BTC", Display: "the bitcoin alphabet"},
				{Value: "flickr", Display: "the Flickr alphabet"},
				{Value: "XRP", Display: "the XPR Ledger alphabet"},
				{Value: string(ra), Display: "random alphabet"},
			}
		}),
	}
}
