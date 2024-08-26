# Overview

This component implements a web application that we can use to retrieve Spotify Access tokens, so we can use them in the backend.

In the future, we want this to be an interface to enable and customize playlist creations for the user.

# Development setup test

Make sure Node 20 and Make are available in your system. Then install dependencies:

```shell
npm install
```

To run this app you need to create your own Spotify app following [these instructions](https://developer.spotify.com/documentation/web-api/tutorials/getting-started#create-an-app).

Then fill the env file (i.e. `.env`) with the corresponding Spotify secrets. Use `.env.sample` to create `.env` with the required variables to run the project.

```text
SPOTIFY_CLIENT_ID=<spotify_client_id>
SPOTIFY_SECRET=<spotify_secret>
SPOTIFY_REDIRECT_URI=<redirect_uri>
```

**Important** You need to put `http://localhost:3000` as a `NEXT_SPOTIFY_REDIRECT_URI` in your local environment to work properly

To run the app, type:

```shell
make run-app
```
