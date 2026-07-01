## 🚀 Descripción
<!-- Describe brevemente los cambios introducidos en este PR y por qué son necesarios. -->


## 📌 Tipo de Cambio
<!-- Marca con una [x] la opción que corresponda -->
- [ ] 🚀 Nueva funcionalidad (feature)
- [ ] 🐛 Corrección de errores (bug fix)
- [ ] 🧹 Refactorización / Limpieza de código (refactor)
- [ ] 📝 Documentación
- [ ] 🧪 Pruebas / Tests
- [ ] ⚙️ CI/CD o Configuración interna

## 📋 Checklist de Calidad (Go-specific)
<!-- Por favor, marca todas las que apliquen antes de solicitar revisión -->
- [ ] Mi código compila localmente (`go build ./...`).
- [ ] He ejecutado los linters y no hay errores (`golangci-lint run`).
- [ ] He agregado pruebas unitarias para esta nueva lógica.
- [ ] Todos los tests existentes y nuevos pasan correctamente (`go test ./...`).
- [ ] No he introducido *data races* (verificado con `go test -race ./...`).
- [ ] Actualicé la documentación (Swagger/OpenAPI, README, Postman si aplica).

## 🧪 ¿Cómo se probó?
<!-- Describe las pruebas que realizaste para verificar tus cambios. Incluye endpoints modificados, payloads de ejemplo o comandos ejecutados. -->

### Ejemplo de Request / Response (si aplica)
```http
METHOD /api/v1/endpoint
Content-Type: application/json

{
  "campo": "valor"
}