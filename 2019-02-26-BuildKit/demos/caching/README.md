# Caching with BuildKit

## Caching package installs

Let's demo the cache mount with a Debian `apt` package cache.

### Without using a cache mount

```console
$ docker build -f Dockerfile.before --target figlet --iidfile before.iid .
$ docker history $(< before.iid)
```

---

### With BuildKit, using a cache mount

First run the updater (usually this would be in the same target!)

```console
$ docker build --target updater --iidfile updater.iid .
$ docker history $(< updater.iid)
```

- note the tiny size of the updater layers, also see the increased size in Build Cache (`docker system df`)

---

Now we can install packages without having to `apt update`:

```console
$ docker build --target figlet --iidfile figlet.iid .
$ docker history $(< figlet.iid)
```
