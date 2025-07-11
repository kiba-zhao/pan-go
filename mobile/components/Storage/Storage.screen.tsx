import {HeaderAction, HeaderActions} from '../App/App.header';
import {Screen} from '../App/App.navigation';
import Icon from '../Common/Icon';
import {ScreenLayout} from '../Common/Layout';
import Text from '../Common/Text';

import {StorageName} from './Storage.constants';

export const StorageScreenName = StorageName;

export const StorageScreen = () => {
  return (
    <ScreenLayout>
      <Text color="textPrimary">Storage Screen</Text>
    </ScreenLayout>
  );
};

export default StorageScreen;

const StorageHeaderActions = () => {
  return (
    <HeaderActions>
      <StorageHeaderSearchAction />
    </HeaderActions>
  );
};

export const StorageAppScreen = () => (
  <Screen
    name={StorageScreenName}
    component={StorageScreen}
    options={{
      title: StorageName.toUpperCase(),
      headerRight: StorageHeaderActions,
    }}
  />
);

const StorageHeaderSearchAction = () => {
  const handlePress = () => {
    console.log('storage search');
  };
  return (
    <HeaderAction onPress={handlePress}>
      <Icon name="search-sharp" size="medium" color="textPrimary" />
    </HeaderAction>
  );
};
