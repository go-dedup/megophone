package megophone

import "testing"

func Test(t *testing.T) {
	Metaphone("michael")
	Metaphone("mchugh")
	Metaphone("chianti")
	Metaphone("caesar")
	Metaphone("czerny")
	Metaphone("bach")
	Metaphone("focaccia")
}
