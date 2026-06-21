package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/binary"
	waProto "go.mau.fi/whatsmeow/proto/waE2E"
	"go.mau.fi/whatsmeow/types/events"
	"google.golang.org/protobuf/proto"
)

var SendLogo = true          
var cachedLogoResp *whatsmeow.UploadResponse 



func getCommandCount() string {
	content, err := os.ReadFile("3.go")
	if err != nil {
		return "500+"
	}

	count := 0
	lines := strings.Split(string(content), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "case \"") && strings.HasSuffix(trimmed, ":") {
			count += strings.Count(trimmed, ",") + 1
		}
	}
	return fmt.Sprintf("%d", count)
}


func sendMainMenu(client *whatsmeow.Client, v *events.Message, settings BotSettings) {
	uptimeStr := getUptimeString(settings.UptimeStart)

	headerText := fmt.Sprintf(`┏━━━〔 👑 👑 👑 %s 👑 👑 👑 〕━━━┈
┃ 👤 *Owner:* %s
┃ ⚙️ *Mode:* %s
┃ ⏱️ *Uptime:* %s
┃ ⚡ *Prefix:* [ %s ]
┃ 📊 *Commands:* %s
┗━━━━━━━━━━━━━━━━━━┈`, GlobalConfig.BotName, GlobalConfig.Developer, strings.ToUpper(settings.Mode), uptimeStr, settings.Prefix, getCommandCount())
	
	mainMenuBody := headerText + "\n\n" + `┍──╼〔 📋 *MAIN MENU LIST* 〕
┗━━━━━━━━━━━━━━━━━━┈`

	
	listParams := map[string]interface{}{
		"title": "Open Menu List",
		"sections": []map[string]interface{}{
			{
				"title": "Select a Menu",
				"rows": []map[string]interface{}{
					{"id": settings.Prefix + "ytmenu", "title": "📺 YOUTUBE MENU", "description": "YouTube Downloader Commands"},
					{"id": settings.Prefix + "ttmenu", "title": "📱 TIKTOK MENU", "description": "TikTok Downloader Commands"},
					{"id": settings.Prefix + "dlmenu", "title": "🌐 DOWNLOAD MENU", "description": "Social Media Download Menu"},
					{"id": settings.Prefix + "aimenu", "title": "🧠 AI CHAT MENU", "description": "Artificial Intelligence Commands"},
					{"id": settings.Prefix + "gpmenu", "title": "🛡️ GROUP MENU", "description": "Group Management Commands"},
					{"id": settings.Prefix + "ownermenu", "title": "⚙️ OWNER MENU", "description": "Bot Owner Settings"},
					{"id": settings.Prefix + "utilmenu", "title": "🛠️ UTILITY TOOLS", "description": "Utility Tools and Commands"},
					{"id": settings.Prefix + "editmenu", "title": "🎨 EDITING ZONE", "description": "Sticker and Media Editing"},
					{"id": settings.Prefix + "aitools", "title": "✨ AI TOOLS ZONE", "description": "Advanced Image and Translation Tools"},
					{"id": settings.Prefix + "artificialintelligencemenu", "title": "🤖 ARTIFICIAL INTELLIGENCE", "description": "61 Commands"},
					{"id": settings.Prefix + "imagegenerationmenu", "title": "🖼️ IMAGE GENERATION", "description": "16 Commands"},
					{"id": settings.Prefix + "animemenu", "title": "🌸 ANIME", "description": "54 Commands"},
					{"id": settings.Prefix + "gamesmenu", "title": "🎮 GAMES", "description": "5 Commands"},
					{"id": settings.Prefix + "imagecreatormenu", "title": "🖌️ IMAGE CREATOR", "description": "7 Commands"},
					{"id": settings.Prefix + "moviesmenu", "title": "🎬 MOVIES", "description": "9 Commands"},
					{"id": settings.Prefix + "searchmenu", "title": "🔍 SEARCH", "description": "17 Commands"},
					{"id": settings.Prefix + "randommenu", "title": "🎲 RANDOM", "description": "48 Commands"},
					{"id": settings.Prefix + "audiomenu", "title": "🎵 AUDIO", "description": "2 Commands"},
					{"id": settings.Prefix + "sportsmenu", "title": "⚽ SPORTS", "description": "3 Commands"},
					{"id": settings.Prefix + "screenshotwebsitemenu", "title": "🖥️ SCREENSHOT WEBSITE", "description": "3 Commands"},
					{"id": settings.Prefix + "stalkmenu", "title": "🕵️ STALK", "description": "6 Commands"},
					{"id": settings.Prefix + "textmakermenu", "title": "📝 TEXT MAKER", "description": "30 Commands"},
					{"id": settings.Prefix + "toolsmenu", "title": "🧰 TOOLS", "description": "45 Commands"},
					{"id": settings.Prefix + "urlshortnermenu", "title": "🔗 URL SHORTNER", "description": "6 Commands"},
					{"id": settings.Prefix + "styletextmenu", "title": "🔤 STYLETEXT", "description": "34 Commands"},
					{"id": settings.Prefix + "texttospeechmenu", "title": "🗣️ TEXT TO SPEECH", "description": "137 Commands"},
					{"id": settings.Prefix + "virtualnumbermenu", "title": "🔢 VIRTUAL NUMBER", "description": "6 Commands"},
				},
			},
		},
	}
	listParamsJSON, _ := json.Marshal(listParams)

	
	urlBtn1Params := map[string]string{
		"display_text": "📢 Join Channel",
		"url":          GlobalConfig.ChannelLink,
		"merchant_url": GlobalConfig.ChannelLink,
	}
	urlBtn1JSON, _ := json.Marshal(urlBtn1Params)

	buttons := []*waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
		{
			Name:             proto.String("single_select"),
			ButtonParamsJSON: proto.String(string(listParamsJSON)),
		},
		{
			Name:             proto.String("cta_url"),
			ButtonParamsJSON: proto.String(string(urlBtn1JSON)),
		},
	}

	interactiveMsg := &waProto.InteractiveMessage{
		Body: &waProto.InteractiveMessage_Body{
			Text: proto.String(mainMenuBody),
		},
		Footer: &waProto.InteractiveMessage_Footer{
			Text: proto.String(fmt.Sprintf("Powered by %s Engine", GlobalConfig.BotName)),
		},
		InteractiveMessage: &waProto.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: &waProto.InteractiveMessage_NativeFlowMessage{
				Buttons: buttons,
			},
		},
	}

		
	if SendLogo {
		if cachedLogoResp == nil {
			imageBytes, err := os.ReadFile("logo.png")
			if err == nil {
				resp, err := client.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
				if err == nil {
					cachedLogoResp = &resp
				}
			}
		}

		if cachedLogoResp != nil {
			interactiveMsg.Header = &waProto.InteractiveMessage_Header{
				HasMediaAttachment: proto.Bool(true),
				Media: &waProto.InteractiveMessage_Header_ImageMessage{
					ImageMessage: &waProto.ImageMessage{
						Mimetype:          proto.String("image/png"),
						URL:               proto.String(cachedLogoResp.URL),
						DirectPath:        proto.String(cachedLogoResp.DirectPath),
						MediaKey:          cachedLogoResp.MediaKey,
						FileEncSHA256:     cachedLogoResp.FileEncSHA256,
						FileSHA256:        cachedLogoResp.FileSHA256,
						FileLength:        proto.Uint64(uint64(cachedLogoResp.FileLength)),
						MediaKeyTimestamp: proto.Int64(time.Now().Unix()), 
					},
				},
			}
		}
	}


	if interactiveMsg.Header == nil {
		interactiveMsg.Header = &waProto.InteractiveMessage_Header{
			Title:              proto.String(fmt.Sprintf("👑 %s 👑", strings.ToUpper(GlobalConfig.BotName))),
			HasMediaAttachment: proto.Bool(false),
		}
	}

	msg := &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: &waProto.Message{
				InteractiveMessage: interactiveMsg,
			},
		},
	}

	bizNode := binary.Node{
		Tag: "biz",
		Content: []binary.Node{
			{
				Tag: "interactive",
				Attrs: map[string]interface{}{"type": "native_flow", "v": "1"},
				Content: []binary.Node{
					{
						Tag:   "native_flow",
						Attrs: map[string]interface{}{"v": "9", "name": "mixed"},
					},
				},
			},
		},
	}
	nodes := []binary.Node{bizNode}
	client.SendMessage(context.Background(), v.Info.Chat, msg, whatsmeow.SendRequestExtra{AdditionalNodes: &nodes})
}


func sendSubMenu(client *whatsmeow.Client, v *events.Message, menuType string, settings BotSettings) {
	uptimeStr := getUptimeString(settings.UptimeStart)

	headerText := fmt.Sprintf(`┏━━━〔 👑 %s 👑 〕━━━┈
┃ 👤 *Owner:* %s
┃ ⚙️ *Mode:* %s
┃ ⏱️ *Uptime:* %s
┃ ⚡ *Prefix:* [ %s ]
┃ 📊 *Commands:* %s
┗━━━━━━━━━━━━━━━━━━┈`, strings.ToUpper(GlobalConfig.BotName), GlobalConfig.Developer, strings.ToUpper(settings.Mode), uptimeStr, settings.Prefix, getCommandCount())

	p := settings.Prefix
	var bodyText string

	switch strings.ToLower(menuType) {
	case "ytmenu":
		bodyText = `
 ╭── ✦ [ *YOUTUBE MENU* ] ✦ ──╮
 │
 │ ➭ *[p]play* / *[p]song* [name]
 │    _Direct HQ Audio Download_
 │
 │ ➭ *[p]video* [name]
 │    _Direct HD Video Download_
 │
 │ ➭ *[p]yt* [link]
 │    _Download YT Video/Audio_
 │
 │ ➭ *[p]yts* [query]
 │    _Search YouTube Videos_
 │
 ╰──────────────────────╯`

	case "ttmenu":
		bodyText = `
 ╭── ✦ [ *TIKTOK MENU* ] ✦ ──╮
 │
 │ ➭ *[p]tt* / *[p]tiktok* [link]
 │    _Download TikTok Video (No WM)_
 │
 │ ➭ *[p]tts* [query]
 │    _Search TikTok Videos_
 │
 ╰──────────────────────╯`

	case "dlmenu":
		bodyText = `
 ╭── ✦ [ *DOWNLOAD MENU* ] ✦ ──╮
 │
 │ ➭ *[p]fb* / *[p]facebook* [link]
 │    _Download Facebook Videos_
 │
 │ ➭ *[p]ig* / *[p]insta* [link]
 │    _Download Instagram Reels/Posts_
 │
 │ ➭ *[p]tw* / *[p]x* [link]
 │    _Download Twitter Videos_
 │
 │ ➭ *[p]spotify* / *[p]apple* [link]
 │    _Download Music directly_
 │
 ╰──────────────────────╯`

	case "aimenu":
		bodyText = `
 ╭── ✦ [ *AI CHAT MENU* ] ✦ ──╮
 │
 │ ➭ *[p]ai* / *[p]ask* / *[p]bot* [query]
 │    _Talk to the ultra smart AI_
 │
 │ ➭ *[p]gpt* / *[p]chatgpt* [query]
 │    _Standard ChatGPT Interface_
 │
 │ ➭ *[p]gemini* / *[p]claude* [query]
 │    _Advanced AI Models_
 │
 ╰──────────────────────╯`

	case "gpmenu":
		bodyText = `
 ╭── ✦ [ *GROUP MENU* ] ✦ ──╮
 │
 │ ➭ *[p]antilink* [on/off]
 │    _Enable/Disable link protection_
 │
 │ ➭ *[p]kick* / *[p]add* [@user]
 │    _Manage Group Members_
 │
 │ ➭ *[p]promote* / *[p]demote* [@user]
 │    _Manage Group Admins_
 │
 │ ➭ *[p]gcname* / *[p]gcdp*
 │    _Change Group Name or Photo_
 │
 │ ➭ *[p]link* / *[p]revoke*
 │    _Get or Reset Group Link_
 │
 ╰──────────────────────╯`

	case "ownermenu":
		bodyText = `
 ╭── ✦ [ *OWNER MENU* ] ✦ ──╮
 │
 │ ➭ *[p]setprefix* [symbol]
 │    _Change the Bot Prefix_
 │
 │ ➭ *[p]mode* [public/private]
 │    _Change Bot Access Mode_
 │
 │ ➭ *[p]pair* [number]
 │    _Generate a pairing code_
 │
 │ ➭ *[p]block* [@user]
 │    _Block a user permanently_
 │
 ╰──────────────────────╯`

	case "utilmenu":
		bodyText = `
┍──╼〔 🛠️ *UTILITY* 〕
│ ⬡ [p]vv
│ ⬡ [p]id
│ ⬡ [p]vc
┕━━━━━━━━━━━━━━━━━━┈`

	case "editmenu":
		bodyText = `
┍──╼〔 🎨 *EDITING ZONE* 〕
│ ⬡ [p]s
│  [p]sticker
│ ⬡ [p]toimg
│ ⬡ [p]togif
│ ⬡ [p]stovideo
│ ⬡ [p]stourl
│ ⬡ [p]stoptt
│ ⬡ [p]fancy
┕━━━━━━━━━━━━━━━━━━┈`

	case "aitools":
		bodyText = `
┍──╼〔 ✨ *AI TOOLS* 〕
│ ⬡ [p]img
│ ⬡ [p]remini
│ ⬡ [p]removebg
│ ⬡ [p]tr
│ ⬡ [p]ss
│ ⬡ [p]google
│ ⬡ [p]weather
┕━━━━━━━━━━━━━━━━━━┈`

	case "artificialintelligencemenu":
		bodyText = `
 ╭── ✦ [ *ARTIFICIAL INTELLIGENCE MENU* ] ✦ ──╮
 │
 │ ➭ *[p]3d* [text]
 │    _Generate AI images using Live3D_
 │
 │ ➭ *[p]ai4chat* [text]
 │    _Chat with AI4Chat assistant_
 │
 │ ➭ *[p]aiappchat* [text]
 │    _Advanced AI chat with vision sup..._
 │
 │ ➭ *[p]aiappgen* [text]
 │    _Generate images using Flux, DALL_
 │
 │ ➭ *[p]dalle* [text]
 │    _Generate high_
 │
 │ ➭ *[p]aichat* [text]
 │    _Chat with AI assistant_
 │
 │ ➭ *[p]aiserv* [text]
 │    _Advanced AI with GPT_
 │
 │ ➭ *[p]quick* [text]
 │    _Generate quick stories with defa..._
 │
 │ ➭ *[p]advanced* [text]
 │    _Generate stories with customizab..._
 │
 │ ➭ *[p]animekill* [image]
 │    _Transform your images into anime..._
 │
 │ ➭ *[p]blackbox* [text]
 │    _Advanced AI chat assistant with ..._
 │
 │ ➭ *[p]borli* [action]
 │    _Chat with AI characters. Search ..._
 │
 │ ➭ *[p]cartoon* [imageurl]
 │    _Turn your photos into anime styl..._
 │
 │ ➭ *[p]copilot* [text]
 │    _Chat with Microsoft Copilot AI_
 │
 │ ➭ *[p]copilotthink* [text]
 │    _Deep thinking mode for complex r..._
 │
 │ ➭ *[p]gpt5* [text]
 │    _Advanced GPT_
 │
 │ ➭ *[p]sch* [text]
 │    _Free AI chat_
 │
 │ ➭ *[p]chatbot* [text]
 │    _AI chatbot with optional web sea..._
 │
 │ ➭ *[p]chatevo* [text]
 │    _Generate images from text prompt..._
 │
 │ ➭ *[p]chatex* [text]
 │    _Chat with Chatex AI (GPT_
 │
 │ ➭ *[p]chatup* [text]
 │    _Chat with AI using ChatUp_
 │
 │ ➭ *[p]prompttocode* [text]
 │    _Generate code from text prompt i..._
 │
 │ ➭ *[p]detectbugs* [code]
 │    _Find and fix bugs in your code_
 │
 │ ➭ *[p]convertcode* [code]
 │    _Convert code between programming..._
 │
 │ ➭ *[p]explaincode* [code]
 │    _Get detailed explanation of any ..._
 │
 │ ➭ *[p]chateverywhere* [text]
 │    _Chat with GPT_
 │
 │ ➭ *[p]chateverywherereset* [userId]
 │    _Reset conversation history for C..._
 │
 │ ➭ *[p]deepquery* [text]
 │    _Advanced AI with deep knowledge_
 │
 │ ➭ *[p]logical* [text]
 │    _For analytical and reasoning que..._
 │
 │ ➭ *[p]creative* [text]
 │    _For creative writing and brainst..._
 │
 │ ➭ *[p]summarize* [text]
 │    _For text summarization and conde..._
 │
 │ ➭ *[p]codebeginner* [text]
 │    _For beginner programming questions_
 │
 │ ➭ *[p]codeadvanced* [text]
 │    _For advanced programming and alg..._
 │
 │ ➭ *[p]dream* [dream]
 │    _Interpret your dreams with AI_
 │
 │ ➭ *[p]deepseekchat* [text]
 │    _Standard DeepSeek chat model_
 │
 │ ➭ *[p]deepseekreasoner* [text]
 │    _DeepSeek model with reasoning ca..._
 │
 │ ➭ *[p]reset*
 │    _Reset DeepSeek conversation history_
 │
 │ ➭ *[p]history*
 │    _Get current conversation history_
 │
 │ ➭ *[p]easemategenerate* [text]
 │    _Generate images from text prompt..._
 │
 │ ➭ *[p]easematechat* [text]
 │    _Chat with AI assistant powered b..._
 │
 │ ➭ *[p]homeplannerchat* [text]
 │    _AI chat assistant with optional ..._
 │
 │ ➭ *[p]homeplannerimage* [text]
 │    _Generate images from text prompt..._
 │
 │ ➭ *[p]homeplannertts* [text]
 │    _Convert text to speech with mult..._
 │
 │ ➭ *[p]homeplannersearch* [text]
 │    _Search the web with AI_
 │
 │ ➭ *[p]homeplanneryt* [text]
 │    _Get AI summary of any YouTube video_
 │
 │ ➭ *[p]img2img* [imageUrl]
 │    _AI_
 │
 │ ➭ *[p]lumo* [text]
 │    _Encrypted AI chat by Proton_
 │
 │ ➭ *[p]chatmusiclyrics* [text]
 │    _Generate song lyrics using AI fr..._
 │
 │ ➭ *[p]chatmusiccreate* [text]
 │    _Create AI_
 │
 │ ➭ *[p]chatmusicstatus* [state]
 │    _Check the status of your music g..._
 │
 │ ➭ *[p]mydreams* [text]
 │    _Generate images with AI_
 │
 │ ➭ *[p]solabiba* [text]
 │    _Free AI chatbot with multiple mo..._
 │
 │ ➭ *[p]photogpt* [text]
 │    _AI image generation with multipl..._
 │
 │ ➭ *[p]photonex* [image]
 │    _Transform images using AI_
 │
 │ ➭ *[p]soraremover* [link]
 │    _Remove watermark from Sora AI ge..._
 │
 │ ➭ *[p]txt2img* [text]
 │    _AI image generator from text pro..._
 │
 │ ➭ *[p]sunlimai* [image]
 │    _AI_
 │
 │ ➭ *[p]video* [text]
 │    _AI video generation from text pr..._
 │
 │ ➭ *[p]saiwriterchat* [text]
 │    _Chat with AI Writer models (gpt_
 │
 │ ➭ *[p]saiwriterimage* [text]
 │    _Generate images with AI Writer_
 │
 │ ➭ *[p]saiwritermodels*
 │    _Get list of available AI Writer ..._
 │
 ╰──────────────────────╯`

	case "imagegenerationmenu":
		bodyText = `
 ╭── ✦ [ *IMAGE GENERATION MENU* ] ✦ ──╮
 │
 │ ➭ *[p]realistic* [text]
 │    _Generate realistic/photographic ..._
 │
 │ ➭ *[p]anime* [text]
 │    _Generate anime_
 │
 │ ➭ *[p]fantasy* [text]
 │    _Generate fantasy/artistic images_
 │
 │ ➭ *[p]cyberpunk* [text]
 │    _Generate cyberpunk/futuristic im..._
 │
 │ ➭ *[p]watercolor* [text]
 │    _Generate watercolor painting sty..._
 │
 │ ➭ *[p]oilpainting* [text]
 │    _Generate oil painting style images_
 │
 │ ➭ *[p]pixelart* [text]
 │    _Generate pixel art style images_
 │
 │ ➭ *[p]sketch* [text]
 │    _Generate sketch/drawing style im..._
 │
 │ ➭ *[p]abstract* [text]
 │    _Generate abstract art images_
 │
 │ ➭ *[p]minimalist* [text]
 │    _Generate minimalist style images_
 │
 │ ➭ *[p]surreal* [text]
 │    _Generate surreal/abstract images_
 │
 │ ➭ *[p]vintage* [text]
 │    _Generate vintage/retro style images_
 │
 │ ➭ *[p]steampunk* [text]
 │    _Generate steampunk style images_
 │
 │ ➭ *[p]horror* [text]
 │    _Generate horror/dark style images_
 │
 │ ➭ *[p]scifi* [text]
 │    _Generate science fiction style i..._
 │
 │ ➭ *[p]popart* [text]
 │    _Generate pop art style images_
 │
 ╰──────────────────────╯`

	case "animemenu":
		bodyText = `
 ╭── ✦ [ *ANIME MENU* ] ✦ ──╮
 │
 │ ➭ *[p]animekillhome*
 │    _Get homepage anime TV listing_
 │
 │ ➭ *[p]animekillhomestatic*
 │    _Get static homepage data_
 │
 │ ➭ *[p]animekillsearch* [text]
 │    _Search for anime by title_
 │
 │ ➭ *[p]animekilldetail* [anime_id]
 │    _Get detailed information about a..._
 │
 │ ➭ *[p]animekillepisodes* [anime_id]
 │    _Get list of episodes for an anime_
 │
 │ ➭ *[p]animekillstream* [anime_id]
 │    _Get video stream URL for an episode_
 │
 │ ➭ *[p]animekillcomments* [anime_id]
 │    _Get comments for an anime_
 │
 │ ➭ *[p]animekillbygenre* [genre]
 │    _Get anime filtered by genre_
 │
 │ ➭ *[p]animekillgenres*
 │    _Get list of all available genres_
 │
 │ ➭ *[p]animekillschedule*
 │    _Get weekly anime schedule_
 │
 │ ➭ *[p]animesearch* [text]
 │    _Search for anime by title/keyword_
 │
 │ ➭ *[p]animedetail* [link]
 │    _Get detailed information and epi..._
 │
 │ ➭ *[p]animedownload* [link]
 │    _Get streaming servers and downlo..._
 │
 │ ➭ *[p]mangahome* [page]
 │    _Get latest manga updates from Ko..._
 │
 │ ➭ *[p]mangasearch* [text]
 │    _Search for manga on Komiku_
 │
 │ ➭ *[p]mangadetail* [id]
 │    _Get detailed information about a..._
 │
 │ ➭ *[p]mangachapter* [chapter_id]
 │    _Get manga chapter images_
 │
 │ ➭ *[p]mangasuggestions* [suggestion_type]
 │    _Get manga suggestions_
 │
 │ ➭ *[p]mangaepisodes* [id]
 │    _Get episodes list of a manga_
 │
 │ ➭ *[p]mangaseries* [id]
 │    _Get series contents of a manga_
 │
 │ ➭ *[p]mangacomments* [content_id]
 │    _Get comments for a manga_
 │
 │ ➭ *[p]mangarankfilters*
 │    _Get ranking filters_
 │
 │ ➭ *[p]mangaranktags*
 │    _Get top tags for ranking_
 │
 │ ➭ *[p]hug*
 │    _Get a random anime hug GIF_
 │
 │ ➭ *[p]slap*
 │    _Get a random anime slap GIF_
 │
 │ ➭ *[p]pat*
 │    _Get a random anime pat GIF_
 │
 │ ➭ *[p]cry*
 │    _Get a random anime cry GIF_
 │
 │ ➭ *[p]kill*
 │    _Get a random anime kill GIF_
 │
 │ ➭ *[p]bite*
 │    _Get a random anime bite GIF_
 │
 │ ➭ *[p]yeet*
 │    _Get a random anime yeet GIF_
 │
 │ ➭ *[p]bully*
 │    _Get a random anime bully GIF_
 │
 │ ➭ *[p]bonk*
 │    _Get a random anime bonk GIF_
 │
 │ ➭ *[p]wink*
 │    _Get a random anime wink GIF_
 │
 │ ➭ *[p]poke*
 │    _Get a random anime poke GIF_
 │
 │ ➭ *[p]nom*
 │    _Get a random anime nom GIF_
 │
 │ ➭ *[p]smile*
 │    _Get a random anime smile GIF_
 │
 │ ➭ *[p]wave*
 │    _Get a random anime wave GIF_
 │
 │ ➭ *[p]awoo*
 │    _Get a random anime awoo GIF_
 │
 │ ➭ *[p]blush*
 │    _Get a random anime blush GIF_
 │
 │ ➭ *[p]smug*
 │    _Get a random anime smug GIF_
 │
 │ ➭ *[p]glomp*
 │    _Get a random anime glomp GIF_
 │
 │ ➭ *[p]happy*
 │    _Get a random anime happy GIF_
 │
 │ ➭ *[p]dance*
 │    _Get a random anime dance GIF_
 │
 │ ➭ *[p]cringe*
 │    _Get a random anime cringe GIF_
 │
 │ ➭ *[p]cuddle*
 │    _Get a random anime cuddle GIF_
 │
 │ ➭ *[p]highfive*
 │    _Get a random anime highfive GIF_
 │
 │ ➭ *[p]handhold*
 │    _Get a random anime handhold GIF_
 │
 │ ➭ *[p]shinobu*
 │    _Get a random anime shinobu GIF_
 │
 │ ➭ *[p]reactions*
 │    _Get anime reaction GIFs/images_
 │
 │ ➭ *[p]webnovelhot*
 │    _Get hot search terms from Webnovel_
 │
 │ ➭ *[p]webnovelrank* [page]
 │    _Get novel rankings from Webnovel_
 │
 │ ➭ *[p]webnovelsearch* [text]
 │    _Search for novels on Webnovel_
 │
 │ ➭ *[p]webnoveldetail* [bid]
 │    _Get detailed information about a..._
 │
 │ ➭ *[p]webnovelchapter* [bid]
 │    _Get chapter content from a novel_
 │
 ╰──────────────────────╯`

	case "gamesmenu":
		bodyText = `
 ╭── ✦ [ *GAMES MENU* ] ✦ ──╮
 │
 │ ➭ *[p]quizcategories*
 │    _Get all available quiz categories_
 │
 │ ➭ *[p]quizguess* [level]
 │    _Guess the correct answer from mu..._
 │
 │ ➭ *[p]quizpuzzle* [level]
 │    _Solve puzzle_
 │
 │ ➭ *[p]quiztruefalse* [level]
 │    _Answer true or false questions_
 │
 │ ➭ *[p]quizrandom* [level]
 │    _Mixed random questions from all ..._
 │
  ╰──────────────────────╯`

	case "imagecreatormenu":
		bodyText = `
 ╭── ✦ [ *IMAGE CREATOR MENU* ] ✦ ──╮
 │
 │ ➭ *[p]image* [text]
 │    _Create simple text images with c..._
 │
 │ ➭ *[p]gif* [text]
 │    _Create animated text GIFs_
 │
 │ ➭ *[p]mp4* [text]
 │    _Create animated text videos_
 │
 │ ➭ *[p]meme* [topText]
 │    _Create memes with top and bottom..._
 │
 │ ➭ *[p]memetext* [text]
 │    _Create text_
 │
 │ ➭ *[p]spongebob* [text]
 │    _Create SpongeBob "How dare you" ..._
 │
 │ ➭ *[p]ttp* [text]
 │    _Text to ttp_
 │
 ╰──────────────────────╯`

	case "moviesmenu":
		bodyText = `
 ╭── ✦ [ *MOVIES MENU* ] ✦ ──╮
 │
 │ ➭ *[p]moviesearch* [text]
 │    _Search for movies by title or ke..._
 │
 │ ➭ *[p]moviedetail* [link]
 │    _Get detailed information about a..._
 │
 │ ➭ *[p]search* [text]
 │    _Search for movies by keyword_
 │
 │ ➭ *[p]suggest* [text]
 │    _Get search suggestions for a key..._
 │
 │ ➭ *[p]detail* [id]
 │    _Get detailed information about a..._
 │
 │ ➭ *[p]recommendations* [id]
 │    _Get recommended movies based on ..._
 │
 │ ➭ *[p]trending* [tabId]
 │    _Get trending movies_
 │
 │ ➭ *[p]home*
 │    _Get home page feed with featured..._
 │
 │ ➭ *[p]countries*
 │    _Get list of available country codes_
 │
 ╰──────────────────────╯`

	case "searchmenu":
		bodyText = `
 ╭── ✦ [ *SEARCH MENU* ] ✦ ──╮
 │
 │ ➭ *[p]android1* [text]
 │    _Search modded game on android1_
 │
 │ ➭ *[p]applemusic* [text]
 │    _Search for songs, artists, and p..._
 │
 │ ➭ *[p]cuaca* [kota]
 │    _Search info cuaca_
 │
 │ ➭ *[p]repos* [text]
 │    _Search for GitHub repositories_
 │
 │ ➭ *[p]users* [text]
 │    _Search for GitHub users_
 │
 │ ➭ *[p]issues* [text]
 │    _Search for GitHub issues_
 │
 │ ➭ *[p]code* [text]
 │    _Search for code on GitHub_
 │
 │ ➭ *[p]imdb* [text]
 │    _Search for movie/series informat..._
 │
 │ ➭ *[p]lyrics* [title]
 │    _Search for song lyrics_
 │
 │ ➭ *[p]nik* [text]
 │    _Search nik_
 │
 │ ➭ *[p]wallpaper* [text]
 │    _Search HD wallpapers_
 │
 │ ➭ *[p]telegram* [text]
 │    _Search for Telegram channels by ..._
 │
 │ ➭ *[p]tggroup* [text]
 │    _Search for Telegram groups and c..._
 │
 │ ➭ *[p]tiktoksearch* [text]
 │    _Search video on tiktok_
 │
 │ ➭ *[p]wagroup* [text]
 │    _Search for WhatsApp groups by ke..._
 │
 │ ➭ *[p]youtube* [text]
 │    _Search video on youtube_
 │
 │ ➭ *[p]ytmonet* [link]
 │    _YouTube monetization checker_
 │
 ╰──────────────────────╯`

	case "randommenu":
		bodyText = `
 ╭── ✦ [ *RANDOM MENU* ] ✦ ──╮
 │
 │ ➭ *[p]akiyama*
 │    _Get random Akiyama anime images_
 │
 │ ➭ *[p]ana*
 │    _Get random Ana anime images_
 │
 │ ➭ *[p]asuna*
 │    _Get random Asuna anime images_
 │
 │ ➭ *[p]ayuzawa*
 │    _Get random Ayuzawa anime images_
 │
 │ ➭ *[p]boruto*
 │    _Get random Boruto anime images_
 │
 │ ➭ *[p]chitoge*
 │    _Get random Chitoge anime images_
 │
 │ ➭ *[p]deidara*
 │    _Get random Deidara anime images_
 │
 │ ➭ *[p]doraemon*
 │    _Get random Doraemon anime images_
 │
 │ ➭ *[p]elaina*
 │    _Get random Elaina anime images_
 │
 │ ➭ *[p]emilia*
 │    _Get random Emilia anime images_
 │
 │ ➭ *[p]erza*
 │    _Get random Erza anime images_
 │
 │ ➭ *[p]hestia*
 │    _Get random Hestia anime images_
 │
 │ ➭ *[p]husbu*
 │    _Get random Husbu anime images_
 │
 │ ➭ *[p]inori*
 │    _Get random Inori anime images_
 │
 │ ➭ *[p]itachi*
 │    _Get random Itachi anime images_
 │
 │ ➭ *[p]kagura*
 │    _Get random Kagura anime images_
 │
 │ ➭ *[p]kaori*
 │    _Get random Kaori anime images_
 │
 │ ➭ *[p]keneki*
 │    _Get random Keneki anime images_
 │
 │ ➭ *[p]kotori*
 │    _Get random Kotori anime images_
 │
 │ ➭ *[p]kurumi*
 │    _Get random Kurumi anime images_
 │
 │ ➭ *[p]madara*
 │    _Get random Madara anime images_
 │
 │ ➭ *[p]megumin*
 │    _Get random Megumin anime images_
 │
 │ ➭ *[p]mikasa*
 │    _Get random Mikasa anime images_
 │
 │ ➭ *[p]miku*
 │    _Get random Miku anime images_
 │
 │ ➭ *[p]minato*
 │    _Get random Minato anime images_
 │
 │ ➭ *[p]naruto*
 │    _Get random Naruto anime images_
 │
 │ ➭ *[p]nekonime*
 │    _Get random Nekonime anime images_
 │
 │ ➭ *[p]nezuko*
 │    _Get random Nezuko anime images_
 │
 │ ➭ *[p]onepiece*
 │    _Get random One Piece anime images_
 │
 │ ➭ *[p]rize*
 │    _Get random Rize anime images_
 │
 │ ➭ *[p]sagiri*
 │    _Get random Sagiri anime images_
 │
 │ ➭ *[p]sakura*
 │    _Get random Sakura anime images_
 │
 │ ➭ *[p]sasuke*
 │    _Get random Sasuke anime images_
 │
 │ ➭ *[p]shinomiya*
 │    _Get random Shinomiya anime images_
 │
 │ ➭ *[p]tsunade*
 │    _Get random Tsunade anime images_
 │
 │ ➭ *[p]yotsuba*
 │    _Get random Yotsuba anime images_
 │
 │ ➭ *[p]yuki*
 │    _Get random Yuki anime images_
 │
 │ ➭ *[p]yumeko*
 │    _Get random Yumeko anime images_
 │
 │ ➭ *[p]art*
 │    _Get random art wallpapers_
 │
 │ ➭ *[p]cyber*
 │    _Get random cyber wallpapers_
 │
 │ ➭ *[p]gamewallpaper*
 │    _Get random game wallpapers_
 │
 │ ➭ *[p]mountain*
 │    _Get random mountain wallpapers_
 │
 │ ➭ *[p]programming*
 │    _Get random programming wallpapers_
 │
 │ ➭ *[p]space*
 │    _Get random space wallpapers_
 │
 │ ➭ *[p]technology*
 │    _Get random technology wallpapers_
 │
 │ ➭ *[p]wallhp*
 │    _Get random mobile wallpapers_
 │
 │ ➭ *[p]wallml*
 │    _Get random Mobile Legends wallpa..._
 │
 │ ➭ *[p]wallmlnime*
 │    _Get random anime wallpapers_
 │
 ╰──────────────────────╯`

	case "audiomenu":
		bodyText = `
 ╭── ✦ [ *AUDIO MENU* ] ✦ ──╮
 │
 │ ➭ *[p]download* [link]
 │    _Download sound as MP3 audio file_
 │
 │ ➭ *[p]nonstick* [type]
 │    _Get sound sources from NonStick_
 │
 ╰──────────────────────╯`

	case "sportsmenu":
		bodyText = `
 ╭── ✦ [ *SPORTS MENU* ] ✦ ──╮
 │
 │ ➭ *[p]football* [detail]
 │    _Get live football matches, score..._
 │
 │ ➭ *[p]basketball* [detail]
 │    _Get live basketball matches, sco..._
 │
 │ ➭ *[p]othersports* [detail]
 │    _Get live matches for other sport..._
 │
 ╰──────────────────────╯`

	case "screenshotwebsitemenu":
		bodyText = `
 ╭── ✦ [ *SCREENSHOT WEBSITE MENU* ] ✦ ──╮
 │
 │ ➭ *[p]webss* [link]
 │    _Take screenshot using WebSS prov..._
 │
 │ ➭ *[p]apiflash* [link]
 │    _Take screenshot using Flash prov..._
 │
 │ ➭ *[p]screenshotlayer* [link]
 │    _Take screenshot using Screenshot..._
 │
 ╰──────────────────────╯`

	case "stalkmenu":
		bodyText = `
 ╭── ✦ [ *STALK MENU* ] ✦ ──╮
 │
 │ ➭ *[p]ffstalk* [id]
 │    _Get info freefire account_
 │
 │ ➭ *[p]igstalk* [user]
 │    _Get info instagram account (@user)_
 │
 │ ➭ *[p]igstalkv2* [user]
 │    _Get info instagram account_
 │
 │ ➭ *[p]ttstalk* [user]
 │    _Get info tiktok account_
 │
 │ ➭ *[p]twitterstalk* [user]
 │    _Get info twitter account_
 │
 │ ➭ *[p]ytstalk* [user]
 │    _Get info youtube account (@user)_
 │
 ╰──────────────────────╯`

	case "textmakermenu":
		bodyText = `
 ╭── ✦ [ *TEXT MAKER MENU* ] ✦ ──╮
 │
 │ ➭ *[p]glitchtext* [text]
 │    _Create digital glitch text effects_
 │
 │ ➭ *[p]writetext* [text]
 │    _Write text on wet glass effect_
 │
 │ ➭ *[p]advancedglow* [text]
 │    _Advanced glow text effects_
 │
 │ ➭ *[p]typographytext* [text]
 │    _Create typography text effect on..._
 │
 │ ➭ *[p]pixelglitch* [text]
 │    _Create pixel glitch text effect_
 │
 │ ➭ *[p]neonglitch* [text]
 │    _Create impressive neon glitch te..._
 │
 │ ➭ *[p]flagtext* [text]
 │    _Nigeria 3D flag text effect_
 │
 │ ➭ *[p]flag3dtext* [text]
 │    _American flag 3D text effect_
 │
 │ ➭ *[p]deletingtext* [text]
 │    _Create eraser deleting text effect_
 │
 │ ➭ *[p]blackpinkstyle* [text]
 │    _Blackpink style logo maker effect_
 │
 │ ➭ *[p]glowingtext* [text]
 │    _Create glowing text effects_
 │
 │ ➭ *[p]underwatertext* [text]
 │    _3D underwater text effect_
 │
 │ ➭ *[p]logomaker* [text]
 │    _Free bear logo maker_
 │
 │ ➭ *[p]cartoonstyle* [text]
 │    _Create cartoon style graffiti te..._
 │
 │ ➭ *[p]papercutstyle* [text]
 │    _Multicolor 3D paper cut style te..._
 │
 │ ➭ *[p]watercolortext* [text]
 │    _Create a watercolor text effect_
 │
 │ ➭ *[p]effectclouds* [text]
 │    _Write text effect clouds in the sky_
 │
 │ ➭ *[p]blackpinklogo* [text]
 │    _Create Blackpink logo online_
 │
 │ ➭ *[p]gradienttext* [text]
 │    _Create 3D gradient text effect_
 │
 │ ➭ *[p]summerbeach* [text]
 │    _Write in sand summer beach_
 │
 │ ➭ *[p]luxurygold* [text]
 │    _Create a luxury gold text effect_
 │
 │ ➭ *[p]multicoloredneon* [text]
 │    _Create multicolored neon light s..._
 │
 │ ➭ *[p]sandsummer* [text]
 │    _Write in sand summer beach_
 │
 │ ➭ *[p]galaxywallpaper* [text]
 │    _Create galaxy wallpaper mobile_
 │
 │ ➭ *[p]style1917* [text]
 │    _1917 style text effect_
 │
 │ ➭ *[p]makingneon* [text]
 │    _Making neon light text effect wi..._
 │
 │ ➭ *[p]royaltext* [text]
 │    _Royal text effect online_
 │
 │ ➭ *[p]freecreate* [text]
 │    _Free create a 3D hologram text e..._
 │
 │ ➭ *[p]galaxystyle* [text]
 │    _Create galaxy style free name logo_
 │
 │ ➭ *[p]lighteffects* [text]
 │    _Create light effects green neon_
 │
 ╰──────────────────────╯`

	case "toolsmenu":
		bodyText = `
 ╭── ✦ [ *TOOLS MENU* ] ✦ ──╮
 │
 │ ➭ *[p]sendemail* [to]
 │    _Send anonymous emails without re..._
 │
 │ ➭ *[p]codeanalyzer* [code]
 │    _Analyze code for security issues..._
 │
 │ ➭ *[p]codeconverter* [code]
 │    _Convert code between programming..._
 │
 │ ➭ *[p]tojavascript* [code]
 │    _Convert code to JavaScript_
 │
 │ ➭ *[p]topython* [code]
 │    _Convert code to Python_
 │
 │ ➭ *[p]tojava* [code]
 │    _Convert code to Java_
 │
 │ ➭ *[p]tocpp* [code]
 │    _Convert code to C++_
 │
 │ ➭ *[p]tophp* [code]
 │    _Convert code to PHP_
 │
 │ ➭ *[p]compiler* [code]
 │    _Compile and execute code in mult..._
 │
 │ ➭ *[p]compilejs* [code]
 │    _Compile and run JavaScript code_
 │
 │ ➭ *[p]compilepython* [code]
 │    _Compile and run Python code_
 │
 │ ➭ *[p]compilejava* [code]
 │    _Compile and run Java code_
 │
 │ ➭ *[p]compilec* [code]
 │    _Compile and run C code_
 │
 │ ➭ *[p]compilecpp* [code]
 │    _Compile and run C++ code_
 │
 │ ➭ *[p]compilecsharp* [code]
 │    _Compile and run C# code_
 │
 │ ➭ *[p]emojiencrypt* [input]
 │    _Encrypt text into emojis using p..._
 │
 │ ➭ *[p]emojidecrypt* [input]
 │    _Decrypt emojis back to text usin..._
 │
 │ ➭ *[p]htmlecnc* [html]
 │    _Encrypt and obfuscate HTML code_
 │
 │ ➭ *[p]htmlbasic* [html]
 │    _Encrypt HTML with basic obfuscation_
 │
 │ ➭ *[p]htmlextended* [html]
 │    _Encrypt HTML with extended security_
 │
 │ ➭ *[p]htmlhigh* [html]
 │    _Encrypt HTML with high security ..._
 │
 │ ➭ *[p]htmlmaximum* [html]
 │    _Encrypt HTML with maximum security_
 │
 │ ⬭ *[p]fdroidsearch* [text]
 │    _Search for apps on F_
 │
 │ ➭ *[p]fdroidpackage* [link]
 │    _Get detailed information about F_
 │
 │ ➭ *[p]fdroidapp* [package]
 │    _Get app details by package name ..._
 │
 │ ➭ *[p]geoip* [ip]
 │    _Get geolocation information for ..._
 │
 │ ➭ *[p]myip*
 │    _Get geolocation information for ..._
 │
 │ ➭ *[p]hostcheck* [domain]
 │    _Get detailed hosting information..._
 │
 │ ➭ *[p]hostchecksimple* [domain]
 │    _Get basic hosting information fo..._
 │
 │ ➭ *[p]html2img* [html]
 │    _Convert HTML to image_
 │
 │ ➭ *[p]html2imgdirect* [html]
 │    _Convert HTML to image_
 │
 │ ➭ *[p]obflow* [code]
 │    _Obfuscate JavaScript code with l..._
 │
 │ ➭ *[p]obfmedium* [code]
 │    _Obfuscate JavaScript code with m..._
 │
 │ ➭ *[p]obfhigh* [code]
 │    _Obfuscate JavaScript code with h..._
 │
 │ ➭ *[p]obfextreme* [code]
 │    _Obfuscate JavaScript code with e..._
 │
 │ ➭ *[p]tiktoktranscript* [link]
 │    _Get transcript from TikTok video_
 │
 │ ➭ *[p]entoid* [text]
 │    _Translate from English to Indone..._
 │
 │ ➭ *[p]idtoen* [text]
 │    _Translate from Indonesian to Eng..._
 │
 │ ➭ *[p]jatoid* [text]
 │    _Translate from Japanese to Indon..._
 │
 │ ➭ *[p]kotoid* [text]
 │    _Translate from Korean to Indonesian_
 │
 │ ➭ *[p]zhtoid* [text]
 │    _Translate from Chinese to Indone..._
 │
 │ ➭ *[p]artoid* [text]
 │    _Translate from Arabic to Indonesian_
 │
 │ ➭ *[p]detectlanguage* [text]
 │    _Detect the language of provided ..._
 │
 │ ➭ *[p]languages*
 │    _Get list of supported languages ..._
 │
 │ ➭ *[p]youtubetranscript* [link]
 │    _Get YouTube video transcript, de..._
 │
 ╰──────────────────────╯`

	case "urlshortnermenu":
		bodyText = `
 ╭── ✦ [ *URL SHORTNER MENU* ] ✦ ──╮
 │
 │ ➭ *[p]dagd* [link]
 │    _Shorten URL using da.gd service_
 │
 │ ➭ *[p]vgd* [link]
 │    _Shorten URL using v.gd service_
 │
 │ ➭ *[p]tinube* [link]
 │    _Shorten URL using tinu.be service_
 │
 │ ➭ *[p]spoome* [link]
 │    _Shorten URL using Spoo.me service_
 │
 │ ➭ *[p]spooemoji* [link]
 │    _Shorten URL with emojis using Sp..._
 │
 │ ➭ *[p]random* [link]
 │    _Shorten URL using random provider_
 │
 ╰──────────────────────╯`

	case "styletextmenu":
		bodyText = `
 ╭── ✦ [ *STYLETEXT MENU* ] ✦ ──╮
 │
 │ ➭ *[p]allstyles* [text]
 │    _Generate all 35+ text styles at ..._
 │
 │ ➭ *[p]circled* [text]
 │    _Convert text to Circled style_
 │
 │ ➭ *[p]circledneg* [text]
 │    _Convert text to Circled (neg) style_
 │
 │ ➭ *[p]fullwidth* [text]
 │    _Convert text to Fullwidth style_
 │
 │ ➭ *[p]mathbold* [text]
 │    _Convert text to Math bold style_
 │
 │ ➭ *[p]mathboldfraktur* [text]
 │    _Convert text to Math bold Fraktu..._
 │
 │ ➭ *[p]mathbolditalic* [text]
 │    _Convert text to Math bold italic..._
 │
 │ ➭ *[p]mathboldscript* [text]
 │    _Convert text to Math bold script..._
 │
 │ ➭ *[p]mathdoublestruck* [text]
 │    _Convert text to Math double_
 │
 │ ➭ *[p]mathmonospace* [text]
 │    _Convert text to Math monospace s..._
 │
 │ ➭ *[p]mathsans* [text]
 │    _Convert text to Math sans style_
 │
 │ ➭ *[p]mathsansbold* [text]
 │    _Convert text to Math sans bold s..._
 │
 │ ➭ *[p]mathsansbolditalic* [text]
 │    _Convert text to Math sans bold i..._
 │
 │ ➭ *[p]mathsansitalic* [text]
 │    _Convert text to Math sans italic..._
 │
 │ ➭ *[p]parenthesized* [text]
 │    _Convert text to Parenthesized style_
 │
 │ ➭ *[p]regionalindicator* [text]
 │    _Convert text to Regional Indicat..._
 │
 │ ➭ *[p]squared* [text]
 │    _Convert text to Squared style_
 │
 │ ➭ *[p]squaredneg* [text]
 │    _Convert text to Squared (neg) style_
 │
 │ ➭ *[p]tag* [text]
 │    _Convert text to Tag style_
 │
 │ ➭ *[p]acute* [text]
 │    _Convert text to A_
 │
 │ ➭ *[p]cjkthai* [text]
 │    _Convert text to CJK+Thai style_
 │
 │ ➭ *[p]curvy1* [text]
 │    _Convert text to Curvy 1 style_
 │
 │ ➭ *[p]curvy2* [text]
 │    _Convert text to Curvy 2 style_
 │
 │ ➭ *[p]curvy3* [text]
 │    _Convert text to Curvy 3 style_
 │
 │ ➭ *[p]fauxcyrillic* [text]
 │    _Convert text to Faux Cyrillic style_
 │
 │ ➭ *[p]fauxethiopic* [text]
 │    _Convert text to Faux Ethiopic style_
 │
 │ ➭ *[p]mathfraktur* [text]
 │    _Convert text to Math Fraktur style_
 │
 │ ➭ *[p]rockdots* [text]
 │    _Convert text to Rock Dots style_
 │
 │ ➭ *[p]smallcaps* [text]
 │    _Convert text to Small Caps style_
 │
 │ ➭ *[p]stroked* [text]
 │    _Convert text to Stroked style_
 │
 │ ➭ *[p]subscript* [text]
 │    _Convert text to Subscript style_
 │
 │ ➭ *[p]superscript* [text]
 │    _Convert text to Superscript style_
 │
 │ ➭ *[p]inverted* [text]
 │    _Convert text to Inverted style_
 │
 │ ➭ *[p]reversed* [text]
 │    _Convert text to Reversed style_
 │
 ╰──────────────────────╯`

	case "texttospeechmenu":
		bodyText = `
 ╭── ✦ [ *TEXT TO SPEECH MENU* ] ✦ ──╮
 │
 │ ➭ *[p]ttsen* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsid* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttses* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsfr* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsde* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsit* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttspt* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsnl* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttspl* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsru* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsja* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsko* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttzhcn* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttszhtw* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsar* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttshi* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsth* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsvi* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttstr* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttssv* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsno* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsda* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsfi* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsel* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttshe* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttscs* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttshu* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsro* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]ttsuk* [text]
 │    _Convert text to speech_
 │
 │ ➭ *[p]xiaofei* [text]
 │    _Chinese female voice with 19 styles_
 │
 │ ➭ *[p]xiaolei* [text]
 │    _Chinese male voice with 10 styles_
 │
 │ ➭ *[p]xiaojie* [text]
 │    _Chinese male voice with 12 styles_
 │
 │ ➭ *[p]xiaohua* [text]
 │    _Chinese female voice_
 │
 │ ➭ *[p]xiaofeng* [text]
 │    _Chinese male voice with 8 styles_
 │
 │ ➭ *[p]xiaoze* [text]
 │    _Chinese male voice_
 │
 │ ➭ *[p]xiaoyuan* [text]
 │    _Chinese female voice with 9 styles_
 │
 │ ➭ *[p]xiaozheng* [text]
 │    _Chinese male voice with 3 styles_
 │
 │ ➭ *[p]xiaoying* [text]
 │    _Chinese female voice with 12 styles_
 │
 │ ➭ *[p]xiaoqing* [text]
 │    _Chinese female voice_
 │
 │ ➭ *[p]xiaoxiang* [text]
 │    _Chinese female voice_
 │
 │ ➭ *[p]xiaoyan* [text]
 │    _Chinese female voice with 3 styles_
 │
 │ ➭ *[p]xianran* [text]
 │    _Chinese child female voice with ..._
 │
 │ ➭ *[p]xiaoxue* [text]
 │    _Chinese female voice_
 │
 │ ➭ *[p]xiaoxuan* [text]
 │    _Chinese child female voice_
 │
 │ ➭ *[p]xiaolu* [text]
 │    _Chinese female voice with 6 styles_
 │
 │ ➭ *[p]xiaowei* [text]
 │    _Chinese male voice with 7 styles_
 │
 │ ➭ *[p]xiaozhe* [text]
 │    _Chinese male voice with advertis..._
 │
 │ ➭ *[p]xiaohao* [text]
 │    _Chinese male voice_
 │
 │ ➭ *[p]xiaoyi* [text]
 │    _Chinese male voice with 5 styles_
 │
 │ ➭ *[p]xiaotao* [text]
 │    _Chinese male voice with 8 styles_
 │
 │ ➭ *[p]xiaoming* [text]
 │    _Chinese male voice_
 │
 │ ➭ *[p]david* [text]
 │    _English US male voice with empat..._
 │
 │ ➭ *[p]layla* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]james* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]joey* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]jennifer* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]john* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]paul* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]xena* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]marcus* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]jacob* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]sam* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]camila* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]amy* [text]
 │    _English US female voice with 11 ..._
 │
 │ ➭ *[p]quincy* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]sally* [text]
 │    _English US female voice with 7 s..._
 │
 │ ➭ *[p]emma* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]ethan* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]michael* [text]
 │    _English US male voice with 11 st..._
 │
 │ ➭ *[p]olivia* [text]
 │    _English US female voice with 16 ..._
 │
 │ ➭ *[p]mia* [text]
 │    _English US female voice with 10 ..._
 │
 │ ➭ *[p]jackson* [text]
 │    _English US male voice with 10 st..._
 │
 │ ➭ *[p]matthew* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]sophia* [text]
 │    _English US female voice with 13 ..._
 │
 │ ➭ *[p]owen* [text]
 │    _English US male voice with 11 st..._
 │
 │ ➭ *[p]beatrice* [text]
 │    _English US female voice with 10 ..._
 │
 │ ➭ *[p]scott* [text]
 │    _English US male voice with 10 st..._
 │
 │ ➭ *[p]ivy* [text]
 │    _English US child female voice_
 │
 │ ➭ *[p]eric* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]kevin* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]hannah* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]katrina* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]victor* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]justin* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]leo* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]grace* [text]
 │    _English US female voice_
 │
 │ ➭ *[p]casey* [text]
 │    _English US neutral voice_
 │
 │ ➭ *[p]dylan* [text]
 │    _English US male voice with conve..._
 │
 │ ➭ *[p]julie* [text]
 │    _English US female voice with con..._
 │
 │ ➭ *[p]thomas* [text]
 │    _English US male voice_
 │
 │ ➭ *[p]freya* [text]
 │    _English UK female voice_
 │
 │ ➭ *[p]max* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]phoebe* [text]
 │    _English UK female voice with che..._
 │
 │ ➭ *[p]noah* [text]
 │    _English UK male voice with chat ..._
 │
 │ ➭ *[p]sophie* [text]
 │    _English UK female voice_
 │
 │ ➭ *[p]isla* [text]
 │    _English UK female voice_
 │
 │ ➭ *[p]theo* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]ella* [text]
 │    _English UK female voice_
 │
 │ ➭ *[p]freddie* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]arthur* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]isabella* [text]
 │    _English UK female voice_
 │
 │ ➭ *[p]evie* [text]
 │    _English UK child female voice_
 │
 │ ➭ *[p]william* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]henry* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]lily* [text]
 │    _English UK female voice_
 │
 │ ➭ *[p]charlie* [text]
 │    _English UK male voice_
 │
 │ ➭ *[p]ttsvoices*
 │    _Get list of all 30 available TTS..._
 │
 │ ➭ *[p]ttsadultfemale1americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultfemale2americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale1americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale2americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale3americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale4americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale5americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale6americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale7americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsadultmale8americanenglishtruvoice* [text]
 │    _Convert text to speech using "Ad..._
 │
 │ ➭ *[p]ttsfemalewhisper* [text]
 │    _Convert text to speech using "Fe..._
 │
 │ ➭ *[p]ttsmalewhisper* [text]
 │    _Convert text to speech using "Ma..._
 │
 │ ➭ *[p]ttsmary* [text]
 │    _Convert text to speech using "Ma..._
 │
 │ ➭ *[p]ttsmaryfortelephone* [text]
 │    _Convert text to speech using "Ma..._
 │
 │ ➭ *[p]ttsmaryinhall* [text]
 │    _Convert text to speech using "Ma..._
 │
 │ ➭ *[p]ttsmaryinspace* [text]
 │    _Convert text to speech using "Ma..._
 │
 │ ➭ *[p]ttsmaryinstadium* [text]
 │    _Convert text to speech using "Ma..._
 │
 │ ➭ *[p]ttsmike* [text]
 │    _Convert text to speech using "Mi..._
 │
 │ ➭ *[p]ttsmikefortelephone* [text]
 │    _Convert text to speech using "Mi..._
 │
 │ ➭ *[p]ttsmikeinhall* [text]
 │    _Convert text to speech using "Mi..._
 │
 │ ➭ *[p]ttsmikeinspace* [text]
 │    _Convert text to speech using "Mi..._
 │
 │ ➭ *[p]ttsmikeinstadium* [text]
 │    _Convert text to speech using "Mi..._
 │
 │ ➭ *[p]ttsrobosoftfive* [text]
 │    _Convert text to speech using "Ro..._
 │
 │ ➭ *[p]ttsrobosoftfour* [text]
 │    _Convert text to speech using "Ro..._
 │
 │ ➭ *[p]ttsrobosoftone* [text]
 │    _Convert text to speech using "Ro..._
 │
 │ ➭ *[p]ttsrobosoftsix* [text]
 │    _Convert text to speech using "Ro..._
 │
 │ ➭ *[p]ttsrobosoftthree* [text]
 │    _Convert text to speech using "Ro..._
 │
 │ ➭ *[p]ttsrobosofttwo* [text]
 │    _Convert text to speech using "Ro..._
 │
 │ ➭ *[p]ttssam* [text]
 │    _Convert text to speech using "Sa..._
 │
 │ ➭ *[p]ttsbonzi* [text]
 │    _Convert text to speech using "Bo..._
 │
 ╰──────────────────────╯`

	case "virtualnumbermenu":
		bodyText = `
 ╭── ✦ [ *VIRTUAL NUMBER MENU* ] ✦ ──╮
 │
 │ ➭ *[p]sms24countries*
 │    _Get available countries from SMS24_
 │
 │ ➭ *[p]sms24numbers* [country]
 │    _Get virtual numbers from SMS24 b..._
 │
 │ ➭ *[p]sms24messages* [number]
 │    _Get messages from SMS24 by number_
 │
 │ ➭ *[p]veepncountries*
 │    _Get available countries from Vee..._
 │
 │ ➭ *[p]veepnnumbers* [country]
 │    _Get virtual numbers from VeePN b..._
 │
 │ ➭ *[p]veepnmessages* [country]
 │    _Get messages from VeePN virtual ..._
 │
 ╰──────────────────────╯`

	default:
		return 
	}
	
	
	bodyText = strings.ReplaceAll(bodyText, "[p]", p)
	
	
	subMenuBody := headerText + "\n" + bodyText

	
	subBtnParams := map[string]string{
		"display_text": "📢 Join Channel",
		"url":          GlobalConfig.ChannelLink,
		"merchant_url": GlobalConfig.ChannelLink,
	}
	subBtnJSON, _ := json.Marshal(subBtnParams)

	subButtons := []*waProto.InteractiveMessage_NativeFlowMessage_NativeFlowButton{
		{
			Name:             proto.String("cta_url"),
			ButtonParamsJSON: proto.String(string(subBtnJSON)),
		},
	}

	subInteractiveMsg := &waProto.InteractiveMessage{
		Body: &waProto.InteractiveMessage_Body{
			Text: proto.String(subMenuBody),
		},
		Footer: &waProto.InteractiveMessage_Footer{
			Text: proto.String(fmt.Sprintf("🔥 POWERED BY %s 🔥", strings.ToUpper(GlobalConfig.BotName))),
		},
		InteractiveMessage: &waProto.InteractiveMessage_NativeFlowMessage_{
			NativeFlowMessage: &waProto.InteractiveMessage_NativeFlowMessage{
				Buttons: subButtons,
			},
		},
	}

		
	if SendLogo {
		if cachedLogoResp == nil {
			imageBytes, err := os.ReadFile("logo.png")
			if err == nil {
				resp, err := client.Upload(context.Background(), imageBytes, whatsmeow.MediaImage)
				if err == nil {
					cachedLogoResp = &resp
				}
			}
		}

		if cachedLogoResp != nil {
			subInteractiveMsg.Header = &waProto.InteractiveMessage_Header{
				HasMediaAttachment: proto.Bool(true),
				Media: &waProto.InteractiveMessage_Header_ImageMessage{
					ImageMessage: &waProto.ImageMessage{
						Mimetype:          proto.String("image/png"),
						URL:               proto.String(cachedLogoResp.URL),
						DirectPath:        proto.String(cachedLogoResp.DirectPath),
						MediaKey:          cachedLogoResp.MediaKey,
						FileEncSHA256:     cachedLogoResp.FileEncSHA256,
						FileSHA256:        cachedLogoResp.FileSHA256,
						FileLength:        proto.Uint64(uint64(cachedLogoResp.FileLength)),
						MediaKeyTimestamp: proto.Int64(time.Now().Unix()), 
					},
				},
			}
		}
	}


	if subInteractiveMsg.Header == nil {
		subInteractiveMsg.Header = &waProto.InteractiveMessage_Header{
			Title:              proto.String(fmt.Sprintf("👑 %s SUB MENU 👑", strings.ToUpper(GlobalConfig.BotName))),
			HasMediaAttachment: proto.Bool(false),
		}
	}

	msg := &waProto.Message{
		ViewOnceMessage: &waProto.FutureProofMessage{
			Message: &waProto.Message{
				InteractiveMessage: subInteractiveMsg,
			},
		},
	}

	bizNode := binary.Node{
		Tag: "biz",
		Content: []binary.Node{
			{
				Tag: "interactive",
				Attrs: map[string]interface{}{"type": "native_flow", "v": "1"},
				Content: []binary.Node{
					{
						Tag:   "native_flow",
						Attrs: map[string]interface{}{"v": "9", "name": "mixed"},
					},
				},
			},
		},
	}
	nodes := []binary.Node{bizNode}
	client.SendMessage(context.Background(), v.Info.Chat, msg, whatsmeow.SendRequestExtra{AdditionalNodes: &nodes})
}
