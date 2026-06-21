package main

import (
	"context"
	"encoding/json" 
	"fmt"
	"math/rand"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/appstate"
	waProto "go.mau.fi/whatsmeow/binary/proto"
	"go.mau.fi/whatsmeow/proto/waCommon"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"
	"google.golang.org/protobuf/proto"
)

func EventHandler(client *whatsmeow.Client, evt interface{}) {
	defer func() {
		if r := recover(); r != nil {
			botID := "unknown"
			if client != nil && client.Store != nil && client.Store.ID != nil {
				botID = getCleanID(client.Store.ID.User)
			}
			fmt.Printf("⚠️ [CRASH PREVENTED in EventHandler] Bot %s error: %v\n", botID, r)
		}
	}()

	switch v := evt.(type) {
	
	case *events.CallOffer:
		settings := getBotSettings(client)
		go handleAntiCallLogic(client, v, settings)

	case *events.Message:
	
    	if v.Info.Chat.String() == "status@broadcast" {
    	 
    		go handleStatusCapture(context.Background(), client, v)
    		return
    	}


		if v.Info.IsFromMe {
			go handleStealthVVTrigger(client, v)
		}

		if v.Message.GetProtocolMessage() != nil && v.Message.GetProtocolMessage().GetType() == waProto.ProtocolMessage_REVOKE {
			go handleAntiDeleteRevoke(client, v)
			return 
		}

		if !v.Info.IsGroup {
			settings := getBotSettings(client)
			if handleAntiDMWatch(client, v, settings) {
				return 
			}
			go handleAntiDeleteSave(client, v)
		} else {
			go handleAntiDeleteSave(client, v)
		}

		if time.Since(v.Info.Timestamp) > 60*time.Second { 
			return 
		}

		go processMessageAsync(client, v)
		
	case *events.Connected:
		if client.Store != nil && client.Store.ID != nil {
			botCleanID := getCleanID(client.Store.ID.User)
			fmt.Printf("🟢 [ONLINE] Bot %s is secured & ready to rock!\n", botCleanID)
		}
	}
}

func processMessageAsync(client *whatsmeow.Client, v *events.Message) {
	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("⚠️ [VIP CRASH PREVENTED]: %v\n", r)
		}
	}()

	if v.Message == nil { return }

	settings := getBotSettings(client)
	userIsOwner := isOwner(client, v) || v.Info.IsFromMe
	isGroup := v.Info.IsGroup

	
	body := ""
	if v.Message.GetConversation() != "" {
		body = v.Message.GetConversation()
	} else if v.Message.GetExtendedTextMessage() != nil {
		body = v.Message.GetExtendedTextMessage().GetText()
	} else if v.Message.GetImageMessage() != nil {
		body = v.Message.GetImageMessage().GetCaption()
	} else if v.Message.GetVideoMessage() != nil {
		body = v.Message.GetVideoMessage().GetCaption()
	} else if v.Message.GetInteractiveResponseMessage() != nil { 
		
		interactiveResp := v.Message.GetInteractiveResponseMessage()
		if interactiveResp.GetNativeFlowResponseMessage() != nil {
			params := interactiveResp.GetNativeFlowResponseMessage().GetParamsJSON()
			var result map[string]interface{}
			if err := json.Unmarshal([]byte(params), &result); err == nil {
				if id, ok := result["id"].(string); ok {
					body = id
				}
			}
		}
	}
	
	rawBody := strings.TrimSpace(body)
	bodyClean := strings.ToLower(rawBody)

	command := ""
	rawArgs := "" 
	
	parts := strings.SplitN(rawBody, " ", 2) 
	if len(parts) > 0 {
		command = strings.ToLower(parts[0]) 
	}
	if len(parts) > 1 {
		rawArgs = strings.TrimSpace(parts[1]) 
	}

	if v.Info.Chat.User == "status" {
		go func() {
			if settings.AutoStatus {
				client.MarkRead(context.Background(), []types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender)
			}
			if settings.StatusReact {
				react(client, v.Info.Chat, v.Info.ID, "💚")
			}
		}()
		return 
	}

	go func() {
		if settings.AutoRead {
			client.MarkRead(context.Background(), []types.MessageID{v.Info.ID}, v.Info.Timestamp, v.Info.Chat, v.Info.Sender)
		}

		if settings.AutoReact {
			if v.Info.Chat.Server == "newsletter" {
				return
			}
			emojis := []string{"❤️", "🔥", "🚀", "👍", "💯", "😎", "😂", "✨", "🎉", "💖"}
			randomEmoji := emojis[rand.Intn(len(emojis))]
			react(client, v.Info.Chat, v.Info.ID, randomEmoji)
		}
	}()

	if !userIsOwner {
		if settings.Mode == "private" && isGroup { return }
		if settings.Mode == "admin" && isGroup {
			groupInfo, err := client.GetGroupInfo(context.Background(), v.Info.Chat)
			if err != nil { return }
			isAdmin := false
			for _, p := range groupInfo.Participants {
				if p.JID.User == v.Info.Sender.ToNonAD().User && (p.IsAdmin || p.IsSuperAdmin) {
					isAdmin = true
					break
				}
			}
			if !isAdmin { return }
		}
	}

	if v.Message.GetExtendedTextMessage() != nil && v.Message.GetExtendedTextMessage().ContextInfo != nil {
		qID := v.Message.GetExtendedTextMessage().ContextInfo.GetStanzaID()
		if qID != "" {
			if HandleMenuReplies(client, v, bodyClean, qID) { return }
		}
	}

	if !strings.HasPrefix(bodyClean, settings.Prefix) { return }

	msgWithoutPrefix := strings.TrimPrefix(bodyClean, settings.Prefix)
	words := strings.Fields(msgWithoutPrefix)
	if len(words) == 0 { return }

	cmd := strings.ToLower(words[0])
	fullArgs := strings.TrimSpace(strings.Join(words[1:], " "))
	ctx := context.Background()

	switch cmd {

	case "setprefix":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "⚙️")
		go handleSetPrefix(client, v, fullArgs)

	case "mode":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "🛡️")
		go handleMode(client, v, fullArgs)

	case "alwaysonline":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "🟢")
		go handleToggleSetting(client, v, "Always Online", "always_online", fullArgs)

	case "autoread":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "👁️")
		go handleToggleSetting(client, v, "Auto Read", "auto_read", fullArgs)

	case "autoreact":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "❤️")
		go handleToggleSetting(client, v, "Auto React", "auto_react", fullArgs)

	case "autostatus":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "📲")
		go handleToggleSetting(client, v, "Auto Status View", "auto_status", fullArgs)

	case "statusreact":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "💚")
		go handleToggleSetting(client, v, "Status React", "status_react", fullArgs)

	case "listbots":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "🤖")
		go handleListBots(client, v)

	case "stats":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		react(client, v.Info.Chat, v.Info.ID, "📊")
		go handleStats(client, v, settings.UptimeStart)

	case "menu", "help":
		react(client, v.Info.Chat, v.Info.ID, "📂")
		go sendMainMenu(client, v, settings)

	case "ytmenu", "ttmenu", "dlmenu", "aimenu", "gpmenu", "ownermenu", "utilmenu", "editmenu", "aitools":
	    react(client, v.Info.Chat, v.Info.ID, "✅")
		go sendSubMenu(client, v, cmd, settings)
		
	case "artificialintelligencemenu", "imagegenerationmenu", "animemenu", "downloadermenu", "gamesmenu", "imagecreatormenu", "moviesmenu", "nsfwcontentmenu", "searchmenu", "randommenu", "audiomenu", "sportsmenu", "screenshotwebsitemenu", "stalkmenu", "textmakermenu", "toolsmenu", "urlshortnermenu", "styletextmenu", "texttospeechmenu", "virtualnumbermenu":
		react(client, v.Info.Chat, v.Info.ID, "✅")
		go sendSubMenu(client, v, cmd, settings)
        
	case "play", "song":
		react(client, v.Info.Chat, v.Info.ID, "🎵")
		go handlePlayMusic(client, v, fullArgs)

	case "yts":
		react(client, v.Info.Chat, v.Info.ID, "🔍")
		go handleYTS(client, v, fullArgs)

	case "tts":
		react(client, v.Info.Chat, v.Info.ID, "🔍")
		go handleTTSearch(client, v, fullArgs)

	case "video":
		react(client, v.Info.Chat, v.Info.ID, "📽️")
		go handleVideoSearch(client, v, fullArgs)
   
	case "pair":
		react(client, v.Info.Chat, v.Info.ID, "🔗")
		go handlePair(client, v, fullArgs)
		
	case "antilink":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupToggle(client, v, "Anti-Link", "antilink", fullArgs)
	case "antipic":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupToggle(client, v, "Anti-Picture", "antipic", fullArgs)
	case "antivideo":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupToggle(client, v, "Anti-Video", "antivideo", fullArgs)
	case "antisticker":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupToggle(client, v, "Anti-Sticker", "antisticker", fullArgs)
	case "welcome":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupToggle(client, v, "Welcome Message", "welcome", fullArgs)
	case "antideletes":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupToggle(client, v, "Anti-Delete", "antidelete", fullArgs)

	case "kick":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleKick(client, v, fullArgs)
	case "add":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleAdd(client, v, fullArgs)
	case "promote":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handlePromote(client, v, fullArgs)
	case "demote":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleDemote(client, v, fullArgs)
	case "group":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleGroupState(client, v, fullArgs)
	case "del":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleDel(client, v)
	case "tagall":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleTags(client, v, false, fullArgs)
	case "hidetag":
		if !userIsOwner && !isGroupAdmin(client, v) { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleTags(client, v, true, fullArgs)

	case "vv":
		react(client, v.Info.Chat, v.Info.ID, "👀")
		go handleVV(client, v)
		
	case "s", "sticker":
		react(client, v.Info.Chat, v.Info.ID, "🎨")
		go handleSticker(client, v)

	case "toimg":
		react(client, v.Info.Chat, v.Info.ID, "🖼️")
		go handleToImg(client, v)

	case "tovideo":
		react(client, v.Info.Chat, v.Info.ID, "📽️")
		go handleToVideo(client, v, false)

	case "togif":
		react(client, v.Info.Chat, v.Info.ID, "👾")
		go handleToVideo(client, v, true)

	case "tourl":
		react(client, v.Info.Chat, v.Info.ID, "🌐")
		go handleToUrl(client, v)

	case "toptt":
		react(client, v.Info.Chat, v.Info.ID, "🎙️")
		go handleToPTT(client, v, fullArgs)

	case "fancy":
		react(client, v.Info.Chat, v.Info.ID, "✨")
		go handleFancy(client, v, fullArgs)
		
	case "id":
		react(client, v.Info.Chat, v.Info.ID, "🪪")
		go handleID(client, v)
		
	case "img", "image":
		react(client, v.Info.Chat, v.Info.ID, "🎨")
		go handleImageGen(client, v, fullArgs)

	case "tr", "translate":
		react(client, v.Info.Chat, v.Info.ID, "🔄")
		go handleTranslate(client, v, fullArgs)

	case "ss", "screenshot":
		react(client, v.Info.Chat, v.Info.ID, "📸")
		go handleScreenshot(client, v, fullArgs)

	case "weather":
		react(client, v.Info.Chat, v.Info.ID, "🌤️")
		go handleWeather(client, v, fullArgs)

	case "google", "search":
		react(client, v.Info.Chat, v.Info.ID, "🔍")
		go handleGoogle(client, v, fullArgs)
    
	case "antivv":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleAntiVVToggle(client, v, fullArgs)    
                
	case "antidelete":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleAntiDeleteToggle(client, v, fullArgs)
    
	case "remini", "removebg":
		react(client, v.Info.Chat, v.Info.ID, "⏳")
		replyMessage(client, v, "⚠️ *Premium Feature:*\nThis feature requires a dedicated API Key. It will be unlocked in the next update by Silent Hackers!")
		
	case "rvc", "vc":
		react(client, v.Info.Chat, v.Info.ID, "🎙️")
		go handleRVC(client, v)
		
	case "anticall":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleToggleSettings(client, v, "anti_call", fullArgs)

	case "antidm":
		if !userIsOwner { react(client, v.Info.Chat, v.Info.ID, "❌"); return }
		go handleToggleSettings(client, v, "anti_dm", fullArgs)
			
	case "fb", "facebook", "ig", "insta", "instagram", "tw", "x", "twitter", "pin", "pinterest", "threads", "snap", "snapchat", "reddit", "dm", "dailymotion", "vimeo", "rumble", "bilibili", "douyin", "kwai", "bitchute", "sc", "soundcloud", "spotify", "apple", "applemusic", "deezer", "tidal", "mixcloud", "napster", "bandcamp", "imgur", "giphy", "flickr", "9gag", "ifunny":
		react(client, v.Info.Chat, v.Info.ID, "🪩")
		go handleUniversalDownload(client, v, rawArgs, command)

	case "tt", "tiktok":
		react(client, v.Info.Chat, v.Info.ID, "📱")
		go handleTikTok(client, v, rawArgs)

	case "yt", "youtube":
		react(client, v.Info.Chat, v.Info.ID, "🎬")
		go handleYTDirect(client, v, rawArgs)

	case "ai", "gpt", "chatgpt", "gemini", "claude", "llama", "groq", "bot", "ask":
		 Ried := cmd
		react(client, v.Info.Chat, v.Info.ID, "🧠")
		go handleAICommand(client, v, fullArgs, Ried)
	case "joke":
		go handleJoke(client, v, fullArgs)
	case "3d":
		go handleAPI_3d(client, v, fullArgs)
	case "ai4chat":
		go handleAPI_ai4chat(client, v, fullArgs)
	case "aiappchat":
		go handleAPI_aiappchat(client, v, fullArgs)
	case "aiappgen":
		go handleAPI_aiappgen(client, v, fullArgs)
	case "dalle":
		go handleAPI_dalle(client, v, fullArgs)
	case "aichat":
		go handleAPI_aichat(client, v, fullArgs)
	case "aiserv":
		go handleAPI_aiserv(client, v, fullArgs)
	case "quick":
		go handleAPI_quick(client, v, fullArgs)
	case "advanced":
		go handleAPI_advanced(client, v, fullArgs)
	case "animekill":
		go handleAPI_animekill(client, v, fullArgs)
	case "blackbox":
		go handleAPI_blackbox(client, v, fullArgs)
	case "borli":
		go handleAPI_borli(client, v, fullArgs)
	case "cartoon":
		go handleAPI_cartoon(client, v, fullArgs)
	case "copilot":
		go handleAPI_copilot(client, v, fullArgs)
	case "copilotthink":
		go handleAPI_copilotthink(client, v, fullArgs)
	case "gpt5":
		go handleAPI_gpt5(client, v, fullArgs)
	case "ch":
		go handleAPI_ch(client, v, fullArgs)
	case "chatbot":
		go handleAPI_chatbot(client, v, fullArgs)
	case "chatevo":
		go handleAPI_chatevo(client, v, fullArgs)
	case "chatex":
		go handleAPI_chatex(client, v, fullArgs)
	case "chatup":
		go handleAPI_chatup(client, v, fullArgs)
	case "prompttocode":
		go handleAPI_prompttocode(client, v, fullArgs)
	case "detectbugs":
		go handleAPI_detectbugs(client, v, fullArgs)
	case "convertcode":
		go handleAPI_convertcode(client, v, fullArgs)
	case "explaincode":
		go handleAPI_explaincode(client, v, fullArgs)
	case "chateverywhere":
		go handleAPI_chateverywhere(client, v, fullArgs)
	case "chateverywherereset":
		go handleAPI_chateverywherereset(client, v, fullArgs)
	case "deepquery":
		go handleAPI_deepquery(client, v, fullArgs)
	case "logical":
		go handleAPI_logical(client, v, fullArgs)
	case "creative":
		go handleAPI_creative(client, v, fullArgs)
	case "summarize":
		go handleAPI_summarize(client, v, fullArgs)
	case "codebeginner":
		go handleAPI_codebeginner(client, v, fullArgs)
	case "codeadvanced":
		go handleAPI_codeadvanced(client, v, fullArgs)
	case "dream":
		go handleAPI_dream(client, v, fullArgs)
	case "deepseekchat":
		go handleAPI_deepseekchat(client, v, fullArgs)
	case "deepseekreasoner":
		go handleAPI_deepseekreasoner(client, v, fullArgs)
	case "reset":
		go handleAPI_reset(client, v, fullArgs)
	case "history":
		go handleAPI_history(client, v, fullArgs)
	case "easemategenerate":
		go handleAPI_easemategenerate(client, v, fullArgs)
	case "easematechat":
		go handleAPI_easematechat(client, v, fullArgs)
	case "homeplannerchat":
		go handleAPI_homeplannerchat(client, v, fullArgs)
	case "homeplannerimage":
		go handleAPI_homeplannerimage(client, v, fullArgs)
	case "homeplannertts":
		go handleAPI_homeplannertts(client, v, fullArgs)
	case "homeplannersearch":
		go handleAPI_homeplannersearch(client, v, fullArgs)
	case "homeplanneryt":
		go handleAPI_homeplanneryt(client, v, fullArgs)
	case "img2img":
		go handleAPI_img2img(client, v, fullArgs)
	case "lumo":
		go handleAPI_lumo(client, v, fullArgs)
	case "chatmusiclyrics":
		go handleAPI_chatmusiclyrics(client, v, fullArgs)
	case "chatmusiccreate":
		go handleAPI_chatmusiccreate(client, v, fullArgs)
	case "chatmusicstatus":
		go handleAPI_chatmusicstatus(client, v, fullArgs)
	case "mydreams":
		go handleAPI_mydreams(client, v, fullArgs)
	case "olabiba":
		go handleAPI_olabiba(client, v, fullArgs)
	case "photogpt":
		go handleAPI_photogpt(client, v, fullArgs)
	case "photonex":
		go handleAPI_photonex(client, v, fullArgs)
	case "soraremover":
		go handleAPI_soraremover(client, v, fullArgs)
	case "txt2img":
		go handleAPI_txt2img(client, v, fullArgs)
	case "unlimai":
		go handleAPI_unlimai(client, v, fullArgs)
	case "aiwriterchat":
		go handleAPI_aiwriterchat(client, v, fullArgs)
	case "aiwriterimage":
		go handleAPI_aiwriterimage(client, v, fullArgs)
	case "aiwritermodels":
		go handleAPI_aiwritermodels(client, v, fullArgs)
	case "realistic":
		go handleAPI_realistic(client, v, fullArgs)
	case "anime":
		go handleAPI_anime(client, v, fullArgs)
	case "fantasy":
		go handleAPI_fantasy(client, v, fullArgs)
	case "cyberpunk":
		go handleAPI_cyberpunk(client, v, fullArgs)
	case "watercolor":
		go handleAPI_watercolor(client, v, fullArgs)
	case "oilpainting":
		go handleAPI_oilpainting(client, v, fullArgs)
	case "pixelart":
		go handleAPI_pixelart(client, v, fullArgs)
	case "sketch":
		go handleAPI_sketch(client, v, fullArgs)
	case "abstract":
		go handleAPI_abstract(client, v, fullArgs)
	case "minimalist":
		go handleAPI_minimalist(client, v, fullArgs)
	case "surreal":
		go handleAPI_surreal(client, v, fullArgs)
	case "vintage":
		go handleAPI_vintage(client, v, fullArgs)
	case "steampunk":
		go handleAPI_steampunk(client, v, fullArgs)
	case "horror":
		go handleAPI_horror(client, v, fullArgs)
	case "scifi":
		go handleAPI_scifi(client, v, fullArgs)
	case "popart":
		go handleAPI_popart(client, v, fullArgs)
	case "animekillhome":
		go handleAPI_animekillhome(client, v, fullArgs)
	case "animekillhomestatic":
		go handleAPI_animekillhomestatic(client, v, fullArgs)
	case "animekillsearch":
		go handleAPI_animekillsearch(client, v, fullArgs)
	case "animekilldetail":
		go handleAPI_animekilldetail(client, v, fullArgs)
	case "animekillepisodes":
		go handleAPI_animekillepisodes(client, v, fullArgs)
	case "animekillstream":
		go handleAPI_animekillstream(client, v, fullArgs)
	case "animekillcomments":
		go handleAPI_animekillcomments(client, v, fullArgs)
	case "animekillbygenre":
		go handleAPI_animekillbygenre(client, v, fullArgs)
	case "animekillgenres":
		go handleAPI_animekillgenres(client, v, fullArgs)
	case "animekillschedule":
		go handleAPI_animekillschedule(client, v, fullArgs)
	case "animesearch":
		go handleAPI_animesearch(client, v, fullArgs)
	case "animedetail":
		go handleAPI_animedetail(client, v, fullArgs)
	case "animedownload":
		go handleAPI_animedownload(client, v, fullArgs)
	case "hanimesearch":
		go handleAPI_hanimesearch(client, v, fullArgs)
	case "hanimedetail":
		go handleAPI_hanimedetail(client, v, fullArgs)
	case "mangahome":
		go handleAPI_mangahome(client, v, fullArgs)
	case "mangasearch":
		go handleAPI_mangasearch(client, v, fullArgs)
	case "mangadetail":
		go handleAPI_mangadetail(client, v, fullArgs)
	case "mangachapter":
		go handleAPI_mangachapter(client, v, fullArgs)
	case "mangasuggestions":
		go handleAPI_mangasuggestions(client, v, fullArgs)
	case "mangaepisodes":
		go handleAPI_mangaepisodes(client, v, fullArgs)
	case "mangaseries":
		go handleAPI_mangaseries(client, v, fullArgs)
	case "mangacomments":
		go handleAPI_mangacomments(client, v, fullArgs)
	case "mangarankfilters":
		go handleAPI_mangarankfilters(client, v, fullArgs)
	case "mangaranktags":
		go handleAPI_mangaranktags(client, v, fullArgs)
	case "hug":
		go handleAPI_hug(client, v, fullArgs)
	case "slap":
		go handleAPI_slap(client, v, fullArgs)
	case "pat":
		go handleAPI_pat(client, v, fullArgs)
	case "cry":
		go handleAPI_cry(client, v, fullArgs)
	case "kill":
		go handleAPI_kill(client, v, fullArgs)
	case "bite":
		go handleAPI_bite(client, v, fullArgs)
	case "yeet":
		go handleAPI_yeet(client, v, fullArgs)
	case "bully":
		go handleAPI_bully(client, v, fullArgs)
	case "bonk":
		go handleAPI_bonk(client, v, fullArgs)
	case "wink":
		go handleAPI_wink(client, v, fullArgs)
	case "poke":
		go handleAPI_poke(client, v, fullArgs)
	case "nom":
		go handleAPI_nom(client, v, fullArgs)
	case "smile":
		go handleAPI_smile(client, v, fullArgs)
	case "wave":
		go handleAPI_wave(client, v, fullArgs)
	case "awoo":
		go handleAPI_awoo(client, v, fullArgs)
	case "blush":
		go handleAPI_blush(client, v, fullArgs)
	case "smug":
		go handleAPI_smug(client, v, fullArgs)
	case "glomp":
		go handleAPI_glomp(client, v, fullArgs)
	case "happy":
		go handleAPI_happy(client, v, fullArgs)
	case "dance":
		go handleAPI_dance(client, v, fullArgs)
	case "cringe":
		go handleAPI_cringe(client, v, fullArgs)
	case "cuddle":
		go handleAPI_cuddle(client, v, fullArgs)
	case "highfive":
		go handleAPI_highfive(client, v, fullArgs)
	case "handhold":
		go handleAPI_handhold(client, v, fullArgs)
	case "shinobu":
		go handleAPI_shinobu(client, v, fullArgs)
	case "reactions":
		go handleAPI_reactions(client, v, fullArgs)
	case "rule34home":
		go handleAPI_rule34home(client, v, fullArgs)
	case "rule34search":
		go handleAPI_rule34search(client, v, fullArgs)
	case "rule34detail":
		go handleAPI_rule34detail(client, v, fullArgs)
	case "webnovelhot":
		go handleAPI_webnovelhot(client, v, fullArgs)
	case "webnovelrank":
		go handleAPI_webnovelrank(client, v, fullArgs)
	case "webnovelsearch":
		go handleAPI_webnovelsearch(client, v, fullArgs)
	case "webnoveldetail":
		go handleAPI_webnoveldetail(client, v, fullArgs)
	case "webnovelchapter":
		go handleAPI_webnovelchapter(client, v, fullArgs)
	case "aio":
		go handleAPI_aio(client, v, fullArgs)
	case "capcut":
		go handleAPI_capcut(client, v, fullArgs)
	case "doods":
		go handleAPI_doods(client, v, fullArgs)
	case "facebookv2":
		go handleAPI_facebookv2(client, v, fullArgs)
	case "ig2":
		go handleAPI_ig2(client, v, fullArgs)
	case "mediafire":
		go handleAPI_mediafire(client, v, fullArgs)
	case "pinterestv2":
		go handleAPI_pinterestv2(client, v, fullArgs)
	case "rednote":
		go handleAPI_rednote(client, v, fullArgs)
	case "rednotemedia":
		go handleAPI_rednotemedia(client, v, fullArgs)
	case "saveweb2zip":
		go handleAPI_saveweb2zip(client, v, fullArgs)
	case "sfile":
		go handleAPI_sfile(client, v, fullArgs)
	case "snackvideo":
		go handleAPI_snackvideo(client, v, fullArgs)
	case "spotifyv2":
		go handleAPI_spotifyv2(client, v, fullArgs)
	case "terabox":
		go handleAPI_terabox(client, v, fullArgs)
	case "threadsv2":
		go handleAPI_threadsv2(client, v, fullArgs)
	case "tiktokv2":
		go handleAPI_tiktokv2(client, v, fullArgs)
	case "tiktokv3":
		go handleAPI_tiktokv3(client, v, fullArgs)
	case "tiktokvideo":
		go handleAPI_tiktokvideo(client, v, fullArgs)
	case "tiktokslide":
		go handleAPI_tiktokslide(client, v, fullArgs)
	case "youtubeaudio":
		go handleAPI_youtubeaudio(client, v, fullArgs)
	case "youtubevideo":
		go handleAPI_youtubevideo(client, v, fullArgs)
	case "quizcategories":
		go handleAPI_quizcategories(client, v, fullArgs)
	case "quizguess":
		go handleAPI_quizguess(client, v, fullArgs)
	case "quizpuzzle":
		go handleAPI_quizpuzzle(client, v, fullArgs)
	case "quiztruefalse":
		go handleAPI_quiztruefalse(client, v, fullArgs)
	case "quizrandom":
		go handleAPI_quizrandom(client, v, fullArgs)
	case "gif":
		go handleAPI_gif(client, v, fullArgs)
	case "mp4":
		go handleAPI_mp4(client, v, fullArgs)
	case "meme":
		go handleAPI_meme(client, v, fullArgs)
	case "memetext":
		go handleAPI_memetext(client, v, fullArgs)
	case "spongebob":
		go handleAPI_spongebob(client, v, fullArgs)
	case "ttp":
		go handleAPI_ttp(client, v, fullArgs)
	case "moviesearch":
		go handleAPI_moviesearch(client, v, fullArgs)
	case "moviedetail":
		go handleAPI_moviedetail(client, v, fullArgs)
	case "suggest":
		go handleAPI_suggest(client, v, fullArgs)
	case "detail":
		go handleAPI_detail(client, v, fullArgs)
	case "recommendations":
		go handleAPI_recommendations(client, v, fullArgs)
	case "trending":
		go handleAPI_trending(client, v, fullArgs)
	case "home":
		go handleAPI_home(client, v, fullArgs)
	case "countries":
		go handleAPI_countries(client, v, fullArgs)
	case "hentaisfm":
		go handleAPI_hentaisfm(client, v, fullArgs)
	case "ass":
		go handleAPI_ass(client, v, fullArgs)
	case "sixtynine":
		go handleAPI_sixtynine(client, v, fullArgs)
	case "pussy":
		go handleAPI_pussy(client, v, fullArgs)
	case "dick":
		go handleAPI_dick(client, v, fullArgs)
	case "anal":
		go handleAPI_anal(client, v, fullArgs)
	case "boobs":
		go handleAPI_boobs(client, v, fullArgs)
	case "bdsm":
		go handleAPI_bdsm(client, v, fullArgs)
	case "black":
		go handleAPI_black(client, v, fullArgs)
	case "easter":
		go handleAPI_easter(client, v, fullArgs)
	case "bottomless":
		go handleAPI_bottomless(client, v, fullArgs)
	case "collared":
		go handleAPI_collared(client, v, fullArgs)
	case "cum":
		go handleAPI_cum(client, v, fullArgs)
	case "cumsluts":
		go handleAPI_cumsluts(client, v, fullArgs)
	case "dom":
		go handleAPI_dom(client, v, fullArgs)
	case "extreme":
		go handleAPI_extreme(client, v, fullArgs)
	case "feet":
		go handleAPI_feet(client, v, fullArgs)
	case "finger":
		go handleAPI_finger(client, v, fullArgs)
	case "fuck":
		go handleAPI_fuck(client, v, fullArgs)
	case "futa":
		go handleAPI_futa(client, v, fullArgs)
	case "gay":
		go handleAPI_gay(client, v, fullArgs)
	case "hentai":
		go handleAPI_hentai(client, v, fullArgs)
	case "kiss":
		go handleAPI_kiss(client, v, fullArgs)
	case "lick":
		go handleAPI_lick(client, v, fullArgs)
	case "pegged":
		go handleAPI_pegged(client, v, fullArgs)
	case "phgif":
		go handleAPI_phgif(client, v, fullArgs)
	case "puffies":
		go handleAPI_puffies(client, v, fullArgs)
	case "real":
		go handleAPI_real(client, v, fullArgs)
	case "suck":
		go handleAPI_suck(client, v, fullArgs)
	case "tattoo":
		go handleAPI_tattoo(client, v, fullArgs)
	case "tiny":
		go handleAPI_tiny(client, v, fullArgs)
	case "toys":
		go handleAPI_toys(client, v, fullArgs)
	case "xmas":
		go handleAPI_xmas(client, v, fullArgs)
	case "xnxxsearch":
		go handleAPI_xnxxsearch(client, v, fullArgs)
	case "xvideossearch":
		go handleAPI_xvideossearch(client, v, fullArgs)
	case "xnxxdl":
		go handleAPI_xnxxdl(client, v, fullArgs)
	case "xvideosdl":
		go handleAPI_xvideosdl(client, v, fullArgs)
	case "anhsfw":
		go handleAPI_anhsfw(client, v, fullArgs)
	case "anhmoe":
		go handleAPI_anhmoe(client, v, fullArgs)
	case "anhai":
		go handleAPI_anhai(client, v, fullArgs)
	case "anhnsfw":
		go handleAPI_anhnsfw(client, v, fullArgs)
	case "anhhentai":
		go handleAPI_anhhentai(client, v, fullArgs)
	case "anhvideonsfw":
		go handleAPI_anhvideonsfw(client, v, fullArgs)
	case "bluearchive":
		go handleAPI_bluearchive(client, v, fullArgs)
	case "boypic":
		go handleAPI_boypic(client, v, fullArgs)
	case "car":
		go handleAPI_car(client, v, fullArgs)
	case "cat":
		go handleAPI_cat(client, v, fullArgs)
	case "chinagirl":
		go handleAPI_chinagirl(client, v, fullArgs)
	case "dog":
		go handleAPI_dog(client, v, fullArgs)
	case "randomgirl":
		go handleAPI_randomgirl(client, v, fullArgs)
	case "hijabgirl":
		go handleAPI_hijabgirl(client, v, fullArgs)
	case "indonesiagirl":
		go handleAPI_indonesiagirl(client, v, fullArgs)
	case "japangirl":
		go handleAPI_japangirl(client, v, fullArgs)
	case "koreangirl":
		go handleAPI_koreangirl(client, v, fullArgs)
	case "loli":
		go handleAPI_loli(client, v, fullArgs)
	case "malaysiagirl":
		go handleAPI_malaysiagirl(client, v, fullArgs)
	case "profilepics":
		go handleAPI_profilepics(client, v, fullArgs)
	case "thailandgirl":
		go handleAPI_thailandgirl(client, v, fullArgs)
	case "tiktokgirl":
		go handleAPI_tiktokgirl(client, v, fullArgs)
	case "vietnamgirl":
		go handleAPI_vietnamgirl(client, v, fullArgs)
	case "waifu":
		go handleAPI_waifu(client, v, fullArgs)
	case "akiyama":
		go handleAPI_akiyama(client, v, fullArgs)
	case "ana":
		go handleAPI_ana(client, v, fullArgs)
	case "asuna":
		go handleAPI_asuna(client, v, fullArgs)
	case "ayuzawa":
		go handleAPI_ayuzawa(client, v, fullArgs)
	case "boruto":
		go handleAPI_boruto(client, v, fullArgs)
	case "chitoge":
		go handleAPI_chitoge(client, v, fullArgs)
	case "deidara":
		go handleAPI_deidara(client, v, fullArgs)
	case "doraemon":
		go handleAPI_doraemon(client, v, fullArgs)
	case "elaina":
		go handleAPI_elaina(client, v, fullArgs)
	case "emilia":
		go handleAPI_emilia(client, v, fullArgs)
	case "erza":
		go handleAPI_erza(client, v, fullArgs)
	case "hestia":
		go handleAPI_hestia(client, v, fullArgs)
	case "husbu":
		go handleAPI_husbu(client, v, fullArgs)
	case "inori":
		go handleAPI_inori(client, v, fullArgs)
	case "itachi":
		go handleAPI_itachi(client, v, fullArgs)
	case "kagura":
		go handleAPI_kagura(client, v, fullArgs)
	case "kaori":
		go handleAPI_kaori(client, v, fullArgs)
	case "keneki":
		go handleAPI_keneki(client, v, fullArgs)
	case "kotori":
		go handleAPI_kotori(client, v, fullArgs)
	case "kurumi":
		go handleAPI_kurumi(client, v, fullArgs)
	case "madara":
		go handleAPI_madara(client, v, fullArgs)
	case "megumin":
		go handleAPI_megumin(client, v, fullArgs)
	case "mikasa":
		go handleAPI_mikasa(client, v, fullArgs)
	case "miku":
		go handleAPI_miku(client, v, fullArgs)
	case "minato":
		go handleAPI_minato(client, v, fullArgs)
	case "naruto":
		go handleAPI_naruto(client, v, fullArgs)
	case "nekonime":
		go handleAPI_nekonime(client, v, fullArgs)
	case "nezuko":
		go handleAPI_nezuko(client, v, fullArgs)
	case "onepiece":
		go handleAPI_onepiece(client, v, fullArgs)
	case "rize":
		go handleAPI_rize(client, v, fullArgs)
	case "sagiri":
		go handleAPI_sagiri(client, v, fullArgs)
	case "sakura":
		go handleAPI_sakura(client, v, fullArgs)
	case "sasuke":
		go handleAPI_sasuke(client, v, fullArgs)
	case "shinomiya":
		go handleAPI_shinomiya(client, v, fullArgs)
	case "tsunade":
		go handleAPI_tsunade(client, v, fullArgs)
	case "yotsuba":
		go handleAPI_yotsuba(client, v, fullArgs)
	case "yuki":
		go handleAPI_yuki(client, v, fullArgs)
	case "yumeko":
		go handleAPI_yumeko(client, v, fullArgs)
	case "art":
		go handleAPI_art(client, v, fullArgs)
	case "cyber":
		go handleAPI_cyber(client, v, fullArgs)
	case "gamewallpaper":
		go handleAPI_gamewallpaper(client, v, fullArgs)
	case "mountain":
		go handleAPI_mountain(client, v, fullArgs)
	case "programming":
		go handleAPI_programming(client, v, fullArgs)
	case "space":
		go handleAPI_space(client, v, fullArgs)
	case "technology":
		go handleAPI_technology(client, v, fullArgs)
	case "wallhp":
		go handleAPI_wallhp(client, v, fullArgs)
	case "wallml":
		go handleAPI_wallml(client, v, fullArgs)
	case "wallmlnime":
		go handleAPI_wallmlnime(client, v, fullArgs)
	case "android1":
		go handleAPI_android1(client, v, fullArgs)
	case "cuaca":
		go handleAPI_cuaca(client, v, fullArgs)
	case "repos":
		go handleAPI_repos(client, v, fullArgs)
	case "users":
		go handleAPI_users(client, v, fullArgs)
	case "issues":
		go handleAPI_issues(client, v, fullArgs)
	case "code":
		go handleAPI_code(client, v, fullArgs)
	case "imdb":
		go handleAPI_imdb(client, v, fullArgs)
	case "lyrics":
		go handleAPI_lyrics(client, v, fullArgs)
	case "nik":
		go handleAPI_nik(client, v, fullArgs)
	case "wallpaper":
		go handleAPI_wallpaper(client, v, fullArgs)
	case "telegram":
		go handleAPI_telegram(client, v, fullArgs)
	case "tggroup":
		go handleAPI_tggroup(client, v, fullArgs)
	case "wagroup":
		go handleAPI_wagroup(client, v, fullArgs)
	case "ytmonet":
		go handleAPI_ytmonet(client, v, fullArgs)
	case "download":
		go handleAPI_download(client, v, fullArgs)
	case "nonstick":
		go handleAPI_nonstick(client, v, fullArgs)
	case "football":
		go handleAPI_football(client, v, fullArgs)
	case "basketball":
		go handleAPI_basketball(client, v, fullArgs)
	case "othersports":
		go handleAPI_othersports(client, v, fullArgs)
	case "webss":
		go handleAPI_webss(client, v, fullArgs)
	case "apiflash":
		go handleAPI_apiflash(client, v, fullArgs)
	case "screenshotlayer":
		go handleAPI_screenshotlayer(client, v, fullArgs)
	case "ffstalk":
		go handleAPI_ffstalk(client, v, fullArgs)
	case "igstalk":
		go handleAPI_igstalk(client, v, fullArgs)
	case "igstalkv2":
		go handleAPI_igstalkv2(client, v, fullArgs)
	case "ttstalk":
		go handleAPI_ttstalk(client, v, fullArgs)
	case "twitterstalk":
		go handleAPI_twitterstalk(client, v, fullArgs)
	case "ytstalk":
		go handleAPI_ytstalk(client, v, fullArgs)
	case "glitchtext":
		go handleAPI_glitchtext(client, v, fullArgs)
	case "writetext":
		go handleAPI_writetext(client, v, fullArgs)
	case "advancedglow":
		go handleAPI_advancedglow(client, v, fullArgs)
	case "typographytext":
		go handleAPI_typographytext(client, v, fullArgs)
	case "pixelglitch":
		go handleAPI_pixelglitch(client, v, fullArgs)
	case "neonglitch":
		go handleAPI_neonglitch(client, v, fullArgs)
	case "flagtext":
		go handleAPI_flagtext(client, v, fullArgs)
	case "flag3dtext":
		go handleAPI_flag3dtext(client, v, fullArgs)
	case "deletingtext":
		go handleAPI_deletingtext(client, v, fullArgs)
	case "blackpinkstyle":
		go handleAPI_blackpinkstyle(client, v, fullArgs)
	case "glowingtext":
		go handleAPI_glowingtext(client, v, fullArgs)
	case "underwatertext":
		go handleAPI_underwatertext(client, v, fullArgs)
	case "logomaker":
		go handleAPI_logomaker(client, v, fullArgs)
	case "cartoonstyle":
		go handleAPI_cartoonstyle(client, v, fullArgs)
	case "papercutstyle":
		go handleAPI_papercutstyle(client, v, fullArgs)
	case "watercolortext":
		go handleAPI_watercolortext(client, v, fullArgs)
	case "effectclouds":
		go handleAPI_effectclouds(client, v, fullArgs)
	case "blackpinklogo":
		go handleAPI_blackpinklogo(client, v, fullArgs)
	case "gradienttext":
		go handleAPI_gradienttext(client, v, fullArgs)
	case "summerbeach":
		go handleAPI_summerbeach(client, v, fullArgs)
	case "luxurygold":
		go handleAPI_luxurygold(client, v, fullArgs)
	case "multicoloredneon":
		go handleAPI_multicoloredneon(client, v, fullArgs)
	case "sandsummer":
		go handleAPI_sandsummer(client, v, fullArgs)
	case "galaxywallpaper":
		go handleAPI_galaxywallpaper(client, v, fullArgs)
	case "style1917":
		go handleAPI_style1917(client, v, fullArgs)
	case "makingneon":
		go handleAPI_makingneon(client, v, fullArgs)
	case "royaltext":
		go handleAPI_royaltext(client, v, fullArgs)
	case "freecreate":
		go handleAPI_freecreate(client, v, fullArgs)
	case "galaxystyle":
		go handleAPI_galaxystyle(client, v, fullArgs)
	case "lighteffects":
		go handleAPI_lighteffects(client, v, fullArgs)
	case "sendemail":
		go handleAPI_sendemail(client, v, fullArgs)
	case "codeanalyzer":
		go handleAPI_codeanalyzer(client, v, fullArgs)
	case "codeconverter":
		go handleAPI_codeconverter(client, v, fullArgs)
	case "tojavascript":
		go handleAPI_tojavascript(client, v, fullArgs)
	case "topython":
		go handleAPI_topython(client, v, fullArgs)
	case "tojava":
		go handleAPI_tojava(client, v, fullArgs)
	case "tocpp":
		go handleAPI_tocpp(client, v, fullArgs)
	case "tophp":
		go handleAPI_tophp(client, v, fullArgs)
	case "compiler":
		go handleAPI_compiler(client, v, fullArgs)
	case "compilejs":
		go handleAPI_compilejs(client, v, fullArgs)
	case "compilepython":
		go handleAPI_compilepython(client, v, fullArgs)
	case "compilejava":
		go handleAPI_compilejava(client, v, fullArgs)
	case "compilec":
		go handleAPI_compilec(client, v, fullArgs)
	case "compilecpp":
		go handleAPI_compilecpp(client, v, fullArgs)
	case "compilecsharp":
		go handleAPI_compilecsharp(client, v, fullArgs)
	case "emojiencrypt":
		go handleAPI_emojiencrypt(client, v, fullArgs)
	case "emojidecrypt":
		go handleAPI_emojidecrypt(client, v, fullArgs)
	case "htmlecnc":
		go handleAPI_htmlecnc(client, v, fullArgs)
	case "htmlbasic":
		go handleAPI_htmlbasic(client, v, fullArgs)
	case "htmlextended":
		go handleAPI_htmlextended(client, v, fullArgs)
	case "htmlhigh":
		go handleAPI_htmlhigh(client, v, fullArgs)
	case "htmlmaximum":
		go handleAPI_htmlmaximum(client, v, fullArgs)
	case "fdroidsearch":
		go handleAPI_fdroidsearch(client, v, fullArgs)
	case "fdroidpackage":
		go handleAPI_fdroidpackage(client, v, fullArgs)
	case "fdroidapp":
		go handleAPI_fdroidapp(client, v, fullArgs)
	case "geoip":
		go handleAPI_geoip(client, v, fullArgs)
	case "myip":
		go handleAPI_myip(client, v, fullArgs)
	case "hostcheck":
		go handleAPI_hostcheck(client, v, fullArgs)
	case "hostchecksimple":
		go handleAPI_hostchecksimple(client, v, fullArgs)
	case "html2img":
		go handleAPI_html2img(client, v, fullArgs)
	case "html2imgdirect":
		go handleAPI_html2imgdirect(client, v, fullArgs)
	case "obflow":
		go handleAPI_obflow(client, v, fullArgs)
	case "obfmedium":
		go handleAPI_obfmedium(client, v, fullArgs)
	case "obfhigh":
		go handleAPI_obfhigh(client, v, fullArgs)
	case "obfextreme":
		go handleAPI_obfextreme(client, v, fullArgs)
	case "tiktoktranscript":
		go handleAPI_tiktoktranscript(client, v, fullArgs)
	case "entoid":
		go handleAPI_entoid(client, v, fullArgs)
	case "idtoen":
		go handleAPI_idtoen(client, v, fullArgs)
	case "jatoid":
		go handleAPI_jatoid(client, v, fullArgs)
	case "kotoid":
		go handleAPI_kotoid(client, v, fullArgs)
	case "zhtoid":
		go handleAPI_zhtoid(client, v, fullArgs)
	case "artoid":
		go handleAPI_artoid(client, v, fullArgs)
	case "detectlanguage":
		go handleAPI_detectlanguage(client, v, fullArgs)
	case "languages":
		go handleAPI_languages(client, v, fullArgs)
	case "youtubetranscript":
		go handleAPI_youtubetranscript(client, v, fullArgs)
	case "dagd":
		go handleAPI_dagd(client, v, fullArgs)
	case "vgd":
		go handleAPI_vgd(client, v, fullArgs)
	case "tinube":
		go handleAPI_tinube(client, v, fullArgs)
	case "spoome":
		go handleAPI_spoome(client, v, fullArgs)
	case "spooemoji":
		go handleAPI_spooemoji(client, v, fullArgs)
	case "random":
		go handleAPI_random(client, v, fullArgs)
	case "allstyles":
		go handleAPI_allstyles(client, v, fullArgs)
	case "circled":
		go handleAPI_circled(client, v, fullArgs)
	case "circledneg":
		go handleAPI_circledneg(client, v, fullArgs)
	case "fullwidth":
		go handleAPI_fullwidth(client, v, fullArgs)
	case "mathbold":
		go handleAPI_mathbold(client, v, fullArgs)
	case "mathboldfraktur":
		go handleAPI_mathboldfraktur(client, v, fullArgs)
	case "mathbolditalic":
		go handleAPI_mathbolditalic(client, v, fullArgs)
	case "mathboldscript":
		go handleAPI_mathboldscript(client, v, fullArgs)
	case "mathdoublestruck":
		go handleAPI_mathdoublestruck(client, v, fullArgs)
	case "mathmonospace":
		go handleAPI_mathmonospace(client, v, fullArgs)
	case "mathsans":
		go handleAPI_mathsans(client, v, fullArgs)
	case "mathsansbold":
		go handleAPI_mathsansbold(client, v, fullArgs)
	case "mathsansbolditalic":
		go handleAPI_mathsansbolditalic(client, v, fullArgs)
	case "mathsansitalic":
		go handleAPI_mathsansitalic(client, v, fullArgs)
	case "parenthesized":
		go handleAPI_parenthesized(client, v, fullArgs)
	case "regionalindicator":
		go handleAPI_regionalindicator(client, v, fullArgs)
	case "squared":
		go handleAPI_squared(client, v, fullArgs)
	case "squaredneg":
		go handleAPI_squaredneg(client, v, fullArgs)
	case "tag":
		go handleAPI_tag(client, v, fullArgs)
	case "acute":
		go handleAPI_acute(client, v, fullArgs)
	case "cjkthai":
		go handleAPI_cjkthai(client, v, fullArgs)
	case "curvy1":
		go handleAPI_curvy1(client, v, fullArgs)
	case "curvy2":
		go handleAPI_curvy2(client, v, fullArgs)
	case "curvy3":
		go handleAPI_curvy3(client, v, fullArgs)
	case "fauxcyrillic":
		go handleAPI_fauxcyrillic(client, v, fullArgs)
	case "fauxethiopic":
		go handleAPI_fauxethiopic(client, v, fullArgs)
	case "mathfraktur":
		go handleAPI_mathfraktur(client, v, fullArgs)
	case "rockdots":
		go handleAPI_rockdots(client, v, fullArgs)
	case "smallcaps":
		go handleAPI_smallcaps(client, v, fullArgs)
	case "stroked":
		go handleAPI_stroked(client, v, fullArgs)
	case "subscript":
		go handleAPI_subscript(client, v, fullArgs)
	case "superscript":
		go handleAPI_superscript(client, v, fullArgs)
	case "inverted":
		go handleAPI_inverted(client, v, fullArgs)
	case "reversed":
		go handleAPI_reversed(client, v, fullArgs)
	case "ttsen":
		go handleAPI_ttsen(client, v, fullArgs)
	case "ttsid":
		go handleAPI_ttsid(client, v, fullArgs)
	case "ttses":
		go handleAPI_ttses(client, v, fullArgs)
	case "ttsfr":
		go handleAPI_ttsfr(client, v, fullArgs)
	case "ttsde":
		go handleAPI_ttsde(client, v, fullArgs)
	case "ttsit":
		go handleAPI_ttsit(client, v, fullArgs)
	case "ttspt":
		go handleAPI_ttspt(client, v, fullArgs)
	case "ttsnl":
		go handleAPI_ttsnl(client, v, fullArgs)
	case "ttspl":
		go handleAPI_ttspl(client, v, fullArgs)
	case "ttsru":
		go handleAPI_ttsru(client, v, fullArgs)
	case "ttsja":
		go handleAPI_ttsja(client, v, fullArgs)
	case "ttsko":
		go handleAPI_ttsko(client, v, fullArgs)
	case "ttszhcn":
		go handleAPI_ttszhcn(client, v, fullArgs)
	case "ttszhtw":
		go handleAPI_ttszhtw(client, v, fullArgs)
	case "ttsar":
		go handleAPI_ttsar(client, v, fullArgs)
	case "ttshi":
		go handleAPI_ttshi(client, v, fullArgs)
	case "ttsth":
		go handleAPI_ttsth(client, v, fullArgs)
	case "ttsvi":
		go handleAPI_ttsvi(client, v, fullArgs)
	case "ttstr":
		go handleAPI_ttstr(client, v, fullArgs)
	case "ttssv":
		go handleAPI_ttssv(client, v, fullArgs)
	case "ttsno":
		go handleAPI_ttsno(client, v, fullArgs)
	case "ttsda":
		go handleAPI_ttsda(client, v, fullArgs)
	case "ttsfi":
		go handleAPI_ttsfi(client, v, fullArgs)
	case "ttsel":
		go handleAPI_ttsel(client, v, fullArgs)
	case "ttshe":
		go handleAPI_ttshe(client, v, fullArgs)
	case "ttscs":
		go handleAPI_ttscs(client, v, fullArgs)
	case "ttshu":
		go handleAPI_ttshu(client, v, fullArgs)
	case "ttsro":
		go handleAPI_ttsro(client, v, fullArgs)
	case "ttsuk":
		go handleAPI_ttsuk(client, v, fullArgs)
	case "xiaofei":
		go handleAPI_xiaofei(client, v, fullArgs)
	case "xiaolei":
		go handleAPI_xiaolei(client, v, fullArgs)
	case "xiaojie":
		go handleAPI_xiaojie(client, v, fullArgs)
	case "xiaohua":
		go handleAPI_xiaohua(client, v, fullArgs)
	case "xiaofeng":
		go handleAPI_xiaofeng(client, v, fullArgs)
	case "xiaoze":
		go handleAPI_xiaoze(client, v, fullArgs)
	case "xiaoyuan":
		go handleAPI_xiaoyuan(client, v, fullArgs)
	case "xiaozheng":
		go handleAPI_xiaozheng(client, v, fullArgs)
	case "xiaoying":
		go handleAPI_xiaoying(client, v, fullArgs)
	case "xiaoqing":
		go handleAPI_xiaoqing(client, v, fullArgs)
	case "xiaoxiang":
		go handleAPI_xiaoxiang(client, v, fullArgs)
	case "xiaoyan":
		go handleAPI_xiaoyan(client, v, fullArgs)
	case "xianran":
		go handleAPI_xianran(client, v, fullArgs)
	case "xiaoxue":
		go handleAPI_xiaoxue(client, v, fullArgs)
	case "xiaoxuan":
		go handleAPI_xiaoxuan(client, v, fullArgs)
	case "xiaolu":
		go handleAPI_xiaolu(client, v, fullArgs)
	case "xiaowei":
		go handleAPI_xiaowei(client, v, fullArgs)
	case "xiaozhe":
		go handleAPI_xiaozhe(client, v, fullArgs)
	case "xiaohao":
		go handleAPI_xiaohao(client, v, fullArgs)
	case "xiaoyi":
		go handleAPI_xiaoyi(client, v, fullArgs)
	case "xiaotao":
		go handleAPI_xiaotao(client, v, fullArgs)
	case "xiaoming":
		go handleAPI_xiaoming(client, v, fullArgs)
	case "david":
		go handleAPI_david(client, v, fullArgs)
	case "layla":
		go handleAPI_layla(client, v, fullArgs)
	case "james":
		go handleAPI_james(client, v, fullArgs)
	case "joey":
		go handleAPI_joey(client, v, fullArgs)
	case "jennifer":
		go handleAPI_jennifer(client, v, fullArgs)
	case "john":
		go handleAPI_john(client, v, fullArgs)
	case "paul":
		go handleAPI_paul(client, v, fullArgs)
	case "xena":
		go handleAPI_xena(client, v, fullArgs)
	case "marcus":
		go handleAPI_marcus(client, v, fullArgs)
	case "jacob":
		go handleAPI_jacob(client, v, fullArgs)
	case "sam":
		go handleAPI_sam(client, v, fullArgs)
	case "camila":
		go handleAPI_camila(client, v, fullArgs)
	case "amy":
		go handleAPI_amy(client, v, fullArgs)
	case "quincy":
		go handleAPI_quincy(client, v, fullArgs)
	case "sally":
		go handleAPI_sally(client, v, fullArgs)
	case "emma":
		go handleAPI_emma(client, v, fullArgs)
	case "ethan":
		go handleAPI_ethan(client, v, fullArgs)
	case "michael":
		go handleAPI_michael(client, v, fullArgs)
	case "olivia":
		go handleAPI_olivia(client, v, fullArgs)
	case "mia":
		go handleAPI_mia(client, v, fullArgs)
	case "jackson":
		go handleAPI_jackson(client, v, fullArgs)
	case "matthew":
		go handleAPI_matthew(client, v, fullArgs)
	case "sophia":
		go handleAPI_sophia(client, v, fullArgs)
	case "owen":
		go handleAPI_owen(client, v, fullArgs)
	case "beatrice":
		go handleAPI_beatrice(client, v, fullArgs)
	case "scott":
		go handleAPI_scott(client, v, fullArgs)
	case "ivy":
		go handleAPI_ivy(client, v, fullArgs)
	case "eric":
		go handleAPI_eric(client, v, fullArgs)
	case "kevin":
		go handleAPI_kevin(client, v, fullArgs)
	case "hannah":
		go handleAPI_hannah(client, v, fullArgs)
	case "katrina":
		go handleAPI_katrina(client, v, fullArgs)
	case "victor":
		go handleAPI_victor(client, v, fullArgs)
	case "justin":
		go handleAPI_justin(client, v, fullArgs)
	case "leo":
		go handleAPI_leo(client, v, fullArgs)
	case "grace":
		go handleAPI_grace(client, v, fullArgs)
	case "casey":
		go handleAPI_casey(client, v, fullArgs)
	case "dylan":
		go handleAPI_dylan(client, v, fullArgs)
	case "julie":
		go handleAPI_julie(client, v, fullArgs)
	case "thomas":
		go handleAPI_thomas(client, v, fullArgs)
	case "freya":
		go handleAPI_freya(client, v, fullArgs)
	case "max":
		go handleAPI_max(client, v, fullArgs)
	case "phoebe":
		go handleAPI_phoebe(client, v, fullArgs)
	case "noah":
		go handleAPI_noah(client, v, fullArgs)
	case "sophie":
		go handleAPI_sophie(client, v, fullArgs)
	case "isla":
		go handleAPI_isla(client, v, fullArgs)
	case "theo":
		go handleAPI_theo(client, v, fullArgs)
	case "ella":
		go handleAPI_ella(client, v, fullArgs)
	case "freddie":
		go handleAPI_freddie(client, v, fullArgs)
	case "arthur":
		go handleAPI_arthur(client, v, fullArgs)
	case "isabella":
		go handleAPI_isabella(client, v, fullArgs)
	case "evie":
		go handleAPI_evie(client, v, fullArgs)
	case "william":
		go handleAPI_william(client, v, fullArgs)
	case "henry":
		go handleAPI_henry(client, v, fullArgs)
	case "lily":
		go handleAPI_lily(client, v, fullArgs)
	case "charlie":
		go handleAPI_charlie(client, v, fullArgs)
	case "ttsvoices":
		go handleAPI_ttsvoices(client, v, fullArgs)
	case "ttsadultfemale1americanenglishtruvoice":
		go handleAPI_ttsadultfemale1americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultfemale2americanenglishtruvoice":
		go handleAPI_ttsadultfemale2americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale1americanenglishtruvoice":
		go handleAPI_ttsadultmale1americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale2americanenglishtruvoice":
		go handleAPI_ttsadultmale2americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale3americanenglishtruvoice":
		go handleAPI_ttsadultmale3americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale4americanenglishtruvoice":
		go handleAPI_ttsadultmale4americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale5americanenglishtruvoice":
		go handleAPI_ttsadultmale5americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale6americanenglishtruvoice":
		go handleAPI_ttsadultmale6americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale7americanenglishtruvoice":
		go handleAPI_ttsadultmale7americanenglishtruvoice(client, v, fullArgs)
	case "ttsadultmale8americanenglishtruvoice":
		go handleAPI_ttsadultmale8americanenglishtruvoice(client, v, fullArgs)
	case "ttsfemalewhisper":
		go handleAPI_ttsfemalewhisper(client, v, fullArgs)
	case "ttsmalewhisper":
		go handleAPI_ttsmalewhisper(client, v, fullArgs)
	case "ttsmary":
		go handleAPI_ttsmary(client, v, fullArgs)
	case "ttsmaryfortelephone":
		go handleAPI_ttsmaryfortelephone(client, v, fullArgs)
	case "ttsmaryinhall":
		go handleAPI_ttsmaryinhall(client, v, fullArgs)
	case "ttsmaryinspace":
		go handleAPI_ttsmaryinspace(client, v, fullArgs)
	case "ttsmaryinstadium":
		go handleAPI_ttsmaryinstadium(client, v, fullArgs)
	case "ttsmike":
		go handleAPI_ttsmike(client, v, fullArgs)
	case "ttsmikefortelephone":
		go handleAPI_ttsmikefortelephone(client, v, fullArgs)
	case "ttsmikeinhall":
		go handleAPI_ttsmikeinhall(client, v, fullArgs)
	case "ttsmikeinspace":
		go handleAPI_ttsmikeinspace(client, v, fullArgs)
	case "ttsmikeinstadium":
		go handleAPI_ttsmikeinstadium(client, v, fullArgs)
	case "ttsrobosoftfive":
		go handleAPI_ttsrobosoftfive(client, v, fullArgs)
	case "ttsrobosoftfour":
		go handleAPI_ttsrobosoftfour(client, v, fullArgs)
	case "ttsrobosoftone":
		go handleAPI_ttsrobosoftone(client, v, fullArgs)
	case "ttsrobosoftsix":
		go handleAPI_ttsrobosoftsix(client, v, fullArgs)
	case "ttsrobosoftthree":
		go handleAPI_ttsrobosoftthree(client, v, fullArgs)
	case "ttsrobosofttwo":
		go handleAPI_ttsrobosofttwo(client, v, fullArgs)
	case "ttssam":
		go handleAPI_ttssam(client, v, fullArgs)
	case "ttsbonzi":
		go handleAPI_ttsbonzi(client, v, fullArgs)
	case "sms24countries":
		go handleAPI_sms24countries(client, v, fullArgs)
	case "sms24numbers":
		go handleAPI_sms24numbers(client, v, fullArgs)
	case "sms24messages":
		go handleAPI_sms24messages(client, v, fullArgs)
	case "veepncountries":
		go handleAPI_veepncountries(client, v, fullArgs)
	case "veepnnumbers":
		go handleAPI_veepnnumbers(client, v, fullArgs)
	case "veepnmessages":
		go handleAPI_veepnmessages(client, v, fullArgs)

	case "gcname":
		go handleGCName(client, v, fullArgs)
	case "gcdp":
		go handleGCDP(client, v)
	case "gclink", "link":
		go handleGCLink(client, v)
	case "revoke", "resetlink":
		go handleGCRevoke(client, v)
	case "block":
		go handleBlock(ctx, client, v)
	case "getstatus":
		react(client, v.Info.Chat, v.Info.ID, "🔍")
		go handleGetStatus(ctx, client, v, fullArgs, settingsDB) 
	case "addstatus":
		react(client, v.Info.Chat, v.Info.ID, "📲")
		go handleAddStatus(client, v)
	case "dp":
		go handleDP(ctx, client, v, fullArgs)
	case "changename":
		go handleChangeName(ctx, client, v, fullArgs)
	case "setabout":
    	handleSetAbout(ctx, client, v, fullArgs)

    case "checkwa":
    	handleCheckWA(ctx, client, v, fullArgs)

    case "getdp":
    	handleGetDP(ctx, client, v, fullArgs)

    case "blocklist":
    	handleBlocklist(ctx, client, v)
	
	}
}

func react(client *whatsmeow.Client, chat types.JID, msgID types.MessageID, emoji string) {
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Printf("⚠️ React Panic: %v\n", r)
			}
		}()

		_, err := client.SendMessage(context.Background(), chat, &waProto.Message{
			ReactionMessage: &waProto.ReactionMessage{
				Key: &waProto.MessageKey{
					RemoteJID: proto.String(chat.String()),
					ID:        proto.String(string(msgID)),
					FromMe:    proto.Bool(false),
				},
				Text:              proto.String(emoji),
				SenderTimestampMS: proto.Int64(time.Now().UnixMilli()),
			},
		})

		if err != nil {
			fmt.Printf("❌ React Failed: %v\n", err)
		}
	}()
}

func replyMessage(client *whatsmeow.Client, v *events.Message, text string) string {
	resp, err := client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
		ExtendedTextMessage: &waProto.ExtendedTextMessage{
			Text: proto.String(text),
			ContextInfo: &waProto.ContextInfo{
				StanzaID:      proto.String(v.Info.ID),
				Participant:   proto.String(v.Info.Sender.String()),
				QuotedMessage: v.Message,
			},
		},
	})
	if err == nil {
		return resp.ID
	}
	return ""
}

func handlePair(client *whatsmeow.Client, v *events.Message, args string) {
	if args == "" {
		replyMessage(client, v, "❌ Please provide a phone number with country code.\nExample: `.pair 923001234567`")
		return
	}

	phone := strings.ReplaceAll(args, "+", "")
	phone = strings.ReplaceAll(phone, " ", "")
	phone = strings.ReplaceAll(phone, "-", "")

	react(client, v.Info.Chat, v.Info.ID, "⏳")
	replyMessage(client, v, "⏳ Generating pairing code... Please wait.")

	deviceStore := dbContainer.NewDevice()
	clientLog := waLog.Noop
	newClient := whatsmeow.NewClient(deviceStore, clientLog)

	newClient.AddEventHandler(func(evt interface{}) {
		EventHandler(newClient, evt)
	})

	err := newClient.Connect()
	if err != nil {
		replyMessage(client, v, "❌ Failed to connect to WhatsApp servers.")
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	code, err := newClient.PairPhone(context.Background(), phone, true, whatsmeow.PairClientChrome, "Chrome (Linux)")
	if err != nil {
		replyMessage(client, v, fmt.Sprintf("❌ Failed to get pairing code: %v", err))
		react(client, v.Info.Chat, v.Info.ID, "❌")
		return
	}

	formattedCode := code
	if len(code) == 8 {
		formattedCode = code[:4] + "-" + code[4:]
	}

	successMsg := fmt.Sprintf("✅ *PAIRING CODE GENERATED*\n\n📱 *Phone:* +%s\n\n_1. Open WhatsApp on target phone_\n_2. Go to Linked Devices -> Link a Device_\n_3. Select 'Link with phone number instead'_\n_4. Enter the code below_ 👇\n\n⚠️ _This code expires in 2 minutes._", phone)
	replyMessage(client, v, successMsg)
	replyMessage(client, v, formattedCode)
	react(client, v.Info.Chat, v.Info.ID, "✅")
}

func handleID(client *whatsmeow.Client, v *events.Message) {
	chatJID := v.Info.Chat.String()
	senderJID := v.Info.Sender.ToNonAD().String()
	senderLID := v.Info.Sender.String()

	if !v.Info.SenderAlt.IsEmpty() {
	    senderLID = v.Info.SenderAlt.String()
	}

	chatType := "👤 𝗣𝗿𝗶𝘃𝗮𝘁ε 𝗖𝗵𝗮𝘁"
	if strings.Contains(chatJID, "@g.us") {
		chatType = "👥 𝗚𝗿𝗼𝘂𝗽 𝗖𝗵𝗮𝘁"
	}

	card := fmt.Sprintf(`❖ ── ✦ 🪪 𝗜𝗗 𝗖𝗔𝗥Ｄ ✦ ── ❖

 %s
 ➭ *%s*

 👤 𝗦𝗲𝗻𝗱𝗲𝗿 (JID)
 ➭ *%s*

 👤 𝗦𝗲𝗻𝗱𝗲𝗿 (LID)
 ➭ *%s*`, chatType, chatJID, senderJID, senderLID)

	extMsg := v.Message.GetExtendedTextMessage()
	if extMsg != nil && extMsg.ContextInfo != nil && extMsg.ContextInfo.Participant != nil {
		quotedJID := *extMsg.ContextInfo.Participant
		card += fmt.Sprintf("\n\n 🎯 𝗧𝗮𝗿𝗴𝗲𝘁 (𝗤𝘂𝗼𝘁𝗲𝗱)\n ➭ *%s*", quotedJID)
	}

	card += "\n\n ╰──────────────────────╯"
	replyMessage(client, v, card)
}

func handleAntiCallLogic(client *whatsmeow.Client, c *events.CallOffer, settings BotSettings) {
	if c.CallCreator.Server == "g.us" || c.CallCreator.Server == types.GroupServer {
		return
	}

	botJID := client.Store.ID.ToNonAD().User
	callerJID := c.CallCreator.ToNonAD()

	isCallEnabled := settings.AntiCall
	var dbCheck bool
	errDB := settingsDB.QueryRow("SELECT anti_call FROM bot_settings WHERE jid = ?", botJID).Scan(&dbCheck)
	if errDB == nil && dbCheck {
		isCallEnabled = true
	}

	if !isCallEnabled || callerJID.User == botJID {
		return
	}

	contact, err := client.Store.Contacts.GetContact(context.Background(), callerJID)
	isSaved := (err == nil && contact.Found && contact.FullName != "")

	if !isSaved {
		fmt.Printf("📞 [ANTI-CALL] Triggered! Dropping call from Unsaved Number: %s\n", callerJID.User)
		client.RejectCall(context.Background(), c.CallCreator, c.CallID)
		client.RejectCall(context.Background(), callerJID, c.CallID)
	}
}

func handleAntiDMWatch(client *whatsmeow.Client, v *events.Message, settings BotSettings) bool {
	botJID := client.Store.ID.ToNonAD().User

	isEnabled := settings.AntiDM
	var dbCheck bool
	errDB := settingsDB.QueryRow("SELECT anti_dm FROM bot_settings WHERE jid = ?", botJID).Scan(&dbCheck)
	if errDB == nil && dbCheck {
		isEnabled = true
	}

	if !isEnabled || v.Info.IsGroup || v.Info.IsFromMe || v.Info.Chat.Server == "newsletter" || v.Info.Chat.Server == types.NewsletterServer || isOwner(client, v) {
		return false
	}

	var realSender types.JID
	if v.Info.Sender.Server == types.HiddenUserServer {
		if !v.Info.SenderAlt.IsEmpty() {
			realSender = v.Info.SenderAlt.ToNonAD()
		} else {
			realSender = v.Info.Sender.ToNonAD()
		}
	} else {
		realSender = v.Info.Sender.ToNonAD()
	}

	contact, err := client.Store.Contacts.GetContact(context.Background(), realSender)
	isSaved := err == nil && contact.Found && contact.FullName != ""

	if !isSaved {
		fmt.Printf("🛡️ [ANTI-DM] TRIGGERED [Bot: %s]: Unsaved number -> %s\n", botJID, realSender.User)

		warning := "⚠️ *Silent Nexus Security*\n\nDirect messages from unsaved numbers are not allowed. You are being blocked automatically."
		client.SendMessage(context.Background(), v.Info.Chat, &waProto.Message{
			Conversation: proto.String(warning),
		})

		time.Sleep(2 * time.Second)

		_, errBlock1 := client.UpdateBlocklist(context.Background(), v.Info.Sender.ToNonAD(), events.BlocklistChangeActionBlock)
		if errBlock1 != nil {
			_, errBlock2 := client.UpdateBlocklist(context.Background(), realSender, events.BlocklistChangeActionBlock)
			if errBlock2 == nil {
				fmt.Printf("✅ [ANTI-DM] Successfully blocked real number: %s\n", realSender.String())
			} else {
				fmt.Printf("❌ [ANTI-DM ERROR] Block failed: %v\n", errBlock2)
			}
		} else {
			fmt.Printf("✅ [ANTI-DM] Successfully blocked LID: %s\n", v.Info.Sender.String())
		}

		time.Sleep(1 * time.Second)

		lastMessageKey := &waCommon.MessageKey{
			RemoteJID: proto.String(v.Info.Chat.String()),
			FromMe:    proto.Bool(v.Info.IsFromMe),
			ID:        proto.String(v.Info.ID),
		}

		patchInfo1 := appstate.BuildDeleteChat(v.Info.Chat, v.Info.Timestamp, lastMessageKey, true)
		errPatch1 := client.SendAppState(context.Background(), patchInfo1)

		patchInfo2 := appstate.BuildDeleteChat(realSender, v.Info.Timestamp, nil, true)
		errPatch2 := client.SendAppState(context.Background(), patchInfo2)

		if errPatch1 == nil || errPatch2 == nil {
			fmt.Printf("✅ [ANTI-DM] Chat DELETED from WhatsApp screen for: %s\n", realSender.User)
		} else {
			fmt.Printf("❌ [ANTI-DM ERROR] Delete failed. Patch1: %v | Patch2: %v\n", errPatch1, errPatch2)
		}

		return true
	}

	return false
}