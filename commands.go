package main

import (
	"strconv"
	"strings"
)

const MAX_LEVEL_RECURSION = 10

// PerformPong response pack with "PONG", or optionally a passed in argument
func PerformPong(args []string) string {
	if len(args) > 0 {
		return stringMsg(args[0])
	}
	return stringMsg("PONG")
}

// PerformGet retrieves a value from the database, if it exists and is not expired. If it is expired, it will be deleted
func PerformGet(args []string) string {
	if len(args) == 0 {
		return errorMsg("no value provided to 'GET'")
	}

	//buscar llave principal
	item := store[args[0]]

	// Check if it's null
	if item == nil {
		return nilBulkStringMsg()
	}

	// Check the expiry
	//now := time.Now()
	//if !item.Expiry.IsZero() && item.Expiry.Before(now) {
	//	store[args[0]] = nil
	//	return nilBulkStringMsg()
	//}
	// Item exists - uso mutex en caso de que este presente Lock
	item.Mutex.RLock()
	defer item.Mutex.RUnlock()

	return jsonMsg(item.Value)
}

func PerformJsonGet(args []string) string {
	if len(args) == 0 {
		return errorMsg("no value provided to 'HGET'")
	}

	if len(args) < 2 {
		return errorMsg("no field provided to 'HGET'")
	}

	//buscar llave principal
	item := store[args[0]]
	parts := strings.Split(args[1], ".")

	// Check if it's null
	if item == nil {
		return nilBulkStringMsg()
	}

	if len(parts) > MAX_LEVEL_RECURSION {
		return errorMsg("depth level is higher, max " + strconv.Itoa(MAX_LEVEL_RECURSION))
	}

	item.Mutex.RLock()
	defer item.Mutex.RUnlock()

	// val := findKey(item, parts)
	val := recurseMap(item.Value, parts)
	if val == nil {
		return nilBulkStringMsg()
	}

	return jsonMsg(val)
}

/*
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
*/
func PerformScan(args []string) string {
	var out []string
	var c int = 0
	for k := range store {
		out = append(out, k)
		c++
	}
	out = append(out, strconv.Itoa(c))
	return arrayMsg(out)
}
