// const menu = document.querySelector("#menu");
// const modal = document.querySelector("#modal");
//
// menu.addEventListener("click", () => {
//     if (menu.classList.contains("inactive")) {
//         modal.showModal();
//         menu.classList.add("active");
//         menu.classList.remove("inactive");
//     } else if (menu.classList.contains("active")) {
//         modal.close();
//         menu.classList.add("inactive");
//         menu.classList.remove("active");
//     }
// });
// profile.js
document.addEventListener('DOMContentLoaded', function() {
    document.body.addEventListener('click', function(event) {
        if (event.target.id === 'close-profile') {
            document.getElementById('profile-slide').classList.add('-translate-x-full');
        } else if (event.target.id === 'edit-profile-btn') {
            document.getElementById('profile-info').classList.add('hidden');
            document.getElementById('edit-profile-btn').classList.add('hidden');
            document.getElementById('edit-profile-form').classList.remove('hidden');
        } else if (event.target.id === 'cancel-edit') {
            document.getElementById('profile-info').classList.remove('hidden');
            document.getElementById('edit-profile-btn').classList.remove('hidden');
            document.getElementById('edit-profile-form').classList.add('hidden');
        }
    });

    document.body.addEventListener('change', function(event) {
        if (event.target.id === 'profile-picture-upload') {
            if (event.target.files && event.target.files[0]) {
                var file = event.target.files[0];
                var maxSize = 20 * 1024 * 1024; // 20MB

                if (file.size > maxSize) {
                    alert('File is too large. Maximum size is 20MB.');
                    event.target.value = '';
                    return;
                }
            }
        }
    });

    htmx.on('htmx:afterRequest', function (evt) {
        if (evt.detail.elt.id === 'profile-picture-upload') {
            if (evt.detail.successful) {
                console.log('Upload successful');
            } else {
                console.error('Upload failed');
                alert('An error occurred during upload. Please try again.');
            }
        }
    });

    htmx.on('htmx:responseError', function (evt) {
        console.error('HTMX response error:', evt.detail.error);
        alert('An error occurred. Please try again.');
    });
});

