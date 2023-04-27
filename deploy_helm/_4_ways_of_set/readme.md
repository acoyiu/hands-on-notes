# Way to use "set"

- 1: (by "values.yaml")
- 2: --set
- 3: --set-file
- 4: --set-string

## will read default "values.yaml" if existed Always

```sh
helm install chartname . \
  --set app.mdname=nginx-b \
  --set-string port=80 \
  -f ./values2.yaml \
  --dry-run
```

## --set-file (mostly for multiline string or dynamic files)

--set-file key=filepath is another variant of --set. It reads the file and use its content ***as a value***

```sh
helm install chartname . --set app.mdname=nginx-b --set-string port=80 -f ./values2.yaml \
  --set-file anno=textAnno.txt \
  --set-file script=textScript.js \
  --dry-run
```
