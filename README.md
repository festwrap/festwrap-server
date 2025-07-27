# Overview

The purpose of this application is to facilitate the creation of customized playlist using Golang.

We are relying on Spotify for storing the playlists and Setlistfm for retrieving the top songs from each artist, though we can support other services in the future.

The UI is located in [this other repository](https://github.com/DanielMoraDC/festwrap-ui).

# Local development

Make sure Go 1.24+, Make and Docker v2 are available in your system.

We use pre-commit for static code analysis. Make sure hooks are installed (i.e. `brew install pre-commit` in MacOS) before contributing:

```shell
make pre-commit-install
```

# Testing the code

You can run the tests by typing:

```shell
make run-tests
```

# Running the code

## First time settings

Before you start, make sure you have added the corresponding required variables in `.env`. Make a copy from the template:

```shell
cp .env.template .env
```

Here are the variables you will need to fill:

- `SPOTIFY_CLIENT_ID`: Your Spotify app client id. Follow [these instructions](https://developer.spotify.com/documentation/web-api/tutorials/getting-started#create-an-app) to create your app.
- `SPOTIFY_CLIENT_SECRET`: Your Spotify app client secret. See previous variable for instructions.
- `SPOTIFY_REFRESH_TOKEN`: Spotify refresh token. See [these instructions](https://developer.spotify.com/documentation/web-api/tutorials/refreshing-tokens) on how to obtain it.
- `FESTWRAP_SETLISTFM_APIKEY`: Your Setlistfm API key. It can be requested [here](https://api.setlist.fm/docs/1.0/index.html) for free for non-commercial projects as this one.


### Run the app

To run the API locally, you can type:

```shell
make run-local-server
```

### Run the app container

Start the container by typing:

```shell
make run-server
```

To stop the container, run in a separate terminal:

```shell
make stop-server
```

# Supporting services

Running the app locally will also start a set of supporting services/jobs:
- `pubsub`: local fake of Google Pubsub.
- `pubsub-consumer`: consumes the messages published into the topics.
- `pubsub-init`: creates the topics and subscriber group.

In order to see what is being published into pubsub, you can consult the `pubsub-consumer` logs in a separate terminal (note that it takes some time to start):

```shell
docker logs integration-pubsub-consumer-1 -f
```

## Calling the API

### Artists search

```shell
curl --location 'http://localhost:8080/artists/search?name=<artist>'
```

### Add songs

Creating a new playlist with setlists for some artists:

```shell
curl -X PUT --location 'http://localhost:8080/playlists' \
      --header 'Content-Type: application/json' \
      --data '{"artists":[{"name": "<artist_name>"}],"playlist":{"name":"<playlist_name>"}}'
```
