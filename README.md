# cityheaven

cli for [cityheaven](https://www.cityheaven.net).

## install

```bash
$ go install github.com/bonnou-shounen/cityheaven/cmd/cityheaven@latest
```

## usage

```bash
$ export CITYHEAVEN_LOGIN=xxxx
$ export CITYHEAVEN_PASSWORD=xxxx

$ cityheaven dump fav casts > casts.txt
$ vim casts.txt  # edit order
$ cityheaven restore fav casts < casts.txt
```
