# TrapAdvisor API
[![Build Status](https://travis-ci.org/victorspringer/trapAdvisor.svg?branch=master)](https://travis-ci.org/victorspringer/trapAdvisor)

This API is the result of the final task for the Information Systems subject "Project and Construction of Systems with DBMS" at Unirio.
It is built with Golang (under Domain Driven Design concept) and MySQL.

## Getting started
Before starting the server, you will need to configure the application domain/port, database settings and Facebook app credentials as well. In order to do it, you must go to the `Makefile` and change the environment variables exported under `start` command, accordingly to your own configuration.

To start the API server just run `make start`, then it will be ready to be accessed, by default, at `localhost:8080`.

### ...