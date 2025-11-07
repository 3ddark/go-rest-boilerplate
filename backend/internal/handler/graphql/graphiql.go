package graphql

const graphiqlHTML = `
<!DOCTYPE html>
<html>
<head>
    <meta charset="UTF-8">
    <title>GraphiQL</title>
    <link href="https://unpkg.com/graphiql@3.0.6/graphiql.min.css" rel="stylesheet" />
    <style>
        body { margin: 0; font-family: system-ui, -apple-system, sans-serif; }
        #auth-toolbar {
            height: 50px;
            background-color: #1b1b1b;
            border-bottom: 1px solid #333;
            display: flex;
            align-items: center;
            padding: 0 20px;
            color: white;
            box-sizing: border-box;
        }
        #auth-toolbar .brand { font-weight: bold; margin-right: 20px; color: #e535ab; }
        #auth-toolbar input {
            flex-grow: 1;
            max-width: 500px;
            padding: 8px 12px;
            border-radius: 4px;
            border: 1px solid #444;
            background: #333;
            color: #fff;
            margin-right: 10px;
        }
        #auth-toolbar button {
            padding: 8px 16px;
            border: none;
            border-radius: 4px;
            cursor: pointer;
            font-weight: 600;
            transition: background 0.2s;
        }
        #btn-save { background: #009688; color: white; margin-right: 10px; }
        #btn-save:hover { background: #00796b; }
        #btn-clear { background: #f44336; color: white; }
        #btn-clear:hover { background: #d32f2f; }
        #graphiql { height: calc(100vh - 50px); }
    </style>
</head>
<body>
    <div id="auth-toolbar">
        <span class="brand">GraphQL API</span>
        <input type="text" id="jwt-token" placeholder="Bearer Token'ınızı buraya yapıştırın..." />
        <button id="btn-save" onclick="saveToken()">Authorize</button>
        <button id="btn-clear" onclick="clearToken()">Logout</button>
        <span id="status-msg" style="margin-left: 15px; color: #8bc34a; opacity: 0; transition: opacity 0.5s;">✓ Kaydedildi</span>
    </div>

    <div id="graphiql"></div>

    <script crossorigin src="https://unpkg.com/react@18.2.0/umd/react.production.min.js"></script>
    <script crossorigin src="https://unpkg.com/react-dom@18.2.0/umd/react-dom.production.min.js"></script>
    <script crossorigin src="https://unpkg.com/graphiql@3.0.6/graphiql.min.js"></script>

    <script>
        const tokenInput = document.getElementById('jwt-token');
        const statusMsg = document.getElementById('status-msg');

        window.addEventListener('load', () => {
            const currentToken = localStorage.getItem('jwt_token');
            if (currentToken) {
                tokenInput.value = currentToken;
            }
        });

        function showStatus(msg) {
            statusMsg.textContent = msg;
            statusMsg.style.opacity = 1;
            setTimeout(() => { statusMsg.style.opacity = 0; }, 2000);
        }

        function saveToken() {
            const token = tokenInput.value.trim();
            if (!token) {
                alert("Lütfen geçerli bir token giriniz.");
                return;
            }
            const cleanToken = token.replace('Bearer ', '');
            localStorage.setItem('jwt_token', cleanToken);
            tokenInput.value = cleanToken;
            showStatus("✓ Token Kaydedildi");
        }

        function clearToken() {
            localStorage.removeItem('jwt_token');
            tokenInput.value = '';
            showStatus("✓ Çıkış Yapıldı");
        }

        const root = ReactDOM.createRoot(document.getElementById('graphiql'));
        
        const customFetcher = async (graphQLParams) => {
            const token = localStorage.getItem('jwt_token');
            const headers = { 'Content-Type': 'application/json' };

            if (token) {
                headers['Authorization'] = 'Bearer ' + token;
            }

            try {
                const response = await fetch('/graphql', {
                    method: 'POST',
                    headers: headers,
                    body: JSON.stringify(graphQLParams),
                });
                return await response.json();
            } catch (error) {
                return { errors: [{ message: "Network Error: " + error.message }] };
            }
        };

        root.render(
            React.createElement(GraphiQL, { 
                fetcher: customFetcher,
                defaultEditorToolsVisibility: true,
            }),
        );
    </script>
</body>
</html>
`
