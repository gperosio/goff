package main

import (
	"strconv"
	"sync"
	"time"
)

type StoreItem struct {
	Value  string
	Expiry time.Time
	Mutex  sync.Mutex
}

var store = make(map[string]*StoreItem)

// PerformPong response pack with "PONG", or optionally a passed in argument
func PerformPong(args []string) string {
	if len(args) > 0 {
		return stringMsg(args[0])
	}
	return stringMsg("PONG")
}

// PerformSet stores a value with an expiry in the database
func PerformSet(args []string) string {
	if len(args) < 2 {
		return errorMsg("invalid syntax provided to 'SET'")
	}
	key := &args[0]
	val := &args[1]

	store[*key] = &StoreItem{
		Value: *val,
	}
	return stringMsg("OK")
}

func PerformSetEx(args []string) string {
	if len(args) < 3 {
		return errorMsg("invalid syntax provided to 'SETEX'")
	}
	key := &args[0]
	val := &args[1]
	dur := &args[2]
	var exp time.Time

	if len(args) > 2 {
		// We have options to parse. Let's parse them!
		expSeconds, err := strconv.Atoi(*dur)
		if err != nil {
			return errorMsg("invalid format provided to 'SETEX'")
		}
		exp = time.Now().Add(time.Second * time.Duration(expSeconds))
	}

	// Default to an hour
	if exp == (time.Time{}) {
		exp = time.Now().Add(time.Hour)
	}

	store[*key] = &StoreItem{
		Value:  *val,
		Expiry: exp,
	}
	return stringMsg("OK")
}

// PerformGet retrieves a value from the database, if it exists and is not expired. If it is expired, it will be deleted
func PerformGet(args []string) string {
	if len(args) == 0 {
		return errorMsg("no value provided to 'GET'")
	}

	item := store[args[0]]

	// Check if it's null
	if item == nil {
		return nilBulkStringMsg()
	}

	// Item exists - enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	// Check the expiry
	now := time.Now()
	if !item.Expiry.IsZero() && item.Expiry.Before(now) {
		store[args[0]] = nil
		return nilBulkStringMsg()
	}

	return stringMsg(store[args[0]].Value)
}

// PerformDel deletes a value from the database, if it exists
func PerformDel(args []string) string {
	if len(args) == 0 {
		return errorMsg("no value provided to 'DEL'")
	}

	item := store[args[0]]

	// Check if it's null
	if item == nil {
		return nilBulkStringMsg()
	}

	// Item exists - enforce mutual exclusion on the expiry operation and retrieval
	item.Mutex.Lock()
	defer item.Mutex.Unlock()

	// Delete the item
	delete(store, args[0])

	return stringMsg("OK")
}
