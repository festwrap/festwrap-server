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

You need a Setlistfm API key to run the server. It can be requested [here](https://api.setlist.fm/docs/1.0/index.html) for free for non-commercial projects as this one.

### Run the app

To run the API locally, you can type:

```shell
make run-local-server
```

You will need to make sure you have added the corresponding required variables in `.env`. Make a copy from the template:

```shell
cp .env.template .env
```

And then fill accordingly.

### Run the app container

To run the app Docker image, first make sure to build the image:

```shell
make build-image
```

Then start the container:

```shell
FESTWRAP_SETLISTFM_APIKEY=<setlistfm_key> make run-server
```

To stop the container:

```shell
make stop-server
```

## Calling the API

All endpoints require passing a Spotify token to authenticate. Note that this expire after some hours, so they need to be refreshed. This can be obtained following instructions in [here](../frontend/README.md).

### Artists search

```shell
curl --location 'http://localhost:8080/artists/search?name=<artist>' \
      --header 'Authorization: Bearer <token>'
```

### Playlist search

```shell
curl --location 'http://localhost:8080/playlists/search?name=<playlist>' \
      --header 'Authorization: Bearer <token>'
```

### Add songs

For adding setlists to existing playlists:

```shell
curl -X POST --location 'http://localhost:8080/playlists/<playlist_id>' \
      --header 'Authorization: Bearer <token>'
      --header 'Content-Type: application/json' \
--data '{"artists":[{"name": "<artist_name>"}]}
```

For creating a new playlist with setlists:

```shell
curl -X PUT --location 'http://localhost:8080/playlists' \
      --header 'Authorization: Bearer <token>'
      --header 'Content-Type: application/json' \
--data '{"artists":[{"name": "<artist_name>"}],"playlist":{"name":"<playlist_name>","description":"<playlist_description>","isPublic":<true_false>}}
```

> [!IMPORTANT]
> `public` refers to whether the playlist is publicly shown in your profile.
> However, the playlist is still publicly available given the playlist id.
> See [this](https://community.spotify.com/t5/Spotify-for-Developers/Api-to-create-a-private-playlist-doesn-t-work/td-p/5407807) for more details.
