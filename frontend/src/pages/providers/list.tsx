import React from 'react';
import { List, useTable, EditButton, ShowButton, DeleteButton } from '@refinedev/antd';
import { Table, Space } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

export const ProviderList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="ID" />
        <Table.Column dataIndex="name" title="Name" />
        <Table.Column dataIndex="consoleUrl" title="Console URL"
          render={(value: string) => value ? (
            <a href={value} target="_blank" rel="noopener noreferrer">
              {value}
            </a>
          ) : '-'}
        />
        <Table.Column dataIndex="notes" title="Notes"
          render={(value: string) => value || '-'}
        />
        <Table.Column
          dataIndex={['_count', 'servers']}
          title="Servers"
          render={(value: number) => value || 0}
        />
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