function sendClientErrorLog(message, source, lineno, colno, error) {
    fetch('/log-client-error', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            type: 'error',
            message,
            source,
            lineno,
            colno,
            stack: error?.stack || null,
            userAgent: navigator.userAgent,
            timestamp: new Date().toISOString()
        })
    });
}

window.onerror = function (message, source, lineno, colno, error) {
    sendClientErrorLog(message, source, lineno, colno, error);
};

window.onunhandledrejection = function (event) {
    fetch('/log-client-error', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json'
        },
        body: JSON.stringify({
            type: 'unhandledrejection',
            reason: event.reason,
            stack: event.reason?.stack || null,
            userAgent: navigator.userAgent,
            timestamp: new Date().toISOString()
        })
    });
};
