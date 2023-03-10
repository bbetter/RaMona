# RaMona (Rada Monitor)

![Ramon Gif](https://media.giphy.com/media/pCg4tODDp38ze/giphy.gif)

RaMona is a command line application which scrapes RSS data from the government law publishing website of "rada.gov.ua" for the latest law updates. The application filters the new laws based on words contained in their title or description. The list of these words can be specified with the triggers parameter followed by a list of space-separated words.

## Installation
To install the app follow this steps:

1. Clone the repository
```bash
git clone https://github.com/yourusername/your-repo.git
```
2. Navigate to directory
```bash
cd your-repo
```
3. Build the application
```bash
go build .
```

## Usage
To run RaMona, simply execute the following command:

```bash
go run main.go --triggers "keyword1 keyword2"
```

## Configuration

RaMona can be configured to send notifications to a Telegram chat. To enable this feature, set the *RANO_TELEGRAM_BOT_TOKEN* environment variables to the appropriate values for your Telegram bot and chat, respectively.

