package main

import "time"

// message representerer en enkelt melding
type message struct {
	Name      string
	Message   string
	When      time.Time
	AvatarURL string
}
