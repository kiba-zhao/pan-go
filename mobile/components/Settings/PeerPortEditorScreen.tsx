import {AppSettingsFields} from '@pango/data';
import {Fragment, useState} from 'react';
import {View} from 'react-native';
import {useNavigation} from '../App/Navigation';
import {HeaderTitle, ScreenLayout} from '../App/ScreenBase';
import Alert from '../Common/Alert';
import Button from '../Common/Button';
import {useTranslation} from '../Common/I18Next';
import {useMutation, useQueryClient} from '../Common/ReactQuery';
import {ScreenRefreshControl} from '../Common/ScreenBase';
import {scale} from '../Common/SizeMatters';
import Text from '../Common/Text';
import TextInput from '../Common/TextInput';
import {useTheme} from '../Common/Theme';
import {save} from '../Spec/AppSettings';
import {QueryKey, useSettings} from './ReactQuery';
import {I18NextProvider} from './ScreenBase';

const PeerPortEditorScreen = () => {
  const {data, refetch} = useSettings();
  return (
    <I18NextProvider>
      <ScreenLayout
        refreshControl={<ScreenRefreshControl onRefresh={refetch} />}>
        <HeaderTitle i18nKey="screen.peerPortEditor.name" />
        <EditorSection value={data?.peerPort} />
      </ScreenLayout>
    </I18NextProvider>
  );
};

export default PeerPortEditorScreen;

const EditorSection = ({value}: {value?: number}) => {
  const {t} = useTranslation();
  const {sizes, colors} = useTheme();

  const [text, setText] = useState(value ? value.toString() : '');

  const handleChange = (text: string) => {
    if (text.length > 0) {
      const num = parseInt(text);
      if (num > 65535) {
        return;
      }
      if (num.toString() !== text) {
        return;
      }
    }
    setText(text);
  };

  const navigation = useNavigation();
  const queryClient = useQueryClient();
  const {mutate, isPending} = useMutation({
    mutationFn: async (peerPort: number) =>
      await save({peerPort} as AppSettingsFields),
    onSuccess: data => {
      queryClient.setQueryData(QueryKey, data);
      navigation.goBack();
    },
    onError: error => {
      Alert.alert(error.name, error.message);
    },
  });

  const handleSubmit = (text: string) => {
    if (text.length <= 0 || isPending) return;
    mutate(parseInt(text));
  };

  return (
    <Fragment>
      <View style={{paddingBottom: scale(sizes.base * 3)}}>
        <TextInput
          style={{
            width: '100%',
            borderColor: colors.divider,
            borderBottomWidth: 1,
          }}
          editable={!isPending}
          maxLength={text === '0' ? 1 : 5}
          value={text}
          placeholder={t('screen.peerPortEditor.placeholder')}
          onChangeText={handleChange}
          onSubmitEditing={() => handleSubmit(text)}
          inputMode="numeric"
        />
        <Text padding={[0, 0.5]} color="textSecondary" size="small">
          {t('screen.peerPortEditor.reminder')}
        </Text>
      </View>
      <Button
        title={t('action.save')}
        disabled={
          isPending || !text || text.length <= 0 || value?.toString() === text
        }
        onPress={() => handleSubmit(text)}
      />
    </Fragment>
  );
};
