/// <reference lib="WebWorker" />
declare const self: ServiceWorkerGlobalScope;
import { add } from 'date-fns'

const passages = [
  {
    id: "1",
    reference: "Psalm 25:1",
    reviewDate: new Date("12-31-2022"),
    level: 3,
    text: `1 To you, O Lord, I lift up my soul.`,
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
            const passage = {
              id: (passages.length + 1).toString(),
              reference: body.reference,
              reviewDate: new Date(),
              level: 0,
              text: body.text,
            };
            passages.push(passage);
            resolve(
              new Response(null, {
                status: 201,
                headers: {
                  location: `/api/passages/${passage.id}`,
                },
              })
            );
          })
        );
        break;
      }
    }
  } else if (url.pathname.match(/^\/api\/passages\/(\w|\d)+$/)) {
    const id = url.pathname.split("/").at(-1);
    const passage = passages.find((passage) => passage.id === id);
    if (!passage) {
      event.respondWith(new Response(null, { status: 404 }));
      return;
    }
    switch (event.request.method) {
      case "GET": {
        event.respondWith(
          new Response(JSON.stringify(passage), {
            status: 200,
            headers: {
              "content-type": "application/js",
            },
          })
        );
        break;
      }
      case "PATCH": {
        new Promise(async (resolve, reject) => {
          const body = await event.request.json() as { review: boolean };
          if (body.review) {
            passage.level += 1;
          }
          passage.reviewDate = add(new Date(), { days: passage.level });
          event.respondWith(
            new Response(null, {
              status: 204,
            })
          );
        })
        break;
      }
    }
  }
});

export {};
