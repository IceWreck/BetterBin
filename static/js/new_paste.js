const toUrlEncoded = (obj) =>
    Object.keys(obj)
        .map((k) => encodeURIComponent(k) + "=" + encodeURIComponent(obj[k]))
        .join("&");

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
            burn: "0",
            discuss: "0",
        }),
    })
        .then(function (response) {
            return response.json();
        })
        .then(function (data) {
            console.log(data);
            if ("paste_id" in data) {
                let pasteLink =
                    "http://" +
                    window.location.host +
                    "/paste/view/" +
                    data["paste_id"];
                message = `Your paste has been created at <a href="${pasteLink}" class="alert-link">${pasteLink}</a>.`;
                newAlert("success", message);
            } else if ("error" in data) {
                newAlert("danger", `Error: ${data["error"]}.`);
            } else {
                newAlert("danger", "An unknown error occurred.");
            }
        });
};

const newAlert = (alertType, message) => {
    const alert = `
        <div class="alert alert-${alertType}" role="alert">
        ${message}
        </div>
        `;
    document
        .getElementById("information-alerts")
        .insertAdjacentHTML("afterbegin", alert);
};
