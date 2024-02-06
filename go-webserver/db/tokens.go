package db

import "time"

func (db *DB) IsTokenRevoked(token string) (bool, error) {
	db.mutex.RLock()
	defer db.mutex.RUnlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return false, err
	}

	_, revoked := dbstruct.RevokedTokens[token]

	return revoked, nil
}

func (db *DB) RevokeToken(token string) (time.Time, error) {
	db.mutex.Lock()
	defer db.mutex.Unlock()

	dbstruct, err := db.loadDB()
	if err != nil {
		return time.Time{}, err
	}

	now := time.Now()
	dbstruct.RevokedTokens[token] = now

	err = db.persist(dbstruct)
	if err != nil {
		return time.Time{}, err
	}

	return now, nil
}
