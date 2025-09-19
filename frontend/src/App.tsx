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