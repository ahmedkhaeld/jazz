package render

import (
	"github.com/CloudyKit/jet/v6"
	"net/http/httptest"
	"testing"
)

var pageData = []struct {
	name          string
	renderer      string
	template      string
	errorExpected bool
	errorMessage  string
}{
	{"go_page", "go", "home.page.tmpl", false, "error rendering go template"},
	{"go_page_no_template", "go", "no-file", true, "no error rendering non-existent go template, when one is expected"},
	{"jet_page", "jet", "home", false, "error rendering jet template"},
	{"jet_page_no_template", "jet", "no-file", true, "no error rendering non-existent jet template, when one is expected"},
	{"invalid_render_engine", "foo", "home", true, "no error rendering with non-existent template engine"},
}

func TestRender_Page(t *testing.T) {
	for _, e := range pageData {
		r, err := getSessionData()
		if err != nil {
			t.Error(err)
		}

		w := httptest.NewRecorder()

		testRenderer.Renderer = e.renderer
		testRenderer.RootPath = "./testdata"
		if e.renderer == "go" {
			tc, _ := testRenderer.CreateTemplateCache()
			testRenderer.TemplateCache = tc
		}

		testRenderer.UseCache = true

		vars := make(jet.VarMap)

		err = testRenderer.Page(w, r, e.template, vars, &TemplateData{})
		if e.errorExpected {
			if err == nil {
				t.Errorf("%s: %s", e.name, e.errorMessage)
			}
		} else {
			if err != nil {
				t.Errorf("%s: %s: %s", e.name, e.errorMessage, err.Error())
			}
		}
	}
}

func TestRender_GoPage(t *testing.T) {
	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"

	err = testRenderer.GoPage(w, r, "home", &TemplateData{})
	if err != nil {
		t.Error("Error rendering page", err)
	}

}

func TestRender_GoTemplate(t *testing.T) {
	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRenderer.Renderer = "go"
	testRenderer.RootPath = "./testdata"
	tc, _ := testRenderer.CreateTemplateCache()
	testRenderer.TemplateCache = tc
	testRenderer.UseCache = true

	err = testRenderer.GoTemplate(w, r, "home.page.tmpl", &TemplateData{})
	if err != nil {
		t.Error("Error rendering page", err)
	}

	err = testRenderer.GoTemplate(w, r, "not-exist.page.tmpl", &TemplateData{})
	if err == nil {

		t.Error("page should not be exists", err)
	}
}

func TestRender_JetPage(t *testing.T) {
	r, err := getSessionData()
	if err != nil {
		t.Error(err)
	}

	w := httptest.NewRecorder()

	testRenderer.Renderer = "jet"
	testRenderer.RootPath = "./testdata"

	vars := make(jet.VarMap)

	err = testRenderer.Page(w, r, "home", vars, &TemplateData{})
	if err != nil {
		t.Error("Error rendering page", err)
	}

}
