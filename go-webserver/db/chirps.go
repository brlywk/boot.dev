package db

import (
	"fmt"
	"log"
	"sort"
)

func (db *DB) CreateChirp(body string) (Chirp, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	log.Println("\nCreating Chirp")

	dbstruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbstruct.Chirps) + 1

	chirp := Chirp{
		Id:   newId,
		Body: body,
	}

	dbstruct.Chirps[newId] = chirp

	err = db.persist(dbstruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

func (db *DB) GetChirps() ([]Chirp, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	sortedChirps := []Chirp{}

	for _, c := range dbstruct.Chirps {
		sortedChirps = append(sortedChirps, c)
	}

	sort.Slice(sortedChirps, func(i, j int) bool {
		return j > i
	})

	return sortedChirps, nil
}

func (db *DB) GetChirp(id int) (Chirp, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	chirp := Chirp{}

	dbstruct, err := db.loadDB()
	if err != nil {
		return chirp, err
	}

	chirp, ok := dbstruct.Chirps[id]
	if !ok {
		return chirp, fmt.Errorf("No chirp found with id: %v", id)
	}

	return chirp, nil
}
