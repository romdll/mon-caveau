const addWineModal = document.getElementById('addWineModal');
const closeModal = document.getElementById('closeModal');
const createByHandBtn = document.getElementById('createByHandBtn');
const scanBottleBtn = document.getElementById('scanBottleBtn');
let activeInput = null;
let suggestionList = null

function addWine() {
    showWineAddOptionsModal();
}

function showWineAddOptionsModal() {
    addWineModal.style.display = 'flex';
    addWineModal.classList.add('show');
}

function hideWineAddOptionsModal() {
    addWineModal.classList.remove('show');
    setTimeout(() => {
        addWineModal.style.display = 'none';
    }, 300); 
}

closeModal.addEventListener('click', hideWineAddOptionsModal);

createByHandBtn.addEventListener('click', function () {
    hideWineAddOptionsModal();
    openHardCreationModal();
});

scanBottleBtn.addEventListener('click', function () {
    alert('Scanner la bouteille sera disponible bientôt.');
});

window.addEventListener('click', function (event) {
    if (event.target === addWineModal) {
        hideWineAddOptionsModal();
    }

    if (event.target !== activeInput && suggestionList) {
        suggestionList.style.display = 'none';
    }
});

function openHardCreationModal() {
    document.getElementById("addHardWineModal").style.display = "flex";
    document.getElementById("addHardWineModal").classList.add('show');
}

function closeHardCreationModal() {
    document.getElementById("addHardWineModal").classList.remove('show');
    setTimeout(() => {
        document.getElementById("addHardWineModal").style.display = 'none';
    }, 300); 
}

async function fetchDomains(query) {
    if (activeInput !== document.getElementById("domainInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const domains = ['Domain A', 'Domain B', 'Domain C'];
    const suggestions = filterSuggestions(domains, query);

    updateSuggestions('domainSuggestions', suggestions, query);
    setActiveInput(document.getElementById("domainInput"));
}

async function fetchRegions(query) {
    if (activeInput !== document.getElementById("regionInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const regions = ['Region A', 'Region B', 'Region C'];
    const suggestions = filterSuggestions(regions, query);

    updateSuggestions('regionSuggestions', suggestions, query);
    setActiveInput(document.getElementById("regionInput"));
}

async function fetchTypes(query) {
    if (activeInput !== document.getElementById("typeInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const types = ['Rouge 1', 'Blanc B', 'Rosé C', 'Rouge 2'];
    const suggestions = filterSuggestions(types, query);

    updateSuggestions('typeSuggestions', suggestions, query);
    setActiveInput(document.getElementById("typeInput"));
}

async function fetchBottleSizes(query) {
    if (activeInput !== document.getElementById("bottleSizeInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const bottleSizes = ['750ml', '1L', '1.5L'];
    const suggestions = filterSuggestions(bottleSizes, query);

    updateSuggestions('bottleSizeSuggestions', suggestions, query);
    setActiveInput(document.getElementById("bottleSizeInput"));
}

function filterSuggestions(allSuggestions, query) {
    if (!query || query && query === "") {
        return allSuggestions;
    }
    return allSuggestions.filter(data => data.includes(query));
}

function highlightMatch(suggestion, query) {
    if (!query) {
        return suggestion;
    }

    const startIdx = suggestion.toLowerCase().indexOf(query.toLowerCase());
    if (startIdx === -1) return suggestion;

    const beforeMatch = suggestion.slice(0, startIdx);
    const matchText = suggestion.slice(startIdx, startIdx + query.length);
    const afterMatch = suggestion.slice(startIdx + query.length);

    if (matchText === "") {
        return afterMatch;
    }

    return `${beforeMatch}<span class="highlight">${matchText}</span>${afterMatch}`;
}

function updateSuggestions(datalistId, suggestions, query) {
    suggestionList = document.getElementById(datalistId);
    suggestionList.innerHTML = ''; 
    suggestionList.style.display = suggestions.length > 0 ? 'block' : 'none';

    suggestions.forEach(suggestion => {
        const suggestionItem = document.createElement('div');
        suggestionItem.classList.add('suggestion-item');

        const highlightedSuggestion = highlightMatch(suggestion, query);

        suggestionItem.innerHTML = highlightedSuggestion;
        suggestionItem.addEventListener('click', () => selectSuggestion(suggestion));

        suggestionList.appendChild(suggestionItem);
    });
}

function selectSuggestion(suggestion) {
    activeInput.value = suggestion;
    suggestionList.style.display = 'none';
}

function setActiveInput(inputElement) {
    activeInput = inputElement;
    suggestionList.style.display = 'block';
    suggestionList.style.width = `${activeInput.offsetWidth}px`;
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
    closeHardCreationModal();
}

function toggleMoreFields() {
    const extraFields = document.getElementById('extraFields');
    const btn = document.getElementById('toggleMoreFieldsBtn');
    
    if (extraFields.style.display === 'none' || extraFields.style.display === '') {
        extraFields.style.display = 'block';
        btn.innerHTML = 'Informations supplémentaires &#x2191;';
    } else {
        extraFields.style.display = 'none';
        btn.innerHTML = 'Informations supplémentaires &#x2193;';
    }
}