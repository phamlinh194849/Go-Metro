# Docs API

Admin Account:
```aiignore
{
    "email": "admin@gometro.vn",
    "password": "!Gometro1234",
    "full_name": "Pham Thi Ngoc Linh",
    "username": "adminGoMetro"
}
```

register /auth/register
```
{
"email": "nguyenmanh@gmail.com",
"password": "@Nguyen231002",
"full_name": "Nguyen Cong Manh",
"username": "manhGoMetro"
}
```
login  /auth/login
```
{
    "email": "nguyenmanh@gmail.com",
    "password": "@Nguyen231002"
}
```

change password /user/password
```aiignore
{
"old_password": "@Nguyen231002",
"new_password": "@Nguyen231003"
}
```