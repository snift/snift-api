package service

import (
    "strings"
    "log"
)

func CalculateProtocolScore(protocol string) (score int, message string) {

	score = -1
	if strings.Compare(protocol,"http") ==0 {
		score = 0
		message = "Website is not Encrypted and hence subjective to Man-in-the-Middle attacks(MITM) and Eavesdropping Attacks."
	} else if strings.Compare(protocol,"https") ==0 {
		score = 5
		message = "From the protocol level, Website is secure."
	} else {
		log.Panic("Protocol Not Found")
		message = "Protocol Not Found"
	}
	return
}