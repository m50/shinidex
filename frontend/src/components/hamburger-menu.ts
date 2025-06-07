export default () => {
    const button = document.getElementById("navbar-button")
    if (!button) {
        console.error("unable to find hamburger menu button", button)
        return
    }
    button.addEventListener('click', buttonOnClick)
}

const buttonOnClick = (e: Event) => {
    const menu = (document.getElementById("navbar-menu") as HTMLDivElement)
    if (!menu) {
        console.error("unable to find hamburger menu", menu)
        return
    }
    const hidden = menu.classList.contains("hidden")
    console.log("click...", menu)
    if (hidden) {
        menu.classList.remove("hidden")
        menu.classList.add("flex")
    } else {
        menu.classList.remove("flex")
        menu.classList.add("hidden")
    }
}