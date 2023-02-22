# jazz

### **A GO module that provides out of the box features help build web applications**
* Render HTML templates
* Session management  
* Easy multiple databases access
* Authentication (web and api)
* Cache capabilities (redis and badger)
* Email Management (SMTP or API) 
* Forget & Reset password management
* File management (upload, download)
* Response management (JSON, XML, CSV)
* Encryption
* CSRF protection
* Form validation
---


#### **Prerequisites**
```
go installed
docker installed
make installed
```
#### **Usage**
* install jazz package into your pc <br>
`go get github.com/ahmedkhaeld/jazz` <br>
* change dir to jazz  <br>
`cd jazz` <br>
* build the cli into dist directory <br>
`make build`  <br>
* change to dist and copy the executable into your desired directory <br>
`cp jazz ~/Desktop` <br>
* cd into Desktop and start new app <br>
`./jazz new coolapp` <br>

* This will create a new project with the following structure
```
newProject
├── data
│   └── models.go
├── handlers
│   └── handlers.go
├── middleware
│   └── middleware.go
├── go.mod
├── go.sum
├── handlers
│   └── handlers.go
├── public
│   ├── icons
│   │   └── style.css
│   └── imgages
│       └── logo.png
├── views
│   └── home.jet
├── .env
├── .gitignore
├── docker-compose.yml
├── go.mod
├── go.sum
├── app.go
├── main.go
├── Makefile
├── readme.md
└── routes.go

```
---


#### **CMD commands**

```
	help                  - show the help commands
	version               - print application version
	new <appName>           - create a new application
	migrate               - runs all up migrations that have not been run previously
	migrate down          - reverses the most recent migration
	migrate reset         - runs all down migrations in reverse order, and then all up migrations
	create migration <name> - creates two new up and down migrations in the migrations folder
	create auth             - creates and runs migrations for authentication tables, and creates models and middleware
	create handler <name>   - creates a stub handler in the handlers directory
	create model <name>     - creates a new model in the data directory
	create mail <name>      - creates two starter mail templates in the mail directory
	create session          - creates a table in the database as a session store
```



