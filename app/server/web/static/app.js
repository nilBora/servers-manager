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

// Initialize hints on modal load and dashboard sorting on content updates
document.body.addEventListener('htmx:afterSwap', function(event) {
    const providerSelect = document.getElementById('provider_id');
    if (providerSelect) {
        updateApiKeyHint(providerSelect);
    }

    // Reinitialize dashboard sorting if container exists
    if (document.getElementById('provider-groups-container')) {
        initDashboardSorting();
    }
});

// Initialize on page load
document.addEventListener('DOMContentLoaded', function() {
    if (document.getElementById('provider-groups-container')) {
        initDashboardSorting();
    }
});

// Dashboard drag-and-drop sorting
const DASHBOARD_ORDER_KEY = 'dashboard_group_order';

function initDashboardSorting() {
    const container = document.getElementById('provider-groups-container');
    if (!container) return;

    // Apply saved order first
    applySavedOrder(container);

    // Setup drag events
    const groups = container.querySelectorAll('.provider-group');
    groups.forEach(group => {
        group.addEventListener('dragstart', handleDragStart);
        group.addEventListener('dragend', handleDragEnd);
        group.addEventListener('dragover', handleDragOver);
        group.addEventListener('dragenter', handleDragEnter);
        group.addEventListener('dragleave', handleDragLeave);
        group.addEventListener('drop', handleDrop);
    });
}

let draggedElement = null;

function handleDragStart(e) {
    draggedElement = this;
    this.classList.add('dragging');
    e.dataTransfer.effectAllowed = 'move';
    e.dataTransfer.setData('text/plain', this.dataset.groupKey);
}

function handleDragEnd(e) {
    this.classList.remove('dragging');
    document.querySelectorAll('.provider-group').forEach(group => {
        group.classList.remove('drag-over');
    });
    draggedElement = null;
}

function handleDragOver(e) {
    e.preventDefault();
    e.dataTransfer.dropEffect = 'move';
}

function handleDragEnter(e) {
    e.preventDefault();
    if (this !== draggedElement) {
        this.classList.add('drag-over');
    }
}

function handleDragLeave(e) {
    this.classList.remove('drag-over');
}

function handleDrop(e) {
    e.preventDefault();
    this.classList.remove('drag-over');

    if (draggedElement && this !== draggedElement) {
        const container = this.parentNode;
        const allGroups = Array.from(container.querySelectorAll('.provider-group'));
        const draggedIdx = allGroups.indexOf(draggedElement);
        const targetIdx = allGroups.indexOf(this);

        if (draggedIdx < targetIdx) {
            container.insertBefore(draggedElement, this.nextSibling);
        } else {
            container.insertBefore(draggedElement, this);
        }

        // Save new order
        saveGroupOrder(container);
    }
}

function saveGroupOrder(container) {
    const groups = container.querySelectorAll('.provider-group');
    const order = Array.from(groups).map(g => g.dataset.groupKey);
    localStorage.setItem(DASHBOARD_ORDER_KEY, JSON.stringify(order));
}

function applySavedOrder(container) {
    const savedOrder = localStorage.getItem(DASHBOARD_ORDER_KEY);
    if (!savedOrder) return;

    try {
        const order = JSON.parse(savedOrder);
        const groups = Array.from(container.querySelectorAll('.provider-group'));
        const groupMap = new Map(groups.map(g => [g.dataset.groupKey, g]));

        // Sort groups according to saved order
        order.forEach(key => {
            const group = groupMap.get(key);
            if (group) {
                container.appendChild(group);
            }
        });

        // Append any new groups that weren't in saved order
        groups.forEach(group => {
            if (!order.includes(group.dataset.groupKey)) {
                container.appendChild(group);
            }
        });
    } catch (e) {
        console.warn('Failed to apply saved dashboard order:', e);
    }
}
