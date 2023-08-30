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
        }).then((rawData) => {
            data = JSON.parse(rawData)
            document.getElementById("form-content").value += `![img](${data.url})`
            Swal.close()
        });
    }
    reader.readAsDataURL(file)
}