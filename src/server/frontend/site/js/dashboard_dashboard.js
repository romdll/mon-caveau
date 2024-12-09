const regionData = {
    labels: ['Bordeaux', 'Bourgogne', 'Champagne', 'Pinot Noir'],
    datasets: [{
        label: 'Vins par Région',
        data: [50, 30, 20, 10],
        backgroundColor: ['#FF5733', '#C70039', '#900C3F', '#581845'],
        borderColor: ['#FF5733', '#C70039', '#900C3F', '#581845'],
        borderWidth: 1
    }]
};

const typeData = {
    labels: ['Rouge', 'Blanc', 'Rosé'],
    datasets: [{
        label: 'Vins par Type',
        data: [60, 40, 20], 
        backgroundColor: ['#FF9F00', '#33A1FF', '#FF33A1'],
        borderColor: ['#FF9F00', '#33A1FF', '#FF33A1'],
        borderWidth: 1
    }]
};

const ctxRegion = document.getElementById('regionChart').getContext('2d');
const regionChart = new Chart(ctxRegion, {
    type: 'pie',
    data: regionData,
    options: {
        responsive: true,
        plugins: {
            legend: {
                position: 'top',
            },
            tooltip: {
                callbacks: {
                    label: function(tooltipItem) {
                        return tooltipItem.label + ': ' + tooltipItem.raw + ' Vins';
                    }
                }
            }
        }
    }
});

const ctxType = document.getElementById('typeChart').getContext('2d');
const typeChart = new Chart(ctxType, {
    type: 'pie',
    data: typeData,
    options: {
        responsive: true,
        plugins: {
            legend: {
                position: 'top',
            },
            tooltip: {
                callbacks: {
                    label: function(tooltipItem) {
                        return tooltipItem.label + ': ' + tooltipItem.raw + ' Vins';
                    }
                }
            }
        }
    }
});

function addWine() {
    alert('Fonction "Ajouter un Vin" activée. À implémenter plus tard!');
}

async function SetupDashboardPage() {
    const response = await fetch("/api/wines/basic")
    const json = await response.json();

    document.getElementById("totalWines").innerText = json["totalWines"];
    document.getElementById("tastedWines").innerText = json["totalWinesDrankSold"];
    document.getElementById("soldDrunkThisMonth").innerText = json["totalWinesDrankSoldThisMonth"];
}