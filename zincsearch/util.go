package zincsearch

import (
	"math/rand"
)

func generateRandomProductName() string {
	var prefixes []string = []string{
		"Cuaderno",
		"Lapicera",
		"Diccionario",
		"Varita",
		"Jornadas",
	}
	var posfixes []string = []string{
		"De Luxe",
		"Super reforzado/a",
		"Lumpurus",
		"De Harray potus",
		"De Disciplina",
	}
	prefixIndex := rand.Intn(len(prefixes))
	posfixIndex := rand.Intn(len(posfixes))
	return prefixes[prefixIndex] + " " + posfixes[posfixIndex]
}
