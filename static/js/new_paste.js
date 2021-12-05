const newPaste = () => {
    const title = document.getElementById("input-paste-title").value;
    const password = document.getElementById("input-paste-password").value;
    const content = document.getElementById("input-paste-content").value;
    const preview = document.getElementById("input-paste-view-type").value;
    const expiry = document.getElementById("input-paste-expiry").value;

    if (content.length < 1) {
        newAlert("danger", "Paste content cannot be empty.");
        return;
    }

    fetch("/paste/new", {
        method: "post",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: toUrlEncoded({
            content: content,
            title: title,
            expiry: expiry,
            password: password,
            discuss: "0",
        }),
    })
        .then(function (response) {
            return response.json();
        })
        .then(function (data) {
            console.log(data);
            if ("id" in data) {
                // console.log(preview)
                let previewURLString = "";
                if (preview == "code") {
                    previewURLString = "?lang=";
                } else if (preview == "markdown") {
                    previewURLString = "?preview=markdown";
                }

                let pasteLink =
                    "http://" +
                    window.location.host +
                    "/paste/view/" +
                    data["id"] +
                    previewURLString;

                // open the view paste page
                window.location.href = pasteLink

                // or you can create a new alert with the link

                // message = `Your paste has been created at <a href="${pasteLink}" class="alert-link">${pasteLink}</a>.`;
                // newAlert("success", message);
            } else if ("error" in data) {
                newAlert("danger", `Error: ${data["error"]}.`);
            } else {
                newAlert("danger", "An unknown error occurred.");
            }
        }).catch((error) => {
        console.log(error);
        newAlert("danger", "An unknown error occurred.");
    });
};

// submit on Ctrl + Enter (can also press the create button)
document.body.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && e.ctrlKey) {
        newPaste();
    }
});
