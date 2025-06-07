import registerHamburgerMenu from './components/hamburger-menu';
import feed from './components/notification-feed';
import { isHTMXEvent } from './htmx-types'

export default () => {
    feed()
    registerHamburgerMenu()
    document.body.addEventListener('htmx:beforeOnLoad', (e: Event) => {
        if (!isHTMXEvent<any>(e)) {
            return
        }

        e.detail.shouldSwap = true
        e.detail.isError = false
    });

    document.body.addEventListener('htmx:afterSwap', (e: Event) => {
        if (!isHTMXEvent<any>(e)) {
            return
        }
        registerHamburgerMenu()
    })
};