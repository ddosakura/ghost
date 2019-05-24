package cmd

import (
	"testing"
)

func tUserVer(t *testing.T, ver string, a int, b int, c int, d int, verCheck string) {
	v := NewVer(ver)
	if v.Major == a && v.Minor == b && v.Revision == c && v.Build == d {
		s := v.String()
		if s == verCheck {
			t.Log(ver, "pass")
		} else {
			t.Errorf("error %v -> %s", v, s)
		}
	} else {
		t.Errorf("error %s -> %v", ver, v)
	}
}

var (
	vcSymbol = map[int]string{
		-1: "<",
		0:  "=",
		1:  ">",
	}
)

func tUserVerCompare(t *testing.T, a string, b string, ans int) {
	c := NewVer(a).Compare(b)
	t.Logf("%s %s %s", a, vcSymbol[c], b)
	if c != ans {
		t.Errorf("error, should be: %s %s %s", a, vcSymbol[ans], b)
	}
}

func TestUtilVer(t *testing.T) {
	tUserVer(t, "1.34.56-3", 1, 34, 56, 3, "1.34.56-3")
	tUserVer(t, "2.5", 2, 5, 0, 0, "2.5.0-0")
	tUserVer(t, "3.4.g-1", 3, 4, 0, 0, "3.4.0-0")
	tUserVer(t, "3.4g.1", 3, 4, 0, 0, "3.4.0-0")
	tUserVer(t, "1a", 1, 0, 0, 0, "1.0.0-0")

	tUserVerCompare(t, "1.5", "1.3", 1)
	tUserVerCompare(t, "1.3", "1.3.0-1", -1)
	tUserVerCompare(t, "0.5", "1.3", -1)
}
