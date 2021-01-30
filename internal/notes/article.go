package notes

import (
	"fmt"
	"io"
	"io/ioutil"

	"github.com/tailscale/hujson"
)

type Article struct {
	Title string

	// Filename will be set after parsing.
	Filename string `json:"-"`

	// URL will be set after parsing.
	URL string `json:"-"`

	// Raw tells the parser not to include the standard header and footer.
	Raw bool

	Tags []string

	ReviewedOn *string `json:"reviewed_on" db:"reviewed_on"`

	ReviewBy *string `json:"review_by" db:"review_by"`

	Extra map[string]string

	Body []byte
}

func ReadArticle(r io.Reader) (Article, error) {
	var a Article
	d := hujson.NewDecoder(r)
	err := d.Decode(&a)
	if err != nil {
		return a, fmt.Errorf("hujson.Decoder.Decode: %w", err)
	}
	a.Body, err = ioutil.ReadAll(d.Buffered())
	if err != nil {
		return a, fmt.Errorf("hujson.Decoder.Buffered+ioutil.ReadAll: %w", err)
	}

	c, err := ioutil.ReadAll(r)
	if err != nil {
		return a, err
	}

	a.Body = append(a.Body, c...)

	return a, err
}

