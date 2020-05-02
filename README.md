# Bookmark API

[![Build Status](https://travis-ci.org/firstthumb/bookmark-api.svg?branch=master)](https://travis-ci.org/firstthumb/bookmark-api)

This is simple bookmark api that you could handle your bookmarks easily

## Getting Started

Please follow [the instructions](https://golang.org/doc/install) to install Go on your computer.

After installing Go, run the following commands to start:

```shell
git clone https://github.com/firstthumb/bookmark-api.git

cd bookmark-api

make run
```

At this time, you have a RESTful API server running at `http://127.0.0.1:8080`. It provides the following endpoints:

- `GET /signin/google`: google auth, creates JWT Token
- `POST /bookmarks`: creates new bookmark
- `GET /bookmarks/:id`: returns the detailed information of an bookmark
- `POST /bookmarks/:id/tags/:tag`: adds tag to the bookmark
- `DELETE /bookmarks/:id/tags/:tag`: deletes tag to the bookmark

## DEMO

Create AccessToken with [api.booklog.link/singin/google](https://api.booklog.link/signin/google)

```shell
curl -X GET https://api.booklog.link/api/v1/bookmarks/{BOOKMARK_ID} -H 'Authorization: Bearer {TOKEN}'

curl -X POST https://api.booklog.link/api/v1/bookmarks -H 'Authorization: Bearer {TOKEN}' \
    -d '{
        "name": "Google",
        "url": "https://www.google.com",
        "tags": ["google", "search", "engine"]
    }'

curl -X GET https://api.booklog.link/api/v1/bookmarks?query=Google -H 'Authorization: Bearer {TOKEN}'

curl -X POST https://api.booklog.link/api/v1/bookmarks/{BOOKMARK_ID}/tags/{TAG} -H 'Authorization: Bearer {TOKEN}'

curl -X DELETE https://api.booklog.link/api/v1/bookmarks/{BOOKMARK_ID}/tags/{TAG} -H 'Authorization: Bearer {TOKEN}'

```

## DynamoDB Structure

| ID                  |         RANGE          |               Action |
| ------------------- | :--------------------: | -------------------: |
| USERNAME-{USERNAME} |    NAME-{NAME}-{ID}    |         SearchByName |
| USERNAME-{USERNAME} |     TAG-{TAG}-{ID}     |          SearchByTag |
| USERNAME-{USERNAME} | CREATED-{CREATED}-{ID} | PartitionByCreatedAt |
| USERNAME-{USERNAME} |     BOOKMARK-{ID}      |        Bookmark Data |

I am planning to use [Lambda Store](https://lambda.store/) for caching.

## Project Layout

```
.
├── cmd                  main applications of the project
│   └── bookmark         the API server application
├── config               configuration files for different environments
├── function             lambda functions
│   ├── lambda           lambda main function for HTTP
│   └─- worker           lambda main function for SQS
│   └─- authorizer       lambda authorizer
├── internal             private application
│   ├── auth             auth features
│   ├── bookmark         bookmark features
│   ├── di               wire configuration
│   ├── entity           entity definitions
│   ├── errors           error types
│   ├── session          session operations
│   └── user             user features
├── pkg                  public library code
│   ├── db               database implementation
│   ├── logger           logger
│   └── utils            utilities
```
