export interface HTMXEvent<T> extends Event {
    detail: T;
}

export interface OobAfterSwapDetail {
    target: HTMLElement;
    elt: HTMLElement;
    xhr: XMLHttpRequest;
};

export function isHTMXEvent<T>(e: Event): e is HTMXEvent<T> {
    return e.type.startsWith("htmx:")
}
