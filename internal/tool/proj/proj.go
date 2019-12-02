package proj

import (
	"errors"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/user"
)

var proj, vimSessions, notes, smartcd string

func init() {
	u, err := user.Current()
	if err != nil {
		panic("Couldn't get current user: " + err.Error())
	}
	vimSessions = u.HomeDir + "/.vvar/sessions"
	notes = u.HomeDir + "/code/notes/content/posts"
	smartcd = u.HomeDir + "/.smartcd/scripts"

	proj = os.Getenv("PROJ")
}

/*
Proj integrates my various project management tools.

Usage is:

```bash
$ cd my/cool-project
$ proj init cool-project
```

The following flags are supported before the project name:

 * `-skip-vim`: does not create vim session
 * `-force-vim`: overwrites any existing vim session
 * `-skip-note`: does not create note
 * `-force-note`: overwrites any existing note
 * `-skip-smartcd`: does not create smartcd enter script
 * `-force-smartcd`: overwrites any existing smartcd enter script

Once a project has been initialized, you can run:

```bash
$ proj vim
```

To open the vim session for that project.

I use [vim sessions][vim], [a notes system][notes], and of course checkouts of
code all over the place.  Proj is meant to make creation of a vim session and a
note easy and eventually allow jumping back and forth between the three.  As of
2019-12-02 it is almost painfully specific to my personal setup, but as I
discover the actual patterns I'll probably generalize.

Proj uses uses [smartcd][smartcd] both as a mechanism and as the means to
add functionality to projects within shell sessions.

[vim]: https://blog.afoolishmanifesto.com/posts/vim-session-workflow/
[notes]: https://blog.afoolishmanifesto.com/posts/a-love-letter-to-plain-text/#notes
[smartcd]: https://github.com/cxreg/smartcd

Command: proj
*/
func Proj(args []string, _ io.Reader) error {
	if len(args) < 2 {
		return fmt.Errorf("usage: %s init | vim | note", args[0])
	}
	switch args[1] {
	case "init":
		return initialize(args[1:])
	case "vim":
		return vim()
	case "note":
		return errors.New("nyi")
	default:
		return errors.New("unknown subcommand " + args[1])
	}

}

// XXX would be nice to make this use exec instead; I know go doens't
// technically support that but I also know it does actually work.
func vim() error {
	if proj == "" {
		return errors.New("cannot infer session without PROJ set")
	}

	vim := exec.Command("vim", "-S", vimSessions+"/"+proj)
	vim.Stdin = os.Stdin
	vim.Stdout = os.Stdout
	vim.Stderr = os.Stderr
	return vim.Run()
}
