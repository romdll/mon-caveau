const addWineModal = document.getElementById('addWineModal');
const closeModal = document.getElementById('closeModal');
const createByHandBtn = document.getElementById('createByHandBtn');
const scanBottleBtn = document.getElementById('scanBottleBtn');
let activeInput = null;
let suggestionList = null;
let intervalListsUpdate = null;

let jsonDomains = null;
let jsonTypes = null;
let jsonBottleSizes = null;
let jsonRegionCountries = null;

let selectedRegion = null;
let selectedCountry = null;

function addWine() {
    showWineAddOptionsModal();
}

function showWineAddOptionsModal() {
    addWineModal.style.display = 'flex';
    addWineModal.style.opacity = '1';
    addWineModal.classList.add('show');
}

function hideWineAddOptionsModal() {
    addWineModal.classList.remove('show');
    addWineModal.style.opacity = '0';
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

async function fetchAll() {
    const responseDomains = await fetch("/api/wines/domains");
    jsonDomains = await responseDomains.json();

    const responseTypes = await fetch("/api/wines/types");
    jsonTypes = await responseTypes.json();

    const responseRegionsCountries = await fetch("/api/wines/countries/regions");
    jsonRegionCountries = await responseRegionsCountries.json();

    if (document.getElementById("countryInput").value === "") {
        const countries = [...new Set(jsonRegionCountries.map(item => item.country))];

        if (countries.length === 1) {
            selectCountry(countries[0]);
        }
    }

    const responseBottleSizes = await fetch("/api/wines/bottles/sizes");
    jsonBottleSizes = await responseBottleSizes.json();
}

function openHardCreationModal() {
    const form = document.getElementById('wineForm');
    form.reset();
    
    document.getElementById("addHardWineModal").style.display = "flex";
    document.getElementById("addHardWineModal").classList.add('show');
    document.getElementById("addHardWineModal").style.opacity = '1';
    fetchAll();
    intervalListsUpdate = setInterval(fetchAll, 60000);
}

function closeHardCreationModal() {
    if (intervalListsUpdate) {
        clearInterval(intervalListsUpdate);
    }

    document.getElementById("addHardWineModal").classList.remove('show');
    document.getElementById("addHardWineModal").style.opacity = '0';
    setTimeout(() => {
        document.getElementById("addHardWineModal").style.display = 'none';
    }, 300);
}

async function fetchDomains(query) {
    if (activeInput !== document.getElementById("domainInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const domains = jsonDomains.map(item => item.name);
    const suggestions = filterSuggestions(domains, query);

    updateSuggestions('domainSuggestions', suggestions, query);
    setActiveInput(document.getElementById("domainInput"));
}

async function fetchRegions(query) {
    if (activeInput !== document.getElementById("regionInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const regionsWithCountries = jsonRegionCountries.map(item => ({
        region: item.name,
        country: item.country
    }));

    const regionCountrySuggestions = regionsWithCountries.map(item => `${item.region} - ${item.country}`);

    const suggestions = filterSuggestions(regionCountrySuggestions, query);
    updateSuggestions('regionSuggestions', suggestions, query);

    setActiveInput(document.getElementById("regionInput"));
}

function selectRegion(regionWithCountry) {
    const [region, country] = regionWithCountry.split(' - ');

    selectedRegion = region;
    document.getElementById("regionInput").value = region;

    const availableCountries = jsonRegionCountries.filter(item => item.name === region);
    if (availableCountries.length === 1) {
        selectedCountry = availableCountries[0].country;
        document.getElementById("countryInput").value = selectedCountry;
    } else {
        selectedCountry = null;
    }

    document.getElementById("regionSuggestions").style.display = 'none';
}

async function fetchCountries(query) {
    if (activeInput !== document.getElementById("countryInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const countries = [...new Set(jsonRegionCountries.map(item => item.country))];
    const suggestions = filterSuggestions(countries, query);

    updateSuggestions('countrySuggestions', suggestions, query);
    setActiveInput(document.getElementById("countryInput"));
}

function selectCountry(country) {
    selectedCountry = country;
    document.getElementById("countryInput").value = country;
    document.getElementById("countrySuggestions").style.display = 'none';

    selectedRegion = null;
    document.getElementById("regionInput").value = "";
}

async function fetchTypes(query) {
    if (activeInput !== document.getElementById("typeInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const types = jsonTypes.map(item => item.name);
    const suggestions = filterSuggestions(types, query);

    updateSuggestions('typeSuggestions', suggestions, query);
    setActiveInput(document.getElementById("typeInput"));
}

async function fetchBottleSizes(query) {
    if (activeInput !== document.getElementById("bottleSizeInput") && suggestionList) {
        suggestionList.style.display = 'none';
    }

    const bottleSizes = jsonBottleSizes.map(item => ({ name: item.name, size: item.size }));

    const suggestions = filterSuggestions(bottleSizes.map(item => `${item.name} (${item.size}ml)`), query);

    updateSuggestions('bottleSizeSuggestions', suggestions, query);
    setActiveInput(document.getElementById("bottleSizeInput"));
}

function updateSliderValue(size) {
    document.getElementById("sliderValueLabel").textContent = `${size}ml`;
    document.getElementById("bottleSizeValue").value = size;
}

function selectSuggestionForBottleSizes(suggestion) {
    const regex = /(.*)\s\((\d+)ml\)/;
    const match = suggestion.match(regex);

    if (match) {
        const name = match[1];
        const size = parseInt(match[2]);

        activeInput.value = name;

        document.getElementById("bottleSizeValue").value = size;

        document.getElementById("bottleSizeSlider").value = size;
        document.getElementById("sliderValueLabel").textContent = `${size}ml`;

        suggestionList.style.display = 'none';
    }
}

function filterSuggestions(allSuggestions, query) {
    if (!query || query && query === "") {
        return allSuggestions;
    }
    return allSuggestions.filter(data => data.toLowerCase().includes(query.toLowerCase()));
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

    if (suggestions.length > 0) {
        suggestions.forEach(suggestion => {
            const suggestionItem = document.createElement('div');
            suggestionItem.classList.add('suggestion-item');

            const highlightedSuggestion = highlightMatch(suggestion, query);

            suggestionItem.innerHTML = highlightedSuggestion;
            suggestionItem.addEventListener('click', () => {
                if (datalistId === 'regionSuggestions') {
                    selectRegion(suggestion);
                } else if (datalistId === 'countrySuggestions') {
                    selectCountry(suggestion);
                } else if (datalistId === 'bottleSizeSuggestions') {
                    selectSuggestionForBottleSizes(suggestion);
                } else {
                    selectSuggestion(suggestion);
                }
            });

            suggestionList.appendChild(suggestionItem);
        });
        suggestionList.style.display = 'block';
    } else {
        suggestionList.style.display = 'none';
    }
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

function openSizeInputModal() {
    const preciseSizeInput = document.getElementById("preciseSizeInput");
    preciseSizeInput.value = document.getElementById("bottleSizeSlider").value;
    document.getElementById("sizeInputModal").style.zIndex = 101;
    document.getElementById("sizeInputModal").style.display = "flex";
    document.getElementById("sizeInputModal").style.opacity = '1';
    document.getElementById("sizeInputModal").classList.add('show');
}

function closeSizeInputModal() {
    document.getElementById("sizeInputModal").classList.remove('show');
    document.getElementById("sizeInputModal").style.opacity = '0';
    setTimeout(() => {
        document.getElementById("sizeInputModal").style.display = "none";
    }, 300);
}

function updateSliderFromInput(value) {
    const slider = document.getElementById("bottleSizeSlider");
    slider.value = value;
    updateSliderValue(value);
}

function updateSliderValue(value) {
    document.getElementById("sliderValueLabel").textContent = `${value}ml`;
    document.getElementById("bottleSizeValue").value = value;
}

function updateBottleSize() {
    const preciseSizeInput = document.getElementById("preciseSizeInput");
    const newSize = preciseSizeInput.value;

    document.getElementById("bottleSizeSlider").value = newSize;
    updateSliderValue(newSize);

    closeSizeInputModal();
}

function getDomainData() {
    const domainName = document.getElementById('domainInput').value;
    const domain = jsonDomains.find(item => item.name === domainName);
    return domain ? { "id": domain.id } : { "name": domainName };
}

function getRegionData() {
    const regionName = document.getElementById('regionInput').value;
    const countryName = document.getElementById('countryInput').value;

    if (regionName && countryName) {
        const region = jsonRegionCountries.find(item => item.name === regionName && item.country === countryName);
        return region ? { "id": region.id } : { "name": regionName, "country": countryName }; 
    }
    return { "name": regionName, "country": countryName }; 
}

function getTypeData() {
    const typeName = document.getElementById('typeInput').value;
    const type = jsonTypes.find(item => item.name === typeName);
    return type ? { "id": type.id } : { "name": typeName };
}

function getBottleSizeData() {
    const bottleSizeName = document.getElementById('bottleSizeInput').value;
    const bottleSizeValue = parseInt(document.getElementById('bottleSizeValue').value);
    const bottleSize = jsonBottleSizes.find(item => item.size === bottleSizeValue);

    return bottleSize ? { "id": bottleSize.id } : { "name": bottleSizeName, "size": bottleSizeValue }; 
}

async function submitWineForm(event) {
    event.preventDefault();

    jsonToSend = {
        "name": document.getElementById('wineName').value, 
        "domain": getDomainData(),
        "region": getRegionData(),
        "type": getTypeData(), 
        "bottle_size": getBottleSizeData(), 
        "vintage": parseInt(document.getElementById('vintage').value), 
        "quantity": parseInt(document.getElementById('quantity').value),
        "buy_price": parseFloat(document.getElementById('buyPrice').value) || null,
        "description": document.getElementById('description').value || null, 
        "image": document.getElementById('image').value || null 
    };

    const response = await fetch("/api/wines/create", {
        method: "POST",
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify(jsonToSend)
    });
    if (response.status === 200) {
        closeHardCreationModal();
    } else {
        console.log(response);
        console.log(await response.json());
    }
}