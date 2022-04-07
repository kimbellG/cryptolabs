package main

import (
	"bytes"
	"flag"
	"fmt"
	"image/jpeg"
	"image/png"
	"io"
	"os"

	"github.com/kimbellG/cryptolabs/shorthand/lsb"
)

var (
	extract = flag.Bool("e", false, "extract mode")

	photo   = flag.String("photo", "", "jpeg photo")
	message = flag.String("m", "", "message to the photo")
	out     = flag.String("out", "", "output file")
)

func main() {
	flag.Parse()

	m()
}

func m() {
	switch {
	case *extract:
		if *photo == "" {
			fmt.Printf("%s: input photo file. -h for help\n", flag.Arg(0))

			os.Exit(-1)
		}

		msg, err := decode(*photo)
		if err != nil {
			fmt.Printf("%s: decode message from photo: %v\n", flag.Arg(0), err)

			os.Exit(-1)
		}

		fmt.Printf("message from photo: %s\n", msg)

	default:
		if *photo == "" || *message == "" || *out == "" {
			fmt.Printf("%s: message, photo or out file isn't correct\n", flag.Arg(0))

			os.Exit(-1)
		}

		if err := encode(*message, *photo, *out); err != nil {
			fmt.Printf("%s: encode message: %v\n", flag.Arg(0), err)

			os.Exit(-1)
		}
	}
}

func encode(message, photo, outfilename string) error {
	in, err := os.Open(photo)
	if err != nil {
		return fmt.Errorf("could not open file with photo: %w", err)
	}

	jp, err := jpeg.Decode(in)
	if err != nil {
		return fmt.Errorf("decode file like jpeg photo: %w", err)
	}

	buf := bytes.NewBuffer(nil)

	err = lsb.Encode(buf, jp, []byte(message))
	if err != nil {
		return fmt.Errorf("encode message: %w", err)
	}

	out, err := os.Create(outfilename)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	_, err = io.Copy(out, buf)
	if err != nil {
		return fmt.Errorf("write to the output file: %w", err)
	}

	return nil
}

func decode(photo string) (string, error) {
	ph, err := os.Open(photo)
	if err != nil {
		return "", fmt.Errorf("open photo file: %w", err)
	}

	jp, err := png.Decode(ph)
	if err != nil {
		return "", fmt.Errorf("decode jpeg photo: %w", err)
	}

	size := lsb.GetMessageSizeFromImage(jp)

	return string(lsb.Decode(size, jp)), nil
}
