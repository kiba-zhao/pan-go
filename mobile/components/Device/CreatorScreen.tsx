import {useRoute} from '../App/Navigation';
import {ScreenLayout} from '../Common/ScreenBase';
import {default as Section} from '../Common/Section';
import {CreatorParam} from './ScreenRoute';

type DeviceScreenState = {};

const DeviceCreatorScreen = () => {
  const route = useRoute();
  const {name, peerId} = route.params as CreatorParam;

  return (
    <ScreenLayout>
      <DeviceFieldsSection />
    </ScreenLayout>
  );
};

export default DeviceCreatorScreen;

const DeviceFieldsSection = () => {
  return <Section> </Section>;
};
