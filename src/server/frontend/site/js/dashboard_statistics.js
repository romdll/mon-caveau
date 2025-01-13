let wineTypeByRegion = {};

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
        if (tx.type === "drank") groupedData[date].drunk += tx.quantity;
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

function updateWineTypeChart(region, parentRegion) {
    let filteredData;
    let titleText;

    if (region === "") {
        filteredData = Object.entries(wineTypeByRegion).reduce((acc, [country, regions]) => {
            Object.entries(regions).forEach(([regionName, wineTypes]) => {
                Object.entries(wineTypes).forEach(([wineTypeName, value]) => {
                    const existing = acc.find(item => item.name === wineTypeName);
                    if (existing) {
                        existing.value += value;
                    } else {
                        acc.push({ name: wineTypeName, value });
                    }
                });
            });
            return acc;
        }, []);

        titleText = 'Distribution des types de vins';
    } else if (parentRegion && parentRegion !== "") {
        const regionData = Object.entries(wineTypeByRegion).reduce((acc, [country, regions]) => {
            if (country === parentRegion && regions[region]) {
                Object.entries(regions[region]).forEach(([wineTypeName, value]) => {
                    acc.push({ name: wineTypeName, value });
                });
            }
            return acc;
        }, []);

        filteredData = regionData;
        titleText = `Distribution des types de vins - ${region} (${parentRegion})`;
    } else {
        const regionData = Object.entries(wineTypeByRegion).reduce((acc, [country, regions]) => {
            if (country === region) {
                Object.entries(regions).forEach(([regionName, wineTypes]) => {
                    Object.entries(wineTypes).forEach(([wineTypeName, value]) => {
                        const existing = acc.find(item => item.name === wineTypeName);
                        if (existing) {
                            existing.value += value;
                        } else {
                            acc.push({ name: wineTypeName, value });
                        }
                    });
                });
            }
            return acc;
        }, []);

        filteredData = regionData;
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
                data: filteredData,
                radius: '65%',
                center: ['50%', '55%']
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

    const wineStatisticsRequest = await fetch("/api/wines/statistics/raw");
    const wineStatistics = await wineStatisticsRequest.json();

    const { top5Domains, wineDistributionPerVintage, wineTypesDistributionPerRegion, userUsedRegionsWithBottlecount } = wineStatistics;

    const groupedByCountry = userUsedRegionsWithBottlecount.reduce((acc, { region, bottle_count }) => {
        const { country, name } = region;

        if (!acc[country]) {
            acc[country] = [];
        }

        acc[country].push({ name, value: bottle_count });

        return acc;
    }, {});

    const regionHierarchy = Object.entries(groupedByCountry).map(([country, regions]) => ({
        name: country,
        children: regions
    }));

    wineTypeByRegion = wineTypesDistributionPerRegion.reduce((acc, { region, wine_types }) => {
        const { country, name: regionName } = region;

        if (!acc[country]) {
            acc[country] = {};
        }

        if (!acc[country][regionName]) {
            acc[country][regionName] = {};
        }

        wine_types.forEach(wineType => {
            const wineTypeName = wineType.name;
            if (!acc[country][regionName][wineTypeName]) {
                acc[country][regionName][wineTypeName] = 0;
            }
            acc[country][regionName][wineTypeName] += wineType.count;
        });

        return acc;
    }, {});

    const options = {
        wineTypes: (wineStatistics && wineTypesDistributionPerRegion) ? {
            title: {
                text: 'Distribution des Types de Vins',
                subtext: 'Répartition des vins selon leur type et leur nombre de bouteilles',
                left: 'center',
            },
            tooltip: { trigger: 'item' },
        } : emptyChart(),
        regions: (wineStatistics && userUsedRegionsWithBottlecount) ? {
            title: {
                text: 'Représentation Territoriale',
                subtext: 'Vue d\'ensemble des régions viticoles et du nombre de bouteilles',
                left: 'center'
            },
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
        } : emptyChart(),
        vintages: (wineStatistics && wineDistributionPerVintage) ? {
            title: {
                text: 'Distributions des Vins par Millésime',
                subtext: 'Répartition des vins selon leur année de production',
                left: 'center'
            },
            tooltip: {},
            xAxis: {
                type: 'category', data: Object.entries(wineDistributionPerVintage).map(([key]) => {
                    return key;
                })
            },
            yAxis: { type: 'value' },
            series: [
                {
                    type: 'bar',
                    data: Object.entries(wineDistributionPerVintage).map(([_, value]) => {
                        return value;
                    })
                }
            ]
        } : emptyChart(),
        domains: (wineStatistics && top5Domains) ? {
            title: {
                text: 'Top 5 Domaines par Nombre de Bouteilles',
                subtext: 'Classement des domaines en fonction du nombre total de bouteilles disponibles',
                left: 'center'
            },
            tooltip: {
                trigger: 'item',
                formatter: '{b}: {c} bouteilles ({d}%)'
            },
            series: [
                {
                    type: 'pie',
                    radius: ['25%', '35%'],
                    label: {
                        formatter: '{b} ({c})',
                        fontSize: 12,
                        color: '#333'
                    },
                    data: Object.entries(top5Domains).map(([key, value]) => {
                        return { name: key, value: value };
                    })
                }
            ]
        } : emptyChart(),
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
                    data: achats,
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
        const ancestors = params.treeAncestors;
        let parentRegion = null;

        if (ancestors && ancestors.length > 1) {
            parentRegion = ancestors[ancestors.length - 2].name;
        }

        const selectedRegion = params.name ? params.name : (params.nodeData ? params.nodeData.name : "");
        updateWineTypeChart(selectedRegion, parentRegion);
    });
}
