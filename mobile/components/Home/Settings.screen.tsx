import Text from '../Common/Text';

import {ScreenLayout} from '../Common/Layout';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

export const SettingsName = `home.settings`;

const SettingsScreen = () => {
  return (
    <ScreenLayout>
      <Text color="textPrimary">Settings Screen</Text>
    </ScreenLayout>
  );
};

export default SettingsScreen;

const SettingsIcon = newTabBarIcon({
  focusedName: 'settings-sharp',
  defaultName: 'settings-outline',
});

export const SettingsHomeScreen = () => (
  <Screen
    name={SettingsName}
    component={SettingsScreen}
    options={{
      title: 'Settings',
      tabBarIcon: SettingsIcon,
    }}
  />
);
