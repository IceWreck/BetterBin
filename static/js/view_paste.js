const customizeURL = () => {
    const currentURL = window.location.href;
    window.location.href =
        window.location.origin +
        "/shortner/new?" +
        toUrlEncoded({ url: currentURL });
};
