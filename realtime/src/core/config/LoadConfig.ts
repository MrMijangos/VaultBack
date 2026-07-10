import "dotenv/config";

import { Config } from "./Config";

export function loadConfig(): Config {
  const databaseUrl = process.env.DATABASE_URL;
  const jwtSecret = process.env.JWT_SECRET;

  if (!databaseUrl) {
    throw new Error("falta la variable de entorno obligatoria: DATABASE_URL");
  }
  if (!jwtSecret) {
    throw new Error("falta la variable de entorno obligatoria: JWT_SECRET");
  }

  return {
    port: process.env.PORT ? parseInt(process.env.PORT, 10) : 8081,
    databaseUrl,
    jwtSecret,
    corsOrigin: process.env.CORS_ORIGIN ?? "*",
  };
}
