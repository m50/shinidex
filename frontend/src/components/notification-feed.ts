import { delay } from "../functions";

function observe(mutationList: MutationRecord[]) {
    for (const mutation of mutationList) {
        if (mutation.type !== "childList") {
            continue;
        }
        const notification = (mutation.addedNodes[0] as HTMLDivElement)
        if (!notification) {
            continue;
        }
        delay(5000).then(async () => {
            notification.classList.add("opacity-0")
            await delay(1000);
            notification.remove();
        })
    }
}

export default () => {
    const feed = document.getElementById("notification-feed");
    if (feed === null) {
        console.log("#notification-feed not found");
        return;
    }
    const observer = new MutationObserver(observe);
    observer.observe(feed, {childList: true});
}