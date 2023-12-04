# Golang WebApp

This repository contains a Golang-based REST API project integrated with a PostgreSQL database. Additionally, it includes configuration files(workflows) for building an AWS AMI using Packer.

## Table of Contents
  

- [Golang WebApp](#golang-webapp)
  - [Table of Contents](#table-of-contents)
  - [Overview](#overview)
  - [Prerequisites](#prerequisites)
  - [Project Structure](#project-structure)
  - [Usage](#usage)
  - [Workflows](#workflows)


  

## Overview

  

This project demonstrates a REST API written in Golang that interacts with a PostgreSQL database. Additionally, it includes a Packer configuration to build an AWS AMI for the application.

  

## Prerequisites

  

Before getting started, make sure you have the following prerequisites installed:

  

- Golang: [Install Golang](https://golang.org/doc/install)

- PostgreSQL: [Install PostgreSQL](https://www.postgresql.org/download/)

- Packer: [Install Packer](https://www.packer.io/docs/install)

- AWS CLI:[Install AWS CLI](https://aws.amazon.com/cli/)

  

## Project Structure
```
.
├── LICENSE
├── README.md
├── api
│   └── handler
│       └── nocache.go
├── controllers
│   ├── assignment_controller.go
│   ├── authentication.go
│   └── health_controller.go
├── data
│   └── users.csv
├── database
│   ├── account.go
│   ├── assignment.go
│   ├── database.go
│   └── seeder.go
├── go.mod
├── go.sum
├── main.go
├── packer.pkr.hcl
├── packer.pkrvars.hcl
├── scripts
│   └── setup.sh
└── test
    └── healthz_test.go
```
## Usage

- To download the application dependencies run `go mod downlaod`
- To build the executable run `go build`
- To run the executable run `./webapp`
- Command to import the certificate 
```
aws acm import-certificate --profile demo \
  --certificate fileb://demo_ecorplabs_me.crt \
  --private-key fileb://private.key
```

## Workflows

- GitHub workflows are triggered when a pull request is created on the main org repository
- Another GitHub workflow is triggered when a pull request is merged into the main branch of the organization repository

