const form = document.getElementById('registerForm');
const password = document.getElementById('password');
const passwordConfirm = document.getElementById('passwordConfirm');
const passwordError = document.getElementById('passwordError');

form.addEventListener('submit', function(event) {
    if (password.value !== passwordConfirm.value) {
        passwordError.textContent = "Пароли не совпадают!";
        event.preventDefault(); // Остановить отправку формы
    } else {
        passwordError.textContent = ""; // Очистить сообщение, если всё ок
    }
});

// Можно добавить динамическую очистку ошибки при изменении любого из полей пароля:
password.addEventListener('input', function() {
    if (password.value === passwordConfirm.value) {
        passwordError.textContent = "";
    }
});
passwordConfirm.addEventListener('input', function() {
    if (password.value === passwordConfirm.value) {
        passwordError.textContent = "";
    }
});

