/// <reference lib="WebWorker" />
declare const self: ServiceWorkerGlobalScope;

const passages = [
  {
    reference: "Psalm 1",
    reviewDate: new Date("12-31-2022"),
    level: 3,
  },
  {
    reference: "Psalm 2",
    reviewDate: new Date("12-30-2022"),
    level: 2,
  },
];

self.addEventListener("fetch", (event: any) => {
  const url = new URL(event.request.url);
  if (url.pathname === "/api/passages") {
    switch (event.request.method) {
      case "GET": {
        event.respondWith(
          new Response(JSON.stringify(passages), {
            status: 200,
            headers: {
              "content-type": "application/js",
            },
          })
        );
        break;
      }
      case "POST": {
        event.respondWith(new Promise(async (resolve, reject) => {
          const body = await event.request.json()
          passages.push({
            reference: body.reference,
            reviewDate: new Date(),
            level: 0
          })
          resolve(new Response('', { status: 201 }))
        }))
        break;
      }
    }
  }
});

export {};
