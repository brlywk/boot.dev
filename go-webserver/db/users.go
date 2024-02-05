package db

import (
	"fmt"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// Create a new user.
// The users password will be hashed before saving
func (db *DB) CreateUser(email string, pwd string) (User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbstruct.Users) + 1
	password, err := bcrypt.GenerateFromPassword([]byte(pwd), -1)
	if err != nil {
		return User{}, err
	}

	user := User{
		Id:       newId,
		Email:    email,
		Password: string(password),
	}

	dbstruct.Users[newId] = user

	err = db.persist(dbstruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}

// Looks up a user by email, and if the user exists, tries to verify
// the users password.
// Returns the user if it exists and is verified, and an error otherwise
func (db *DB) VerifyUser(email string, password string) (User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	user := User{}

	user, err := db.GetUserByEmail(email)
	if err != nil {
		return user, fmt.Errorf("No user exists with email '%v'", email)
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return user, fmt.Errorf("Password incorrect")
	}

	return user, nil
}

func (db *DB) GetUserById(id int) (User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbstruct.Users[id]
	if !ok {
		return user, fmt.Errorf("No user exists with id '%v'", id)
	}

	return user, nil
}

func (db *DB) GetUserByEmail(email string) (User, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user := User{}

	for _, u := range dbstruct.Users {
		if strings.ToLower(u.Email) == strings.ToLower(email) {
			return u, nil
		}
	}

	return user, fmt.Errorf("No user found with email %v", email)
}

func (db *DB) UpdateUser(u User) (User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	user, ok := dbstruct.Users[u.Id]
	if !ok {
		return user, fmt.Errorf("no user exists with email '%v'", u.Email)
	}

	// check if password has changed as we don't want to encrypt the encrypted password
	// and accidentally overwrite the old password with that
	samePwd := user.Password == u.Password

	user.Email = u.Email
	user.IsChirpyRed = u.IsChirpyRed

	if !samePwd {
		password, err := bcrypt.GenerateFromPassword([]byte(u.Password), -1)
		if err != nil {
			return User{}, err
		}

		user.Password = string(password)
	}

	dbstruct.Users[u.Id] = user

	err = db.persist(dbstruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
