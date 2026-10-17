package record

import (
	"math"
	"testing"
)

func TestJaroWinkler(t *testing.T) {
	// Reference values from Winkler's papers.
	for _, c := range []struct {
		a, b string
		want float64
	}{
		{"MARTHA", "MARHTA", 0.961},
		{"DWAYNE", "DUANE", 0.840},
		{"DIXON", "DICKSONX", 0.813},
		{"SMITH", "SMITH", 1},
		{"SMITH", "", 0},
		{"ABC", "XYZ", 0},
	} {
		if got := JaroWinkler(c.a, c.b); math.Abs(got-c.want) > 0.001 {
			t.Errorf("JaroWinkler(%q, %q) = %.3f, want %.3f", c.a, c.b, got, c.want)
		}
	}
}

func TestJaroWinklerSymmetric(t *testing.T) {
	names := []string{"JON", "JOHN", "JONATHAN", "GARCIA", "GRACIA", "GARZA", "A", "AB"}
	for _, a := range names {
		for _, b := range names {
			if x, y := JaroWinkler(a, b), JaroWinkler(b, a); math.Abs(x-y) > 1e-12 {
				t.Errorf("%s/%s: %.4f vs %.4f", a, b, x, y)
			}
		}
	}
}

func TestJaroWinklerLongInput(t *testing.T) {
	long := string(make([]byte, 100))
	if got := JaroWinkler(long+"A", long+"B"); got < 0 || got > 1 {
		t.Errorf("got %f", got)
	}
}

func TestSoundex(t *testing.T) {
	for in, want := range map[string]string{
		"ROBERT": "R163", "RUPERT": "R163", "ASHCRAFT": "A261",
		"TYMCZAK": "T522", "PFISTER": "P236", "HONEYMAN": "H555",
		"LEE": "L000", "": "",
	} {
		if got := Soundex(in); got != want {
			t.Errorf("Soundex(%q) = %q, want %q", in, got, want)
		}
	}
}

func TestOSA(t *testing.T) {
	for _, c := range []struct {
		a, b string
		want int
	}{
		{"RYAN", "RYSN", 1}, {"AMY", "MY", 1}, {"GARCIA", "GRACIA", 1}, {"CA", "AC", 1},
		{"ANNA", "ANA", 1}, {"MARK", "MARK", 0}, {"", "ABC", 3}, {"JOHN", "MARY", 4},
		{"KITTEN", "SITTING", 3},
	} {
		if got := OSA(c.a, c.b); got != c.want {
			t.Errorf("OSA(%s, %s) = %d, want %d", c.a, c.b, got, c.want)
		}
		if got := OSA(c.b, c.a); got != c.want {
			t.Errorf("OSA(%s, %s) = %d, want %d", c.b, c.a, got, c.want)
		}
	}
}
