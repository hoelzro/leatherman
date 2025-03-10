// Code generated by maint/generate-dispatch. DO NOT EDIT.

package main

import (
	"io"

	"github.com/frioux/leatherman/internal/tool/allpurpose/clocks"
	"github.com/frioux/leatherman/internal/tool/allpurpose/csv"
	"github.com/frioux/leatherman/internal/tool/allpurpose/debounce"
	"github.com/frioux/leatherman/internal/tool/allpurpose/dumpmozlz4"
	"github.com/frioux/leatherman/internal/tool/allpurpose/expandurl"
	"github.com/frioux/leatherman/internal/tool/allpurpose/fn"
	"github.com/frioux/leatherman/internal/tool/allpurpose/groupbydate"
	"github.com/frioux/leatherman/internal/tool/allpurpose/minotaur"
	"github.com/frioux/leatherman/internal/tool/allpurpose/netrcpassword"
	"github.com/frioux/leatherman/internal/tool/allpurpose/pomotimer"
	"github.com/frioux/leatherman/internal/tool/allpurpose/replaceunzip"
	"github.com/frioux/leatherman/internal/tool/allpurpose/srv"
	"github.com/frioux/leatherman/internal/tool/allpurpose/toml"
	"github.com/frioux/leatherman/internal/tool/allpurpose/uni"
	"github.com/frioux/leatherman/internal/tool/allpurpose/yaml"
	"github.com/frioux/leatherman/internal/tool/chat/automoji"
	"github.com/frioux/leatherman/internal/tool/chat/slack"
	"github.com/frioux/leatherman/internal/tool/leatherman/update"
	"github.com/frioux/leatherman/internal/tool/mail/email"
	"github.com/frioux/leatherman/internal/tool/misc/backlight"
	"github.com/frioux/leatherman/internal/tool/misc/bamboo"
	"github.com/frioux/leatherman/internal/tool/misc/desktop"
	"github.com/frioux/leatherman/internal/tool/misc/img"
	"github.com/frioux/leatherman/internal/tool/misc/prependhist"
	"github.com/frioux/leatherman/internal/tool/misc/smlist"
	"github.com/frioux/leatherman/internal/tool/misc/status"
	"github.com/frioux/leatherman/internal/tool/misc/steamsrv"
	"github.com/frioux/leatherman/internal/tool/misc/twilio"
	"github.com/frioux/leatherman/internal/tool/misc/wuphf"
	"github.com/frioux/leatherman/internal/tool/notes/amygdala"
	"github.com/frioux/leatherman/internal/tool/notes/brainstem"
	"github.com/frioux/leatherman/internal/tool/notes/now"
	"github.com/frioux/leatherman/internal/tool/notes/proj"
	"github.com/frioux/leatherman/internal/tool/notes/zine"
)

func init() {
	Dispatch = map[string]func([]string, io.Reader) error{
		"addrs":                email.Addrs,
		"addrspec-to-tabs":     email.ToTabs,
		"alluni":               uni.All,
		"amygdala":             amygdala.Amygdala,
		"auto-emote":           automoji.Run,
		"backlight":            backlight.Run,
		"brainstem":            brainstem.Brainstem,
		"clocks":               clocks.Run,
		"csv2json":             csv.ToJSON,
		"csv2md":               csv.ToMarkdown,
		"debounce":             debounce.Run,
		"draw":                 img.Draw,
		"dump-mozlz4":          dumpmozlz4.Run,
		"email2json":           email.ToJSON,
		"expand-url":           expandurl.Run,
		"export-bamboohr":      bamboo.ExportDirectory,
		"export-bamboohr-tree": bamboo.ExportOrgChart,
		"fn":                   fn.Run,
		"group-by-date":        groupbydate.Run,
		"media-remote":         desktop.MediaRemote,
		"minotaur":             minotaur.Run,
		"name2rune":            uni.ToRune,
		"netrc-password":       netrcpassword.Run,
		"notes":                now.Serve,
		"pomotimer":            pomotimer.Run,
		"prepend-hist":         prependemojihist.Run,
		"proj":                 proj.Proj,
		"render-mail":          email.Render,
		"replace-unzip":        replaceunzip.Run,
		"slack-deaddrop":       slack.Deaddrop,
		"slack-open":           slack.Open,
		"slack-status":         slack.Status,
		"sm-list":              smlist.Run,
		"srv":                  srv.Serve,
		"status":               status.Status,
		"steamsrv":             steamsrv.Serve,
		"toml2json":            toml.ToJSON,
		"twilio":               twilio.Twilio,
		"uni":                  uni.Describe,
		"update":               update.Update,
		"wuphf":                wuphf.Wuphf,
		"yaml2json":            yaml.ToJSON,
		"zine":                 zine.Run,

		"help":    Help,
		"version": Version,
		"explode": Explode,
	}
}
