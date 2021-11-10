package main

import (
	"github.com/zalando/go-keyring"
)

const (
	keyRingAppID = "synodsm" // used as ID in the keyring / keychain
)

func saveCredentials(user string, password string) error {

	// try to get credentials, if exist, delete and make new
	// is not exist create new credetials

	_, err := keyring.Get(keyRingAppID, user)
	if err != nil {

		err = keyring.Set(keyRingAppID, user, password)
		if err != nil {
			return err
		}

	} else {

		err = keyring.Delete(keyRingAppID, user)
		if err != nil {
			return err
		}

		err = keyring.Set(keyRingAppID, user, password)
		if err != nil {
			return err
		}
	}

	return nil

}

func getPassword(user string) (string, error) {

	secret, err := keyring.Get(keyRingAppID, user)

	return secret, err

}
