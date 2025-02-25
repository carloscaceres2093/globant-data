
## Requisitos

Antes de iniciar asegurate de tener lo siguiente para hacer que todo funcione correctamente

- [Golang 1.17+](https://yunopayments.atlassian.net/wiki/spaces/TECH/pages/2326537/Golang+lenguaje+de+programacion)

---

## Ejecutando el proyecto

Una vez configuradas todas las variables de entorno y teniendo los requisitos, estamos listo para ejecutar el proyecto, pero antes asegurate de tener todas las dependencias del proyecto con el comando

```sh
   go mod tidy
``` 

### Tests
Para ejecutar los test, que se encuentran dentro del directorio `tests`
```sh
   go test ./...
``` 

### Run
Para ejecutar el proyecto de manera normal se usa
```sh
   go run cmd/api/main.go
``` 

---
## Información sobre el proyecto

[Estructura de directorios](https://yunopayments.atlassian.net/wiki/spaces/TECH/pages/2392095/Estructura+de+carpetas+en+servicios+de+golang)