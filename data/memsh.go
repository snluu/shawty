package data

import (
	"errors"
	"go.3fps.com/shawty/utils"
	"time"
)

// NewMemSh creates and initialize a new MemoryShawties instance
func NewMemSh(r utils.Rand) *MemSh {
	return &MemSh{
		1,
		make([]*Shawty, 0, 20),
		r,
	}
}

// MemSh implements the Shawties interface and stores all the 
// Shawty instances in memory. Do not use this in production or your data will be lost.
type MemSh struct {
	nextID uint64
	data   []*Shawty
	random utils.Rand
}

func (ms *MemSh) GetByID(id uint64, r string) (*Shawty, error) {
	for _, s := range ms.data {
		if s.ID == id && s.Rand == r {
			return s, nil
		}
	}

	return nil, errors.New("Cannot find Shawty")
}

func (ms *MemSh) GetByUrl(url string) (*Shawty, error) {
	for _, s := range ms.data {
		if s.Url == url {
			return s, nil
		}
	}

	return nil, errors.New("Cannot find Shawty")
}

func (ms *MemSh) Create(r, url string) (*Shawty, error) {
	if r == "" {
		r = utils.ToSafeBase(uint64(ms.random.Byte()) % utils.BaseLen)
	}
	sh := &Shawty{ms.nextID, r, 0, url, time.Now()}
	ms.data = append(ms.data, sh)
	ms.nextID++
	return sh, nil
}

func (ms *MemSh) GetOrCreate(url string) (*Shawty, error) {
	sh, err := ms.GetByUrl(url)
	if err != nil {
		sh, err = ms.Create("", url)
	}
	return sh, err
}

func (ms *MemSh) IncHits(id uint64) error {
	for _, s := range ms.data {
		if s.ID == id {
			s.Hits++
			return nil
		}
	}

	return errors.New("Cannot find Shawty")
}
