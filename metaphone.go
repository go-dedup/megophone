package megophone

/*
The original Metaphone algorithm was published in 1990 as an improvement over
the Soundex algorithm. Like Soundex, it was limited to English-only use. The
Metaphone algorithm does not produce phonetic representations of an input word
or name; rather, the output is an intentionally approximate phonetic
representation. The approximate encoding is necessary to account for the way
speakers vary their pronunciations and misspell or otherwise vary words and
names they are trying to spell.
The Double Metaphone phonetic encoding algorithm is the second generation of
the Metaphone algorithm. Its implementation was described in the June 2000
issue of C/C++ Users Journal. It makes a number of fundamental design
improvements over the original Metaphone algorithm.
It is called "Double" because it can return both a primary and a secondary code
for a string; this accounts for some ambiguous cases as well as for multiple
variants of surnames with common ancestry. For example, encoding the name
"Smith" yields a primary code of SM0 and a secondary code of XMT, while the
name "Schmidt" yields a primary code of XMT and a secondary code of SMT--both
have XMT in common.
Double Metaphone tries to account for myriad irregularities in English of
Slavic, Germanic, Celtic, Greek, French, Italian, Spanish, Chinese, and other
origin. Thus it uses a much more complex ruleset for coding than its
predecessor; for example, it tests for approximately 100 different contexts of
the use of the letter C alone.
This script implements the Double Metaphone algorithm (c) 1998, 1999 originally
implemented by Lawrence Philips in C++. It was further modified in C++ by Kevin
Atkinson (http {//aspell.net/metaphone/). It was translated to C by Maurice
Aubrey <maurice@hevanet.com> for use in a Perl extension. A Python version was
created by Andrew Collins on January 12, 2007, using the C source
(http {//www.atomodo.com/code/double-metaphone/metaphone.py/view). It was also
translated to Go by Adele Dewey-Lopez <adele@seed.co> using Atkinson's C++ source.
  Updated 2007-02-14 - Found a typo in the 'gh' section (0.1.1)
  Updated 2007-12-17 - Bugs fixed in 'S', 'Z', and 'J' sections (0.2;
                       Chris Leong)
  Updated 2009-03-05 - Various bug fixes against the reference C++
                       implementation (0.3; Matthew Somerville)
  Updated 2012-07    - Fixed long lines, added more docs, changed names,
                       reformulated as objects, fixed a bug in 'G'
                       (0.4; Duncan McGreggor)
  Updated 2013-06    - Enforced unicode literals (0.5; Ian Beaver)
*/

import "fmt"

type phoneticData struct {
	t               string
	cur             int
	isSlavoGermanic bool
	metaphone1      string
	metaphone2      string
}

func (p *phoneticData) matchesAny(pos int, matches ...string) bool {
	if len(matches) == 0 {
		return true
	}
	// out of bounds
	if p.cur+pos < 0 {
		return false
	}

	for i, str := range matches {
		size := len(matches[i])
		if p.t[p.cur+pos:p.cur+size+pos] == str {
			return true
		}
	}

	return false
}

func (p *phoneticData) add(phoneme ...string) {
	if len(phoneme) > 0 {
		p.metaphone1 += phoneme[0]
		if len(phoneme) > 1 {
			p.metaphone2 += phoneme[1]
		} else {
			p.metaphone2 += phoneme[0]
		}
	}
}

func (p *phoneticData) skip(skipBy int) {
	p.cur += skipBy
}

func (p *phoneticData) isVowel(pos int) bool {
	return p.matchesAny(pos, "a", "e", "i", "o", "u", "y")
}

func (p *phoneticData) b() {
	p.add("p")
	// skip double b
	if p.t[p.cur+1] == 'b' {
		p.skip(1)
	}
}

func (p *phoneticData) ç() {
	p.add("s")
}

func (p *phoneticData) c() {

	if p.cur > 1 && !p.isVowel(-2) && p.matchesAny(-1, "ach") && !p.matchesAny(2, "i") &&
		(!p.matchesAny(2, "e") || p.matchesAny(-2, "acher")) {
		// various germanic
		p.add("k")
		p.skip(1)
	} else if p.cur == 0 && p.matchesAny(0, "caesar") {
		// special case: "caesar"
		p.add("s")
		p.skip(1)
	} else if p.matchesAny(0, "chia") {
		// italian "chianti"
		p.add("k")
		p.skip(1)
	} else if p.matchesAny(0, "ch") {
		// ch
		if p.cur > 0 && p.matchesAny(0, "chae") {
			// find "michael"
			p.add("k")
		} else if p.cur == 0 && !p.matchesAny(0, "chore") &&
			p.matchesAny(1, "harac", "haris", "hor", "hym", "hia", "hem") {
			// greek roots
			p.add("k")
		} else if p.matchesAny(-p.cur, "van ", "von ", "sch") ||
			p.matchesAny(-2, "orches", "archit", "orchid") || p.matchesAny(2, "t", "s") ||
			((p.matchesAny(-1, "a", "e", "o", "u") || p.cur == 0) &&
				p.matchesAny(2, "l", "r", "n", "m", "b", "h", "f", "v", "w", " ")) {
			// germanic greek or otherwise "ch" for "kh" sound
			// "architect" but not "arch", "orchestra" or "orchid"
			// e.g., "watchler", "wechsler", but not "tichner"
			p.add("k")
		} else if p.cur > 0 {
			if p.matchesAny(-p.cur, "mc") {
				// e.g. "McHugh"
				p.add("k")
			} else {
				p.add("x", "k")
			}
		} else {
			p.add("x")
		}
		p.skip(1)
	} else if p.matchesAny(0, "cz") && !p.matchesAny(-2, "wicz") {
		// e.g. "czerny"
		p.add("s", "x")
		p.skip(1)
	} else if p.matchesAny(1, "cia") {
		// e.g. "focaccia"
		p.add("x")
		p.skip(2)
	}
}

func Metaphone(s string) (string, string) {

	// initialize
	var p *phoneticData
	p = &phoneticData{}

	// pad string
	// normalize
	p.t = s + "     "

	if p.matchesAny(0, "gn", "kn", "pn", "wr", "ps") {
		p.skip(2)
	}

	if p.matchesAny(0, "x") {
		p.add("s")
	}

	for p.cur < len(p.t) {
		next := p.t[p.cur]
		//fmt.Println(p.cur, ": ", string(next))
		switch next {
		case 'a', 'e', 'i', 'o', 'u', 'y':
			if p.cur == 0 {
				p.add("a")
			}
		case 'b':
			p.b()
		case 'ç':
			p.ç()
		case 'c':
			p.c()
			// case 'd':
			// 	p.b()
			// case 'f':
			// 	p.b()
			// case 'g':
			// 	p.b()
			// case 'h':
			// 	p.b()
		}
		p.cur++

	}

	fmt.Println("First: ", p.metaphone1, "\tSecond: ", p.metaphone2, "\t Original: ", s)

	return p.metaphone1, p.metaphone2
}
