# Docker Compose
## Background
We usually create a single container application such as Docker when we create an application. Previous tasks enabled us to create a single container application using Docker, but you need to learn how to make multiplex container applications. Today's lesson aims to create multiplex container applications using docker-compose.

## Task
Based on [ref](https://qiita.com/muroya2355/items/d48c384a4a82c7ed34ae), please build a multiplex container including Golang and PostgreSQL using Docker Compose. (In your feature, you can build an environment for React-Go application.)

- Build a Postgres server and Golang container using Docker
  - [Dokcer ComposeでGoとPostgreSQLの開発環境構築](https://qiita.com/muroya2355/items/d48c384a4a82c7ed34ae)
- Build a Golang server using the below script
  - main.go
```
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type TestUser struct {
	UserID   int
	Password string
}
func main() {
	http.HandleFunc("/a", handler)
	http.ListenAndServe(":8181", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	var Db *sql.DB

	Db, err := sql.Open("postgres", "host=postgres user=app_user password=password dbname=app_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	sql := "SELECT user_id, user_password FROM TEST_USER WHERE user_id=$1;"
	pstatement, err := Db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	queryID := 1
	var testUser TestUser
	err = pstatement.QueryRow(queryID).Scan(&testUser.UserID, &testUser.Password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, testUser.Password)
}

```

- Open a port for Golang server
The steps are shown below. You can solve the assignment by reading the reference if you don't want to see them.

<details><summary> Answer </summary>

### Step0: Install Golang
https://go.dev/doc/install

### Step1: Create golang + postgreSQL
If you don't know GOPATH, you can get GOPATH by executing the below command.
```
go env GOPATH
```

#### The directory architecture

    GOPATH
    ├──docker-compose.yml
    ├──docker
    │  |──golang
    │  │   └─Dockerfile
    │  └─postgres
    │      ├─Dockerfile
    │      └─init
    │          ├─1_create_table.sql
    │          └─2_insert_testdata.sql
    └─go
       └─main.go

#### Create files
- docker-compose.yml
```
version: '3'
services:
  postgres:
    container_name: postgres
    build:
      context: .
      dockerfile: ./docker/postgres/Dockerfile
    environment:
      - POSTGRES_USER=app_user
      - POSTGRES_PASSWORD=password
      - POSTGRES_DB=app_db
  app:
    container_name: app
    depends_on:
      - postgres
    build:
      context: .
      dockerfile: ./docker/golang/Dockerfile
    environment:
      - GOPATH=/go
    volumes:
      - ./go:/go/src/app/go/
    command: go run main.go
```
- docker/golang/Dockerfile
```
FROM golang:1.9
WORKDIR /go/src/denki/go
ADD ./go .
RUN go get github.com/lib/pq
```
- docker/postgres/Dockerfile
```
FROM postgres:latest
RUN localedef -i ja_JP -c -f UTF-8 -A /usr/share/locale/locale.alias ja_JP.UTF-8
ENV LANG ja_JP.UTF-
COPY ./docker/postgres/init/*.sql /docker-entrypoint-initdb.d/
```
- docker/postgres/init/1_create_table.sql
```
CREATE TABLE TEST_USER (
    user_id BIGINT PRIMARY KEY,
    user_password VARCHAR(20) NOT NULL
);
```
- docker/postgres/init/2_insert_testdata.sql
```
DELETE FROM TEST_USER;
INSERT INTO TEST_USER VALUES(1, 'password1');
INSERT INTO TEST_USER VALUES(2, 'password2');
```
- docker/go/main.go
```
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/lib/pq"
)
type TestUser struct {
    UserID   int
    Password string
}

func main() {
    var Db *sql.DB
    Db, err := sql.Open("postgres", "host=postgres user=app_user password=password dbname=app_db sslmode=disable")
    if err != nil {
        log.Fatal(err)
    }
    sql := "SELECT user_id, user_password FROM TEST_USER WHERE user_id=$1;"
    pstatement, err := Db.Prepare(sql)
    if err != nil {
        log.Fatal(err)
    }

    queryID := 1
    var testUser TestUser

    err = pstatement.QueryRow(queryID).Scan(&testUser.UserID, &testUser.Password)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Println(testUser.UserID, testUser.Password)
}
```
#### Execute docker-compose
```
$ docker-compose build
$ docker-compose up
```
<details><summary> Error </summary>

- Modify the golang dockerfile
If the error has been ocurred, you should replace below command.
  - Error
```
executor failed running [/bin/sh -c go get github.com/lib/pq]: exit code: 2
ERROR: Service 'app' failed to build : Build failed 
```
  - Dockerfile
```
# Specify the base image for the go app.
FROM golang:1.15

WORKDIR /go/src/github.com/postgres-go
# Copy everything from this project into the filesystem of the container.
COPY . .
# Obtain the package needed to run code. Alternatively use GO Modules. 
RUN go get -u github.com/lib/pq

WORKDIR /go/src/yuta/go
# Add the contents of the host OS . /go contents to the working directory
ADD ./go .
```
</details>

### Step2 Build a server in the golang container
Copy a below code and paste it in /go/main.go
```
package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"

	_ "github.com/lib/pq"
)

type TestUser struct {
	UserID   int
	Password string
}
func main() {
	http.HandleFunc("/a", handler)
	http.ListenAndServe(":8181", nil)
}

func handler(w http.ResponseWriter, r *http.Request) {

	var Db *sql.DB

	Db, err := sql.Open("postgres", "host=postgres user=app_user password=password dbname=app_db sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}
	sql := "SELECT user_id, user_password FROM TEST_USER WHERE user_id=$1;"
	pstatement, err := Db.Prepare(sql)
	if err != nil {
		log.Fatal(err)
	}
	queryID := 1
	var testUser TestUser
	err = pstatement.QueryRow(queryID).Scan(&testUser.UserID, &testUser.Password)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, testUser.Password)
}

```
We can complete building a server.
The detail of explanation is described in [this](https://note.com/miso_hijiki/n/nf816bdf23430)

But, we didn't access the server as below in local terminal.
```
% curl localhost:8181/a  
curl: (7) Failed to connect to localhost port 8181: Connection refused
```
### Step3 Open a port of the golang server

Add ports under services.app in docker-compose.yaml
```
services:
  app:
    ports:
      - "8181:8181"
```

The task will be terminated when the following message is displayed
```
% curl localhost:8181/a                           
password1
```

</details>
