function addWine() {
    showModal();
}

const addWineModal = document.getElementById('addWineModal');
const closeModal = document.getElementById('closeModal');
const createByHandBtn = document.getElementById('createByHandBtn');
const scanBottleBtn = document.getElementById('scanBottleBtn');

function showModal() {
    addWineModal.style.display = 'flex';
    addWineModal.classList.add('show');
}

function hideModal() {
    addWineModal.classList.remove('show');
    setTimeout(() => {
        addWineModal.style.display = 'none';
    }, 300); 
}

closeModal.addEventListener('click', hideModal);

createByHandBtn.addEventListener('click', function () {
    hideModal();
    openModal();
});

scanBottleBtn.addEventListener('click', function () {
    alert('Scanner la bouteille sera disponible bientôt.');
});

window.addEventListener('click', function (event) {
    if (event.target === addWineModal) {
        hideModal();
    }
});

function openModal() {
    document.getElementById("addHardWineModal").style.display = "flex";
}

// Close modal
function closeModalFunc() {
    document.getElementById("addHardWineModal").style.display = "none";
}

async function fetchDomains(query) {
    const domains = ['Domain A', 'Domain B', 'Domain C'];
    const suggestions = domains.filter(domain => domain.toLowerCase().includes(query.toLowerCase()));

    updateSuggestions('domainSuggestions', suggestions);
}

async function fetchRegions(query) {
    const regions = ['Region A', 'Region B', 'Region C'];
    const suggestions = regions.filter(region => region.toLowerCase().includes(query.toLowerCase()));

    updateSuggestions('regionSuggestions', suggestions);
}

async function fetchTypes(query) {
    const types = ['Type A', 'Type B', 'Type C']; 
    const suggestions = types.filter(type => type.toLowerCase().includes(query.toLowerCase()));

    updateSuggestions('typeSuggestions', suggestions);
}

async function fetchBottleSizes(query) {
    const bottleSizes = ['750ml', '1L', '1.5L'];
    const suggestions = bottleSizes.filter(size => size.includes(query));

    updateSuggestions('bottleSizeSuggestions', suggestions);
}

function updateSuggestions(datalistId, suggestions) {
    const datalist = document.getElementById(datalistId);
    datalist.innerHTML = '';
    suggestions.forEach(suggestion => {
        const option = document.createElement('option');
        option.value = suggestion;
        datalist.appendChild(option);
    });
}

function submitWineForm(event) {
    event.preventDefault();

    const wineData = new FormData(document.getElementById('wineForm'));

    if (wineData.get('domain') && !document.getElementById('domainSuggestions').querySelector('option[value="' + wineData.get('domain') + '"]')) {
        wineData.append('newDomain', wineData.get('domain'));
    }

    if (wineData.get('region') && !document.getElementById('regionSuggestions').querySelector('option[value="' + wineData.get('region') + '"]')) {
        wineData.append('newRegion', wineData.get('region'));
    }

    if (wineData.get('type') && !document.getElementById('typeSuggestions').querySelector('option[value="' + wineData.get('type') + '"]')) {
        wineData.append('newType', wineData.get('type'));
    }

    if (wineData.get('bottleSize') && !document.getElementById('bottleSizeSuggestions').querySelector('option[value="' + wineData.get('bottleSize') + '"]')) {
        wineData.append('newBottleSize', wineData.get('bottleSize'));
    }

    console.log("Wine Data Submitted:", wineData);
    closeModalFunc();
}