package record

import "testing"

func TestClean(t *testing.T) {
	for in, want := range map[string]string{
		"O'Brien":     "OBRIEN",
		"o brien":     "OBRIEN",
		"Mary-Jane":   "MARYJANE",
		"  smith  ":   "SMITH",
		"José":        "JOSE",
		"Muñoz":       "MUNOZ",
		"":            "",
		"de la Cruz2": "DELACRUZ",
	} {
		if got := Clean(in); got != want {
			t.Errorf("Clean(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestDigits(t *testing.T) {
	if got := Digits("(555) 010-2234"); got != "5550102234" {
		t.Errorf("got %q", got)
	}
}

func TestDateKnown(t *testing.T) {
	if (Date{1980, 0, 3}).Known() || !(Date{1980, 2, 3}).Known() {
		t.Error("Known is wrong")
	}
}
