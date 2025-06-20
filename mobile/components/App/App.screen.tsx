import {
  default as HomeScreen,
  HomeScreenName,
  HomeScreenOptions,
} from '../Home/Home.screen';

import {
  default as DeviceScreen,
  DeviceScreenName,
  DeviceScreenOptions,
} from '../Device/Device.screen';

import {
  default as SettingsScreen,
  SettingsScreenName,
  SettingsScreenOptions,
} from '../Settings/Settings.screen';

import {SafeAreaProvider} from 'react-native-safe-area-context';
import {Navigator, Screen} from './App.navigation';

const AppScreen = () => {
  return (
    <SafeAreaProvider>
      <Navigator initialRouteName={HomeScreenName}>
        <Screen
          name={HomeScreenName}
          component={HomeScreen}
          options={HomeScreenOptions}
        />
        <Screen
          name={DeviceScreenName}
          component={DeviceScreen}
          options={DeviceScreenOptions}
        />
        <Screen
          name={SettingsScreenName}
          component={SettingsScreen}
          options={SettingsScreenOptions}
        />
      </Navigator>
    </SafeAreaProvider>
  );
};

export default AppScreen;
