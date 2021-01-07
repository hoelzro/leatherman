package automoji

import (
	"encoding/json"
	"io/ioutil"
	"regexp"

	"github.com/frioux/leatherman/internal/dropbox"
)

type JSONRE struct {
	*regexp.Regexp
}

func (r *JSONRE) UnmarshalJSON(b []byte) error {
	var (
		err error
		str string
	)

	if err := json.Unmarshal(b, &str); err != nil {
		return err
	}

	r.Regexp, err = regexp.Compile(str)
	return err
}

type matcher struct {
	// required will disable the use of randomness and always
	// insert the specified response
	Required bool

	// JSONRE matches against input
	JSONRE JSONRE

	// emoji to add to the response; use just the emoji (☃) to add standard
	// responses, or a string (:constanza:920343920483092840) to use one of
	// the custom emoji.  To get the id in the custom one above go on discord
	// and type `\:whatever:` and you'll see the id printed out.
	Emoji string
}

func (m matcher) MatchString(s string) bool {
	return m.JSONRE.MatchString(s)
}

func loadMatchers(dbCl dropbox.Client, path string) ([]matcher, error) {
	r, err := dbCl.Download(path)
	if err != nil {
		return nil, err
	}

	var m []matcher

	d := json.NewDecoder(r)
	if err := d.Decode(&m); err != nil {
		return nil, err
	}

	return m, nil
}

func loadLua(dbCl dropbox.Client, path string) (string, error) {
	r, err := dbCl.Download(path)
	if err != nil {
		return "", err
	}

	b, err := ioutil.ReadAll(r)
	if err != nil {
		return "", err
	}

	return string(b), err
}
