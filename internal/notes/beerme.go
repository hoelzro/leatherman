package notes

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"math/rand"
	"regexp"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/frioux/amygdala/internal/twilio"
	"golang.org/x/xerrors"
)

var isItem = regexp.MustCompile(`^\s?\*\s+(.*?)\s*$`)
var mdLink = regexp.MustCompile(`^\[(.*)\]\((.*)\)$`)

var errNone = errors.New("never found anything")

func beerMe(r io.Reader) (string, error) {
	s := bufio.NewScanner(r)

	o := []string{}
	for s.Scan() {
		m := isItem.FindStringSubmatch(s.Text())
		if len(m) != 2 {
			continue
		}
		o = append(o, m[1])
	}

	if len(o) == 0 {
		return "", errNone
	}

	rand.Shuffle(len(o), func(i, j int) { o[i], o[j] = o[j], o[i] })

	fmt.Println(mdLink.FindStringSubmatch(o[0]))
	if l := mdLink.FindStringSubmatch(o[0]); len(l) == 3 {
		return fmt.Sprintf("[%s]( %s )", l[1], l[2]), nil
	}

	return o[0], nil
}

func inspireMe(cl dropbox.Client) func(_ string, _ []twilio.Media) (string, error) {
	return func(_ string, _ []twilio.Media) (string, error) {
		r, err := cl.Download("/notes/content/posts/inspiration.md")
		if err != nil {
			return personality.Err(), xerrors.Errorf("dropbox.Download: %w", err)
		}
		n, err := beerMe(r)
		if err != nil {
			return personality.Err(), err
		}
		return n, nil
	}
}
