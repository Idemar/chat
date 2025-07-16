package main

type room struct {
	// videresend er en kanal som inneholder innkommende meldinger
	// som skal videresendes til de andre klientene.
	forward chan []byte

	// join  er en kanal for klienter som ønsker å bli med i rommet.
	join chan *client

	// Leave er en kanal for klienter som ønsker å forlate rommet.
	leave chan *client

	// clients holder alle nåværende klienter i dette rommet.
	clients map[*client]bool
}
