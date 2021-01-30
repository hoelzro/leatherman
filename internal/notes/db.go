package notes

import (
	"database/sql"
	"sync"

	"github.com/jmoiron/sqlx"
	_ "modernc.org/sqlite"
)

const Schema = `
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
`

type DB struct {
	*sqlx.DB
	insertTags *sql.Stmt

	stmtCacheMu *sync.Mutex
	stmtCache   map[string]*sqlx.Stmt
}

func NewDB() (*DB, error) {
	var (
		dbh *sqlx.DB
		err error
	)

	dbh, err = sqlx.Open("sqlite", "file:.posts.db?mode=memory&_sync=OFF&_journal=OFF&_vacuum=0")
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

	if _, err := dbh.Exec(Schema); err != nil {
		return nil, err
	}

	var success bool
	defer func() {
		if !success {
			dbh.Close()
		}
	}()

	d := &DB{DB: dbh, stmtCache: map[string]*sqlx.Stmt{}, stmtCacheMu: &sync.Mutex{}}
	d.insertTags, err = d.Prepare(`INSERT INTO article_tag (id, tag) VALUES (?, ?)`)
	if err != nil {
		return nil, err
	}

	success = true
	return d, nil
}

func (d *DB) PrepareCached(sql string) (*sqlx.Stmt, error) {
	d.stmtCacheMu.Lock()
	defer d.stmtCacheMu.Unlock()
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

func (d *DB) InsertArticle(a Article) error {
	stmt, err := d.PrepareCached(`INSERT INTO articles (
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

func (d *DB) LoadArticle(name string) (Article, error) {
	stmt, err := d.PrepareCached(`
	SELECT rowid, title, url, filename, reviewed_on, review_by, body
	FROM articles
	WHERE filename = ?
	`)
	if err != nil {
		return Article{}, err
	}

	var ret struct {
		Article
		RowID int
	}

	if err := stmt.Get(&ret, name); err != nil {
		return Article{}, err
	}

	tagsStmt, err := d.PrepareCached(`SELECT tag FROM article_tag WHERE id = ?`)
	if err != nil {
		return Article{}, err
	}

	if err := tagsStmt.Select(&ret.Tags, ret.RowID); err != nil {
		return Article{}, err
	}

	return ret.Article, nil
}

func (d *DB) DeleteArticle(name string) error {
	tagStmt, err := d.PrepareCached(`DELETE FROM article_tag WHERE id IN (SELECT rowid FROM articles WHERE filename = ?)`)
	if err != nil {
		return err
	}

	if _, err := tagStmt.Exec(name); err != nil {
		return err
	}

	stmt, err := d.PrepareCached(`DELETE FROM articles WHERE filename = ?`)
	if err != nil {
		return err
	}

	if _, err := stmt.Exec(name); err != nil {
		return err
	}

	return nil
}

func (d *DB) ReplaceArticle(a Article) (err error) {
	if err := d.DeleteArticle(a.Filename); err != nil {
		return err
	}
	if err := d.InsertArticle(a); err != nil {
		return err
	}

	return nil
}
