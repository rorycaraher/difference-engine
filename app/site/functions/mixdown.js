export async function onRequestPost(context) {
    const url = context.env.WORKERS_URL;
    if (!url) {
        return new Response("WORKERS_URL is not configured", { status: 500 });
    }
    return fetch(url, {
        method: "POST",
        headers: context.request.headers,
        body: context.request.body,
    });
}
