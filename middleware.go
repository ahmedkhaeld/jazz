package jazz

import (
	"github.com/justinas/nosurf"
	"net/http"
	"path"
	"strconv"
	"strings"
	"time"
)

/* Middleware:
You can use middleware to inspect the request and make decisions
based on its content before passing it along to the next handler.

For example,
the middleware could respond to the client with an error if the handler
requires authentication and an unauthenticated client sent the request.

Middleware can also collect metrics, log requests, or control access to
resources

*/

func (j *Jazz) SessionLoad(next http.Handler) http.Handler {
	j.InfoLog.Println("Session Load called")
	return j.Session.LoadAndSave(next)
}

func (j *Jazz) NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)
	secure, _ := strconv.ParseBool(j.settings.cookie.secure)

	csrfHandler.ExemptGlob("/api/*")

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   secure,
		SameSite: http.SameSiteStrictMode,
		Domain:   j.settings.cookie.domain,
	})

	return csrfHandler
}

func (j *Jazz) Trace(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		j.InfoLog.Printf("%s\t%s\t%s\tduration:%s", r.Method, r.RequestURI, r.Proto, time.Now().Sub(start))
	})
}

func (j *Jazz) RestrictPrefix(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		for _, p := range strings.Split(path.Clean(r.URL.Path), "/") {
			if strings.HasPrefix(p, "prefix") {
				http.Error(w, "Not Found", http.StatusNotFound)
				return
			}
		}
		next.ServeHTTP(w, r)
	},
	)
}
