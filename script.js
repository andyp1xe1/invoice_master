// Get the popup elements
const popupOverlay = document.getElementById('popup-overlay');
const signupLink = document.querySelector('.signup-link');

// Open popup when the "Sign Up" link is clicked
signupLink.addEventListener('click', function (e) {
    e.preventDefault(); // Prevent the default anchor behavior
    popupOverlay.style.display = 'flex';
});

// Close popup when clicking outside the popup
popupOverlay.addEventListener('click', function (e) {
    if (e.target === popupOverlay) {
        popupOverlay.style.display = 'none';
    }
});
