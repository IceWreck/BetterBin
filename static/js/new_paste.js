// Create a new paste.
// If openImmidiately is true then the user will be redirected to new paste.
// Otherwise a green/success alertbox will show new link.
const newPaste = (openImmidiately) => {
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


                // if paste is burn on reading then open immidiately would render it useless
                if (!openImmidiately || expiry === "burn") {

                    // create a new alert with the link

                    message = `Your paste has been created at <a href="${pasteLink}" class="alert-link">${pasteLink}</a>.`;
                    newAlert("success", message);

                } else {
                    // open the view paste page
                    window.location.href = pasteLink
                }

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

// Submit on Ctrl + Enter (can also press the create button).
document.body.addEventListener('keydown', (e) => {
    if (e.key === 'Enter' && e.ctrlKey) {
        newPaste(true);
    }
});


// On pressing tab, add relevant space.
// This is still different from a real tab because real tab gets you to next
// nearest multiple of TAB_SIZE and not just add TAB_SIZE to cursor postiton.
const contentArea = document.getElementById("input-paste-content");
contentArea.onkeydown = (event) => {
    if (event.key === 'Tab') {
        console.log("tabekey pressed")
        event.preventDefault();
        // add 1 tab = 4 spaces
        contentArea.value += "    ";
        contentArea.focus();
    }
}
