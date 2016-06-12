# gcredstash

## Description

This is a port of [CredStash](https://github.com/fugue/credstash) to Go.

gcredstash manages credentials using AWS Key Management Service (KMS) and DynamoDB.

[![Build Status](https://travis-ci.org/winebarrel/gcredstash.svg?branch=master)](https://travis-ci.org/winebarrel/gcredstash)

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
usage: gcredstash get [-v VERSION] [-n] credential [context [context ...]]

$ gcredstash -h getall
usage: gcredstash getall [context [context ...]]

$ gcredstash -h list
usage: gcredstash list

$ gcredstash -h put
usage: gcredstash put [-k KEY] [-v VERSION] [-a] credential value [context [context ...]]

$ gcredstash -h setup
usage: credstash setup
```

## Example

```
$ gcredstash put foo.bar 100
foo.bar has been stored

$ gcredstash put foo.baz 200
foo.baz has been stored

$ gcredstash get foo.bar
100

$ gcredstash get foo.*
{
  "foo.bar": "100",
  "foo.baz": "200"
}
```

```
// DynamoDB data
> select all * from credential-store \G
[
  {
    "contents": "wlpc",
    "hmac": "a925335f7f313e400ed54702f739f1f4ffddd6ff1722fa9ac1e2b6d4e24d5096",
    "key": "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMWB1+YqVMNVT+V5dGAgEQgFtj6aGqRmg+wJwDGPk1kRduGoX6rtyUhm116wSmkQA2SXdPzAr2NcY02/joiiqzu534QQSwpOF+oKIkfLXaaNZCCWQkki94EE+EpkiVeFxcoqAdIaHf7FzwKz1A",
    "name": "foo.baz",
    "version": "0000000000000000001"
  },
  {
    "contents": "yUBx",
    "hmac": "cf6a6ef2458356996ac26de9bf384acce400a367b4d00a42e0e4dd44c8560b99",
    "key": "CiDY1vsR456LEdoL3+0p+PrTCleoqi/sutbDfJZNiUSpphLLAQEBAQB42Nb7EeOeixHaC9/tKfj60wpXqKov7LrWw3yWTYlEqaYAAACiMIGfBgkqhkiG9w0BBwaggZEwgY4CAQAwgYgGCSqGSIb3DQEHATAeBglghkgBZQMEAS4wEQQMccvp6R6qUho35bCEAgEQgFumGPEIHX7B2KgU6S2vaoEOJKX84pGKe0ydMh1r+rMWEZGd5si61FZ76YlgM0X6rnO5qlLK6SGUHhA0whzi7R7Zpbc9euBXYWFYQeMRU9jpDh7H/bhP2fa7BtNV",
    "name": "foo.bar",
    "version": "0000000000000000001"
  }
]
```

## Put from stdin

```
$ echo 300 | gcredstash put xxx.zzz -
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

### CentOS

```sh
yum install https://github.com/winebarrel/gcredstash/releases/download/vN.N.N/gcredstash-N.N.N-x.el6.x86_64.rpm
```

## Setup

* `IAM > Encryption Keys`
  * Create Encryption Key: `Alias`: `credstash`
* Run `gcredstash setup`

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
