# archgregator
A Cosmos-chain data aggregator for Archway network.

# Installation

## Install Golang

Go 1.16 is required for Archway.

If you haven't already, download and install Go. See the official [go.dev documentation](https://golang.org/doc/install). Make sure your `GOBIN` and `GOPATH` are setup.

## Get the Archway source code

Retrieve the source code from the official [archway-network/archway](https://github.com/archway-network/archway) GitHub repository.

```
git clone https://github.com/NuclearBlock/archgregator
cd archgregator
git checkout main
```

## Build the Archgregator binary

You can build with:

```
make install
```

This command installs the `archgregator` to your `GOPATH`.


## Init home folder and prepare conf file

```
archgregator init
```
This command creates ~/.archgregator folder where you have to place a config.yaml file (Please see config.yaml.example as a reference)

Recommended mode - 'remote' node. (local mode not yet completely ready) 


## Run Postgres

```
docker compose up
```
This command runs docker container with Postgres and creates all nessessary tables 


## Run Parser

```
archgregator start
```
This command runs Archgregator to parse blocks from RPC node.


To use collected data please see our ExpressJS/ReactJS solution - github.com/NuclearBlock/archgregator_front

