import express from 'express';
import { SpotifyCredentials } from './credentials.js'
import { SpotifyAuthConfig, redirectToSpotifyAuth, requestSpotifyToken } from './authorization.js'

const app = express();
const port = process.env.APP_PORT || 3000;

app.set('view engine', 'ejs');

const LOGIN_RANDOM_STRING_LEN = 16
const AUTH_URL = 'https://accounts.spotify.com/authorize?'
const TOKEN_URL = 'https://accounts.spotify.com/api/token'
const SCOPE = 'playlist-modify-private playlist-modify-public playlist-read-private';

const REDIRECT_ENDPOINT = process.env.REDIRECT_ENDPOINT
const REDIRECT_URI = process.env.REDIRECT_URI
const SPOTIFY_CLIENT_ID = process.env.SPOTIFY_CLIENT_ID
const SPOTIFY_SECRET = process.env.SPOTIFY_SECRET

if (!REDIRECT_ENDPOINT || !REDIRECT_URI || !SPOTIFY_CLIENT_ID || !SPOTIFY_SECRET) {
  console.error('Missing essential environment variables');
  process.exit(1);
}

const spotifyAuthConfig = new SpotifyAuthConfig(AUTH_URL, TOKEN_URL, REDIRECT_URI, SCOPE)
const spotifyCredentials = new SpotifyCredentials(SPOTIFY_CLIENT_ID, SPOTIFY_SECRET)


app.get('/', (_, response) => {
  response.render("index");
}); 

app.get('/login', (_, response) => {
  redirectToSpotifyAuth(response, spotifyAuthConfig, spotifyCredentials.getClientId(), LOGIN_RANDOM_STRING_LEN)
}); 

app.get(REDIRECT_ENDPOINT, async (request, response) => {

  const authCode = request.query.code;
  const error = request.query.error;

  if (error || !authCode) {
    console.log(`Error requesting access token. ${error}`);
    return response.render("auth_error");
  }

  try {
    const accessToken = await requestSpotifyToken(spotifyAuthConfig, authCode, spotifyCredentials.getBase64Secret())
    response.render("token", {token: accessToken});
  } catch (error) {
    console.log(`Error requesting access token. ${error}`);
    response.render("auth_error");
  }
}); 

app.listen(port, () => {
  console.log(`Server is running on http://localhost:${port}`);
});
