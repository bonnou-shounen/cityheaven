package cmd

type optGlobal struct {
	Debug bool `hidden:"" env:"CITYHEAVEN_DEBUG"`
}

type optAuth struct {
	LoginID  string `short:"l" name:"login" env:"CITYHEAVEN_LOGIN"`
	Password string `short:"p" env:"CITYHEAVEN_PASSWORD"`
}

type optDump struct {
	JSON bool `hidden:"" help:"as JSON"`
}

type optShop struct {
	Area string `env:"CITYHEAVEN_AREA" default:"tokyo" help:"area part in shop URL"`
	Shop string `xor:"shop-url" help:"name part in shop URL"`
	URL  string `xor:"shop-url" help:"the shop URL"`
}

type optDumpCast struct {
	NoFav bool `help:"skip counting favorites"`
}

type CLI struct {
	optGlobal
	Dump struct {
		optDump
		Fav struct {
			optAuth
			Casts struct {
				optDumpCast
				DumpFavoriteCasts `cmd:""`
			} `cmd:""`
		} `cmd:""`
		Follow struct {
			optAuth
			Casts struct {
				optDumpCast
				DumpFollowingCasts `cmd:""`
			} `cmd:""`
		} `cmd:""`
		Shop struct {
			optShop
			Casts struct {
				optDumpCast
				DumpShopCasts `cmd:""`
			} `cmd:""`
		} `cmd:""`
	} `cmd:""`
	Restore struct {
		Fav struct {
			optAuth
			Casts struct {
				RestoreFavoriteCasts `cmd:""`
			} `cmd:""`
		} `cmd:""`
	} `cmd:""`
	Version PrintVersion `cmd:"" hidden:""`
}
