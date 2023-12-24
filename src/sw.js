const cacheName = "BibleMemory";
const dbName = "bible-memory";

/**
 * @param {Event} event
 */
async function handleInstall(event) {
  const cache = await caches.open(cacheName);
  cache.addAll([
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

const passages = [
  { reference: "Gen 1:1-1:5", level: 1, reviewDate: new Date() },
  { reference: "Gen 1:6-1:10", level: 3, reviewDate: new Date() },
];

/**
 * @param {Event} event
 */
async function handleFetch(event) {
  const url = new URL(event.request.url);

  console.log(url.pathname);
  switch (event.request.method) {
    case "GET": {
      switch (url.pathname) {
        case "/": {
          console.log("request");
          const passages = await db.getAll();
          return renderHTML(generateHomePage({ passages }));
        }
        case "/passages/new": {
          return renderHTML(generateNewPassagePage());
        }
        default: {
          const cache = await caches.open(cacheName);

          const file = await cache.match(event.request, { ignoreSearch: true });
          if (file) {
            return file;
          } else {
            return fetch(event.request);
          }
        }
      }
    }
    case "POST": {
      switch (url.pathname) {
        case "/passages/new": {
          const data = await event.request.formData();
          await db.insert({
            reference: `${data.get("book")} ${data.get(
              "start_chapter"
            )}:${data.get("start_verse")}-${data.get("end_chapter")}:${data.get(
              "end_verse"
            )}`,
            level: 0,
          });
          return redirect("/");
        }
        default: {
          return new Response("", {
            status: 405,
          });
        }
      }
    }
    case "DELETE": {
      if (url.pathname.match(/\/passages\/\w+$/)) {
        const id = url.pathname.split("/").at(-1);
        await db.delete(parseInt(id));
        return new Response("", {
          status: 200,
          headers: {
            "Content-Type": "text/html",
          },
        });
      } else {
        return new Response("", {
          status: 405,
        });
      }
    }
    default: {
      return new Response("", {
        status: 405,
      });
    }
  }
}

/**
 * @param {string} html
 */
function renderHTML(html) {
  return new Response(html, {
    status: 200,
    headers: {
      "Content-Type": "text/html",
    },
  });
}

/**
 * @param {string} path
 */
function redirect(path) {
  return new Response("", {
    status: 302,
    headers: {
      Location: path,
    },
  });
}

/**
 * @param {{ head?: string; body?: string; title?: string }} options
 */
function generatePage(options) {
  return `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <meta name="theme-color" content="#317EFB" />

  <link rel='manifest' href='/manifest.json'>

  <title>${options.title ? `${options.title} | ` : ""}Bible Memory</title>

  <script src="https://unpkg.com/htmx.org@1.9.10"></script>
  <script>
    if ('serviceWorker' in navigator) {
      navigator.serviceWorker.register('/sw.js', { scope: '/' }).then((registration) => {
        console.log(registration)
        if (!registration.active) {
          location.reload()
        }
      })
    }
  </script>
  ${options.head ?? ""}
</head>
<body>
  ${options.body ?? ""}
</body>
</html>`;
}

function generateNewPassagePage() {
  return generatePage({
    body: `
      <h1>New Passage</h1>
      <form hx-boost="true" action="/passages/new" method="POST">
        <fieldset>
          <legend>REFERENCE</legend>
          <div>
            <label for="reference-book">BOOK</label>
            <input
              id="reference-book"
              name="book"
              autoComplete="off"
              required
            />
          </div>
          <div>
            <label for="reference-start-chapter">Chapter</label>
            <input
              id="reference-start-chapter"
              name="start_chapter"
              type="number"
              autoComplete="off"
              step=1
              min=1
              required
            />
          </div>
          <div>:</div>
          <div>
            <label for="reference-start-verse">Verse</label>
            <input
              id="reference-start-verse"
              name="start_verse"
              type="number"
              autoComplete="off"
              step=1
              min=1
              required
            />
          </div>
          <div>&mdash;</div>
          <div>
            <label for="reference-end-chapter">Chapter</label>
            <input
              id="reference-end-chapter"
              name="end_chapter"
              type="number"
              autoComplete="off"
              step=1
              min=1
              required
            />
          </div>
          <div>:</div>
          <div>
            <label for="reference-end-verse">Verse</label>
            <input
              id="reference-end-verse"
              name="end_verse"
              type="number"
              autoComplete="off"
              step=1
              min=1
              required
            />
          </div>
        </fieldset>
        <div>
          <label for="text">TEXT</label>
          <textarea
            id="text"
            name="text"
            required
          ></textarea>
        </div>
        <div className="mt-4">
          <button>Save</button>
        </div>
      </form>
    `,
  });
}

function generateHomePage(context) {
  return generatePage({
    body: `
      <h1>Passages</h1>
      <table>
        <thead>
          <tr>
            <td>PASSAGE</td>
            <td>LEVEL</td>
            <td>NEXT REVIEW</td>
            <td></td>
          </tr>
            <tr>
              <td colSpan={4}>
                <a href="/passages/new">
                  + Add Passage
                </a>
              </td>
            </tr>
        </thead>
        <tbody>
          ${context.passages?.map(
            (passage) => `<tr>
            <td>${passage.reference}</td>
            <td>${passage.level}</td>
            <td>${passage.reviewDate}</td>
            <td><button type="button" hx-delete="/passages/${passage.id}" hx-target="closest tr" hx-swap="outerHTML">Delete</button></td>
          </tr>`
          )}
        </tbody>
      </table>
    `,
  });
}

function generateNotFoundPage() {
  return generatePage({
    title: "Not Found",
    body: "Not Found",
  });
}

async function openDb() {
  const db = await new Promise((res, rej) => {
    const request = indexedDB.open(dbName, 2);
    request.onerror = (event) => {
      rej(new Error(`Error opening database ${dbName}`));
    };
    request.onsuccess = (event) => {
      res(event.target.result);
    };

    request.onupgradeneeded = (event) => {
      console.log("upgrade");
      event.currentTarget.result.createObjectStore("passages", {
        keyPath: "id",
        autoIncrement: true,
      });
    };
  });

  db.onerror = (event) => {
    console.error(`Database error: ${event.target.errorCode}`);
  };

  return {
    getAll() {
      return new Promise((res, rej) => {
        let passages = [];
        db
          .transaction("passages")
          .objectStore("passages")
          .openCursor().onsuccess = (event) => {
          const cursor = event.target.result;
          if (cursor) {
            passages.push(cursor.value);
            cursor.continue();
          } else {
            res(passages);
          }
        };
      });
    },
    async insert(data) {
      return new Promise((resolve, reject) => {
        db
          .transaction("passages", "readwrite")
          .objectStore("passages")
          .add(data).onsuccess = (event) => {
          resolve();
        };
      });
    },
    async delete(id) {
      return new Promise((resolve, reject) => {
        db
          .transaction("passages", "readwrite")
          .objectStore("passages")
          .delete(id).onsuccess = (event) => {
          resolve();
        };
      });
    },
  };
}

// Ensures db upgrades when needed
let db;
openDb().then((_db) => (db = _db));

self.addEventListener("install", (e) => e.waitUntil(handleInstall(e)));

self.addEventListener("activate", (event) => event.waitUntil(clients.claim()));

self.addEventListener("fetch", (event) =>
  event.respondWith(handleFetch(event))
);
