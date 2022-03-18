# rest-api-tutorial

# user-service

# REST API


GET /users -- list of users -- 200, 404, 500
GET /users/:id -- user by id -- 200, 404, 500
POST /users/:id --  create user -- 204, 404, Header Location: url
PUT /users/:id -- fully update user -- 204/200, 404, 400, 500
PATCH /users/:id -- partialy update user -- 204/200, 404, 400, 500
DELETE /users/:id -- delete user by id -- 204, 404, 400
