package cmd

type Option struct {
	Login    string `short:"l" help:"mail address"`
	Password string `short:"p" help:"password"`
}

type Arg struct {
	Option
	Dump struct {
		Fav struct {
			Casts DumpFavoriteCasts `cmd:""`
		} `cmd:""`
		Shop struct {
			Casts DumpShopCasts `cmd:""`
		} `cmd:""`
	} `cmd:""`
	Restore struct {
		Fav struct {
			Casts RestoreFavoriteCasts `cmd:""`
		} `cmd:""`
	} `cmd:""`
	Version PrintVersion `cmd:"" hidden:""`
}
