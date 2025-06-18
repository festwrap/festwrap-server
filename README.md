# Overview

The purpose of this application is to facilitate the creation of customized playlist using Golang.

We are relying on Spotify for storing the playlists and Setlistfm for retrieving the top songs from each artist, though we can support other services in the future.

The UI is located in [this other repository](https://github.com/DanielMoraDC/festwrap-ui).

# Local development

Make sure Go 1.24+ and Make are available in your system.

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

## Running the API

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

To run the app Docker image, first make sure to build the image:

```shell
make build-image
```

Then start the container:

```shell
make run-server
```

To stop the container:

```shell
make stop-server
```

## Calling the API

### Artists search

```shell
curl --location 'http://localhost:8080/artists/search?name=<artist>'
```

### Add songs

For adding setlists to existing playlists:

```shell
curl -X POST --location 'http://localhost:8080/playlists/<playlist_id>' \
      --header 'Content-Type: application/json' \
      --data '{"artists":[{"name": "<artist_name>"}]}
```

For creating a new playlist with setlists:

```shell
curl -X PUT --location 'http://localhost:8080/playlists' \
      --header 'Content-Type: application/json' \
      --data '{"artists":[{"name": "<artist_name>"}],"playlist":{"name":"<playlist_name>","description":"<playlist_description>","isPublic":<true_false>}}
```

> [!IMPORTANT]
> `public` refers to whether the playlist is publicly shown in your profile.
> However, the playlist is still publicly available given the playlist id.
> See [this](https://community.spotify.com/t5/Spotify-for-Developers/Api-to-create-a-private-playlist-doesn-t-work/td-p/5407807) for more details.
