package data

import (
	_ "code.google.com/p/go-mysql-driver/mysql"
	"database/sql"
	"../utils"
	"go.3fps.com/utils/log"
	"time"
)

// MySh is the MySQL implementation of Shawties interface
type MySh struct {
	db     *sql.DB
	random utils.Rand
}

// NewMySh yields a new MySh instance with pre-open database connection
func NewMySh(random utils.Rand, dataSrc string) (*MySh, error) {
	var sh = new(MySh)
	sh.random = random
	var err = sh.open("mysql", dataSrc)
	return sh, err
}

// open opens a new database connection
func (sh *MySh) open(driverName, dataSrc string) error {
	var db, err = sql.Open(driverName, dataSrc)
	sh.db = db
	return err
}

// Close closes the database connection
func (sh *MySh) Close() error {
	return sh.db.Close()
}

func (sh *MySh) GetByUrl(url string) (*Shawty, error) {
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

func (sh *MySh) GetByID(id uint64, r string) (*Shawty, error) {
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

func (sh *MySh) Create(r string, url string) (*Shawty, error) {
	if r == "" {
		r = utils.ToSafeBase(uint64(sh.random.Byte()) % utils.BaseLen)
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
	log.Infof("Created link ID %d", id)
	return s, nil
}

func (sh *MySh) GetOrCreate(url string) (*Shawty, error) {
	var s, err = sh.GetByUrl(url)
	if err != nil {
		// try creating
		s, err = sh.Create("", url)
		if err != nil {
			return nil, err
		}
	}

	return s, err
}

func (sh *MySh) IncHits(id uint64) error {
	stmt, err := sh.db.Prepare("update `shawties` set `Hits` = `Hits` + 1 where `ID` = ?")
	if err != nil {
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}
