# SIS
Simple image server
I came up with this because i want a way to share files over the internet but only allow a set amount of accesses before the file expires and bescomes inaccessible.
this has two endpoints:

1. POST /add/:count, (expects a file in multiform file with name: datafile) and the Authorization header set to env ACCESS_TOKEN and will take a new file, further count in the path sets the amount of allowed accesses, 0 means infinite, should be a positiv decimal number OR 0.
2. GET /d/:id, takes the id given back by a sucessful add and will return if the max hit count wasnt reached yet, if the correct Auth header is provided this will also return the file beyond the limit.

## Running
Get docker/docker-compose
Run `docker-compose up --build` in the root

# License
As with everything i do, this is licensed as Free software under the GPL 2.0
