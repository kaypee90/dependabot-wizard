import Button from '@mui/material/Button';

import { observer } from 'mobx-react-lite';
import { configStore } from '../store/configuration-store';

export const GithubInfoButton: React.FC = observer(() => {

  const handleOpen = () => {
    if (configStore.hasSavedRepoDetails) {
      configStore.setShowGithubInfo(true);
    } else {
      configStore.setShowGithubRepoDetails(true);
    }
  }

  return (
    <Button color="inherit" onClick={handleOpen}>Github Info</Button>
  );
});

