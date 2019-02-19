## BuildKit Demo - using `--ssh`

First we try without the `--ssh` option:

```bash
docker build -t sshdemo .
```

This crashes and burns.

Now try with the default ssh agent setting:

```bash
docker build -t sshdemo --ssh=default .
```

Success!

Some other things to try:

- fail fast by adding `required`
- specify the key by fingerprint
- give a different ID
