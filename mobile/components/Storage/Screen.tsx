import {Fragment, lazy} from 'react';
import {Screen} from '../App/Navigation';

const StorageAppScreen = () => <Fragment>{ExplorerAppScreen()}</Fragment>;
export default StorageAppScreen;

const ExplorerScreen = lazy(() => import('./ExplorerScreen'));
export const ExplorerScreenName = 'storage.explorer';
const ExplorerAppScreen = () => (
  <Screen
    name={ExplorerScreenName}
    component={ExplorerScreen}
    options={{
      title: 'Storage',
    }}
  />
);
