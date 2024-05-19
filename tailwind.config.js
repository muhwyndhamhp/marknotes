import { screens as _screens } from 'tailwindcss/defaultTheme';
import plugin from 'tailwindcss/plugin';

export const darkMode = 'class';
export const content = ["./pub/**/*.templ", "./src/**/*.js"];
export const theme = {
    screens: {
        'xs': '320px',
        ..._screens,
    },
    extend: {},
};
export const plugins = [
    plugin(function({ addBase }) {
        addBase({
            'html': {
                fontSize: "16px", // Default font size, which is typically 16px
                '@screen md': {
                    fontSize: "14px", // Font size on medium (md) breakpoint, which is typically 14px
                },
            },
        });
    }),
    require('@tailwindcss/typography'),
    require('daisyui')
];
export const daisyui = {
    themes: [
        "cupcake",
        "sunset",
    ],
    darkTheme: "sunset",
};
