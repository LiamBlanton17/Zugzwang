$(() => {

    // Handle starting game
    $('#start-game')?.click(async (e) => {
        e.preventDefault();

        // Verify the inputs before starting the game
        const nameEl = $('#name');
        const eloEl = $('#elo');
        const name = String(nameEl?.val()) || '';
        const elo = Number(eloEl?.val()) || 0;

        // Require name to be between 3 and 50 characters
        if (!name) {
            alert('Please enter your name');
            return;
        }
        if (name.length > 50 || name.length < 3) {
            alert('Your name must be between 3 and 50 characters (sorry if it\'s not)');
            return;
        }

        // Require ELO to be between 100 and 3200
        if (!elo) {
            alert('Please enter your elo (or just an estimate)')
            return;
        }
        if (elo < 100 || elo > 3200) {
            alert('Your elo must be between 100 and 3000 (you\'re not that good)')
            return;
        }

        // Save ELO and name to local storage
        localStorage.setItem('name', name);
        localStorage.setItem('elo', elo);

        
        // Make a request to the backend to setup the game
        let gameId;
        try {
            const response = await fetch('/setup', {
                method: "POST",
                body: JSON.stringify({name, elo}),
            });

            // Check resposne ok
            if (!response.ok) {
                throw new Error(`Response threw error: ${response.statusText}`);
            }

            // Get response json
            const json = await response.json();

            // Get game id
            gameId = json['gameId'];

            // Check game id
            if (!gameId) {
                throw new Error('Response did not provide a game id');
            }

        } catch (e) {
            alert('Something went wrong starting the game. Please try again.');
            console.log(e);
        }

        // Start the game via a websocket
        startGameWS(gameId);
    });

    // Function to start the game via a websocket
    async function startGameWS(gameId) {

        // Define the websocket connection
        const ws = new WebSocket(`ws://localhost/start/${gameId}`);

        // Define the open handling
        ws.addEventListener('open', (e) => {

        });

        // Define the message handling
        ws.addEventListener('message', () => {

        });

        // Define the error handling
        ws.addEventListener('error', () => {
            
        });

        // Define the close handling
        ws.addEventListener('close', () => {

        });

    }

});
