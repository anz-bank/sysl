package seqs

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

func OutputPlantuml(output, plantuml, umlInput string) {
	l := len(output)
	mode := output[l-3:]

	if plantuml == "" {
		plantuml = os.Getenv("SYSL_PLANTUML")
	}

	switch mode {
	case "png", "svg":
		plantuml = fmt.Sprintf("%s/%s/%s", plantuml, mode, DeflateAndEncode([]byte(umlInput)))
		out, _ := sendHttpRequest(plantuml)
		ioutil.WriteFile(output, out, os.ModePerm)

	case "uml":
		output := output[:l-3]
		output += "puml"
		ioutil.WriteFile(output, []byte(umlInput), os.ModePerm)

	default:
		log.Errorf("Extension %s not supported. Valid extensions: svg, png, uml.", mode)
	}
}

func sendHttpRequest(url string) ([]byte, error) {
	resp, err := http.Get(url)
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
func DeflateAndEncode(text []byte) string {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		errors.Errorf("Unable to encode []byte, Error:%s", err.Error())
	}
	zw.Write(text)
	zw.Close()
	return encode(buf.Bytes())
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
func encode3bytes(buf *bytes.Buffer, b1, b2, b3 byte) {
	buf.WriteByte(encode6bit(0x3F & (b1 >> 2)))
	buf.WriteByte(encode6bit(0x3F & (((b1 & 0x3) << 4) | (b2 >> 4))))
	buf.WriteByte(encode6bit(0x3F & (((b2 & 0xF) << 2) | (b3 >> 6))))
	buf.WriteByte(encode6bit(0x3F & b3))
}

func encode6bit(b byte) byte {
	// 6 bit makes value 0 to 63. The func maps 0-63 to characters
	// '0'-'9', 'A'-'Z', 'a'-'z', '-', '_'. '?' should never be reached.
	if b > 63 {
		panic("unexpected character!")
	}
	return "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz-_"[b]
}
