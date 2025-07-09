import Text from '../Common/Text';

import {Layout} from '../Common/Layout';
import {Screen} from './Home.navigation';
import {newTabBarIcon} from './Home.tab';

export const SettingsName = `settings`;

const SettingsScreen = () => {
  return (
    <Layout>
      <Text color="textPrimary">Settings Screen</Text>
    </Layout>
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
      title: SettingsName.toUpperCase(),
      tabBarIcon: SettingsIcon,
    }}
  />
);
