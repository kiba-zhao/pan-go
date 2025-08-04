import type {CursorFields, Device} from '@pango/data';
import debounce from 'lodash.debounce';
import type {Dispatch, PropsWithChildren} from 'react';
import {
  createContext,
  Fragment,
  useCallback,
  useContext,
  useMemo,
  useReducer,
  useState,
} from 'react';
import {ActivityIndicator, VirtualizedList} from 'react-native';

import {CommonActions, useNavigation} from '../App/Navigation';
import Icon from '../Common/Icon';
import {FlexRowLayout} from '../Common/Layout';
import Paper from '../Common/Paper';
import type {
  InfiniteData,
  QueryKey as ReactQueryKey,
} from '../Common/ReactQuery';
import {useInfiniteQuery, useQueryClient} from '../Common/ReactQuery';
import {SectionItem, type SectionItemVariant} from '../Common/Section';
import {scale} from '../Common/SizeMatters';
import Text from '../Common/Text';
import TextInput from '../Common/TextInput';
import {useTheme} from '../Common/Theme';
import type {FetchDevicesResults} from '../Spec/Device';
import {fetchDevices} from '../Spec/Device';
import {QueryKey} from './ReactQuery';
import {DeviceEditorScreenName} from './ScreenRoute';

type ScreenState = {
  q?: string;
};
type ScreenStateAction = ScreenState;

const reducer = (state: ScreenState, action: ScreenStateAction) => ({
  ...state,
  ...action,
});

const Context = createContext<ScreenState>({});
const DispatchContext = createContext<Dispatch<ScreenStateAction> | null>(null);

const ScreenStateProvider = ({children}: PropsWithChildren<{}>) => {
  const [state, dispatch] = useReducer(reducer, {});
  return (
    <Context.Provider value={state}>
      <DispatchContext.Provider value={dispatch}>
        {children}
      </DispatchContext.Provider>
    </Context.Provider>
  );
};

const DeviceSearchScreen = () => {
  return (
    <ScreenStateProvider>
      <DeviceSearchView />
    </ScreenStateProvider>
  );
};

export default DeviceSearchScreen;

const DeviceItemLimit = 14;
type DeviceListData = [Device[], number];
const DeviceSearchView = () => {
  const {q} = useContext(Context);

  const {data, fetchNextPage, hasNextPage, refetch, isRefetching, isFetching} =
    useInfiniteQuery<
      FetchDevicesResults,
      Error,
      InfiniteData<FetchDevicesResults, CursorFields<string>>,
      ReactQueryKey,
      CursorFields<string>
    >({
      queryKey: [...QueryKey, {q}],
      initialPageParam: {_limit: DeviceItemLimit},
      getNextPageParam: lastPage =>
        lastPage[0]?.next !== void 0
          ? {
              _cursor: lastPage[0]?.next,
              _limit: DeviceItemLimit,
            }
          : void 0,
      queryFn: async ({pageParam}) =>
        await fetchDevices({
          q,
          ...pageParam,
        }),
      enabled: q !== void 0 && q.length > 0,
    });

  const [entities, total, isDirty] = useMemo(() => {
    if (q === void 0 || q.length <= 0 || data === void 0) {
      return [[], 0, false, ''] as [Device[], number, boolean, string];
    }
    const {pages} = data;
    return pages.reduce<[Device[], number, boolean, string]>(
      (results, page, idx) => {
        const [meta, entities] = page;
        const [matrix, total, isDirty, _tag] = results;

        if (entities.length > 0) {
          matrix.push(...entities);

          results[1] = total + entities.length;
        }

        if (isDirty) {
          return results;
        }

        if (_tag.length > 0) {
          results[2] = meta.tag !== _tag;
        } else {
          results[3] = meta.tag;
        }

        return results;
      },
      [[], 0, false, ''],
    );
  }, [data, q]);

  const handleEndReached = useMemo(() => {
    if (!hasNextPage || isFetching) return;
    return () => fetchNextPage();
  }, [isFetching, hasNextPage, fetchNextPage]);

  return (
    <VirtualizedList<Device>
      ListHeaderComponent={() => <DeviceSearchHeader isDirty={isDirty} />}
      ListEmptyComponent={DeviceEmpty}
      ListFooterComponent={hasNextPage ? DeviceLoading : void 0}
      data={[entities, total]}
      renderItem={info => <DeviceItem device={info.item} />}
      getItemCount={data => (data as DeviceListData)[1]}
      getItem={(data: DeviceListData, idx: number) => data[0][idx]}
      keyExtractor={info => `device-item-${info.id}-${info.updatedAt}`}
      refreshing={isRefetching}
      onRefresh={refetch}
      onEndReached={handleEndReached}
      windowSize={21}
    />
  );
};

type DeviceItemProps = {device: Device; variant?: SectionItemVariant};
const DeviceItem = ({device, variant = 'row'}: DeviceItemProps) => {
  const navigation = useNavigation();
  const handlePress = () => {
    navigation.dispatch(
      CommonActions.navigate(DeviceEditorScreenName, {
        id: device.id,
      }),
    );
  };
  return (
    <SectionItem variant={variant} gap={1} onPress={handlePress}>
      <Icon name="desktop-outline" color="textPrimary" />
      <FlexRowLayout style={{alignItems: 'center'}}>
        <Text color="textPrimary" style={{flex: 1}}>
          {device.name}
        </Text>
        <Icon name="chevron-forward-outline" size="small" color="textPrimary" />
      </FlexRowLayout>
    </SectionItem>
  );
};

type DeviceSearchHeaderProps = {isDirty?: boolean};
const DeviceSearchHeader = ({isDirty}: DeviceSearchHeaderProps) => {
  return (
    <Fragment>
      <Paper
        bgColor={'transparent'}
        padding={1}
        gap={1}
        style={{flexDirection: 'row', alignItems: 'center'}}>
        <DeviceSearchInput />
        <DeviceCancelAction />
      </Paper>
      {!!isDirty && <DeviceRefechAction />}
    </Fragment>
  );
};

const DeviceCancelAction = () => {
  const navigation = useNavigation();

  const handleCancel = () => {
    navigation.dispatch(CommonActions.goBack());
  };
  return (
    <Text color="primary" onPress={handleCancel}>
      取消
    </Text>
  );
};

const DeviceSearchInput = () => {
  const {q} = useContext(Context);
  const dispatch = useContext(DispatchContext);

  const [inputValue, setInputValue] = useState(q || '');

  const handlePressClear = () => {
    handleChangeText('');
  };

  const handleSearch = useCallback(
    debounce((text: string) => {
      dispatch?.({q: text});
    }, 500),
    [],
  );

  const handleChangeText = (text: string) => {
    setInputValue(text);
    handleSearch(text);
  };

  return (
    <Paper
      padding={[0, 1, 0, 1]}
      style={{
        flex: 1,
        height: '100%',
        flexDirection: 'row',
        alignItems: 'center',
      }}>
      <Icon name="search-outline" color="textPrimary" />
      <TextInput
        placeholder="搜索"
        style={{flex: 1}}
        value={inputValue}
        onChangeText={handleChangeText}
      />
      {inputValue.length > 0 && (
        <Icon
          name="close-circle"
          color="textPrimary"
          onPress={handlePressClear}
        />
      )}
    </Paper>
  );
};

const DeviceEmpty = () => {
  return (
    <Paper
      bgColor={'transparent'}
      padding={[8, 0]}
      style={{alignItems: 'center'}}>
      <Text font="bold" color="textDisabled" size="title">
        No Content
      </Text>
    </Paper>
  );
};

const DeviceLoading = () => {
  const {sizes} = useTheme();
  return <ActivityIndicator size={scale(sizes.base) * 3} />;
};

const DeviceRefechAction = () => {
  const {q} = useContext(Context);
  const queryClient = useQueryClient();

  const handlePress = () => {
    queryClient.refetchQueries({
      queryKey: [...QueryKey, {q}],
      type: 'active',
    });
  };
  return (
    <Paper
      bgColor={'transparent'}
      padding={[0, 0, 1, 0]}
      style={{alignItems: 'center', alignContent: 'center'}}>
      <Text color="primary" onPress={handlePress}>
        获取最新列表
      </Text>
    </Paper>
  );
};
