package toml

import (
	"bytes"
	"encoding/json"
	"io"
	"os"

	parser "github.com/BurntSushi/toml"
	"golang.org/x/xerrors"
)

/*
ToJSON reads [TOML](https://github.com/toml-lang/toml) on stdin and writes JSON
on stdout.


```bash
$ echo 'foo = "bar"` | toml2json
{"foo":"bar"}
```

Command: toml2json
*/
func ToJSON(_ []string, stdin io.Reader) error {
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, stdin); err != nil {
		return xerrors.Errorf("io.Copy: %w", err)
	}
	var ret interface{}
	if err := parser.Unmarshal(buf.Bytes(), &ret); err != nil {
		return xerrors.Errorf("toml.Unmarshal: %w", err)
	}

	e := json.NewEncoder(os.Stdout)
	if err := e.Encode(ret); err != nil {
		return xerrors.Errorf("json.Encode: %w", err)
	}

	return nil
}
