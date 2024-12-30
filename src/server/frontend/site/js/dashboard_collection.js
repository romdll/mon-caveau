let winesPerPage = 8;
let currentPage = 1;
let totalWines = 0;
let searchQuery = "";

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

async function loadWines(page) {
    const wineListContainer = document.querySelector('.wine-list');
    wineListContainer.innerHTML = '';

    const placeholders = Array.from({ length: winesPerPage }, () => {
        const placeholder = document.createElement('div');
        placeholder.classList.add('wine-item', 'loading-placeholder');
        placeholder.innerHTML = `
            <div class="wine-info">
                <h4 class="loading-placeholder-text">Loading...</h4>
                <p class="loading-placeholder-text"><strong>Type:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Région:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Taille:</strong> Loading...</p>
                <p class="loading-placeholder-text"><strong>Quantité:</strong> Loading...</p>
            </div>`;
        wineListContainer.appendChild(placeholder);
        return placeholder;
    });

    try {
        const response = await fetch(`/api/wines?page=${page}&limit=${winesPerPage}&search=${encodeURIComponent(searchQuery)}`);
        const data = await response.json();
        const { wines, total } = data;
        totalWines = total;

        wines.forEach((wine, index) => {
            const domainName = jsonDomains.find(domain => domain.id === wine.domain_id)?.name || "Unknown Domain";
            const typeName = jsonTypes.find(type => type.id === wine.type_id)?.name || "Unknown Type";
            const region = jsonRegionCountries.find(region => region.id === wine.region_id);
            const regionName = region ? `${region.name} (${region.country})` : "Unknown Region";
            const bottleSize = jsonBottleSizes.find(size => size.id === wine.bottle_size_id);
            const bottleSizeName = bottleSize ? `${bottleSize.name} (${bottleSize.size} ml)` : "Unknown Size";

            const wineItem = document.createElement('div');
            wineItem.classList.add('wine-item');
            wineItem.innerHTML = `
                <div class="wine-info">
                    <h4>${wine.name} (${wine.vintage})</h4>
                    <p><strong>Domain:</strong> ${domainName}</p>
                    <p><strong>Type:</strong> ${typeName}</p>
                    <p><strong>Region:</strong> ${regionName}</p>
                    <p><strong>Size:</strong> ${bottleSizeName}</p>
                    <p><strong>Quantity:</strong> ${wine.quantity}</p>
                </div>
                <div class="wine-item-actions">
                    <button class="edit" onclick="editWine(${wine.id})" disabled>Modifier</button>
                    <button class="delete" onclick="askDeleteWine(${wine.id})">Supprimer</button>
                </div>`;

            const img = document.createElement('img');
            img.src = wine.image || "/v1/images/no_photo_generic.svg";
            img.alt = `Image of ${wine.name}`;
            img.onerror = () => {
                img.src = "/v1/images/no_photo_generic.svg";
            };
            wineItem.insertBefore(img, wineItem.firstChild);

            placeholders[index].replaceWith(wineItem);
        });

        if (wines.length < placeholders.length) {
            placeholders.slice(wines.length).forEach(placeholder => placeholder.remove());
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

function confirmDelete() {
    if (wineToDelete !== null) {
        const index = filteredWines.findIndex(w => w.id === wineToDelete);
        if (index !== -1) {
            filteredWines.splice(index, 1);
            loadWines(currentPage);
            setupPagination();
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
        winesPerPage = 8;
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