import React from 'react';
import { Show, TextField, DateField } from '@refinedev/antd';
import { Typography } from 'antd';
import { useShow, IResourceComponentsProps } from '@refinedev/core';

const { Title } = Typography;

export const CostSnapshotShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>ID</Title>
      <TextField value={record?.id} />
      <Title level={5}>Server</Title>
      <TextField value={record?.server?.name || '-'} />
      <Title level={5}>Hostname</Title>
      <TextField value={record?.server?.hostname || '-'} />
      <Title level={5}>Month</Title>
      <TextField
        value={record?.month ? new Date(record.month).toLocaleDateString('uk-UA', {
          year: 'numeric',
          month: 'long'
        }) : '-'}
      />
      <Title level={5}>Monthly Cost</Title>
      <TextField value={record?.costMonth ? `$${record.costMonth.toFixed(2)}` : '$0.00'} />
      <Title level={5}>Source</Title>
      <TextField value={record?.source || '-'} />
      <Title level={5}>Created At</Title>
      <DateField value={record?.createdAt} />
      <Title level={5}>Updated At</Title>
      <DateField value={record?.updatedAt} />
    </Show>
  );
};