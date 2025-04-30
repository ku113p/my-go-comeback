# CryptoPlatform Bot

A Telegram bot for monitoring cryptocurrency prices and managing notifications. This project is part of a Go learning series.

## 🚀 Features

- **Cryptocurrency Price Monitoring**: Fetches real-time cryptocurrency prices using the CoinMarketCap API.
- **User Notifications**: Allows users to set notifications for price thresholds.
- **Telegram Integration**: Interacts with users via Telegram commands and callback queries.
- **In-Memory Database**: Stores user data and notifications temporarily.

## 🛠 Setup

### Prerequisites

1. **Go**: Ensure you have Go installed on your system.
2. **Telegram Bot Token**: Obtain a bot token from [BotFather](https://core.telegram.org/bots#botfather).
3. **CoinMarketCap API Key**: Sign up at [CoinMarketCap](https://coinmarketcap.com/) to get an API key.

### Environment Variables

Set the following environment variables:

- `TG_API_TOKEN`: Your Telegram bot token.
- `CMC_API_KEY`: Your CoinMarketCap API key.

Example:

```bash
export TG_API_TOKEN="your-telegram-bot-token"
export CMC_API_KEY="your-coinmarketcap-api-key"
```

## 🏃‍♂️ Running the Bot

1. Navigate to the project directory:

   ```bash
   cd chapterts/02practice/03cryptoplatform
   ```

2. Run the bot:

   ```bash
   go run .
   ```

The bot will start in polling mode and interact with users via Telegram.

## 📖 Commands

### User Commands

- `/help`: Displays help information about the bot.
- `/add <symbol> <sign> <amount>`: Adds a notification for a cryptocurrency price.  
  Example: `/add BTC > 50000`
- `/list`: Lists all active notifications.

### Callback Queries

- **View Notification Details**: Click on a notification to view more details.
- **Delete Notification**: Confirm or cancel the deletion of a notification.

## 🔧 Project Structure

```
chapterts/02practice/03cryptoplatform/
├── app/                # Application-level utilities
├── collectors/         # Background tasks (e.g., price collection)
├── coinmarketcap/      # CoinMarketCap API integration
├── db/                 # In-memory database implementation
├── models/             # Data models
├── telegram/           # Telegram bot logic
│   ├── handlers/       # Command and callback query handlers
│   ├── middleware/     # Middleware for request handling
│   ├── services/       # Business logic services
│   └── view/           # Telegram UI components (keyboards, etc.)
├── utils/              # Utility functions (e.g., logging)
└── main.go             # Entry point
```

## 📚 Further Development

- Add persistent storage (e.g., PostgreSQL or MongoDB).
- Implement webhook mode for Telegram bot.
- Add unit tests for handlers and services.
- Enhance error handling and logging.

## 🧑‍💻 Contributing

Contributions are welcome! Feel free to submit issues or pull requests.

## 📜 License

This project is licensed under the MIT License. See the [LICENSE](../../LICENSE) file for details.