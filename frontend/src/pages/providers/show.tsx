import React from 'react';
import { Show, TextField, DateField } from '@refinedev/antd';
import { Typography, Table, Space, Tag } from 'antd';
import { useShow, IResourceComponentsProps } from '@refinedev/core';

const { Title } = Typography;

export const ProviderShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>ID</Title>
      <TextField value={record?.id} />
      <Title level={5}>Name</Title>
      <TextField value={record?.name} />
      <Title level={5}>Console URL</Title>
      {record?.consoleUrl ? (
        <a href={record.consoleUrl} target="_blank" rel="noopener noreferrer">
          {record.consoleUrl}
        </a>
      ) : (
        <TextField value="-" />
      )}
      <Title level={5}>Notes</Title>
      <TextField value={record?.notes || '-'} />
      <Title level={5}>Created At</Title>
      <DateField value={record?.createdAt} />
      <Title level={5}>Updated At</Title>
      <DateField value={record?.updatedAt} />

      {record?.servers && (
        <>
          <Title level={5}>Servers ({record.servers.length})</Title>
          <Table
            dataSource={record.servers}
            rowKey="id"
            pagination={false}
            size="small"
          >
            <Table.Column dataIndex="name" title="Name" />
            <Table.Column dataIndex="hostname" title="Hostname" />
            <Table.Column
              dataIndex="status"
              title="Status"
              render={(value: string) => (
                <Tag color={value === 'ACTIVE' ? 'green' : value === 'STANDBY' ? 'orange' : 'red'}>
                  {value}
                </Tag>
              )}
            />
            <Table.Column dataIndex="location" title="Location" />
          </Table>
        </>
      )}
    </Show>
  );
};