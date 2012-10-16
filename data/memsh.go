package data

import (
	"errors"
	"time"
)

// NewMemSh creates and initialize a new MemoryShawties instance
func NewMemSh() *MemSh {
	return &MemSh{
		0,
		make([]Shawty, 20),
	}
}

// MemSh implements the Shawties interface and stores all the 
// Shawty instances in memory. Do not use this in production or your data will be lost.
type MemSh struct {
	nextID uint64
	data   []Shawty
}

func (ms *MemSh) GetByID(id uint64, r string) (*Shawty, error) {
	for _, s := range ms.data {
		if s.ID == id && s.Rand == r {
			return &s, nil
		}
	}

	return nil, errors.New("Cannot find Shawty")
}

func (ms *MemSh) GetByUrl(url string) (*Shawty, error) {
	for _, s := range ms.data {
		if s.Url == url {
			return &s, nil
		}
	}

	return nil, errors.New("Cannot find Shawty")
}

func (ms *MemSh) Create(r, url string) (*Shawty, error) {
	sh := &Shawty{ms.nextID, r, 0, url, time.Now()}
	ms.nextID++
	return sh, nil
}

func (ms *MemSh) IncHits(id uint64) error {
	for _, s := range ms.data {
		if s.ID == id {
			s.Hits++
			return
		}
	}

	return errors.New("Cannot find Shawty")
}