package jazz

import (
	"crypto/rand"
	"fmt"
	"regexp"
	"runtime"
	"time"
)

func (j *Jazz) LoadTime(start time.Time) {
	elapsed := time.Since(start)
	pc, _, _, _ := runtime.Caller(1)
	funcObj := runtime.FuncForPC(pc)
	runtimeFunc := regexp.MustCompile(`^.*\.(.*)$`)
	name := runtimeFunc.ReplaceAllString(funcObj.Name(), "$1")
	j.InfoLog.Println(fmt.Sprintf("Load Time: %s took %s", name, elapsed))
}

// RandomString generates a random string length n from values in the const randomString
func (j *Jazz) RandomString(n int) string {
	s, r := make([]rune, n), []rune(randomString)

	for i := range s {
		p, _ := rand.Prime(rand.Reader, len(r))
		x, y := p.Uint64(), uint64(len(r))
		s[i] = r[x%y]

	}

	return string(s)
}
