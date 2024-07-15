const menu = document.querySelector("#menu");
const modal = document.querySelector("#modal");

menu.addEventListener("click", () => {
    if (menu.classList.contains("inactive")) {
        modal.showModal();
        menu.classList.add("active");
        menu.classList.remove("inactive");
    } else if (menu.classList.contains("active")) {
        modal.close();
        menu.classList.add("inactive");
        menu.classList.remove("active");
    }
});
