window.headerUpload = async function(ev, url) {
    ev.preventDefault()
    if(ev.dataTransfer.files.length === 0) {
        return
    }

    let file = ev.dataTransfer.files[0]

    Swal.showLoading()

    const formData = new FormData()
    formData.append("file", file)

    return fetch(url + "?size=600", {
        method: "POST",
        body: formData,
        contentType: "multipart/form-data"
    })
        .then((response) => {
            return response.text()
        })
}
