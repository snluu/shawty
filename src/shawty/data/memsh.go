package data

import (
	"errors"
	"shawty/utils"
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

func (ms *MemSh) Create(r, url, creatorIP string) (*Shawty, error) {
	if r == "" {
		r = utils.ToSafeBase(uint64(ms.random.Byte()) % utils.BaseLen)
	}
	sh := &Shawty{
		ID:        ms.nextID,
		Rand:      r,
		Hits:      0,
		Url:       url,
		CreatorIP: creatorIP,
		CreatedOn: time.Now()}
	ms.data = append(ms.data, sh)
	ms.nextID++
	return sh, nil
}

func (ms *MemSh) GetOrCreate(url, creatorIP string) (*Shawty, error) {
	sh, err := ms.GetByUrl(url)
	if err != nil {
		sh, err = ms.Create("", url, creatorIP)
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

func (ms *MemSh) NumLinks(creatorIP string, t time.Time) (uint32, error) {
	timestamp := t.Unix()
	n := uint32(0)
	for _, s := range ms.data {
		if s.CreatorIP == creatorIP && s.CreatedOn.Unix() >= timestamp {
			n++
		}
	}
	return n, nil
}

func (this *MemSh) Close() error {
	return nil
}
