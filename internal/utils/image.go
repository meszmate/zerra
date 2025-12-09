package utils

import (
	"bytes"
	"image"
	"io"
	"mime/multipart"
	"path/filepath"
	"strings"

	"github.com/disintegration/imaging"
	"github.com/meszmate/zerra/internal/errx"
)

func IsValidJPEG(file multipart.File) *errx.Error {
	_, format, err := image.DecodeConfig(file)
	if err != nil {
		return errx.ErrImageFormat
	}
	if format != "jpeg" {
		return errx.ErrImageFormat
	}
	return nil
}

func resizeTo(img image.Image, w, h int) *image.NRGBA {
	return imaging.Fill(img, w, h, imaging.Center, imaging.Lanczos)
}

func GetJPG(f *multipart.FileHeader, maxSize, w, h int) (io.Reader, *errx.Error) {
	ext := strings.ToLower(filepath.Ext(f.Header.Get("name")))
	if ext != ".jpg" && ext != ".jpeg" {
		return nil, errx.ErrImageFormat
	}

	file, err := f.Open()
	if err != nil {
		return nil, errx.ErrImage
	}
	defer file.Close()

	if err := IsValidJPEG(file); err != nil {
		return nil, err
	}

	lreader := io.LimitReader(file, int64(maxSize))

	img, err := imaging.Decode(lreader)
	if err != nil {
		return nil, errx.ErrImage
	}

	img = resizeTo(img, w, h)
	buf := new(bytes.Buffer)
	if err := imaging.Encode(buf, img, imaging.JPEG, imaging.JPEGQuality(90)); err != nil {
		return nil, errx.ErrImage
	}

	return bytes.NewReader(buf.Bytes()), nil
}
