package notes

import (
	"crypto/sha1"
	"encoding/hex"
	"strings"
	"time"

	"github.com/frioux/amygdala/internal/dropbox"
	"github.com/frioux/amygdala/internal/personality"
	"github.com/frioux/amygdala/internal/reminders"
	"github.com/frioux/amygdala/internal/twilio"
	"golang.org/x/xerrors"
)

// remind format:
// remind me (?:to )xyz (at $time|in $duration)
func remind(cl dropbox.Client) func(string, []twilio.Media) (string, error) {
	return func(message string, media []twilio.Media) (string, error) {
		when, what, err := reminders.Parse(time.Now(), message)
		if err != nil {
			return personality.UserErr(err), err
		}

		for _, m := range media {
			what += " " + m.URL
		}

		sha := sha1.Sum([]byte(what))
		id := hex.EncodeToString(sha[:])
		ts := when.Format(time.RFC3339)
		path := "/notes/.alerts/" + ts + "_" + id + ".txt"

		buf := strings.NewReader(what)

		up := dropbox.UploadParams{Path: path, Autorename: true}
		if err := cl.Create(up, buf); err != nil {
			return personality.Err(), xerrors.Errorf("dropbox.Create: %w", err)
		}

		return personality.Ack() + "; will remind you @ " + ts, nil
	}
}
