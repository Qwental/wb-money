const express = require('express');
const path = require('path');
const {createProxyMiddleware} = require('http-proxy-middleware');

const app = express();
const PORT = 3000;

// Для парсинга JSON в POST-запросах
app.use(express.json());

//  Раздача статики
app.use(express.static(path.join(__dirname, 'public')));

// CORS для gRPC-Web
app.use((req, res, next) => {
    res.header('Access-Control-Allow-Origin', '*');
    res.header('Access-Control-Allow-Methods', 'GET, POST, PUT, DELETE, OPTIONS');
    res.header('Access-Control-Allow-Headers', 'Origin, X-Requested-With, Content-Type, Accept, Authorization, X-Grpc-Web, X-User-Agent');
    res.header('Access-Control-Expose-Headers', 'Grpc-Status, Grpc-Message, Grpc-Encoding, Grpc-Accept-Encoding');

    if (req.method === 'OPTIONS') {
        res.sendStatus(200);
    } else {
        next();
    }
});

app.use('/grpc', createProxyMiddleware({
    target: 'http://money-service:8080',  // Используем имя сервиса и порт gRPC-Web
    changeOrigin: true,
    pathRewrite: {'^/grpc': ''},
    onProxyReq: (proxyReq, req) => {
        console.log(`[PROXY] ${req.method} ${req.url} → http://money-service:8080${req.url.replace('/grpc', '')}`);
        if (req.headers['content-type']) {
            proxyReq.setHeader('Content-Type', req.headers['content-type']);
        }
        if (req.headers['x-grpc-web']) {
            proxyReq.setHeader('X-Grpc-Web', req.headers['x-grpc-web']);
        }
        if (req.headers['x-user-agent']) {
            proxyReq.setHeader('X-User-Agent', req.headers['x-user-agent']);
        }
    },
    onProxyRes: (proxyRes, req, res) => {
        console.log(`[PROXY] Response: ${proxyRes.statusCode} for ${req.method} ${req.url}`);
    },
    onError: (err, req, res) => {
        console.error(`[PROXY] Error for ${req.method} ${req.url}:`, err.message);
        res.status(502).json({error: 'Proxy error', details: err.message});
    }
}));

// API заглушка
app.get('/api/savings/:userId', (req, res) => {
    res.status(501).json({error: "Use gRPC endpoint instead"});
});

// Приём JS-ошибок от клиента
app.post('/log-client-error', (req, res) => {
    const {type, message, source, lineno, colno, stack, reason, userAgent, timestamp} = req.body;

    console.error(`\nClient JS ${type === 'unhandledrejection' ? 'Unhandled Promise' : 'Error'} at ${timestamp}`);
    if (type === 'unhandledrejection') {
        console.error(`Reason: ${JSON.stringify(reason)}`);
    } else {
        console.error(`Message: ${message}`);
        console.error(`Source: ${source}:${lineno}:${colno}`);
        if (stack) console.error(`Stack:\n${stack}`);
    }
    console.error(`UserAgent: ${userAgent}`);
    res.status(204).end();
});

// 🔹 Запуск сервера
const server = app.listen(PORT, () => {
    console.log(`Server running at http://localhost:${PORT}`);
    console.log(`Proxying /grpc/* → http://localhost:8080`);
});

//Graceful shutdown & error catch
function shutdown() {
    console.log('\n🧹 Gracefully shutting down...');
    server.close(() => {
        console.log('Server closed');
        process.exit(0);
    });
    setTimeout(() => {
        console.error('Force shutdown (timeout)');
        process.exit(1);
    }, 5000);
}

process.on('SIGINT', shutdown);
process.on('SIGTERM', shutdown);

// Для nodemon
process.once('SIGUSR2', () => {
    shutdown();
    process.kill(process.pid, 'SIGUSR2');
});

// Глобальные ошибки
process.on('uncaughtException', (err) => {
    console.error('Uncaught Exception:', err);
    shutdown();
});

process.on('unhandledRejection', (reason) => {
    console.error('Unhandled Rejection:', reason);
    shutdown();
});