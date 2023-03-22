# Install

{{< hint "info" >}}

This document walks through the installation of `kwokctl` and `kwok` binaries.

{{< /hint >}}

## Homebrew

On Linux/MacOS systems you can install kwok/kwokctl via [brew](https://formulae.brew.sh/formula/kwok):

``` bash
brew install kwok
```

## Binary Releases

### Install `kwokctl`

``` bash
wget -O kwokctl -c "https://github.com/{{< variable "repo" >}}/releases/download/v{{< variable "version" >}}/kwokctl-$(go env GOOS)-$(go env GOARCH)"
chmod +x kwokctl
sudo mv kwokctl /usr/local/bin/kwokctl
```

### Install `kwok`

``` bash
wget -O kwok -c "https://github.com/{{< variable "repo" >}}/releases/download/v{{< variable "version" >}}/kwok-$(go env GOOS)-$(go env GOARCH)"
chmod +x kwok
sudo mv kwok /usr/local/bin/kwok
```
