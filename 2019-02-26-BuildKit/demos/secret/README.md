## BuildKit Demo - using `--secret`

```bash
docker build -t secretdemo \
    --secret id=thesecret,src=../secrets/ejson.key \
    .
```

The image should run like this:

```console
$ docker run --rm secretdemo
The password is swordfish
```
