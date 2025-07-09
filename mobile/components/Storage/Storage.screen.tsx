import {HeaderAction, HeaderActions} from '../App/App.header';
import {Screen} from '../App/App.navigation';
import Icon from '../Common/Icon';
import {Layout} from '../Common/Layout';
import Text from '../Common/Text';

import {StorageName} from './Storage.constants';

export const StorageScreenName = StorageName;

export const StorageScreen = () => {
  return (
    <Layout>
      <Text color="textPrimary">Storage Screen</Text>
    </Layout>
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
