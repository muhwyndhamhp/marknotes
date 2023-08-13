# What is this?
This is a teeny tiny template repo for my own personal use. This setup meant for a fullstack Go project with templating support with HTMX. 

This repo's goal is to provide opinionated, but tiny template that ready to be used to develop go web apps. I would only provide what's important (for me) and omit what's not. 

## What this repo can handle out of the box
- Hot Reload via [Air](https://github.com/cosmtrek/air) just `make-run`
- Quick Local dev setup including DB via `make setup-local`
- ORM ready using [GORM](https://gorm.io/)
- Routing ready using [Labstack Echo](https://github.com/labstack/echo)
- Exposed directory for resources in the `dist` directory

## What this repo will never handle
- Deployment beyond simple Dockerfile
- Testing

## Prerequisites
- [Air](https://github.com/cosmtrek/air)
- [Docker](https://docs.docker.com/get-started/)

## Getting Started
Create `.env` file (see `env.sample`), then run `make local-setup` and `make run`. That's it :)
