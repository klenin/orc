ORC
===

Course work.

**ORC** helps participants to simplify interaction with
organizers of various competitions, events, tournaments, seminars,
conferences, etc.

The system has a web interface and provides the ability to generate
different kinds of registration forms, allows to organize an individual
or a group registration of participants in events of any type.

# Installation

## Requirements

[Go][1] v1.5.1 or higher.

[PostgreSQL][2] v9.3 or higher.

## Getting code and running

- Install [git][3].
- Configure your `GOROOT` and `GOPATH` environment variables.
- Install godep:

```console
    $ go get git@github.com:tools/godep.git
```

- Get a copy of the repository:


```console
    $ godep get git@github.com:klenin/orc.git
```

- Configure Your environment as described in following sections.

- Run:

```console
    $ $GOPATH/bin/orc
```
There is also `run.sh` file, providing some ease of local test running e.g. database creation and filling it with local test data.


Test user credentials to login the system (`number` in `[0-19]`):

    Login: user<number>
    Password: secret<number>

## External dependencies

- [package pq][4] ([godoc](http://godoc.org/github.com/lib/pq))
- [package securecookie][5] ([godoc](http://godoc.org/github.com/gorilla/securecookie))

All external dependenies are vendored with project's source code.
Consequently, there is no need to install them separately.

## Preparation of database

Create a database `orc`:

    CREATE DATABASE orc;

Create a user account called `admin` and grant permission for database called `orc`:

    CREATE USER admin WITH PASSWORD 'admin';
    GRANT ALL PRIVILEGES ON DATABASE orc to admin;

## Setting Admin credentials

In order to send emails with team invitations, login confirmations etc., You need to configure Admin credentials.
In this context, admin is the person who represents project's service.

Set following environment variables fo executable:
ADMIN_NAME -- admin's name
EMAIL -- the email to send messages from
EMAIL_PASSWORD -- mailbox password
SMTP_SERVER -- mailing servise SMTP address
SMTP_PORT -- mailing service SMTP port
SERVER_URI -- server URI to use with mailing templates

###example:

```console
export \
    ADMIN_NAME="Admin" \
    EMAIL="example@gmail.com" \
    SMTP_SERVER="smtp.gmail.com" \
    SMTP_PORT="587" \
    EMAIL_PASSWORD="password" \
    SERVER_URI="https://server/link/"
```

Administrator credentials to login the system:

    Login: admin
    Password: password

## Database and port anvironment variables

Set the local postgres connection string for the database called `orc`:

```console
    $ export DATABASE_URL="user=admin host=localhost dbname=orc password=password sslmode=disable"
```

Set the port (5000 by default):

```console
    $ export PORT="6543"
```

## Running as heroku app

- install heroku and follow it's [documentation guidelines](https://devcenter.heroku.com/articles/getting-started-with-go#introduction)
- configure [heroku environment variables](https://devcenter.heroku.com/articles/config-vars) before running!

## Configuring [Apache][6]

### Example using Reverse Proxies/Gateways

Redirection requests from localhost:8080 to localhost:6543.

Install modules [mod_proxy][7] and [mod_proxy_http][8]. Uncomment lines in `httpd.conf`:

```
    # httpd.conf
    LoadModule proxy_module modules/mod_proxy.so
    LoadModule proxy_http_module modules/mod_proxy_http.so
```

Add to httpd-vhosts.conf:

```
    # httpd-vhosts.conf

    <VirtualHost localhost:6543>
        DocumentRoot "path/to/orc"
        ServerName localhost:6543

        # Other directives here

    </VirtualHost>

    <VirtualHost *:80>
        DocumentRoot "path/to/localhost"
        ServerName localhost

        ProxyPass /js http://localhost:6543/js
        ProxyPass /css http://localhost:6543/css
        ProxyPass /img http://localhost:6543/img

        ProxyPass /examplecontroller/controlleraction/args http://localhost:6543/examplecontroller/controlleraction/args

        # Other directives here

        ProxyPassReverse /same/path/ http://localhost:6543

    </VirtualHost>
```

[1]: https://golang.org
[2]: http://www.postgresql.org
[3]: http://git-scm.com
[4]: https://github.com/lib/pq
[5]: http://www.gorillatoolkit.org/pkg/securecookie
[6]: http://httpd.apache.org
[7]: http://httpd.apache.org/docs/2.2/mod/mod_proxy.html
[8]: http://httpd.apache.org/docs/2.2/mod/mod_proxy_http.html
