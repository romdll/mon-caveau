function timeAgo(date) {
    const now = new Date();
    const targetDate = new Date(date);
    const diffMs = targetDate - now;
    const diffSeconds = Math.abs(Math.floor(diffMs / 1000));
    const diffMinutes = Math.floor(diffSeconds / 60);
    const diffHours = Math.floor(diffMinutes / 60);
    const diffDays = Math.floor(diffHours / 24);
    const diffMonths = Math.floor(diffDays / 30); 
    const diffYears = Math.floor(diffMonths / 12);

    const isFuture = diffMs > 0;

    if (diffYears > 0) return isFuture ? `dans ${diffYears} an${diffYears > 1 ? "s" : ""}` : `il y a ${diffYears} an${diffYears > 1 ? "s" : ""}`;
    if (diffMonths > 0) return isFuture ? `dans ${diffMonths} mois` : `il y a ${diffMonths} mois`;
    if (diffDays > 0) return isFuture ? `dans ${diffDays} jour${diffDays > 1 ? "s" : ""}` : `il y a ${diffDays} jour${diffDays > 1 ? "s" : ""}`;
    if (diffHours > 0) return isFuture ? `dans ${diffHours} heure${diffHours > 1 ? "s" : ""}` : `il y a ${diffHours} heure${diffHours > 1 ? "s" : ""}`;
    if (diffMinutes > 0) return isFuture ? `dans ${diffMinutes} minute${diffMinutes > 1 ? "s" : ""}` : `il y a ${diffMinutes} minute${diffMinutes > 1 ? "s" : ""}`;
    return isFuture ? `dans quelques secondes` : `il y a quelques secondes`;
}

async function SetupAccountPage() {
    const accountData = {
        accountKey: "ACC12345XYZ",
        email: "utilisateur@exemple.com",
        name: "Jean",
        surname: "Dupont",
    };

    const sessionData = [
        {
            id: 1,
            createdAt: "2025-01-01T10:00:00",
            lastActivity: "2025-01-05T10:30:00",
            expiresAt: "2025-01-10T10:00:00"
        },
        {
            id: 2,
            createdAt: "2025-01-03T12:00:00",
            lastActivity: "2025-01-05T11:00:00",
            expiresAt: "2025-01-08T12:00:00"
        }
    ];

    document.getElementById("accountKey").value = accountData.accountKey;
    document.getElementById("email").value = accountData.email;
    document.getElementById("name").value = accountData.name;
    document.getElementById("surname").value = accountData.surname;

    const tbody = document.querySelector("#sessionsTable tbody");
    sessionData.forEach(session => {
        const row = document.createElement("tr");
        row.innerHTML = `
            <td>${session.id}</td>
            <td>${timeAgo(session.createdAt)}</td>
            <td>${timeAgo(session.lastActivity)}</td>
            <td>${timeAgo(session.expiresAt)}</td>
            <td><button onclick="deleteSession(${session.id})">Supprimer</button></td>
        `;
        tbody.appendChild(row);
    });
}

function deleteSession(sessionId) {
    alert(`Session ${sessionId} supprimée.`);
}