import {AppSettings, AppSettingsFields} from '@pango/data';
import {useMemo, useState} from 'react';
import {ActivityIndicator, VirtualizedList} from 'react-native';
import {HeaderActions, HeaderTitle, ScreenViewLayout} from '../App/ScreenBase';
import Alert from '../Common/Alert';
import {TextFieldModal} from '../Common/Field';
import {useTranslation} from '../Common/I18Next';
import Icon from '../Common/Icon';
import Paper from '../Common/Paper';
import {useMutation, useQueryClient} from '../Common/ReactQuery';
import type {SectionItemVariant} from '../Common/Section';
import {SectionItem} from '../Common/Section';
import type {TextProps} from '../Common/Text';
import Text from '../Common/Text';
import TextInput from '../Common/TextInput';
import {save} from '../Spec/AppSettings';
import {QueryKey, useSettings} from './ReactQuery';
import {I18NextProvider} from './ScreenBase';

type AddressInfo = {
  value: string;
  index: number;
};
const BroadcastAddressEditorScreen = () => {
  const [selectedInfo, setSelectedInfo] = useState<AddressInfo | undefined>(
    void 0,
  );

  const {data, refetch, isRefetching} = useSettings();
  const [filterText, setFilterText] = useState<string>('');

  const addressInfoList = useMemo(() => {
    let list = (data?.broadcastAddress || []).map((value, index) => ({
      value,
      index,
    }));

    if (filterText.length > 0) {
      list = list.filter(item => item.value.includes(filterText));
    }
    return list;
  }, [data?.broadcastAddress, filterText]);

  const handleNew = () => {
    setSelectedInfo({index: -1, value: ''});
  };

  const handleEdit = (info: AddressInfo) => {
    setSelectedInfo(info);
  };

  const handleModalClose = () => {
    setSelectedInfo(void 0);
  };

  return (
    <I18NextProvider>
      <HeaderTitle i18nKey="screen.broadcastAddressEditor.name" />
      <BroadcastAddressHeaderActions onNewPress={handleNew} />
      <ScreenViewLayout>
        <BroadcastAddressModal
          onClose={handleModalClose}
          info={selectedInfo}
          addressList={data?.broadcastAddress}
        />
        <VirtualizedList<AddressInfo>
          ListHeaderComponent={
            <BroadcastAddressHeader
              value={filterText}
              onChange={setFilterText}
            />
          }
          data={addressInfoList}
          renderItem={info => (
            <BroadcastAddressItem
              addresses={
                data?.broadcastAddress as AppSettings['broadcastAddress']
              }
              info={info.item}
              onPress={() => handleEdit(info.item)}
            />
          )}
          getItemCount={data => data.length}
          getItem={(data, index) => data[index]}
          keyExtractor={info => `broadcast-address-item-${info.value}`}
          refreshing={isRefetching}
          onRefresh={refetch}
          windowSize={21}
        />
      </ScreenViewLayout>
    </I18NextProvider>
  );
};

export default BroadcastAddressEditorScreen;

type BroadcastAddressHeaderActionsProps = {
  onNewPress?: TextProps['onPress'];
};
const BroadcastAddressHeaderActions = ({
  onNewPress,
}: BroadcastAddressHeaderActionsProps) => {
  return (
    <HeaderActions>
      <BroadcastAddressNewHeaderAction onPress={onNewPress} />
    </HeaderActions>
  );
};

const BroadcastAddressNewHeaderAction = ({
  onPress,
}: {
  onPress?: TextProps['onPress'];
}) => {
  const {t} = useTranslation();
  return (
    <Text color="primary" font="bold" onPress={onPress}>
      {t('action.new')}
    </Text>
  );
};

type BroadcastAddressHeaderProps = BroadcastAddressSearchInputProps;
const BroadcastAddressHeader = ({
  value,
  onChange,
}: BroadcastAddressHeaderProps) => {
  return (
    <Paper
      bgColor={'transparent'}
      padding={1}
      gap={1}
      style={{flexDirection: 'row', alignItems: 'center'}}>
      <BroadcastAddressSearchInput value={value} onChange={onChange} />
    </Paper>
  );
};

type BroadcastAddressSearchInputProps = {
  value?: string;
  onChange?: (value: string) => void;
};
const BroadcastAddressSearchInput = ({
  value,
  onChange,
}: BroadcastAddressSearchInputProps) => {
  const {t} = useTranslation();
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
        placeholder={t('screen.broadcastAddressEditor.searchPlaceholder')}
        style={{flex: 1}}
        value={value}
        onChangeText={onChange}
      />
    </Paper>
  );
};

type BroadcastAddressItemProps = {
  addresses: AppSettings['broadcastAddress'];
  info?: AddressInfo;
  variant?: SectionItemVariant;
  onPress?: () => void;
};
const BroadcastAddressItem = ({
  addresses,
  info,
  variant = 'row',
  onPress,
}: BroadcastAddressItemProps) => {
  const {t} = useTranslation();

  const queryClient = useQueryClient();
  const {mutate, isPending} = useMutation({
    mutationFn: async (field: AppSettings['broadcastAddress']) =>
      await save({broadcastAddress: field} as AppSettingsFields),
    onSuccess: data => {
      queryClient.setQueryData(QueryKey, data);
    },
    onError: error => {
      Alert.alert(error.name, error.message);
    },
  });

  const handleConfirm = () => {
    const addresses_ = addresses.slice();
    addresses_.splice(info?.index || 0, 1);
    mutate(addresses_);
  };

  const handleDelete = () => {
    Alert.alert(
      t('alert.warningTitle'),
      t('screen.broadcastAddressEditor.deleteAlertMessage', {
        address: info?.value,
      }),
      [
        {text: t('action.cancel')},
        {text: t('action.ok'), onPress: handleConfirm},
      ],
    );
  };

  return (
    <SectionItem variant={variant} gap={1} onPress={onPress}>
      <Icon name="create-outline" color="textPrimary" />
      <Text color="textPrimary" style={{flex: 1}}>
        {info?.value}
      </Text>
      {isPending ? (
        <ActivityIndicator />
      ) : (
        <Text color="primary" onPress={handleDelete}>
          {t('action.delete')}
        </Text>
      )}
    </SectionItem>
  );
};

type BroadcastAddressModalProps = {
  info?: AddressInfo;
  addressList?: AppSettings['broadcastAddress'];
  onClose?: () => void;
};
const BroadcastAddressModal = ({
  info,
  addressList,
  onClose,
}: BroadcastAddressModalProps) => {
  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async (broadcastAddress: string[]) =>
      await save({broadcastAddress} as AppSettingsFields),
    onSuccess: data => {
      queryClient.setQueryData(QueryKey, data);
      onClose?.();
    },
    onError: error => {
      Alert.alert(error.name, error.message);
    },
  });

  const handleSubmit = (address: string) => {
    if (address.length <= 0) return;
    let addresses = addressList?.slice(0) || [];
    if (info?.index === void 0 || info?.index < 0) {
      addresses.push(address);
    } else {
      addresses.splice(info?.index || 0, 1, address);
    }
    console.log(3333, addresses);
    mutate(addresses);
  };

  const handleValid = (text?: string): boolean | Error => {
    if (text === void 0 || text.length <= 0) return false;
    if (addressList?.includes(text)) return false;
    return true;
  };

  return (
    <TextFieldModal
      visible={!!info}
      value={info?.value}
      labelI18nKey="screen.broadcastAddressEditor.addressModal.label"
      onValid={handleValid}
      onSubmit={handleSubmit}
      onClose={onClose}
      keyboardType="numeric"
      placeholderI18nKey="screen.broadcastAddressEditor.addressModal.example"
      helperI18nKey="screen.broadcastAddressEditor.addressModal.helper"
    />
  );
};
