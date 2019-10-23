# Traefik 2.0 Forward Auth + Auth Host example

A realistic demo that combines:

- `traefik:v2.0`
- `mesosphere/traefik-forward-auth:1.0.4`
- `oryd/hydra:v1`
- `albinodrought/creamy-home-auth-provider`
- `albinodrought/creamy-videos`

## Usage

1. Boot the services: `docker-compose up`

2. Wait for migrations to finish and for services to stabilize. (you're usually good-to-go after seeing `"GET /.well-known/openid-configuration HTTP/1.1" 200`)

3. Make an OAuth client using `./test-001-make-hydra-client.sh`

4. Make a user account using `./test-002-make-user.sh` 

5. Open the demo service using `./test-003-open-login.sh` or by navigating to [http://creamy-videos.test.localhost](http://creamy-videos.test.localhost)

6. You should have been automatically redirected to [http://auth-provider.test.localhost/](http://auth-provider.test.localhost/). There should be a purple login form. Relax, take a breather, everything is probably working.

7. Login with the username `username` and the password `password`.

8. You should have been automatically redirected back to [http://creamy-videos.test.localhost](http://creamy-videos.test.localhost), but this time you should see the [creamy-videos](https://github.com/AlbinoDrought/creamy-videos) UI instead of errors or any login form. It worked!
