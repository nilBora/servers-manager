import React from 'react';
import { Create, useForm } from '@refinedev/antd';
import { Form, Input } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

export const ProviderCreate: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps } = useForm();

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Name"
          name="name"
          rules={[
            {
              required: true,
              message: 'Please input the provider name!',
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item label="Console URL" name="consoleUrl">
          <Input placeholder="https://console.provider.com" />
        </Form.Item>
        <Form.Item label="Notes" name="notes">
          <Input.TextArea rows={4} placeholder="Additional notes about this provider" />
        </Form.Item>
      </Form>
    </Create>
  );
};