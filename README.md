# Nihon Vocabulary Server

## Environment (.env)

- MONGO_URI
- DATABASE_NAME
- VOCABULARY_COLLECTION
- USER_COLLECTION
- SECRET
- PORT

## Docker

```cmd
docker build -t nihon-vocabulary .
docker image ls
docker run -d -p 8080:8080 --name nv01 --env-file .\.env nihon-vocabulary
```

## APIs

### Authentication

- Login `POST /auth/login`
- Register `POST /auth/register`
- Reset Password `POST /auth/password`

### Current User related

- Show info `GET /user`
- Update info `PUT /user`
- Change password `PUT /user/password`
- Delete account `Delete /user`

### Vocabulary related

- Create Vocabulary `POST /vocabulary`
- Show All Vocabularies `GET /vocabulary`
- Show Vocabulary `GET /vocabulary/:id`
- Update Vocabulary `PUT /vocabulary/:id`
- Delete Vocabulary `DELETE /vocabulary/:id`

## Reference

- [Fiber](http://gofiber.io/)
- [MongoDB Atlas](https://www.mongodb.com/atlas/database)
- [JWT](https://github.com/golang-jwt/jwt)
- [bcrypt](https://pkg.go.dev/golang.org/x/crypto/bcrypt)
