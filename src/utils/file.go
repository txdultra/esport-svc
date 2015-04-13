package utils

import (
	"bufio"
	"bytes"
	"code.google.com/p/go.image/bmp"
	"code.google.com/p/go.image/tiff"
	"errors"
	"fmt"
	"image"
	_ "image/color"
	_ "image/gif"
	"image/jpeg"
	"image/png"
	"io"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"strings"
)

func IsImage(extName string) bool {
	switch extName {
	case "gif", "png", "jpg", "jpeg", "bmp", "tif", "tiff":
		return true
	default:
		return false
	}
}

func ByteToImage(data []byte) *image.Image {
	if len(data) == 0 {
		return nil
	}
	buf := bytes.NewBuffer(data)
	img, _, err := image.Decode(buf)
	if err != nil {
		return nil
	}
	return &img
}

func ImageToBytes(img image.Image, filename string) (data []byte, err error) {
	format := strings.ToLower(filepath.Ext(filename))
	okay := false
	for _, ext := range []string{".jpg", ".jpeg", ".png", ".tif", ".tiff", ".bmp"} {
		if format == ext {
			okay = true
			break
		}
	}
	if okay {
		buf := &bytes.Buffer{}
		switch format {
		case ".jpg", ".jpeg":
			var rgba *image.RGBA
			if nrgba, ok := img.(*image.NRGBA); ok {
				if nrgba.Opaque() {
					rgba = &image.RGBA{
						Pix:    nrgba.Pix,
						Stride: nrgba.Stride,
						Rect:   nrgba.Rect,
					}
				}
			}
			if rgba != nil {
				err = jpeg.Encode(buf, rgba, &jpeg.Options{Quality: 95})
			} else {
				err = jpeg.Encode(buf, img, &jpeg.Options{Quality: 95})
			}

		case ".png":
			err = png.Encode(buf, img)
		case ".tif", ".tiff":
			err = tiff.Encode(buf, img, &tiff.Options{Compression: tiff.Deflate, Predictor: true})
		case ".bmp":
			err = bmp.Encode(buf, img)
		}
		if err == nil {
			data = buf.Bytes()
		}
		return
	}
	return nil, fmt.Errorf(`unsupported image format: "%s"`, format)
}

func FileName(file string) string {
	if len(file) == 0 {
		return ""
	}
	idx := strings.Index(file, ".")
	if idx <= 0 {
		return file
	}
	return file[0:idx]
}

func FileExtName(file string) string {
	ext := strings.ToLower(path.Ext(file))
	if len(ext) > 0 {
		return ext[1:]
	}
	return ""
}

// SelfPath gets compiled executable file absolute path
func SelfPath() string {
	path, _ := filepath.Abs(os.Args[0])
	return path
}

// SelfDir gets compiled executable file directory
func SelfDir() string {
	return filepath.Dir(SelfPath())
}

// FileExists reports whether the named file or directory exists.
func FileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}

// Search a file in paths.
// this is often used in search config file in /etc ~/
func SearchFile(filename string, paths ...string) (fullpath string, err error) {
	for _, path := range paths {
		if fullpath = filepath.Join(path, filename); FileExists(fullpath) {
			return
		}
	}
	err = errors.New(fullpath + " not found in paths")
	return
}

// like command grep -E
// for example: GrepFile(`^hello`, "hello.txt")
// \n is striped while read
func GrepFile(patten string, filename string) (lines []string, err error) {
	re, err := regexp.Compile(patten)
	if err != nil {
		return
	}

	fd, err := os.Open(filename)
	if err != nil {
		return
	}
	lines = make([]string, 0)
	reader := bufio.NewReader(fd)
	prefix := ""
	isLongLine := false
	for {
		byteLine, isPrefix, er := reader.ReadLine()
		if er != nil && er != io.EOF {
			return nil, er
		}
		if er == io.EOF {
			break
		}
		line := string(byteLine)
		if isPrefix {
			prefix += line
			continue
		} else {
			isLongLine = true
		}

		line = prefix + line
		if isLongLine {
			prefix = ""
		}
		if re.MatchString(line) {
			lines = append(lines, line)
		}
	}
	return lines, nil
}
