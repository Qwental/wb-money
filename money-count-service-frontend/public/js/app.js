// 💡 Создаём gRPC клиент
const client = new proto.money_service.MoneyServiceClient("http://localhost:3000/grpc");

// 🎯 Основная логика
document.addEventListener('DOMContentLoaded', () => {
    const calculateBtn = document.getElementById('calculateBtn');
    const userIdInput = document.getElementById('userId');

    if (!calculateBtn || !userIdInput) {
        console.error('❌ Элементы кнопки или input не найдены!');
        return;
    }

    calculateBtn.addEventListener('click', calculateSavings);
    userIdInput.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') calculateSavings();
    });

    console.log('✅ Приложение инициализировано');
});

async function calculateSavings() {
    const userId = document.getElementById('userId').value;
    console.log('[UI] Кнопка "Рассчитать экономию" нажата');
    console.log(`[UI] Введённый userId: "${userId}"`);

    if (!userId) {
        console.warn('[UI] userId пустой — показываем ошибку');
        showError('Пожалуйста, введите ваш ID пользователя');
        return;
    }

    hideElements();
    showLoading();

    try {
        console.log('[RPC] Отправка запроса в gRPC...');
        const data = await getSavings(userId);
        console.log('[RPC] Получен ответ:', data);
        hideLoading();
        showResult(data);
    } catch (error) {
        console.error('[RPC] Ошибка при запросе:', error);
        hideLoading();
        showError(error.message || 'Ошибка при получении данных');
    }
}

function getSavings(userId) {
    const request = new proto.money_service.GetSavingsRequest();
    request.setUserId(Number(userId));

    return new Promise((resolve, reject) => {
        client.getSavings(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else {
                resolve({
                    totalSavings: response.getTotalSavings(),
                    currency: response.getCurrency(),
                });
            }
        });
    });
}
function showResult(data) {
    document.getElementById('savingsAmount').textContent = `${data.totalSavings.toLocaleString()} ₽`;
    document.getElementById('totalOrders').textContent = '-';
    document.getElementById('totalSpent').textContent = '-';

    document.getElementById('result').style.display = 'block';
    document.getElementById('stats').style.display = 'grid';
}

function showError(message) {
    const errorElement = document.getElementById('error');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
}

function showLoading() {
    document.getElementById('loading').style.display = 'block';
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
}

function hideElements() {
    document.getElementById('result').style.display = 'none';
    document.getElementById('stats').style.display = 'none';
    document.getElementById('error').style.display = 'none';
}
