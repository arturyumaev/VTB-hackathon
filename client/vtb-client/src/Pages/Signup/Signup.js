import React from 'react'
import { Form, Input, Button, Space, Typography, Alert } from 'antd';
import { HttpHandler } from '../../static/js/sendHttpRequest';
import { Redirect } from 'react-router-dom'

const { Title } = Typography;

const formItemLayout = {
  labelCol: {
    xs: { span: 24 },
    sm: { span: 12 },
  },
  wrapperCol: {
    xs: { span: 24 },
    sm: { span: 24 },
  },
};

export const Signup = () => {
  const [form] = Form.useForm()
  const [showAlert, setShowAlert] = React.useState({
    type: '',
    msg: '',
    show: false,
    loginSuccess: false
  })
  const httpHandler = new HttpHandler

  React.useEffect(() => {
    httpHandler.getProfile()
      .then(response => response.json())
      .then((data) => {
        if (data['login']) // Уже были залогинены
          setTimeout(() => setShowAlert(state => ({ ...state, loginSuccess: true })), 500)
      })
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
    httpHandler.registerNewUser(values)
      .then(response => response.json())
      .then((data) => {
        if (data['error']) {
          setShowAlert({
            type: 'error',
            msg: 'Пользователь с такими данными уже существует',
            show: true,
          })
          
          disableAlert()
        }

        if (data['status'] && data['status'] == 'OK') {
          setShowAlert({
            type: 'success',
            msg: 'Регистрация прошла успешно',
            show: true,
          })

          disableAlert()
        }
      });
  };

  if (showAlert.loginSuccess) {
    console.log(document.cookie)
    return <Redirect to='/portal'/>
  }

  return (
    <div className="d-flex flex-row justify-content-center">
      <Space align="center" direction="vertical" style={{ width: "100%", backgroundColor: '#EAEDF5' }}>
        <img className="mt-3" src="https://www.vtb.ru/-/media/Images/Events/Logo/vtb120-028av01-o_logotypes_svg.svg" width="150"/>
        <Title className="mt-2" level={3}>Регистрация</Title>
        <div className="mt-3 p-4 mb-5">
          <Form
            {...formItemLayout}
            form={form}
            name="register"
            className="login-form"
            initialValues={{ remember: true }}
            onFinish={onFinish}
            scrollToFirstError
          >
            
            {/* Имя */}
            <Form.Item
              name="name"
              rules={[
                { type: 'string', message: 'Введенное имя некорректное!', },
                { required: true, message: 'Введите имя!',},
              ]}
            >
              <Input placeholder="Имя"/>
            </Form.Item>

            {/* Логин */}
            <Form.Item
              name="login"
              rules={[
                { required: true },
              ]}
            >
              <Input placeholder="Логин"/>
            </Form.Item>
            
            { /* E-mail */}
            <Form.Item
              name="email"
              rules={[
                { type: 'email', message: 'Введенный E-mail некорректный!', },
                { required: true, message: 'Введите E-mail!',},
              ]}
            >
              <Input placeholder="E-mail"/>
            </Form.Item>

            {/* Пароль */}
            <Form.Item
              name="password"
              rules={[ { required: true, message: 'Введите пароль!', }, ]}
              hasFeedback
            >
              <Input.Password placeholder="Пароль"/>
            </Form.Item>

            <Form.Item
              name="confirm"
              dependencies={['password']}
              hasFeedback
              rules={[
                { required: true, message: 'Подтвердите пароль!', },
                ({ getFieldValue }) => ({ validator(rule, value) {
                  if (!value || getFieldValue('password') === value) 
                    return Promise.resolve();
                  
                  return Promise.reject('Пароли не совпадают!');
                },
              }),]}
            >
              <Input.Password placeholder="Повторите пароль"/>
            </Form.Item>

            <Form.Item>
              <Button type="primary" htmlType="submit" style={{ width: '100%' }}>
                Отправить
              </Button>
            </Form.Item>

          </Form>
          {showAlert.show && 
            <Alert message={showAlert.msg} type={showAlert.type} />
          }
        </div>
      </Space>
    </div>
  )
}