# snippetbox

A Web-Application created in Golang, following "Let's Go! 2nd Edition" by Alex Edwards.

## Overview

Snippetbox is a web-app where people can post snippets of code or text and share it later (like gist of pastebin).

## Architecture

- MySQL is used as a database with tables for snippets, users and web sessions
- user authentication by username, email and password
- user passwords are stored as hash function in database
- middleware is used to add security headers to all http requests and check if request coming from authorized user
- to chain middleware function we use [alice](https://pkg.go.dev/github.com/justinas/alice@v1.2.0) package
- implementation of helpers that contain error handling, decoder of HTML forms, HTML template renderer
- validation of entered data by user in snippet form and authorization
- using self-signed TLS certificates
- tests for routes and other functions
