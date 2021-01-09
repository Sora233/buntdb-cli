# Buntdb-cli

an interactive [buntdb](https://github.com/tidwall/buntdb) shell client

----

### Install

* Download from [Release](https://github.com/Sora233/buntdb-cli/releases)

### Build from source

*It's **recommended** to use binary release*

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
    * show
    * keys
    * use

You can provide -h flag for command to print help message.
![get](https://user-images.githubusercontent.com/11474360/104104364-81e09e00-52e2-11eb-8863-391420bf6064.jpg)

### TODO

- [ ] create index