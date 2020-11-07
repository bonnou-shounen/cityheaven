# cityheaven

cli for [cityheaven](https://www.cityheaven.net).

## install

```bash
$ go get -u github.com/bonnou-shounen/cityheaven/cmd/cityheaven
```

## usage

```bash
$ export CITYHEAVEN_LOGIN=xxxx
$ export CITYHEAVEN_PASSWORD=xxxx

$ cityheaven dump fav casts > casts.txt
$ vim casts.txt  # edit order
$ cityheaven restore fav casts < casts.txt
```
