# Overview

The purpose of this application is to facilitate the creation of customized playlist using Golang.

We are relying on Spotify for storing the playlists and Setlistfm for retrieving the top songs from each artist, though we can support other services in the future.

The UI is located in [this other repository](https://github.com/DanielMoraDC/festwrap-ui).

# Local development

Make sure Go 1.22+ and Make are available in your system.

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

To run the API, you can type:

```shell
go run cmd/main.go
```

Then you can query it locally by typing:

```shell
curl --location 'http://localhost:8080/artists/search?name=<artist>' \
      --header 'Authorization: Bearer <token>'
```

And you will need to fill for the following variables:
- `<artist>`: the artist to search for.
- `<token>`: Spotify token to access the API. Note that this expire after some hours, so they need to be refreshed. This can be obtained following instructions in [here](../frontend/README.md).


## Add songs to playlist

We can add recent setlist songs from an artist using the `main` file:

```shell
go \
    run cmd/scripts/add_songs_to_playlist.go \
    --spotify-token <spotify_token> \
    --setlistfm-key <setlistfm_key> \
    --artist <artist> \
    --playlist-id <playlist_id>
```

Here we explain the parameters to provide and how to get them:
- `<spotify_token>`: See above.
- `<setlistfm_key>`: Setlistfm API token to obtain the latest setlist for an artist. It can be requested [here](https://api.setlist.fm/docs/1.0/index.html) for free for non-commercial projects as this one.
- `<artist>`: The artist to request songs from.
- `<playlist_id>`: Identifier of the playlist to add songs to. To obtain that, go to your playlist and click the "..." button. Then go to `Share -> Copy link to playlist`. The copied content will look like this: `https://open.spotify.com/playlist/<playlist_id>?<params>`.
