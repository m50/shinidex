import { delay } from "../functions"

const pathRegexp = /\/dex\/(?!new)\w+|\/pokemon/

export default () => {
    const jumpTo = (document.getElementById('jump-to') as HTMLInputElement)
    const app = document.getElementById('app')
    if (!jumpTo) {
        console.log("unable to find jumpTo")
        return
    }
    jumpTo.addEventListener("keyup", changeEvent(jumpTo))
    if (!pathRegexp.test(document.location.href)) {
        jumpTo.parentElement?.parentElement?.remove()
        app?.classList.remove("mt-32")
        app?.classList.add("mt-16")
    } else {
        app?.classList.remove("mt-16")
        app?.classList.add("mt-32")
    }
}

const changeEvent = (jumpTo: HTMLInputElement) => (e: Event) => {
    const value = jumpTo.value
    const pkmn = document.getElementById(value)
    if (!pkmn) return 
    pkmn.scrollIntoView({
        behavior: "smooth",
        block: "center",
    })
    pkmn.classList.add("bg-indigo-300")
    delay(3000).then(() => {
        pkmn.classList.remove("bg-indigo-300")
    })
}