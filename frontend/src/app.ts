import feed from './components/notification-feed';
import { isHTMXEvent } from './htmx-types'

export default () => {
    feed();
    document.body.addEventListener('htmx:beforeOnLoad', (e: Event) => {
        if (!isHTMXEvent<any>(e)) {
            return;
        }

        e.detail.shouldSwap = true;
        e.detail.isError = false;
    });
};