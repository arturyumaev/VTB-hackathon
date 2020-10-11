import React from 'react'
import { HttpHandler } from '../../static/js/sendHttpRequest'
import { Redirect } from 'react-router-dom'
import { Typography, InputNumber, Input, Button } from 'antd'
import { FontAwesomeIcon } from '@fortawesome/react-fontawesome'
import { faCreditCard } from '@fortawesome/free-regular-svg-icons'
import { faPaypal } from '@fortawesome/free-brands-svg-icons'
import { faExclamationTriangle } from '@fortawesome/free-solid-svg-icons'

const { Title } = Typography
const paymentIcon = faCreditCard
let token = ''
const cardStyle = {
  border: '0px solid black',
  backgroundColor: '#F8F9FA'
}

export default function Portal() {
  const [state, setState] = React.useState({ login: true })
  const [userData, setUserData] = React.useState({ loaded: false })
  const [accessTokenSuccess, setAccessTokenSuccess] = React.useState(false)
  const [userBalance, setUserBalance] = React.useState(0)
  const [moneyToSend, setMoneyToSend] = React.useState(0)
  const [verificationFailed, setVerificationFailed] = React.useState(false)

  const [code, setCode] = React.useState('')

  const handleSendCode = () => {
    httpHandler.sendVerificationCode(token, code)
      .then(r => r.json())
      .then((d) => {
        if (d['allowed']) {
          setVerificationFailed(false)
          setTimeout(() => {
            loadPaymentInfo()
          }, 1000)
        }
      })
  }

  const httpHandler = new HttpHandler()

  const requstAccessToken = () => {
    httpHandler.getAccessToken()
      .then(r => r.json())
      .then(tokenData => {
        if (tokenData['accessToken']) {
          token = tokenData['accessToken']
          setAccessTokenSuccess(true)
        }

        if (tokenData['error']) {
          setVerificationFailed(true)
        }
      })
  }

  const checkHasAccess = () => {
    httpHandler.getProfile()
      .then(r => r.json())
      .then((data) => {
        if (!data['login']) {
          setState(state => ({ ...state, login: false }))
        } else {
          setUserData(state => ({ ...state, ...data, loaded: true, }))
          requstAccessToken()
        }
      })
  }

  const checkFirstPageVisit = () => {
    if (window.localStorage) {
      if (!localStorage.getItem('firstLoad')) {
        localStorage['firstLoad'] = true
        window.location.reload()
      } else {
        localStorage.removeItem('firstLoad')
      }
    }
  }

  const loadPaymentInfo = () => {
    requstAccessToken()
    setTimeout(() => {
      httpHandler.loadUserBalance(token)
      .then(r => r.json())
      .then((data) => {
        data['balance'] && setUserBalance(data['balance'])
      })
    }, 1000)
  }

  const handleMakeDeposit = () => {
    requstAccessToken()
    setTimeout(() => {
      httpHandler.makeDeposit(moneyToSend, token)
        .then(r => r.json())
        .then(() => loadPaymentInfo())
    }, 1000)
  }

  /* Redirect and reload */
  React.useEffect(() => {
    checkHasAccess()
    checkFirstPageVisit()
  }, [])

  React.useEffect(() => {
    if (accessTokenSuccess) {
      loadPaymentInfo()
    }
    console.log('tokenData изменилась!', accessTokenSuccess)
  }, [accessTokenSuccess])

  if (!state.login) {
    console.log('redirected from portal with no login')
    return <Redirect to='/login'/>
  }

  if (verificationFailed) {
    return (
      <div>
        <div className="d-flex flex-column align-items-center">
          <Title className="mb-4 mt-5" level={4}>
            <FontAwesomeIcon icon={faExclamationTriangle} style={{ color: 'red' }}/>
            &nbsp;
            На Вашем аккаунте замечена подозрительная активность
          </Title>
          <div className="mb-2">
            На Ваш электронный адрес было направлено письмо с кодом подтвреждения операции 
          </div>
          <div>
            Введите код подтверждения ниже, он действителен 10 минут 
          </div>
          <Input placeholder="6-значный код" style={{ width: 200 }} className="mb-2 mt-5" onChange={e => setCode(e.target.value)}/>
          <Button type="primary" className="mb-5 mt-2" onClick={handleSendCode}>Подтвердить</Button>
        </div>
      </div>
      )
  }

  return (
    <>
      {userData.loaded &&
        <div className="container mt-4 d-flex flex-column">
          <Title className="mb-5" level={3}>Добро пожаловать, {userData.name}</Title>
          {accessTokenSuccess &&
            <div>
              {/* Balance */}
              <div className="p-3 mb-3" style={cardStyle}>
                <div>
                  <Title level={5}>Состояние счета</Title>
                </div>
                <div>
                  <FontAwesomeIcon icon={paymentIcon} />&nbsp;Visa Classic
                </div>
                <div>
                  <Title level={5}>{userBalance},00 ₽</Title>
                </div>
              </div>

              {/* Send money */}
              <div className="p-3 mb-3" style={cardStyle}>
                <div>
                  <Title level={5}>Перевести деньги клиенту ВТБ</Title>
                </div>
                <div className="d-flex flex-row">
                  <Input placeholder="Найти пользователя" style={{ width: 200 }}/>
                  &nbsp;
                  &nbsp;
                  <InputNumber
                    defaultValue={0}
                    min={0}
                    max={userBalance}
                    formatter={value => `₽ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                    parser={value => value.replace(/\₽\s?|(,*)/g, '')}
                    onChange={() => {}}
                  />
                  &nbsp;
                  &nbsp;
                  <Button type="primary">Запросить перевод</Button>
                </div>
              </div>
              
              {/* Get money */}
              <div className="p-3 mb-3" style={cardStyle}>
                <div>
                  <Title level={5}>
                    <FontAwesomeIcon icon={faPaypal} />
                    &nbsp;
                    Пополнить счет ВТБ из PayPal
                  </Title>
                </div>
                <div className="d-flex flex-row">
                  <InputNumber
                    defaultValue={0}
                    min={0}
                    max={5000}
                    formatter={value => `₽ ${value}`.replace(/\B(?=(\d{3})+(?!\d))/g, ',')}
                    parser={value => value.replace(/\₽\s?|(,*)/g, '')}
                    onChange={v => setMoneyToSend(v)}
                  />
                  &nbsp;
                  &nbsp;
                  <Button type="primary" onClick={handleMakeDeposit}>Пополнить</Button>
                </div>
              </div>

            </div>
          }
        </div>
      }
    </>
  )
}
