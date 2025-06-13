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

    if (!userId || isNaN(userId) || Number(userId) <= 0) {
        console.warn('[UI] Некорректный userId — показываем ошибку');
        showError('Пожалуйста, введите корректный ID пользователя (положительное число)');
        return;
    }

    hideElements();
    showLoading();

    try {
        console.log('[RPC] Отправка запроса в gRPC...');
        const data = await getSavings(userId);
        console.log('[RPC] Получен ответ:', data);
        hideLoading();

        // Проверяем статус ответа
        if (data.status === 'OK') {
            showResult(data);
        } else {
            showError(getErrorMessage(data.status, data.message));
        }
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
                console.error('[gRPC] Ошибка соединения:', err);
                reject(err);
                return;
            }

            // Получаем все поля из protobuf ответа с fallback-значениями
            const status = response.getStatus ? response.getStatus() : 0; // Default to OK
            const statusName = getStatusName(status);

            const result = {
                status: statusName,
                statusCode: status,
                totalSavings: response.getTotalSavings ? response.getTotalSavings() : 0,
                currency: response.getCurrency ? response.getCurrency() : 'RUB',
                totalPurchases: response.getTotalPurchases ? response.getTotalPurchases() : 0,
                wbCardPurchases: response.getWbCardPurchases ? response.getWbCardPurchases() : 0,
                message: response.getMessage ? response.getMessage() : ''
            };

            // Логируем, если ожидаемые поля отсутствуют
            if (!response.getStatus) console.warn('[gRPC] Поле status отсутствует в ответе');
            if (!response.getTotalPurchases) console.warn('[gRPC] Поле total_purchases отсутствует в ответе');
            if (!response.getWbCardPurchases) console.warn('[gRPC] Поле wb_card_purchases отсутствует в ответе');

            console.log('[gRPC] Полученные данные:', result);
            resolve(result);
        });
    });
}

// Преобразуем числовой статус в читаемое название
function getStatusName(statusCode) {
    const statusMap = {
        0: 'OK',
        1: 'USER_NOT_FOUND',
        2: 'NO_PURCHASES',
        3: 'DB_ERROR',
        4: 'INVALID_REQUEST',
        5: 'UNAUTHORIZED',
        6: 'UNKNOWN_ERROR'
    };
    return statusMap[statusCode] || 'UNKNOWN_STATUS';
}

// Возвращаем человекочитаемое сообщение об ошибке
function getErrorMessage(status, serverMessage) {
    const errorMessages = {
        'USER_NOT_FOUND': 'Ошибка:',
        'NO_PURCHASES': 'У пользователя нет покупок',
        'DB_ERROR': 'Ошибка базы данных',
        'INVALID_REQUEST': 'Некорректный запрос',
        'UNAUTHORIZED': 'Нет доступа к данным пользователя',
        'UNKNOWN_ERROR': 'Неизвестная ошибка'
    };

    const baseMessage = errorMessages[status] || 'Неизвестная ошибка';
    return serverMessage ? `${baseMessage}: ${serverMessage}` : baseMessage;
}

function showResult(data) {
    // Основная сумма экономии
    const savingsAmount = document.getElementById('savingsAmount');
    const savingsValue = data.totalSavings;

    // Определяем цвет и знак для отображения
    if (savingsValue > 0) {
        savingsAmount.textContent = `+${savingsValue.toLocaleString()} ${data.currency}`;
        savingsAmount.style.color = '#22c55e'; // зеленый для экономии
        savingsAmount.classList.add('positive');
        savingsAmount.classList.remove('negative');
    } else if (savingsValue < 0) {
        savingsAmount.textContent = `${savingsValue.toLocaleString()} ${data.currency}`;
        savingsAmount.style.color = '#ef4444'; // красный для переплат
        savingsAmount.classList.add('negative');
        savingsAmount.classList.remove('positive');
    } else {
        savingsAmount.textContent = `${savingsValue.toLocaleString()} ${data.currency}`;
        savingsAmount.style.color = '#6b7280'; // серый для нуля
        savingsAmount.classList.remove('positive', 'negative');
    }

    // Статистика покупок
    document.getElementById('totalOrders').textContent = data.totalPurchases.toLocaleString();
    document.getElementById('totalSpent').textContent = data.wbCardPurchases.toLocaleString();

    // Дополнительное сообщение от сервера (если есть)
    const messageElement = document.getElementById('serverMessage');
    if (data.message && messageElement) {
        messageElement.textContent = data.message;
        messageElement.style.display = 'block';
    } else if (messageElement) {
        messageElement.style.display = 'none';
    }

    // Показываем результат
    document.getElementById('result').style.display = 'block';
    document.getElementById('stats').style.display = 'grid';

    console.log('[UI] Результат отображён:', {
        savings: savingsValue,
        currency: data.currency,
        totalPurchases: data.totalPurchases,
        wbCardPurchases: data.wbCardPurchases
    });
}

function showError(message) {
    const errorElement = document.getElementById('error');
    errorElement.textContent = message;
    errorElement.style.display = 'block';
    console.log('[UI] Ошибка отображена:', message);
}

function showLoading() {
    document.getElementById('loading').style.display = 'block';
    console.log('[UI] Показан индикатор загрузки');
}

function hideLoading() {
    document.getElementById('loading').style.display = 'none';
    console.log('[UI] Скрыт индикатор загрузки');
}

function hideElements() {
    document.getElementById('result').style.display = 'none';
    document.getElementById('stats').style.display = 'none';
    document.getElementById('error').style.display = 'none';

    // Скрываем дополнительное сообщение если есть
    const messageElement = document.getElementById('serverMessage');
    if (messageElement) {
        messageElement.style.display = 'none';
    }

    console.log('[UI] Все элементы результата скрыты');
}