document.addEventListener('DOMContentLoaded', () => {
    // Select elements
    const sidebar = document.getElementById('sidebar');
    const toggleSidebar = document.getElementById('toggle-sidebar');
    const mainContent = document.getElementById('main-content');

    // Handle sidebar collapse/expand on desktop
    sidebar.addEventListener('click', (event) => {
        if (window.innerWidth > 768 && event.target.id !== 'toggle-sidebar') {
            sidebar.classList.toggle('collapsed');
            if (sidebar.classList.contains('collapsed')) {
                mainContent.style.marginLeft = '80px';  // Narrow content when sidebar is collapsed
            } else {
                mainContent.style.marginLeft = '250px';  // Restore full content width
            }
        }
    });

    // Toggle sidebar for mobile
    toggleSidebar.addEventListener('click', () => {
        sidebar.classList.toggle('open');
        mainContent.style.marginLeft = sidebar.classList.contains('open') ? '250px' : '0';
    });
});
