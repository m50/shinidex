import registerHamburgerMenu from './components/hamburger-menu';
import registerNotificationFeed from './components/notification-feed';
import registerJumpTo from './components/jump-to';
import { isHTMXEvent } from './htmx-types'

export default () => {
    registerNotificationFeed()
    registerHamburgerMenu()
    registerJumpTo()
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
        registerJumpTo()
    })
};