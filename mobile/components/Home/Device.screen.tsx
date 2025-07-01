import {Text} from 'react-native';

import {Layout} from '../Common/Layout';

import {HeaderActionIcon} from '../Common/Icon';
import {HeaderAction, HeaderActions} from './Home.header';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

export const DeviceName = `device`;

const DeviceScreen = () => {
  return (
    <Layout>
      <Text>Device Screen</Text>
    </Layout>
  );
};

export default DeviceScreen;

const DeviceIcon = newTabBarIcon({
  focusedName: 'radio-sharp',
  defaultName: 'radio-outline',
});

const DeviceHeaderActions = () => {
  return (
    <HeaderActions>
      <DeviceScreenSearchButton />
    </HeaderActions>
  );
};

export const DeviceHomeScreen = () => (
  <Screen
    name={DeviceName}
    component={DeviceScreen}
    options={{
      title: DeviceName.toUpperCase(),
      tabBarIcon: DeviceIcon,
      headerRight: DeviceHeaderActions,
    }}
  />
);

const DeviceScreenSearchButton = () => {
  const handlePress = () => {
    console.log('search');
  };
  return (
    <HeaderAction onPress={handlePress}>
      <HeaderActionIcon name="search-sharp" />
    </HeaderAction>
  );
};
