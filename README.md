ORC
===

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

Administrator credentials to login the system:

    Login: admin
    Password: password

## External dependencies

- [package pq][4] ([godoc](http://godoc.org/github.com/lib/pq))
- [package securecookie][5] ([godoc](http://godoc.org/github.com/gorilla/securecookie))

All external dependenies are vendored with project's source code.
Consequently, there is no need to install them separately.

## Preparation of database

Create a database `orc`:

```sql
    CREATE DATABASE orc;
```

Create a user account called `admin` and grant permission for database called `orc`:

```sql
    CREATE USER admin WITH PASSWORD 'admin';
    GRANT ALL PRIVILEGES ON DATABASE orc to admin;
```

## Configuration

Configuration is a set of key-value string pairs, which should be specified in some of following ways

1. Preparing the configuration file, which is `~/.orcrc`, where `~` means home path, by default
    * Configuration file should consist of nonempty lines in format `<key>=<value>`.
    * File path could be overriden by `ORC_CONFIG_PATH` environment variable.
2. Setting up configuration environment variables. The name of variable is <key> prefixed with `ORC_` string.

Configuration could be set up simultaneously throught both ways. In this case, environment variables have higher priority.

### Setting Admin credentials

In order to send emails with team invitations, login confirmations etc., You need to configure Admin credentials.
In this context, admin is the person who represents project's service.

Set following config keys:
`ADMIN_NAME` -- admin's name
`EMAIL` -- the email to send messages from
`EMAIL_PASSWORD` -- mailbox password
`SMTP_SERVER` -- mailing servise SMTP address
`SMTP_PORT` -- mailing service SMTP port
`SERVER_URI` -- server URI to use with mailing templates


### Database and port anvironment variables

In order to connect to posgres database, `DATABASE_URL` config key should be set;
Set the port (5000 by default) through `PORT` config key:


###example:

```console
$ export \
    ORC_ADMIN_NAME="Admin Name" \
    ORC_EMAIL="example@gmail.com" \
    ORC_SMTP_SERVER="smtp.gmail.com" \
    ORC_SMTP_PORT="587" \
    ORC_EMAIL_PASSWORD="password" \
    ORC_SERVER_URI="https://server/link/" \
    ORC_DATABASE_URL="user=admin host=localhost dbname=orc password=password sslmode=disable" \
    ORC_CONFIG_PATH="./server_config"
$ echo -e "PORT=6543\nEMAIL_PASSWORD="qwerty" > ./server_config
```

In this example two config keys are set throught `server_config` file of `cwd`, but the value of `EMAIL_PASSWORD` key would be loaded from `ORC_EMAIL_PASSWORD` env var, as it has higher priority.

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
