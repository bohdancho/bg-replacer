package codecs

import (
	"bytes"
	"fmt"
	"image"
	"image/jpeg"
	"image/png"
	"io"
	"strings"
)

type ImageContentType string

const (
	imageTypePng  ImageContentType = "image/png"
	imageTypeJpeg ImageContentType = "image/jpeg"
)

type Encoder func(w io.Writer, m image.Image) error
type Decoder func(r io.Reader) (image.Image, error)

func (contentType ImageContentType) Extension() string {
	s := string(contentType)
	return "." + strings.TrimPrefix(s, "image/")
}

var (
	encoders = map[ImageContentType]Encoder{
		imageTypePng: png.Encode,
		imageTypeJpeg: func(w io.Writer, m image.Image) error {
			return jpeg.Encode(w, m, &jpeg.Options{Quality: 20})
		},
	}
	decoders = map[ImageContentType]Decoder{
		imageTypePng:  png.Decode,
		imageTypeJpeg: jpeg.Decode,
	}
)

func DecodeImage(r io.Reader, contentType ImageContentType) (image.Image, error) {
	decoder := decoders[contentType]
	return decoder(r)
}

func EncodeImage(img image.Image, contentType ImageContentType) (blob []byte, err error) {
	encoder := encoders[contentType]

	var buf bytes.Buffer
	err = encoder(&buf, img)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func AssertSupportedImageType(str string) (ImageContentType, error) {
	contentType := ImageContentType(str)
	switch contentType {
	case imageTypePng, imageTypeJpeg:
		return contentType, nil
	}
	return "", fmt.Errorf("unsupported contentType: %s", str)
}
