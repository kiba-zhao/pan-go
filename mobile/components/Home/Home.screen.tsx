import {Screen} from '../App/App.navigation';

import {HomeName} from './Home.constants';
import {Navigator} from './Home.navigation';

import {useTheme} from '../Common/Theme';
import {AppHomeScreen, AppName} from './App.screen';
import {DeviceHomeScreen} from './Device.screen';
import {SettingsHomeScreen} from './Settings.screen';

export const HomeScreenName = HomeName;

const HomeScreen = () => {
  const {colors, sizes} = useTheme();
  return (
    <Navigator
      initialRouteName={AppName}
      screenOptions={{
        headerStyle: {backgroundColor: colors.background},
        headerShadowVisible: false,
        headerTitleStyle: {
          fontSize: sizes.text + sizes.base * 0.5,
        },
      }}>
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
