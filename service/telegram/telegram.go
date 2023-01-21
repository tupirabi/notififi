package telegram

import (
	"context"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

const (
	ModeMarkdown = tgbotapi.ModeMarkdown
	ModeHTML     = tgbotapi.ModeHTML
)

var parseMode = ModeHTML

type Telegram struct {
	client  *tgbotapi.BotAPI
	chatIDs []int64
}

func New(apiToken string) (*Telegram, error) {
	client, err := tgbotapi.NewBotAPI(apiToken)
	if err != nil {
		return nil, err
	}

	t := &Telegram{
		client:  client,
		chatIDs: []int64{},
	}

	return t, nil
}

func (t *Telegram) SetClient(client *tgbotapi.BotAPI) {
	t.client = client
}

func (t *Telegram) SetParseMode(mode string) {
	parseMode = mode
}

func (t *Telegram) AddReceivers(chatIDs ...int64) {
	t.chatIDs = append(t.chatIDs, chatIDs...)
}

func (t Telegram) Send(ctx context.Context, subject, message string) error {
	fullMessage := subject + "\n" + message // Treating subject as message title

	msg := tgbotapi.NewMessage(0, fullMessage)
	msg.ParseMode = parseMode

	for _, chatID := range t.chatIDs {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
			msg.ChatID = chatID
			_, err := t.client.Send(msg)
			if err != nil {
				return errors.Wrapf(err, "failed to send message to Telegram chat '%d'", chatID)
			}
		}
	}

	return nil
}
