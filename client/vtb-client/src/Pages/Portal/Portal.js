import React from 'react'
import { HttpHandler } from '../../static/js/sendHttpRequest'
import { Redirect } from 'react-router-dom'

export default function Portal() {
  const [state, setState] = React.useState({ login: true })
  const httpHandler = new HttpHandler()

  React.useEffect(() => {
    httpHandler.getProfile()
      .then((response) => response.json())
      .then((data) => {
        if (!data['login'])
          setState(state => ({ ...state, login: false }))
      })

    if (window.localStorage) {
      if (!localStorage.getItem('firstLoad')) {
        localStorage['firstLoad'] = true
        window.location.reload()
      } else {
        localStorage.removeItem('firstLoad')
      }
    }
  }, [])

  if (!state.login) {
    console.log('redirected from portal with no login')
    return <Redirect to='/login'/>
  }

  return (
    <>
      <div>
        Portal
      </div>
      <div>
        Cookie: {document.cookie}
      </div>
    </>
  )
}
