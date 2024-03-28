package codecs

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
)

type SupportedImageFormats string

const (
	imageFormatPng  SupportedImageFormats = "image/png"
	imageFormatJpeg SupportedImageFormats = "image/jpeg"
)

type Encoder func(w io.Writer, m image.Image) error
type Decoder func(r io.Reader) (image.Image, error)

var (
	encoders = map[SupportedImageFormats]Encoder{
		imageFormatPng: png.Encode,
		imageFormatJpeg: func(w io.Writer, m image.Image) error {
			return jpeg.Encode(w, m, &jpeg.Options{Quality: 20})
		},
	}
	decoders = map[SupportedImageFormats]Decoder{
		imageFormatPng:  png.Decode,
		imageFormatJpeg: jpeg.Decode,
	}
)

func DecodeImage(r io.Reader, format SupportedImageFormats) (image.Image, error) {
	decoder := decoders[format]
	return decoder(r)
}

func EncodeImage(img image.Image, format SupportedImageFormats) (blob []byte, err error) {
	encoder := encoders[format]

	var buf bytes.Buffer
	err = encoder(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func AssertSupportedFormat(str string) (SupportedImageFormats, error) {
	format := SupportedImageFormats(str)
	switch format {
	case imageFormatPng, imageFormatJpeg:
		return format, nil
	}
	return "", fmt.Errorf("unsupported format: %s", str)
}
