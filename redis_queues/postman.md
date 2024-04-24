## LPUSH: Listenin sol tarafına bir veya birden fazla öğe eklemek için POST isteği

- **Method**: POST
- **URL**: `http://localhost:3000/task/lpush`
- **Body** (Raw, JSON):
    ```json
    {
        "item": "New task"
    }
    ```

## RPUSH: Listenin sağ tarafına bir veya birden fazla öğe eklemek için POST isteği

- **Method**: POST
- **URL**: `http://localhost:3000/task/rpush`
- **Body** (Raw, JSON):
    ```json
    {
        "item": "New task"
    }
    ```

## LPOP: Listenin sol tarafındaki bir öğeyi kaldırmak için GET isteği

- **Method**: GET
- **URL**: `http://localhost:3000/task/lpop`

## RPOP: Listenin sağ tarafındaki bir öğeyi kaldırmak için GET isteği

- **Method**: GET
- **URL**: `http://localhost:3000/task/rpop`

## LRANGE: Belirtilen aralıktaki öğeleri almak için GET isteği

- **Method**: GET
- **URL**: `http://localhost:3000/task/lrange`

## LINDEX: Belirtilen dizindeki öğeyi almak için GET isteği

- **Method**: GET
- **URL**: `http://localhost:3000/task/lindex/{index}` _(index değerini belirli bir tam sayıya değiştirin)_

## LINSERT: Belirli bir öğeyi belirli bir öğeden önce veya sonra eklemek için POST isteği

- **Method**: POST
- **URL**: `http://localhost:3000/task/linsert`
- **Body** (Raw, JSON):
    ```json
    {
        "before": "New task",
        "after": "Inserted task"
    }
    ```

## LLEN: Listenin uzunluğunu almak için GET isteği

- **Method**: GET
- **URL**: `http://localhost:3000/task/llen`

## LREM: Belirtilen değeri belirli sayıda listeden kaldırmak için POST isteği

- **Method**: POST
- **URL**: `http://localhost:3000/task/lrem`
- **Body** (Raw, JSON):
    ```json
    {
        "value": "New task",
        "count": 1
    }
    ```

## LSET: Belirli bir dizindeki öğeyi değiştirmek için POST isteği

- **Method**: POST
- **URL**: `http://localhost:3000/task/lset/{index}` _(index değerini belirli bir tam sayıya değiştirin)_
- **Body** (Raw, JSON):
    ```json
    {
        "item": "Updated task"
    }
    ```

## LTRIM: Listeyi belirtilen aralıkta kırpma için POST isteği

- **Method**: POST
- **URL**: `http://localhost:3000/task/ltrim`
- **Body** (Raw, JSON):
    ```json
    {
        "start": 0,
        "stop": 2
    }
    ```
