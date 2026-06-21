package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"net/http"
	"regexp"
	"strings"
	"time"
	"net/url"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"

	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"google.golang.org/protobuf/proto"

)

type AIMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type AISession struct {
	SenderID string
	Messages []AIMessage
	BotLID   string
}

var aiCache = make(map[string]AISession)

func uploadImageToTelegraph(client *whatsmeow.Client, v *events.Message) string {
	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg == nil || extMsg.ContextInfo == nil || extMsg.ContextInfo.QuotedMessage == nil {
		return ""
	}

	quoted := extMsg.ContextInfo.QuotedMessage
	if img := quoted.GetImageMessage(); img != nil {
		data, err := client.Download(context.Background(), img)
		if err != nil { return "" }

		body := &bytes.Buffer{}
		writer := multipart.NewWriter(body)
		part, _ := writer.CreateFormFile("file", "image.jpg")
		part.Write(data)
		writer.Close()

		req, _ := http.NewRequest("POST", "https://telegra.ph/upload", body)
		req.Header.Set("Content-Type", writer.FormDataContentType())

		httpClient := &http.Client{Timeout: 10 * time.Second}
		resp, err := httpClient.Do(req)
		if err != nil { return "" }
		defer resp.Body.Close()

		var tResp []struct { Src string `json:"src"` }
		json.NewDecoder(resp.Body).Decode(&tResp)
		if len(tResp) > 0 {
			return "https://telegra.ph" + tResp[0].Src
		}
	}
	return ""
}

func handleAICommand(client *whatsmeow.Client, v *events.Message, query string, cmd string) {
	if query == "" {
		replyMessage(client, v, "❌ *Error:* Please ask a question.\nExample: `.ai Explain BS CS concepts.`")
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "🧠")

	persona := fmt.Sprintf(`You are a silent ai made by %s.
RULES:
1. CASUAL CHAT: Be friendly and short.
2. HAPPY MODE: Send smile emojis and happy responses.
3. HINDI BLOCK: STRICTLY NO HINDI (Devanagari). Reply only in English or Roman Urdu.
4. SHORT ANSWER: Always keep answers short.
5. MEMORY: Always keep context in mind.`, GlobalConfig.Developer)

	imgUrl := uploadImageToTelegraph(client, v)

	session := AISession{
		SenderID: v.Info.Sender.User,
		BotLID:   getCleanID(client.Store.ID.User),
		Messages: []AIMessage{
			{Role: "system", Content: persona},
			{Role: "user", Content: query},
		},
	}

	if imgUrl != "" {
	    session.Messages[len(session.Messages)-1].Content += " [IMAGE_URL: " + imgUrl + "]"
	}

	go processAndSendAI(client, v, session)
}

func processAndSendAI(client *whatsmeow.Client, v *events.Message, session AISession) {
	react(client, v.Info.Chat, v.Info.ID, "⏳")

	imgUrl := ""
	lastMsg := session.Messages[len(session.Messages)-1].Content
	if strings.Contains(lastMsg, "[IMAGE_URL: ") {
		parts := strings.Split(lastMsg, "[IMAGE_URL: ")
		if len(parts) == 2 {
			session.Messages[len(session.Messages)-1].Content = strings.TrimSpace(parts[0])
			imgUrl = strings.TrimSuffix(parts[1], "]")
		}
	}

	var compiledPrompt strings.Builder
	for _, msg := range session.Messages {
		if msg.Role == "system" {
			compiledPrompt.WriteString(msg.Content + "\n\n")
		} else if msg.Role == "user" {
			compiledPrompt.WriteString("User: " + msg.Content + "\n")
		} else if msg.Role == "assistant" {
			compiledPrompt.WriteString("AI: " + msg.Content + "\n")
		}
	}

	var aiReplyText string
	if imgUrl != "" {
		apiURL := fmt.Sprintf("%s/ai/aiappchat?prompt=%s&image=%s", GlobalConfig.APIBaseURL, url.QueryEscape(compiledPrompt.String()), url.QueryEscape(imgUrl))
		resp, err := http.Get(apiURL)
		if err == nil {
			defer resp.Body.Close()
			var result map[string]interface{}
			if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
				if data, ok := result["data"].(string); ok {
					aiReplyText = data
				}
			}
		}
	} else {
		requestBody := map[string]string{
			"key":    "silent-ai",
			"prompt": compiledPrompt.String() + "AI:",
		}

		jsonData, _ := json.Marshal(requestBody)
		req, _ := http.NewRequest("POST", "https://silent-ai-pro-phi.vercel.app/api/ask", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")

		httpClient := &http.Client{Timeout: 90 * time.Second}
		resp, err := httpClient.Do(req)

		if err == nil {
			defer resp.Body.Close()
			reader := bufio.NewReader(resp.Body)
			var rawResponse strings.Builder
			for {
				line, err := reader.ReadString('\n')
				if err != nil { break }
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "data: ") {
					jsonStr := strings.TrimPrefix(line, "data: ")
					var dataChunk struct {
						Type string `json:"type"`
						Text string `json:"text"`
					}
					if err := json.Unmarshal([]byte(jsonStr), &dataChunk); err == nil && dataChunk.Type == "text" {
						rawResponse.WriteString(dataChunk.Text)
					}
				}
			}
			aiReplyText = strings.TrimSpace(rawResponse.String())
		}

		if aiReplyText == "" {
		    apiURL := fmt.Sprintf("%s/ai/chatup?prompt=%s", GlobalConfig.APIBaseURL, url.QueryEscape(compiledPrompt.String()))
		    resp2, err2 := http.Get(apiURL)
		    if err2 == nil {
		        defer resp2.Body.Close()
		        var result map[string]interface{}
		        if err := json.NewDecoder(resp2.Body).Decode(&result); err == nil {
		            if data, ok := result["data"].(string); ok {
		                aiReplyText = data
		            }
		        }
		    }
		}
	}

	aiReplyText = strings.ReplaceAll(aiReplyText, "**", "*")
	reHeaders := regexp.MustCompile("(?m)^#{1,6}\\s+(.*)$")
	aiReplyText = reHeaders.ReplaceAllString(aiReplyText, "*$1*")

	if aiReplyText != "" {
		msgID := replyMessage(client, v, aiReplyText)
		session.Messages = append(session.Messages, AIMessage{Role: "assistant", Content: aiReplyText})

		if msgID != "" {
			aiCache[msgID] = session
			go func(id string) {
				time.Sleep(1 * time.Hour)
				delete(aiCache, id)
			}(msgID)
		}
		react(client, v.Info.Chat, v.Info.ID, "✅")
	} else {
		replyMessage(client, v, "❌ API Offline or Error occurred.")
		react(client, v.Info.Chat, v.Info.ID, "❌")
	}
}

func HandleAIChatReply(client *whatsmeow.Client, v *events.Message, bodyClean string, qID string) bool {
	if session, ok := aiCache[qID]; ok {
	    senderClean := strings.Split(v.Info.Sender.User, "@")[0]
	    sessionSenderClean := strings.Split(session.SenderID, "@")[0]

		if senderClean == sessionSenderClean {
			session.Messages = append(session.Messages, AIMessage{Role: "user", Content: bodyClean})
		
			if len(session.Messages) > 15 {
				session.Messages = append([]AIMessage{session.Messages[0]}, session.Messages[len(session.Messages)-14:]...)
			}
			
			go processAndSendAI(client, v, session)
			return true
		}
	}
	return false
}

func getCleanID(jidStr string) string {
	if jidStr == "" { return "unknown" }
	parts := strings.Split(jidStr, "@")
	if len(parts) == 0 { return "unknown" }
	return strings.TrimSpace(parts[0])
}

func handleJoke(client *whatsmeow.Client, v *events.Message, args string) {
    targetJID := v.Info.Sender.User

	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg != nil && extMsg.ContextInfo != nil && extMsg.ContextInfo.Participant != nil {
		targetJID = *extMsg.ContextInfo.Participant
	} else if extMsg != nil && extMsg.ContextInfo != nil && len(extMsg.ContextInfo.MentionedJID) > 0 {
		targetJID = extMsg.ContextInfo.MentionedJID[0]
	}

    react(client, v.Info.Chat, v.Info.ID, "🤣")

    cleanTarget := strings.Split(targetJID, "@")[0]

    prompt := fmt.Sprintf("Make a very funny, savage but friendly short joke in Roman Urdu about the user @%s. The joke should be extremely funny and easy to read for a layman.", cleanTarget)

    apiURL := fmt.Sprintf("%s/ai/chatup?prompt=%s", GlobalConfig.APIBaseURL, url.QueryEscape(prompt))
    resp, err := http.Get(apiURL)
    if err == nil {
        defer resp.Body.Close()
        var result map[string]interface{}
        if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
            if data, ok := result["data"].(string); ok {
                replyMsg := &waProto.Message{
					ExtendedTextMessage: &waProto.ExtendedTextMessage{
						Text: proto.String(data),
						ContextInfo: &waProto.ContextInfo{
							StanzaID:      proto.String(v.Info.ID),
							Participant:   proto.String(v.Info.Sender.String()),
							QuotedMessage: v.Message,
							MentionedJID:  []string{targetJID},
						},
					},
				}
				client.SendMessage(context.Background(), v.Info.Chat, replyMsg)
                react(client, v.Info.Chat, v.Info.ID, "✅")
                return
            }
        }
    }

    replyMessage(client, v, "❌ Could not generate joke right now.")
}
