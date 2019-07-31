package notes

import (
	"errors"
	"regexp"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/twilio"
)

type rule struct {
	*regexp.Regexp
	action func(string, []twilio.Media) (string, error)
}

// Rules has an ordered list of regexp+callback rules
type Rules []rule

// NewRules creates default rule set
func NewRules(token string) (Rules, error) {
	cl, err := dropbox.NewClient(dropbox.Client{Token: token})
	if err != nil {
		return nil, err
	}
	return []rule{
		{Regexp: regexp.MustCompile(`(?i)^\s*inspire\s+me\s*$`), action: inspireMe(cl)},
		{Regexp: regexp.MustCompile(`(?i)^\s*remind\s+me\s*`), action: remind(cl)},
		{Regexp: deferPattern, action: deferMessage(cl)},
		{Regexp: regexp.MustCompile(``), action: todo(cl)},
	}, nil
}

var errNoRule = errors.New("no rules matched")

// Dispatch selects and runs rules based on input
func (rules Rules) Dispatch(input string, media []twilio.Media) (string, error) {
	for _, r := range rules {
		if !r.MatchString(input) {
			continue
		}
		return r.action(input, media)
	}

	return "", errNoRule
}
