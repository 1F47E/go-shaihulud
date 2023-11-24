import { useEffect, useRef, useState } from 'react';
import { TextInput, TextInputProps, ActionIcon, Button, Paper, Text, Stack, Container, Textarea, Skeleton, SimpleGrid, Grid, useMantineTheme, rem } from '@mantine/core';
import { IconSearch, IconArrowRight } from '@tabler/icons-react';

type Message = {
    text: string;
    user: string;
};

const PRIMARY_COL_HEIGHT = rem(500);



export default function Chat() {
    const [messages, setMessages] = useState<Message[]>([]);
    const [newMessage, setNewMessage] = useState('');
    const messagesEndRef = useRef(null);

    const scrollToBottom = () => {
        const currentElement = messagesEndRef.current as unknown as HTMLElement | null;
        if (currentElement) {
            currentElement.scrollIntoView({ behavior: "smooth" });
        }
    };

    const [lastUser, setLastUser] = useState('Guest');

    const messageInputRef = useRef<HTMLTextAreaElement>(null);

    useEffect(() => {
        scrollToBottom();
        // Focus on message input after sending a message and on load
        // messageInputRef.current?.focus();
        focusTextarea()
    }, [messages]);

    const handleSend = () => {
        if (newMessage.trim()) {
            const nextUser = lastUser === 'You' ? 'Guest' : 'You';
            setLastUser(nextUser);
            setMessages([...messages, { text: newMessage, user: nextUser }]);
            setNewMessage('');
            // Focus on the textarea after sending a message
            // messageInputRef.current?.focus();
        focusTextarea()
        }
    };
    const focusTextarea = () => {
        const textarea = document.getElementById('chat-textarea') as HTMLTextAreaElement;
        textarea?.focus();
    };

    const handleInput = (value: string) => {
        setNewMessage(value)
    }

    const theme = useMantineTheme();

    return (
        <Container 
        style={{ 
            // border: '1px solid #f00',
            height: '100vh' 
            }}>
            <Grid>
                <Grid.Col>
                    <Stack 
                    justify="flex-end"
                    style={{ 
                        // border: '1px solid #00f',
                        height: '90vh'
                    }}>
                        <Paper style={{
                            overflowY: 'auto',
                            scrollbarWidth: 'none',
                            flexGrow: 1,
                            padding: 10
                        }}
                        >
                            {messages.map((message, index) => (
                                <Paper
                                    key={index}
                                    shadow="xs"
                                    radius="xl"
                                    p="md"
                                    style={{
                                        paddingLeft: '2rem',
                                        maxWidth: '70%',
                                        margin: '15px',
                                        color: message.user === 'You' ? theme.colors.blue[0] : theme.colors.gray[0],
                                        backgroundColor: message.user === 'You' ? theme.colors.blue[6] : theme.colors.gray[9],
                                        borderColor: message.user === 'You' ? theme.colors.blue[6] : theme.colors.gray[6],
                                    }}
                                >
                                    <Text 
                                    style={{ 
                                        marginBottom: '0px', 
                                        fontWeight: 500, 
                                        textAlign: 'left'
                                        }}>
                                        {message.user}
                                    </Text>
                                    <Text size="sm" 
                                    style={{ 
                                        textAlign: 'left'
                                        }}
                                        ><pre>{message.text}</pre></Text>
                                </Paper>
                            ))}
                            <div ref={messagesEndRef} />
                        </Paper>
                    </Stack>
                    <Stack 
                    justify="flex-end"
                    style={{ 
                        }}>
                    <Paper style={{ marginTop: 'md' }}>
                        <Textarea
                            id="chat-textarea"
                            placeholder="Type your message"
                            value={newMessage}
                            onChange={(event) => handleInput(event.currentTarget.value)}
                            onKeyDown={(event) => {
                                if (event.key === 'Enter' && event.shiftKey) {
                                    handleInput(event.currentTarget.value);
                                } else if (event.key === 'Enter') {
                                    handleSend();
                                }
                            }}
                            minRows={2}
                            autosize
                            radius="lg"
                            rightSectionWidth={90}
                            rightSection={
                                <Button
                                    radius="lg"
                                    color={theme.primaryColor}
                                    variant="filled"
                                    onClick={handleSend}
                                >
                                    Send
                                </Button>
                            }
                            autoFocus
                        />
                    </Paper>
                </Stack>
            </Grid.Col>
        </Grid>
        </Container >

    );

    // dissabledc 
    return (
        <Container style={{ display: 'flex', flexDirection: 'column', height: '100vh' }}>


            <Stack style={{ overflowY: 'scroll', flexGrow: 1 }}>
                {messages.map((message, index) => (
                    <Paper
                        key={index}
                        shadow="xs"
                        p="md"
                        style={{
                            maxWidth: '70%',
                            margin: '10px',
                            backgroundColor: message.user === 'You' ? '#blue' : '#grey',
                        }}
                    >
                        <Text style={{ marginBottom: '5px', fontWeight: 500 }}>
                            {message.user}
                        </Text>
                        <Text size="sm">{message.text}</Text>
                    </Paper>
                ))}
            </Stack>
            <Stack style={{ position: 'fixed', bottom: 0, left: 0, right: 0, padding: '10px' }}>
                <Textarea
                    placeholder="Type your message"
                    value={newMessage}
                    onChange={(event) => setNewMessage(event.currentTarget.value)}
                    style={{ marginBottom: '10px' }}
                />
                <Button onClick={handleSend}>Send</Button>
            </Stack>
        </Container>
    );
}
