# cityheaven

cli for [cityheaven](https://www.cityheaven.net).

## install

```bash
$ go install github.com/bonnou-shounen/cityheaven/cmd/cityheaven@latest
```

## usage

sort "my girls" order
```bash
$ export CITYHEAVEN_LOGIN=xxxx
$ export CITYHEAVEN_PASSWORD=xxxx

$ cityheaven dump fav casts > fav-casts.txt
$ vim fav-casts.txt  # edit order
$ cityheaven restore fav casts < fav-casts.txt
```

popular 10 casts on the shop
```bash
$ cityheaven dump shop casts --shop=xxxx | sort -k3nr | head -10
```

new 5 casts there
```
$ cityheaven dump shop casts --shop=yyyy --no-fav | sort -k1nr | head -5
```

casts you may talk with
```bash
$ cityheaven dump follow casts --mutual > favtalk-casts.txt
```
