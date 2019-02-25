package utils

import "log"

// Writer checks and validates the response
func Writer(n int, err error) {
	if err != nil {
		log.Fatal("Error Occur while writing response: ", err)
	}
}
