ORC
===

Course work.

**ORC** helps to simplify the interaction of the participants with
the organizers of the various competitions, events, tournaments, seminars,
conferences, etc.

Installation
------------------------

### Requirements

[Go][1] v1.2 or higher.

[PostgreSQL][2] v9.3 or higher.

### Getting Source Files

Install [git][3]. Get a copy of repository:

    $ git clone git@github.com:GuraYaroslava/orc.git

### Install Gorilla

[Gorilla][4] is a web toolkit for the Go programming language.

- [package securecookie][5]

        $ go get github.com/gorilla/securecookie

- [package sessions][6]

        $ go get github.com/gorilla/sessions

### Running

    $ go build && orc.exe

[1]: https://golang.org
[2]: http://www.postgresql.org
[3]: http://git-scm.com
[4]: http://www.gorillatoolkit.org/
[5]: http://www.gorillatoolkit.org/pkg/securecookie
[6]: http://www.gorillatoolkit.org/pkg/sessions
