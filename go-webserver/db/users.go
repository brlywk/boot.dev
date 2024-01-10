package db

func (db *DB) CreateUser(email string) (User, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return User{}, err
	}

	newId := len(dbstruct.Users) + 1

	user := User{
		Id:    newId,
		Email: email,
	}

	dbstruct.Users[newId] = user

	err = db.persist(dbstruct)
	if err != nil {
		return User{}, err
	}

	return user, nil
}
