package main

import "database/sql"
import "os"
import "time"
import _ "code.google.com/p/go-mysql-driver/mysql"
import "math/rand"
import "errors"

type Shawty struct {
	ID        uint64
	Rand      string
	Hits      uint64
	Url       string
	CreatedOn time.Time
}

// Shawties is the proxy to manage Shawty records in the database
type Shawties struct {
	db *sql.DB
}

// ShortID constructs a short ID
func ShortID(id uint64, r string) string {
	return r + toSafeBase(id)
}

// FullID deconstructs a full ID
func FullID(str string) (uint64, string, error) {
	if len(str) < 2 {
		return 0, "", errors.New("Cannot deconstruct " + str)
	}
	return toDec(str[1:]), str[:1], nil
}

// NewShawties yields a new Shawties instance with pre-open database connection
func NewShawties() (*Shawties, error) {
	var sh = new(Shawties)
	var err = sh.Open("mysql", os.Getenv("SHAWTY_DB"))
	return sh, err
}

// Open opens a new database connection
func (sh *Shawties) Open(driverName, dataSrc string) error {
	var db, err = sql.Open(driverName, dataSrc)
	sh.db = db
	return err
}

// Close closes the database connection
func (sh *Shawties) Close() error {
	return sh.db.Close()
}

// get fetches a Shawty instance based on the URL
func (sh *Shawties) getByUrl(url string) (*Shawty, error) {
	var stmt, err = sh.db.Prepare("select `ID`, `Rand`, `Hits`, `Url`, `CreatedOn` from `shawties` where `Url` = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var row = stmt.QueryRow(url)
	var s = new(Shawty)
	var t int64
	err = row.Scan(&s.ID, &s.Rand, &s.Hits, &s.Url, &t)
	if err != nil {
		return nil, err
	}
	s.CreatedOn = time.Unix(t, 0)
	return s, nil
}

// get fetches a Shawty instance by the ID and a random 
func (sh *Shawties) GetByID(id uint64, r string) (*Shawty, error) {
	var stmt, err = sh.db.Prepare("select `ID`, `Rand`, `Hits`, `Url`, `CreatedOn` from `shawties` where `ID` = ? and `Rand` = ? limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var row = stmt.QueryRow(id, r)
	var s = new(Shawty)
	var t int64
	err = row.Scan(&s.ID, &s.Rand, &s.Hits, &s.Url, &t)
	if err != nil {
		return nil, err
	}
	s.CreatedOn = time.Unix(t, 0)
	return s, nil
}

// create inserts a new record into the database
func (sh *Shawties) create(r string, url string) (*Shawty, error) {
	if r == "" {
		r = toSafeBase(uint64(rand.Int63n(int64(baseLen))))
	}

	var stmt, err = sh.db.Prepare("insert ignore into `shawties`(`Rand`, `Hits`, `Url`, `CreatedOn`) values(?, 0, ?, ?)")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()
	var now = time.Now()
	result, err := stmt.Exec(r, url, now.Unix())
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	var s = &Shawty{uint64(id), r, 0, url, now}
	return s, nil
}

func (sh *Shawties) GetOrCreate(url string) (*Shawty, error) {
	var s, err = sh.getByUrl(url)
	if err != nil {
		// try creating
		s, err = sh.create("", url)
		if err != nil {
			return nil, err
		}
	}

	return s, err
}

func (sh *Shawties) IncHits(id uint64) error {
	stmt, err := sh.db.Prepare("update `shawties` set `Hits` = `Hits` + 1 where `ID` = ?") 
	if err != nil {
		return err
	}
	defer stmt.Close();
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
