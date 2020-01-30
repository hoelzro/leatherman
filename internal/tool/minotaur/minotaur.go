package minotaur

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"time"

	"github.com/fsnotify/fsnotify"
)

/*
Run watches one or more directories (before the `--`) and runs a script when
events in those directories occur.  The script receives the events as arguments,
so you can exit early if only irrelevant files changed.

```bash
minotaur -suppress-args -include internal -ignore yaml \
   ~/code/leatherman -- \
   go test ~/code/leatherman/...
```

The arguments are of the form `$event\t$filename`; for example `CREATE	x.pl`.
As far as I know the valid events are;

 * `CHMOD`
 * `CREATE`
 * `REMOVE`
 * `RENAME`
 * `WRITE`

Include the `-suppress-args` flag to disable the arguments.

The events are deduplicated and also debounced, so your script will never fire
more often than once a second.  If events are happening every half second the
debouncing will cause the script to never run.

The underlying library supports emitting multiple events in a single line (ie
`CREATE|CHMOD`) though I've not seen that in Linux.

`minotaur` reëmits all output (both stderr and stdout) of the passed script to
standard out, so you could make a script like this to experiment with the
events with timestamps:

```bash
#!/bin/sh

for x in "$@"; do
   echo "$x"
done | ts
```

You can do all kinds of interesting things in the script, for example you could
verify that the events deserve a restart, then restart a service, then block till
the service can serve traffic, then restart some other related service.

The `-include` and `-ignore` arguments are optional; by default `-include` is
empty, so matches everything, and `-ignore` matches `.git`.  You can also pass
`-verbose` to include output about minotaur itself, like which directories it's
watching.

The flag `-run-at-start` will run the script before there are any events.

Command: minotaur
*/
func Run(args []string, _ io.Reader) error {
	args = args[1:]

	c, err := parseFlags(args)
	if err != nil {
		return fmt.Errorf("parseFlags: %w", err)
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return fmt.Errorf("fsnotify.NewWatcher: %w", err)
	}
	defer watcher.Close()

	done := make(chan bool)
	var timeout <-chan time.Time
	events := make(map[string]bool)

	if c.runAtStart {
		timeout = time.After(0)
	}

	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				if event.Op&fsnotify.Create == fsnotify.Create {
					stat, err := os.Stat(event.Name)
					if err != nil {
						fmt.Fprintf(os.Stderr, "Couldn't stat created thing: %s\n", err)
					} else if stat.IsDir() && c.include.MatchString(event.Name) && !c.ignore.MatchString(event.Name) {
						err := watcher.Add(event.Name)
						if err != nil {
							fmt.Fprintf(os.Stderr, "failed to watch %s: %s\n", event.Name, err)
						} else if c.verbose {
							fmt.Fprintf(os.Stderr, "watching %s\n", event.Name)
						}
					}
				}

				events[event.Op.String()+"\t"+event.Name] = true
				timeout = time.After(time.Second)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				fmt.Println("error:", err)
			case <-timeout:
				s := make([]string, 0, len(c.script)+len(events))
				s = append(s, c.script...)
				if !c.suppressArgs {
					for e := range events {
						s = append(s, e)
					}
				}
				events = make(map[string]bool)
				cmd := exec.Command(s[0], s[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				err := cmd.Run()
				if err != nil && c.verbose {
					fmt.Fprintf(os.Stderr, "script (%q) failed: %s\n", s, err)
				}
			}

		}
	}()

	for _, path := range c.dirs {
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}

			if !info.IsDir() {
				return nil
			}

			if c.ignore.MatchString(path) {
				return filepath.SkipDir
			}
			if !c.include.MatchString(path) {
				return nil
			}

			if c.verbose {
				fmt.Fprintln(os.Stderr, "watching "+path)
			}
			if err := watcher.Add(path); err != nil {
				return fmt.Errorf("fsnotify.Watcher.Add: %w", err)
			}
			return nil
		})
		if err != nil {
			return fmt.Errorf("filepath.Walk: %w", err)
		}
	}
	<-done

	return nil
}
