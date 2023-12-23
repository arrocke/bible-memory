const cacheName = "BibleMemory";

/**
 * @param {Event} event
 */
async function handleInstall(event) {
  const cache = await caches.open(cacheName);
  cache.addAll([
    "/",
    "/index.html",
    "/manifest.json",
    "/icons/icon-16x16.png",
    "/icons/icon-32x32.png",
    "/icons/icon-72x72.png",
    "/icons/icon-96x96.png",
    "/icons/icon-128x128.png",
    "/icons/icon-144x144.png",
    "/icons/icon-192x192.png",
    "/icons/icon-384x384.png",
    "/icons/icon-512x512.png",
    "https://unpkg.com/htmx.org@1.9.10",
  ]);
}

/**
 * @param {Event} event
 */
async function handleFetch(event) {
  const url = new URL(event.request.url);

  console.log(url.pathname);
  switch (url.pathname) {
    case "/root": {
      return new Response("<div #app>test</div>", {
        status: 200,
        headers: {
          "Content-Type": "text/html",
        },
      });
    }
    default: {
      const cache = await caches.open(cacheName);

      const file = await cache.match(event.request, { ignoreSearch: true });
      if (file) {
        return file;
      } else {
        return cache.match("./index.html");
      }
    }
  }
}

self.addEventListener("install", (e) => e.waitUntil(handleInstall(e)));

self.addEventListener("activate", (event) => event.waitUntil(clients.claim()));

self.addEventListener("fetch", (event) =>
  event.respondWith(handleFetch(event))
);
