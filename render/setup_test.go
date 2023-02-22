package render

import (
	"github.com/CloudyKit/jet/v6"
	"github.com/alexedwards/scs/v2"
	"net/http"
	"os"
	"testing"
	"time"
)

var v = jet.NewSet(
	jet.NewOSFileSystemLoader("./testdata/views"),
	jet.InDevelopmentMode(),
)

var sess *scs.SessionManager

var testRenderer = Render{
	Renderer: "",
	RootPath: "",
	JetViews: v,
}

func TestMain(m *testing.M) {
	sess = scs.New()
	sess.Lifetime = 24 * time.Hour
	sess.Cookie.Persist = true
	sess.Cookie.Secure = false
	sess.Cookie.SameSite = http.SameSiteLaxMode
	testRenderer.Session = sess
	os.Exit(m.Run())
}

// getSession builds a http request that has session data(context)
func getSessionData() (*http.Request, error) {

	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	// get the context from the request
	ctx := r.Context()
	// put session data in ctx
	ctx, _ = sess.Load(ctx, r.Header.Get("X-Session"))
	// put the context back into the request
	r = r.WithContext(ctx)

	return r, nil
}
