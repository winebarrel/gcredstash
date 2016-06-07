# gcredstash

## Description

This is a port of [CredStash](https://github.com/fugue/credstash).

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

$ gcredstash -h get
usage: gcredstash get [-v VERSION] credential [context [context ...]]

$ gcredstash -h getall
usage: gcredstash getall [-v VERSION] [context [context ...]]

$ gcredstash -h list
usage: gcredstash list

$ gcredstash -h put
usage: gcredstash put [-k KEY] [-v VERSION] [-a] credential value [context [context ...]]

$ gcredstash -h setup
usage: credstash setup
```
