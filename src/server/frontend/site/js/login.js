function generateAccountKey() {
    fetch('/api/register', {
        method: 'GET'
    })
    .then(response => response.json())
    .then(data => {
        const modal = document.getElementById('registerModal');
        const keyElement = document.getElementById('newAccountKey');
        keyElement.textContent = `Votre nouvelle clé : ${data.key}`;
        modal.style.display = 'flex';

        document.getElementById('saveKeyButton').onclick = () => {
            navigator.clipboard.writeText(data.key)
                .then(() => alert('Clé copiée dans le presse-papiers !'))
                .catch(() => alert('Échec de la copie de la clé.'));
        };
    })
    .catch(error => {
        console.error('Error:', error);
        showErrorModal('Échec de la génération de la clé.');
    });
}

function offerRedirectIfAlreadyLoggedIn(loggedIn) {
    if (loggedIn) {
        showLoggedInModal();
    }
}

function showLoggedInModal() {
    const modal = document.getElementById('loggedInModal');
    const countdownElement = document.getElementById('countdown');
    let countdown = 3;
    modal.style.display = 'flex';

    const countdownInterval = setInterval(() => {
        countdown--;
        countdownElement.textContent = countdown;

        if (countdown <= 0) {
            clearInterval(countdownInterval);
            window.location.replace("/v1/dashboard"); 
        }
    }, 1000);
}

function closeModal(modalId) {
    document.getElementById(modalId).style.display = 'none';
}

document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    
    const accountKey = document.getElementById('account_key').value;
    const rememberMe = document.getElementById('rememberMe').checked;

    const data = { 
        account_key: accountKey,
        remember_me: rememberMe
    };

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