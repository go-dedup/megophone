package megophone

import "testing"

func Test(t *testing.T) {
	Metaphone("Michael")
	Metaphone("McHugh")
	Metaphone("Chianti")
	Metaphone("Caesar")
	Metaphone("Czerny")
	Metaphone("Bach")
	Metaphone("focaccia")
	Metaphone("accident")
	Metaphone("bacci")
	Metaphone("Mac Gregor")
}
