# Post Comment Service

В данный сервисе разработана система для добавления и чтения постов и комментариев с использованием GraphQL, аналогичную комментариям к постам на популярных платформах, таких как Хабр или Reddit.

## Характеристики системы постов

•	Можно просмотреть список постов.
•	Можно просмотреть пост и комментарии под ним.
•	Пользователь, написавший пост, может запретить оставление комментариев к своему посту.

## Характеристики системы комментариев к постам:
•	Комментарии организованы иерархически, позволяя вложенность без ограничений.
•	Длина текста комментария ограничена до, например, 2000 символов.
•	Система пагинации для получения списка комментариев.

### Запросы
- Получение списка всех постов:
```graphql

  query {
    posts {
      id
      title
      content
      author
    }
  }
```
- Получение конкретного поста по ID, с возможностью пагинации комментариев:
```
graphql
query {
     post(id: 1, commentsOffset: 0, commentsLimit: 10) {
        id
        title
        content
        author
        comments {
            id
            content
            author
        }
    }
}
```
- Создание нового поста:
```
graphql
mutation {
  createPost(title: "New Post", content: "This is a new post.", author: "Author") {
    id
    title
    content
    author
  }
}
```
- Обновление существующего поста:
```
graphql
mutation {
  updatePost(id: 1, title: "Updated Post", content: "This is an updated post.") {
    id
    title
    content
    author
  }
}
```
- Отключение комментариев для поста:
```
graphql
mutation {
  disableComments(postId: 1) {
    id
    title
    content
    author
    commentsDisabled
  }
}
```
- Создание нового комментария:
```
graphql
mutation {
  createComment(postId: 1, author: "Commenter", content: "This is a comment.") {
    id
    postId
    author
    content
  }
}
```
- Обновление существующего комментария:
```
graphql
mutation {
  updateComment(id: 1, content: "This is an updated comment.") {
    id
    postId
    author
    content
  }
}
```