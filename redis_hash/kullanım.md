## Görev Oluşturma (POST):

- **URL**: `http://localhost:3000/tasks`
- **Method**: POST
- **Body**: JSON (raw)
    ```json
    {
        "header": "Yeni Görev Başlığı",
        "description": "Yeni görev açıklaması"
    }
    ```

## Tüm Görevleri Listeleme (GET):

- **URL**: `http://localhost:3000/tasks`
- **Method**: GET

## Belirli Bir Görevi Getirme (GET):

- **URL**: `http://localhost:3000/tasks/{id}`
- **Method**: GET
- **Not**: `{id}` kısmını oluşturduğunuz bir görevin ID'siyle değiştirin.

## Görev Güncelleme (PUT):

- **URL**: `http://localhost:3000/tasks/{id}`
- **Method**: PUT
- **Body**: JSON (raw)
    ```json
    {
        "header": "Güncellenmiş Görev Başlığı",
        "description": "Güncellenmiş görev açıklaması"
    }
    ```
- **Not**: `{id}` kısmını güncellemek istediğiniz görevin ID'siyle değiştirin.

## Görev Silme (DELETE):

- **URL**: `http://localhost:3000/tasks/{id}`
- **Method**: DELETE
- **Not**: `{id}` kısmını silmek istediğiniz görevin ID'siyle değiştirin.
