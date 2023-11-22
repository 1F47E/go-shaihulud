import { useState } from 'react';
import { Textarea, TextInput, Button, LoadingOverlay, Space} from '@mantine/core';
import { IconKey, IconLock } from '@tabler/icons-react';

export default function ChatJoin() {
  const [authKey, setAuthKey] = useState('');
  const [password, setPassword] = useState('');
  const [loading, setLoading] = useState(false);

  const handleJoin = () => {
    setLoading(true);
  };

  const handleCreate = () => {
    setLoading(true);
  }

  return (
    <div>
      <LoadingOverlay visible={loading} />
      <Textarea
        leftSection={<IconKey />}
        placeholder="AUTH KEY"
        value={authKey}
        onChange={(event) => setAuthKey(event.currentTarget.value)}
        size="md"
        radius="md"
        style={{ marginBottom: '1em' }}
        minRows={8}
      />
      <TextInput
        leftSection={<IconLock />}
        placeholder="PASSWORD"
        type="password"
        value={password}
        onChange={(event) => setPassword(event.currentTarget.value)}
        size="md"
        radius="md"
        style={{ marginBottom: '1em' }} 
      />
      <Button
        variant="filled"
        color="violet"
        onClick={handleJoin}
        disabled={loading}
        size="md"
        radius="xl"
        style={{ marginBottom: '4em', width: '100%' }} 
      >
        Join
      </Button>

      <Space h="md"/>

      <Button
        variant="filled"
        color="green"
        onClick={handleCreate}
        disabled={loading}
        size="sm"
        radius="xl"
      >
        Create new chat
      </Button>
    </div>
  );
}
