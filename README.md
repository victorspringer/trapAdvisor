# TrapAdvisor API
[![Build Status](https://travis-ci.org/victorspringer/trapAdvisor.svg?branch=master)](https://travis-ci.org/victorspringer/trapAdvisor)

This API is the result of the final task for the Information Systems subject "Project and Construction of Systems with DBMS" at Unirio.
It is built with Golang (under Domain Driven Design concept) and MySQL.

## Getting started
Before starting the server, you will probably need to configure the database settings properly. In order to do it, you must go to `database/database.go` and change the constants (`user, password, protocol, address, port`) accordingly to your own configuration.
You will need also to create a new Facebook app and put your `ClientID` and `ClientSecret` at the `NewService` method located at `authenticating/service.go`.

To start the API server just run `make start`, then it will be ready to be accessed at `localhost:8080`.

### ...