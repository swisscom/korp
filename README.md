# korp
A command line tool for pushing docker images into a corporate registry based on Kubernetes yaml files

## Installation

1. Download the [latest release](https://github.com/swisscom/korp/releases) and unpack it
2. Copy the `korp` binary to your path

## Usage

### Scan your yaml files

Calling `korp scan -f <path to yaml files> -r <corporate registry name>` will scan your yaml files for image references and will create a `kustomization.yaml` file specifying the necessary replacements.
