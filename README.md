
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
cat kustomization.yaml
```

`example`
```
korp scan -f ./sample-yaml
cat kustomization.yaml
```

### Pull

`TODO`

### Push

`TODO`

### Patch

`TODO`

---

## TODOs

- [x] scan
- [ ] pull
- [ ] push
- [ ] patch `TBD`
