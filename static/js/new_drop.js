const newDrop = () => {
    newAlert("primary", "Uploading...");

    const title = document.getElementById("title-input").value;
    let uploadedFile = document.getElementById("upload-file").files[0];
    let formData = new FormData();
    formData.append("title", title);
    formData.append("upload", uploadedFile);

    fetch("/drop/new", {
        method: "post",
        headers: {
            Accept: "application/json",
            // "Content-Type": "multipart/form-data",
        },
        body: formData,
    })
        .then(function (response) {
            return response.json();
        })
        .then(function (data) {
            console.log(data);
            if ("id" in data) {
                let downloadLink =
                    "http://" +
                    window.location.host +
                    "/drop/dl/" +
                    data["id"];
                message = `The uploaded file is at <a href="${downloadLink}" class="alert-link">${downloadLink}</a>.`;
                newAlert("success", message);
            } else if ("error" in data) {
                newAlert("danger", `Error: ${data["error"]}.`);
            } else {
                newAlert("danger", "An unknown error occurred.");
            }
        });
};
