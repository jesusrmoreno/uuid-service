package impl

import (
	"fmt"

	"github.com/xyproto/simplebolt"
)

// BoltStore ...
type BoltStore struct {
	db       *simplebolt.Database
	stores   map[string]*simplebolt.List
	dbStores *simplebolt.List
}

// NewBoltStore ...
func NewBoltStore(path string) (*BoltStore, error) {
	db, err := simplebolt.New(path)
	if err != nil {
		return nil, err
	}

	spacesDB, err := simplebolt.NewList(db, "___namespaces___")
	if err != nil {
		return nil, err
	}

	b := &BoltStore{
		db:       db,
		stores:   map[string]*simplebolt.List{},
		dbStores: spacesDB,
	}

	spaces, err := spacesDB.GetAll()
	if err != nil {
		return nil, err
	}
	for _, ns := range spaces {
		fmt.Println(ns)
		l, err := simplebolt.NewList(db, ns)
		if err != nil {
			return nil, err
		}

		b.stores[ns] = l
	}
	// b.dbStores.ForEach(func(ns []byte, ignore []byte) {
	// 	fmt.Println(string(ignore))
	// })

	return b, nil
}

// Store ...
func (b BoltStore) Store(ns, id string) (bool, error) {
	// If the store exists then we just add to it
	if b.stores[ns] != nil {
		if err := b.stores[ns].Add(id); err != nil {
			return false, err
		}
		return true, nil
	}

	// Otherwise we create it and then add to it.
	if err := b.dbStores.Add(ns); err != nil {
		return false, err
	}
	l, err := simplebolt.NewList(b.db, ns)
	if err != nil {
		return false, err
	}
	b.stores[ns] = l
	if err := b.stores[ns].Add(id); err != nil {
		return false, err
	}
	return true, nil
}

// Contains ...
func (b BoltStore) Contains(ns, id string) (bool, error) {
	if b.stores[ns] != nil {
		keys, err := b.stores[ns].GetAll()
		if err != nil {
			return false, err
		}
		for _, key := range keys {
			if id == key {
				return true, nil
			}
		}
	}
	return false, nil
}

// All ...
func (b BoltStore) All(ns string) ([]string, error) {
	if b.stores[ns] != nil {
		ids, err := b.stores[ns].GetAll()
		if err != nil {
			return nil, err
		}

		return ids, nil
	}
	return []string{}, nil
}
