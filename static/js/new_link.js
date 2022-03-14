const newLink = () => {
    const longURL = document.getElementById("shortner-input").value;
    const preferredID = document.getElementById("preferred-id").value;

    if (longURL.length < 1) {
        newAlert("danger", "Input URL cannot be empty.");
        return;
    }

    fetch("/shortner/new", {
        method: "post",
        headers: {
            Accept: "application/json",
            "Content-Type": "application/x-www-form-urlencoded",
        },
        body: toUrlEncoded({
            url: longURL,
            id: preferredID,
        }),
    })
        .then(function (response) {
            return response.json();
        })
        .then(function (data) {
            console.log(data);
            if ("id" in data) {
                let shortLink =
                    "http://" + window.location.host + "/s/" + data["id"];
                message = `The short link is <a href="${shortLink}" class="alert-link">${shortLink}</a>.`;
                newAlert("success", message);
            } else if ("error" in data) {
                newAlert("danger", `Error: ${data["error"]}.`);
            } else {
                newAlert("danger", "An unknown error occurred.");
            }
        })
        .catch((error) => {
            console.log(error);
            newAlert("danger", "An unknown error occurred.");
        });
};
