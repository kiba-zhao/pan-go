import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';
import {withLazy} from '../Common/Component';
import {HeaderBackHomeAction} from '../Home/ScreenBase';

const DeviceAppScreen = () => (
  <Fragment>
    {DeviceEditorAppScreen()}
    {DeviceCreatorAppScreen()}
  </Fragment>
);

export default DeviceAppScreen;

const DeviceEditorScreen = lazy(() => import('./EditorScreen'));
const DeviceEditorHeaderActions = lazy(() =>
  import('./EditorScreen').then(m => ({default: m.DeviceEditorHeaderActions})),
);
export const DeviceEditorScreenName = 'device.editor';

const DeviceEditorAppScreen = () => (
  <Screen
    name={DeviceEditorScreenName}
    component={DeviceEditorScreen}
    options={{
      title: 'Device Settings',
      headerRight: DeviceEditorHeaderRight,
      headerLeft: HeaderBackHomeAction,
    }}
  />
);

const DeviceEditorHeaderRight = withLazy(DeviceEditorHeaderActions);

const DeviceCreatorScreen = lazy(() => import('./CreatorScreen'));
export const DeviceCreatorScreenName = 'device.creator';

const DeviceCreatorAppScreen = () => (
  <Screen
    name={DeviceCreatorScreenName}
    component={DeviceCreatorScreen}
    options={{
      title: 'New Device',
    }}
  />
);
