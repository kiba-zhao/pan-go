import {HeaderAction, HeaderActions} from '../App/App.header';
import {Screen} from '../App/App.navigation';
import Icon from '../Common/Icon';
import {ScreenLayout} from '../Common/Layout';
import Text from '../Common/Text';

import {DeviceStorageName} from './DeviceStorage.constants';

export const DeviceStorageScreenName = DeviceStorageName;

export const DeviceStorageScreen = () => {
  return (
    <ScreenLayout>
      <Text color="textPrimary">Device Storage Screen</Text>
    </ScreenLayout>
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
      <Icon name="search-sharp" size="medium" color="textPrimary" />
    </HeaderAction>
  );
};
