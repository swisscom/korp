# korp

A command line tool for pushing docker images into a different Docker registry based on given Kubernetes yaml files.

## Installation

1. Download the [latest release](https://github.com/swisscom/korp/releases) and unpack it
2. Add the `korp` binary to your PATH

---

## Usage

### Scan

```
korp scan -f <path to yaml files> -r <docker registry name>
```

`example`

```
korp scan -f ./sample-yaml -r
cat kustomization.yaml
```

### Pull

```
korp pull -k <path to kustomization.yaml>
```

`example`

```
korp scan -f ./sample-yaml
korp pull
docker images
```

### Push

```
korp pull -k <path to kustomization.yaml>
```

`example`

```
docker run -d --name registry --restart always -p 5000:5000 registry
korp scan -f ./sample-yaml -r "localhost:5000"
korp pull
korp push
```

### Patch

`TODO`

### All

`TODO`

---

## TODOs

### features

- [x] scan
- [x] pull
- [x] push
  - [ ] push to registry with auth
- [ ] patch
- [ ] all

### testing

- [ ] scan
- [ ] pull
- [ ] push
- [ ] patch
- [ ] whole

### improvements

- [x] rename and split utils package
- [ ] debug flag / env-var
  - [ ] using cli framework
- [ ] config file
  - [ ] using cli framework

---

## Issues

- CRD image references not recognized
