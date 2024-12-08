document.getElementById('toggle-btn').addEventListener('click', function() {
    document.getElementById('sidebar').classList.toggle('active');
    this.classList.toggle('open');
});

// Example content change based on sidebar selection
document.getElementById('dashboardLink').addEventListener('click', function() {
    setActiveMenu(this);
    document.getElementById('contentArea').innerHTML = `
        <h2>Bienvenue sur le Tableau de bord</h2>
        <p>Voici les informations importantes sur votre caveau.</p>
        <div class="example-data">
            <h3>Exemple de données</h3>
            <ul>
                <li>Total de vins: 120</li>
                <li>Total de vins dégustés: 75</li>
                <li>Vins en attente de dégustation: 45</li>
            </ul>
        </div>
    `;
});

document.getElementById('collectionLink').addEventListener('click', function() {
    setActiveMenu(this);
    document.getElementById('contentArea').innerHTML = `
        <h2>Votre Collection de Vins</h2>
        <p>Voici une vue détaillée de votre collection de vins.</p>
        <div class="example-data">
            <h3>Exemple de Collection</h3>
            <ul>
                <li>Vin 1: Bordeaux 2015</li>
                <li>Vin 2: Champagne Brut 2018</li>
                <li>Vin 3: Pinot Noir 2017</li>
            </ul>
        </div>
    `;
});

document.getElementById('statsLink').addEventListener('click', function() {
    setActiveMenu(this);
    document.getElementById('contentArea').innerHTML = `
        <h2>Statistiques</h2>
        <p>Analyse des données de votre caveau.</p>
        <div class="example-data">
            <h3>Exemple de Statistiques</h3>
            <ul>
                <li>Vins par région: Bordeaux (50), Bourgogne (30), Champagne (20)</li>
                <li>Vins par type: Rouge (60), Blanc (40), Rosé (20)</li>
            </ul>
        </div>
    `;
});

// Function to set active menu item
function setActiveMenu(selectedLink) {
    const items = document.querySelectorAll('.sidebar-item');
    items.forEach(item => item.classList.remove('active')); // Remove active from all
    selectedLink.classList.add('active'); // Add active to clicked link
}
