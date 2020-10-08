import React from 'react'
import { Form, Input, Button, Checkbox } from 'antd';
import { UserOutlined, LockOutlined } from '@ant-design/icons';

export const Login = () => {
  const onFinish = values => {
    console.log('Received values of form: ', values);
  };

  return (
    <div className="container-fluid d-flex flex-row justify-content-center">
      <div className="mt-5 p-4">
        <Form
          name="normal_login"
          className="login-form"
          initialValues={{ remember: true }}
          onFinish={onFinish}
        >
          <Form.Item
            name="username"
            rules={[{ required: true, message: 'Please input your Username!' }]}
          >
            <Input prefix={<UserOutlined className="site-form-item-icon" />} placeholder="Username" />
          </Form.Item>

          <Form.Item
            name="password"
            rules={[{ required: true, message: 'Please input your Password!' }]}
          >
            <Input
              prefix={<LockOutlined className="site-form-item-icon" />}
              type="password"
              placeholder="Password"
            />
          </Form.Item>

          <Form.Item>
            <Button type="primary" htmlType="submit" className="login-form-button" style={{ width: '100%' }}>
              Log in
            </Button>
            <div className="mt-3">Or <a href="">register now!</a></div>
          </Form.Item>
        </Form>
      </div>
    </div>
  )
}