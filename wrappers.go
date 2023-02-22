package jazz

import (
	"context"
	"net/http"
)

func (j *Jazz) Get(s string, h http.HandlerFunc) {
	j.Routes.Get(s, h)
}

func (j *Jazz) Post(s string, h http.HandlerFunc) {
	j.Routes.Post(s, h)
}

func (j *Jazz) Use(m ...func(http.Handler) http.Handler) {
	j.Routes.Use(m...)
}

func (j *Jazz) Page(w http.ResponseWriter, r *http.Request, tmpl string, variables, data interface{}) error {
	return j.Render.Page(w, r, tmpl, variables, data)
}

func (j *Jazz) SessionPut(ctx context.Context, key string, val interface{}) {
	j.Session.Put(ctx, key, val)
}

func (j *Jazz) SessionHas(ctx context.Context, key string) bool {
	return j.Session.Exists(ctx, key)
}

func (j *Jazz) SessionGet(ctx context.Context, key string) interface{} {
	return j.Session.Get(ctx, key)
}

func (j *Jazz) SessionRemove(ctx context.Context, key string) {
	j.Session.Remove(ctx, key)
}

func (j *Jazz) SessionRenew(ctx context.Context) error {
	return j.Session.RenewToken(ctx)
}

func (j *Jazz) SessionDestroy(ctx context.Context) error {
	return j.Session.Destroy(ctx)
}
