#!/bin/sh

xdg-open "http://localhost:4444/oauth2/auth?client_id=auth-code-client&scope=openid+offline&response_type=code&state=stateless"