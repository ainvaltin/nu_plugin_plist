package main

import (
	"context"
	"fmt"
	"strings"

	"howett.net/plist"

	"github.com/ainvaltin/nu-plugin"
	"github.com/ainvaltin/nu-plugin/syntaxshape"
	"github.com/ainvaltin/nu-plugin/types"
)

func toPlist() *nu.Command {
	cmd := &nu.Command{
		Signature: nu.PluginSignature{
			Name:             "to plist",
			Category:         "Formats",
			Desc:             `Convert Nushell Value to property list.`,
			SearchTerms:      []string{"plist", "GNU Step", "Open Step", "xml"},
			InputOutputTypes: []nu.InOutTypes{{In: types.Any(), Out: types.Binary()}, {In: types.Any(), Out: types.String()}},
			Named: []nu.Flag{
				{Long: "format", Short: "f", Shape: syntaxshape.String(), Default: &nu.Value{Value: "bin"}, Desc: "Which plist format to use: xml, gnu[step], open[step]. Any other value will mean that binary format will be used."},
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
	outFmt := plistFormat(call)
	prettyFmt := prettyFormat(call)

	switch in := call.Input.(type) {
	case nil:
		return nil
	case nu.Value:
		v, err := toPlistValue(in, outFmt, prettyFmt)
		if err != nil {
			return err
		}
		return call.ReturnValue(ctx, v)
	case <-chan nu.Value:
		out, err := call.ReturnListStream(ctx)
		if err != nil {
			return err
		}
		defer close(out)
		for v := range in {
			v, err := toPlistValue(v, outFmt, prettyFmt)
			if err != nil {
				return err
			}
			out <- v
		}
		return nil
	default:
		return fmt.Errorf("unsupported input type %T", call.Input)
	}
}

func toPlistValue(v nu.Value, outFmt int, prettyFmt bool) (nu.Value, error) {
	var buf []byte
	var err error
	if prettyFmt && outFmt != plist.BinaryFormat {
		buf, err = plist.MarshalIndent(fromValue(v), outFmt, "\t")
	} else {
		buf, err = plist.Marshal(fromValue(v), outFmt)
	}
	if err != nil {
		return nu.Value{}, fmt.Errorf("encoding %T as plist: %w", v.Value, err)
	}
	if outFmt == plist.BinaryFormat {
		return nu.Value{Value: buf}, nil
	}
	return nu.Value{Value: string(buf)}, nil
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

func plistFormat(call *nu.ExecCommand) int {
	v, _ := call.FlagValue("format")
	switch strings.ToLower(v.Value.(string)) {
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

func prettyFormat(call *nu.ExecCommand) bool {
	v, _ := call.FlagValue("pretty")
	return v.Value.(bool)
}
