package db

import (
	"fmt"
	"log"
	"sort"
	"strings"
)

const (
	SortOrderAscending  = "asc"
	SortOrderDescending = "desc"
)

// Create a new chirp and save it in the database
func (db *DB) CreateChirp(body string, authorId int) (Chirp, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return Chirp{}, err
	}

	newId := len(dbstruct.Chirps) + 1

	chirp := Chirp{
		Id:       newId,
		Body:     body,
		AuthorId: authorId,
	}

	dbstruct.Chirps[newId] = chirp

	err = db.persist(dbstruct)
	if err != nil {
		return Chirp{}, err
	}

	return chirp, nil
}

// Helper to sort chirps either ASCending or DESCending
func sortChirps(chirps *[]Chirp, order string) {
	order = strings.ToLower(order)
	if order != SortOrderAscending && order != SortOrderDescending {
		return
	}

	var sortFn func(int, int) bool
	asc := func(i, j int) bool {
		return (*chirps)[j].Id > (*chirps)[i].Id
	}
	desc := func(i, j int) bool {
		return (*chirps)[j].Id < (*chirps)[i].Id
	}

	if order == "asc" {
		sortFn = asc
	} else {
		sortFn = desc
	}

	sort.Slice(*chirps, sortFn)
}

// Retrieves all chirps from the database
func (db *DB) GetChirps(order string) ([]Chirp, error) {
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

	sortChirps(&sortedChirps, order)
	log.Printf("[GetChirps]\tOrder: %v\n%v", order, sortedChirps)

	return sortedChirps, nil
}

func (db *DB) GetChirpsByAuthorId(authorId int, order string) ([]Chirp, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return []Chirp{}, err
	}

	sortedChirpsByAuthor := []Chirp{}

	for _, c := range dbstruct.Chirps {
		if c.AuthorId == authorId {
			sortedChirpsByAuthor = append(sortedChirpsByAuthor, c)
		}
	}

	sortChirps(&sortedChirpsByAuthor, order)
	log.Printf("[GetChirps]\tOrder: %v\n%v", order, sortedChirpsByAuthor)

	return sortedChirpsByAuthor, nil
}

// Retreives a single chirp from the database
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

func (db *DB) DeleteChirp(id int, authorId int) error {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return err
	}

	chirp, ok := dbstruct.Chirps[id]
	if !ok {
		return fmt.Errorf("No chirp found with id: %v", id)
	}

	if chirp.AuthorId != authorId {
		return fmt.Errorf("You are not allowed to delete this chirp")
	}

	delete(dbstruct.Chirps, id)

	err = db.persist(dbstruct)
	if err != nil {
		return err
	}

	return nil
}
