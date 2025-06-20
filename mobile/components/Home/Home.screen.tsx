import IonIcon from '@react-native-vector-icons/ionicons';
import {ComponentProps} from 'react';
import {Text} from 'react-native';
import type {ScreenPropsOptions, TabBarIconType} from '../App/App.navigation';
import {ScreenLayout} from '../Template/Screen.layout';
import {HomeName, HomeTitle} from './Home.constants';

export const HomeScreenName = HomeName;

const HomeScreen = () => {
  return (
    <ScreenLayout>
      <Text>Home Screen123</Text>
    </ScreenLayout>
  );
};
export default HomeScreen;

type HomeScreenIconProps = ComponentProps<TabBarIconType>;
const HomeScreenIcon = ({focused, size, color}: HomeScreenIconProps) => (
  <IonIcon
    name={focused ? 'apps-sharp' : 'apps-outline'}
    size={size}
    color={color}
  />
);

export const HomeScreenOptions: ScreenPropsOptions = {
  title: HomeTitle,
  tabBarIcon: HomeScreenIcon,
};
