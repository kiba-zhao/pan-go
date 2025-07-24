import {HeaderButton} from '@react-navigation/elements';

import type {ViewStyle} from 'react-native';
import {StyleSheet, View} from 'react-native';

import {Fragment, type ComponentProps} from 'react';
import {CommonActions, useNavigation} from '../App/Navigation';
import Icon from '../Common/Icon';
import {scale} from '../Common/SizeMatters';
import {withTheme} from '../Common/StyleSheet';
import {useTheme} from '../Common/Theme';

export const HeaderAction = withTheme(
  HeaderButton,
  ({sizes}) =>
    ({
      padding: scale(sizes.base),
    } as ViewStyle),
);

export const HeaderActions = withTheme(
  View,
  ({sizes}) =>
    ({
      alignItems: 'flex-end',
      flexDirection: 'row',
      marginRight: scale(sizes.base) * -1,
    } as ViewStyle),
);

type HeaderBackActionProps = Omit<
  ComponentProps<typeof HeaderButton>,
  'children'
> & {
  action?: CommonActions.Action;
};
export const HeaderBackAction = ({
  action = CommonActions.goBack(),
  onPress,
  style,
  ...props
}: HeaderBackActionProps) => {
  const {sizes} = useTheme();
  const navigation = useNavigation();
  const handlePress = () => {
    onPress?.();
    navigation.dispatch(action);
  };
  return (
    <Fragment>
      <HeaderButton
        style={StyleSheet.compose(
          {
            padding: scale(sizes.base),
            marginLeft: scale(sizes.base) * -1,
          },
          style,
        )}
        onPress={handlePress}
        {...props}>
        <Icon name="arrow-back-sharp" color="textPrimary" size="medium" />
      </HeaderButton>
      <View style={{width: scale(sizes.base * 2)}} />
    </Fragment>
  );
};
