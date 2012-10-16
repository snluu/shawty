package data

// Shawties defines the data access interface to manage Shawty instances
type Shawties interface {
	// GetByID fetches a Shawty instance by the ID and the random digit
	GetByID(id uint64, r string) (*Shawty, error)
	// GetByUrl featches a Shawty instance based on the URL
	GetByUrl(url string) (*Shawty, error)
	// Create makes and stores a new Shawty instance. r is the random digit
	Create(r, url string) (*Shawty, error)
	// GetOrCreate tries to get a Shawty instance based on the URL and creates one if it does not exist
	GetOrCreate(url string) (*Shawty, error)
	// IncHits increases the hit count for the Shawty with the specific ID
	IncHits(id uint64) error
}
