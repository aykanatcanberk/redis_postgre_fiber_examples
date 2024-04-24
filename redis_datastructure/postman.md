## Strings (Dizgiler)

- **GET /get/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/get/mykey`

- **POST /set**
  - JSON Veri:
    ```json
    {
        "key": "mykey",
        "value": "myvalue"
    }
    ```

## Lists (Listeler)

- **POST /lpush**
  - JSON Veri:
    ```json
    {
        "key": "mylist",
        "value": "item1"
    }
    ```

- **GET /lrange/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/lrange/mylist`

## Sets (Kümeler)

- **POST /sadd**
  - JSON Veri:
    ```json
    {
        "key": "myset",
        "member": "value1"
    }
    ```

- **GET /smembers/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/smembers/myset`

## Sorted Sets (Sıralı Kümeler)

- **POST /zadd**
  - JSON Veri:
    ```json
    {
        "key": "mysortedset",
        "score": 10,
        "member": "value1"
    }
    ```

- **GET /zrange/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/zrange/mysortedset`

## Hashes (Hash Tabloları)

- **POST /hset**
  - JSON Veri:
    ```json
    {
        "key": "myhash",
        "field": "field1",
        "value": "value1"
    }
    ```

- **GET /hgetall/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/hgetall/myhash`

## Bitmaps (Bit Eşlemeleri)

- **POST /setbit**
  - JSON Veri:
    ```json
    {
        "key": "mybitmap",
        "offset": 10,
        "value": 1
    }
    ```

- **GET /getbit/{key}**
  - Parametre: key (anahtar), offset
  - Örnek URL: `http://localhost:8080/getbit/mybitmap?offset=10`

## HyperLogLogs

- **POST /pfadd**
  - JSON Veri:
    ```json
    {
        "key": "myhyperloglog",
        "elements": ["item1", "item2", "item3"]
    }
    ```

- **GET /pfcount/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/pfcount/myhyperloglog`

## Geospatial Indexes (Coğrafi Dizinler)

- **POST /geoadd**
  - JSON Veri:
    ```json
    {
        "key": "mygeospatial",
        "longitude": 30.12345,
        "latitude": 40.6789,
        "member": "location1"
    }
    ```

- **GET /geopos/{key}**
  - Parametre: key (anahtar)
  - Örnek URL: `http://localhost:8080/geopos/mygeospatial`
