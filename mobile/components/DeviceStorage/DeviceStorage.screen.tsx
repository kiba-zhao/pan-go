import {HeaderAction, HeaderActions} from '../App/App.header';
import {Screen} from '../App/App.navigation';
import {HeaderActionIcon} from '../Common/Icon';
import {Layout} from '../Common/Layout';
import {Text} from '../Common/Text';

import {DeviceStorageName} from './DeviceStorage.constants';

export const DeviceStorageScreenName = DeviceStorageName;

export const DeviceStorageScreen = () => {
  return (
    <Layout>
      <Text>Device Storage Screen</Text>
    </Layout>
  );
};

export default DeviceStorageScreen;

const DeviceStorageHeaderActions = () => {
  return (
    <HeaderActions>
      <DeviceStorageHeaderSearchAction />
    </HeaderActions>
  );
};

export const DeviceStorageAppScreen = () => (
  <Screen
    name={DeviceStorageScreenName}
    component={DeviceStorageScreen}
    options={{
      title: DeviceStorageName.toUpperCase(),
      headerRight: DeviceStorageHeaderActions,
    }}
  />
);

const DeviceStorageHeaderSearchAction = () => {
  const handlePress = () => {
    console.log('storage search');
  };
  return (
    <HeaderAction onPress={handlePress}>
      <HeaderActionIcon name="search-sharp" />
    </HeaderAction>
  );
};
