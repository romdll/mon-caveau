// Creates nice colors but a bit the same and maybe this time too much flashy 
function generateColorFromString(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        hash = str.charCodeAt(i) + ((hash << 5) - hash);
    }

    const hue = Math.abs(hash) % 360;
    const saturation = 70 + (Math.abs(hash) % 30);
    const lightness = 50 + (Math.abs(hash) % 20);

    return `hsl(${hue}, ${saturation}%, ${lightness}%)`;
}

async function SetupDashboardPage() {
    const response = await fetch("/api/wines/basic");
    const json = await response.json();

    if (json["error"]) {
        const error = json["error"];
    }

    const winesCountPerRegions = json["winesCountPerRegions"];
    const winesCountPerTypes = json["winesCountPerTypes"];
    const last4Transactions = json["last4Transactions"];

    document.getElementById("totalWines").innerText = json["totalWines"] + " (" + json["totalCurrentBottles"] + ")";
    document.getElementById("addedWines").innerText = json["realTotalBottlesAdded"];
    document.getElementById("tastedWines").innerText = json["totalWinesDrankSold"];

    const regionData = Object.keys(winesCountPerRegions).map(region => ({
        name: region,
        value: winesCountPerRegions[region],
    }));

    const regionColorPalette = Object.keys(winesCountPerRegions).map(generateColorFromString);

    const regionChartElement = document.getElementById('regionChart')
    const regionChart = echarts.init(regionChartElement);

    const regionChartOption = regionData.length > 0 ? {
        title: {
            text: "Répartition des Bouteilles par Région / Département",
            textStyle: {
                color: '#000000'
            },
            left: 'center'
        },
        tooltip: {
            confine: true,
            trigger: 'item',
            formatter: '{b}: {c} bouteilles ({d}%)',
        },
        legend: {
            show: false
        },
        series: [
            {
                name: 'Nombre de bouteilles',
                type: 'pie',
                radius: '50%',
                data: regionData,
                // color: regionColorPalette,
                label: {
                    formatter: '{b}: {c} ({d}%)',
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    } : {
        title: {
            textStyle: {
                color: "grey",
                fontSize: 20
            },
            text: "Aucune donnée",
            left: "center",
            top: "center"
        }
    };

    if (regionData.length === 0) {
        regionChartElement.className = 'no-data-chart';
    }
    regionChart.setOption(regionChartOption);

    const typeChartElement = document.getElementById('typeChart')
    const typeChart = echarts.init(typeChartElement);

    const typeData = Object.keys(winesCountPerTypes).map(type => ({
        name: type,
        value: winesCountPerTypes[type],
    }));

    const typeColorPalette = Object.keys(winesCountPerTypes).map(generateColorFromString);

    const typeChartOption = typeData.length > 0 ? {
        title: {
            text: "Répartition des Bouteilles par Type",
            textStyle: {
                color: '#000000'
            },
            left: 'center'
        },
        tooltip: {
            confine: true,
            trigger: 'item',
            formatter: '{b}: {c} bouteilles ({d}%)',
        },
        legend: {
            show: false,
        },
        series: [
            {
                name: 'Nombre de bouteilles',
                type: 'pie',
                radius: '50%',
                data: typeData,
                // color: typeColorPalette,
                label: {
                    formatter: '{b}: {c} ({d}%)',
                },
                emphasis: {
                    itemStyle: {
                        shadowBlur: 10,
                        shadowOffsetX: 0,
                        shadowColor: 'rgba(0, 0, 0, 0.5)'
                    }
                }
            }
        ]
    } : {
        title: {
            textStyle: {
                color: "grey",
                fontSize: 20
            },
            text: "Aucune donnée",
            left: "center",
            top: "center"
        }
    };

    if (typeData.length === 0) {
        typeChartElement.className = 'no-data-chart';
    }
    typeChart.setOption(typeChartOption);

    new ResizeObserver(() => regionChart.resize()).observe(regionChartElement);
    new ResizeObserver(() => typeChart.resize()).observe(typeChartElement);

    const transactionContainer = document.querySelector('.recent-transactions ul');
    transactionContainer.innerHTML = "";
    if (last4Transactions["data"] && last4Transactions["data"].length > 0) {
        last4Transactions["data"].forEach(transaction => {
            const wineName = last4Transactions["winesIdToName"][transaction.wine_id];

            const transactionTypeClass = transaction.type === 'added' ? 'add' : 'sell';
            const transactionIcon = transaction.type === 'added' ? 'fa-plus-circle' : 'fa-minus-circle';

            const transactionDate = new Date(transaction.date);
            const formattedDate = transactionDate.toLocaleString("fr-FR", {
                weekday: 'long',
                year: 'numeric',
                month: 'long',
                day: 'numeric',
                hour: '2-digit',
                minute: '2-digit',
                second: '2-digit',
                hour12: false,
            });

            const quantityWord = transaction.quantity > 1 ? 'bouteilles' : 'bouteille';
            const actionWord = transaction.type === 'added'
                ? (transaction.quantity > 1 ? 'ajoutées' : 'ajoutée')
                : (transaction.type === 'sell' || transaction.type === 'sold' ? (transaction.quantity > 1 ? 'vendues' : 'vendue') : 'dégustée');

            const transactionHTML = `
                <li class="transaction-item ${transactionTypeClass}">
                    <span class="icon"><i class="fas ${transactionIcon}"></i></span>
                    <div class="transaction-info">
                        <strong>${transaction.quantity} ${quantityWord}</strong> de ${wineName} 
                        ${actionWord} (${formattedDate})
                    </div>
                </li>
            `;

            transactionContainer.innerHTML += transactionHTML;
        });
    } else {
        transactionContainer.innerHTML = '<div class="no-data">Aucune transaction récente</div>';
    }
}