package main

import (
	"errors"
)

var store = make(map[string]string)

func Put(key string, value string) error {
	store[key] = value
	return nil
}

var ErrorNoSuchKey = errors.New("key not found")

func Get(key string) (string, error) {
	if val, ok := store[key]; ok {
		return val, nil
	}
	return "", ErrorNoSuchKey
}

func Delete(key string) error {
	delete(store, key)
	return nil
}
