package render

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/CloudyKit/jet/v6"
	"github.com/ahmedkhaeld/jazz/forms"
	"github.com/alexedwards/scs/v2"
	"github.com/justinas/nosurf"
	"html/template"
	"net/http"
	"path/filepath"
	"strings"
)

// Render holds the data types needed to render a template
type Render struct {
	Renderer      string //what rendering engine to use
	RootPath      string //root path of the template files
	Secure        bool   //is the app running on https
	Debug         bool   //is the app running in debug mode
	UseCache      bool
	TemplateCache map[string]*template.Template
	Port          string              //port number
	ServerName    string              //server name
	JetViews      *jet.Set            //jet views
	Session       *scs.SessionManager //session manager
}

// TemplateData holds the types of data passed from handlers to web template pages
type TemplateData struct {
	IsAuth      bool                   //is the user logged-in
	UserID      int                    //id of the logged-in user
	Role        string                 //role of the logged-in user
	RoleID      int                    //role id of the logged-in user
	AccessLevel int                    //access level of the logged-in user
	IntData     map[string]int         //int data to be passed to the template
	StringData  map[string]string      //string data to be passed to the template
	FloatData   map[string]float64     //float data to be passed to the template
	Data        map[string]interface{} //data to be passed to the template
	CSRFToken   string                 //cross site request forgery token
	Port        string                 //port number
	ServerName  string                 //server name
	Secure      bool                   //is the app running on https
	Error       string                 //error message
	Warning     string                 //warning message
	Flash       string                 //flash message
	Form        *forms.Form
}

// defaultData add default data to every request
func (rn *Render) defaultData(td *TemplateData, r *http.Request) *TemplateData {
	td.ServerName = rn.ServerName
	td.Secure = rn.Secure
	td.Port = rn.Port
	//generates a token to be injected to the template data
	td.CSRFToken = nosurf.Token(r)

	//check if session has key userID,
	//meaning that key holds a session value of logged-in user
	//then set IsAuth to true
	if rn.Session.Exists(r.Context(), "userID") {
		td.IsAuth = true
		td.UserID = rn.Session.GetInt(r.Context(), "userID")
		td.Role = rn.Session.GetString(r.Context(), "role")
		td.AccessLevel = rn.Session.GetInt(r.Context(), "accessLevel")
		td.RoleID = rn.Session.GetInt(r.Context(), "roleID")
	}
	td.Error = rn.Session.PopString(r.Context(), "error")
	td.Flash = rn.Session.PopString(r.Context(), "flash")
	td.Warning = rn.Session.PopString(r.Context(), "warning")

	return td
}

// Page decide to render either a jet or go page
//
// takes in w to write to the response writer; r to access the request; tmplName the name of the view to render;
// variables and data to pass to the view
func (rn *Render) Page(w http.ResponseWriter, r *http.Request, tmplName string, variables, data interface{}) error {
	switch strings.ToLower(rn.Renderer) {
	case "go":
		return rn.GoTemplate(w, r, tmplName, data)
	case "jet":
		return rn.JetPage(w, r, tmplName, variables, data)

	default:
		//case no engine fall through to return err
	}
	return ErrNoEngine
}

func (rn *Render) JetPage(w http.ResponseWriter, r *http.Request, tmplName string, variables, data interface{}) error {

	//load the views path and set the mode
	var views *jet.Set
	//in development mode we want to reload the templates on every request
	if rn.Debug {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rn.RootPath)),
			jet.InDevelopmentMode(),
		)
	}
	//in production mode we want to load the templates once
	if !rn.Debug {
		views = jet.NewSet(
			jet.NewOSFileSystemLoader(fmt.Sprintf("%s/views", rn.RootPath)),
		)
	}
	rn.JetViews = views

	var vars jet.VarMap
	//check if variables is nil and create a new map
	if variables == nil {
		vars = make(jet.VarMap)
	}
	//if not nil then cast it to jet.VarMap
	if variables != nil {
		vars = variables.(jet.VarMap)
	}

	// declare td that's an empty TemplateData
	td := &TemplateData{}
	//check if data is not nil and cast data to TemplateData and populate that to td
	if data != nil {
		td = data.(*TemplateData)

	}

	//add default data to every request
	td = rn.defaultData(td, r)

	//get the template from the jet views
	t, err := rn.JetViews.GetTemplate(fmt.Sprintf("%s.jet", tmplName))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNoPage, err)
	}

	//execute the template with the data and variables to the response writer
	if err = t.Execute(w, vars, td); err != nil {
		return fmt.Errorf("%w: %s", ErrExecPage, err)
	}
	return nil
}

// GoPage render a go template page
func (rn *Render) GoPage(w http.ResponseWriter, r *http.Request, tmplName string, data interface{}) error {
	//parse template files
	tmpl, err := template.ParseFiles(fmt.Sprintf("%s/views/%s.page.tmpl", rn.RootPath, tmplName))
	if err != nil {
		return fmt.Errorf("%w: %s", ErrNoPage, err)
	}
	//create a template data instance to hold passed data
	td := &TemplateData{}

	//check if data is not nil and cast data to TemplateData and assign to td
	if data != nil {
		td = data.(*TemplateData)

	}

	//add default data from the request to td
	td = rn.defaultData(td, r)

	err = tmpl.Execute(w, &td)
	if err != nil {
		return fmt.Errorf("%w: %s", ErrExecPage, err)
	}
	return nil
}

func (rn *Render) GoTemplate(w http.ResponseWriter, r *http.Request, tmplName string, data interface{}) error {
	var templateCache map[string]*template.Template
	//avoid loading the cache every time we display a single page when calling GoPage func that is responsible
	//to render the view of a request for production mode
	if rn.UseCache {
		templateCache = rn.TemplateCache
	} else {
		templateCache, _ = rn.CreateTemplateCache()

	}

	theTemplate, ok := templateCache[tmplName]
	if !ok {
		return errors.New("template is not in the views folder ")
	}

	buf := new(bytes.Buffer)

	//td an instance to hold the passed data from the handler
	td := &TemplateData{}

	//check if data is not nil and cast data to TemplateData and assign to td
	if data != nil {
		td = data.(*TemplateData)

	}
	td = rn.defaultData(td, r)

	_ = theTemplate.Execute(buf, &td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Error writing to the browser", err)
		return err
	}

	return nil
}

// CreateTemplateCache load the templates that inside the views folder when the application starts;
// it creates a map of the pages names inside the views folder as its index to quickly be retrieved and uniquely stored
// and when it is going to look the page name up, it gives a fully parsed and ready to use Template
func (rn *Render) CreateTemplateCache() (map[string]*template.Template, error) {

	//cache will hold all the existing template at the start of the application
	cache := map[string]*template.Template{}

	//get all the templates' full path inside the views folder that start with any name and end with "page.tmpl"
	pages, err := filepath.Glob(fmt.Sprintf("%s/views/*.page.tmpl", rn.RootPath))
	if err != nil {
		return cache, err
	}

	//loop through all pages full path then extract the base name which is the actual template name
	for _, page := range pages {
		tmplName := filepath.Base(page)

		//create a new template of the provided template name with its path supported with the available custom functions
		t, err := template.New(tmplName).Funcs(functions).ParseFiles(page)
		if err != nil {
			return cache, err
		}
		//look inside the views if there is any layouts
		matches, err := filepath.Glob(fmt.Sprintf("%s/views/*.layout.tmpl", rn.RootPath))
		if err != nil {
			return cache, err
		}

		//if there are matches, then parse them
		if len(matches) > 0 {
			t, err = t.ParseGlob(fmt.Sprintf("%s/views/*.layout.tmpl", rn.RootPath))
			if err != nil {
				return cache, err
			}
		}
		// append the created template to the cache
		cache[tmplName] = t
	}
	return cache, nil
}
