function backPress() {
    history.back();
};

document.onload = function () {
    checkTheme()
}

//TODO: Does not work 
document.addEventListener('checkTheme', function(evt){
    checkTheme()
})

function checkTheme() {
    if (localStorage.theme === 'dark'
        || (!('theme' in localStorage)
            && window.matchMedia('(prefers-color-scheme: dark)').matches)) {
        document.documentElement.classList.add('dark')
    } else {
        document.documentElement.classList.remove('dark')
    }
}