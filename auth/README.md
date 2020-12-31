# Authentication Modules - Getting Started

1. Google OIDC
* Initialize `config.yml` using config.example as reference for the _google_ field
* If doing development on _local_ environment, add the following line on go.mod file `replace example.com => ./`
* `go mod init example.com`
* `go run main.go`

2. Github OIDC
* Initialize `config.yml` using config.example as reference for the _github_ field