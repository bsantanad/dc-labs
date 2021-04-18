# User Guide

Welcome to the DPIP System, this first stage of the project will let you
create users, and with a token (given by the api), they will be able to
* remove user
* check user status
* upload images

It's pretty easy to use and pretty straght forward. So let's get to it.

## Installation

Good news! We didn't use any external framework or libraries, so everything
you'll need comes with the standard libraries and modules from Golang.

Obviosuly you need to have [golang installed][install], but after that you
are good to go :)

Just clone the repository and _cd_ into the project dir.

## Start the Server

So in order to start the server, you'll need to do
```bash
go run main.go
```
Pretty easy right? The API is up and running

## Enpoints

Will be using `curl(1)` for the examples, if you need help with the tool
you can always do:
```bash
man curl
```

### /login - POST

Quick example:
```bash
curl -X POST -u <username>:<password> localhost:8080/login
```
Output:
```bash
{
  "message": "Hi <username>, welcome to the DPIP System",
  "token": <token>
}
```

This endpoint will create a user and assign it a token. When called it will
return a greeting and the token. Be sure to copy the token, will use it later.

### /status - GET

Quick example:
```bash
curl -H "Authorization: Bearer <token>" localhost:8080/status
```
Output:
```bash
{
  "message": "Hi <username>, the DPIP System is Up and Running",
  "time": "2021-04-18 16:40:31.562031608 +0000 UTC"
}
```
**Note:** The time will change according to when you created the user. We use UTC to
avoid confussions when using it in different timezones.

This endpoint will show you the status of the user, if not logged in it will
display a message telling the user that it didn't find the token.

### upload - POST
Quick example:
```bash
curl -X POST -F 'data=@<path to image>' \
             -H "Authorization: Bearer <token>" \
             localhost:8080/upload
```
Output:
```bash
{
  "message": "An image has been successfully uploaded :)",
  "filename": <filename>,
  "size": "<size of file> bytes"
}
```

This endpoint will upload an image to the user with the token sent. It will
return the file name, and the size in bytes. A user can have more that one
image uploaded

### logout - DELETE

Quick example:
```bash
curl -X DELETE -H "Authorization: Bearer <token>" localhost:8080/logout
```
Output:
```bash
{
  "message": "Bye <username>, your token has been revoked"
}
```

This will remove the user and the token from the _"database"_, this means that
this user will loose all their images and it's token. The token will become
invalid an will not longer be a valid token in the API

## Extra Notes

If you try to run the endpoints with the wrong http method, you'll
get a 404. If you try to call a none existing enpoint you'll be redirected to
index.

