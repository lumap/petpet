export function isImage(url: string) {
    return /\.(jpg|jpeg|png)$/.test(url.toLowerCase());
}
