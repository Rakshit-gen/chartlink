package synth

import (
	"math"
	"math/rand/v2"
)

// Common US surnames, most common first. Real surname frequencies have a
// steep head and a very long tail; the head comes from this list and the
// tail is made up (see surname).
var commonLast = []string{
	"SMITH", "JOHNSON", "WILLIAMS", "BROWN", "JONES", "GARCIA", "MILLER", "DAVIS",
	"RODRIGUEZ", "MARTINEZ", "HERNANDEZ", "LOPEZ", "GONZALEZ", "WILSON", "ANDERSON",
	"THOMAS", "TAYLOR", "MOORE", "JACKSON", "MARTIN", "LEE", "PEREZ", "THOMPSON",
	"WHITE", "HARRIS", "SANCHEZ", "CLARK", "RAMIREZ", "LEWIS", "ROBINSON", "WALKER",
	"YOUNG", "ALLEN", "KING", "WRIGHT", "SCOTT", "TORRES", "NGUYEN", "HILL", "FLORES",
	"GREEN", "ADAMS", "NELSON", "BAKER", "HALL", "RIVERA", "CAMPBELL", "MITCHELL",
	"CARTER", "ROBERTS", "GOMEZ", "PHILLIPS", "EVANS", "TURNER", "DIAZ", "PARKER",
	"CRUZ", "EDWARDS", "COLLINS", "REYES", "STEWART", "MORRIS", "MORALES", "MURPHY",
	"COOK", "ROGERS", "GUTIERREZ", "ORTIZ", "MORGAN", "COOPER", "PETERSON", "BAILEY",
	"REED", "KELLY", "HOWARD", "RAMOS", "KIM", "COX", "WARD", "RICHARDSON", "WATSON",
	"BROOKS", "CHAVEZ", "WOOD", "JAMES", "BENNETT", "GRAY", "MENDOZA", "RUIZ", "HUGHES",
	"PRICE", "ALVAREZ", "CASTILLO", "SANDERS", "PATEL", "MYERS", "LONG", "ROSS", "FOSTER",
	"JIMENEZ", "POWELL", "JENKINS", "PERRY", "RUSSELL", "SULLIVAN", "BELL", "COLEMAN",
	"BUTLER", "HENDERSON", "BARNES", "GONZALES", "FISHER", "VASQUEZ", "SIMMONS", "ROMERO",
	"JORDAN", "PATTERSON", "ALEXANDER", "HAMILTON", "GRAHAM", "REYNOLDS", "GRIFFIN",
	"WALLACE", "MORENO", "WEST", "COLE", "HAYES", "BRYANT", "HERRERA", "GIBSON", "ELLIS",
	"TRAN", "MEDINA", "AGUILAR", "STEVENS", "MURRAY", "FORD", "CASTRO", "MARSHALL", "OWENS",
	"HARRISON", "FERNANDEZ", "MCDONALD", "WOODS", "WASHINGTON", "KENNEDY", "WELLS",
	"VARGAS", "HENRY", "CHEN", "FREEMAN", "WEBB", "TUCKER", "GUZMAN", "BURNS", "CRAWFORD",
	"OLSON", "SIMPSON", "PORTER", "HUNTER", "GORDON", "MENDEZ", "SILVA", "SHAW", "SNYDER",
	"MASON", "DIXON", "MUNOZ", "HUNT", "HICKS", "HOLMES", "PALMER", "WAGNER", "BLACK",
	"ROBERTSON", "BOYD", "ROSE", "STONE", "SALAZAR", "FOX", "WARREN", "MILLS", "MEYER",
	"RICE", "SCHMIDT", "GARZA", "DANIELS", "FERGUSON", "NICHOLS", "STEPHENS", "SOTO",
	"WEAVER", "RYAN", "GARDNER", "PAYNE", "GRANT", "DUNN", "KELLEY", "SPENCER", "HAWKINS",
}

var maleFirst = []string{
	"JAMES", "ROBERT", "JOHN", "MICHAEL", "DAVID", "WILLIAM", "RICHARD", "JOSEPH",
	"THOMAS", "CHRISTOPHER", "CHARLES", "DANIEL", "MATTHEW", "ANTHONY", "MARK", "DONALD",
	"STEVEN", "ANDREW", "PAUL", "JOSHUA", "KENNETH", "KEVIN", "BRIAN", "GEORGE", "TIMOTHY",
	"RONALD", "JASON", "EDWARD", "JEFFREY", "RYAN", "JACOB", "GARY", "NICHOLAS", "ERIC",
	"JONATHAN", "STEPHEN", "LARRY", "JUSTIN", "SCOTT", "BRANDON", "BENJAMIN", "SAMUEL",
	"GREGORY", "ALEXANDER", "PATRICK", "FRANK", "RAYMOND", "JACK", "DENNIS", "JERRY",
	"JOSE", "LUIS", "CARLOS", "JUAN", "JORGE", "MIGUEL", "FRANCISCO", "ALEJANDRO",
	"WEI", "MINH", "RAJ", "ARJUN", "AHMED", "OMAR", "DMITRI", "IVAN", "KWAME", "HIROSHI",
}

var femaleFirst = []string{
	"MARY", "PATRICIA", "JENNIFER", "LINDA", "ELIZABETH", "BARBARA", "SUSAN", "JESSICA",
	"SARAH", "KAREN", "LISA", "NANCY", "BETTY", "MARGARET", "SANDRA", "ASHLEY", "KIMBERLY",
	"EMILY", "DONNA", "MICHELLE", "CAROL", "AMANDA", "DOROTHY", "MELISSA", "DEBORAH",
	"STEPHANIE", "REBECCA", "SHARON", "LAURA", "CYNTHIA", "KATHLEEN", "AMY", "ANGELA",
	"SHIRLEY", "ANNA", "BRENDA", "PAMELA", "EMMA", "NICOLE", "HELEN", "SAMANTHA",
	"KATHERINE", "CHRISTINE", "DEBRA", "RACHEL", "CATHERINE", "CAROLYN", "JANET", "RUTH",
	"MARIA", "ROSA", "GUADALUPE", "CARMEN", "ANA", "SOFIA", "LUCIA", "VICTORIA",
	"MEI", "LINH", "PRIYA", "ANANYA", "FATIMA", "AISHA", "OLGA", "NATASHA", "AMARA", "YUKI",
}

// zipf returns cumulative weights for rank^-s over n items.
func zipf(n int, s float64) []float64 {
	cum := make([]float64, n)
	t := 0.0
	for i := range n {
		t += math.Pow(float64(i+1), -s)
		cum[i] = t
	}
	for i := range cum {
		cum[i] /= t
	}
	return cum
}

func pick(r *rand.Rand, cum []float64) int {
	u := r.Float64()
	lo, hi := 0, len(cum)-1
	for lo < hi {
		m := (lo + hi) / 2
		if cum[m] < u {
			lo = m + 1
		} else {
			hi = m
		}
	}
	return lo
}

var (
	lastCum   = zipf(len(commonLast), 0.5)
	maleCum   = zipf(len(maleFirst), 0.8)
	femaleCum = zipf(len(femaleFirst), 0.8)
)

var onsets = []string{"B", "D", "K", "L", "M", "N", "P", "R", "S", "T", "V", "Z", "CH", "SH", "BR", "KR", "ST", "TR", "GR", "H"}
var vowels = []string{"A", "E", "I", "O", "U", "AI", "OU", "EI"}
var codas = []string{"", "N", "R", "S", "K", "L", "M", "NG", "RT", "LD"}
var endings = []string{"", "SKI", "SON", "OV", "ELLI", "MANN", "ARO", "IAN", "OWITZ", "ER", "INI", "AKIS"}

// tailShare is the fraction of people with a rare surname. Around a third
// of the US population has a surname outside the top few thousand.
const tailShare = 0.35

// surname draws a last name: from the common list by frequency, or a
// made-up rare one that almost nobody else shares.
func surname(r *rand.Rand) string {
	if r.Float64() >= tailShare {
		return commonLast[pick(r, lastCum)]
	}
	s := ""
	for range 2 + r.IntN(2) {
		s += onsets[r.IntN(len(onsets))] + vowels[r.IntN(len(vowels))] + codas[r.IntN(len(codas))]
	}
	return s + endings[r.IntN(len(endings))]
}

func given(r *rand.Rand, sex byte) string {
	if sex == 'M' {
		return maleFirst[pick(r, maleCum)]
	}
	return femaleFirst[pick(r, femaleCum)]
}
