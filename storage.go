package eid

type Storage interface {
	Load() (map[string]uint64, error)
	Save(map[string]uint64) error
}
