#### **Render Capabilities**
* display a template
```go
type Handlers struct {
    *jazz.Jazz
    data.Models
}

func (h *Handlers) Home(w http.ResponseWriter, r *http.Request) {
	err := h.Render.Page(w, r, "home", nil, nil)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}

```
* display a template with data
```go
func (h *Handlers) Sessions(w http.ResponseWriter, r *http.Request) {
	 
	myData := "data"
	myVars:="vars"
	
	
	//set myData to templateData data field
	d:=templateData{
        Data:myData,
    }

	//set myVars values to a jet Vars key
	vars := make(jet.VarMap)
	vars.Set("foo", myVars)

	err := h.Render.Page(w, r, "home", vars, d)
	if err != nil {
		h.ErrorLog.Println("error rendering:", err)
	}
}
```

