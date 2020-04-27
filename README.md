# ToDoGo

A To-Do app written in Go.

Each To-Do is stored as a row in a PostgreSQL database.

## Schema

```
"ID":          integer,
"Title":       string,
"Description": string,
"Status":      boolean,
"Created":     timestamp,
"Modified":    timestamp,
```

* **ID** is mandatory and is generated automatically (except with PUT).
* **Title** and **Description** can be empty strings.
* **Status** is `false` when a To-Do is first created, unless it is set to `true` in the POST request.
* **Created** and **Modified** are added automatically. **Modified** is updated whenever the To-Do is updated (NB: this is true even if the update contains no changes).

## REST API

| endpoint    | verb   | returns                     | description |
|-------------|--------|-----------------------------|-------------|
| /_status    | GET    | `{ "status": description }` | Returns the service status, either "ready" or "no database connection".  Returns a 503 error if the service is down. |
| /todo       | POST   | `{ "id": id }`              | Creates a To-Do and returns its ID. The JSON payload can include "Title" (string), "Description" (string), and "Status" (boolean, defaults to `false`). |
| /todos      | GET    | []ToDo                      | Returns all the To-Dos as a JSON array. |
| /todos/{id} | GET    | ToDo                        | Returns the To-Do with that ID as JSON. |
| /todos/{id} | PUT    | `{ "id": id }`              | Updates the To-Do with that ID, or creates a new one. The JSON payload can include "Title" (string), "Description" (string), and "Status" (boolean, optional). Any fields not set will use their default. |
| /todos/{id} | PATCH  | `{ "updated": count }`      | Updates the To-Do with that ID. The JSON payload should include "Title" (string), "Description" (string), and "Status" (boolean, optional). |
| /todos/{id} | DELETE | `{ "deleted": count }`      | Deletes the To-Do with that ID. Returns the count of deleted To-Dos (currently 0 or 1). |

For GET and PATCH, a 404 error is returned if there is no To-Do with that ID.

## Remaining work

* add filters and sorting to the `/todos` endpoint
* decide whether Title should be a required field
* consider PUT/DELETE at `/todo/{id}/done` as a way to toggle the status
* create formal REST API documentation (Swagger?, OpenAPI?)
* prettier schema documentation
* don't have the database password in plain text in the repo
* add unit and API tests

## Running for development

To work on the application, you will need an instance of Postgres. The simplest way to do this is to have one running in Docker with the database created.

    $ docker run --name db -e POSTGRES_PASSWORD=mysecretpassword -p 5432:5432 -d postgres
    $ psql -h localhost -U postgres -f init.sql

## Running locally for use

1. Build the ToDoGo Docker image:

        $ docker build -t todogo .

2. Run the app with `docker-compose`:

        $ docker-compose -f docker-compose.yml up -d

The app is accessible at port `:8080`.

## Running on GCP

1. Connect to the VM via SSH.

2. Start up the docker-compose container:
```sh
$ docker run docker/compose:1.25.5 version
```

3. Start up the Todogo service:
```sh
$ docker run --rm \
    -v /var/run/docker.sock:/var/run/docker.sock \
    -v "$PWD:$PWD" \
    -w="$PWD" \
    docker/compose:1.25.5 up -d
```