# GOFF

## Resumen

Una implementacion del servicio de Redis (esta implementacion https://github.com/HDT3213/godis/) para poder acceder a valores de configuracion, feature flags u otro valor que tenga formato JSON, y donde la actualizacion de los mismo es minima

## Comandos implementados
- `PING`
- `GET`
- `HGET`
- `KEYS`

## Dependencias

- Go 1.23 or higher

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

```bash
GET codigo_de_error_por_region
```

### HGET

```bash
HGET 531eed34tvfy7b73a818a234 access_select_country
```
