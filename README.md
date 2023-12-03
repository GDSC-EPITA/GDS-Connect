# GDS-Connect
A networking app for students that allows them to anonymously meet other students in their school based on their common interests.

## Goal of this repo
This repo should contain the Backend that will be exposed to communicate with the Frontend part developped by the [GDSC-PLM](https://gdsc.community.dev/pamantasan-ng-lungsod-ng-maynila/) team!
It is a public repository, do not hesitate to contribute using Merge Requests if you want to help the project grow!

## Tech Stack

The tech stack for this project isn't fully determined yet, but we are probably going to use the [Go programming language](https://go.dev/) with the [Gin](https://gin-gonic.com/) framework to manage the server.
This server will interact with a [Firebase](https://firebase.google.com/) database in order to easily handle authentication as well as real-time messaging.

## Contributing to the project

To contribute to the project, please comply to the following guidelines :

### Branches and Issues
* **Never** commit directly to the `master` branch. The branch used for development is the `dev` branch
* To implement a new feature, please open an issue and implement it inside a branch that has the same nomenclature
* After a feature implementation, create a Merge Request for the `dev` branch and wait for a maintainer to approve it

### Commit conventions
* Commit names should consist of a flag in this list: `[feature|fix|remove|refactor]` and a short description
* Commit description have to describe everything that was implemented/done in this commit

## Usage

### Generate the Swagger for the API

The Swagger is generated through comments in the handlers for the API endpoints.
Run `swag init` at the root of the project. This will generate the Swagger by parsing the comments on the API endpoints and place the files in the `docs/` folder.
The swagger will be accessible at `http://localhost:3000/api/swagger/index.html`.
