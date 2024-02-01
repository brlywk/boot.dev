package db

import (
	"encoding/json"
	"log"
	"os"
	"sync"
	"time"
)

// ----- Types -----------------------------------

type DB struct {
	path  string
	mutex *sync.RWMutex
}

type Chirp struct {
	Id       int    `json:"id"`
	Body     string `json:"body"`
	AuthorId int    `json:"author_id"`
}

type User struct {
	Id          int    `json:"id"`
	Email       string `json:"email"`
	Password    string `json:"password"`
	IsChirpyRed bool   `json:"is_chirpy_red"`
}

type DBStructure struct {
	Chirps        map[int]Chirp `json:"chirps"`
	Users         map[int]User  `json:"users"`
	RevokedTokens map[string]time.Time
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
		Chirps:        map[int]Chirp{},
		Users:         map[int]User{},
		RevokedTokens: map[string]time.Time{},
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
