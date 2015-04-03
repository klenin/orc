ORC
===

Course work.

**ORC** helps to participants simplify interaction with
organizers of various competitions, events, tournaments, seminars,
conferences, etc.

Installation
------------

### Requirements

[Go][1] v1.2 or higher.

[PostgreSQL][2] v9.3 or higher.

### Getting Source Files

Install [git][3]. Get a copy of repository:

    $ git clone git@github.com:GuraYaroslava/orc.git

### Install Packages

- [package pq][4] ([godoc](http://godoc.org/github.com/lib/pq))

        $ go get github.com/lib/pq

- [package securecookie][5] ([godoc](http://godoc.org/github.com/gorilla/securecookie))

        $ go get github.com/gorilla/securecookie

### Running

Create database `orc`:

    CREATE DATABASE orc;

Create a user account called `admin` and grant permission for database called `orc`:

    CREATE USER admin WITH PASSWORD 'admin';
    GRANT ALL PRIVILEGES ON DATABASE orc to admin;

Local postgres connection string for a database called `orc`:

    $ export DATABASE_URL="user=admin host=localhost dbname=orc password=password sslmode=disable"

Port (port is 5000 by default):

    $ export PORT="6543"

Run:
    $ go build && orc.exe

Run with downloading test data:

    $ go build && orc.exe -test-data=true

Print all server routers:

    $ console routers

[1]: https://golang.org
[2]: http://www.postgresql.org
[3]: http://git-scm.com
[4]: https://github.com/lib/pq
[5]: http://www.gorillatoolkit.org/pkg/securecookie
