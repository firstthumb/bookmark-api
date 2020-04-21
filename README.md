# Bookmark API
[![Build Status](https://travis-ci.org/firstthumb/bookmark-api.svg?branch=master)](https://travis-ci.org/firstthumb/bookmark-api)

This is simple bookmark api that you could handle your bookmarks easily

Features to be implemented
* CRUD operations with bookmark
* Authenticate users with Google Auth 

## Getting Started

Please follow [the instructions](https://golang.org/doc/install) to install Go on your computer. 

After installing Go, run the following commands to start:

```shell
git clone https://github.com/firstthumb/bookmark-api.git

cd bookmark-api

make run
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

* `POST /bookmarks`: creates new bookmark
* `GET /bookmarks/:id`: returns the detailed information of an bookmark
* `PUT /bookmarks/:id`: updates an existing bookmark
* `DELETE /bookmarks/:id`: deletes an bookmark
* `POST /bookmarks/:id/tags/:tag`: adds tag to the bookmark
* `DELETE /bookmarks/:id/tags/:tag`: deletes tag to the bookmark

## Project Layout
 
```
.
├── cmd                  main applications of the project
│   └── bookmark         the API server application
├── config               configuration files for different environments
├── function             lambda functions
│   ├── lambda           lambda main function for HTTP
│   └─- worker           lambda main function for SQS
├── internal             private application
│   ├── bookmark         bookmark features
│   ├── di               wire configuration
│   ├── entity           entity definitions
│   └── errors           error types
├── pkg                  public library code
│   ├── db               database implementation
│   ├── logger           logger
│   └── utils            utilities 
```