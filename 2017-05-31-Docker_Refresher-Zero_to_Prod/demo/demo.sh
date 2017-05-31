#!/bin/sh

# First we build it
docker build -t figlet .

# Then we run it
docker run figlet

# Now let's change an environment variable and run it again
docker run -e THING=Ottawa figlet
