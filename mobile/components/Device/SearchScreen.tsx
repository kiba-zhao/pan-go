import {Fragment} from 'react';
import {VirtualizedList} from 'react-native';

import {CommonActions, useNavigation} from '../App/Navigation';
import Icon from '../Common/Icon';
import Paper from '../Common/Paper';
import Text from '../Common/Text';
import TextInput from '../Common/TextInput';

const DeviceSearchScreen = () => {
  return (
    <VirtualizedList
      ListHeaderComponent={DeviceSearchHeader}
      ListEmptyComponent={DeviceEmpty}
      renderItem={DeviceItem}
      getItemCount={() => 0}
      getItem={() => ''}
    />
  );
};

export default DeviceSearchScreen;

const DeviceItem = () => {
  return (
    <Fragment>
      <Text>DeviceItem</Text>
    </Fragment>
  );
};

const DeviceSearchHeader = () => {
  const navigation = useNavigation();

  const handleCancel = () => {
    navigation.dispatch(CommonActions.goBack());
  };
  return (
    <Paper
      bgColor={'transparent'}
      padding={1}
      gap={1}
      style={{flexDirection: 'row', alignItems: 'center'}}>
      <Paper
        padding={[0, 1, 0, 1]}
        style={{
          flex: 1,
          height: '100%',
          flexDirection: 'row',
          alignItems: 'center',
        }}>
        <Icon name="search-outline" color="textPrimary" />
        <TextInput placeholder="搜索" style={{flex: 1}} />
        <Icon name="close-circle" color="textPrimary" />
      </Paper>
      <Text color="primary" onPress={handleCancel}>
        取消
      </Text>
    </Paper>
  );
};

const DeviceEmpty = () => {
  return (
    <Paper padding={[6, 0]} style={{alignItems: 'center'}}>
      <Text font="bold" color="textDisabled" size="title">
        No Content
      </Text>
    </Paper>
  );
};
