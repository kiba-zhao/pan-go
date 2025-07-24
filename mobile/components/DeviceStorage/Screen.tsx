import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';

const DeviceStorageAppScreen = () => <Fragment>{ExplorerAppScreen()}</Fragment>;

export default DeviceStorageAppScreen;

const ExplorerScreen = lazy(() => import('./ExplorerScreen'));
export const ExplorerScreenName = 'deviceStorage.explorer';

const ExplorerAppScreen = () => (
  <Screen
    name={ExplorerScreenName}
    component={ExplorerScreen}
    options={{
      title: 'Device Strorage',
    }}
  />
);
