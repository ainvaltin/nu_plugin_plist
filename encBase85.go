package main

import (
	"bytes"
	"context"
	"encoding/ascii85"
	"fmt"
	"io"

	"github.com/ainvaltin/nu-plugin"
)

func encodeBase85() *nu.Command {
	return &nu.Command{
		Sig: nu.PluginSignature{
			Name:             "encode base85",
			Category:         "Hash",
			Usage:            `Encode binary as Ascii85, also called Base85.`,
			UsageEx:          `Implements the ascii85 data encoding as used in the btoa tool and Adobe's PostScript and PDF document formats.`,
			SearchTerms:      []string{"Ascii85", "Base85"},
			InputOutputTypes: [][]string{{"Binary", "String"}, {"String", "String"}},
			/*Named: []nu.Flag{
				{Long: "add-markers", Short: "m", Desc: "wrap output in <~ and ~> symbols"},
			},*/
			AllowMissingExamples: true,
		},
		Exm: nu.Examples{
			{Description: `Convert binary to base85`, Example: `0x[0102030405] | encode base85`, Result: &nu.Value{Value: `!<N?+"T`}},
			{Description: `Convert an string to base85`, Example: `'Some Data' | encode base85`, Result: &nu.Value{Value: `;f?Ma+@KX[@/`}},
		},
		OnRun: encodeBase85Handler,
	}
}

func encodeBase85Handler(ctx context.Context, exec *nu.ExecCommand) error {
	out, err := exec.ReturnRawStream(ctx)
	if err != nil {
		return fmt.Errorf("creating response stream: %w", err)
	}
	defer out.Close()

	enc := ascii85.NewEncoder(out)
	defer enc.Close()

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
		_, err = io.Copy(enc, bytes.NewBuffer(buf))
		return err
	case io.ReadCloser:
		_, err = io.Copy(enc, in)
		return err
	default:
		return fmt.Errorf("unsupported input type %T", exec.Input)
	}
}
