import React from 'react';
import { Edit, useForm } from '@refinedev/antd';
import { Form, Input } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

export const PersonEdit: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps } = useForm();

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Name"
          name="name"
          rules={[
            {
              required: true,
              message: 'Please input the person name!',
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Email"
          name="email"
          rules={[
            {
              type: 'email',
              message: 'Please enter a valid email!',
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item label="Telegram" name="telegram">
          <Input placeholder="@username or username" />
        </Form.Item>
      </Form>
    </Edit>
  );
};