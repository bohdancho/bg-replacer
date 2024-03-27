package codecs

import (
	"bytes"
	"encoding/base64"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type ImageFormat string

const (
	imageFormatPng  ImageFormat = "image/png"
	imageFormatJpeg ImageFormat = "image/jpeg"
)

var ErrUnsupportedImageFormat = errors.New("unsupported image format")

type Encoder func(w io.Writer, m image.Image) error
type Decoder func(r io.Reader) (image.Image, error)

var (
	encoders = map[ImageFormat]Encoder{
		imageFormatPng: png.Encode,
		imageFormatJpeg: func(w io.Writer, m image.Image) error {
			return jpeg.Encode(w, m, &jpeg.Options{Quality: 20})
		},
	}
	decoders = map[ImageFormat]Decoder{
		imageFormatPng:  png.Decode,
		imageFormatJpeg: jpeg.Decode,
	}
)

func DecodeImage(s string) (image.Image, ImageFormat, error) {
	header, content, found := strings.Cut(s, ";base64,")
	if !found {
		return nil, "", errors.New("no ';base64,' separator found")
	}
	b, err := base64.StdEncoding.DecodeString(content)
	if err != nil {
		return nil, "", err
	}

	format := ImageFormat(strings.TrimPrefix(header, "data:"))

	decoder := decoders[format]
	if decoder == nil {
		return nil, "", ErrUnsupportedImageFormat
	}

	var decodedImage image.Image
	decodedImage, err = decoder(bytes.NewReader(b))
	if err != nil {
		return nil, "", err
	}

	return decodedImage, format, nil
}

func EncodeImage(img image.Image, format ImageFormat) (imgSrc string, err error) {
	encoder := encoders[format]
	if encoder == nil {
		return "", ErrUnsupportedImageFormat
	}

	var buf bytes.Buffer
	err = encoder(&buf, img)
	if err != nil {
		return "", err
	}

	image64 := base64.StdEncoding.EncodeToString(buf.Bytes())
	header := fmt.Sprintf("data:image/%s;base64,", format)

	responseBytes := append([]byte(header), image64...)
	return string(responseBytes), nil
}
