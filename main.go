package main

import (
	"io"
	"log"

	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client" // Import the client subpackage
	"github.com/emersion/go-message/mail"
)

func main() {
	// 1. Подключение к IMAP серверу Gmail через SSL
	// Use imap.gmail.com:993 for TLS connection to Gmail
	c, err := client.DialTLS("imap.gmail.com:993", nil)
	if err != nil {
		log.Fatalf("Ошибка подключения: %v", err)
	}
	defer c.Logout()

	// 2. Авторизация (вместо "app_password" вставьте 16-значный пароль приложения)
	if err := c.Login("xadmy2771@gmail.com", "dmwq xrrd ykiw tgui"); err != nil {
		log.Fatalf("Ошибка авторизации: %v", err)
	}

	// 3. Выбор папки "Входящие"
	mbox, err := c.Select("INBOX", false)
	if err != nil {
		log.Fatalf("Ошибка выбора папки: %v", err)
	}

	if mbox.Messages == 0 {
		log.Println("В папке нет писем.")
		return
	}

	// 4. Запрос последнего письма
	seqset := new(imap.SeqSet)
	seqset.AddNum(mbox.Messages) // Берем индекс последнего сообщения

	// Запрашиваем структуру и тело письма
	var section imap.BodySectionName
	items := []imap.FetchItem{imap.FetchEnvelope, section.FetchItem()}
	messages := make(chan *imap.Message, 1)

	go func() {
		if err := c.Fetch(seqset, items, messages); err != nil {
			log.Fatalf("Ошибка получения письма: %v", err)
		}
	}()

	msg := <-messages
	if msg == nil {
		log.Fatal("Сервер не вернул сообщение.")
	}

	// 5. Парсинг заголовков (Тема и Отправитель)
	log.Printf("От: %s", msg.Envelope.From[0].Address())
	log.Printf("Тема: %s", msg.Envelope.Subject)

	// 6. Парсинг текста письма
	r := msg.GetBody(&section)
	if r == nil {
		log.Fatal("Тело письма пустое.")
	}

	mr, err := mail.CreateReader(r)
	if err != nil {
		log.Fatalf("Ошибка создания reader: %v", err)
	}

	for {
		p, err := mr.NextPart()
		if err == io.EOF {
			break
		} else if err != nil {
			log.Fatal(err)
		}

		switch h := p.Header.(type) {
		case *mail.InlineHeader:
			contentType, _, _ := h.ContentType()
			// Читаем только обычный текст (text/plain)
			if contentType == "text/plain" {
				b, _ := io.ReadAll(p.Body)
				log.Printf("Текст письма:\n%s", string(b))
			}
		}
	}
}
