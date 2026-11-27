# GoFF

## Resumen

Servicio de Feature Flag, que implementa el protocolo de Redis (tomado de https://github.com/HDT3213/godis/) escrito en Go.
Se puede utilizar para acceder a valores de configuracion del backend, opciones para dispositivos clientes, o cualquier tipo de esquema que tenga formato JSON, y donde la actualizacion de los datos es minima.

## Comandos implementados
- `PING`
- `GET`
- `HGET`
- `KEYS`

## Dependencias

- Go 1.26 or higher

## Correr el servidor
En desarrollo: 
```bash
go run .
```

## Uso
Se puede testear con el cliente de redis

```bash
redis-cli -p 6382
```

Ejemplos de uso 

### GET

```
127.0.0.1:6382> GET 5060292302201
```

### HGET

```
127.0.0.1:6382> HGET 5060292302201 products.categories_old
```
