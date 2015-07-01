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

### Configuring [Apache][6]

Install modules [mod_proxy][7] and [mod_proxy_http][8]. Uncomment lines in `httpd.conf`:

    LoadModule proxy_module modules/mod_proxy.so
    LoadModule proxy_http_module modules/mod_proxy_http.so

Add to httpd-vhosts.conf:

    <VirtualHost localhost:6543>
        DocumentRoot "path/to/orc"
        ServerName localhost:6543

        # Other directives here

    </VirtualHost>

    <VirtualHost *:80>
        DocumentRoot "path/to/localhost"
        ServerName localhost

        ProxyPass /examplecontroller/controlleraction/args http://localhost:6543/examplecontroller/controlleraction/args

    </VirtualHost>

[1]: https://golang.org
[2]: http://www.postgresql.org
[3]: http://git-scm.com
[4]: https://github.com/lib/pq
[5]: http://www.gorillatoolkit.org/pkg/securecookie
[6]: http://httpd.apache.org
[7]: http://httpd.apache.org/docs/2.2/mod/mod_proxy.html
[8]: http://httpd.apache.org/docs/2.2/mod/mod_proxy_http.html
