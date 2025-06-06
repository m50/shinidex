import { delay } from "../functions";

function observe(mutationList: MutationRecord[]) {
    const feed = document.getElementById("notification-feed")
    const k = feed?.childNodes.length ?? 0

    for (const mutation of mutationList) {
        if (mutation.type !== "childList") {
            continue
        }
        const notification = (mutation.addedNodes[0] as HTMLDivElement)
        if (!notification) {
            continue
        }
        console.log(k)
        if (k > 5) {
            const firstElm = (feed?.lastElementChild as HTMLDivElement)
            firstElm.classList.add("opacity-0")
            firstElm.remove()
        }
        delay(3000).then(async () => {
            notification.classList.add("opacity-0")
            await delay(1000)
            notification.remove()
        })
    }
}

export default () => {
    const feed = document.getElementById("notification-feed")
    if (feed === null) {
        console.log("#notification-feed not found")
        return
    }
    const observer = new MutationObserver(observe)
    observer.observe(feed, {childList: true})
}