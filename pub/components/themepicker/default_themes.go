package pub_themepicker

type ThemeOption struct {
	Label  string
	Theme  string
	IsDark bool
}

var DefaultThemes = []ThemeOption{
	{
		Label:  "Cupcake",
		Theme:  "cupcake",
		IsDark: false,
	},
	{
		Label:  "Sunset",
		Theme:  "sunset",
		IsDark: true,
	},
}
