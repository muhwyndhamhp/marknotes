const colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');

// Listen for changes in the color scheme preference
colorSchemeQuery.addEventListener('change', (event) => {
   setMkTheme(event.matches)
});

window.addEventListener('load', () => {
   setMkTheme(colorSchemeQuery.matches)
});

function setMkTheme(isDark, theme) {
   if (isDark) {
      theme = localStorage.getItem('mk-theme-dark')
      if (theme === null) {
         theme = 'night'
         localStorage.setItem('mk-theme-dark', theme)
      }
   } else {
      theme = localStorage.getItem('mk-theme-light')
      if (theme === null) {
         theme = 'fantasy'
         localStorage.setItem('mk-theme-light', theme)
      }
   }

   const baseElement = document.documentElement;
   baseElement.setAttribute('data-theme', theme)
}
