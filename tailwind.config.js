const defaultTheme = require('tailwindcss/defaultTheme')

module.exports = {
  content: ["./public/**/*.{html,js}"],
  theme: {
    screens: {
        'xs': '320px',
      ...defaultTheme.screens,
    },
    fontFamily: {
      'sans': ['JetBrains Mono', 'sans-serif'],
      'serif': ['Spicy Kebab', 'serif'],
      'mono': ['JetBrains Mono', 'mono'],
    },
    extend: {},
  },
  plugins: [],
}