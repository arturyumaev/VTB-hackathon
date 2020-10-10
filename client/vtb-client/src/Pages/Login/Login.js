import React from 'react'
import { Form, Input, Button } from 'antd';
import { Typography } from 'antd';
import { Alert } from 'antd'
import { UserOutlined, LockOutlined } from '@ant-design/icons';
import { HttpHandler } from '../../static/js/sendHttpRequest'
import { Redirect } from 'react-router-dom'

const { Title } = Typography;

export const Login = () => {
  const [ showAlert, setShowAlert ] = new React.useState({
    type: '',
    msg: '',
    show: false,
    loginSuccess: false
  })

  const httpHandler = new HttpHandler()

  React.useEffect(() => {
    httpHandler.getProfile()
      .then(response => response.json())
      .then((data) => {
        if (data['login']) // Уже были залогинены
          setTimeout(() => setShowAlert(state => ({ ...state, loginSuccess: true })), 500)
      })
    if (document.cookie) {
      setShowAlert(prevState => ({ ...prevState, loginSuccess: true }))
    }
  }, [])

  const disableAlert = () => {
    setTimeout(() => {
      setShowAlert({
        type: '',
        msg: '',
        show: false,
      })
    }, 3500)
  }

  const onFinish = values => {
    httpHandler.logIn(values)
      .then(response => response.json())
      .then((data) => {
        if (data['error']) {
          setShowAlert(prevState => ({
            ...prevState,
            type: 'error',
            msg: 'Неверный логин или пароль',
            show: true,
          }))
          
          disableAlert()
        }

        if (data['status'] && data['status'] == 'OK') {
          setShowAlert({
            type: 'success',
            msg: 'Данные верные',
            show: true,
            loginSuccess: true
          })

          disableAlert()
        }
      });
  };

  if (showAlert.loginSuccess) {
    console.log(document.cookie)
    // document.location.reload()
    return <Redirect to='/portal'/>
  }

  return (
    <div className="container d-flex flex-column align-items-center" style={{ backgroundColor: '#EAEDF5' }}>
      <img className="mt-3" src="https://www.vtb.ru/-/media/Images/Events/Logo/vtb120-028av01-o_logotypes_svg.svg" width="150"/>
      <Title className="mt-2" level={3}>Войдите или зарегестрируйтесь</Title>
      <div className="mt-3 p-4 mb-5">
        <Form name="normal_login" className="login-form" initialValues={{ remember: true }} onFinish={onFinish}>
          
          <Form.Item name="login" rules={[{ required: true, message: 'Введите логин!' }]}>
            <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Логин" />
          </Form.Item>

          <Form.Item name="password" rules={[{ required: true, message: 'Введите пароль!' }]}>
            <Input
              prefix={<LockOutlined className="site-form-item-icon" />}
              type="password"
              placeholder="Пароль"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" className="login-form-button" style={{ width: '100%' }}>
              Войти
            </Button>
            <div className="mt-3">
              <a href="/signup">Регистрация</a>
            </div>
          </Form.Item>
        </Form>
        {showAlert.show && 
          <Alert message={showAlert.msg} type={showAlert.type} />
        }
      </div>
    </div>
  )
}