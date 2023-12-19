package info

import "kantoku/common/data/record"

type Storage struct {
	records  record.Storage
	settings Settings
}

func NewStorage(records record.Storage, settings Settings) *Storage {
	return &Storage{
		records:  records,
		settings: settings,
	}
}

func (s *Storage) Get(id string) Info {
	return Info{
		storage: s,
		id:      id,
	}
}

func (s *Storage) Records() record.Storage {
	return s.records
}

func (s *Storage) IdProperty() string {
	return s.settings.IdProperty
}
