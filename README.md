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

### Run the app

To run the API locally, you can type:

```shell
go run cmd/main.go
```

### Run the app container

To run the app Docker image, first make sure to build the image:

```shell
make build-image
```

Then start the container:

```shell
FESTWRAP_SETLISTFM_APIKEY=<setlistfm_key> make run-server
```

The Setlistfm API key can be requested [here](https://api.setlist.fm/docs/1.0/index.html) for free for non-commercial projects as this one.

To stop the container:

```shell
make stop-server
```

## Calling the API

Once the API is up, you can query it locally by typing:

```shell
curl --location 'http://localhost:8080/artists/search?name=<artist>' \
      --header 'Authorization: Bearer <token>'
```

And you will need to fill for the following variables:
- `<artist>`: the artist to search for.
- `<token>`: Spotify token to access the API. Note that this expire after some hours, so they need to be refreshed. This can be obtained following instructions in [here](../frontend/README.md).
