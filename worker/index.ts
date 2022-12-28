/// <reference lib="WebWorker" />
declare const self: ServiceWorkerGlobalScope;

const passages = [
  {
    id: "1",
    reference: "Psalm 1",
    reviewDate: new Date("12-31-2022"),
    level: 3,
  },
  {
    id: "2",
    reference: "Psalm 2",
    reviewDate: new Date("12-30-2022"),
    level: 2,
  },
];

self.addEventListener("fetch", (event) => {
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
        event.respondWith(
          new Promise(async (resolve, reject) => {
            const body = await event.request.json();
            passages.push({
              id: passages.length.toString(),
              reference: body.reference,
              reviewDate: new Date(),
              level: 0,
            });
            resolve(new Response("", { status: 201 }));
          })
        );
        break;
      }
    }
  } else if (url.pathname.match(/^\/api\/passages\/(\w|\d)+$/)) {
    switch (event.request.method) {
      case "GET": {
        const id = url.pathname.split("/").at(-1);
        const passage = passages.find((passage) => passage.id === id);
        if (passage) {
          event.respondWith(
            new Response(JSON.stringify(passage), {
              status: 200,
              headers: {
                "content-type": "application/js",
              },
            })
          );
        } else {
          event.respondWith(
            new Response("", { status: 404 })
          );
        }
        break;
      }
    }
  }
});

export {};
