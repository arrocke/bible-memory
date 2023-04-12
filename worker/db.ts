import { DBSchema, openDB } from "idb";

export interface DbPassage {
  id: number;
  reference: string;
  text: string;
  level: number;
  reviewDate?: Date;
}

interface DbSchema extends DBSchema {
  passages: {
    key: number;
    value: DbPassage;
  };
}

export function open() {
  let migratePassages = false;
  const dbPromise = openDB<DbSchema>("bible-memory", 2, {
    upgrade(db, prevVersion, newVersion, transaction) {
      if (prevVersion < 2) {
        migratePassages = true;
        transaction.objectStore("passages").name = "passages-old";
      }
      db.createObjectStore("passages", {
        keyPath: "id",
        autoIncrement: true,
      });
    },
  });

  dbPromise.then(async (db) => {
    if (migratePassages) {
      const passages = await db.getAll("passages-old" as any);

      for (const { id, ...passage } of passages) {
        await db.put("passages", passage);
      }
    }
  });

  return {
    passages: {
      async insert(data: Omit<DbPassage, "id">): Promise<DbPassage> {
        const id = await (await dbPromise).put("passages", data as DbPassage);
        return { ...data, id };
      },
      async update(passage: DbPassage): Promise<void> {
        await (await dbPromise).put("passages", passage);
      },
      async getById(id: number): Promise<DbPassage | undefined> {
        return await (await dbPromise).get("passages", id);
      },
      async getAll(): Promise<DbPassage[]> {
        return await (await dbPromise).getAll("passages");
      },
      async delete(id: number): Promise<void> {
        return await (await dbPromise).delete("passages", id);
      },
    },
  };
}
