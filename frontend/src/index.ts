import htmx from 'htmx.org';
import app from './app'

// Init HTMX
declare global {
    interface Window { htmx: typeof htmx }
}
window.htmx = htmx;

// Init app
document.addEventListener("DOMContentLoaded", app);