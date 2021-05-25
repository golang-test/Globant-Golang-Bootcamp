# test_task

Запуск с помощью docker-compose up

Затем настройка миграции: \
export POSTGRESQL_URL="postgres://user:123@localhost:5432/postgres?sslmode=disable" \
migrate -database ${POSTGRESQL_URL} -path ./Database/service/migrations up

Пример обращений к сервису:
* Создать книгу \
curl -X POST http://localhost:8080/create -H "Content-type: application/json" -d '{ "name": "book", "price": 2, "genre":1, "amount":5}'
* Обновить книгу \
  curl -X POST http://localhost:8080/update?id=2 -H "Content-type: application/json" -d '{ "name": "book", "price": 2, "genre":1, "amount":5}'
* Получить книгу 
  curl -X GET http://localhost:8080/get_book?id=2 \
  Книга вернется с жанром, указанным текстом
* Получить все книги по запросу
  curl -X GET http://localhost:8080/get_all -H "Content-type: application/json" -d '{ "min___price": 2, "max___price": 10, "names":"book", "genre":1}'
* Удалить книгу
  curl -X POST http://localhost:8080/delete?id=2
