Just the events
===============

Readable and auditable filter for just listening docker events

First terminal:

    go build
    ./just-the-events

Second terminal:

    docker -H unix:///tmp/just-the-events.sock events

Third terminal, do something with docker.

Licence
-------

3 terms BSD licence, © 2016 Mathieu Lecarme
