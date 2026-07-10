# Vault API

Backend de Vault: microservicios independientes detrás de un gateway de nginx, con la base de datos PostgreSQL como único componente compartido.

```
                              Cliente (app móvil / web)
                                        │
                                        ▼
                          ┌─────────────────────────┐
                          │   gateway (nginx:80)     │
                          │   punto de entrada único  │
                          └────────────┬─────────────┘
                       /api/*          │        /ws
                    ┌──────────────────┴──────────────────┐
                    ▼                                       ▼
          ┌─────────────────────┐                 ┌───────────────────────┐
          │   api (Go, :8080)    │                 │ realtime (Node, :8081) │
          │   REST               │                 │  WebSocket             │
          └──────────┬───────────┘                 └───────────┬───────────┘
                     │                                          │
                     │  lee/escribe                LISTEN/NOTIFY │
                     │                              (tabla notifications)
                     └──────────────────┬───────────────────────┘
                                        ▼
                              ┌───────────────────┐
                              │     PostgreSQL      │
                              │  (esquema compartido) │
                              └───────────────────┘
```

- **`gateway/`** — nginx. Único punto de entrada externo (puerto `8000` vía Docker); enruta `/api/*` hacia `api` y `/ws` hacia `realtime`. Los servicios internos no se exponen directamente fuera de la red de contenedores.
- **`api/`** (Go) — API REST: autenticación, usuarios, productos (assets), negocios, posts, comentarios, reseñas, notificaciones, certificados blockchain, etc.
- **`realtime/`** (Node.js + TypeScript) — servicio de WebSocket para notificaciones en tiempo real. No expone REST; solo escucha la tabla `notifications` vía `LISTEN/NOTIFY` de Postgres y reenvía al usuario conectado.

Cada uno es un proyecto independiente (su propio `go.mod` / `package.json`, su propio `Dockerfile`, su propio ciclo de build y deploy). No es un monolito: lo único que comparten es la base de datos.

---

## Estructura del repositorio

```
Vault/
  api/                        # microservicio REST (Go)
    src/
      core/                   # utilidades transversales (config, seguridad, http, middleware)
      features/<feature>/
        domain/
          entities/           # modelos de dominio
          dto/request|response/ # contratos de entrada/salida
          repositories/        # puertos (interfaces)
        application/          # casos de uso
        infrastructure/
          adapters/           # implementaciones concretas (Postgres, Cloudinary...)
          controllers/         # entrada HTTP
          router/              # registro de rutas
          dependencies.go      # wiring
    main.go                   # solo arranca: carga config, arma dependencias, escucha
    init.sql                  # esquema completo de la base de datos
    Dockerfile
  realtime/                   # microservicio WebSocket (Node.js + TypeScript)
    src/
      core/                   # config, seguridad (mismo rol que en api/)
      features/notifications/
        domain/                # entidad + puertos (NotificationRepository, ConnectionRegistry)
        application/           # caso de uso (BroadcastNotificationUseCase)
        infrastructure/
          adapters/            # Postgres LISTEN/NOTIFY, registro de conexiones en memoria
          websocket/            # servidor WebSocket
          Dependencies.ts       # wiring
    src/index.ts               # solo arranca
    Dockerfile
  gateway/
    nginx.conf                 # reverse proxy: /api/* -> api, /ws -> realtime
    Dockerfile
```

Ambos servicios (`api/` y `realtime/`) siguen la misma arquitectura por capas (hexagonal / puertos y adaptadores) con SRP estricto: un archivo = una responsabilidad completa (un tipo con su constructor y sus métodos, o una función).

---

## Cómo levantar todo

Cada servicio es independiente: su propio `Dockerfile`, sin orquestador. Se puede correr cada uno directo con el runtime (`go run`, `npm run dev`) o construir su imagen con `docker build` — no hay `docker-compose`.

### Opción 1: cada servicio directo (sin Docker)

Requiere una instancia de PostgreSQL propia corriendo.

**`api/`** (Go):
```bash
cd api
cp .env.example .env   # completar credenciales reales
go run .
```
Variables (`api/.env`): `APP_PORT` (default `8080`), `DB_HOST`/`DB_PORT`/`DB_NAME`/`DB_USER`/`DB_PASSWORD`/`DB_SSL`, `JWT_SECRET`, `CORS_ORIGIN` (default `*`), `COOKIE_SECURE` (default `false`), `CLOUDINARY_CLOUD_NAME`/`CLOUDINARY_API_KEY`/`CLOUDINARY_API_SECRET`. El esquema ([`api/init.sql`](api/init.sql)) se aplica automáticamente al arrancar (`CREATE TABLE IF NOT EXISTS`, idempotente).

**`realtime/`** (Node.js + TypeScript):
```bash
cd realtime
npm install
cp .env.example .env   # JWT_SECRET debe ser IDENTICO al de api/.env
npm run dev             # desarrollo, con recarga automática
npm run build && npm start   # produccion
```
Variables (`realtime/.env`): `PORT` (default `8081`), `DATABASE_URL` (misma base que `api`), `JWT_SECRET` (**debe ser igual** al de `api`), `CORS_ORIGIN` (default `*`).

Cada servicio se prueba en su propio puerto (`8080` y `8081`).

### Opción 2: cada servicio como contenedor (Docker, sin orquestador)

```bash
# 1. Red compartida, para que los contenedores se resuelvan por nombre
docker network create vault_net

# 2. PostgreSQL
docker run -d --name postgres --network vault_net --network-alias postgres \
  -e POSTGRES_USER=vault_user -e POSTGRES_PASSWORD=changeme -e POSTGRES_DB=vault \
  postgres:16-alpine

# 3. api
docker build -t vault-api ./api
docker run -d --name api --network vault_net --network-alias api \
  --env-file api/.env -e DB_HOST=postgres \
  vault-api

# 4. realtime
docker build -t vault-realtime ./realtime
docker run -d --name realtime --network vault_net --network-alias realtime \
  --env-file realtime/.env -e DATABASE_URL="postgres://vault_user:changeme@postgres:5432/vault?sslmode=disable" \
  vault-realtime

# 5. gateway (nginx), unico puerto publicado al host
docker build -t vault-gateway ./gateway
docker run -d --name gateway --network vault_net -p 8000:80 \
  vault-gateway
```

Con esto, `api` y `realtime` no se publican al host — solo el gateway (`http://localhost:8000`), igual que en producción. El gateway usa `nginx:alpine` con la config fija de [`gateway/nginx.conf`](gateway/nginx.conf): escucha en el puerto `80` y apunta a los hosts `api:8080` y `realtime:8081`.

### Despliegue en Railway

Cada carpeta (`api/`, `realtime/`, `gateway/`) es un servicio separado de Railway (mismo proyecto, cada uno apuntando a su propio `Dockerfile` con "Root Directory"), más un plugin de PostgreSQL. Railway le da a cada servicio del proyecto un host privado tipo `<nombre-del-servicio>.railway.internal`, alcanzable solo por los demás servicios del mismo proyecto.

Como `gateway/nginx.conf` es un archivo fijo (no una plantilla), antes de desplegar hay que editarlo a mano con los nombres reales que le pongas a tus servicios en Railway:

```nginx
upstream vault_api {
    server <nombre-del-servicio-api>.railway.internal:8080;
}

upstream vault_realtime {
    server <nombre-del-servicio-realtime>.railway.internal:8081;
}
```

Y en cada servicio:

- **`api`**: variables de entorno = las de `api/.env.example`, usando el host/credenciales del PostgreSQL de Railway, más `APP_PORT=8080` fijo (Railway también le inyecta un `PORT` dinámico automáticamente; `APP_PORT` fijo hace que el backend lo ignore).
- **`realtime`**: variables = las de `realtime/.env.example`, más `PORT=8081` fijo (sobreescribe el que Railway inyecta automáticamente). `JWT_SECRET` debe ser idéntico al de `api`.
- **`gateway`**: sin variables de entorno especiales. En la configuración de red de Railway para este servicio, hay que indicar que el puerto interno de la app es `80` (el que escucha nginx) y generar el dominio público ahí. Es el único de los tres servicios con dominio público — `api` y `realtime` quedan solo en la red privada del proyecto.

---

## Autenticación

`POST /api/v1/auth/login` responde con el usuario y coloca el JWT en una cookie **HttpOnly** llamada `vault_token` (no accesible desde JavaScript, `SameSite=Lax`). Todas las rutas protegidas leen esa cookie.

Para clientes que no manejan cookies automáticamente (apps nativas, el WebSocket), el token también puede enviarse:
- En la conexión WebSocket: como query param `?token=<jwt>`.

El JWT contiene `user_id` (uuid) y `role`, expira a las 24 horas. Es el mismo secreto (`JWT_SECRET`) en `api` y `realtime`, así que un token emitido por el login de Go es válido directamente contra el WebSocket de Node.

### Formato de error

Todas las respuestas de error de `api` tienen esta forma:

```json
{ "error": "mensaje descriptivo" }
```

---

## Referencia de la API REST

Prefijo `/api/v1`. Rutas mostradas relativas a ese prefijo; con el gateway levantado, la base es `http://localhost:8000/api/v1`, en modo manual es `http://localhost:8080/api/v1`.

🔒 = requiere cookie de sesión válida · 🔓 = público · 👤 = requiere ser el dueño del recurso (si no, responde `404` para no revelar existencia)

### Auth

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/auth/login` | 🔓 | `{ email, password }` → usuario + cookie `vault_token` |

### Users

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/users` | 🔓 | Registro. `{ name, email, password, avatar_url?, role? }` (`role` ∈ `usuario`\|`vendedor`\|`restaurador`\|`servicio`\|`admin`, default `usuario`) |
| GET | `/users` | 🔒 | Lista todos los usuarios |
| GET | `/users/{id}` | 🔒 | Un usuario por id |
| PUT | `/users/{id}` | 🔒 | Actualiza `{ name, avatar_url, role }` (no email/password) |
| DELETE | `/users/{id}` | 🔒 | Elimina un usuario |
| PUT | `/users/{id}/image` | 🔒 | `multipart/form-data`, campo `image` (máx. 5MB) → sube a Cloudinary y actualiza `avatar_url` |

Respuesta de usuario (`UserResponse`):
```json
{ "id": "uuid", "name": "...", "email": "...", "avatar_url": "...", "role": "usuario" }
```

### Assets (productos)

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/assets` | 🔒 | Crea un producto para el usuario autenticado |
| GET | `/assets` | 🔓 | Lista todos los productos |
| GET | `/assets/{id}` | 🔓 | Un producto con sus fotos |
| PUT | `/assets/{id}` | 🔒👤 | Actualiza |
| DELETE | `/assets/{id}` | 🔒👤 | Elimina |
| POST | `/assets/{id}/photos` | 🔒👤 | `multipart/form-data`, campo `image` → sube a Cloudinary, agrega a `asset_photos` |

Body de creación/actualización:
```json
{
  "name": "Air Jordan 1",
  "category": "sneakers",
  "brand": "Nike",
  "purchase_value": 250.5,
  "condition": "nuevo",
  "purchase_date": "2024-05-10",
  "store_origin": "Nike Store",
  "notes": "Edición limitada"
}
```
`category` ∈ `sneakers`\|`gorras`\|`relojes`\|`lentes`\|`carteras`\|`bolsos`\|`pulsos`\|`bisuteria`\|`coleccionables`\|`otros`. `condition` ∈ `nuevo`\|`seminuevo`\|`usado` (default `nuevo`).

### Businesses

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/businesses` | 🔒 | Crea el negocio del usuario (uno por usuario) |
| GET | `/businesses` | 🔓 | Lista todos |
| GET | `/businesses/{id}` | 🔓 | Uno por id |
| PUT | `/businesses/{id}` | 🔒👤 | Actualiza |
| DELETE | `/businesses/{id}` | 🔒👤 | Elimina |

Body: `{ "name", "type": "restaurador"|"servicio", "description", "location" }`

### Maintenance logs

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/maintenance-logs` | 🔒👤* | Crea un registro de servicio sobre un `asset` propio |
| GET | `/maintenance-logs?asset_id=<uuid>` | 🔓 | Lista los registros de un producto |
| GET | `/maintenance-logs/{id}` | 🔓 | Uno por id |
| PUT | `/maintenance-logs/{id}` | 🔒👤* | Actualiza |
| DELETE | `/maintenance-logs/{id}` | 🔒👤* | Elimina |

\* La propiedad se valida contra el dueño del `asset` relacionado, no del log.

Body: `{ "asset_id", "provider_id"?, "type": "mantenimiento"|"restauracion", "subtype", "cost"?, "performed_at"?, "notes"? }`

### Blockchain certificates

Registros inmutables (sin update/delete).

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/blockchain-certificates` | 🔒👤* | Certifica un `asset` propio |
| GET | `/blockchain-certificates?asset_id=<uuid>` | 🔓 | Lista certificados de un producto |
| GET | `/blockchain-certificates/{id}` | 🔓 | Uno por id |

Body: `{ "asset_id", "tx_id" (único), "asset_hash", "action": "REGISTERED"|"MAINTAINED"|"RESTORED"|"TRANSFERRED", "network"?: "testnet"|"mainnet" }`

### Posts

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/posts` | 🔒 | Crea `{ content, asset_id? }` |
| GET | `/posts` | 🔓 | Feed público (solo visibles) |
| GET | `/posts/{id}` | 🔓 | Uno con sus fotos |
| PUT | `/posts/{id}` | 🔒👤 | Actualiza `{ content }` |
| DELETE | `/posts/{id}` | 🔒👤 | Elimina |
| POST | `/posts/{id}/photos` | 🔒👤 | `multipart/form-data`, campo `image` |
| POST | `/posts/{id}/likes` | 🔒 | Da like (idempotente) |
| DELETE | `/posts/{id}/likes` | 🔒 | Quita el like |

### Comments

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/posts/{id}/comments` | 🔒 | Comenta un post `{ content }` |
| GET | `/posts/{id}/comments` | 🔓 | Lista comentarios visibles de un post |
| DELETE | `/comments/{id}` | 🔒👤 | Elimina |

### Reviews

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/reviews` | 🔒 | Reseña a un proveedor `{ provider_id, content }` |
| GET | `/reviews?provider_id=<uuid>` | 🔓 | Reseñas de un proveedor |
| GET | `/reviews/{id}` | 🔓 | Una por id |
| DELETE | `/reviews/{id}` | 🔒👤 | Elimina |
| POST | `/reviews/{id}/likes` | 🔒 | Da like |
| DELETE | `/reviews/{id}/likes` | 🔒 | Quita el like |

### Notifications

| Método | Ruta | Acceso | Descripción |
|---|---|---|---|
| POST | `/notifications` | 🔒 | Crea una notificación propia (dispara el push en tiempo real, ver abajo) |
| GET | `/notifications` | 🔒 | Lista las notificaciones del usuario autenticado |
| PUT | `/notifications/{id}/read` | 🔒👤 | Marca como leída |
| DELETE | `/notifications/{id}` | 🔒👤 | Elimina |

Body de creación: `{ "type", "subtype", "title", "body", "data"? }` — ver `type`/`subtype` válidos en [`api/init.sql`](api/init.sql).

---

## WebSocket de notificaciones en tiempo real

```
ws://localhost:8000/ws?token=<jwt>     # via gateway (recomendado)
ws://localhost:8081/ws?token=<jwt>     # directo a realtime (modo manual)
```

- Requiere el mismo JWT que emite `POST /api/v1/auth/login` (por query param `?token=`, o por cookie `vault_token` si el cliente la reenvía automáticamente).
- Sin token válido, el servidor rechaza el *upgrade* con `401`.
- No hay que enviar nada tras conectar: el servidor solo empuja mensajes cuando se crea una notificación nueva para ese usuario (`INSERT` en la tabla `notifications`, sin importar qué servicio la haya insertado).
- Cada mensaje es un JSON con la notificación completa:

```json
{
  "id": "uuid",
  "userId": "uuid",
  "type": "comunidad",
  "subtype": "likes_post",
  "title": "...",
  "body": "...",
  "data": {},
  "read": false,
  "createdAt": "2026-07-10T07:46:46.039074Z"
}
```

### Cómo funciona por dentro

1. Un `INSERT` en `notifications` (desde cualquier servicio, hoy solo `api`) dispara un trigger de Postgres que ejecuta `pg_notify('new_notification', <fila como JSON>)`.
2. `realtime` mantiene una conexión persistente con `LISTEN new_notification`.
3. Al recibir el evento, busca los sockets abiertos del `user_id` de la notificación y les reenvía el mensaje.

Esto mantiene a los dos servicios totalmente desacoplados: `api` no sabe que `realtime` existe, y viceversa.

---

## Ejemplo de flujo completo (vía gateway)

```bash
# 1. Registro + login
curl -X POST http://localhost:8000/api/v1/users -H "Content-Type: application/json" \
  -d '{"name":"Ana","email":"ana@vault.test","password":"password123"}'

curl -c cookies.txt -X POST http://localhost:8000/api/v1/auth/login -H "Content-Type: application/json" \
  -d '{"email":"ana@vault.test","password":"password123"}'

# 2. Conectar al WebSocket con el token de esa cookie, por el mismo puerto del gateway
TOKEN=$(grep vault_token cookies.txt | awk '{print $NF}')
websocat "ws://localhost:8000/ws?token=$TOKEN"

# 3. En otra terminal, crear una notificación — llega al instante por el socket
curl -b cookies.txt -X POST http://localhost:8000/api/v1/notifications -H "Content-Type: application/json" \
  -d '{"type":"comunidad","subtype":"likes_post","title":"Hola","body":"Notificacion en vivo"}'
```
