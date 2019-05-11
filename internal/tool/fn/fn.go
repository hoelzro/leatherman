package fn

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"

	"github.com/frioux/leatherman/pkg/shellquote"
	"golang.org/x/xerrors"
)

var dir = os.Getenv("HOME") + "/code/dotfiles/bin"

/*
Run creates persistent functions by actually writing scripts.  Example usage:

```
fn count-users 'wc -l < /etc/passwd'
```

Command: fn
*/
func Run(args []string, _ io.Reader) error {
	if len(args) < 3 {
		fmt.Fprintf(os.Stderr, "Usage: %s $scriptname [-f] $command $tokens\n", args[0])
		os.Exit(1)
	}

	script := dir + "/" + args[1]

	if args[2] == "-f" {
		os.Remove(script)
		args = append(args[:2], args[3:]...)
	}

	var body string
	if len(args[2:]) == 1 {
		body = args[2]
	} else {
		var err error
		body, err = shellquote.Quote(args[2:])
		if err != nil {
			return xerrors.Errorf("Couldn't quote args to script script: %w", err)
		}
	}

	// If script exists or we can't stat it
	stat, err := os.Stat(script)
	if stat != nil {
		return xerrors.Errorf("Script ("+script+") already exists: %w", err)
	} else if !os.IsNotExist(err) {
		return xerrors.Errorf("Couldn't stat new script: %w", err)
	}

	if err := ioutil.WriteFile(script, []byte("#!/bin/sh\n\n"+body+"\n"), 0755); err != nil {
		return xerrors.Errorf("Couldn't create new script: %w", err)
	}

	return nil
}
