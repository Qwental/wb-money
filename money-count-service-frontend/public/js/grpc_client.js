const { MoneyServiceClient } = require('./proto/money_service_grpc_web_pb');
const { UserRequest } = require('./proto/money_service_pb');

// Создаем клиент
const client = new MoneyServiceClient('http://localhost:3000/grpc');

async function getSavings(userId) {
    const request = new UserRequest();
    request.setUserId(userId);

    return new Promise((resolve, reject) => {
        client.getSavings(request, {}, (err, response) => {
            if (err) {
                reject(err);
            } else {
                resolve({
                    potentialSavings: response.getPotentialSavings(),
                    totalOrders: response.getTotalOrders(),
                    totalSpent: response.getTotalSpent()
                });
            }
        });
    });
}

module.exports = { getSavings };