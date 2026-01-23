// Modal management
function showModal() {
    document.getElementById('modal-backdrop').classList.add('show');
}

function hideModal() {
    document.getElementById('modal-backdrop').classList.remove('show');
    document.getElementById('modal-content').innerHTML = '';
}

function hideConfirmModal() {
    document.getElementById('confirm-modal').classList.remove('show');
}

// Confirm delete dialog
function confirmDelete(itemName, deleteUrl, targetSelector) {
    const modal = document.getElementById('confirm-modal');
    const itemNameEl = document.getElementById('confirm-item-name');
    const deleteBtn = document.getElementById('confirm-delete-btn');

    itemNameEl.textContent = itemName;

    // Remove old event listeners by cloning
    const newDeleteBtn = deleteBtn.cloneNode(true);
    deleteBtn.parentNode.replaceChild(newDeleteBtn, deleteBtn);

    newDeleteBtn.addEventListener('click', function() {
        htmx.ajax('DELETE', deleteUrl, {
            target: targetSelector,
            swap: 'innerHTML'
        }).then(function() {
            hideConfirmModal();
            // Trigger event for dashboard stats refresh
            htmx.trigger(document.body, 'serverUpdated');
        });
    });

    modal.classList.add('show');
}

// Status dropdown toggle
function toggleStatusMenu(button) {
    const menu = button.nextElementSibling;

    // Close all other dropdowns
    document.querySelectorAll('.dropdown-menu.show').forEach(function(el) {
        if (el !== menu) {
            el.classList.remove('show');
        }
    });

    menu.classList.toggle('show');
}

// Close dropdowns when clicking outside
document.addEventListener('click', function(event) {
    if (!event.target.closest('.status-dropdown')) {
        document.querySelectorAll('.dropdown-menu.show').forEach(function(el) {
            el.classList.remove('show');
        });
    }
});

// Close modal on escape key
document.addEventListener('keydown', function(event) {
    if (event.key === 'Escape') {
        hideModal();
        hideConfirmModal();
    }
});

// HTMX event listeners
document.body.addEventListener('htmx:afterRequest', function(event) {
    // Close dropdowns after status change
    document.querySelectorAll('.dropdown-menu.show').forEach(function(el) {
        el.classList.remove('show');
    });
});

// Handle theme from system preference
if (!document.cookie.includes('theme=')) {
    const prefersDark = window.matchMedia('(prefers-color-scheme: dark)').matches;
    if (prefersDark) {
        document.documentElement.setAttribute('data-theme', 'dark');
    }
}

// Title case helper for templates
function title(str) {
    return str.charAt(0).toUpperCase() + str.slice(1);
}
