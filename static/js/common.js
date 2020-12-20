const toUrlEncoded = (obj) =>
    Object.keys(obj)
        .map((k) => encodeURIComponent(k) + "=" + encodeURIComponent(obj[k]))
        .join("&");

const newAlert = (alertType, message) => {
    const alert = `
                <div class="alert alert-${alertType}" role="alert">
                ${message}
                </div>
                `;
    document
        .getElementById("information-alerts")
        // .insertAdjacentHTML("afterbegin", alert);
        .innerHTML = alert;
};
