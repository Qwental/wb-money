// Импорт сгенерированных файлов
const proto = require('./gen/money_service_pb.js');
const service = require('./gen/money_service_grpc_web_pb.js');

// Кладём в глобальный объект для доступа из app.js
window.proto = window.proto || {};
window.proto.money_service = {
    GetSavingsRequest: proto.GetSavingsRequest,
    GetSavingsResponse: proto.GetSavingsResponse,
    MoneyServiceClient: service.MoneyServiceClient,
};
