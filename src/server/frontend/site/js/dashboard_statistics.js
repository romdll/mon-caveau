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

async function SetupStatisticsPage() {
    const charts = {
        wineTypes: echarts.init(document.getElementById('wineTypes')), 
        regions: echarts.init(document.getElementById('regions')),
        vintages: echarts.init(document.getElementById('vintages')),
        domains: echarts.init(document.getElementById('domains')),
        transactions: echarts.init(document.getElementById('transactions'))
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
        transactions: {
            title: {
                text: 'Suivi des Transactions et du Stock de Vins',
                subtext: 'Achats, Ventes, Consommations et Stock',
                left: 'center'
            },
            tooltip: {
                trigger: 'axis',
                formatter: function (params) {
                    let result = `${params[0].name}<br/>`;
                    params.forEach(param => {
                        result += `${param.seriesName}: ${param.value} bouteilles<br/>`;
                    });
                    return result;
                }
            },
            legend: {
                data: ['Achats', 'Ventes', 'Consommations', 'Stock'],
                orient: 'horizontal', 
                bottom: 10,            
                left: 'center'        
            },
            xAxis: {
                type: 'category',
                data: ['2024-01-01', '2024-01-02', '2024-01-03', '2024-01-04', '2024-01-05'], 
            },
            yAxis: {
                type: 'value'
            },
            series: [
                {
                    name: 'Achats',
                    type: 'line',
                    data: [30, 50, 40, 60, 80] 
                },
                {
                    name: 'Ventes',
                    type: 'line',
                    data: [10, 0, 20, 5, 0]
                },
                {
                    name: 'Consommations',
                    type: 'line',
                    data: [5, 10, 0, 12, 7] 
                },
                {
                    name: 'Stock',
                    type: 'line',
                    data: (function () {
                        let stock = 100; 
                        const stockData = [];
                        const purchases = [30, 50, 40, 60, 80];
                        const sales = [10, 0, 20, 5, 0]; 
                        const consumptions = [5, 10, 0, 12, 7];
                        
                        for (let i = 0; i < purchases.length; i++) {
                            stock += purchases[i] - sales[i] - consumptions[i];
                            stockData.push(stock);
                        }
                        return stockData;
                    })()
                }
            ]
        }
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
