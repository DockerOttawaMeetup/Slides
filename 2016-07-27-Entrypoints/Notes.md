layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Using `ENTRYPOINT`s for fun and profit!

_Dave Henderson, July 2016_
---

## What's an entrypoint?

--

### It's an instruction

```Dockerfile
FROM alpine

RUN apk update \
  && apk add figlet \
  && rm -rf /var/cache/apk/*

* ENTRYPOINT ["figlet"]
```

???
Demo this!
--

### It's an argument

```console
$ docker run --entrypoint "echo" hairyhenderson/figlet
```

_(technically, `--entrypoint` is just an override for the instruction)_

???
Demo this too!
---

## What's an entrypoint?

### From the docs...

> An `ENTRYPOINT` allows you to configure a container that will run as an executable.

Useful for:

```bash
alias figlet="docker run figlet"
```

???
So, this means I can do `docker run hairyhenderson/figlet hello world` and it'd be just like running `figlet hello world` inside the container.

Instead of running an executable in a container, this approach is more like running the container _as_ an executable. Subtle difference, but useful.
---

### `ENTRYPOINT` and `CMD`

- `CMD` can be used alone, but when overridden you lose the default
- `CMD` can be used with `ENTRYPOINT`, but is used as the argument list

```Dockerfile
# exec form
ENTRYPOINT [ "figlet" ]
# shell form
ENTRYPOINT figlet
# exec form
CMD [ "--help" ]
# shell form (broken!)
CMD --help
```
--
Combinations:

```bash
# exec/exec
figlet --help
# shell/exec
/bin/sh -c figlet --help
# exec/shell
figlet /bin/sh -c --help
# shell/shell
/bin/sh -c figlet /bin/sh -c --help
```

⚠️ _probably best to just not use the shell form_

???
These are similar, but different. Think of `ENTRYPOINT` as an optional _prefix_ to `CMD` for building flexible images. Used together, `ENTRYPOINT` is the command, and `CMD` is the list of arguments.

If you want, you can use _only_ `CMD`, but it can be more difficult to reuse an image if the command is complex and has lots of options.
---

## Use cases
--

### Single-purpose commands

```console
$ alias figlet=docker run hairyhenderson/figlet
$ figlet hello world
...
```

???
See `figlet` earlier... Using `alias` can make this super-powerful. In fact [some people](https://github.com/jfrazelle/dotfiles/blob/master/.dockerfunc) use this approach (and similar) to containerize most of their environment.

---
## Use cases

### Runtime pre-configuration

_Let's get meta!_

```Dockerfile
FROM nginx:1-alpine

# install gomplate
ENV GOMPLATE_VER v1.0.0
ADD https://github.com/hairyhenderson/gomplate/releases/download/$GOMPLATE_VER/gomplate_linux-amd64-slim /usr/local/bin/gomplate
RUN chmod a+rx /usr/local/bin/gomplate

ENV TITLE My Title

COPY index.html.tmpl /
COPY docker-entrypoint.sh /

ENTRYPOINT [ "/docker-entrypoint.sh" ]

CMD ["-g", "daemon off;"]
```
_From [hairyhenderson/remarkjs](https://github.com/hairyhenderson/dockerfiles/tree/master/remarkjs)_

???
This is probably the most powerful use-case.

Let's say you want to give a presentation using `remark.js`. You don't want to have to set up a webserver somewhere, so naturally you're going to run it in a container.

Since the content of the slides aren't part of the image,
we'll need to mount them in. But what if we want to customize the title of the page? We'll need some pre-config for this.

---

## Use cases

#### The `docker-entrypoint.sh` script

```sh
#!/bin/ash
set -e

if [ "${1:0:1}" = '-' ]; then
  set -- nginx "$@"
fi

if [ "$1" = 'nginx' ]; then
  WEBROOT=/usr/share/nginx/html
  gomplate < /index.html.tmpl > $WEBROOT/index.html
  if [ -f /slides.md.tmpl ]; then
    gomplate < /slides.md.tmpl > $WEBROOT/slides.md
  else
    cp /slides.md $WEBROOT/slides.md
  fi
fi

exec "$@"
```
_From [hairyhenderson/remarkjs](https://github.com/hairyhenderson/dockerfiles/tree/master/remarkjs)_

???
Here's a pretty simple entrypoint script. It uses the `gomplate` tool to do templating for `index.html` so that we can set the page title.
Also it'll optionally template the slides themselves if they're provided as a template.

Extra notes:
- we can still run any command other than `nginx` and it'll be run as expected
- common to run things like `gosu` to de-escalate privileges

---

## Use cases

- Pulling credentials and certificates from Vault (or other services)
  - _Vine shoulda done this_ 😲
- Grabbing cloud metadata (like AWS EC2 metadata) - [`gomplate`](https://github.com/hairyhenderson/gomplate) useful for this
- Some more complex examples from official images:
  - [Postgres](https://github.com/docker-library/postgres/blob/master/docker-entrypoint.sh)
  - [Cassandra](https://github.com/docker-library/cassandra/blob/master/3.7/docker-entrypoint.sh)

???
If you look at the `docker-entrypoint.sh` scripts for many of the official images, you'll see the same basic patterns. It's a pre-requisite to being accepted as an official image.

---

## Writing a good entrypoint script

- fail fast - use `set -e`

--

- only use one default command

--

- assume args starting with `-` are flags to the default command

--

- watch out for multiple processes (signal handling can be tricky!)

--

  - but you can always start other processes in the background

--

- don't do _too_ much in the entrypoint - increases service start-time

--

- use `exec` so the final process gets signals and takes over PID 1

--

- pass through non-default args unmodified

---

## Thank you!
