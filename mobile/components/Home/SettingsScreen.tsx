import {AppSettings} from '@pango/data';
import {
  createContext,
  Dispatch,
  Fragment,
  PropsWithChildren,
  useReducer,
} from 'react';
import {CommonActions, useNavigation} from '../App/Navigation';
import {IconFieldItem, TextFieldItem} from '../Common/Field';
import {useTranslation} from '../Common/I18Next';
import {useQuery} from '../Common/ReactQuery';
import {ScreenLayout, ScreenRefreshControl} from '../Common/ScreenBase';
import Section, {SectionHeader, SectionItem} from '../Common/Section';
import Text from '../Common/Text';
import {QueryKey as SettingsQueryKey} from '../Settings/ReactQuery';
import {NameEditorScreenName, QRCodeScreenName} from '../Settings/ScreenRoute';
import {load as loadSettings} from '../Spec/AppSettings';

/**
 *  define context for editor screen
 */

type ScreenState = {
  nameModalVisible?: boolean;
};

type ScreenAction = ScreenState;
const reducer = (state: ScreenState, action: ScreenAction) => ({
  ...state,
  ...action,
});

const Context = createContext<ScreenState>({});
const DispatchContext = createContext<Dispatch<ScreenAction> | null>(null);

const ScreenProvider = ({children}: PropsWithChildren<{}>) => {
  const [state, dispatch] = useReducer(reducer, {});
  return (
    <Context.Provider value={state}>
      <DispatchContext.Provider value={dispatch}>
        {children}
      </DispatchContext.Provider>
    </Context.Provider>
  );
};

/**
 * end define context
 */

const SettingsScreen = () => {
  const {refetch, data} = useQuery({
    queryKey: SettingsQueryKey,
    queryFn: loadSettings,
    placeholderData: _ => _,
  });

  const handleRefresh = () => {
    refetch();
  };

  return (
    <ScreenLayout
      style={{paddingHorizontal: 0}}
      refreshControl={<ScreenRefreshControl onRefresh={handleRefresh} />}>
      <BaseSection settings={data} />
      <NetSection settings={data} />
    </ScreenLayout>
  );
};

export default SettingsScreen;

type BaseSectionProps = {settings?: AppSettings};
const BaseSection = ({settings}: BaseSectionProps) => {
  const {t} = useTranslation();
  const navigation = useNavigation();

  const handleQRCodePress = () => {
    navigation.dispatch(CommonActions.navigate(QRCodeScreenName));
  };

  const handleNamePress = () => {
    navigation.dispatch(CommonActions.navigate(NameEditorScreenName));
  };

  return (
    <Fragment>
      <Section>
        <SectionItem variant="row" onPress={handleQRCodePress}>
          <IconFieldItem
            label={t('screen.settings.sections.qrcode')}
            disabled={!settings}
            name="qr-code-outline"
          />
        </SectionItem>
      </Section>
      <Section>
        <SectionItem variant="row" onPress={handleNamePress}>
          <TextFieldItem
            label={t('screen.settings.sections.name')}
            text={settings?.name}
            editable
          />
        </SectionItem>
        <SectionItem variant="row">
          <IconFieldItem
            label={t('screen.settings.sections.peerId')}
            disabled={!settings}
            name="open-outline"
          />
        </SectionItem>
      </Section>
    </Fragment>
  );
};

type NetSectionProps = BaseSectionProps;
const NetSection = ({settings}: NetSectionProps) => {
  const {t} = useTranslation();
  return (
    <Section>
      <SectionHeader>
        <Text size="small" color="textSecondary">
          {t('screen.settings.sections.netSection')}
        </Text>
      </SectionHeader>
      <SectionItem variant="row">
        <TextFieldItem
          label={t('screen.settings.sections.peerPort')}
          text={settings?.peerPort ? settings.peerPort.toString() : ''}
          editable
        />
      </SectionItem>
      <SectionItem variant="row">
        <IconFieldItem
          label={t('screen.settings.sections.broadcastAddress')}
          disabled={
            !settings ||
            !settings.broadcastAddress ||
            settings.broadcastAddress.length <= 0
          }
          name="chevron-forward-outline"
        />
      </SectionItem>
      <SectionItem variant="row">
        <IconFieldItem
          label={t('screen.settings.sections.publicAddress')}
          disabled={
            !settings ||
            !settings.publicAddress ||
            settings.publicAddress.length <= 0
          }
          name="chevron-forward-outline"
        />
      </SectionItem>
    </Section>
  );
};
