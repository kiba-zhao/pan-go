import type {CursorFields, Storage} from '@pango/data';
import {Fragment, useMemo, useState} from 'react';
import {ActivityIndicator, VirtualizedList} from 'react-native';
import {HeaderTitle, ScreenViewLayout} from '../App/ScreenBase';
import {useTranslation} from '../Common/I18Next';
import {RowLayout} from '../Common/Layout';
import type {QueryKey as ReactQueryKey} from '../Common/ReactQuery';
import {InfiniteData, useInfiniteQuery} from '../Common/ReactQuery';
import {ScreenEmpty} from '../Common/ScreenBase';
import {scale} from '../Common/SizeMatters';
import Text from '../Common/Text';
import {useTheme} from '../Common/Theme';
import type {FetchStoragesResults} from '../Spec/Storage';
import {fetchStorages} from '../Spec/Storage';
import {QueryKey} from './ReactQuery';
import {I18NextProvider} from './ScreenBase';

const StorageItemLimit = 15;

type ExplorerListData = [Storage[], number];
const ExplorerScreen = () => {
  const [q, setQ] = useState('');

  const {data, fetchNextPage, hasNextPage, refetch, isRefetching, isFetching} =
    useInfiniteQuery<
      FetchStoragesResults,
      Error,
      InfiniteData<FetchStoragesResults, CursorFields<string>>,
      ReactQueryKey,
      CursorFields<string>
    >({
      queryKey: [...QueryKey, {q}],
      initialPageParam: {_limit: StorageItemLimit},
      getNextPageParam: lastPage =>
        lastPage[0]?.next !== void 0
          ? {
              _cursor: lastPage[0]?.next,
              _limit: StorageItemLimit,
            }
          : void 0,
      queryFn: async ({pageParam}) =>
        await fetchStorages({
          q,
          ...pageParam,
        }),
    });

  const [entities, total, isDirty] = useMemo(() => {
    if (data === void 0) {
      return [[], 0, false, ''] as [Storage[], number, boolean, string];
    }
    const {pages} = data;
    return pages.reduce<[Storage[], number, boolean, string]>(
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
    <I18NextProvider>
      <HeaderTitle i18nKey="screen.explorer.name" />
      <ScreenViewLayout>
        <VirtualizedList<Storage>
          ListHeaderComponent={() => (
            <ExplorerHeader
              isDirty={isDirty}
              total={total}
              onRefresh={refetch}
            />
          )}
          ListEmptyComponent={ExplorerEmpty}
          ListFooterComponent={hasNextPage ? ExplorerLoading : void 0}
          data={[entities, total]}
          renderItem={info => <ExplorerItem info={info.item} />}
          getItemCount={data => (data as ExplorerListData)[1]}
          getItem={(data: ExplorerListData, idx: number) => data[0][idx]}
          keyExtractor={info => `storage-item-${info.id}-${info.updatedAt}`}
          refreshing={isRefetching}
          onRefresh={refetch}
          onEndReached={handleEndReached}
          windowSize={21}
        />
      </ScreenViewLayout>
    </I18NextProvider>
  );
};

export default ExplorerScreen;

const ExplorerHeader = ({
  isDirty,
  total,
  onRefresh,
}: {
  isDirty: boolean;
  total: number;
  onRefresh?: () => void;
}) => {
  const {t} = useTranslation();
  const {sizes} = useTheme();
  return (
    <Fragment>
      <RowLayout
        style={{
          paddingHorizontal: scale(sizes.base),
          alignItems: 'center',
          justifyContent: 'space-between',
        }}>
        <Text color="textSecondary" size="small" font="bold">
          {t('screen.explorer.totalDesc', {num: total})}
        </Text>
        <Text color="primary" size="small">
          {t('screen.explorer.selectFiles')}
        </Text>
      </RowLayout>
      {isDirty && (
        <Text color="primary" size="small" onPress={onRefresh}>
          {t('screen.explorer.refresh')}
        </Text>
      )}
    </Fragment>
  );
};

const ExplorerEmpty = () => {
  const {t} = useTranslation();
  return <ScreenEmpty>{t('screen.explorer.empty')}</ScreenEmpty>;
};

const ExplorerLoading = () => {
  const {sizes} = useTheme();
  return <ActivityIndicator size={scale(sizes.base) * 3} />;
};

const ExplorerItem = ({info}: {info: Storage}) => {
  return <Text>{info.name}</Text>;
};
