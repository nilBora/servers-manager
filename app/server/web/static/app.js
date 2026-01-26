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

// Update API key hint based on selected provider
function updateApiKeyHint(select) {
    const selectedOption = select.options[select.selectedIndex];
    const ident = selectedOption.getAttribute('data-ident') || '';

    const robotHint = document.getElementById('api-key-hint-robot');
    const cloudHint = document.getElementById('api-key-hint-cloud');

    if (!robotHint || !cloudHint) return;

    // Hide all hints first
    robotHint.style.display = 'none';
    cloudHint.style.display = 'none';

    // Show relevant hint
    if (ident === 'hetzner_robot') {
        robotHint.style.display = 'block';
    } else if (ident === 'hetzner_cloud') {
        cloudHint.style.display = 'block';
    }
}

// Initialize hints on modal load
document.body.addEventListener('htmx:afterSwap', function(event) {
    const providerSelect = document.getElementById('provider_id');
    if (providerSelect) {
        updateApiKeyHint(providerSelect);
    }
});
