import React from 'react';
import { List, useTable, EditButton, ShowButton, DeleteButton } from '@refinedev/antd';
import { Table, Space, Tag } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

export const ServerList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="ID" />
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column dataIndex="hostname" title="Hostname" />
        <Table.Column dataIndex="ipPublic" title="Public IP" />
        <Table.Column dataIndex="ipPrivate" title="Private IP" />
        <Table.Column dataIndex="port" title="Port" />
        <Table.Column
          dataIndex="status"
          title="Status"
          render={(value: string) => (
            <Tag
              color={
                value === 'ACTIVE'
                  ? 'green'
                  : value === 'STANDBY'
                  ? 'orange'
                  : value === 'TO_DECOM'
                  ? 'red'
                  : 'gray'
              }
            >
              {value}
            </Tag>
          )}
        />
        <Table.Column
          dataIndex={['provider', 'name']}
          title="Provider"
          render={(value: string) => value || '-'}
        />
        <Table.Column
          dataIndex={['owner', 'name']}
          title="Owner"
          render={(value: string) => value || '-'}
        />
        <Table.Column dataIndex="purpose" title="Purpose" />
        <Table.Column dataIndex="location" title="Location" />
        <Table.Column dataIndex="os" title="OS" />
        <Table.Column
          title="Actions"
          dataIndex="actions"
          render={(_, record: any) => (
            <Space>
              <EditButton hideText size="small" recordItemId={record.id} />
              <ShowButton hideText size="small" recordItemId={record.id} />
              <DeleteButton hideText size="small" recordItemId={record.id} />
            </Space>
          )}
        />
      </Table>
    </List>
  );
};