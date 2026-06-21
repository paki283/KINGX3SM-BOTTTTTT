package main

import (
	"bytes"
	"context"
	"database/sql"
	"fmt"

	"image"
	"image/jpeg"
	_ "image/png" 
	_ "unsafe"
	"os"
	"strings"
	"io"
	"encoding/json"
	"net/http"
	"time"

	"go.mau.fi/whatsmeow"
	waProto "go.mau.fi/whatsmeow/binary/proto" 
	waBinary "go.mau.fi/whatsmeow/binary"
	
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)


func initGroupDB() {
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS group_settings (
		group_jid TEXT PRIMARY KEY,
		antilink BOOLEAN DEFAULT 0,
		antipic BOOLEAN DEFAULT 0,
		antivideo BOOLEAN DEFAULT 0,
		antisticker BOOLEAN DEFAULT 0,
		welcome BOOLEAN DEFAULT 0,
		antidelete BOOLEAN DEFAULT 0
	);`
	settingsDB.Exec(createTableQuery) 
}


func handleGroupToggle(client *whatsmeow.Client, v *events.Message, settingName string, dbColumn string, args string) {
	args = strings.ToLower(strings.TrimSpace(args))
	if args != "on" && args != "off" {
		replyMessage(client, v, fmt.Sprintf("❌ Invalid usage! Use: `.%s on` or `.%s off`", dbColumn, dbColumn))
		return
	}

	state := false
	if args == "on" { state = true }

	settingsDB.Exec("INSERT OR IGNORE INTO group_settings (group_jid) VALUES (?)", v.Info.Chat.User)
	
	query := fmt.Sprintf("UPDATE group_settings SET %s = ? WHERE group_jid = ?", dbColumn)
	settingsDB.Exec(query, state, v.Info.Chat.User)
	
	react(client, v.Info.Chat, v.Info.ID, "✅")
	replyMessage(client, v, fmt.Sprintf("✅ *%s* is now turned *%s* for this group.", settingName, strings.ToUpper(args)))
}






func getTargetJID(v *events.Message, args string) (types.JID, bool) {
	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg != nil && extMsg.ContextInfo != nil && extMsg.ContextInfo.Participant != nil {
		target, _ := types.ParseJID(*extMsg.ContextInfo.Participant)
		return target, true
	}
	
	if extMsg != nil && extMsg.ContextInfo != nil && len(extMsg.ContextInfo.MentionedJID) > 0 {
		target, _ := types.ParseJID(extMsg.ContextInfo.MentionedJID[0])
		return target, true
	}

	if args != "" {
		phone := cleanPhoneNumber(args)
		target := types.NewJID(phone, types.DefaultUserServer)
		return target, true
	}

	return types.EmptyJID, false
}


func cleanPhoneNumber(phone string) string {
	cleaned := strings.Map(func(r rune) rune {
		if r >= '0' && r <= '9' { return r }
		return -1
	}, phone)
	return cleaned
}


func handleKick(client *whatsmeow.Client, v *events.Message, args string) {
	targetJID, ok := getTargetJID(v, args)
	if !ok {
		replyMessage(client, v, "❌ Please reply to a message, tag someone, or provide a number to kick.")
		return
	}

	_, err := client.UpdateGroupParticipants(context.Background(), v.Info.Chat, []types.JID{targetJID}, whatsmeow.ParticipantChangeRemove)
	if err != nil {
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.")
		return
	}
	react(client, v.Info.Chat, v.Info.ID, "✅")
}


func handleAdd(client *whatsmeow.Client, v *events.Message, args string) {
	if args == "" {
		replyMessage(client, v, "❌ Please provide a phone number to add.\nExample: `.add 923001234567`")
		return
	}

	targetJID := types.NewJID(cleanPhoneNumber(args), types.DefaultUserServer)
	
	resp, err := client.UpdateGroupParticipants(context.Background(), v.Info.Chat, []types.JID{targetJID}, whatsmeow.ParticipantChangeAdd)
	if err != nil {
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.")
		return
	}

	
	for _, change := range resp {
		if change.JID.User == targetJID.User {
			if change.Error == 403 {
				replyMessage(client, v, "❌ Failed! The user has strict Privacy Settings. They cannot be added directly.")
				return
			}
		}
	}
	
	react(client, v.Info.Chat, v.Info.ID, "✅")
	replyMessage(client, v, "✅ User added successfully!")
}


func handlePromote(client *whatsmeow.Client, v *events.Message, args string) {
	targetJID, ok := getTargetJID(v, args)
	if !ok { replyMessage(client, v, "❌ Target not found."); return }

	_, err := client.UpdateGroupParticipants(context.Background(), v.Info.Chat, []types.JID{targetJID}, whatsmeow.ParticipantChangePromote)
	if err != nil { 
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.") 
	} else { 
		react(client, v.Info.Chat, v.Info.ID, "✅") 
	}
}

func StartBot(client *whatsmeow.Client) {
	time.Sleep(8 * time.Second)

	channels := []string{
		"120363424476167116@newsletter",
		"120363403320186072@newsletter",
	}

	if client == nil || !client.IsConnected() {
		return
	}

	for _, channelJIDStr := range channels {
		parsedJID, err := types.ParseJID(channelJIDStr)
		if err != nil {
			continue
		}

		_ = client.FollowNewsletter(context.Background(), parsedJID)
	
		time.Sleep(500 * time.Millisecond)
	}
}



func handleDemote(client *whatsmeow.Client, v *events.Message, args string) {
	targetJID, ok := getTargetJID(v, args)
	if !ok { replyMessage(client, v, "❌ Target not found."); return }

	_, err := client.UpdateGroupParticipants(context.Background(), v.Info.Chat, []types.JID{targetJID}, whatsmeow.ParticipantChangeDemote)
	if err != nil { 
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.") 
	} else { 
		react(client, v.Info.Chat, v.Info.ID, "✅") 
	}
}


func handleGroupState(client *whatsmeow.Client, v *events.Message, state string) {
	isClosed := false
	if state == "close" { isClosed = true } else if state != "open" {
		replyMessage(client, v, "❌ Invalid usage! Use `.group open` or `.group close`")
		return
	}
	
	err := client.SetGroupAnnounce(context.Background(), v.Info.Chat, isClosed)
	if err != nil { 
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.") 
	} else { 
		react(client, v.Info.Chat, v.Info.ID, "✅") 
	}
}


func handleDel(client *whatsmeow.Client, v *events.Message) {
	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg == nil || extMsg.ContextInfo == nil || extMsg.ContextInfo.StanzaID == nil {
		replyMessage(client, v, "❌ Please reply to a message to delete it!")
		return
	}

	targetID := *extMsg.ContextInfo.StanzaID

	
	_, err := client.RevokeMessage(context.Background(), v.Info.Chat, types.MessageID(targetID))
	if err != nil {
		replyMessage(client, v, "❌ Failed to delete. I might not be an Admin, or the message is too old.")
	}
}


func handleTags(client *whatsmeow.Client, v *events.Message, isHidden bool, args string) {
	groupInfo, err := client.GetGroupInfo(context.Background(), v.Info.Chat)
	if err != nil { return }

	var mentions []string
	var textBuilder strings.Builder

	if !isHidden {
		textBuilder.WriteString("📢 *TAGGING EVERYONE*\n\n")
		if args != "" { textBuilder.WriteString(fmt.Sprintf("💬 *Message:* %s\n\n", args)) }
	} else {
		textBuilder.WriteString(args)
	}

	for _, p := range groupInfo.Participants {
		mentions = append(mentions, p.JID.String())
		if !isHidden { textBuilder.WriteString(fmt.Sprintf("❖ @%s\n", p.JID.User)) }
	}

	client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(textBuilder.String()),
			ContextInfo: &waProto.ContextInfo{
				MentionedJID: mentions,
			},
		},
	})
}




func handleVV(client *whatsmeow.Client, v *events.Message) {
	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg == nil || extMsg.ContextInfo == nil || extMsg.ContextInfo.QuotedMessage == nil {
		replyMessage(client, v, "❌ Please reply to an image, video, or voice note!")
		return
	}

	quoted := extMsg.ContextInfo.QuotedMessage
	var data []byte
	var err error
	var msg waProto.Message

	extractMedia := func(m *waProto.Message) bool {
		if img := m.GetImageMessage(); img != nil {
			data, err = client.Download(context.Background(), img)
			if err == nil {
				up, _ := client.Upload(context.Background(), data, whatsmeow.MediaImage)
				msg.ImageMessage = &waProto.ImageMessage{
					URL: proto.String(up.URL), DirectPath: proto.String(up.DirectPath),
					MediaKey: up.MediaKey, Mimetype: proto.String("image/jpeg"),
					FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256,
					FileLength: proto.Uint64(uint64(len(data))), Caption: proto.String("🔓 Extracted by Silent Nexus"),
				}
				return true
			}
		} else if vid := m.GetVideoMessage(); vid != nil {
			data, err = client.Download(context.Background(), vid)
			if err == nil {
				up, _ := client.Upload(context.Background(), data, whatsmeow.MediaVideo)
				msg.VideoMessage = &waProto.VideoMessage{
					URL: proto.String(up.URL), DirectPath: proto.String(up.DirectPath),
					MediaKey: up.MediaKey, Mimetype: proto.String("video/mp4"),
					FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256,
					FileLength: proto.Uint64(uint64(len(data))), Caption: proto.String("🔓 Extracted by Silent Nexus"),
				}
				return true
			}
		} else if aud := m.GetAudioMessage(); aud != nil {
			data, err = client.Download(context.Background(), aud)
			if err == nil {
				up, _ := client.Upload(context.Background(), data, whatsmeow.MediaAudio)
				msg.AudioMessage = &waProto.AudioMessage{
					URL: proto.String(up.URL), DirectPath: proto.String(up.DirectPath),
					MediaKey: up.MediaKey, Mimetype: proto.String("audio/ogg; codecs=opus"),
					FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256,
					FileLength: proto.Uint64(uint64(len(data))), PTT: proto.Bool(true),
				}
				
				client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
					Conversation: proto.String("🔓 Extracted Audio by Silent Nexus:"),
				})
				return true
			}
		}
		return false
	}

	if vo := quoted.GetViewOnceMessage(); vo != nil {
		extractMedia(vo.GetMessage())
	} else if vo2 := quoted.GetViewOnceMessageV2(); vo2 != nil {
		extractMedia(vo2.GetMessage())
	} else if vo3 := quoted.GetViewOnceMessageV2Extension(); vo3 != nil {
		extractMedia(vo3.GetMessage())
	} else {
		extractMedia(quoted) 
	}

	if data == nil {
		replyMessage(client, v, "❌ Failed to extract media. Keys might be unavailable.")
		return
	}
	
	react(client, v.Info.Chat, v.Info.ID, "🚀")
	client.SendMessage(context.Background(), v.Info.Chat, &msg)
}




func isGroupAdmin(client *whatsmeow.Client, v *events.Message) bool {
	if !strings.Contains(v.Info.Chat.String(), "@g.us") {
		return false
	}

	groupInfo, err := client.GetGroupInfo(context.Background(), v.Info.Chat)
	if err != nil {
		return false
	}

	senderNum := v.Info.Sender.ToNonAD().User

	for _, participant := range groupInfo.Participants {
		if participant.JID.User == senderNum && (participant.IsAdmin || participant.IsSuperAdmin) {
			return true 
		}
	}

	return false
}

func handleGCName(client *whatsmeow.Client, v *events.Message, args string) {
	if !v.Info.IsGroup {
		replyMessage(client, v, "❌ This command can only be used in groups.")
		return
	}
	if args == "" {
		replyMessage(client, v, "❌ Please provide a new group name.")
		return
	}

	err := client.SetGroupName(context.Background(), v.Info.Chat, args)
	if err != nil {
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.")
	} else {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, "✅ Group name changed successfully!")
	}
}

func handleGCDP(client *whatsmeow.Client, v *events.Message) {
	if !v.Info.IsGroup {
		replyMessage(client, v, "❌ This command can only be used in groups.")
		return
	}

	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg == nil || extMsg.ContextInfo == nil || extMsg.ContextInfo.QuotedMessage == nil {
		replyMessage(client, v, "❌ Please reply to an image to set as group picture.")
		return
	}

	quoted := extMsg.ContextInfo.QuotedMessage
	if img := quoted.GetImageMessage(); img != nil {
		data, err := client.Download(context.Background(), img)
		if err == nil {
			_, err = client.SetGroupPhoto(context.Background(), v.Info.Chat, data)
			if err == nil {
				react(client, v.Info.Chat, v.Info.ID, "✅")
				replyMessage(client, v, "✅ Group picture changed successfully!")
			} else {
				replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.")
			}
		}
	} else {
		replyMessage(client, v, "❌ Please reply to a valid image.")
	}
}

func handleGCLink(client *whatsmeow.Client, v *events.Message) {
	if !v.Info.IsGroup {
		replyMessage(client, v, "❌ This command can only be used in groups.")
		return
	}

	link, err := client.GetGroupInviteLink(context.Background(), v.Info.Chat, false)
	if err != nil {
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.")
	} else {
		react(client, v.Info.Chat, v.Info.ID, "🔗")
		replyMessage(client, v, fmt.Sprintf("🔗 *Group Link:*\n%s", link))
	}
}

func handleGCRevoke(client *whatsmeow.Client, v *events.Message) {
	if !v.Info.IsGroup {
		replyMessage(client, v, "❌ This command can only be used in groups.")
		return
	}

	link, err := client.GetGroupInviteLink(context.Background(), v.Info.Chat, true)
	if err != nil {
		replyMessage(client, v, "❌ Action Failed! I am probably not an Admin.")
	} else {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, fmt.Sprintf("✅ *Group Link Revoked!*\nNew Link: %s", link))
	}
}



func handleAddStatus(client *whatsmeow.Client, v *events.Message) {
	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg == nil || extMsg.ContextInfo == nil || extMsg.ContextInfo.QuotedMessage == nil {
		replyMessage(client, v, "❌ Please reply to an image or video to add it to your status!")
		return
	}

	quoted := extMsg.ContextInfo.QuotedMessage
	var data []byte
	var err error
	var msg waProto.Message

	if img := quoted.GetImageMessage(); img != nil {
		data, err = client.Download(context.Background(), img)
		if err == nil {
			up, _ := client.Upload(context.Background(), data, whatsmeow.MediaImage)
			msg.ImageMessage = &waProto.ImageMessage{
				URL: proto.String(up.URL), DirectPath: proto.String(up.DirectPath),
				MediaKey: up.MediaKey, Mimetype: proto.String("image/jpeg"),
				FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256,
				FileLength: proto.Uint64(uint64(len(data))),
			}
		}
	} else if vid := quoted.GetVideoMessage(); vid != nil {
		data, err = client.Download(context.Background(), vid)
		if err == nil {
			up, _ := client.Upload(context.Background(), data, whatsmeow.MediaVideo)
			msg.VideoMessage = &waProto.VideoMessage{
				URL: proto.String(up.URL), DirectPath: proto.String(up.DirectPath),
				MediaKey: up.MediaKey, Mimetype: proto.String("video/mp4"),
				FileEncSHA256: up.FileEncSHA256, FileSHA256: up.FileSHA256,
				FileLength: proto.Uint64(uint64(len(data))),
			}
		}
	} else if aud := quoted.GetAudioMessage(); aud != nil {
        replyMessage(client, v, "❌ Voice status not yet supported via bot.")
        return
    }

	if data == nil {
		replyMessage(client, v, "❌ Failed to extract media.")
		return
	}

	statusJID := types.NewJID("status", types.BroadcastServer)
	_, sendErr := client.SendMessage(context.Background(), statusJID, &msg)

	if sendErr == nil {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, "✅ Status updated successfully!")
	} else {
		replyMessage(client, v, "❌ Failed to update status.")
	}
}


func handleDP(ctx context.Context, client *whatsmeow.Client, v *events.Message, args string) {
	contextInfo := v.Message.GetExtendedTextMessage().GetContextInfo()
	if contextInfo == nil || contextInfo.GetQuotedMessage() == nil {
		replyMessage(client, v, "❌ *Raw Error:*\n```\nPlease reply to an image to set it as DP.\n```")
		return
	}

	quotedMsg := contextInfo.GetQuotedMessage()
	imageMsg := quotedMsg.GetImageMessage()
	if imageMsg == nil {
		replyMessage(client, v, "❌ *Raw Error:*\n```\nThe quoted message is not an image.\n```")
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "⏳")

	
	imageData, err := client.Download(ctx, imageMsg)
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ *Raw Error (Download):*\n```\n%v\n```", err))
		return
	}

	img, _, err := image.Decode(bytes.NewReader(imageData))
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ *Raw Error (Decode):*\n```\n%v\n```", err))
		return
	}

	
	bounds := img.Bounds()
	width, height := bounds.Dx(), bounds.Dy()
	size := width
	if height < size {
		size = height
	}

	x0 := bounds.Min.X + (width-size)/2
	y0 := bounds.Min.Y + (height-size)/2

	targetSize := 640
	if size < 640 {
		targetSize = size 
	}

	scaledImg := image.NewRGBA(image.Rect(0, 0, targetSize, targetSize))
	for y := 0; y < targetSize; y++ {
		for x := 0; x < targetSize; x++ {
			srcX := x0 + (x * size / targetSize)
			srcY := y0 + (y * size / targetSize)
			scaledImg.Set(x, y, img.At(srcX, srcY))
		}
	}

	var jpegBuf bytes.Buffer
	err = jpeg.Encode(&jpegBuf, scaledImg, &jpeg.Options{Quality: 70})
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ *Raw Error (Encode):*\n```\n%v\n```", err))
		return
	}

	
	targetClient := client 
	targetNumber := strings.TrimSpace(args)

	if targetNumber != "" {
		cleanTarget := getCleanID(targetNumber)

		clientsMutex.RLock()
		botClient, exists := activeClients[cleanTarget]
		clientsMutex.RUnlock()

		if exists && botClient != nil {
			targetClient = botClient 
		} else {
			replyMessage(client, v, fmt.Sprintf("❌ *Raw Error:*\n```\nBot %s is not active or not found in memory.\n```", targetNumber))
			return 
		}
	}

	
	_, err = targetClient.SetGroupPhoto(ctx, types.EmptyJID, jpegBuf.Bytes())
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ *Raw Error:*\n```\n%v\n```", err))
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "✅")
	if targetNumber != "" {
		replyMessage(client, v, fmt.Sprintf("✅ Profile picture successfully updated for bot *%s*", targetNumber))
	} else {
		replyMessage(client, v, "✅ My profile picture has been successfully updated!")
	}
}



func handleChangeName(ctx context.Context, client *whatsmeow.Client, v *events.Message, args string) {
	if !isOwner(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
	if args == "" {
		replyMessage(client, v, "❌ Please provide a new name.\nExample: `.setname MyBot`")
		return
	}

	client.Store.PushName = args
	
	err := client.Store.Save(ctx)
	if err != nil {
		replyMessage(client, v, "❌ Failed to save name in store.")
		return
	}

	
	err = client.SendPresence(ctx, types.PresenceAvailable)
	if err != nil {
		replyMessage(client, v, "❌ Failed to broadcast new profile name to WhatsApp.")
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "✅")
	replyMessage(client, v, "✅ Profile name successfully changed to *"+args+"*")
}





func handleGetStatus(ctx context.Context, client *whatsmeow.Client, v *events.Message, args string, db *sql.DB) {
	if args == "" {
		replyMessage(client, v, "❌ Please provide a phone number.\nExample: `.getstatus 923001234567`")
		return
	}

	cleanNumber := strings.TrimSpace(args)
	cleanNumber = strings.ReplaceAll(cleanNumber, "+", "")
	targetJID := cleanNumber + "@s.whatsapp.net"

	oneDayAgo := time.Now().Unix() - (24 * 60 * 60)
	query := "SELECT media_path, media_type, caption FROM cached_statuses WHERE sender_jid = ? AND timestamp >= ?"
	rows, err := db.QueryContext(ctx, query, targetJID, oneDayAgo)
	if err != nil {
		replyMessage(client, v, "❌ Database error: "+err.Error())
		return
	}
	defer rows.Close()

	count := 0
	for rows.Next() {
		var mediaPath, mediaType, caption string
		if err := rows.Scan(&mediaPath, &mediaType, &caption); err != nil {
			continue
		}

		fileBytes, err := os.ReadFile(mediaPath)
		if err != nil {
			continue 
		}

		var mType whatsmeow.MediaType
		if strings.Contains(mediaType, "image") {
			mType = whatsmeow.MediaImage
		} else if strings.Contains(mediaType, "video") {
			mType = whatsmeow.MediaVideo
		} else {
			replyMessage(client, v, fmt.Sprintf("📝 *Text Status:* \n\n%s", caption))
			count++
			continue
		}

		uploaded, err := client.Upload(ctx, fileBytes, mType)
		if err != nil {
			continue
		}

		
		var msg waProto.Message
		if mType == whatsmeow.MediaImage {
			msg.ImageMessage = &waProto.ImageMessage{
				Caption:       proto.String(caption),
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mediaType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileBytes))),
			}
		} else if mType == whatsmeow.MediaVideo {
			msg.VideoMessage = &waProto.VideoMessage{
				Caption:       proto.String(caption),
				URL:           proto.String(uploaded.URL),
				DirectPath:    proto.String(uploaded.DirectPath),
				MediaKey:      uploaded.MediaKey,
				Mimetype:      proto.String(mediaType),
				FileEncSHA256: uploaded.FileEncSHA256,
				FileSHA256:    uploaded.FileSHA256,
				FileLength:    proto.Uint64(uint64(len(fileBytes))),
			}
		}

		_, err = client.SendMessage(ctx, v.Info.Chat, &msg)
		if err == nil {
			count++
		}
	}

	if count == 0 {
		replyMessage(client, v, "❌ No cached statuses found for this contact within the last 24 hours.")
	} else {
		react(client, v.Info.Chat, v.Info.ID, "✅")
	}
}

func handleStatusCapture(ctx context.Context, client *whatsmeow.Client, v *events.Message) {
	botJID := types.NewJID(client.Store.ID.User, client.Store.ID.Server).String()
	msgID := v.Info.ID
	timestamp := v.Info.Timestamp.Unix()
	
	var err error 

	
	targetSenderJID := LIDToJID(ctx, client, v.Info.Sender)

	senderJIDStr := targetSenderJID.ToNonAD().String()

	var mediaPath string
	var mediaType string
	var caption string
	var fileBytes []byte

	if img := v.Message.GetImageMessage(); img != nil {
		mediaType = img.GetMimetype()
		caption = img.GetCaption()
		fileBytes, err = client.Download(ctx, img)
	} else if vid := v.Message.GetVideoMessage(); vid != nil {
		mediaType = vid.GetMimetype()
		caption = vid.GetCaption()
		fileBytes, err = client.Download(ctx, vid)
	} else if txt := v.Message.GetConversation(); txt != "" {
		mediaType = "text/plain"
		caption = txt
	} else if extTxt := v.Message.GetExtendedTextMessage(); extTxt != nil {
		mediaType = "text/plain"
		caption = extTxt.GetText()
	}

	if len(fileBytes) > 0 && err == nil {
		_ = os.MkdirAll("./data/statuses", os.ModePerm)
		ext := "jpg"
		if strings.Contains(mediaType, "video") {
			ext = "mp4"
		}
		mediaPath = fmt.Sprintf("./data/statuses/%s.%s", msgID, ext)
		_ = os.WriteFile(mediaPath, fileBytes, 0644)
	}

	query := `INSERT INTO cached_statuses (msg_id, bot_jid, sender_jid, media_path, media_type, caption, timestamp) 
			  VALUES (?, ?, ?, ?, ?, ?, ?) 
			  ON CONFLICT(msg_id) DO NOTHING;`
			  
	_, err = settingsDB.ExecContext(ctx, query, msgID, botJID, senderJIDStr, mediaPath, mediaType, caption, timestamp)
}


func handleGetDP(ctx context.Context, client *whatsmeow.Client, v *events.Message, args string) {
	var targetJID types.JID

	
	if args == "" {
		targetJID = v.Info.Chat
	} else {
		cleanNumber := strings.TrimSpace(args)
		cleanNumber = strings.ReplaceAll(cleanNumber, "+", "")
		targetJID = types.NewJID(cleanNumber, types.LegacyUserServer)
	}

	
	react(client, v.Info.Chat, v.Info.ID, "⏳")

	
	params := &whatsmeow.GetProfilePictureParams{Preview: false}
	picInfo, err := client.GetProfilePictureInfo(ctx, targetJID, params)
	if err != nil || picInfo == nil || picInfo.URL == "" {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	
	req, err := http.NewRequestWithContext(ctx, "GET", picInfo.URL, nil)
	if err != nil {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		if resp != nil {
			resp.Body.Close()
		}
		return
	}
	defer resp.Body.Close()

	imgBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	
	uploaded, err := client.Upload(ctx, imgBytes, whatsmeow.MediaImage)
	if err != nil {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	
	var msg waProto.Message
	msg.ImageMessage = &waProto.ImageMessage{
		URL:           proto.String(uploaded.URL),
		DirectPath:    proto.String(uploaded.DirectPath),
		MediaKey:      uploaded.MediaKey,
		Mimetype:      proto.String("image/jpeg"),
		FileEncSHA256: uploaded.FileEncSHA256,
		FileSHA256:    uploaded.FileSHA256,
		FileLength:    proto.Uint64(uint64(len(imgBytes))),
	}

	
	_, err = client.SendMessage(ctx, v.Info.Chat, &msg)
	if err != nil {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	
	react(client, v.Info.Chat, v.Info.ID, "✅")
}


func handleCheckWA(ctx context.Context, client *whatsmeow.Client, v *events.Message, args string) {
	var numbers []string

	
	cleanNum := func(val string) string {
		val = strings.TrimSpace(val)
		val = strings.ReplaceAll(val, "+", "")
		val = strings.ReplaceAll(val, "-", "")
		val = strings.ReplaceAll(val, " ", "")
		return val
	}

	
	contextInfo := v.Message.GetExtendedTextMessage().GetContextInfo()
	var docMsg *waProto.DocumentMessage
	if contextInfo != nil && contextInfo.GetQuotedMessage() != nil {
		docMsg = contextInfo.GetQuotedMessage().GetDocumentMessage()
	}

	if docMsg != nil && (strings.Contains(docMsg.GetMimetype(), "text/plain") || strings.HasSuffix(docMsg.GetFileName(), ".txt")) {
		react(client, v.Info.Chat, v.Info.ID, "⏳")
		
		fileBytes, err := client.Download(ctx, docMsg)
		if err != nil {
			react(client, v.Info.Chat, v.Info.ID, "❌")
			return
		}

		lines := strings.Split(string(fileBytes), "\n")
		for _, line := range lines {
			if num := cleanNum(line); num != "" {
				numbers = append(numbers, num)
			}
		}
	} else {
		
		if args == "" {
			react(client, v.Info.Chat, v.Info.ID, "❌")
			return
		}
		react(client, v.Info.Chat, v.Info.ID, "⏳")

		lines := strings.Split(args, "\n")
		for _, line := range lines {
			line = strings.ReplaceAll(line, ",", " ")
			subFields := strings.Fields(line)
			for _, field := range subFields {
				if num := cleanNum(field); num != "" {
					numbers = append(numbers, num)
				}
			}
		}
	}

	if len(numbers) == 0 {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	var activeList []string
	var inactiveList []string

	
	for _, num := range numbers {
		
		time.Sleep(2500 * time.Millisecond)

		res, err := client.IsOnWhatsApp(ctx, []string{num})
		if err != nil {
			continue 
		}

		if len(res) > 0 && res[0].IsIn {
			activeList = append(activeList, num)
		} else {
			inactiveList = append(inactiveList, num)
		}
	}

	
	sendFileResult := func(filename, content, caption string) {
		bytesData := []byte(content)
		if len(bytesData) == 0 {
			bytesData = []byte("No numbers filtered in this category.")
		}

		uploaded, err := client.Upload(ctx, bytesData, whatsmeow.MediaDocument)
		if err != nil {
			return
		}

		var msg waProto.Message
		msg.DocumentMessage = &waProto.DocumentMessage{
			URL:           proto.String(uploaded.URL),
			DirectPath:    proto.String(uploaded.DirectPath),
			MediaKey:      uploaded.MediaKey,
			Mimetype:      proto.String("text/plain"),
			FileName:      proto.String(filename),
			FileEncSHA256: uploaded.FileEncSHA256,
			FileSHA256:    uploaded.FileSHA256,
			FileLength:    proto.Uint64(uint64(len(bytesData))),
			Caption:       proto.String(caption),
		}
		_, _ = client.SendMessage(ctx, v.Info.Chat, &msg)
	}

	
	activeRaw := strings.Join(activeList, "\n")
	inactiveRaw := strings.Join(inactiveList, "\n")

	sendFileResult("active_whatsapp.txt", activeRaw, fmt.Sprintf("✅ Active WhatsApp Numbers: %d", len(activeList)))
	sendFileResult("inactive_whatsapp.txt", inactiveRaw, fmt.Sprintf("❌ Non-Active Numbers: %d", len(inactiveList)))

	react(client, v.Info.Chat, v.Info.ID, "✅")
}


func handleSetAbout(ctx context.Context, client *whatsmeow.Client, v *events.Message, args string) {
	if !isOwner(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
	if args == "" { react(client, v.Info.Chat, v.Info.ID, "❌"); return }

	err := client.SetStatusMessage(ctx, args)
	if err != nil {
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "✅")
}

func LIDToJID(ctx context.Context, client *whatsmeow.Client, lid types.JID) types.JID {
	
	if lid.Server != types.HiddenUserServer { 
		return lid
	}

	
	altJID, err := client.Store.GetAltJID(ctx, lid)
	if err == nil && !altJID.IsEmpty() {
		return altJID
	}

	
	fmt.Printf("\n⚠️ [LIDToJID] Local store map missing for LID: %s. Initiating live server sync...\n", lid.String())
	
	
	fmt.Printf("📤 [RAW OUTGOING REQUEST TO WHATSAPP] Target JID Payload Struct: %+v\n", lid)

	
	resp, err := client.GetUserInfo(ctx, []types.JID{lid})
	if err != nil {
		
		fmt.Printf("❌ [RAW NETWORK ERROR] WhatsApp server rejected USync/UserInfo query: %v\n", err)
	} else if resp != nil {
		
		respJSON, _ := json.MarshalIndent(resp, "", "  ")
		fmt.Printf("📥 [RAW INCOMING RESPONSE FROM WHATSAPP]\n%s\n", string(respJSON))
	}

	if err == nil {
		
		updatedAltJID, err := client.Store.GetAltJID(ctx, lid)
		if err == nil && !updatedAltJID.IsEmpty() {
			fmt.Printf("✅ [LIDToJID SUCCESS] Safely resolved and cached real JID: %s\n\n", updatedAltJID.String())
			return updatedAltJID
		}
	}

	
	fmt.Printf("❌ [LIDToJID FAILED] Network sync completed but no phone mapping found for LID: %s\n\n", lid.String())
	return lid
}


func handleBlocklist(ctx context.Context, client *whatsmeow.Client, v *events.Message) {
	if !isOwner(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }

	
	list, err := client.GetBlocklist(ctx)
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ *Raw Error (GetBlocklist):*\n```\n%v\n```", err))
		return
	}

	if len(list.JIDs) == 0 {
		replyMessage(client, v, "📝 Your blocklist is empty.")
		return
	}

	var sb strings.Builder
	sb.WriteString(fmt.Sprintf("🚫 *Bot Blocklist (%d Real Numbers Resolved):*\n\n", len(list.JIDs)))
	
	for i, jid := range list.JIDs {
		
		resolvedJID := LIDToJID(ctx, client, jid)
		
		sb.WriteString(fmt.Sprintf("%d. +%s\n", i+1, resolvedJID.User))
	}

	replyMessage(client, v, sb.String())
}


func sendCustomBlockNode(ctx context.Context, client *whatsmeow.Client, currentDHash string, targetJID types.JID) error {
	iqID := fmt.Sprintf("silent_perfect_block_%d", time.Now().UnixNano())
	
	blockNode := waBinary.Node{
		Tag: "iq",
		Attrs: waBinary.Attrs{
			"id":    iqID,
			"to":    types.ServerJID,
			"type":  "set",
			"xmlns": "blocklist",
		},
		Content: []waBinary.Node{{
			Tag: "list",
			Attrs: waBinary.Attrs{
				"dhash": currentDHash,
			},
			Content: []waBinary.Node{{
				Tag: "item",
				Attrs: waBinary.Attrs{
					"action": "block",
					"jid":    targetJID,
				},
			}},
		}},
	}

	return whatsmeowSendNode(client, ctx, blockNode)
}


func isUserSuccessfullyBlocked(targetLID, targetJID types.JID, blockedJIDs []types.JID) bool {
	for _, jid := range blockedJIDs {
		if jid.User == targetLID.User || jid.User == targetJID.User {
			return true
		}
	}
	return false
}


func handleBlock(ctx context.Context, client *whatsmeow.Client, v *events.Message) {
	if !isOwner(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }

	targetJID := v.Info.Chat
	if v.Info.IsGroup {
		targetJID = v.Info.Sender
	}

	
	realPhoneJID := LIDToJID(ctx, client, targetJID)

	
	initialList, err := client.GetBlocklist(ctx)
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ *Failed to fetch current blocklist:* `%v`", err))
		return
	}
	currentDHash := initialList.DHash

	
	
	
	_, errNativeLID := client.UpdateBlocklist(ctx, targetJID, events.BlocklistChangeActionBlock)
	time.Sleep(1200 * time.Millisecond)
	
	check1, _ := client.GetBlocklist(ctx)
	if isUserSuccessfullyBlocked(targetJID, realPhoneJID, check1.JIDs) {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, fmt.Sprintf("✅ *Block Successful!*\n\n• *Method Worked:* `Native WhatsMeow (LID)`\n• *Target:* `%s`", targetJID.String()))
		return
	}

	
	
	
	_, errNativeJID := client.UpdateBlocklist(ctx, realPhoneJID, events.BlocklistChangeActionBlock)
	time.Sleep(1200 * time.Millisecond)
	
	check2, _ := client.GetBlocklist(ctx)
	if isUserSuccessfullyBlocked(targetJID, realPhoneJID, check2.JIDs) {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, fmt.Sprintf("✅ *Block Successful!*\n\n• *Method Worked:* `Native WhatsMeow (Phone JID)`\n• *Target:* `%s`", realPhoneJID.String()))
		return
	}

	
	
	
	_ = sendCustomBlockNode(ctx, client, currentDHash, targetJID)
	time.Sleep(1200 * time.Millisecond)
	
	check3, _ := client.GetBlocklist(ctx)
	if isUserSuccessfullyBlocked(targetJID, realPhoneJID, check3.JIDs) {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, fmt.Sprintf("✅ *Block Successful!*\n\n• *Method Worked:* `Custom XML list-Wrapper (LID)`\n• *Target:* `%s`", targetJID.String()))
		return
	}

	
	
	
	_ = sendCustomBlockNode(ctx, client, currentDHash, realPhoneJID)
	time.Sleep(1200 * time.Millisecond)
	
	check4, _ := client.GetBlocklist(ctx)
	if isUserSuccessfullyBlocked(targetJID, realPhoneJID, check4.JIDs) {
		react(client, v.Info.Chat, v.Info.ID, "✅")
		replyMessage(client, v, fmt.Sprintf("✅ *Block Successful!*\n\n• *Method Worked:* `Custom XML list-Wrapper (Phone JID)`\n• *Target:* `%s`", realPhoneJID.String()))
		return
	}

	
	
	
	failSummary := fmt.Sprintf(
		"❌ *All 4 Block Strategies Failed!*\n\n• *Native LID Error:* `%v`\n• *Native JID Error:* `%v`\n\nℹ️ _چاروں پے لوڈز سرور نے مسترد کر دیے۔ ٹرمینل پر سیشن ہینڈ شیک یا کنکشن ڈراپس چیک کریں!_", 
		errNativeLID, errNativeJID,
	)
	replyMessage(client, v, failSummary)
}


func whatsmeowSendNode(cli *whatsmeow.Client, ctx context.Context, node waBinary.Node) error
