import React from 'react';
import { Create, useForm, useSelect } from '@refinedev/antd';
import { Form, Input, Select, InputNumber } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

const { Option } = Select;

export const ServerCreate: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps } = useForm();

  const { selectProps: providerSelectProps } = useSelect({
    resource: 'providers',
    optionLabel: 'name',
    optionValue: 'id',
  });

  const { selectProps: ownerSelectProps } = useSelect({
    resource: 'people',
    optionLabel: 'name',
    optionValue: 'id',
  });

  return (
    <Create saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Name"
          name="name"
          rules={[
            {
              required: true,
              message: 'Please input the server name!',
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item
          label="Hostname"
          name="hostname"
          rules={[
            {
              required: true,
              message: 'Please input the hostname!',
            },
          ]}
        >
          <Input />
        </Form.Item>
        <Form.Item label="Public IP Address" name="ipPublic">
          <Input />
        </Form.Item>
        <Form.Item label="Private IP Address" name="ipPrivate">
          <Input />
        </Form.Item>
        <Form.Item label="Port" name="port">
          <InputNumber min={1} max={65535} defaultValue={22} />
        </Form.Item>
        <Form.Item label="Username" name="username">
          <Input />
        </Form.Item>
        <Form.Item label="Password" name="password">
          <Input.Password />
        </Form.Item>
        <Form.Item label="SSH Key" name="sshKey">
          <Input.TextArea rows={4} />
        </Form.Item>
        <Form.Item label="Status" name="status">
          <Select defaultValue="ACTIVE">
            <Option value="ACTIVE">Active</Option>
            <Option value="STANDBY">Standby</Option>
            <Option value="TO_DECOM">To Decommission</Option>
          </Select>
        </Form.Item>
        <Form.Item label="Purpose" name="purpose">
          <Select defaultValue="DEV">
            <Option value="PROD">Production</Option>
            <Option value="STAGING">Staging</Option>
            <Option value="DEV">Development</Option>
            <Option value="TEST">Test</Option>
          </Select>
        </Form.Item>
        <Form.Item label="Billing Type" name="billingType">
          <Select defaultValue="MONTHLY">
            <Option value="HOURLY">Hourly</Option>
            <Option value="MONTHLY">Monthly</Option>
            <Option value="SPOT">Spot</Option>
          </Select>
        </Form.Item>
        <Form.Item label="Provider" name="providerId">
          <Select {...providerSelectProps} placeholder="Select a provider" />
        </Form.Item>
        <Form.Item label="Owner" name="ownerId">
          <Select {...ownerSelectProps} placeholder="Select an owner" />
        </Form.Item>
        <Form.Item label="Location" name="location">
          <Input />
        </Form.Item>
        <Form.Item label="Monthly Cost Estimate" name="costMonthEstimated">
          <InputNumber min={0} precision={2} />
        </Form.Item>
        <Form.Item label="Decommission Date" name="decommissionAt">
          <Input type="date" />
        </Form.Item>
        <Form.Item label="Account" name="account">
          <Input />
        </Form.Item>
        <Form.Item label="Operating System" name="os">
          <Input />
        </Form.Item>
        <Form.Item label="CPU" name="cpu">
          <Input />
        </Form.Item>
        <Form.Item label="RAM" name="ram">
          <Input />
        </Form.Item>
        <Form.Item label="Storage" name="storage">
          <Input />
        </Form.Item>
        <Form.Item label="Description" name="description">
          <Input.TextArea rows={3} />
        </Form.Item>
        <Form.Item label="Tags" name="tags">
          <Input placeholder="Comma separated tags" />
        </Form.Item>
      </Form>
    </Create>
  );
};