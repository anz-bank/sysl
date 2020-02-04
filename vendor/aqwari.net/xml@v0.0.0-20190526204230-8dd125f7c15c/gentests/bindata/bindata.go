package bindata

import (
	"bytes"
	"encoding/base64"
	"encoding/hex"
	"encoding/xml"
)

type Bindata struct {
	HexData  []byte `xml:"tns hexData"`
	B64Data  []byte `xml:"tns b64Data"`
	Filename string `xml:"tns filename"`
}

func (t *Bindata) MarshalXML(e *xml.Encoder, start xml.StartElement) error {
	type T Bindata
	var layout struct {
		*T
		HexData *xsdHexBinary    `xml:"tns hexData"`
		B64Data *xsdBase64Binary `xml:"tns b64Data"`
	}
	layout.T = (*T)(t)
	layout.HexData = (*xsdHexBinary)(&layout.T.HexData)
	layout.B64Data = (*xsdBase64Binary)(&layout.T.B64Data)
	return e.EncodeElement(layout, start)
}
func (t *Bindata) UnmarshalXML(d *xml.Decoder, start xml.StartElement) error {
	type T Bindata
	var overlay struct {
		*T
		HexData *xsdHexBinary    `xml:"tns hexData"`
		B64Data *xsdBase64Binary `xml:"tns b64Data"`
	}
	overlay.T = (*T)(t)
	overlay.HexData = (*xsdHexBinary)(&overlay.T.HexData)
	overlay.B64Data = (*xsdBase64Binary)(&overlay.T.B64Data)
	return d.DecodeElement(&overlay, &start)
}

type xsdBase64Binary []byte

func (b *xsdBase64Binary) UnmarshalText(text []byte) (err error) {
	*b, err = base64.StdEncoding.DecodeString(string(text))
	return
}
func (b xsdBase64Binary) MarshalText() ([]byte, error) {
	var buf bytes.Buffer
	enc := base64.NewEncoder(base64.StdEncoding, &buf)
	enc.Write([]byte(b))
	enc.Close()
	return buf.Bytes(), nil
}

type xsdHexBinary []byte

func (b *xsdHexBinary) UnmarshalText(text []byte) (err error) {
	*b, err = hex.DecodeString(string(text))
	return
}
func (b xsdHexBinary) MarshalText() ([]byte, error) {
	n := hex.EncodedLen(len(b))
	buf := make([]byte, n)
	hex.Encode(buf, []byte(b))
	return buf, nil
}
