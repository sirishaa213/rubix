package storage

import (
	"github.com/EnsurityTechnologies/adapter"
	"github.com/EnsurityTechnologies/config"
	"github.com/EnsurityTechnologies/uuid"
)

type StorageDB struct {
	ad *adapter.Adapter
}

func NewStorageDB(cfg *config.Config) (*StorageDB, error) {
	ad, err := adapter.NewAdapter(cfg)
	if err != nil {
		return nil, err
	}
	s := &StorageDB{
		ad: ad,
	}
	return s, nil
}

// Init will initialize storage
func (s *StorageDB) Init(storageName string, value interface{}) error {
	return s.ad.InitTable(storageName, value)
}

// Write will write into storage
func (s *StorageDB) Write(storageName string, value interface{}) error {
	return s.ad.Create(storageName, value)
}

// Update will update the storage
func (s *StorageDB) Update(stroageName string, value interface{}, querryString string, querryVaule ...interface{}) error {
	return s.ad.UpdateNew(uuid.Nil, stroageName, querryString, value, querryVaule...)
}

// Delete will delet the data from the storage
func (s *StorageDB) Delete(stroageName string, value interface{}, querryString string, querryVaule ...interface{}) error {
	return s.ad.DeleteNew(uuid.Nil, stroageName, querryString, value, querryVaule...)
}

// Read will read from the storage
func (s *StorageDB) Read(stroageName string, value interface{}, querryString string, querryVaule ...interface{}) error {
	return s.ad.FindNew(uuid.Nil, stroageName, querryString, value, querryVaule...)
}

// Close will close the stroage BD
func (s *StorageDB) Close() error {
	return s.ad.GetDB().Close()
}
