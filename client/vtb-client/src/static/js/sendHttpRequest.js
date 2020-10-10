export class HttpHandler {
  constructor() {
    this.protocol = 'https'
    this.hostId = '109.234.38.216'
    this.domain = 'vtbpayment.ru'
    this.portAuth = '8080'
    this.portPayment = '8081'
  }

  registerNewUser = data => {
    var path = '/auth/register'
    return this._sendHttpRequest(path, this.portAuth, 'POST', data)
  }

  logIn = data => {
    var path = '/auth/login/self'
    return this._sendHttpRequest(path, this.portAuth, 'POST', data)
  }

  logOut = () => {
    var path = '/auth/logout'
    return this._sendHttpRequest(path, this.portAuth, 'POST', undefined)
  }

  getProfile = () => this._sendHttpRequest('/auth/profile', '8080', 'GET', {})

  _sendHttpRequest = (path, port, requestMethod, data) => {
    const httpData = {
      method: requestMethod,
      headers: {
        'Accept': 'application/json',
        'Content-Type': 'application/json'
      },
      credentials: 'include'
    }

    if (requestMethod == 'POST')
      httpData['body'] = JSON.stringify(data)

    if (requestMethod == 'GET')
     
    
    console.log(httpData)
    return fetch(`${this.protocol}://${this.domain}:${port}${path}`, httpData)
  }
}