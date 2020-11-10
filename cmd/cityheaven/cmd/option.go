package cmd

type Option struct {
	Login    string `short:"l" help:"Mail address."`
	Password string `short:"p" help:"Password."`
}

//nolint:maligned
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
