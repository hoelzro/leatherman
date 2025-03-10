package email

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/mail"
	"os"
	"path/filepath"
)

type email struct {
	Header map[string]string
}

func toJSON(r io.Reader, w io.Writer) error {
	enc := json.NewEncoder(w)

	e, err := mail.ReadMessage(r)
	if err != nil {
		return fmt.Errorf("mail.ReadMessage: %w", err)
	}

	dec := new(mime.WordDecoder)

	eml := email{Header: make(map[string]string)}
	for k := range e.Header {
		header, err := dec.DecodeHeader(e.Header.Get(k))
		if err != nil {
			panic(err)
		}
		eml.Header[k] = header
	}

	err = enc.Encode(eml)
	if err != nil {
		return fmt.Errorf("json.Encode: %w", err)
	}

	return nil
}

func toJSONFromFile(path string, w io.Writer) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return fmt.Errorf("os.Open: %w", err)
	}
	defer file.Close()

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%s", r)
		}
	}()

	if err := toJSON(file, w); err != nil {
		return fmt.Errorf("%w (%s)", err, path)
	}
	return nil
}

func ToJSON(args []string, stdin io.Reader) error {
	if len(args) < 2 {
		return errors.New("please pass one or more globs")
	}

	for _, glob := range args[1:] {
		matches, err := filepath.Glob(glob)
		if err != nil {
			return fmt.Errorf("filepath.Glob: %w", err)
		}

		for _, path := range matches {
			if err := toJSONFromFile(path, os.Stdout); err != nil {
				fmt.Fprintln(os.Stderr, err)
			}
		}
	}

	return nil
}
