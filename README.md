Описание операций
1. Query
Получить список постов
query {
  posts {
    id
    title
    content
    allowComments
  }
}


Описание:
Возвращает список всех постов с их основными полями.

Ответ:

{
  "data": {
    "posts": [
      {
        "id": "1",
        "title": "Первый пост",
        "content": "Пример содержимого",
        "allowComments": true
      }
    ]
  }
}

Получить один пост по ID с комментариями
query {
  post(id: "1") {
    id
    title
    content
    comments(limit: 10, offset: 0) {
      id
      content
      parentId
      replies(limit: 5, offset: 0) {
        id
        content
      }
    }
  }
}


Описание:
Возвращает пост по его id, включая комментарии и вложенные ответы.
Для комментариев предусмотрена пагинация:

limit — количество комментариев за один запрос,

offset — смещение.

2. Mutation

Все мутации защищены директивой @auth.
Для их выполнения пользователь должен быть аутентифицирован (например, через токен, cookie или header).

Создать пост
mutation {
  createPost(title: "Новый пост", content: "Текст поста") {
    id
    title
    content
    allowComments
  }
}


Описание:
Создает новый пост. По умолчанию комментарии к нему разрешены.

Добавить комментарий
mutation {
  addComment(postId: "1", content: "Отличная статья!") {
    id
    postId
    content
  }
}


Описание:
Добавляет комментарий к посту.
Параметр parentId можно указать, чтобы оставить ответ на другой комментарий.

Пример с ответом на комментарий:

mutation {
  addComment(postId: "1", parentId: "2", content: "Согласен!") {
    id
    parentId
    content
  }
}

Запретить комментарии к посту
mutation {
  disallowComments(postId: "1")
}


Описание:
Запрещает оставление комментариев к выбранному посту.
Только автор поста может выполнить эту операцию.

3. Subscription
Подписка на новые комментарии к посту
subscription {
  commentAdded(postId: "1") {
    id
    postId
    parentId
    content
  }
}
