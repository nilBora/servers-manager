import React from 'react';
import { Show, TextField, TagField, DateField } from '@refinedev/antd';
import { Typography, Space, Tag } from 'antd';
import { useShow, IResourceComponentsProps } from '@refinedev/core';

const { Title } = Typography;

export const ServerShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>ID</Title>
      <TextField value={record?.id} />
      <Title level={5}>Name</Title>
      <TextField value={record?.name} />
      <Title level={5}>Hostname</Title>
      <TextField value={record?.hostname} />
      <Title level={5}>Public IP</Title>
      <TextField value={record?.ipPublic} />
      <Title level={5}>Private IP</Title>
      <TextField value={record?.ipPrivate} />
      <Title level={5}>Port</Title>
      <TextField value={record?.port} />
      <Title level={5}>Username</Title>
      <TextField value={record?.username} />
      <Title level={5}>Status</Title>
      <Tag
        color={
          record?.status === 'ACTIVE'
            ? 'green'
            : record?.status === 'STANDBY'
            ? 'orange'
            : record?.status === 'TO_DECOM'
            ? 'red'
            : 'gray'
        }
      >
        {record?.status}
      </Tag>
      <Title level={5}>Purpose</Title>
      <TextField value={record?.purpose} />
      <Title level={5}>Billing Type</Title>
      <TextField value={record?.billingType} />
      <Title level={5}>Provider</Title>
      <TextField value={record?.provider?.name} />
      <Title level={5}>Owner</Title>
      <TextField value={record?.owner?.name} />
      <Title level={5}>Monthly Cost Estimate</Title>
      <TextField value={record?.costMonthEstimated ? `$${record.costMonthEstimated}` : '-'} />
      <Title level={5}>Account</Title>
      <TextField value={record?.account} />
      <Title level={5}>Location</Title>
      <TextField value={record?.location} />
      <Title level={5}>Operating System</Title>
      <TextField value={record?.os} />
      <Title level={5}>CPU</Title>
      <TextField value={record?.cpu} />
      <Title level={5}>RAM</Title>
      <TextField value={record?.ram} />
      <Title level={5}>Storage</Title>
      <TextField value={record?.storage} />
      <Title level={5}>Description</Title>
      <TextField value={record?.description} />
      <Title level={5}>Tags</Title>
      <TextField value={record?.tags} />
      <Title level={5}>Created At</Title>
      <DateField value={record?.createdAt} />
      <Title level={5}>Updated At</Title>
      <DateField value={record?.updatedAt} />
    </Show>
  );
};