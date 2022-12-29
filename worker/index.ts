/// <reference lib="WebWorker" />
declare const self: ServiceWorkerGlobalScope;
import { add } from "date-fns";
import { open } from "./db";

const db = open();

self.addEventListener("fetch", (event) => {
  const url = new URL(event.request.url);
  if (url.pathname === "/api/passages") {
    event.respondWith(
      new Promise(async (resolve, reject) => {
        switch (event.request.method) {
          case "GET": {
            console.log(await db.passages.getAll())
            return resolve(
              new Response(JSON.stringify(await db.passages.getAll()), {
                status: 200,
                headers: {
                  "content-type": "application/js",
                },
              })
            );
          }
          case "POST": {
            const body = await event.request.json();
            const passage = await db.passages.insert(body);
            return resolve(
              new Response(null, {
                status: 201,
                headers: {
                  location: `/api/passages/${passage.id}`,
                },
              })
            );
          }
        }
      })
    );
  } else if (url.pathname.match(/^\/api\/passages\/(\w|\d)+$/)) {
    event.respondWith(
      new Promise(async (resolve, reject) => {
        const id = url.pathname.split("/").at(-1)!;
        const passage = await db.passages.getById(id);
        if (!passage) {
          return resolve(new Response(null, { status: 404 }));
        }
        switch (event.request.method) {
          case "GET": {
            return resolve(
              new Response(JSON.stringify(passage), {
                status: 200,
                headers: {
                  "content-type": "application/js",
                },
              })
            );
          }
          case "PATCH": {
            const body = (await event.request.json()) as { review: boolean };
            if (body.review) {
              passage.level += 1;
            }
            passage.reviewDate = add(new Date(), { days: passage.level });
            await db.passages.update(passage);
            return resolve(new Response(null, { status: 204 }));
          }
        }
      })
    );
  }
});

export {};
