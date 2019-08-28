
# korp
A command line tool for pushing docker images into a different Docker registry based on given Kubernetes yaml files.

## Versioning

Current version: `0.3.2`

---

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
korp scan -f ./samples/yaml
cat kustomization.yaml
```

### Pull

```
korp pull -k <path to kustomization.yaml>
```

`example`

```
korp scan -f ./samples/yaml
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
korp scan -f ./samples/yaml -r "localhost:5000"
korp pull
korp push
```

### Patch

`TODO`

### All

`TODO`

### Debug mode

```
korp -d
```

### Autocompletion

Source the `autocomplete-scripts/*_autocomplete` file in your `.bashrc | .zshrc` file while setting the `PROG` variable to the name of your program.

#### Method 1
```
go build .
source <(./korp autocompletion zsh)
./korp
# now play with tab
```

#### Method 2
```
go build .
PROG=korp source autocomplete-scripts/zsh_autocomplete
./korp
# now play with tab
```

---

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
- [ ] all
- [ ] autocompletion

### general

- [x] rename and split utils package
- [x] debug flag / env-var
  - [x] using cli framework
- [ ] config file (`TBD how`)
- [x] shell autocompletion
- [ ] testing
- [ ] replace logrus with zap or zerolog
- [ ] fix image-search regex to incluse CRDs
- [ ] release 1.0.0

---

## Known issues

- CRD image references not recognized (to be fixed in the next release)
