const passages = [
  {
    reference: "Psalm 1",
    reviewDate: new Date("12-31-2022"),
  },
  {
    reference: "Psalm 2",
    reviewDate: new Date("12-30-2022"),
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
            reviewDate: new Date()
          })
          resolve(new Response('', { status: 201 }))
        }))
        break;
      }
    }
  }
});

export {};
