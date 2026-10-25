 INSERT INTO users (
  email,
  password_hash,
  name
 ) VALUES (
  $1,
  $2,
  $3
 )
 RETURNING *;

 -- name: GetUserByID :one
 SELECT * FROM users WHERE id = $1 LIMIT 1;

 -- name: GetUserByEmail :one
 SELECT * FROM users WHERE emai = $1 LIMIT 1;