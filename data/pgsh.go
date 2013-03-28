package data

import (
	"database/sql"
	log "github.com/3fps/log2go"
	"github.com/3fps/shawty/utils"
	_ "github.com/lib/pq"
	"time"
)

// PgSh is the MySQL implementation of Shawties interface
type PgSh struct {
	db     *sql.DB
	random utils.Rand
}

// NewPgSh yields a new PgSh instance with pre-open database connection
func NewPgSh(random utils.Rand, dataSrc string) (*PgSh, error) {
	var sh = new(PgSh)
	sh.random = random
	var err = sh.open("postgres", dataSrc)
	if err != nil {
		log.Error("Fail to create postgresql")
	}
	return sh, err
}

// open opens a new database connection
func (sh *PgSh) open(driverName, dataSrc string) error {
	var db, err = sql.Open(driverName, dataSrc)
	sh.db = db
	return err
}

// Close closes the database connection
func (sh *PgSh) Close() error {
	return sh.db.Close()
}

func (sh *PgSh) GetByUrl(url string) (*Shawty, error) {
	var stmt, err = sh.db.Prepare("select ID, Rand, Hits, Url, CreatedOn from shawties where Url = $1 limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var row = stmt.QueryRow(url)
	var s = new(Shawty)
	var t int64
	err = row.Scan(&s.ID, &s.Rand, &s.Hits, &s.Url, &t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.CreatedOn = time.Unix(t, 0)
	return s, nil
}

func (sh *PgSh) GetByID(id uint64, r string) (*Shawty, error) {
	var stmt, err = sh.db.Prepare("select ID, Rand, Hits, Url, CreatedOn from shawties where ID = $1 and Rand = $2 limit 1")
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	var row = stmt.QueryRow(id, r)
	var s = new(Shawty)
	var t int64
	err = row.Scan(&s.ID, &s.Rand, &s.Hits, &s.Url, &t)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	s.CreatedOn = time.Unix(t, 0)
	return s, nil
}

func (sh *PgSh) Create(r string, url, creatorIP string) (*Shawty, error) {
	if r == "" {
		r = utils.ToSafeBase(uint64(sh.random.Byte()) % utils.BaseLen)
	}
	var stmt, err = sh.db.Prepare("insert into shawties(Rand, Hits, Url, CreatorIP, CreatedOn) values($1, 0, $2, $3, $4) returning ID")
	if err != nil {
		log.Error(err)
		return nil, err
	}
	defer stmt.Close()
	var now = time.Now()
	row := stmt.QueryRow(r, url, creatorIP, now.Unix())

	var id uint64
	err = row.Scan(&id)
	if err != nil {
		log.Error(err)
		return nil, err
	}
	var s = &Shawty{
		ID:        uint64(id),
		Rand:      r,
		Hits:      0,
		Url:       url,
		CreatorIP: creatorIP,
		CreatedOn: now}
	log.Infof("Created link ID %d", id)
	return s, nil
}

func (sh *PgSh) GetOrCreate(url, creatorIP string) (*Shawty, error) {
	var s, err = sh.GetByUrl(url)
	if err != nil {
		// try creating
		s, err = sh.Create("", url, creatorIP)
		if err != nil {
			return nil, err
		}
	}

	return s, err
}

func (sh *PgSh) IncHits(id uint64) error {
	stmt, err := sh.db.Prepare("update shawties set Hits = Hits + 1 where ID = $1")
	if err != nil {
		log.Error(err)
		return err
	}
	defer stmt.Close()
	_, err = stmt.Exec(id)
	if err != nil {
		return err
	}
	return nil
}

func (sh *PgSh) NumLinks(creatorIP string, t time.Time) (uint32, error) {
	var stmt, err = sh.db.Prepare("select count(ID) NumLinks from shawties where CreatorIP = $1 and CreatedOn >= $2")
	if err != nil {
		log.Error(err)
		return uint32(0), err
	}
	defer stmt.Close()

	var row = stmt.QueryRow(creatorIP, t.Unix())
	n := uint32(0)
	err = row.Scan(&n)
	if err != nil {
		log.Error(err)
		return uint32(0), err
	}
	return n, nil
}
