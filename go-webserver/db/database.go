package db

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sort"
	"sync"
)

// ----- Types -----------------------------------

type DB struct {
	path  string
	mutex *sync.RWMutex
}

type Chirp struct {
	Id   int    `json:"id"`
	Body string `json:"body"`
}

type DBStructure struct {
	Chirps map[int]Chirp `json:"chirps"`
}

// ----- Methods ---------------------------------

// Creates a new DB connection
// DB will be saved in file database.json, and created if necessary
func NewDB(path string) (*DB, error) {
	db := DB{
		path:  path,
		mutex: &sync.RWMutex{},
	}

	err := db.spawnDB()
	if err != nil {
		return nil, err
	}

	return &db, nil
}

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

// ----- Helper ----------------------------------

// Spawns the DB file if it doesn't exist
func (db *DB) spawnDB() error {
	// check if file exists, if not create
	if _, err := os.Stat(db.path); err != nil {
		file, err := os.Create(db.path)
		if err != nil {
			return err
		}
		defer file.Close()
	}
	return nil
}

// Read DB file into memory
func (db *DB) loadDB() (DBStructure, error) {
	dbStruct := DBStructure{
		Chirps: map[int]Chirp{},
	}

	fileData, err := os.ReadFile(db.path)
	if err != nil {
		log.Println("Error reading content of database file")
		return dbStruct, err
	}

	// if the database did not exist, we need to to provide
	// some minimal valid json to the unmarshaller...
	if len(fileData) < 1 {
		log.Println("Empty file found")
		fileData = []byte(string("{}"))
	}

	err = json.Unmarshal(fileData, &dbStruct)
	if err != nil {
		return dbStruct, err
	}

	return dbStruct, nil
}

// Persist database
func (db *DB) persist(dbstructure DBStructure) error {
	jsonData, err := json.Marshal(dbstructure)
	if err != nil {
		return err
	}

	err = os.WriteFile(db.path, jsonData, 0666)
	if err != nil {
		return err
	}

	return nil
}
