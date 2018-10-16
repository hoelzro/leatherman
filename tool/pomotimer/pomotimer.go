package pomotimer

import (
	"fmt"
	"io"
	"os"
	"os/exec"
	"time"
)

const clear = "\r\x1b[J"

// Run starts a timer for 25m or the duration expressed in the first
// argument.
func Run(args []string, stdin io.Reader) error {
	timer, _ := time.ParseDuration("25m")
	if len(args) > 1 {
		var err error
		timer, err = time.ParseDuration(args[1])
		if err != nil {
			fmt.Fprintf(os.Stderr, "pomotimer: couldn't parse duration: %s", err)
			os.Exit(1)
		}
	}

	initialSeconds := int(timer.Seconds())

	setProcessName("PT" + formatTime(initialSeconds))
	fmt.Print("[p]ause [r]eset abort[!]\n\n")

	// disable input buffering
	exec.Command("stty", "-F", "/dev/tty", "cbreak", "min", "1").Run()
	// do not display entered characters on the screen
	exec.Command("stty", "-F", "/dev/tty", "-echo").Run()

	// restore the echoing state when exiting
	defer exec.Command("stty", "-F", "/dev/tty", "echo").Run()
	defer exec.Command("tmux", "setw", "automatic-rename", "on").Run()

	c := time.Tick(1 * time.Second)
	kb := make(chan string)
	go kbChan(kb, stdin)

	running := true
	secondsRemaining := initialSeconds
LOOP:
	for {
		select {
		case <-c:
			if secondsRemaining > 1 {
				if !running {
					continue
				}
				if secondsRemaining%30 == 0 {
					setProcessName("PT" + formatTime(secondsRemaining))
				}
				fmt.Print(clear+formatTime(secondsRemaining), " remaining")
				secondsRemaining--
			} else {
				fmt.Println(clear + "Take a break!\a")
				break LOOP
			}
		case key := <-kb:
			if key == "p" {
				running = !running
			} else if key == "r" {
				secondsRemaining = initialSeconds
			} else if key == "!" {
				fmt.Println(clear + "aborting timer!")
				break LOOP
			}
		}
	}

	return nil
}

func kbChan(keys chan string, stdin io.Reader) {
	var b = make([]byte, 1)
	for {
		stdin.Read(b)
		keys <- string(b)
	}
}

func formatTime(s int) string {
	return fmt.Sprintf("%02d:%02d", s/60, s%60)
}
