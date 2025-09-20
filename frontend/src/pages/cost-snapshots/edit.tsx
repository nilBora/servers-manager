import React from 'react';
import { Edit, useForm, useSelect } from '@refinedev/antd';
import { Form, Input, InputNumber, Select, DatePicker } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

const { Option } = Select;

export const CostSnapshotEdit: React.FC<IResourceComponentsProps> = () => {
  const { formProps, saveButtonProps } = useForm();

  const { selectProps: serverSelectProps } = useSelect({
    resource: 'servers',
    optionLabel: 'name',
    optionValue: 'id',
  });

  return (
    <Edit saveButtonProps={saveButtonProps}>
      <Form {...formProps} layout="vertical">
        <Form.Item
          label="Server"
          name="serverId"
          rules={[
            {
              required: true,
              message: 'Please select a server!',
            },
          ]}
        >
          <Select {...serverSelectProps} placeholder="Select a server" />
        </Form.Item>
        <Form.Item
          label="Month"
          name="month"
          rules={[
            {
              required: true,
              message: 'Please select the month!',
            },
          ]}
        >
          <DatePicker.MonthPicker style={{ width: '100%' }} />
        </Form.Item>
        <Form.Item
          label="Monthly Cost"
          name="costMonth"
          rules={[
            {
              required: true,
              message: 'Please input the monthly cost!',
            },
          ]}
        >
          <InputNumber
            min={0}
            precision={2}
            style={{ width: '100%' }}
            placeholder="0.00"
            addonBefore="$"
          />
        </Form.Item>
        <Form.Item label="Source" name="source">
          <Input placeholder="e.g., AWS Bill, Manual Entry, etc." />
        </Form.Item>
      </Form>
    </Edit>
  );
};