package zine

import (
	"database/sql"
	"os"

	"github.com/frioux/leatherman/internal/notes"
	"github.com/jmoiron/sqlx"
)

type db struct {
	*sqlx.DB
	insertTags *sql.Stmt

	// stmtCache is not safe for concurrent access.
	stmtCache map[string]*sqlx.Stmt
}

func newDB() (*db, error) {
	var (
		dbh *sqlx.DB
		err error
	)

	if err := os.Remove(".posts.db"); err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	dbh, err = sqlx.Open("sqlite", "file:.posts.db?_sync=OFF&_journal=OFF&_vacuum=0")
	if err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		PRAGMA journal_mode = OFF
	`); err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		PRAGMA synchronous = OFF
	`); err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		PRAGMA auto_vacuum = OFF
	`); err != nil {
		return nil, err
	}

	if _, err := dbh.Exec(`
		CREATE TABLE articles (
			title,
			url,
			filename,
			reviewed_on NULLABLE,
			review_by NULLABLE,
			body
		);
		CREATE TABLE article_tag ( id, tag );
		CREATE VIEW _ ( id, title, url, filename, body, reviewed_on, review_by, tag) AS
			SELECT a.rowid, title, url, filename, body, reviewed_on, review_by, tag
			FROM articles a
			JOIN article_tag at ON a.rowid = at.id;
	`); err != nil {
		return nil, err
	}

	var success bool
	defer func() {
		if !success {
			dbh.Close()
		}
	}()

	d := &db{DB: dbh, stmtCache: map[string]*sqlx.Stmt{}}
	d.insertTags, err = d.Prepare(`INSERT INTO article_tag (id, tag) VALUES (?, ?)`)
	if err != nil {
		return nil, err
	}

	success = true
	return d, nil
}

func (d *db) prepareCached(sql string) (*sqlx.Stmt, error) {
	if stmt, ok := d.stmtCache[sql]; ok {
		return stmt, nil
	}

	stmt, err := d.Preparex(sql)
	if err != nil {
		return nil, err
	}

	d.stmtCache[sql] = stmt
	return stmt, nil
}

func (d *db) insertArticle(a notes.Article) error {
	stmt, err := d.prepareCached(`INSERT INTO articles (
		title, url, filename, reviewed_on, review_by, body
	) VALUES (?, ?, ?, ?, ?, ?)`)
	if err != nil {
		return err
	}
	r, err := stmt.Exec(a.Title, a.URL, a.Filename, a.ReviewedOn, a.ReviewBy, string(a.Body))
	if err != nil {
		return err
	}

	id, err := r.LastInsertId()
	if err != nil {
		return err
	}
	for _, tag := range a.Tags {
		if _, err := d.insertTags.Exec(id, tag); err != nil {
			return err
		}
	}

	return nil
}
