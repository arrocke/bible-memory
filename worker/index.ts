/// <reference lib="WebWorker" />
declare const self: ServiceWorkerGlobalScope;

const passages = [
  {
    id: "1",
    reference: "Psalm 25",
    reviewDate: new Date("12-31-2022"),
    level: 3,
    text: `1 To you, O Lord, I lift up my soul.
2 O my God, in you I trust;
    let me not be put to shame;
    let not my enemies exult over me.
3 Indeed, none who wait for you shall be put to shame;
    they shall be ashamed who are wantonly treacherous.

4 Make me to know your ways, O Lord;
    teach me your paths.
5 Lead me in your truth and teach me,
    for you are the God of my salvation;
    for you I wait all the day long.

6 Remember your mercy, O Lord, and your steadfast love,
    for they have been from of old.
7 Remember not the sins of my youth or my transgressions;
    according to your steadfast love remember me,
    for the sake of your goodness, O Lord!

8 Good and upright is the Lord;
    therefore he instructs sinners in the way.
9 He leads the humble in what is right,
    and teaches the humble his way.
10 All the paths of the Lord are steadfast love and faithfulness,
    for those who keep his covenant and his testimonies.

11 For your name's sake, O Lord,
    pardon my guilt, for it is great.
12 Who is the man who fears the Lord?
    Him will he instruct in the way that he should choose.
13 His soul shall abide in well-being,
    and his offspring shall inherit the land.
14 The friendship of the Lord is for those who fear him,
    and he makes known to them his covenant.
15 My eyes are ever toward the Lord,
    for he will pluck my feet out of the net.

16 Turn to me and be gracious to me,
    for I am lonely and afflicted.
17 The troubles of my heart are enlarged;
    bring me out of my distresses.
18 Consider my affliction and my trouble,
    and forgive all my sins.

19 Consider how many are my foes,
    and with what violent hatred they hate me.
20 Oh, guard my soul, and deliver me!
    Let me not be put to shame, for I take refuge in you.
21 May integrity and uprightness preserve me,
    for I wait for you.

22 Redeem Israel, O God,
    out of all his troubles.`
  }
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
              text: ''
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
