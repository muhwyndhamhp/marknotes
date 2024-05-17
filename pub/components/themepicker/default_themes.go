package pub_themepicker

type ThemeOption struct {
	Label  string
	Theme  string
	IsDark bool
}

var DefaultThemes = []ThemeOption{
	{
		Label:  "Autumn",
		Theme:  "autumn",
		IsDark: false,
	},
	{
		Label:  "Sunset",
		Theme:  "sunset",
		IsDark: true,
	},
}
