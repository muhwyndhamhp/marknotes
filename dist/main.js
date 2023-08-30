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

const url = "https://api.cloudinary.com/v1_1/dmhi9axrr/image/upload";
function allowDrop(ev) {
    ev.preventDefault()
}

function upload(ev) {
    ev.preventDefault()
    console.log(ev)
    file = ev.dataTransfer.files[0]

    console.log(file)

    Swal.showLoading()
    reader = new FileReader();
    reader.onload = function (event) {
        const formData = new FormData()
        formData.append("file", event.target.result)
        formData.append("upload_preset", "unsigned_bucket")

        fetch(url, {
            method: "POST",
            body: formData
        }).then((response) => {
            return response.text()
        }).then(afterUpload);
    }
    reader.readAsDataURL(file)
}

function afterUpload(rawData) {
    data = JSON.parse(rawData)
    elm = document.getElementById("form-content")
    imgMark = `![img](${data.url})`
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