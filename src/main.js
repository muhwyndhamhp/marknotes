import _hyperscript from 'hyperscript.org';

_hyperscript.browserInit()

const colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');

colorSchemeQuery.addEventListener('change', () => {
   window.initialState()
});

window.addEventListener('load', () => {
  if (localStorage.getItem('failed-hx-req')) {
    window.navigateFailedReq()
  }
  window.initialState()
});

window.initialState = function() {
  console.log(localStorage.getItem('dark-mode'))

   const darkToggle = document.getElementById('dark-toggle')
   if (darkToggle) {
     darkToggle.checked = !(localStorage.getItem('dark-mode') === 'true') && !colorSchemeQuery.matches
   }

   window.toggleDarkMode(darkToggle.checked, darkToggle.checked === !colorSchemeQuery.matches)
   window.setMkTheme(null, !darkToggle.checked)
}

window.toggleDarkMode = function(isChecked, save = true) {
   if (!isChecked) {
      document.documentElement.classList.add("dark")
      document.documentElement.classList.remove("light")
   } else {
      document.documentElement.classList.remove("dark")
      document.documentElement.classList.add("light")
   }

  if (save && (isChecked === !colorSchemeQuery.matches)) {
    localStorage.removeItem('dark-mode')
  } else if (save){
    localStorage.setItem('dark-mode', !isChecked)
  }
}

window.setMkTheme = function(theme, isDark) {
   console.log(isDark)
   if (isDark === null) {
      isDark = colorSchemeQuery.matches
   }
   if (isDark) {
      if (theme === null || theme === undefined) {
         theme = localStorage.getItem('mk-theme-dark')
         if (theme === null || theme === undefined) {
            theme = 'sunset'
         }
      }
      localStorage.setItem('mk-theme-dark', theme)
   } else {
      if (theme === null || theme === undefined) {
         theme = localStorage.getItem('mk-theme-light')
         if (theme === null || theme === undefined) {
            theme = 'cupcake'
         }
      }
      localStorage.setItem('mk-theme-light', theme)
   }

   const baseElement = document.documentElement;
   baseElement.setAttribute('data-theme', theme)
}
