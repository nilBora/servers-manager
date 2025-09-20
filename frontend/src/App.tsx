import { Refine } from '@refinedev/core';
import {
  notificationProvider,
  ThemedLayoutV2,
  ErrorComponent,
  RefineThemes,
} from '@refinedev/antd';
import { BrowserRouter, Routes, Route, Outlet, Navigate } from 'react-router-dom';
import routerBindings, {
  NavigateToResource,
  UnsavedChangesNotifier,
  DocumentTitleHandler,
} from '@refinedev/react-router-v6';
import dataProvider from '@refinedev/simple-rest';
import { ConfigProvider } from 'antd';

import '@refinedev/antd/dist/reset.css';

import { ServerList, ServerShow, ServerEdit, ServerCreate } from './pages/servers';
import { ProviderList, ProviderShow, ProviderEdit, ProviderCreate } from './pages/providers';
import { PersonList, PersonShow, PersonEdit, PersonCreate } from './pages/people';
import { CostSnapshotList, CostSnapshotShow, CostSnapshotEdit, CostSnapshotCreate } from './pages/cost-snapshots';

function App() {
  return (
    <BrowserRouter>
      <ConfigProvider theme={RefineThemes.Blue}>
        <Refine
          dataProvider={dataProvider('http://localhost:3000')}
          notificationProvider={notificationProvider}
          routerProvider={routerBindings}
          resources={[
            {
              name: 'servers',
              list: '/servers',
              create: '/servers/create',
              edit: '/servers/edit/:id',
              show: '/servers/show/:id',
              meta: {
                canDelete: true,
              },
            },
            {
              name: 'providers',
              list: '/providers',
              create: '/providers/create',
              edit: '/providers/edit/:id',
              show: '/providers/show/:id',
              meta: {
                canDelete: true,
              },
            },
            {
              name: 'people',
              list: '/people',
              create: '/people/create',
              edit: '/people/edit/:id',
              show: '/people/show/:id',
              meta: {
                canDelete: true,
              },
            },
            {
              name: 'cost-snapshots',
              list: '/cost-snapshots',
              create: '/cost-snapshots/create',
              edit: '/cost-snapshots/edit/:id',
              show: '/cost-snapshots/show/:id',
              meta: {
                canDelete: true,
              },
            },
          ]}
          options={{
            syncWithLocation: true,
            warnWhenUnsavedChanges: true,
          }}
        >
          <Routes>
            <Route
              element={
                <ThemedLayoutV2>
                  <Outlet />
                </ThemedLayoutV2>
              }
            >
              <Route index element={<NavigateToResource resource="servers" />} />
              <Route path="/servers">
                <Route index element={<ServerList />} />
                <Route path="create" element={<ServerCreate />} />
                <Route path="edit/:id" element={<ServerEdit />} />
                <Route path="show/:id" element={<ServerShow />} />
              </Route>
              <Route path="/providers">
                <Route index element={<ProviderList />} />
                <Route path="create" element={<ProviderCreate />} />
                <Route path="edit/:id" element={<ProviderEdit />} />
                <Route path="show/:id" element={<ProviderShow />} />
              </Route>
              <Route path="/people">
                <Route index element={<PersonList />} />
                <Route path="create" element={<PersonCreate />} />
                <Route path="edit/:id" element={<PersonEdit />} />
                <Route path="show/:id" element={<PersonShow />} />
              </Route>
              <Route path="/cost-snapshots">
                <Route index element={<CostSnapshotList />} />
                <Route path="create" element={<CostSnapshotCreate />} />
                <Route path="edit/:id" element={<CostSnapshotEdit />} />
                <Route path="show/:id" element={<CostSnapshotShow />} />
              </Route>
              <Route path="*" element={<ErrorComponent />} />
            </Route>
          </Routes>
          <UnsavedChangesNotifier />
          <DocumentTitleHandler />
        </Refine>
      </ConfigProvider>
    </BrowserRouter>
  );
}

export default App;