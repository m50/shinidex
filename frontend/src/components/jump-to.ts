import { delay } from "../functions"

const pathRegexp = /\/dex\/(?!new|\w+\/edit)\w+|\/pokemon/

export default () => {
    console.log(pathRegexp.test(document.location.href))
    const jumpTo = (document.getElementById('jump-to') as HTMLInputElement)
    const app = document.getElementById('app')
    if (!jumpTo) {
        console.log("unable to find jumpTo")
        return
    }

    const urlParams = new URLSearchParams(window.location.search);
    const jumpToParamVal = urlParams.get('jump-to');
    if (!!jumpToParamVal) {
        jumpTo.value = jumpToParamVal
        findPokemon(jumpToParamVal)
    }

    jumpTo.addEventListener("keyup", changeEvent(jumpTo))
    jumpTo.addEventListener("change", changeEvent(jumpTo))
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
    console.log(isKeyboardEvent(e), isKeyboardEvent(e) && e.key)
    if (!isKeyboardEvent(e) || (isKeyboardEvent(e) && e.key == "Enter")) {
        e.preventDefault()
    }
    const val = jumpTo.value
    setQueryParam(val)
    findPokemon(val) 
}

function setQueryParam(value: string) {
    const url = new URL(window.location.href)
    if (value != "") {
        url.searchParams.set("jump-to", value);
    } else {
        url.searchParams.delete("jump-to")
    }
    history.pushState(null, '', url);
}

function findPokemon(value: string) {
    value = value.toLocaleLowerCase()
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

function isKeyboardEvent(e: Event | KeyboardEvent): e is KeyboardEvent {
    return (typeof((e as KeyboardEvent).key) != "undefined") 
}
