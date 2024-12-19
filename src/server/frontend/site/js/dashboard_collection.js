const allWines = [
    { id: 1, name: "Bordeaux 2015", region: "Bordeaux", vintage: 2015, type: "Red", size: "750ml", quantity: 10, image: "/v1/images/no_photo_generic.svg" },
    { id: 2, name: "Champagne Brut 2018", region: "Champagne", vintage: 2018, type: "Sparkling", size: "750ml", quantity: 8, image: "https://via.placeholder.com/200" },
    { id: 3, name: "Pinot Noir 2017", region: "Burgundy", vintage: 2017, type: "Red", size: "375ml", quantity: 5, image: "https://via.placeholder.com/200" },
    { id: 4, name: "Merlot 2016", region: "California", vintage: 2016, type: "Red", size: "750ml", quantity: 3, image: "https://via.placeholder.com/200" },
    { id: 5, name: "Cabernet Sauvignon 2019", region: "California", vintage: 2019, type: "Red", size: "1500ml", quantity: 7, image: "https://via.placeholder.com/200" },
];

let filteredWines = [...allWines];
const winesPerPage = 3;
let currentPage = 1;
let wineToDelete = null;

function SetupCollectionPage() {
    document.getElementById("collectionContent").style.display = "block";
    loadWines(currentPage);
    setupPagination();
}

function loadWines(page) {
    const startIndex = (page - 1) * winesPerPage;
    const endIndex = startIndex + winesPerPage;
    const winesToShow = filteredWines.slice(startIndex, endIndex);

    const wineListContainer = document.querySelector('.wine-list');
    wineListContainer.innerHTML = '';

    winesToShow.forEach(wine => {
        const wineItem = document.createElement('div');
        wineItem.classList.add('wine-item');
        wineItem.innerHTML = `
            <img src="${wine.image}" alt="${wine.name}">
            <h4>${wine.name} (${wine.vintage})</h4>
            <p><strong>Type:</strong> ${wine.type}</p>
            <p><strong>Region:</strong> ${wine.region}</p>
            <p><strong>Taille:</strong> ${wine.size}</p>
            <p><strong>Quantité:</strong> ${wine.quantity}</p>
            <div class="wine-item-actions">
                <button class="edit" onclick="editWine(${wine.id})">Modifier</button>
                <button class="delete" onclick="askDeleteWine(${wine.id})">Supprimer</button>
            </div>
        `;
        wineListContainer.appendChild(wineItem);
    });
}

function setupPagination() {
    const totalPages = Math.ceil(filteredWines.length / winesPerPage);
    const paginationContainer = document.getElementById('pagination');
    paginationContainer.innerHTML = '';

    for (let i = 1; i <= totalPages; i++) {
        const button = document.createElement('button');
        button.innerText = i;
        button.onclick = () => changePage(i);
        paginationContainer.appendChild(button);
    }
}

function changePage(page) {
    currentPage = page;
    loadWines(page);
    setupPagination();
}

function editWine(wineId) {
    alert(`Editing wine with ID: ${wineId}`);
}

function askDeleteWine(wineId) {
    wineToDelete = wineId;
    document.getElementById("deleteModal").style.display = "flex";
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
    document.getElementById("deleteModal").style.display = "none";
}

function searchWines() {
    const searchTerm = document.getElementById("searchBar").value.toLowerCase();
    filteredWines = allWines.filter(wine =>
        wine.name.toLowerCase().includes(searchTerm) ||
        wine.region.toLowerCase().includes(searchTerm) ||
        wine.type.toLowerCase().includes(searchTerm) ||
        wine.size.toLowerCase().includes(searchTerm)
    );
    currentPage = 1;  // Reset to the first page after search
    loadWines(currentPage);
    setupPagination();
}