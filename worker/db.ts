import { openDB } from 'idb';

export interface DbPassage {
  id: string
  reference: string
  text: string
  level: number
  reviewDate?: Date
}

interface DbSchema {
  passages: {
    key: string
    value: DbPassage
  }
}

export function open() {
  const dbPromise = openDB<DbSchema>('bible-memory', 1, {
    upgrade(db) {
      db.createObjectStore('passages');
    },
  })

  return {
    passages: {
      async insert({ reference, text }: Pick<DbPassage, 'reference' | 'text'>): Promise<DbPassage> {
        const id = (1 + await (await dbPromise).count('passages')).toString()
        const passage = {
          id,
          reference,
          text,
          level: 0
        }
        await (await dbPromise).put('passages', passage, id)
        return passage
      },
      async update(passage: DbPassage): Promise<void> {
        await (await dbPromise).put('passages', passage, passage.id)
      },
      async getById(id: string): Promise<DbPassage> {
        return await (await dbPromise).get('passages', id)
      },
      async getAll(): Promise<DbPassage[]> {
        return await (await dbPromise).getAll('passages')
      }
    }
  }
}