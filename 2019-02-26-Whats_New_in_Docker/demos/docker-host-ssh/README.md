# Using remote Docker hosts with SSH

```console
$ docker -H ssh://demobox.hairyhenderson.ca version
```

Or, more permanent:

```console
$ export DOCKER_HOST=ssh://demobox.hairyhenderson.ca
$ docker version
```
