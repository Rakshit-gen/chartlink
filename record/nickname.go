package record

import "slices"

// nicknames maps a nickname to the formal names it is short for. Front desks
// register whatever the patient says, so BOB at one clinic is ROBERT at
// another. Kept to common US names; a hospital would load its own list.
var nicknames = map[string][]string{
	"AL": {"ALBERT", "ALFRED", "ALAN", "ALEXANDER"}, "ALEX": {"ALEXANDER", "ALEXANDRA"},
	"ANDY": {"ANDREW"}, "ANNIE": {"ANN", "ANNA"}, "BECKY": {"REBECCA"},
	"BEN": {"BENJAMIN"}, "BETH": {"ELIZABETH"}, "BETTY": {"ELIZABETH"},
	"BILL": {"WILLIAM"}, "BILLY": {"WILLIAM"}, "BOB": {"ROBERT"}, "BOBBY": {"ROBERT"},
	"CATHY": {"CATHERINE"}, "CHARLIE": {"CHARLES"}, "CHRIS": {"CHRISTOPHER", "CHRISTINE", "CHRISTINA"},
	"CHUCK": {"CHARLES"}, "DAN": {"DANIEL"}, "DANNY": {"DANIEL"}, "DAVE": {"DAVID"},
	"DEBBIE": {"DEBORAH"}, "DICK": {"RICHARD"}, "DON": {"DONALD"}, "ED": {"EDWARD"},
	"EDDIE": {"EDWARD"}, "FRANK": {"FRANCIS", "FRANKLIN"}, "FRED": {"FREDERICK"},
	"GREG": {"GREGORY"}, "JACK": {"JOHN"}, "JAKE": {"JACOB"}, "JEFF": {"JEFFREY"},
	"JENNY": {"JENNIFER"}, "JERRY": {"GERALD", "JEROME"}, "JIM": {"JAMES"}, "JIMMY": {"JAMES"},
	"JOE": {"JOSEPH"}, "JOEY": {"JOSEPH"}, "JOHNNY": {"JOHN"}, "JON": {"JONATHAN"},
	"KATE": {"KATHERINE", "CATHERINE"}, "KATHY": {"KATHERINE"}, "KATIE": {"KATHERINE"},
	"KEN": {"KENNETH"}, "KIM": {"KIMBERLY"}, "LARRY": {"LAWRENCE"}, "LIZ": {"ELIZABETH"},
	"MAGGIE": {"MARGARET"}, "MATT": {"MATTHEW"}, "MEG": {"MARGARET"}, "MIKE": {"MICHAEL"},
	"NANCY": {"ANN", "ANNE"}, "NATE": {"NATHAN", "NATHANIEL"}, "NICK": {"NICHOLAS"},
	"PAT": {"PATRICK", "PATRICIA"}, "PATTY": {"PATRICIA"}, "PEGGY": {"MARGARET"},
	"PETE": {"PETER"}, "RICH": {"RICHARD"}, "RICK": {"RICHARD"}, "ROB": {"ROBERT"},
	"RON": {"RONALD"}, "SAM": {"SAMUEL", "SAMANTHA"}, "SANDY": {"SANDRA"},
	"STEVE": {"STEVEN", "STEPHEN"}, "SUE": {"SUSAN"}, "SUSIE": {"SUSAN"},
	"TERRY": {"TERRENCE", "TERESA"}, "TOM": {"THOMAS"}, "TOMMY": {"THOMAS"},
	"TONY": {"ANTHONY"}, "VICKY": {"VICTORIA"}, "WILL": {"WILLIAM"},
	"PEPE": {"JOSE"}, "PACO": {"FRANCISCO"}, "LUPE": {"GUADALUPE"},
}

// Nicknames returns every nickname in the table and the formal names it
// stands for, so callers can make one up without reaching into the map.
func Nicknames() map[string][]string { return nicknames }

// Variants returns a cleaned first name with every name it could stand in
// for: its formal names if it is a nickname, every nickname of those, and
// its own nicknames if it is a formal name. JIMMY gives JAMES and JIM.
func Variants(name string) []string {
	formal := append([]string{name}, nicknames[name]...)
	out := slices.Clone(formal)
	for nick, fs := range nicknames {
		for _, f := range fs {
			if slices.Contains(formal, f) && !slices.Contains(out, nick) {
				out = append(out, nick)
			}
		}
	}
	return out
}

// SameGiven reports whether two cleaned first names could be the same
// person's: equal, one a nickname of the other, or two nicknames of one
// formal name (BOB and ROB).
func SameGiven(a, b string) bool {
	if a == b {
		return true
	}
	for _, x := range append(nicknames[a], a) {
		for _, y := range append(nicknames[b], b) {
			if x == y {
				return true
			}
		}
	}
	return false
}
