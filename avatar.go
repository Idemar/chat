package main

import "errors"

// ErrNoAvatarURL er feilen som returneres når
// Avatar-forekomsten ikke kan oppgi en avatar-URL.
var ErrNoAvatarURL = errors.New("chat: Unable to get an avatar URL")

// Avatar representerer typer som kan representere
// brukerprofilbilder.
type Avatar interface {
	// GetAvatarURL henter avatar-URL-en for den angitte klienten,
	// eller returnerer en feilmelding hvis noe går galt.
	// ErrNoAvatarURL returneres hvis objektet ikke klarer å hente
	// en URL for den angitte klienten.
	GetAvatarURL(c *client) (string, error)
}
