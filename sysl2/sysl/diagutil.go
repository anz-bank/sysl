package main

import (
	"bytes"
	"compress/zlib"
	"fmt"
	"image/jpeg"
	"image/png"
	"net/http"
	"os"

	"github.com/sirupsen/logrus"
)

func OutputPlantuml(mode, url string, umlInput []byte) {
	if url == "" {
		url = os.Getenv("SYSL_PLANTUML")
	}
	url = fmt.Sprintf("%s/%s/%s", url, mode, DeflateAndEncode(umlInput))
	resp, err := http.Get(url)
	if err != nil {
		logrus.Errorf("Unable to create http request to %s, Error:", url, err.Error())
	}
	defer resp.Body.Close()

	img, err := png.Decode(resp.Body)
	if err != nil {
		logrus.Errorf("Unable to create image from response, Error:", err.Error())
	}

	f, err := os.Create("_.jpg")
	if err != nil {
		logrus.Errorf("Unable to create file, Error:", err.Error())
	}
	defer f.Close()
	jpeg.Encode(f, img, nil)
}

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
	if b < 10 {
		return byte(48 + b)
	}
	b -= 10
	if b < 26 {
		return byte(65 + b)
	}
	b -= 26
	if b < 26 {
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
