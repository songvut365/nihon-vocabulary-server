# Nihon Vocabulary Server

## API

### Authentication

- Login `POST /auth/login`
- Register `POST /auth/register`
- Reset Password `POST /auth/password`

### Current User related

- Logout `DELETE /auth/logout`
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
