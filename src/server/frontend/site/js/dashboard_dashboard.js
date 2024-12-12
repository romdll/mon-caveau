// Gets a bit lame colors and not easy to recognize
function generateColorFromString(str) {
    let hash = 0;
    for (let i = 0; i < str.length; i++) {
        hash = (hash << 5) - hash + str.charCodeAt(i);
    }

    const red = Math.abs((hash & 0xFF0000) >> 16) % 128 + 128;
    const green = Math.abs((hash & 0x00FF00) >> 8) % 128 + 128;
    const blue = Math.abs(hash & 0x0000FF) % 128 + 128;

    const hexColor = ((1 << 24) | (red << 16) | (green << 8) | blue).toString(16).slice(1).toUpperCase();

    return `#${hexColor}`;
}

async function SetupDashboardPage() {
    const response = await fetch("/api/wines/basic");
    const json = await response.json();

    const winesCountPerRegions = json["winesCountPerRegions"];
    const winesCountPerTypes = json["winesCountPerTypes"];
    const last4Transactions = json["last4Transactions"];

    document.getElementById("totalWines").innerText = json["totalWines"];
    document.getElementById("tastedWines").innerText = json["totalWinesDrankSold"];
    document.getElementById("soldDrunkThisMonth").innerText = json["totalWinesDrankSoldThisMonth"];

    const regionData = Object.keys(winesCountPerRegions).map(region => ({
        name: region,
        value: winesCountPerRegions[region],
    }));

    // const regionColorPalette = Object.keys(winesCountPerRegions).map(generateColorFromString);

    const regionChartElement = document.getElementById('regionChart')
    const regionChart = echarts.init(regionChartElement);

    const regionChartOption = regionData.length > 0 ? {
        title: {
            text: "Répartition des Vins par Région",
            textStyle: {
                color: '#000000'
            },
            left: 'center'
        },
        tooltip: {
            confine: true,
            trigger: 'item'
        },
        legend: {
            show: false,
        },
        series: [
            {
                name: 'Nombre de bouteilles',
                type: 'pie',
                radius: '50%',
                data: regionData,
                // color: regionColorPalette,
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
            text: "Répartition des Vins par Type",
            textStyle: {
                color: '#000000'
            },
            left: 'center'
        },
        tooltip: {
            trigger: 'item'
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
            const wineName = last4Transactions["wineIdToName"][transaction.wine_id];
            
            const transactionTypeClass = transaction.type === 'add' ? 'add' : 'sell';
            const transactionIcon = transaction.type === 'add' ? 'fa-plus-circle' : 'fa-minus-circle';
            
            const transactionDate = new Date(transaction.date);
            const formattedDate = transactionDate.toLocaleDateString("fr-FR");  
    
            const transactionHTML = `
                <li class="transaction-item ${transactionTypeClass}">
                    <span class="icon"><i class="fas ${transactionIcon}"></i></span>
                    <div class="transaction-info">
                        <strong>${transaction.quantity} bouteille${transaction.quantity > 1 ? 's' : ''}</strong> de ${wineName} ${transaction.type === 'add' ? 'ajoutée' : 'vendue / dégustée'} (${formattedDate})
                    </div>
                </li>
            `;
    
            transactionContainer.innerHTML += transactionHTML;
        });
    } else {
        transactionContainer.innerHTML = '<div class="no-data">Aucune transaction récente</div>';
    }
}