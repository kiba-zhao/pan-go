import type {Device} from '@pango/data';
import {HeaderAction} from '../App/ScreenBase';
import Icon from '../Common/Icon';
import {useIsFetching, useQueryClient} from '../Common/ReactQuery';
import {QueryKey} from './ReactQuery';

export const DeviceRefreshHeaderAction = ({id}: {id: Device['id']}) => {
  const isFetching = useIsFetching({
    queryKey: [...QueryKey, id],
  });

  const queryClient = useQueryClient();
  const handlePress = () => {
    if (isFetching > 0) return;
    queryClient.refetchQueries({
      queryKey: [...QueryKey, id],
      type: 'active',
    });
  };

  return (
    <HeaderAction onPress={handlePress} disabled={isFetching > 0}>
      <Icon
        name="refresh-outline"
        size="medium"
        color={isFetching > 0 ? 'textDisabled' : 'textPrimary'}
      />
    </HeaderAction>
  );
};
