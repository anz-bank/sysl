package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func OutputPlantuml(output, plantuml, umlInput string) {
	mode := output[len(output)-3:]

	if plantuml == "" {
		plantuml = os.Getenv("SYSL_PLANTUML")
	}

	switch mode {
	case "png", "svg":
		plantuml = fmt.Sprintf("%s/%s/%s", plantuml, mode, DeflateAndEncode([]byte(umlInput)))
		resp, err := http.Get(plantuml)
		if err != nil {
			logrus.Errorf("Unable to create http request to %s, Error:", plantuml, err.Error())
		}
		defer resp.Body.Close()

		out, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logrus.Errorf("Unable to read from response, Error:", err.Error())
		}
		createFile(output, out)

	case "uml":
		createFile(output, []byte(umlInput))

	default:
		logrus.Errorf("Extension %s not supported. Valid extensions: svg, png, uml.", mode)
	}
}

func createFile(output string, out []byte) {
	f, err := os.Create(output)
	if err != nil {
		logrus.Errorf("Unable to create file, Error:", err.Error())
	}
	_, err = f.Write(out)
	defer f.Close()

	if err != nil {
		logrus.Errorf("Unable to create file, Error:", err.Error())
	}
}

//
// The functions below ported from https://github.com/dougn/python-plantuml/blob/master/plantuml.py
//
func DeflateAndEncode(text []byte) string {
	var buf bytes.Buffer
	zw, err := zlib.NewWriterLevel(&buf, zlib.BestCompression)
	if err != nil {
		logrus.Errorf("Unable to encode []byte, Error:", err.Error())
	}
	zw.Write(text)
	zw.Flush()
	zw.Close()
	return encode(buf.Bytes())
}

func encode(data []byte) string {
	var buf bytes.Buffer
	for i := 0; i < len(data); i += 3 {
		if i+2 == len(data) {
			encode3bytes(&buf, data[i], data[i+1], 0)
		} else if i+1 == len(data) {
			encode3bytes(&buf, data[i], 0, 0)
		} else {
			encode3bytes(&buf, data[i], data[i+1], data[i+2])
		}
	}
	return buf.String()
}

// 3 bytes takes 24 bits. This splits 24 bits into 4 bytes of which lower 6-bit takes account.
func encode3bytes(buf *bytes.Buffer, b1, b2, b3 byte) {
	c1 := b1 >> 2
	c2 := ((b1 & 0x3) << 4) | (b2 >> 4)
	c3 := ((b2 & 0xF) << 2) | (b3 >> 6)
	c4 := b3 & 0x3F

	buf.WriteByte(encode6bit(c1 & 0x3F))
	buf.WriteByte(encode6bit(c2 & 0x3F))
	buf.WriteByte(encode6bit(c3 & 0x3F))
	buf.WriteByte(encode6bit(c4 & 0x3F))
}

func encode6bit(b byte) byte {
	// 6 bit makes value 0 to 63. The func maps 0-63 to characters
	// '0'-'9', 'A'-'Z', 'a'-'z', '-', '_'. '?' should never be reached.
	if b < 10 {
		// 48 -> '0'
		return byte(48 + b)
	}
	b -= 10

	if b < 26 {
		// 65 -> 'A'
		return byte(65 + b)
	}
	b -= 26

	if b < 26 {
		// 97 -> 'a'
		return byte(97 + b)
	}
	b -= 26

	if b == 0 {
		return '-'
	}
	if b == 1 {
		return '_'
	}
	return '?'

}
