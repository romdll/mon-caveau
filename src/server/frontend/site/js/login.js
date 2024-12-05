function generateAccountKey() {
    alert("Une nouvelle clé d'accès a été générée. (Ceci est simulé pour le moment.)");
}

document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const accountKey = document.getElementById('account_key').value;

    const data = { account_key: accountKey };

    fetch('/api/login', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(data)
    })
    .then(response => response.json())
    .then(data => {
        window.location.replace("/site/dashboard")
    })
    .catch((error) => {
        console.error('Error:', error);
    });
});