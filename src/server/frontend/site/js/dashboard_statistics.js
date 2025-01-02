const fakeData = {
    wineTypes: [
        { name: 'Red', value: 40 },
        { name: 'White', value: 30 },
        { name: 'Rosé', value: 20 },
        { name: 'Sparkling', value: 10 }
    ],
    regions: [
        { name: 'France', value: 50 },
        { name: 'Italy', value: 30 },
        { name: 'Spain', value: 20 },
        { name: 'USA', value: 10 }
    ],
    vintages: [1980, 1990, 1995, 2000, 2010, 2015, 2018, 2020],
    domains: [
        { name: 'Chateau Margaux', value: 15, bottles: 300 },
        { name: 'Dom Perignon', value: 12, bottles: 450 },
        { name: 'Sassicaia', value: 10, bottles: 500 },
        { name: 'Penfolds', value: 8, bottles: 250 },
        { name: 'Vega Sicilia', value: 5, bottles: 100 }
    ],
    transactions: [
        [0, 0, 5], [0, 1, 10], [0, 2, 15],
        [1, 0, 20], [1, 1, 25], [1, 2, 30]
    ],
    topValues: [
        { name: 'Sassicaia 2015', value: 500 },
        { name: 'Chateau Latour 2000', value: 450 },
        { name: 'Dom Perignon 2012', value: 400 },
        { name: 'Penfolds Grange 2018', value: 350 },
        { name: 'Opus One 2017', value: 300 }
    ]
};

const regionHierarchy = [
    {
        name: 'France',
        children: [
            { name: 'Bordeaux', value: 20 },
            { name: 'Champagne', value: 15 },
            { name: 'Bourgogne', value: 15 }
        ]
    },
    {
        name: 'Italy',
        children: [
            { name: 'Tuscany', value: 18 },
            { name: 'Piedmont', value: 12 }
        ]
    },
    {
        name: 'Spain',
        children: [
            { name: 'Rioja', value: 10 },
            { name: 'Catalonia', value: 10 }
        ]
    }
];

const wineTypeByRegion = {
    France: { Red: 20, White: 10, Rosé: 5, Sparkling: 5 },
    Bordeaux: { Red: 10, White: 5, Rosé: 2 },
    Champagne: { Sparkling: 15 }
};

const charts = {
    wineTypes: echarts.init(document.getElementById('wineTypes')),
    regions: echarts.init(document.getElementById('regions')),
    vintages: echarts.init(document.getElementById('vintages')),
    domains: echarts.init(document.getElementById('domains')),
    transactions: echarts.init(document.getElementById('transactions'))
};

function processCumulativeTransactions(transactions, accountCreationDate) {
    const groupedData = {};

    const creationDate = accountCreationDate.split(' ')[0]; 

    transactions.forEach(tx => {
        const date = tx.date.split('T')[0];

        if (!groupedData[date]) {
            groupedData[date] = { added: 0, removed: 0, drunk: 0, stock: 0 };
        }

        if (tx.type === "added") groupedData[date].added += tx.quantity;
        if (tx.type === "removed") groupedData[date].removed += tx.quantity;
        if (tx.type === "drunk") groupedData[date].drunk += tx.quantity;
        groupedData[date].stock += tx.quantity;
    });

    const dates = [];
    const achats = [];
    const ventes = [];
    const consommations = [];
    const stockData = [];

    if (!groupedData[creationDate]) {
        groupedData[creationDate] = { added: 0, removed: 0, drunk: 0, stock: 0 };
    }

    dates.push(creationDate);
    achats.push(0);
    ventes.push(0);
    consommations.push(0);
    stockData.push(0);

    Object.keys(groupedData).sort().forEach(date => {
        const dayData = groupedData[date];
        let cumulativeAdded = achats[achats.length - 1];
        let cumulativeRemoved = ventes[ventes.length - 1];
        let cumulativeDrunk = consommations[consommations.length - 1];
        let cumulativeStock = stockData[stockData.length - 1];

        cumulativeAdded += dayData.added;
        cumulativeRemoved += dayData.removed;
        cumulativeDrunk += dayData.drunk;
        cumulativeStock += dayData.added - dayData.removed - dayData.drunk;

        dates.push(date);
        achats.push(cumulativeAdded);
        ventes.push(cumulativeRemoved);
        consommations.push(cumulativeDrunk);
        stockData.push(cumulativeStock);
    });

    return { dates, achats, ventes, consommations, stockData };
};

function updateWineTypeChart(region) {
    let filteredData;
    let titleText;

    if (region === "") {
        filteredData = Object.entries(wineTypeByRegion).reduce((acc, [regionKey, wineTypes]) => {
            Object.entries(wineTypes).forEach(([name, value]) => {
                const existing = acc.find(item => item.name === name);
                if (existing) {
                    existing.value += value;
                } else {
                    acc.push({ name, value });
                }
            });
            return acc;
        }, []);

        titleText = 'Distribution des types de vins';
    } else {
        filteredData = Object.entries(wineTypeByRegion[region] || {}).map(([name, value]) => ({ name, value }));
        titleText = `Distribution des types de vins - ${region}`;
    }

    charts.wineTypes.setOption({
        title: {
            text: titleText,
            left: 'center'
        },
        series: [
            {
                type: 'pie',
                data: filteredData
            }
        ]
    });
}

function emptyChart() {
    return {
        title: {
            textStyle: {
                color: "grey",
                fontSize: 20
            },
            text: "Aucune donnée",
            left: "center",
            top: "center"
        }
    }
}

async function SetupStatisticsPage() {
    const wineTransactionsRequest = await fetch("/api/wines/transactions");
    const wineTransactions = await wineTransactionsRequest.json();

    let dates = [], achats = [], ventes = [], consommations = [], stockData = [];
    if (wineTransactions && wineTransactions.data && wineTransactions.accountCreationDate) {
        ({ dates, achats, ventes, consommations, stockData } = processCumulativeTransactions(wineTransactions.data, wineTransactions.accountCreationDate));
    }

    const options = {
        wineTypes: {
            title: { text: 'Distribution des Types de Vins', left: 'center' },
            tooltip: { trigger: 'item' },
        },
        regions: {
            title: { text: 'Représentation Territoriale', left: 'center' },
            tooltip: {
                trigger: 'item',
                formatter: '{b}: {c}'
            },
            series: [
                {
                    type: 'treemap',
                    data: regionHierarchy,
                    label: {
                        show: true,
                        formatter: function (params) {
                            if (params.data.children) {
                                const totalValue = params.data.children.reduce((sum, child) => sum + child.value, 0);
                                return `${params.data.name} (${totalValue} bouteilles)`;
                            } else {
                                return `${params.data.name}: ${params.value}`;
                            }
                        },
                        position: 'inside',
                        fontSize: 12,
                        color: '#ffffff'
                    },
                    upperLabel: {
                        show: true,
                        formatter: function (params) {
                            if (params.data.children) {
                                if (!params.data.name) {
                                    return `Total Global (${params.data.children.reduce((sum, child) => sum + child.value, 0)} bouteilles)`;
                                }
                                return `${params.data.name} (${params.data.children.reduce((sum, child) => sum + child.value, 0)} bouteilles)`;
                            }
                            return '';
                        },
                    }
                }
            ]
        },
        vintages: {
            title: { text: 'Distributions des Vins par Millésime', left: 'center' },
            tooltip: {},
            xAxis: { type: 'category', data: fakeData.vintages },
            yAxis: { type: 'value' },
            series: [
                {
                    type: 'bar',
                    data: fakeData.vintages.map(() => Math.floor(Math.random() * 50))
                }
            ]
        },
        domains: {
            title: {
                text: 'Top 5 Domaines par Nombre de Bouteilles',
                left: 'center'
            },
            tooltip: {
                trigger: 'item',
                formatter: '{b}: {c} bouteilles ({d}%)'
            },
            series: [
                {
                    type: 'pie',
                    radius: ['40%', '70%'],
                    label: {
                        formatter: '{b} ({c})',
                        fontSize: 14,
                        color: '#333'
                    },
                    data: fakeData.domains.map(d => ({
                        name: d.name,
                        value: d.bottles
                    }))
                }
            ]
        },
        transactions: (wineTransactions && wineTransactions.data && dates.length > 0) ? {
            title: {
                text: 'Suivi des Transactions et du Stock de Vins',
                subtext: 'Cumul des Achats, Ventes, Consommations et Stock',
                left: 'center'
            },
            tooltip: {
                trigger: 'axis',
                confine: true, 
                formatter: function (params) {
                    let result = `<div style="padding: 10px; font-family: Arial, sans-serif; font-size: 14px;">`;
        
                    const formattedDate = new Intl.DateTimeFormat('fr-FR').format(new Date(params[0].name));
                    result += `<h4 style="margin: 0; color: #333;">Date: ${formattedDate}</h4><br/>`;
        
                    let stockValue = null;
        
                    params.forEach(param => {
                        if (param.seriesName === 'Stock') {
                            stockValue = param.value;
                        }
                    });
        
                    if (stockValue !== null) {
                        result += `<h4 style="margin: 0; color: #000;">Stock: ${stockValue} bouteilles</h4><br/>`;
                    }
        
                    params.forEach(param => {
                        if (param.seriesName !== 'Stock' && param.value !== 0) {
                            result += `<div style="color: #555;">${param.seriesName}: <strong>${param.value}</strong> bouteilles</div>`;
                        }
                    });
        
                    result += `</div>`;
                    return result;
                }
            },
            legend: {
                data: ['Achats (Cumulatif)', 'Ventes (Cumulatif)', 'Consommations (Cumulatif)', 'Stock'],
                orient: 'horizontal',
                bottom: 10,
                left: 'center'
            },
            dataZoom: [
                {
                    type: 'slider',
                    start: 0,
                    end: 100,
                    height: 20,
                    bottom: 50,
                },
                {
                    type: 'inside',
                    start: 0,
                    end: 100,
                }
            ],
            xAxis: {
                type: 'category',
                data: dates,
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: 'Achats (Cumulatif)',
                    type: 'line',
                    data: achats
                },
                {
                    name: 'Ventes (Cumulatif)',
                    type: 'line',
                    data: ventes
                },
                {
                    name: 'Consommations (Cumulatif)',
                    type: 'line',
                    data: consommations
                },
                {
                    name: 'Stock',
                    type: 'line',
                    data: stockData
                }
            ]
        } : emptyChart()
    };

    updateWineTypeChart("");

    Object.keys(charts).forEach(chart => {
        charts[chart].setOption(options[chart]);

        const observer = new ResizeObserver(() => charts[chart].resize());
        observer.observe(document.getElementById(chart));
    });

    charts.regions.on('click', function (params) {
        const selectedRegion = params.name ? params.name : (params.nodeData ? params.nodeData.name : "");
        updateWineTypeChart(selectedRegion);
    });
}
