import React from 'react';
import { List, useTable, EditButton, ShowButton, DeleteButton, DateField } from '@refinedev/antd';
import { Table, Space } from 'antd';
import { IResourceComponentsProps } from '@refinedev/core';

export const CostSnapshotList: React.FC<IResourceComponentsProps> = () => {
  const { tableProps } = useTable({
    syncWithLocation: true,
  });

  return (
    <List>
      <Table {...tableProps} rowKey="id">
        <Table.Column dataIndex="id" title="ID" />
        <Table.Column
          dataIndex={['server', 'name']}
          title="Server"
          render={(value: string) => value || '-'}
        />
        <Table.Column
          dataIndex={['server', 'hostname']}
          title="Hostname"
          render={(value: string) => value || '-'}
        />
        <Table.Column
          dataIndex="month"
          title="Month"
          render={(value: string) => new Date(value).toLocaleDateString('uk-UA', {
            year: 'numeric',
            month: 'long'
          })}
        />
        <Table.Column
          dataIndex="costMonth"
          title="Cost"
          render={(value: number) => `$${value?.toFixed(2) || '0.00'}`}
        />
        <Table.Column dataIndex="source" title="Source"
          render={(value: string) => value || '-'}
        />
        <Table.Column
          dataIndex="createdAt"
          title="Created"
          render={(value: string) => <DateField value={value} />}
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