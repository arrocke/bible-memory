/// <reference lib="WebWorker" />
declare const self: ServiceWorkerGlobalScope;
import { add, format } from "date-fns";
import { open } from "./db";

const db = open();

const REVIEW_DAY_MAP = [1, 1, 2, 3, 5, 8, 13, 21, 34, 55]

function getDate(date: Date): Date {
  return new Date(format(date, 'yyy-MM-dd'))
}

self.addEventListener("fetch", (event) => {
  const url = new URL(event.request.url);
  if (url.pathname === "/api/passages") {
    event.respondWith(
      new Promise(async (resolve, reject) => {
        switch (event.request.method) {
          case "GET": {
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
            const passage = await db.passages.insert({
              ...body,
              level: 0
            });
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
            const body = (await event.request.json()) as { review?: boolean, level?: number, reference?: string, text?: string, reviewDate?: string };
            if (typeof body.review === 'boolean') {
              if (body.review) {
                passage.level = Math.min(REVIEW_DAY_MAP.length - 1, passage.level + 1);
              } else {
                passage.level = Math.ceil(passage.level / 2)
              }
              passage.reviewDate = getDate(add(new Date(), { days: REVIEW_DAY_MAP[passage.level] }));
            }
            if (typeof body.level === 'number') {
              passage.level = body.level
              if (passage.level > 0 && !passage.reviewDate) {
                passage.reviewDate = getDate(add(new Date(), { days: REVIEW_DAY_MAP[passage.level] }));
              }
            }
            if (typeof body.reference === 'string') {
              passage.reference = body.reference
            }
            if (typeof body.text === 'string') {
              passage.text = body.text
            }
            if (typeof body.reviewDate === 'string') {
              passage.reviewDate = getDate(new Date(body.reviewDate))
            }
            await db.passages.update(passage);
            return resolve(new Response(null, { status: 204 }));
          }
          case "DELETE": {
            await db.passages.delete(id)
            return resolve(new Response(null, { status: 204 }))
          }
        }
      })
    );
  }
});

export {};
