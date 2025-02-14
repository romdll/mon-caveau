let winesPerPage = 6;
let currentPage = 1;
let totalWines = 0;
let searchQuery = "";
let filterPreferredDates = false;

async function togglePreferredDatesFilter() {
    filterPreferredDates = !filterPreferredDates;
    const filterButton = document.getElementById("filterPreferredDatesButton");
    if (filterPreferredDates) {
        filterButton.textContent = "Afficher tous les vins";
        filterButton.classList.add("active");
    } else {
        filterButton.textContent = "Afficher uniquement les vins à leur apogée";
        filterButton.classList.remove("active");
    }
    await loadWines(currentPage);
}

async function fetchAll() {
    const responseDomains = await fetch("/api/wines/domains");
    jsonDomains = await responseDomains.json();

    const responseTypes = await fetch("/api/wines/types");
    jsonTypes = await responseTypes.json();

    const responseRegionsCountries = await fetch("/api/wines/countries/regions");
    jsonRegionCountries = await responseRegionsCountries.json();

    const responseBottleSizes = await fetch("/api/wines/bottles/sizes");
    jsonBottleSizes = await responseBottleSizes.json();
}

async function SetupCollectionPage() {
    document.getElementById("collectionContent").style.display = "block";
    await fetchAll();
    await loadWines(currentPage);
    new ResizeObserver(function () { adjustPaginationForPhone() }).observe(document.body);
}

async function adjustQuantity(wineId, change) {
    const quantityElement = document.querySelector(`.wine-item .quantity[data-id="${wineId}"]`);
    if (!quantityElement) return;

    const originalQuantity = parseInt(quantityElement.textContent, 10) || 0;
    const newQuantity = originalQuantity + change;

    try {
        quantityElement.textContent = newQuantity;

        quantityElement.classList.add('quantity-updated');
        setTimeout(() => quantityElement.classList.remove('quantity-updated'), 500);

        const response = await fetch(`/api/wines/${wineId}/adjust-quantity`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ change }),
        });

        if (!response.ok) {
            throw new Error('Impossible de changer la quantité du vin.');
        }

        const updatedWine = await response.json();
        if (updatedWine.quantity === 1 || updatedWine.quantity === 0) {
            loadWines(currentPage);
            setupPagination();
        } else {
            quantityElement.textContent = updatedWine.quantity;
        }

    } catch (error) {
        console.error('Error adjusting wine quantity:', error);
        quantityElement.textContent = originalQuantity;
        alert(error);
    }
}

async function loadWines(page) {
    const wineListContainer = document.querySelector('.wine-list');
    wineListContainer.innerHTML = '';

    const placeholders = Array.from({ length: winesPerPage }, () => {
        const placeholder = document.createElement('div');
        placeholder.classList.add('wine-item', 'loading-placeholder');
        placeholder.innerHTML = `
            <div class="wine-info">
                <h4 class="loading-placeholder-text">Loading...</h4>
                <p class="loading-placeholder-text"><strong>Domaine:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Type:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Région / Département:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Pays:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Taille:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Quantité:</strong> Loading...</p>
            </div>`;
        wineListContainer.appendChild(placeholder);
        return placeholder;
    });

    try {
        const response = await fetch(`/api/wines?page=${page}&limit=${winesPerPage}&search=${encodeURIComponent(searchQuery)}&filterPreferredDates=${filterPreferredDates}`);
        const data = await response.json();
        const { wines, total } = data;
        totalWines = total;

        if (wines) {
            wines.forEach((wine, index) => {
                const domainName = jsonDomains.find(domain => domain.id === wine.domain_id)?.name || "Unknown Domain";
                const typeName = jsonTypes.find(type => type.id === wine.type_id)?.name || "Unknown Type";
                const region = jsonRegionCountries.find(region => region.id === wine.region_id);
                const regionName = region ? region.name : "Unknown Region";
                const countryName = region ? region.country : "Unknown Country";
                const bottleSize = jsonBottleSizes.find(size => size.id === wine.bottle_size_id);
                const bottleSizeName = bottleSize ? `${bottleSize.name} (${bottleSize.size} ml)` : "Unknown Size";
    
                const isInactive = wine.quantity === 0;
    
                const today = new Date();
                const currentYear = today.getFullYear();
    
                const startDate = wine.preferred_start_date ? new Date(wine.preferred_start_date) : null;
                const endDate = wine.preferred_end_date ? new Date(wine.preferred_end_date) : null;
    
                const startYear = startDate ? startDate.getFullYear() : null;
                const endYear = endDate ? endDate.getFullYear() : null;
    
                const isWithinPreferredDates = startYear && endYear && currentYear >= startYear && currentYear <= endYear;
                const isPastPreferredDates = endYear && currentYear > endYear;

                const wineItem = document.createElement('div');

                wineItem.classList.add('wine-item');
                if (isInactive) {
                    wineItem.classList.add('wine-inactive');
                }
                if (isWithinPreferredDates) {
                    wineItem.classList.add('highlight-preferred');
                }
                if (isPastPreferredDates) {
                    wineItem.classList.add('highlight-past-preferred');
                }
    
                let dateInfo = '';
                if (startYear || endYear) {
                    dateInfo = `
                        <p><strong>Période de Dégustation:</strong>
                            ${startYear ? `À partir de ${startYear}` : ''}
                            ${endYear ? ` jusqu'en ${endYear}` : ''}
                        </p>`;
                }

                wineItem.innerHTML = `
                    <div class="quick-actions">
                        <button class="decrement" onclick="adjustQuantity(${wine.id}, -1)" ${isInactive ? "disabled" : ""}>-1</button>
                        <button class="increment" onclick="adjustQuantity(${wine.id}, 1)">+1</button>
                    </div>
                    <img src="${wine.image || '/v1/images/no_photo_generic.svg'}" alt="Image of ${wine.name}" onerror="this.src='/v1/images/no_photo_generic.svg';">
                    <div class="wine-info">
                        <h4>${wine.name} (${wine.vintage})</h4>
                        <p><strong>Domaine:</strong> ${domainName}</p>
                        <p><strong>Type:</strong> ${typeName}</p>
                        <p><strong>Région / Département:</strong> ${regionName}</p>
                        <p><strong>Pays:</strong> ${countryName}</p>
                        <p><strong>Taille:</strong> ${bottleSizeName}</p>
                        <p><strong>Quantité:</strong> <span class="quantity" data-id="${wine.id}">${wine.quantity}</span></p>
                        ${dateInfo}
                    </div>
                    <div class="wine-item-actions">
                        <button class="edit" onclick="editWine(${wine.id})" ${isInactive ? "disabled" : ""}>Modifier</button>
                        <button class="delete" onclick="askDeleteWine(${wine.id})" ${wine.quantity === 0 ? "disabled" : ""}>Supprimer</button>
                    </div>
                    ${isInactive ? `<span class="wine-badge">Vin Inactif</span>` : ``}
                    ${isWithinPreferredDates ? `<span class="wine-badge preferred-badge">À son apogée</span>` : ``}
                    ${isPastPreferredDates ? `<span class="wine-badge expired-badge">Période de dégustation dépassée</span>` : ``}
                `;
    
                placeholders[index].replaceWith(wineItem);
            });
        }

        if (wines && wines.length < placeholders.length) {
            placeholders.slice(wines.length).forEach(placeholder => placeholder.remove());
        } else if (!wines) {
            placeholders.forEach(placeholder => placeholder.remove());
        }

        setupPagination();
    } catch (error) {
        console.error('Erreur lors du chargement des vins:', error);

        placeholders.forEach(placeholder => placeholder.remove());
    }
}

function setupPagination() {
    const totalPages = Math.ceil(totalWines / winesPerPage);
    const paginationContainer = document.getElementById('pagination');
    paginationContainer.innerHTML = '';

    if (totalPages <= 1) {
        paginationContainer.style.display = "none";
        return
    };

    paginationContainer.style.display = "flex";

    const maxDisplayedPages = 5;
    let startPage = Math.max(1, currentPage - Math.floor(maxDisplayedPages / 2));
    let endPage = Math.min(totalPages, startPage + maxDisplayedPages - 1);

    if (endPage - startPage + 1 < maxDisplayedPages) {
        startPage = Math.max(1, endPage - maxDisplayedPages + 1);
    }

    if (startPage > 1) {
        appendPaginationButton(paginationContainer, 1);
        if (startPage > 2) appendEllipsis(paginationContainer);
    }

    for (let i = startPage; i <= endPage; i++) {
        appendPaginationButton(paginationContainer, i);
    }

    if (endPage < totalPages - 1) {
        appendEllipsis(paginationContainer);
    }
    if (endPage < totalPages) {
        appendPaginationButton(paginationContainer, totalPages);
    }

    adjustPaginationForPhone();
}

function appendPaginationButton(container, page) {
    const button = document.createElement('button');
    button.innerText = page;
    button.classList.toggle('active', page === currentPage);
    button.onclick = () => changePage(page);
    container.appendChild(button);
}

function appendEllipsis(container) {
    const ellipsis = document.createElement('span');
    ellipsis.innerText = '...';
    ellipsis.classList.add('ellipsis');
    container.appendChild(ellipsis);
}

function changePage(page) {
    currentPage = page;
    loadWines(page);
}

function searchWines() {
    searchQuery = document.getElementById("searchBar").value.toLowerCase();
    currentPage = 1;
    loadWines(currentPage);
}

function askDeleteWine(wineId) {
    wineToDelete = wineId;
    document.getElementById("deleteModal").style.display = "flex";
    document.getElementById("deleteModal").style.opacity = '1';
    document.getElementById("deleteModal").classList.add("show");
}

async function confirmDelete() {
    if (wineToDelete !== null) {
        try {
            const response = await fetch(`/api/wines/${wineToDelete}`, { method: 'DELETE' });
            const responseJson = await response.json();

            if (responseJson.status === "success") {
                loadWines(currentPage);
                setupPagination();
            }
        } catch (error) {
            console.log("Impossible de suprimer le vin: ", error);
        }
        wineToDelete = null;
    }
    closeDeleteModal();
}

function closeDeleteModal() {
    document.getElementById("deleteModal").classList.remove('show');
    document.getElementById("deleteModal").style.opacity = '0';
    setTimeout(() => {
        document.getElementById("deleteModal").style.display = 'none';
    }, 300);
}

let oldInnerWidth = null;

function adjustPaginationForPhone() {
    if (oldInnerWidth === null || (window.innerWidth !== oldInnerWidth)) {
        oldInnerWidth = window.innerWidth;
    } else {
        return;
    }

    const paginationContainer = document.getElementById('pagination');
    if (window.innerWidth <= 768) {
        winesPerPage = 3;
        paginationContainer.style.position = "fixed";
        paginationContainer.style.bottom = "10px";
        paginationContainer.style.left = "50%";
        paginationContainer.style.transform = "translateX(-50%)";
        paginationContainer.style.display = "block";
        paginationContainer.style.justifyContent = "center";
        paginationContainer.style.gap = "5px";
    } else {
        winesPerPage = 6;
        paginationContainer.style.position = "static";
        paginationContainer.style.bottom = "";
        paginationContainer.style.left = "";
        paginationContainer.style.transform = "none";
        paginationContainer.style.display = "";
        paginationContainer.style.justifyContent = "";
        paginationContainer.style.gap = "";
    }

    loadWines(currentPage);
}

function closeEditWineModal() {
    document.getElementById("editWineModal").classList.remove('show');
    document.getElementById("editWineModal").style.opacity = '0';
    setTimeout(() => {
        document.getElementById("editWineModal").style.display = 'none';
    }, 300);
}

function getEditDomainData() {
    const domainName = document.getElementById('editDomainInput').value;
    const domain = jsonDomains ? jsonDomains.find(item => item.name === domainName) : null;
    return domain ? { "id": domain.id } : { "name": domainName };
}

function getEditRegionData() {
    const regionName = document.getElementById('editRegionInput').value;
    const countryName = document.getElementById('editCountryInput').value;

    if (regionName && countryName && jsonRegionCountries) {
        const region = jsonRegionCountries.find(item => item.name === regionName && item.country === countryName);
        return region ? { "id": region.id } : { "name": regionName, "country": countryName };
    }
    return { "name": regionName, "country": countryName };
}

function getEditTypeData() {
    const typeName = document.getElementById('editTypeInput').value;
    const type = jsonTypes ? jsonTypes.find(item => item.name === typeName) : null;
    return type ? { "id": type.id } : { "name": typeName };
}

function getEditBottleSizeData() {
    const bottleSizeName = document.getElementById('editBottleSizeInput').value;
    const bottleSizeValue = parseInt(document.getElementById('editBottleSizeValue').value);
    const bottleSize = jsonBottleSizes ? jsonBottleSizes.find(item => item.size === bottleSizeValue) : null;

    return bottleSize ? { "id": bottleSize.id } : { "name": bottleSizeName, "size": bottleSizeValue };
}

async function submitEditWineForm(event) {
    event.preventDefault();

    const startYear = document.getElementById("editPreferredStartDate").value;
    const endYear = document.getElementById("editPreferredEndDate").value;

    if ((startYear && !endYear) || (!startYear && endYear)) {
        alert("Les champs 'Date de Début' et 'Date de Fin' doivent être remplis si un des champs est fourni.");
        return false;
    }

    let startDate = null;
    let endDate = null;

    if (startYear && endYear) {
        const start = `${startYear}-01-01`;
        const end = `${endYear}-12-31`;

        if (new Date(start) > new Date(end)) {
            alert("'Date de Début' ne peut pas être après 'Date de Fin'.");
            return false;
        }

        startDate = start;
        endDate = end;
    }

    const jsonToSend = {
        "id": parseInt(document.getElementById("editWineId").value),
        "name": document.getElementById("editWineName").value,
        "domain": getEditDomainData(),
        "region": getEditRegionData(),
        "type": getEditTypeData(),
        "bottle_size": getEditBottleSizeData(),
        "vintage": parseInt(document.getElementById("editVintage").value),
        "quantity": parseInt(document.getElementById("editQuantity").value),
        "buy_price": parseFloat(document.getElementById("editBuyPrice").value) || null,
        "description": document.getElementById("editDescription").value || null,
        "image": document.getElementById("editImage").value || null,
        "preferred_start_date": startDate,
        "preferred_end_date": endDate
    };

    const response = await fetch("/api/wines/edit", {
        method: "POST",
        headers: {
            "Content-Type": "application/json"
        },
        body: JSON.stringify(jsonToSend)
    });

    if (response.status === 200) {
        closeEditWineModal();
        refresh();
    } else {
        console.log(response);
    }
}

function updateEditSliderValue(size) {
    document.getElementById("editSliderValueLabel").textContent = `${size}ml`;
    document.getElementById("editBottleSizeValue").value = size;
    document.getElementById("editBottleSizeSlider").value = size;
}

async function editWine(wineId) {
    try {
        const response = await fetch(`/api/wines/${wineId}`);
        if (!response.ok) {
            throw new Error("Failed to fetch wine details");
        }

        const wine = await response.json();

        const domainName = jsonDomains.find(domain => domain.id === wine.domain_id)?.name || "Unknown Domain";
        const typeName = jsonTypes.find(type => type.id === wine.type_id)?.name || "Unknown Type";
        const region = jsonRegionCountries.find(region => region.id === wine.region_id);
        const bottleSize = jsonBottleSizes.find(size => size.id === wine.bottle_size_id);

        document.getElementById("editWineId").value = wineId;
        document.getElementById("editWineName").value = wine.name;
        document.getElementById("editDomainInput").value = domainName;
        document.getElementById("editCountryInput").value = region.country;
        document.getElementById("editRegionInput").value = region.name;
        document.getElementById("editTypeInput").value = typeName;
        document.getElementById("editBottleSizeInput").value = bottleSize.name;
        updateEditSliderValue(bottleSize.size);
        document.getElementById("editVintage").value = wine.vintage;
        document.getElementById("editQuantity").value = wine.quantity;
        document.getElementById("editBuyPrice").value = wine.buy_price || "";
        document.getElementById("editDescription").value = wine.description || "";
        document.getElementById("editImage").value = wine.image || "";
        document.getElementById("editPreferredStartDate").value = wine.preferred_start_date ? wine.preferred_start_date.split("-")[0] : "";
        document.getElementById("editPreferredEndDate").value = wine.preferred_end_date ? wine.preferred_end_date.split("-")[0] : "";

        document.getElementById("editWineModal").style.display = "flex";
        document.getElementById("editWineModal").classList.add("show");
        document.getElementById("editWineModal").style.opacity = '1';
    } catch (error) {
        console.error(error);
    }
}