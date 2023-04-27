## Sample .drone.yml, NOT ".yaml"

```yaml
# Remember File EXTENTION IS = ".yml", NOT ".yaml"
kind: pipeline
type: docker
name: default

steps:
  - name: test
    image: alpine
    commands:
      - echo hello
      - echo world
```

## Skip CI Convention

在 commit 信息中只要包含了下面几个关键词就会跳过 CI，不会触发 CI Build

```
[skip ci]
[ci skip]
[no ci]
[skip actions]
[actions skip]
```

## Stop pipeline after

```yaml
- exit 78
```
