import {Screen} from '../App/App.navigation';

import {HomeName} from './Home.constants';
import {Navigator} from './Home.navigation';

import {AppHomeScreen, AppName} from './App.screen';
import {DeviceHomeScreen} from './Device.screen';
import {SettingsHomeScreen} from './Settings.screen';

export const HomeScreenName = HomeName;

const HomeScreen = () => {
  return (
    <Navigator initialRouteName={AppName}>
      {AppHomeScreen()}
      {DeviceHomeScreen()}
      {SettingsHomeScreen()}
    </Navigator>
  );
};
export default HomeScreen;

export const HomeAppScreen = () => (
  <Screen
    name={HomeScreenName}
    component={HomeScreen}
    options={{headerShown: false}}
  />
);
