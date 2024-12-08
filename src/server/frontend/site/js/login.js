function generateAccountKey() {
    alert("Une nouvelle clé d'accès a été générée. (Ceci est simulé pour le moment.)");
}

function offerRedirectIfAlreadyLoggedIn(loggedIn) {
    if (loggedIn) {
        showLoggedInModal();
    }
}

function showLoggedInModal() {
    // Show modal
    const modal = document.getElementById('loggedInModal');
    const countdownElement = document.getElementById('countdown');
    let countdown = 3;
    modal.style.display = 'flex';

    // Update countdown every second
    const countdownInterval = setInterval(() => {
        countdown--;
        countdownElement.textContent = countdown;

        if (countdown <= 0) {
            clearInterval(countdownInterval);
            window.location.replace("/dashboard");
        }
    }, 1000);
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
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
    .then(response => {
        if (response.ok) {
            return response.json();
        } else {
            throw new Error('Clé de compte invalide.');
        }
    })
    .then(data => {
        const modalTitle = document.getElementById('loggedInTitle');
        modalTitle.textContent = "Votre compte est valide !"
        showLoggedInModal();
    })
    .catch((error) => {
        console.error('Error:', error);
        showErrorModal(error.message); 
    });
});

function showErrorModal(message) {
    const modal = document.getElementById('errorModal');
    modal.querySelector('p').textContent = message;
    modal.style.display = 'flex';
}