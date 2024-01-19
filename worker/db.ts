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
  const dbPromise = openDB<DbSchema>("bible-memory", 2, {
    upgrade(db) {
      db.createObjectStore("passages", {
        keyPath: "id",
        autoIncrement: true,
      });
    },
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
