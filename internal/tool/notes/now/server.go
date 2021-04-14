package now

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"strings"
	"time"

	"github.com/yuin/goldmark"
	"github.com/yuin/goldmark/extension"
	"github.com/yuin/goldmark/parser"
	"github.com/yuin/goldmark/renderer/html"

	"github.com/frioux/leatherman/internal/lmfs"
	"github.com/frioux/leatherman/internal/lmhttp"
	"github.com/frioux/leatherman/internal/notes"
	"github.com/frioux/leatherman/internal/selfupdate"
)

// 3. list and sort all files with `todo-` prefix
//   * list(db) ([]file, error)
//   * render

func handlerAddItem(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		i := req.Form.Get("item")
		if i == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing item parameter")
			return nil
		}

		b, err := fs.ReadFile(fss, nowPath)
		if err != nil {
			return err
		}

		b, err = addItem(bytes.NewReader(b), time.Now(), i)
		if err != nil {
			return err
		}

		if err := lmfs.WriteFile(fss, nowPath, b, 0); err != nil {
			return err
		}

		rw.Header().Add("Location", "/")
		rw.WriteHeader(303)

		b, err = parseNow(bytes.NewReader(b), time.Now())
		if err != nil {
			return err
		}

		v := &HTMLVars{Zine: z, Title: "now"}
		if err := mdwn.Convert(b, v); err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "simple.html", v)
	})
}

func handlerFavicon() http.Handler {
	return http.HandlerFunc(func(rw http.ResponseWriter, _ *http.Request) {
		rw.Header().Add("Content-Type", "image/svg+xml")
		fmt.Fprintln(rw, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 100 100"><text y=".9em" font-size="90">☕</text></svg>`)
	})
}

func handlerList(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		stmt, err := z.Preparex(`SELECT title, url FROM articles ORDER BY title`)
		if err != nil {
			return err
		}

		v := listVars{HTMLVars: &HTMLVars{Zine: z}}
		if err := stmt.Select(&v.Articles); err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "list.html", v)
	})
}

func handlerQ(z *notes.Zine, mdwn goldmark.Markdown) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		v := qVars{HTMLVars: &HTMLVars{Zine: z}}
		q := req.URL.Query().Get("q")
		if q == "" {
			q = "SELECT * FROM articles"
		}
		var err error
		v.Records, err = z.Q(q)
		if err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "q.html", v)
	})
}

func handlerRoot(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if req.URL.Path == "/" {
			b, err := fs.ReadFile(fss, nowPath)
			if err != nil {
				return err
			}

			b, err = parseNow(bytes.NewReader(b), time.Now())
			if err != nil {
				return err
			}

			v := &HTMLVars{Zine: z, Title: "now"}
			if err := mdwn.Convert(b, v); err != nil {
				return err
			}

			if err := tpl.ExecuteTemplate(rw, "simple.html", v); err != nil {
				return err
			}
		}

		f := strings.TrimSuffix(strings.TrimPrefix(req.URL.Path, "/"), "/") + ".md"
		a, err := z.LoadArticle(z.DB, f)
		if err != nil {
			return fmt.Errorf("LoadArticle: %w", err)
		}

		v := articleVars{
			HTMLVars:     &HTMLVars{Zine: z},
			ArticleTitle: a.Title,
			Filename:     f,
		}

		b, err := z.Render(a)
		if err != nil {
			return fmt.Errorf("Render: %w", err)
		}
		buf := bytes.NewBuffer(b)
		if _, err := io.Copy(v, buf); err != nil {
			return err
		}
		if err := tpl.ExecuteTemplate(rw, "article.html", v); err != nil {
			return err
		}
		return nil
	})
}

func handlerToggle(z *notes.Zine, fss fs.FS, mdwn goldmark.Markdown, nowPath string) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		if err := req.ParseForm(); err != nil {
			return err
		}

		c := req.Form.Get("chunk")
		if c == "" {
			rw.WriteHeader(400)
			fmt.Fprint(rw, "missing chunk parameter")
			return nil
		}

		b, err := fs.ReadFile(fss, nowPath)
		if err != nil {
			return err
		}

		b, err = toggleNow(bytes.NewReader(b), time.Now(), c)
		if err != nil {
			return err
		}

		if err := lmfs.WriteFile(fss, nowPath, b, 0); err != nil {
			return err
		}

		rw.Header().Add("Location", "/")
		rw.WriteHeader(303)

		b, err = parseNow(bytes.NewReader(b), time.Now())
		if err != nil {
			return err
		}

		v := &HTMLVars{Zine: z, Title: "now"}
		if err := mdwn.Convert(b, v); err != nil {
			return err
		}

		return tpl.ExecuteTemplate(rw, "simple.html", v)

	})
}

func handlerUpdate(z *notes.Zine, fss fs.FS) http.Handler {
	return lmhttp.HandlerFunc(func(rw http.ResponseWriter, req *http.Request) error {
		switch req.Method {
		case "GET":
			f := req.URL.Query().Get("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			b, err := fs.ReadFile(fss, f)
			if err != nil {
				return err
			}

			v := updateVars{
				HTMLVars: &HTMLVars{Zine: z},
				File:     f,
				Content:  string(b),
			}

			return tpl.ExecuteTemplate(rw, "update.html", v)
		case "POST":
			f := req.FormValue("file")
			if f == "" {
				return errors.New("file parameter required")
			}

			b := req.FormValue("value")
			if b == "" {
				return errors.New("value parameter required")
			}

			b = strings.ReplaceAll(b, "\r", "") // unix files only!
			if err := lmfs.WriteFile(fss, f, []byte(b), 0); err != nil {
				return err
			}
			rw.Header().Add("Location", "/"+strings.TrimSuffix(f, ".md"))
			rw.WriteHeader(303)
			fmt.Fprint(rw, "Successfully updated")
			return nil
		}
		return errors.New("invalid method for /update")
	})
}

func server(fss fs.FS, z *notes.Zine, generation *chan bool) (http.Handler, error) {
	mux := http.NewServeMux()

	mdwn := goldmark.New(
		goldmark.WithParserOptions(
			parser.WithAutoHeadingID(),
		),
		goldmark.WithRendererOptions(
			html.WithUnsafe(),
		),
		goldmark.WithExtensions(
			extension.Strikethrough,
			extension.Table,
		),
	)

	const nowPath = "now.md"

	mux.Handle("/favicon", handlerFavicon())

	mux.Handle("/version/", selfupdate.Handler)

	mux.Handle("/", handlerRoot(z, fss, mdwn, nowPath))

	mux.Handle("/list", handlerList(z, fss, mdwn))

	mux.Handle("/q", handlerQ(z, mdwn))

	mux.Handle("/sup", handlerSup(z))

	mux.Handle("/update", handlerUpdate(z, fss))

	mux.Handle("/toggle", handlerToggle(z, fss, mdwn, nowPath))

	mux.Handle("/add-item", handlerAddItem(z, fss, mdwn, nowPath))

	return autoReload(mux, generation), nil
}
