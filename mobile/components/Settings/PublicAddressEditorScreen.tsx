import {AppSettings, AppSettingsFields} from '@pango/data';
import {useMemo, useState} from 'react';
import {ActivityIndicator, StyleSheet, VirtualizedList} from 'react-native';
import {HeaderActions, HeaderTitle, ScreenViewLayout} from '../App/ScreenBase';
import Alert from '../Common/Alert';
import {TextFieldModal} from '../Common/Field';
import {useTranslation} from '../Common/I18Next';
import Icon from '../Common/Icon';
import Paper from '../Common/Paper';
import {useMutation, useQueryClient} from '../Common/ReactQuery';
import {SectionItem, SectionItemVariant} from '../Common/Section';
import Text, {TextProps} from '../Common/Text';
import TextInput from '../Common/TextInput';
import {save} from '../Spec/AppSettings';
import {QueryKey, useSettings} from './ReactQuery';
import {I18NextProvider} from './ScreenBase';

type AddressInfo = {
  value: string;
  index: number;
};
const PublicAddressEditorScreen = () => {
  const [selectedInfo, setSelectedInfo] = useState<AddressInfo | undefined>(
    void 0,
  );

  const {data, refetch, isRefetching} = useSettings();
  const [filterText, setFilterText] = useState<string>('');

  const addressInfoList = useMemo(() => {
    let list = (data?.publicAddress || []).map((value, index) => ({
      value,
      index,
    }));

    if (filterText.length > 0) {
      list = list.filter(item => item.value.includes(filterText));
    }
    return list;
  }, [data?.publicAddress, filterText]);

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
      <HeaderTitle i18nKey="screen.publicAddressEditor.name" />
      <PublicAddressHeaderActions onNew={handleNew} />
      <ScreenViewLayout>
        <PublicAddressModal
          info={selectedInfo}
          addressList={data?.publicAddress}
          onClose={handleModalClose}
        />
        <VirtualizedList<AddressInfo>
          ListHeaderComponent={
            <PublicAddressHeader value={filterText} onChange={setFilterText} />
          }
          data={addressInfoList}
          renderItem={info => (
            <PublicAddressItem
              addresses={data?.publicAddress as AppSettings['publicAddress']}
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

export default PublicAddressEditorScreen;

type PublicAddressHeaderActionsProps = {
  onNew: () => void;
};
const PublicAddressHeaderActions = ({
  onNew,
}: PublicAddressHeaderActionsProps) => {
  return (
    <HeaderActions>
      <PublicAddressNewHeaderAction onPress={onNew} />
    </HeaderActions>
  );
};

type PublicAddressNewHeaderActionProps = {
  onPress?: TextProps['onPress'];
};
const PublicAddressNewHeaderAction = ({
  onPress,
}: PublicAddressNewHeaderActionProps) => {
  const {t} = useTranslation();
  return (
    <Text color="primary" font="bold" onPress={onPress}>
      {t('action.new')}
    </Text>
  );
};

type PublicAddressHeaderProps = PublicAddressSearchInputProps;
const PublicAddressHeader = ({value, onChange}: PublicAddressHeaderProps) => {
  return (
    <Paper
      bgColor={'transparent'}
      padding={1}
      gap={1}
      style={styles.headerPaper}>
      <PublicAddressSearchInput value={value} onChange={onChange} />
    </Paper>
  );
};

type PublicAddressSearchInputProps = {
  value?: string;
  onChange?: (value: string) => void;
};

const PublicAddressSearchInput = ({
  value,
  onChange,
}: PublicAddressSearchInputProps) => {
  const {t} = useTranslation();
  return (
    <Paper padding={[0, 1, 0, 1]} style={styles.searchInputPaper}>
      <Icon name="search-outline" color="textPrimary" />
      <TextInput
        placeholder={t('screen.publicAddressEditor.searchExample')}
        style={styles.searchInput}
        value={value}
        onChangeText={onChange}
      />
    </Paper>
  );
};

type PublicAddressItemProps = {
  variant?: SectionItemVariant;
  info: AddressInfo;
  onPress?: () => void;
  addresses: AppSettings['publicAddress'];
};
const PublicAddressItem = ({
  info,
  onPress,
  addresses,
  variant = 'row',
}: PublicAddressItemProps) => {
  const {t} = useTranslation();

  const queryClient = useQueryClient();
  const {mutate, isPending} = useMutation({
    mutationFn: async (field: AppSettings['publicAddress']) =>
      await save({publicAddress: field} as AppSettingsFields),
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
      t('screen.publicAddressEditor.deleteAlertMessage', {
        address: info.value,
      }),
      [
        {
          text: t('action.cancel'),
        },
        {
          text: t('action.ok'),
          onPress: handleConfirm,
        },
      ],
    );
  };

  return (
    <SectionItem variant={variant} gap={1} onPress={onPress}>
      <Icon name="create-outline" color="textPrimary" />
      <Text color="textPrimary" style={styles.sectionItemText}>
        {info.value}
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

type PublicAddressModalProps = {
  info?: AddressInfo;
  addressList?: AppSettings['publicAddress'];
  onClose?: () => void;
};
const PublicAddressModal = ({
  info,
  addressList,
  onClose,
}: PublicAddressModalProps) => {
  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async (publicAddress: string[]) =>
      await save({publicAddress} as AppSettingsFields),
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
      labelI18nKey="screen.publicAddressEditor.addressModal.label"
      placeholderI18nKey="screen.publicAddressEditor.addressModal.example"
      helperI18nKey="screen.publicAddressEditor.addressModal.helper"
      onSubmit={handleSubmit}
      onValid={handleValid}
      onClose={onClose}
    />
  );
};

const styles = StyleSheet.create({
  headerPaper: {
    flexDirection: 'row',
    alignItems: 'center',
  },
  searchInputPaper: {
    flex: 1,
    height: '100%',
    flexDirection: 'row',
    alignItems: 'center',
  },
  searchInput: {
    flex: 1,
  },
  sectionItemText: {
    flex: 1,
  },
});
