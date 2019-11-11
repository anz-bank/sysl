package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	"github.com/nofun97/afero"
)

func OutputPlantuml(output, plantuml, umlInput string, fs afero.Fs) error {
	l := len(output)
	mode := output[l-3:]
	
	switch mode {
	case "png", "svg":
		encoded, err := DeflateAndEncode([]byte(umlInput))
		if err != nil {
			return err
		}
		plantuml = fmt.Sprintf("%s/%s/%s", plantuml, mode, encoded)
		out, err := sendHTTPRequest(plantuml)
		if err != nil {
			return err
		}
		return errors.Wrapf(afero.WriteFile(fs, output, out, os.ModePerm), "writing %q", output)

	case "uml":
		output := output[:l-3]
		output += "puml"
		return errors.Wrapf(afero.WriteFile(fs, output, []byte(umlInput), os.ModePerm), "writing %q", output)

	default:
		return fmt.Errorf("extension must be svg, png or uml, not %#v", mode)
	}
}

func sendHTTPRequest(url string) ([]byte, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return nil, errors.Errorf("Unable to create http request to %s, Error:%s", url, err.Error())
	}
	defer resp.Body.Close()

	out, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Errorf("Unable to read body.")
	}
	return out, nil
}

// The functions below ported from https://github.com/dougn/python-plantuml/blob/master/plantuml.py
func DeflateAndEncode(text []byte) (string, error) {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		return "", err
	}
	if _, err := zw.Write(text); err != nil {
		return "", err
	}
	if err := zw.Close(); err != nil {
		return "", err
	}
	return encode(buf.Bytes()), nil
}

func encode(data []byte) string {
	var buf bytes.Buffer
	i := 0
	for wholeTripleBytes := len(data) / 3 * 3; i < wholeTripleBytes; i += 3 {
		encode3bytes(&buf, data[i], data[i+1], data[i+2])
	}
	switch len(data) - i {
	case 1:
		encode3bytes(&buf, data[i], 0, 0)
	case 2:
		encode3bytes(&buf, data[i], data[i+1], 0)
	}
	return buf.String()
}

// 3 bytes takes 24 bits. This splits 24 bits into 4 bytes of which lower 6-bit takes account.
func encode3bytes(w io.ByteWriter, b1, b2, b3 byte) {
	if err := w.WriteByte(encode6bit(0x3F & (b1 >> 2))); err != nil {
		panic(err)
	}
	if err := w.WriteByte(encode6bit(0x3F & (((b1 & 0x3) << 4) | (b2 >> 4)))); err != nil {
		panic(err)
	}
	if err := w.WriteByte(encode6bit(0x3F & (((b2 & 0xF) << 2) | (b3 >> 6)))); err != nil {
		panic(err)
	}
	if err := w.WriteByte(encode6bit(0x3F & b3)); err != nil {
		panic(err)
	}
}

func encode6bit(b byte) byte {
	// 6 bit makes value 0 to 63. The func maps 0-63 to characters
	// '0'-'9', 'A'-'Z', 'a'-'z', '-', '_'. '?' should never be reached.
	if b > 63 {
		panic("unexpected character!")
	}
	return "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"[b]
}
