import React, { useState } from 'react';
import Avatar from '@mui/material/Avatar';
import Button from '@mui/material/Button';
import CssBaseline from '@mui/material/CssBaseline';
import TextField from '@mui/material/TextField';
import FormControlLabel from '@mui/material/FormControlLabel';
import Checkbox from '@mui/material/Checkbox';
import Link from '@mui/material/Link';
import Grid from '@mui/material/Grid';
import Box from '@mui/material/Box';
import LockOutlinedIcon from '@mui/icons-material/LockOutlined';
import Typography from '@mui/material/Typography';
import Container from '@mui/material/Container';

export default function Register() {
    const [email, setEmail] = useState("");
    const [password, setPassword] = useState("");
    const [passwordConfirm, setPasswordConfirm] = useState("");
    const [login, setLogin] = useState("");
    const [message, setMessage] = useState("");

    const handleSubmit = (event) => {
        event.preventDefault();
        const data = new FormData(event.currentTarget);

        try {
            const response = fetch("http://localhost:8080/register", {
                method: "POST",
                headers: {
                    "Content-Type": "application/json",
                },
                body: JSON.stringify({ login, email, password, passwordConfirm }),
            });

            if (response.ok) {
                setMessage("Авторизация успешна!");
            } else {
                const text = response.text();
                setMessage(`Ошибка: ${text}`);
            }
        } catch (error) {
            setMessage(`Ошибка сети: ${error.message}`);
        }
    }

    return (
        <Container component="main" maxWidth="xs">
            <CssBaseline />
            <Box
                sx={{
                    marginTop: 8,
                    display: 'flex',
                    flexDirection: 'column',
                    alignItems: 'center',
                }}
            >
                <Avatar sx={{ m: 1, bgcolor: 'secondary.main' }}>
                    <LockOutlinedIcon />
                </Avatar>
                <Typography component="h1" variant="h5">
                    Регистрация
                </Typography>
                <Box component="form" onSubmit={handleSubmit} noValidate sx={{ mt: 1 }}>
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        onChange={e => setEmail(e.target.value)}
                        id="email"
                        label="Почтовый адрес"
                        name="email"
                        autoComplete="email"
                        autoFocus
                    />
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        onChange={e => setLogin(e.target.value)}
                        id="login"
                        label="Имя пользователя"
                        name="login"
                        autoComplete="login"
                        autoFocus
                    />
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        onChange={e => setPassword(e.target.value)}
                        name="password"
                        label="Пароль"
                        type="password"
                        id="password"
                        autoComplete="current-password"
                    />
                    <TextField
                        margin="normal"
                        required
                        fullWidth
                        onChange={e => setPasswordConfirm(e.target.value)}
                        name="passwordConfirm"
                        label="Подтверждение пароля"
                        type="password"
                        id="passwordConfirm"
                        autoComplete="current-password"
                    />
                    <FormControlLabel
                        control={<Checkbox value="remember" color="primary" />}
                        label="Запомнить меня"
                    />
                    <Button
                        type="submit"
                        fullWidth
                        variant="contained"
                        sx={{ mt: 3, mb: 2 }}
                    >
                        Зарегистрироваться
                    </Button>
                    <Link href="/auth" variant="body2">
                        {"Уже есть аккаунт?"}
                    </Link>
                </Box>
            </Box>
        </Container>
    );
}
