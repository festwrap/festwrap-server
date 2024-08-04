export class SpotifyCredentials {
  constructor(clientId, secret) {
    this.clientId = clientId;
    this.secret = secret;
  }

  getClientId() {
    return this.clientId;
  }

  getBase64Secret() {
    return new Buffer.from(this.clientId + ':' + this.secret).toString('base64')
  }

}
