import "dotenv/config";

import { Config } from "./Config";

// pg-connection-string lee sslmode de la URL y arma su propio objeto ssl,
// que luego pisa el `ssl: {...}` que le pasamos al Client (se mezclan con
// Object.assign y el de la URL queda último). Quitamos sslmode acá para
// que nuestra config de ssl sea la única fuente de verdad.
function stripSslModeParam(databaseUrl: string): string {
  const url = new URL(databaseUrl);
  url.searchParams.delete("sslmode");
  return url.toString();
}

export function loadConfig(): Config {
  const rawDatabaseUrl = process.env.DATABASE_URL;
  const jwtSecret = process.env.JWT_SECRET;

  if (!rawDatabaseUrl) {
    throw new Error("falta la variable de entorno obligatoria: DATABASE_URL");
  }
  if (!jwtSecret) {
    throw new Error("falta la variable de entorno obligatoria: JWT_SECRET");
  }

  return {
    port: process.env.PORT ? parseInt(process.env.PORT, 10) : 8081,
    databaseUrl: stripSslModeParam(rawDatabaseUrl),
    jwtSecret,
    corsOrigin: process.env.CORS_ORIGIN ?? "*",
  };
}
