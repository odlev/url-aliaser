This is a back-end application for creating and following link aliases. Authentication is implemented via a jwt token. 
Queries are created via json,
to get jwt token go to the address /login and in json format pass:
"user": "user",
"password": "secret_key"

to save url:
/save - in json format pass
"url": "your_url",
"alias": "your alias or empty fieild for random 4-letter alias (lowercase)"

to delete url:
/delete - 
similarly

to alias redirection - enter the alias in the address bar

Надо будет доделать регистрацию пользователей в бд и через бд передавать юзер айди для генерации токена. Сейчас там чисто рандомайзер для айди)
