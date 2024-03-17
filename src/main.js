const colorSchemeQuery = window.matchMedia('(prefers-color-scheme: dark)');

colorSchemeQuery.addEventListener('change', () => {
   window.initialState()
});


window.addEventListener('load', () => {
   window.initialState()
});

window.initialState = function() {
   const darkToggle = document.getElementById('dark-toggle')
   if (darkToggle) {
      darkToggle.checked = !colorSchemeQuery.matches
   }
   const darkToggleMobile = document.getElementById('dark-toggle-mobile')
   if (darkToggleMobile) {
      darkToggle.checked = !colorSchemeQuery.matches
   }

   window.toggleDarkMode(darkToggle.checked)
   window.setMkTheme(null, null)
}

window.toggleDarkMode = function(isChecked) {
   if (!isChecked) {
      document.documentElement.classList.add("dark")
      document.documentElement.classList.remove("light")
   } else {
      document.documentElement.classList.remove("dark")
      document.documentElement.classList.add("light")
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
            theme = 'night'
         }
      }
      localStorage.setItem('mk-theme-dark', theme)
   } else {
      if (theme === null || theme === undefined) {
         theme = localStorage.getItem('mk-theme-light')
         if (theme === null || theme === undefined) {
            theme = 'fantasy'
         }
      }
      localStorage.setItem('mk-theme-light', theme)
   }

   const baseElement = document.documentElement;
   baseElement.setAttribute('data-theme', theme)
}
