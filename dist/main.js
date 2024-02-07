function backPress() {
   history.back();
};

document.onload = function () {
   checkTheme()
}

//TODO: Does not work 
document.addEventListener('checkTheme', function (evt) {
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

function allowDrop(ev) {
   ev.preventDefault()
}

function upload(ev, url) {
   ev.preventDefault()
   file = ev.dataTransfer.files[0]

   Swal.showLoading()

   const formData = new FormData()
   formData.append("file", file)

   fetch(url, {
      method: "POST",
      body: formData,
      contentType: "multipart/form-data"
   }).then((response) => {
      return response.text()
   }).then(afterUpload);
}

function afterUpload(rawData) {
   data = JSON.parse(rawData)
   elm = document.getElementById("form-content")
   imgMark = `![img](${data.data.url})`
   if (elm.selectionStart || elm.selectionStart == '0') {
      var startPos = elm.selectionStart;
      var endPos = elm.selectionEnd;

      elm.value = elm.value.substring(0, startPos)
         + imgMark
         + elm.value.substring(endPos, elm.value.length);
   }
   else {
      elm.value += imgMark
   }

   Swal.close()
}
