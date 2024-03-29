# Buntdb-cli

an interactive [buntdb](https://github.com/tidwall/buntdb) shell client

![ci](https://github.com/Sora233/buntdb-cli/workflows/ci/badge.svg)
[![Coverage Status](https://coveralls.io/repos/github/Sora233/buntdb-cli/badge.svg?branch=master)](https://coveralls.io/github/Sora233/buntdb-cli?branch=master)

----

### Install

* Download from [Release](https://github.com/Sora233/buntdb-cli/releases)

### Build from source

*It's **recommended** to use binary release*

- go >= 1.13

```shell
go get -u -v github.com/Sora233/buntdb-cli
```

### Usage

**WARN: DO NOT use write command when other buntdb program is running, as multi write can destroy the buntdb file**

![Demo](https://user-images.githubusercontent.com/11474360/104103798-07fae580-52df-11eb-8030-e5d5ff3d80fe.jpg)

* Support Command
    * get
    * set
    * del
    * ttl
    * rbegin (begin a readonly transaction)
    * rwbegin (begin a read/write transaction)
    * commit
    * rollback
    * show
    * keys
    * search
    * use
    * shrink
    * save

You can provide -h flag for command to print help message.
![get](https://user-images.githubusercontent.com/11474360/104104364-81e09e00-52e2-11eb-8863-391420bf6064.jpg)

### TODO

- [ ] create index (Index is memory-only, You need to create index everytime you restart, so I am considering whether to impl it)
