const urlParams = new URLSearchParams(window.location.search);
const preview = urlParams.get("preview");
if (preview === "markdown") {
    document.getElementById("paste-preview").classList.remove("view-plaintext");
    document.getElementById("paste-preview").innerHTML = marked(pasteContents);
}
