package cmd

import (
	"regexp"
	"strconv"
)

var (
	// Major[.Minor[.R[-B]]]
	regVer = regexp.MustCompile("([0-9]+)(.([0-9]+)(.([0-9]+)(-([0-9]+))?)?)?")
)

// Ver Utils
type Ver struct {
	Major    int // 主版本
	Minor    int // 子版本
	Revision int // 修订版本
	Build    int // 构建版本
}

// NewVer from string
func NewVer(s string) (ver *Ver) {
	v := regVer.FindStringSubmatch(s)
	var (
		Major    int
		Minor    int
		Revision int
		Build    int
		e        error
	)
	defer func() {
		_ = recover()
		ver = &Ver{
			Major,
			Minor,
			Revision,
			Build,
		}
	}()
	if v == nil {
		return
	}

	a := v[1]
	i := v[3]
	r := v[5]
	b := v[7]
	if Major, e = strconv.Atoi(a); e != nil {
		panic(e)
	}
	if Minor, e = strconv.Atoi(i); e != nil {
		panic(e)
	}
	if Revision, e = strconv.Atoi(r); e != nil {
		panic(e)
	}
	if Build, e = strconv.Atoi(b); e != nil {
		panic(e)
	}
	return
}

func (v *Ver) String() (s string) {
	if v.Major < 0 {
		return
	}
	s = strconv.Itoa(v.Major)

	if v.Minor < 0 {
		return
	}
	s += "." + strconv.Itoa(v.Minor)
	if v.Revision < 0 {
		return
	}
	s += "." + strconv.Itoa(v.Revision)

	if v.Build < 0 {
		return
	}
	s += "-" + strconv.Itoa(v.Build)
	return
}

// Compare version string
func (v *Ver) Compare(ver string) int {
	return v.CompareV(NewVer(ver))
}

// CompareV Version
//   after ->  1
//   equal ->  0
//   befer -> -1
func (v *Ver) CompareV(ver *Ver) int {
	if v.Major > ver.Major {
		return 1
	}
	if v.Major < ver.Major {
		return -1
	}

	if v.Minor > ver.Minor {
		return 1
	}
	if v.Minor < ver.Minor {
		return -1
	}

	if v.Revision > ver.Revision {
		return 1
	}
	if v.Revision < ver.Revision {
		return -1
	}

	if v.Build > ver.Build {
		return 1
	}
	if v.Build < ver.Build {
		return -1
	}

	return 0
}

//var (
//	// [user:pass@]host[:port]
//regUpURL=regexp.MustCompile("^(([^:@]+)(:([^:@]+))?@)?([^:@]+)(:([0-9]+))?$")
//)
//// UpURL Utils
//type UpURL struct {
//	User string
//	Pass string
//	Host string
//	Port string
//}
//
//// NewUpURL from s
//func NewUpURL(s string) *UpURL {
//	v := regUpURL.FindStringSubmatch(s)
//	if v == nil {
//		return nil
//	}
//	return &UpURL{
//		User: v[2],
//		Pass: v[4],
//		Host: v[5],
//		Port: v[7],
//	}
//}
