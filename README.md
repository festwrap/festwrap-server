# Overview

This repository is serving two purposes:

- Practise Golang skills and improve familiarity.
- Facilitate the creation of custom playlists for the musical events I attend.

We are relying on Spotify for storing the playlists and Setlistfm for retrieving the top songs from each artist, though we can support other services in the future.

For this initial version I have allowed some technical debt to slip in (e.g. code duplication), but I plan to remove it in incoming iterations. I keep a small backlog of tasks in this [board](https://trello.com/b/Iv1rrKwE/festwrap).

# Local development

We use pre-commit for static code analysis. Make sure hooks are installed (i.e. `brew install pre-commit` in MacOS) before contributing:

```shell
pre-commit install
pre-commit install --hook-type commit-msg
```

# Components

This application has two main components:

- [Frontend](./frontend). For now it is an auxilar web application we use to retrieve the Spotify access token. In the future, we want this to be an interface that supports the customization of the playlist to create.
- [Backend](./backend). At this moment, it contains a set of classes that implement the basic logic for setlist retrieval and playlist modification. In the future, we want the backend to be an API that provides those features to the frontend.

See further details in each of the components folder.
