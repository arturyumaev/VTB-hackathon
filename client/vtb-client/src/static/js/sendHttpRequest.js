import Fingerprint2 from '@fingerprintjs/fingerprintjs'

var options = { excludes: {
  touchSupport: true,
  webdriver: true,
  fonts: true,
  audio: true,
  hasLiedBrowser: true,
  hasLiedOs: true,
  hasLiedResolution: true,
  hasLiedLanguages: true,
  webgl: true,
  canvas: true,
  cpuClass: true,
  plugins: true,
  enumerateDevices: true,
  fontsFlash: true,
  adBlock: true,
  doNotTrack: true,
  addBehavior: true,
  openDatabase: true,
  localStorage: true,
  sessionStorage: true,
  pixelRatio: true,
  colorDepth: true,
  indexedDb: true,
}}

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

  getAccessToken = () => {
    var path = '/auth/accessToken'
    var browserFingerprint = ''
    Fingerprint2.get(options, function (components) {
      var values = components.map(c => c.value)
      browserFingerprint = Fingerprint2.x64hash128(values.join(''), 31)
    })
    var data = {scope: ["pay", "addMoney", "balance"], analytics: { fingerprint: browserFingerprint }}
    return this._sendHttpRequest(path, this.portAuth, 'POST', data)
  }

  sendVerificationCode = (token, code) => {
    var path = '/auth/verify'
    var browserFingerprint = ''
    Fingerprint2.get(options, function (components) {
      var values = components.map(c => c.value)
      browserFingerprint = Fingerprint2.x64hash128(values.join(''), 31)
    })
    var data = {code, analytics: { fingerprint: browserFingerprint }}
    return this._sendHttpRequest(path, '8080', 'POST', data)
  }

  loadUserBalance = token => {
    var path = '/payment/balance'
    return this._sendHttpRequest(path, '8081', 'GET', undefined, [[ 'Authorization', `Bearer ${token}` ]])
  }

  makeDeposit = (amount, token) => {
    var path = '/payment/addMoney'
    return this._sendHttpRequest(path, '8081', 'POST', { amount }, [[ 'Authorization', `Bearer ${token}` ]])
  }

  getProfile = () => this._sendHttpRequest('/auth/profile', '8080', 'GET', {})

  _sendHttpRequest = (path, port, requestMethod, data, headers) => {
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

    if (headers) {
      headers.map((pair) => {
        httpData['headers'][pair[0]] = pair[1]
      })
    }

    console.log(httpData)
    
    console.log(httpData)
    return fetch(`${this.protocol}://${this.domain}:${port}${path}`, httpData)
  }
}