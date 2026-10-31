// server/utils/db.ts
import postgres from 'postgres'

const connectionString = process.env.DATABASE_URL

export const sql = postgres(connectionString)