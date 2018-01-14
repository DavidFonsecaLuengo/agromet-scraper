# Agromet Scraper

Scraper para extraer datos desde http://agromet.inia.cl/estaciones.php

El scraper almacena (y carga) las estaciones en el archivo stations.json. Este archivo
contiene el id, nombre y coordenadas de cada estación. En la primera ejecución genera
stations.json con:

```bash
$ ./agromet-scraper -upd
```

El scraper recupera tres variables:
* temperatura del aire
* humedad relativa
* velocidad del viento

Las mediciones de cada estación se almacenan en un archivo JSON con nombre igual al id
de la estación (por ejemplo 25-inia.json para la estación Butalcura).

Puedes obtener los datos de una estación para un rango de fechas, indicando el id de la
estación:

```bash
$ ./agromet-scraper -station=25-inia -from=13-01-2018 -to=13-01-2018
```

o pasar el flag all para descargar los datos de todas las estaciones en stations.json:

```bash
$ ./agromet-scraper -all -from=13-01-2018 -to=13-01-2018 -wt=30s
```

El parámetro wt indica al scraper cuánto debe esperar entre cada request (por defecto
1 minuto).
