layout: true

background-image: url(https://raw.githubusercontent.com/DockerOttawaMeetup/Slides/master/ottawa-docker-logo.jpg)
background-size: 64px
background-position: right 12px bottom 40px

---
class: middle, center

# Docker `HEALTHCHECK`s

_Dave Henderson, January 2017_
---

```Dockerfile
FROM nginx:alpine

RUN apk --no-cache add curl

HEALTHCHECK --interval=5s CMD curl --fail http://localhost:80/ || exit 1
```

--
- only the last healthcheck is used

--
- `--interval` defaults to 30s (time between checks)
- `--timeout` defaults to 30s (check fails if it takes longer than...)
- `--retries` defaults to 3 (unhealthy after this many failed checks)

--
- must exit `0` or `1`

---

## Thank you!
