package main

import (
	"context"
	"fmt"
	"io"

	"howett.net/plist"

	"github.com/ainvaltin/nu-plugin"
)

func fromPlist() *nu.Command {
	return &nu.Command{
		Signature: nu.PluginSignature{
			Name:                 "from plist",
			Category:             "Formats",
			SearchTerms:          []string{"plist", "GNU Step", "Open Step", "xml"},
			InputOutputTypes:     [][]string{{"Binary", "Any"}, {"String", "Any"}},
			Usage:                `Convert from 'property list' format to Nushell Value.`,
			AllowMissingExamples: true,
		},
		Examples: nu.Examples{
			{
				Description: `Convert an Open Step array to list of Nu values`,
				Example:     `'(10,foo,)' | from plist`,
				Result:      &nu.Value{Value: []nu.Value{{Value: 10}, {Value: "foo"}}},
			},
		},
		OnRun: fromPlistHandler,
	}
}

func fromPlistHandler(ctx context.Context, call *nu.ExecCommand) error {
	switch in := call.Input.(type) {
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
			return fmt.Errorf("unsupported input value type %T", data)
		}
		var v any
		if _, err := plist.Unmarshal(buf, &v); err != nil {
			return fmt.Errorf("decoding input as plist: %w", err)
		}
		rv, err := asValue(v)
		if err != nil {
			return fmt.Errorf("converting to Value: %w", err)
		}
		return call.ReturnValue(ctx, rv)
	case io.Reader:
		// decoder wants io.ReadSeeker so we need to read to buf.
		// could read just enough that the decoder can detect the
		// format and stream the rest?
		buf, err := io.ReadAll(in)
		if err != nil {
			return fmt.Errorf("reding input: %w", err)
		}
		var v any
		if _, err := plist.Unmarshal(buf, &v); err != nil {
			return fmt.Errorf("decoding input as plist: %w", err)
		}
		rv, err := asValue(v)
		if err != nil {
			return fmt.Errorf("converting to Value: %w", err)
		}
		return call.ReturnValue(ctx, rv)
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}

func asValue(v any) (_ nu.Value, err error) {
	switch in := v.(type) {
	case uint64, float64, bool, string, []byte:
		return nu.Value{Value: in}, nil
	case []any:
		lst := make([]nu.Value, len(in))
		for i := 0; i < len(in); i++ {
			if lst[i], err = asValue(in[i]); err != nil {
				return nu.Value{}, err
			}
		}
		return nu.Value{Value: lst}, nil
	case map[string]any:
		rec := nu.Record{}
		for k, v := range in {
			if rec[k], err = asValue(v); err != nil {
				return nu.Value{}, err
			}
		}
		return nu.Value{Value: rec}, nil
	default:
		return nu.Value{}, fmt.Errorf("unsupported value type %T", in)
	}
}
