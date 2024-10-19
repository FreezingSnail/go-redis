package internal

import (
	"errors"
	"fmt"
	"sync"

	tsLog "kvStore/internal/transactionlog"
)

var logger tsLog.TransactionLogger

var store = struct {
	sync.RWMutex
	m map[string]string
}{m: make(map[string]string)}

func Put(key string, value string) error {
	store.Lock()
	store.m[key] = value
	store.Unlock()
	return nil
}

var ErrorNoSuchKey = errors.New("key not found")

func Get(key string) (string, error) {
	store.RLock()
	val, ok := store.m[key]
	store.RUnlock()

	if !ok {
		return "", ErrorNoSuchKey
	}

	return val, nil
}

func Delete(key string) error {
	store.Lock()
	delete(store.m, key)
	store.Unlock()
	return nil
}

func InitializeTransactionLog() error {
	logger, err := tsLog.NewFileTransactionLogger("transaction.log")
	if err != nil {
		return fmt.Errorf("failed to create event logger: %w", err)
	}

	events, errors := logger.ReadEvents()
	e, ok := tsLog.Event{}, true

	for ok && err == nil {
		select {
		case err, ok = <-errors:
		case e, ok = <-events:
			switch e.EventType {
			case tsLog.EventDelete:
				err = Delete(e.Key)
			case tsLog.EventPut:
				err = Put(e.Key, e.Value)
			}
		}
	}

	logger.Run()
	return err
}
