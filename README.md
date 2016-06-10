# gcredstash

## Description

This is a port of [CredStash](https://github.com/fugue/credstash) to Go.

gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

## Usage

```
usage: gcredstash [--version] [--help] <command> [<args>]

Available commands are:
    delete    Delete a credential from the store
    get       Get a credential from the store
    getall    Get all credentials from the store
    list      list credentials and their version
    put       Put a credential into the store
    setup     setup the credential store
```

```
$ gcredstash -h delete
usage: gcredstash delete [-v VERSION] credential

$ gcredstash -h env
usage: gcredstash env [-v VERSION] [-p PREFIX] credential [context [context ...]]

$ gcredstash -h get
usage: gcredstash get [-v VERSION] [-n] credential [context [context ...]]

$ gcredstash -h getall
usage: gcredstash getall [-v VERSION] [context [context ...]]

$ gcredstash -h list
usage: gcredstash list

$ gcredstash -h put
usage: gcredstash put [-k KEY] [-v VERSION] [-a] credential value [context [context ...]]

$ gcredstash -h setup
usage: credstash setup
```

## Set to environment variables

```
$ gcredstash get xxx.*
{
  "xxx.xxx": "100",
  "xxx.yyy": "200"
}

$ gcredstash env xxx.*
export XXX_YYY=200
export XXX_XXX=100

$ gcredstash env xxx.* -p xxx.
export YYY=200
export XXX=100

$ eval $(gcredstash env xxx.*)
```

## Installation

see https://github.com/winebarrel/gcredstash/releases.

### OS X

```sh
brew install https://raw.githubusercontent.com/winebarrel/gcredstash/master/homebrew/gcredstash.rb
```

### Ubuntu

```sh
wget -q -O- https://github.com/winebarrel/gcredstash/releases/download/vN.N.N/gcredstash_N.N.N_amd64.deb | dpkg -i -
```

## Environment variables

```sh
export AWS_REGION=...
export AWS_ACCESS_KEY_ID=...
export AWS_SECRET_ACCESS_KEY=...

# default: credential-store
#export GCREDSTASH_TABLE=...


# default: alias/credstash
#export GCREDSTASH_KMS_KEY"),
```
