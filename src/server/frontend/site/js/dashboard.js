document.getElementById('toggle-btn').addEventListener('click', function() {
    document.getElementById('sidebar').classList.toggle('active');
    this.classList.toggle('open');
});

function setActiveMenu(selectedLink) {
    const items = document.querySelectorAll('.sidebar-item');
    items.forEach(item => item.classList.remove('active')); 
    selectedLink.classList.add('active');
}

function showSection(sectionId) {
    const sections = document.querySelectorAll('.content-section');
    sections.forEach(section => {
        if (section.id === sectionId) {
            section.style.display = 'block'; 
        } else {
            section.style.display = 'none';
        }
    });
}

document.querySelectorAll('.sidebar-item').forEach(item => {
    item.addEventListener('click', async function() {
        // TODO ask the user if he is sure to quit while creating a wine
        const addWineModal = document.getElementById("addWineModal");
        const addHardWineModal = document.getElementById("addHardWineModal");
        if (addWineModal.style.display === "flex"
            || addHardWineModal.style.display === "flex"
        ) {
            if (addWineModal.style.display === "flex") {
                hideWineAddOptionsModal();
            } else {
                closeHardCreationModal();
            }
        }

        if (window.innerWidth <= 768) {
            document.getElementById('sidebar').classList.remove('active');
            document.getElementById('toggle-btn').classList.remove('open');
        }

        document.getElementById('loadingOverlay').style.display = 'block';
        document.getElementById('loadingIndicator').style.display = 'block';

        const sectionId = item.id.replace('Link', 'Content');
        await fetchSpecificData(sectionId);

        document.getElementById('loadingOverlay').style.display = 'none';
        document.getElementById('loadingIndicator').style.display = 'none';

        setActiveMenu(item);
        showSection(sectionId);
    });
});

async function fetchSpecificData(sectionId) {
    switch(sectionId) {
        case 'dashboardContent':
            await SetupDashboardPage();
            break;
        case 'collectionContent':
            await SetupCollectionPage();
            break;
        case 'statsContent':
            await SetupStatisticsPage();
            break;
        case 'accountContent':
            break;
    }
}

document.getElementById('dashboardLink').click();