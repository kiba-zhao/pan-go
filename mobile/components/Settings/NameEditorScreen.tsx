import type {AppSettingsFields} from '@pango/data';
import {Fragment, useEffect, useState} from 'react';
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

const NameEditorScreen = () => {
  const {data, refetch} = useSettings();

  return (
    <I18NextProvider>
      <ScreenLayout
        refreshControl={<ScreenRefreshControl onRefresh={refetch} />}>
        <HeaderTitle i18nKey="screen.nameEditor.name" />
        <EditorSection value={data?.name} />
      </ScreenLayout>
    </I18NextProvider>
  );
};

export default NameEditorScreen;

type EditorSectionProps = {
  value?: string;
};
const EditorSection = ({value}: EditorSectionProps) => {
  const {sizes, colors} = useTheme();

  const {t} = useTranslation();
  const navigation = useNavigation();
  const [text, setText] = useState(value);

  useEffect(() => {
    setText(value);
  }, [value]);

  const queryClient = useQueryClient();
  const {mutate} = useMutation({
    mutationFn: async (name: string) => await save({name} as AppSettingsFields),
    onSuccess: data => {
      queryClient.setQueryData(QueryKey, data);
      navigation.goBack();
    },
    onError: error => {
      Alert.alert(error.name, error.message);
    },
  });

  const handleSubmit = (name?: string) => {
    if (!name) return;
    mutate(name);
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
          value={text}
          onChangeText={setText}
          onSubmitEditing={() => handleSubmit(text)}
          maxLength={50}
        />
        <Text padding={[0, 0.5]} color="textSecondary" size="small">
          {t('screen.nameEditor.helper')}
        </Text>
      </View>
      <Button
        title={t('action.save')}
        disabled={!text || text.length <= 0 || value === text}
        onPress={() => handleSubmit(text)}
      />
    </Fragment>
  );
};
