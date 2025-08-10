import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';
import {
  DeviceCreatorScreenName,
  DeviceEditorScreenName,
  DeviceSearchScreenName,
} from './ScreenRoute';

const DeviceAppScreen = () => (
  <Fragment>
    {DeviceEditorAppScreen()}
    {DeviceCreatorAppScreen()}
    {DeviceSearchAppScreen()}
  </Fragment>
);

export default DeviceAppScreen;

const DeviceEditorScreen = lazy(() => import('./EditorScreen'));

const DeviceEditorAppScreen = () => (
  <Screen
    name={DeviceEditorScreenName}
    component={DeviceEditorScreen}
    options={{
      title: 'Device Settings',
      animation: 'slide_from_right',
    }}
  />
);

const DeviceCreatorScreen = lazy(() => import('./CreatorScreen'));

const DeviceCreatorAppScreen = () => (
  <Screen
    name={DeviceCreatorScreenName}
    component={DeviceCreatorScreen}
    options={{
      title: 'New Device',
    }}
  />
);

const DeviceSearchScreen = lazy(() => import('./SearchScreen'));

const DeviceSearchAppScreen = () => (
  <Screen
    name={DeviceSearchScreenName}
    component={DeviceSearchScreen}
    options={{
      headerShown: false,
    }}
  />
);
