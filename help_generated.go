package main

var readme = []byte(`# Leatherman - fREW's favorite multitool

This is a little project simply to make trivial tools in Go effortless for my
personal usage.  These tools are almost surely of low utility to most people,
but may be instructive nonetheless.

[I have CI/CD to build this into a single
binary](https://github.com/frioux/leatherman/blob/master/.travis.yml) and [an
` + "`" + `explode` + "`" + ` tool that builds
symlinks](https://github.com/frioux/leatherman/blob/master/explode.go) for each
tool in the busybox style.

[I have automation in my
dotfiles](https://github.com/frioux/dotfiles/blob/bef8303c19e2cefac7dfbec420ad8d45b95415b8/install.sh#L133-L141)
to pull the latest binary at install time and run the ` + "`" + `explode` + "`" + ` tool.

## Installation

Here's a copy pasteable script to install the leatherman on OSX or Linux:

` + "`" + `` + "`" + `` + "`" + ` bash
OS=$([ $(uname -s) = "Darwin" ] && echo "-osx")
LMURL="$(curl -s https://api.github.com/repos/frioux/leatherman/releases/latest |
   grep browser_download_url |
   cut -d '"' -f 4 |
   grep -F leatherman${OS}.xz )"
mkdir -p ~/bin
curl -sL "$LMURL" > ~/bin/leatherman.xz
xz -d -f ~/bin/leatherman.xz
chmod +x ~/bin/leatherman
~/bin/leatherman explode
` + "`" + `` + "`" + `` + "`" + `

This asssumes that ` + "`" + `~/bin` + "`" + ` is in your path.  The ` + "`" + `explode` + "`" + ` command will create a
symlink for each of the tools listed below.

## Usage

Each tool takes different args, but to run a tool you can either use a symlink
(presumably created by ` + "`" + `explode` + "`" + `):

` + "`" + `` + "`" + `` + "`" + ` bash
$ echo "---\nfoo: 1" | yaml2json
{"foo":1}
` + "`" + `` + "`" + `` + "`" + `

or use it as a subcommand:

` + "`" + `` + "`" + `` + "`" + ` bash
echo "---\nfoo: 1" | leatherman yaml2json
{"foo":1}
` + "`" + `` + "`" + `` + "`" + `

## Tools

### allpurpose

#### ` + "`" + `alluni` + "`" + `

` + "`" + `alluni` + "`" + ` prints all unicode character names.

` + "`" + `` + "`" + `` + "`" + `bash
$ alluni | grep DENTISTRY
DENTISTRY SYMBOL LIGHT VERTICAL AND TOP RIGHT
DENTISTRY SYMBOL LIGHT VERTICAL AND BOTTOM RIGHT
DENTISTRY SYMBOL LIGHT VERTICAL WITH CIRCLE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL WITH CIRCLE
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL WITH CIRCLE
DENTISTRY SYMBOL LIGHT VERTICAL WITH TRIANGLE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL WITH TRIANGLE
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL WITH TRIANGLE
DENTISTRY SYMBOL LIGHT VERTICAL AND WAVE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL WITH WAVE
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL WITH WAVE
DENTISTRY SYMBOL LIGHT DOWN AND HORIZONTAL
DENTISTRY SYMBOL LIGHT UP AND HORIZONTAL
DENTISTRY SYMBOL LIGHT VERTICAL AND TOP LEFT
DENTISTRY SYMBOL LIGHT VERTICAL AND BOTTOM LEFT
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `clocks` + "`" + `

` + "`" + `clocks` + "`" + ` shows my personal, digital, wall of clocks.  Pass one or more timezone names
to choose which timezones are shown.

` + "`" + `` + "`" + `` + "`" + `bash
clocks Africa/Johannesburg America/Los_Angeles Europe/Copenhagen
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `csv2json` + "`" + `

` + "`" + `csv2json` + "`" + ` reads CSV on stdin and writes JSON on stdout; first line of input is the
header, and thus the keys of the JSON.

#### ` + "`" + `csv2md` + "`" + `

` + "`" + `csv2md` + "`" + ` reads CSV on stdin and writes Markdown on stdout.

#### ` + "`" + `debounce` + "`" + `

` + "`" + `debounce` + "`" + ` debounces input from stdin to stdout

The default lockout time is one second, you can override that with the
` + "`" + `--lockoutTime` + "`" + ` argument.  By default the trailing edge triggers output, so
output is emitted after there is no input for the lockout time.  You can change
this behavior by passing the ` + "`" + `--leadingEdge` + "`" + ` flag.

#### ` + "`" + `dump-mozlz4` + "`" + `

` + "`" + `dump-mozlz4` + "`" + ` dumps the contents of a ` + "`" + `mozlz4` + "`" + ` (aka ` + "`" + `jsonlz4` + "`" + `) file commonly used by
Firefox.  Just takes the name of the file to dump and writes to standard out.

#### ` + "`" + `expand-url` + "`" + `

` + "`" + `expand-url` + "`" + ` reads text on STDIN and writes the same text back, converting any links to
Markdown links, with the title of the page as the title of the link.

#### ` + "`" + `fn` + "`" + `

` + "`" + `fn` + "`" + ` creates persistent functions by actually writing scripts.  Example usage:

` + "`" + `` + "`" + `` + "`" + `
fn count-users 'wc -l < /etc/passwd'
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `group-by-date` + "`" + `

` + "`" + `group-by-date` + "`" + ` creates time series data by counting lines and grouping them by a given date
format.  takes dates on stdin in format -i, will group them by format -g, and
write them in format -o.

#### ` + "`" + `minotaur` + "`" + `

` + "`" + `minotaur` + "`" + ` watches one or more directories (before the ` + "`" + `--` + "`" + `) and runs a script when
events in those directories occur.

` + "`" + `` + "`" + `` + "`" + `bash
minotaur -include-args -include internal -ignore yaml \
   ~/code/leatherman -- \
   go test ~/code/leatherman/...
` + "`" + `` + "`" + `` + "`" + `

If the ` + "`" + `-include-args` + "`" + ` flag is set, the script receives the events as
arguments, so you can exit early if only irrelevant files changed.

The arguments are of the form ` + "`" + `$event\t$filename` + "`" + `; for example ` + "`" + `CREATE	x.pl` + "`" + `.
As far as I know the valid events are;

 * ` + "`" + `CHMOD` + "`" + `
 * ` + "`" + `CREATE` + "`" + `
 * ` + "`" + `REMOVE` + "`" + `
 * ` + "`" + `RENAME` + "`" + `
 * ` + "`" + `WRITE` + "`" + `

The events are deduplicated and also debounced, so your script will never fire
more often than once a second.  If events are happening every half second the
debouncing will cause the script to never run.

The underlying library supports emitting multiple events in a single line (ie
` + "`" + `CREATE|CHMOD` + "`" + `) though I've not seen that in Linux.

` + "`" + `minotaur` + "`" + ` reëmits all output (both stderr and stdout) of the passed script to
standard out, so you could make a script like this to experiment with the
events with timestamps:

` + "`" + `` + "`" + `` + "`" + `bash
#!/bin/sh

for x in "$@"; do
   echo "$x"
done | ts
` + "`" + `` + "`" + `` + "`" + `

You can do all kinds of interesting things in the script, for example you could
verify that the events deserve a restart, then restart a service, then block till
the service can serve traffic, then restart some other related service.

The ` + "`" + `-include` + "`" + ` and ` + "`" + `-ignore` + "`" + ` arguments are optional; by default ` + "`" + `-include` + "`" + ` is
empty, so matches everything, and ` + "`" + `-ignore` + "`" + ` matches ` + "`" + `.git` + "`" + `.  You can also pass
` + "`" + `-verbose` + "`" + ` to include output about minotaur itself, like which directories it's
watching.

The flag ` + "`" + `-no-run-at-start` + "`" + ` will not the the script until there are any events.

The flag ` + "`" + `-report` + "`" + ` will decorate output with a text wrapper to clarify when the
script is run.

#### ` + "`" + `name2rune` + "`" + `

` + "`" + `name2rune` + "`" + ` takes the name of a unicode character and prints out any found
characters.

#### ` + "`" + `netrc-password` + "`" + `

` + "`" + `netrc-password` + "`" + ` prints password for the passed hostname and login.

` + "`" + `` + "`" + `` + "`" + `bash
$ netrc-password google.com me@gmail.com
supersecretpassword
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `pomotimer` + "`" + `

` + "`" + `pomotimer` + "`" + ` starts a timer for 25m or the duration expressed in the first
argument.

` + "`" + `` + "`" + `` + "`" + `
pomotimer 2.5m
` + "`" + `` + "`" + `` + "`" + `

or

` + "`" + `` + "`" + `` + "`" + `
pomotimer 3m12s
` + "`" + `` + "`" + `` + "`" + `

Originally a timer for use with [the pomodoro][1] [technique][2].  Handy timer in any case
since you can pass it arbitrary durations, pause it, reset it, and see it's
progress.

[1]: https://blog.afoolishmanifesto.com/posts/the-pomodoro-technique/
[2]: https://blog.afoolishmanifesto.com/posts/the-pomodoro-technique-three-years-later/

#### ` + "`" + `replace-unzip` + "`" + `

` + "`" + `replace-unzip` + "`" + ` extracts zipfiles, but does not extract ` + "`" + `.DS_Store` + "`" + ` or ` + "`" + `__MACOSX/` + "`" + `.
Automatically extracts into a directory named after the zipfile if there is not
a single root for all files in the zipfile.

#### ` + "`" + `srv` + "`" + `

` + "`" + `srv` + "`" + ` will serve a directory over http, injecting javascript to have pages
reload when files change.

It takes an optional dir to serve, the default is ` + "`" + `.` + "`" + `.

` + "`" + `` + "`" + `` + "`" + `bash
$ srv ~
Serving /home/frew on [::]:21873
` + "`" + `` + "`" + `` + "`" + `

You can pass -port if you care to choose the listen port.

It will set up file watchers and trigger page reloads (via SSE,) this
functionality can be disabled with -no-autoreload.

` + "`" + `` + "`" + `` + "`" + `bash
$ srv -port 8080 -no-autoreload ~
Serving /home/frew on [::]:8080
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `toml2json` + "`" + `

` + "`" + `toml2json` + "`" + ` reads [TOML](https://github.com/toml-lang/toml) on stdin and writes JSON
on stdout.


` + "`" + `` + "`" + `` + "`" + `bash
$ echo 'foo = "bar"` + "`" + ` | toml2json
{"foo":"bar"}
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `uni` + "`" + `

` + "`" + `uni` + "`" + ` will describe the characters in the args.

` + "`" + `` + "`" + `` + "`" + `bash
$ uni ⢾
'⢾' @ 10430 aka BRAILLE PATTERN DOTS-234568 ( graphic | printable | symbol )
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `yaml2json` + "`" + `

` + "`" + `yaml2json` + "`" + ` reads YAML on stdin and writes JSON on stdout.

### chat

#### ` + "`" + `auto-emote` + "`" + `

` + "`" + `auto-emote` + "`" + ` comments to discord and reacts to all messages with vaguely related emoji.

The following env vars should be set:

 * LM_DROPBOX_TOKEN should be set to load a responses.json.
 * LM_BOT_LUA_PATH should be set to the location of lua to process emoji data within dropbox.
 * LM_DISCORD_TOKEN should be set for this to actually function.

Here's an example of lua code that works for this:

` + "`" + `` + "`" + `` + "`" + `lua
function add(w)
        es:addoptional(w:char())
end

function addtoken(tok)
        if tok == "" then
                return
        end

        bn = turtleemoji.findbyname(tok)
        if not bn == nil then
                add(bn)
        end

        for i, te in ipairs(turtleemoji.searchbycategory(tok)) do
                add(te)
        end

        for i, te in ipairs(turtleemoji.searchbykeyword(tok)) do
                add(te)
        end

        if not es:optionallen() == 0 then
                return
        end

        -- this returns so many results for basic words that we only
        -- use it if nothing else found anything for the whole message
        for i, te in ipairs(turtleemoji.search(tok)) do
                add(te)
        end
end

for i, w in ipairs(es:words()) do
        addtoken(w)
end
` + "`" + `` + "`" + `` + "`" + `

The lua code has a global var called ` + "`" + `es` + "`" + ` (for emoji set) and an imported
package called ` + "`" + `turtleemoji` + "`" + `.  ` + "`" + `es` + "`" + ` is how you access the current message,
currently added emoji, etc.  Here are the methods on ` + "`" + `es` + "`" + `:

#### ` + "`" + `es:optional()` + "`" + ` // table of string to bool

Returns a copy of the optional emoji.  Modifications of the table will not
affect the final result; other methods should be used for modification.

#### ` + "`" + `es:addoptional("💀")` + "`" + `

Adds an emoji to randomly include in the reaction.

#### ` + "`" + `es:hasoptional("💀")` + "`" + ` // bool

Returns true of the passed emoji is in the list of optional emoji to include
(at random) on the reaction.

#### ` + "`" + `es:removeoptional("💀")` + "`" + `

Remove the passed emoji from the optionally included emoji.

#### ` + "`" + `es:required()` + "`" + ` // table of required emoji

Returns a copy of the required emoji.  Modifications of the table will not
affect the final result; other methods should be used for modification.

#### ` + "`" + `es:hasrequired("💀")` + "`" + ` // bool

Returns true if the passed emoji is going to be included in the reaction.

#### ` + "`" + `es:addrequired("💀")` + "`" + `

Add an emoji to the reaction.

#### ` + "`" + `es:removerequired("💀")` + "`" + `

Remove an emoji that is going to be included in the reaction.

#### ` + "`" + `es:message()` + "`" + ` // string

Returns the message that triggered the reaction.

#### ` + "`" + `es:messagematches("regexp")` + "`" + ` // bool

True if the message matches the passed regex.
[Docs for regex syntax are here](https://golang.org/pkg/regexp/syntax/).

#### ` + "`" + `es:words()` + "`" + ` // table of tokenized words

Returns a copy of the tokenized words.  Tokenization of words happens on all
non-alpha characters and the message is lowerecased.

#### ` + "`" + `es:hasword("word")` + "`" + ` // bool

True if the word is included in the message.

#### ` + "`" + `es:len()` + "`" + ` // int

Total count of items in optional and required.

All of the following are thin veneers atop
[github.com/hackebrot/turtle](https://github.com/hackebrot/turtle):

 * ` + "`" + `turtle.findbyname("skull")` + "`" + ` // turtleemoji
 * ` + "`" + `turtle.findbychar("💀")` + "`" + ` // turtleemoji
 * ` + "`" + `turtle.searchbycategory("people")` + "`" + ` // table of turtleemoji
 * ` + "`" + `turtle.searchbykeyword("animal")` + "`" + ` // table of turtleemoji
 * ` + "`" + `turtle.search("foo")` + "`" + ` // table of turtleemoji
 * ` + "`" + `turtleemoji#name()` + "`" + ` // string
 * ` + "`" + `turtleemoji#category()` + "`" + ` // string
 * ` + "`" + `turtleemoji#char()` + "`" + ` // string
 * ` + "`" + `turtleemoji#haskeyword("keyword")` + "`" + ` // bool

#### ` + "`" + `slack-deaddrop` + "`" + `

` + "`" + `slack-deaddrop` + "`" + ` allows sending messages to a slack channel without looking at slack.
Typical usage is probably something like:

` + "`" + `` + "`" + `` + "`" + `bash
$ slack-deaddrop -channel general -text 'good morning!'
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `slack-open` + "`" + `

` + "`" + `slack-open` + "`" + ` opens a channel, group message, or direct message by name:

` + "`" + `` + "`" + `` + "`" + `bash
$ slack-open -channel general
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `slack-status` + "`" + `

` + "`" + `slack-status` + "`" + ` sets the current users's status.

` + "`" + `` + "`" + `` + "`" + `bash
$ slack-status -text "working for the weekend" -emoji :guitar:
` + "`" + `` + "`" + `` + "`" + `

### leatherman

#### ` + "`" + `update` + "`" + `

` + "`" + `update` + "`" + ` checks to see if there's an update from github and installs it if there
is.  If LM_GH_TOKEN is set to a personal access token this can be called more
frequently without exhausting github api limits.

### mail

#### ` + "`" + `addrs` + "`" + `

` + "`" + `addrs` + "`" + ` sorts the addresses passed on stdin (in the mutt addrbook format, see
` + "`" + `addrspec-to-tabs` + "`" + `) and sorts them based on how recently they were used, from
the glob passed on the arguments.  The tool exists so that you can create an
address list either with an export tool (like ` + "`" + `goobook` + "`" + `), a subset of your sent
addresses, or whatever else, and then you can sort it based on your sent mail
folder.

` + "`" + `` + "`" + `` + "`" + ` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" >sortedaddrs.txt
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `addrspec-to-tabs` + "`" + `

` + "`" + `addrspec-to-tabs` + "`" + ` converts email addresses from the standard format (` + "`" + `"Hello Friend"
<foo@bar>` + "`" + `) from stdin to the mutt address book format, ie tab separated fields,
on stdout.

` + "`" + `` + "`" + `` + "`" + ` bash
$ <someaddrs.txt addrs "$HOME/mail/gmail.sent/cur/*" | addrspec-to-tabs >addrbook.txt
` + "`" + `` + "`" + `` + "`" + `

This tool ignores the comment because, after actually auditing my addressbook,
most comments are incorrectly recognized by all tools. (for example:
` + "`" + `<5555555555@vzw.com> (555) 555-5555` + "`" + ` should not have a comment of ` + "`" + `(555)` + "`" + `.)

It exists to be combined with ` + "`" + `addrs` + "`" + ` and mutt.

#### ` + "`" + `email2json` + "`" + `

` + "`" + `email2json` + "`" + ` produces a JSON representation of an email from a list of globs.  Only
headers are currently supported, patches welcome to support bodies.

` + "`" + `` + "`" + `` + "`" + `bash
$ email2json '/home/frew/var/mail/mitsi/cur/*' | head -1 | jq .
{
  "Header": {
    "Content-Type": "multipart/alternative; boundary=00163642688b8ef3070464661533",
    "Date": "Thu, 5 Mar 2009 15:45:17 -0600",
    "Delivered-To": "xxx",
    "From": "fREW Schmidt <xxx>",
    "Message-Id": "<fb3648c60903051345o728960f5l6cfb9e1f324bbf50@mail.gmail.com>",
    "Mime-Version": "1.0",
    "Received": "by 10.103.115.8 with HTTP; Thu, 5 Mar 2009 13:45:17 -0800 (PST)",
    "Subject": "STATION",
    "To": "xyzzy@googlegroups.com"
  }
}
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `render-mail` + "`" + `

` + "`" + `render-mail` + "`" + ` is a scrappy tool to render email with a Local-Date included, if Date is
not already in local time.

### misc

#### ` + "`" + `backlight` + "`" + `

` + "`" + `backlight` + "`" + ` is a faster version of ` + "`" + `xbacklight` + "`" + ` by directly writing to ` + "`" + `/sys` + "`" + `.  Example:

Increase by 10%:

` + "`" + `` + "`" + `` + "`" + ` bash
$ backlight 10
` + "`" + `` + "`" + `` + "`" + `

Decrease by 10%:

` + "`" + `` + "`" + `` + "`" + ` bash
$ backlight -10
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `export-bamboohr` + "`" + `

` + "`" + `export-bamboohr` + "`" + ` exports entire company directory as JSON.

#### ` + "`" + `export-bamboohr-tree` + "`" + `

` + "`" + `export-bamboohr-tree` + "`" + ` exports company org chart as JSON.

#### ` + "`" + `prepend-hist` + "`" + `

` + "`" + `prepend-hist` + "`" + ` prints out deduplicated lines from the history file in reverse order and
then prints out the lines from STDIN, filtering out what's already been printed.

` + "`" + `` + "`" + `` + "`" + `bash
$ alluni | prefix-hist ~/.uni_history
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `sm-list` + "`" + `

` + "`" + `sm-list` + "`" + ` lists all of the available [Sweet Maria's](https://www.sweetmarias.com/) coffees
as json documents per line.  Here's how you might see the top ten coffees by
rating:

` + "`" + `` + "`" + `` + "`" + `bash
$ sm-list | jq -r '[.Score, .Title, .URL ] | @tsv' | sort -n | tail -10
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `status` + "`" + `

` + "`" + `status` + "`" + ` runs a little web server that surfaces status information related to how
I'm using the machine.  For example, it can say which window is active, what
firefox tabs are loaded, if the screen is locked, etc.  The main benefit of the
tool is that it caches the values returned.

In the background, it interact swith the [blink(1)](http://blink1.thingm.com/).
It turns the light green when I'm in a meeting and red when audio is playing.

#### ` + "`" + `twilio` + "`" + `

` + "`" + `twilio` + "`" + ` allows interacting with a service that recieves callbacks from twilio
for testing.

It takes four arguments:

 * ` + "`" + `-endpoint` + "`" + `: the url to hit (` + "`" + `http://localhost:8080/twilio` + "`" + `, for example)
 * ` + "`" + `-auth` + "`" + `: the auth token to use
 * ` + "`" + `-message` + "`" + `: the message to send
 * ` + "`" + `-from` + "`" + `: the phone number the message is from (` + "`" + `+15555555555` + "`" + `, for example)

Run ` + "`" + `twilio -help` + "`" + ` to see the defaults.

` + "`" + `` + "`" + `` + "`" + `bash
$ twilio -message "the building is on fire!"
` + "`" + `` + "`" + `` + "`" + `

#### ` + "`" + `wuphf` + "`" + `

` + "`" + `wuphf` + "`" + ` sends alerts via both ` + "`" + `wall` + "`" + ` and [pushover](https://pushover.net).  All
arguments are concatenated to produce the sent message.

The following environment variables should be set:

 * LM_PUSHOVER_TOKEN
 * LM_PUSHOVER_USER
 * LM_PUSHOVER_DEVICE

` + "`" + `` + "`" + `` + "`" + `bash
$ wuphf 'the shoes are on sale'
` + "`" + `` + "`" + `` + "`" + `

### notes

#### ` + "`" + `amygdala` + "`" + `

` + "`" + `amygdala` + "`" + ` is an automated assistant.  The main docs for it are in [amygdala.mdwn
in the root of the leatherman
repo](https://github.com/frioux/leatherman/blob/main/amygdala.mdwn).

Configure it with the following environment variables:

 * LM_DROPBOX_TOKEN
 * LM_MY_CEL
 * LM_TWILIO_URL

Takes an optional ` + "`" + `-port` + "`" + ` argument to specify which port to listen on.

#### ` + "`" + `brainstem` + "`" + `

` + "`" + `brainstem` + "`" + ` allows interacting with amygdala without using any of the server
components, typically for testing the personality etc, but can also be used as
a lightweight amygdala instance.

#### ` + "`" + `notes` + "`" + `

` + "`" + `notes` + "`" + ` provides a web interface to parts of my notes stored in
Dropbox.  This should eventually be merged into Amygdala.

#### ` + "`" + `proj` + "`" + `

` + "`" + `proj` + "`" + ` integrates my various project management tools.

Usage is:

` + "`" + `` + "`" + `` + "`" + `bash
$ cd my/cool-project
$ proj init cool-project
` + "`" + `` + "`" + `` + "`" + `

The following flags are supported before the project name:

 * ` + "`" + `-skip-vim` + "`" + `: does not create vim session
 * ` + "`" + `-force-vim` + "`" + `: overwrites any existing vim session
 * ` + "`" + `-skip-note` + "`" + `: does not create note
 * ` + "`" + `-force-note` + "`" + `: overwrites any existing note
 * ` + "`" + `-skip-smartcd` + "`" + `: does not create smartcd enter script
 * ` + "`" + `-force-smartcd` + "`" + `: overwrites any existing smartcd enter script

Once a project has been initialized, you can run:

` + "`" + `` + "`" + `` + "`" + `bash
$ proj vim
` + "`" + `` + "`" + `` + "`" + `

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

#### ` + "`" + `zine` + "`" + `

` + "`" + `zine` + "`" + ` does read only operations on notes.

## Debugging

In an effort to make debugging simpler, I've created three ways to see what
` + "`" + `leatherman` + "`" + ` is doing:

### Tracing

` + "`" + `LMTRACE=$somefile` + "`" + ` will write an execution trace to ` + "`" + `$somefile` + "`" + `; look at that with ` + "`" + `go tool trace $somefile` + "`" + `

Since so many of the tools are short lived my assumption is that the execution
trace will be the most useful.

### Profiling

` + "`" + `LMPROF=$somefile` + "`" + ` will write a cpu profile to ` + "`" + `$somefile` + "`" + `; look at that with ` + "`" + `go tool pprof -http localhost:10123 $somefile` + "`" + `

If you have a long running tool, the pprof http endpoint is exposed on
` + "`" + `localhost:6060/debug/pprof` + "`" + ` but picks a random port if that port is in use; the
port can be overridden by setting ` + "`" + `LMHTTPPROF=$someport` + "`" + `.

## Auxiliary Tools

Some tools are annoying to have in the main leatherman tool.  Maybe they pull
in deps that are huge or need cgo, but in any case I try to keep them separate.
For now, these tools are under ` + "`" + `leatherman/cmd` + "`" + ` and must be built and run
separately.  At some point I may come up with a policy around building or naming these,
but for now they are simply extra tools.
`)

var commandReadme map[string][]byte

func init() {
	commandReadme = map[string][]byte{
		"addrs": readme[11950:12448],

		"addrspec-to-tabs": readme[12448:13027],

		"alluni": readme[1583:2410],

		"amygdala": readme[16024:16402],

		"auto-emote": readme[7632:11198],

		"backlight": readme[13888:14091],

		"brainstem": readme[16402:16610],

		"clocks": readme[2410:2627],

		"csv2json": readme[2627:2768],

		"csv2md": readme[2768:2843],

		"debounce": readme[2843:3188],

		"dump-mozlz4": readme[3188:3372],

		"email2json": readme[13027:13743],

		"expand-url": readme[3372:3549],

		"export-bamboohr": readme[14091:14176],

		"export-bamboohr-tree": readme[14176:14264],

		"fn": readme[3549:3686],

		"group-by-date": readme[3686:3905],

		"minotaur": readme[3905:5725],

		"name2rune": readme[5725:5831],

		"netrc-password": readme[5831:5997],

		"notes": readme[16610:16747],

		"pomotimer": readme[5997:6490],

		"prepend-hist": readme[14264:14506],

		"proj": readme[16747:18067],

		"render-mail": readme[13743:13878],

		"replace-unzip": readme[6490:6721],

		"slack-deaddrop": readme[11198:11419],

		"slack-open": readme[11419:11554],

		"slack-status": readme[11554:11700],

		"sm-list": readme[14506:14784],

		"srv": readme[6721:7205],

		"status": readme[14784:15242],

		"toml2json": readme[7205:7381],

		"twilio": readme[15242:15701],

		"uni": readme[7381:7544],

		"update": readme[11716:11940],

		"wuphf": readme[15701:16013],

		"yaml2json": readme[7544:7622],

		"zine": readme[18067:18124],
	}
}
