# Creamy Home Auth Provider

<a href="https://hub.docker.com/r/albinodrought/creamy-home-auth-provider">
  <img alt="albinodrought/creamy-home-auth-provider Docker Pulls" src="https://img.shields.io/docker/pulls/albinodrought/creamy-home-auth-provider">
</a>
<a href="https://github.com/AlbinoDrought/creamy-home-auth-provider/blob/master/LICENSE">
  <img alt="AGPL-3.0 License" src="https://img.shields.io/github/license/AlbinoDrought/creamy-home-auth-provider">
</a>

Simple auth provider usable alongside [ory/hydra](https://github.com/ory/hydra), intended for home use.

## Building

### Without Docker

```
go get -d -v
go build
```

### With Docker

`docker build -t albinodrought/creamy-home-auth-provider .`

## Running

```
CREAMY_ACCOUNTS_DIRECTORY="./accounts" \
CREAMY_EMAIL_DOMAIN="example.com" \
CREAMY_HTTP_PORT=7000 \
CREAMY_HYDRA_URL="http://hydra.localhost:4445" \
CREMAY_PATH_TO_STATIC_RESOURCES="./sample-login/static" \
CREAMY_PATH_TO_LOGIN_PAGE="./sample-login/login.html" \
CREAMY_PATH_TO_LOGIN_FAILED_PAGE="./sample-login/failed.html" \
./creamy-home-auth-provider serve
```

- `CREAMY_ACCOUNTS_DIRECTORY`: path where accounts data is stored, defaults to `./accounts`

- `CREAMY_EMAIL_DOMAIN`: domain to append to usernames for emails, defaults to `example.com`

- `CREAMY_HTTP_PORT`: port to listen on, defaults to `7000`

- `CREAMY_HYDRA_URL`: path to your Hydra admin API, defaults to `http://hydra.localhost:4445`

- `CREAMY_PATH_TO_STATIC_RESOURCES`: path to mount on the server at `/static`, defaults to `./sample-login/static`

- `CREAMY_PATH_TO_LOGIN_PAGE`: path to HTML login page, defaults to `./sample-login/login.html`. Should contain a form with `username` and `password` fields that `POST`s to itself (`<form method="POST">`)

- `CREAMY_PATH_TO_LOGIN_FAILED_PAGE`: path to HTML login failed page, defaults to `./sample-login/failed.html`

## Managing Accounts

- Add account: `docker run --rm -v $PWD/accounts:/accounts albinodrought/creamy-home-auth-provider add foo bar`

- Remove account: `docker run --rm -v $PWD/accounts:/accounts albinodrought/creamy-home-auth-provider remove foo`

## Developing

`docker-compose up --build` to boot a local dev environment
