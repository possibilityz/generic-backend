# Authentication Modules - Getting Started

1. Google OIDC
* Initialize `config.yml` using config.example as reference
* If doing development on _local_ environment, add the following line on go.mod file `replace example.com => ./`
* `go mod init example.com`
* `go run cmd/google/main.go`