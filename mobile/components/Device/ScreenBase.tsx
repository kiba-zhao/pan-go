import {CommonActions, useNavigation} from '../App/Navigation';
import {HeaderAction} from '../App/ScreenBase';
import Icon from '../Common/Icon';
import {DeviceSearchScreenName} from './ScreenRoute';

export const DeviceSearchHeaderAction = ({disabled}: {disabled?: boolean}) => {
  const navigation = useNavigation();

  const handlePress = () => {
    navigation.dispatch(CommonActions.navigate(DeviceSearchScreenName));
  };

  return (
    <HeaderAction onPress={handlePress}>
      <Icon
        name="search-sharp"
        color={disabled ? 'textDisabled' : 'textPrimary'}
        size="medium"
      />
    </HeaderAction>
  );
};
