
document.addEventListener("DOMContentLoaded", async function() {
    const response = await fetch("/api/logout");
    if (response && response.status === 200) {
        document.getElementById("logout-status").innerText = "La déconnexion a réussi.";
    } else {
        document.getElementById("logout-status").innerText = "Un problème est survenu lors de la déconnexion.";
    }
});