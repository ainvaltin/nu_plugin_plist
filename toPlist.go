package main

import (
	"context"
	"fmt"
	"strings"

	"howett.net/plist"

	"github.com/ainvaltin/nu-plugin"
)

func toPlist() *nu.Command {
	cmd := &nu.Command{
		Signature: nu.PluginSignature{
			Name:             "to plist",
			Category:         "Formats",
			Usage:            `Convert Nushell Value to property list.`,
			SearchTerms:      []string{"plist", "GNU Step", "Open Step", "xml"},
			InputOutputTypes: [][]string{{"Any", "Binary"}, {"Any", "String"}},
			Named: []nu.Flag{
				{Long: "format", Short: "f", Arg: "String", Default: &nu.Value{Value: "bin"}, Desc: "Which plist format to use: xml, gnu[step], open[step]. Any other value will mean that binary format will be used."},
				{Long: "pretty", Short: "p", Desc: "If this switch is set output is 'pretty printed'. Only makes sense with text based formats, ignored for binary."},
			},
			AllowMissingExamples: true,
		},
		Examples: nu.Examples{
			{Description: `Convert an record to GNU Step format`, Example: `{foo: 10} | to plist -f gnu`, Result: &nu.Value{Value: `{foo=<*I10>;}`}},
			{Description: `Convert an array to Open Step format`, Example: `[10 foo] | to plist -f open`, Result: &nu.Value{Value: `(10,foo,)`}},
		},
		OnRun: toPlistHandler,
	}
	return cmd
}

func toPlistHandler(ctx context.Context, call *nu.ExecCommand) error {
	switch in := call.Input.(type) {
	case nu.Empty:
		return nil
	case nu.Value:
		outFmt := plistFormat(call.Named)
		var buf []byte
		var err error
		if prettyFormat(call.Named) && outFmt != plist.BinaryFormat {
			buf, err = plist.MarshalIndent(fromValue(in), outFmt, "\t")
		} else {
			buf, err = plist.Marshal(fromValue(in), outFmt)
		}
		if err != nil {
			return fmt.Errorf("encoding %T as plist: %w", in.Value, err)
		}
		if outFmt == plist.BinaryFormat {
			return call.ReturnValue(ctx, nu.Value{Value: buf})
		}
		return call.ReturnValue(ctx, nu.Value{Value: string(buf)})
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}

func fromValue(v nu.Value) any {
	switch vt := v.Value.(type) {
	case []nu.Value:
		lst := make([]any, len(vt))
		for i := 0; i < len(vt); i++ {
			lst[i] = fromValue(vt[i])
		}
		return lst
	case nu.Record:
		rec := map[string]any{}
		for k, v := range vt {
			rec[k] = fromValue(v)
		}
		return rec
	}
	return v.Value
}

func plistFormat(flags nu.NamedParams) int {
	switch strings.ToLower(flags.StringValue("format", "binary")) {
	case "xml":
		return plist.XMLFormat
	case "gnu", "gnustep":
		return plist.GNUStepFormat
	case "open", "openstep":
		return plist.OpenStepFormat
	default:
		return plist.BinaryFormat
	}
}

func prettyFormat(flags nu.NamedParams) bool {
	_, ok := flags["pretty"]
	return ok
}
