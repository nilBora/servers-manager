import React from 'react';
import { Show, TextField, DateField } from '@refinedev/antd';
import { Typography, Table, Tag } from 'antd';
import { useShow, IResourceComponentsProps } from '@refinedev/core';

const { Title } = Typography;

export const PersonShow: React.FC<IResourceComponentsProps> = () => {
  const { queryResult } = useShow();
  const { data, isLoading } = queryResult;

  const record = data?.data;

  return (
    <Show isLoading={isLoading}>
      <Title level={5}>ID</Title>
      <TextField value={record?.id} />
      <Title level={5}>Name</Title>
      <TextField value={record?.name} />
      <Title level={5}>Email</Title>
      {record?.email ? (
        <a href={`mailto:${record.email}`}>{record.email}</a>
      ) : (
        <TextField value="-" />
      )}
      <Title level={5}>Telegram</Title>
      <TextField value={record?.telegram || '-'} />
      <Title level={5}>Created At</Title>
      <DateField value={record?.createdAt} />
      <Title level={5}>Updated At</Title>
      <DateField value={record?.updatedAt} />

      {record?.serversOwned && (
        <>
          <Title level={5}>Owned Servers ({record.serversOwned.length})</Title>
          <Table
            dataSource={record.serversOwned}
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
            <Table.Column
              dataIndex={['provider', 'name']}
              title="Provider"
              render={(value: string) => value || '-'}
            />
            <Table.Column dataIndex="location" title="Location" />
          </Table>
        </>
      )}
    </Show>
  );
};