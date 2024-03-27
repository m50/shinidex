import htmx from 'htmx.org';

// Init HTMX
declare global {
    interface Window { htmx: typeof htmx }
}
window.htmx = htmx;