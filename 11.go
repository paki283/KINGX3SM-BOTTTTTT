package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/types/events"
)


func handleGenericTextAPI(client *whatsmeow.Client, v *events.Message, args string, endpoint string, paramName string, cmdName string) {
	if args == "" {
		replyMessage(client, v, fmt.Sprintf("❌ Please provide a query.\nExample: `.%s something`", cmdName))
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "⏳")

	apiURL := fmt.Sprintf("%s%s?%s=%s", GlobalConfig.APIBaseURL, endpoint, paramName, url.QueryEscape(args))
	resp, err := http.Get(apiURL)
	if err != nil {
		replyMessage(client, v, "❌ API Request Failed.")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		
		if data, ok := result["data"].(string); ok {
			replyMessage(client, v, data)
			react(client, v.Info.Chat, v.Info.ID, "✅")
			return
		} else if result["result"] != nil {
			strRes := fmt.Sprintf("%v", result["result"])
			replyMessage(client, v, strRes)
			react(client, v.Info.Chat, v.Info.ID, "✅")
			return
		}
	}

	replyMessage(client, v, "❌ Failed to parse response from server.")
}

func handleAnime(client *whatsmeow.Client, v *events.Message, args string) {
    handleGenericTextAPI(client, v, args, "/anime/animedetail", "url", "anime")
}

func handleManga(client *whatsmeow.Client, v *events.Message, args string) {
    handleGenericTextAPI(client, v, args, "/anime/manga-detail", "id", "manga")
}

func handleMovie(client *whatsmeow.Client, v *events.Message, args string) {
    handleGenericTextAPI(client, v, args, "/detail", "id", "movie")
}

func handleGithub(client *whatsmeow.Client, v *events.Message, args string) {
    handleGenericTextAPI(client, v, args, "/search/githubuser", "query", "github")
}


func handleRandomMedia(client *whatsmeow.Client, v *events.Message, endpoint string, isVideo bool) {
	react(client, v.Info.Chat, v.Info.ID, "⏳")

	apiURL := fmt.Sprintf("%s%s", GlobalConfig.APIBaseURL, endpoint)
	resp, err := http.Get(apiURL)
	if err != nil {
		replyMessage(client, v, "❌ Failed to fetch media.")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if dataURL, ok := result["data"].(string); ok {
			go downloadAndSend(client, v, dataURL, "video") 
			return
		} else if resultURL, ok := result["url"].(string); ok {
		    go downloadAndSend(client, v, resultURL, "video")
			return
		}
	}
	replyMessage(client, v, "❌ Media not found or API changed.")
}

func handleWaifu(client *whatsmeow.Client, v *events.Message) {
    handleRandomMedia(client, v, "/random/waifu", false)
}

func handleElaina(client *whatsmeow.Client, v *events.Message) {
    handleRandomMedia(client, v, "/random/anime/elaina", false)
}

func handleMountain(client *whatsmeow.Client, v *events.Message) {
    handleRandomMedia(client, v, "/random/anime/mountain", false)
}

func handleHentai(client *whatsmeow.Client, v *events.Message) {
    if !v.Info.IsGroup {
        handleRandomMedia(client, v, "/nsfw/hentai", false)
    } else {
        replyMessage(client, v, "❌ NSFW commands are disabled in groups.")
    }
}

func handlePinterestSearch(client *whatsmeow.Client, v *events.Message, args string) {
    if args == "" {
		replyMessage(client, v, "❌ Please provide a search query.")
		return
	}

	react(client, v.Info.Chat, v.Info.ID, "⏳")

	apiURL := fmt.Sprintf("%s/search/pinterest?query=%s", GlobalConfig.APIBaseURL, url.QueryEscape(args))
	resp, err := http.Get(apiURL)
	if err != nil {
		replyMessage(client, v, "❌ API Request Failed.")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if data, ok := result["data"].([]interface{}); ok && len(data) > 0 {
		    limit := len(data)
		    if limit > 3 { limit = 3 }
		    for i := 0; i < limit; i++ {
		        imgStr, _ := data[i].(string)
		        go downloadAndSend(client, v, imgStr, "video")
		    }
			return
		}
	}
	replyMessage(client, v, "❌ No results found.")
}

func handleCodeExplanation(client *whatsmeow.Client, v *events.Message, args string) {
    if args == "" {
		replyMessage(client, v, "❌ Please provide code to explain.")
		return
	}
    react(client, v.Info.Chat, v.Info.ID, "🧠")

	apiURL := fmt.Sprintf("%s/ai/explaincode?code=%s", GlobalConfig.APIBaseURL, url.QueryEscape(args))
	resp, err := http.Get(apiURL)
	if err != nil {
		replyMessage(client, v, "❌ API Request Failed.")
		return
	}
	defer resp.Body.Close()

	var result map[string]interface{}
	if err := json.NewDecoder(resp.Body).Decode(&result); err == nil {
		if data, ok := result["data"].(string); ok {
			replyMessage(client, v, data)
			react(client, v.Info.Chat, v.Info.ID, "✅")
			return
		}
	}
	replyMessage(client, v, "❌ Failed to explain code.")
}











































































































































































































































































































































































































































































































































































































func handleAPI_3d(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/3d", "prompt", "3d")
}

func handleAPI_ai4chat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/ai4chat", "prompt", "ai4chat")
}

func handleAPI_aiappchat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/aiappchat", "prompt", "aiappchat")
}

func handleAPI_aiappgen(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/aiappgen", "prompt", "aiappgen")
}

func handleAPI_dalle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/dalle", "prompt", "dalle")
}

func handleAPI_aichat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/aichat", "prompt", "aichat")
}

func handleAPI_aiserv(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/aiserv", "prompt", "aiserv")
}

func handleAPI_quick(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/quick", "text", "quick")
}

func handleAPI_advanced(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/advanced", "text", "advanced")
}

func handleAPI_animekill(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/animekill", "image", "animekill")
}

func handleAPI_blackbox(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/blackbox", "prompt", "blackbox")
}

func handleAPI_borli(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/borli", "action", "borli")
}

func handleAPI_cartoon(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/cartoon", "imageurl", "cartoon")
}

func handleAPI_copilot(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/copilot", "text", "copilot")
}

func handleAPI_copilotthink(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/copilot-think", "text", "copilotthink")
}

func handleAPI_gpt5(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/gpt-5", "text", "gpt5")
}

func handleAPI_ch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/ch", "q", "ch")
}

func handleAPI_chatbot(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatbot", "text", "chatbot")
}

func handleAPI_chatevo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatevo", "prompt", "chatevo")
}

func handleAPI_chatex(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatex", "text", "chatex")
}

func handleAPI_chatup(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatup", "prompt", "chatup")
}

func handleAPI_prompttocode(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/prompttocode", "prompt", "prompttocode")
}

func handleAPI_detectbugs(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/detectbugs", "code", "detectbugs")
}

func handleAPI_convertcode(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/convertcode", "code", "convertcode")
}

func handleAPI_explaincode(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/explaincode", "code", "explaincode")
}

func handleAPI_chateverywhere(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chateverywhere", "text", "chateverywhere")
}

func handleAPI_chateverywherereset(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chateverywhere-reset", "userId", "chateverywherereset")
}

func handleAPI_deepquery(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/deepquery", "prompt", "deepquery")
}

func handleAPI_logical(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/logical", "text", "logical")
}

func handleAPI_creative(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/creative", "text", "creative")
}

func handleAPI_summarize(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/summarize", "text", "summarize")
}

func handleAPI_codebeginner(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/code-beginner", "text", "codebeginner")
}

func handleAPI_codeadvanced(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/code-advanced", "text", "codeadvanced")
}

func handleAPI_dream(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/dream", "dream", "dream")
}

func handleAPI_deepseekchat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/deepseekchat", "prompt", "deepseekchat")
}

func handleAPI_deepseekreasoner(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/deepseekreasoner", "prompt", "deepseekreasoner")
}

func handleAPI_reset(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/ai/reset", "query", "reset")
}

func handleAPI_history(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/ai/history", "query", "history")
}

func handleAPI_easemategenerate(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/easemate-generate", "prompt", "easemategenerate")
}

func handleAPI_easematechat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/easemate-chat", "prompt", "easematechat")
}

func handleAPI_homeplannerchat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/homeplanner-chat", "prompt", "homeplannerchat")
}

func handleAPI_homeplannerimage(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/homeplanner-image", "prompt", "homeplannerimage")
}

func handleAPI_homeplannertts(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/homeplanner-tts", "prompt", "homeplannertts")
}

func handleAPI_homeplannersearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/homeplanner-search", "prompt", "homeplannersearch")
}

func handleAPI_homeplanneryt(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/homeplanner-yt", "prompt", "homeplanneryt")
}

func handleAPI_img2img(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/img2img", "imageUrl", "img2img")
}

func handleAPI_lumo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/lumo", "q", "lumo")
}

func handleAPI_chatmusiclyrics(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatmusic-lyrics", "prompt", "chatmusiclyrics")
}

func handleAPI_chatmusiccreate(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatmusic-create", "prompt", "chatmusiccreate")
}

func handleAPI_chatmusicstatus(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/chatmusic-status", "state", "chatmusicstatus")
}

func handleAPI_mydreams(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/mydreams", "prompt", "mydreams")
}

func handleAPI_olabiba(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/olabiba", "prompt", "olabiba")
}

func handleAPI_photogpt(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/photogpt", "prompt", "photogpt")
}

func handleAPI_photonex(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/photonex", "image", "photonex")
}

func handleAPI_soraremover(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/soraremover", "url", "soraremover")
}

func handleAPI_txt2img(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/txt2img", "prompt", "txt2img")
}

func handleAPI_unlimai(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/unlimai", "image", "unlimai")
}

func handleAPI_video(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/video", "prompt", "video")
}

func handleAPI_aiwriterchat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/aiwriter-chat", "prompt", "aiwriterchat")
}

func handleAPI_aiwriterimage(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/aiwriter-image", "prompt", "aiwriterimage")
}

func handleAPI_aiwritermodels(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/ai/aiwriter-models", "query", "aiwritermodels")
}

func handleAPI_realistic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/realistic", "prompt", "realistic")
}

func handleAPI_anime(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/anime", "prompt", "anime")
}

func handleAPI_fantasy(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/fantasy", "prompt", "fantasy")
}

func handleAPI_cyberpunk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/cyberpunk", "prompt", "cyberpunk")
}

func handleAPI_watercolor(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/watercolor", "prompt", "watercolor")
}

func handleAPI_oilpainting(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/oil-painting", "prompt", "oilpainting")
}

func handleAPI_pixelart(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/pixel-art", "prompt", "pixelart")
}

func handleAPI_sketch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/sketch", "prompt", "sketch")
}

func handleAPI_abstract(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/abstract", "prompt", "abstract")
}

func handleAPI_minimalist(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/minimalist", "prompt", "minimalist")
}

func handleAPI_surreal(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/surreal", "prompt", "surreal")
}

func handleAPI_vintage(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/vintage", "prompt", "vintage")
}

func handleAPI_steampunk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/steampunk", "prompt", "steampunk")
}

func handleAPI_horror(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/horror", "prompt", "horror")
}

func handleAPI_scifi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/sci-fi", "prompt", "scifi")
}

func handleAPI_popart(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ai/pop-art", "prompt", "popart")
}

func handleAPI_animekillhome(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/animekill-home", "query", "animekillhome")
}

func handleAPI_animekillhomestatic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/animekill-homestatic", "query", "animekillhomestatic")
}

func handleAPI_animekillsearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animekill-search", "query", "animekillsearch")
}

func handleAPI_animekilldetail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animekill-detail", "anime_id", "animekilldetail")
}

func handleAPI_animekillepisodes(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animekill-episodes", "anime_id", "animekillepisodes")
}

func handleAPI_animekillstream(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animekill-stream", "anime_id", "animekillstream")
}

func handleAPI_animekillcomments(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animekill-comments", "anime_id", "animekillcomments")
}

func handleAPI_animekillbygenre(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animekill-bygenre", "genre", "animekillbygenre")
}

func handleAPI_animekillgenres(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/animekill-genres", "query", "animekillgenres")
}

func handleAPI_animekillschedule(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/animekill-schedule", "query", "animekillschedule")
}

func handleAPI_animesearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animesearch", "query", "animesearch")
}

func handleAPI_animedetail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animedetail", "url", "animedetail")
}

func handleAPI_animedownload(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/animedownload", "url", "animedownload")
}

func handleAPI_hanimesearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/hanime-search", "query", "hanimesearch")
}

func handleAPI_hanimedetail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/hanime-detail", "id", "hanimedetail")
}

func handleAPI_mangahome(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-home", "page", "mangahome")
}

func handleAPI_mangasearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-search", "query", "mangasearch")
}

func handleAPI_mangadetail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-detail", "id", "mangadetail")
}

func handleAPI_mangachapter(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-chapter", "chapter_id", "mangachapter")
}

func handleAPI_mangasuggestions(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-suggestions", "suggestion_type", "mangasuggestions")
}

func handleAPI_mangaepisodes(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-episodes", "id", "mangaepisodes")
}

func handleAPI_mangaseries(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-series", "id", "mangaseries")
}

func handleAPI_mangacomments(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/manga-comments", "content_id", "mangacomments")
}

func handleAPI_mangarankfilters(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/manga-rank-filters", "query", "mangarankfilters")
}

func handleAPI_mangaranktags(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/manga-rank-tags", "query", "mangaranktags")
}

func handleAPI_hug(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/hug", "query", "hug")
}

func handleAPI_slap(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/slap", "query", "slap")
}

func handleAPI_pat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/pat", "query", "pat")
}

func handleAPI_cry(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/cry", "query", "cry")
}

func handleAPI_kill(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/kill", "query", "kill")
}

func handleAPI_bite(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/bite", "query", "bite")
}

func handleAPI_yeet(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/yeet", "query", "yeet")
}

func handleAPI_bully(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/bully", "query", "bully")
}

func handleAPI_bonk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/bonk", "query", "bonk")
}

func handleAPI_wink(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/wink", "query", "wink")
}

func handleAPI_poke(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/poke", "query", "poke")
}

func handleAPI_nom(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/nom", "query", "nom")
}

func handleAPI_smile(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/smile", "query", "smile")
}

func handleAPI_wave(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/wave", "query", "wave")
}

func handleAPI_awoo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/awoo", "query", "awoo")
}

func handleAPI_blush(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/blush", "query", "blush")
}

func handleAPI_smug(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/smug", "query", "smug")
}

func handleAPI_glomp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/glomp", "query", "glomp")
}

func handleAPI_happy(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/happy", "query", "happy")
}

func handleAPI_dance(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/dance", "query", "dance")
}

func handleAPI_cringe(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/cringe", "query", "cringe")
}

func handleAPI_cuddle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/cuddle", "query", "cuddle")
}

func handleAPI_highfive(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/highfive", "query", "highfive")
}

func handleAPI_handhold(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/handhold", "query", "handhold")
}

func handleAPI_shinobu(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/shinobu", "query", "shinobu")
}

func handleAPI_reactions(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/reactions", "query", "reactions")
}

func handleAPI_rule34home(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/rule34-home", "query", "rule34home")
}

func handleAPI_rule34search(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/rule34-search", "query", "rule34search")
}

func handleAPI_rule34detail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/rule34-detail", "url", "rule34detail")
}

func handleAPI_webnovelhot(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/anime/webnovel-hot", "query", "webnovelhot")
}

func handleAPI_webnovelrank(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/webnovel-rank", "page", "webnovelrank")
}

func handleAPI_webnovelsearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/webnovel-search", "query", "webnovelsearch")
}

func handleAPI_webnoveldetail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/webnovel-detail", "bid", "webnoveldetail")
}

func handleAPI_webnovelchapter(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/anime/webnovel-chapter", "bid", "webnovelchapter")
}

func handleAPI_aio(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/aio", "url", "aio")
}

func handleAPI_capcut(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/capcut", "url", "capcut")
}

func handleAPI_doods(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/doods", "url", "doods")
}

func handleAPI_douyin(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/douyin", "url", "douyin")
}

func handleAPI_facebook(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/facebook", "url", "facebook")
}

func handleAPI_facebookv2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/facebookv2", "url", "facebookv2")
}

func handleAPI_ig2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/ig2", "url", "ig2")
}

func handleAPI_instagram(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/instagram", "url", "instagram")
}

func handleAPI_mediafire(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/mediafire", "url", "mediafire")
}

func handleAPI_pinterest(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/pinterest", "url", "pinterest")
}

func handleAPI_pinterestv2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/pinterestV2", "url", "pinterestv2")
}

func handleAPI_rednote(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/rednote", "url", "rednote")
}

func handleAPI_rednotemedia(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/rednote-media", "url", "rednotemedia")
}

func handleAPI_saveweb2zip(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/saveweb2zip", "url", "saveweb2zip")
}

func handleAPI_sfile(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/sfile", "url", "sfile")
}

func handleAPI_snackvideo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/snackvideo", "url", "snackvideo")
}

func handleAPI_soundcloud(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/soundcloud", "url", "soundcloud")
}

func handleAPI_spotify(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/spotify", "url", "spotify")
}

func handleAPI_spotifyv2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/spotifyv2", "url", "spotifyv2")
}

func handleAPI_terabox(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/terabox", "url", "terabox")
}

func handleAPI_threads(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/threads", "url", "threads")
}

func handleAPI_threadsv2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/threadsV2", "url", "threadsv2")
}

func handleAPI_tiktok(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/tiktok", "url", "tiktok")
}

func handleAPI_tiktokv2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/tiktokV2", "url", "tiktokv2")
}

func handleAPI_tiktokv3(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/tiktokV3", "url", "tiktokv3")
}

func handleAPI_tiktokvideo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/tiktokvideo", "url", "tiktokvideo")
}

func handleAPI_tiktokslide(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/tiktokslide", "url", "tiktokslide")
}

func handleAPI_twitter(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/twitter", "url", "twitter")
}

func handleAPI_youtubeaudio(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/youtube-audio", "url", "youtubeaudio")
}

func handleAPI_youtubevideo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/download/youtube-video", "url", "youtubevideo")
}

func handleAPI_quizcategories(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/game/quizcategories", "query", "quizcategories")
}

func handleAPI_quizguess(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/game/quizguess", "level", "quizguess")
}

func handleAPI_quizpuzzle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/game/quizpuzzle", "level", "quizpuzzle")
}

func handleAPI_quiztruefalse(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/game/quiztruefalse", "level", "quiztruefalse")
}

func handleAPI_quizrandom(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/game/quizrandom", "level", "quizrandom")
}

func handleAPI_image(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/image", "text", "image")
}

func handleAPI_gif(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/gif", "text", "gif")
}

func handleAPI_mp4(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/mp4", "text", "mp4")
}

func handleAPI_meme(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/meme", "topText", "meme")
}

func handleAPI_memetext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/memeText", "text", "memetext")
}

func handleAPI_spongebob(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/spongebob", "text", "spongebob")
}

func handleAPI_ttp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/imagecreator/ttp", "text", "ttp")
}

func handleAPI_moviesearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/moviesearch", "query", "moviesearch")
}

func handleAPI_moviedetail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/moviedetail", "url", "moviedetail")
}

func handleAPI_search(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search", "q", "search")
}

func handleAPI_suggest(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/suggest", "q", "suggest")
}

func handleAPI_detail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/detail", "id", "detail")
}

func handleAPI_recommendations(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/recommendations", "id", "recommendations")
}

func handleAPI_trending(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/trending", "tabId", "trending")
}

func handleAPI_home(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/home", "query", "home")
}

func handleAPI_countries(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/countries", "query", "countries")
}

func handleAPI_hentaisfm(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/hentai-sfm", "query", "hentaisfm")
}

func handleAPI_ass(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/ass", "query", "ass")
}

func handleAPI_sixtynine(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/sixtynine", "query", "sixtynine")
}

func handleAPI_pussy(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/pussy", "query", "pussy")
}

func handleAPI_dick(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/dick", "query", "dick")
}

func handleAPI_anal(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/anal", "query", "anal")
}

func handleAPI_boobs(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/boobs", "query", "boobs")
}

func handleAPI_bdsm(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/bdsm", "query", "bdsm")
}

func handleAPI_black(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/black", "query", "black")
}

func handleAPI_easter(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/easter", "query", "easter")
}

func handleAPI_bottomless(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/bottomless", "query", "bottomless")
}

func handleAPI_collared(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/collared", "query", "collared")
}

func handleAPI_cum(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/cum", "query", "cum")
}

func handleAPI_cumsluts(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/cumsluts", "query", "cumsluts")
}

func handleAPI_dom(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/dom", "query", "dom")
}

func handleAPI_extreme(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/extreme", "query", "extreme")
}

func handleAPI_feet(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/feet", "query", "feet")
}

func handleAPI_finger(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/finger", "query", "finger")
}

func handleAPI_fuck(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/fuck", "query", "fuck")
}

func handleAPI_futa(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/futa", "query", "futa")
}

func handleAPI_gay(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/gay", "query", "gay")
}

func handleAPI_hentai(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/hentai", "query", "hentai")
}

func handleAPI_kiss(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/kiss", "query", "kiss")
}

func handleAPI_lick(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/lick", "query", "lick")
}

func handleAPI_pegged(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/pegged", "query", "pegged")
}

func handleAPI_phgif(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/phgif", "query", "phgif")
}

func handleAPI_puffies(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/puffies", "query", "puffies")
}

func handleAPI_real(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/real", "query", "real")
}

func handleAPI_suck(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/suck", "query", "suck")
}

func handleAPI_tattoo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/tattoo", "query", "tattoo")
}

func handleAPI_tiny(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/tiny", "query", "tiny")
}

func handleAPI_toys(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/toys", "query", "toys")
}

func handleAPI_xmas(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/nsfw/xmas", "query", "xmas")
}

func handleAPI_xnxxsearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/nsfw/xnxx-search", "query", "xnxxsearch")
}

func handleAPI_xvideossearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/nsfw/xvideos-search", "query", "xvideossearch")
}

func handleAPI_xnxxdl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/nsfw/xnxx-dl", "url", "xnxxdl")
}

func handleAPI_xvideosdl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/nsfw/xvideos-dl", "url", "xvideosdl")
}

func handleAPI_anhsfw(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anhsfw", false)
}

func handleAPI_anhmoe(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anhmoe", false)
}

func handleAPI_anhai(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anhai", false)
}

func handleAPI_anhnsfw(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anhnsfw", false)
}

func handleAPI_anhhentai(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anhhentai", false)
}

func handleAPI_anhvideonsfw(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anhvideonsfw", false)
}

func handleAPI_bluearchive(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/bluearchive", false)
}

func handleAPI_boypic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/boypic", false)
}

func handleAPI_car(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/car", false)
}

func handleAPI_cat(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/cat", false)
}

func handleAPI_chinagirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/chinagirl", false)
}

func handleAPI_dog(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/dog", false)
}

func handleAPI_randomgirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/randomgirl", false)
}

func handleAPI_hijabgirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/hijabgirl", false)
}

func handleAPI_indonesiagirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/indonesiagirl", false)
}

func handleAPI_japangirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/japangirl", false)
}

func handleAPI_koreangirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/koreangirl", false)
}

func handleAPI_loli(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/loli", false)
}

func handleAPI_malaysiagirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/malaysiagirl", false)
}

func handleAPI_profilepics(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/profilepics", false)
}

func handleAPI_thailandgirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/thailandgirl", false)
}

func handleAPI_tiktokgirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/tiktokgirl", false)
}

func handleAPI_vietnamgirl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/vietnamgirl", false)
}

func handleAPI_waifu(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/waifu", false)
}

func handleAPI_akiyama(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/akiyama", false)
}

func handleAPI_ana(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/ana", false)
}

func handleAPI_asuna(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/asuna", false)
}

func handleAPI_ayuzawa(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/ayuzawa", false)
}

func handleAPI_boruto(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/boruto", false)
}

func handleAPI_chitoge(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/chitoge", false)
}

func handleAPI_deidara(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/deidara", false)
}

func handleAPI_doraemon(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/doraemon", false)
}

func handleAPI_elaina(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/elaina", false)
}

func handleAPI_emilia(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/emilia", false)
}

func handleAPI_erza(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/erza", false)
}

func handleAPI_hestia(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/hestia", false)
}

func handleAPI_husbu(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/husbu", false)
}

func handleAPI_inori(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/inori", false)
}

func handleAPI_itachi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/itachi", false)
}

func handleAPI_kagura(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/kagura", false)
}

func handleAPI_kaori(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/kaori", false)
}

func handleAPI_keneki(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/keneki", false)
}

func handleAPI_kotori(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/kotori", false)
}

func handleAPI_kurumi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/kurumi", false)
}

func handleAPI_madara(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/madara", false)
}

func handleAPI_megumin(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/megumin", false)
}

func handleAPI_mikasa(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/mikasa", false)
}

func handleAPI_miku(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/miku", false)
}

func handleAPI_minato(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/minato", false)
}

func handleAPI_naruto(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/naruto", false)
}

func handleAPI_nekonime(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/nekonime", false)
}

func handleAPI_nezuko(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/nezuko", false)
}

func handleAPI_onepiece(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/onepiece", false)
}

func handleAPI_rize(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/rize", false)
}

func handleAPI_sagiri(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/sagiri", false)
}

func handleAPI_sakura(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/sakura", false)
}

func handleAPI_sasuke(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/sasuke", false)
}

func handleAPI_shinomiya(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/shinomiya", false)
}

func handleAPI_tsunade(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/tsunade", false)
}

func handleAPI_yotsuba(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/yotsuba", false)
}

func handleAPI_yuki(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/yuki", false)
}

func handleAPI_yumeko(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/yumeko", false)
}

func handleAPI_art(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/art", false)
}

func handleAPI_cyber(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/cyber", false)
}

func handleAPI_gamewallpaper(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/gamewallpaper", false)
}

func handleAPI_mountain(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/mountain", false)
}

func handleAPI_programming(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/programming", false)
}

func handleAPI_space(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/space", false)
}

func handleAPI_technology(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/technology", false)
}

func handleAPI_wallhp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/wallhp", false)
}

func handleAPI_wallml(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/wallml", false)
}

func handleAPI_wallmlnime(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleRandomMedia(client, v, "/random/anime/wallmlnime", false)
}

func handleAPI_android1(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/android1", "q", "android1")
}

func handleAPI_applemusic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/applemusic", "q", "applemusic")
}

func handleAPI_cuaca(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/cuaca", "kota", "cuaca")
}

func handleAPI_repos(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/repos", "query", "repos")
}

func handleAPI_users(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/users", "query", "users")
}

func handleAPI_issues(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/issues", "query", "issues")
}

func handleAPI_code(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/code", "query", "code")
}

func handleAPI_imdb(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/imdb", "query", "imdb")
}

func handleAPI_lyrics(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/lyrics", "title", "lyrics")
}

func handleAPI_nik(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/nik", "q", "nik")
}

func handleAPI_wallpaper(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/wallpaper", "query", "wallpaper")
}

func handleAPI_telegram(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/telegram", "query", "telegram")
}

func handleAPI_tggroup(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/tggroup", "query", "tggroup")
}

func handleAPI_tiktoksearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/tiktoksearch", "q", "tiktoksearch")
}

func handleAPI_wagroup(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/wagroup", "query", "wagroup")
}

func handleAPI_youtube(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/youtube", "q", "youtube")
}

func handleAPI_ytmonet(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/search/ytmonet", "url", "ytmonet")
}

func handleAPI_download(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/sound/download", "url", "download")
}

func handleAPI_nonstick(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/sound/nonstick", "type", "nonstick")
}

func handleAPI_football(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/sports/football", "detail", "football")
}

func handleAPI_basketball(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/sports/basketball", "detail", "basketball")
}

func handleAPI_othersports(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/sports/othersports", "detail", "othersports")
}

func handleAPI_webss(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ssweb/webss", "url", "webss")
}

func handleAPI_apiflash(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ssweb/apiFlash", "url", "apiflash")
}

func handleAPI_screenshotlayer(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/ssweb/screenshotLayer", "url", "screenshotlayer")
}

func handleAPI_ffstalk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/stalk/ffstalk", "id", "ffstalk")
}

func handleAPI_igstalk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/stalk/igstalk", "user", "igstalk")
}

func handleAPI_igstalkv2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/stalk/igstalkV2", "user", "igstalkv2")
}

func handleAPI_ttstalk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/stalk/ttstalk", "user", "ttstalk")
}

func handleAPI_twitterstalk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/stalk/twitterstalk", "user", "twitterstalk")
}

func handleAPI_ytstalk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/stalk/ytstalk", "user", "ytstalk")
}

func handleAPI_glitchtext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/glitchtext", "text", "glitchtext")
}

func handleAPI_writetext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/writetext", "text", "writetext")
}

func handleAPI_advancedglow(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/advancedglow", "text", "advancedglow")
}

func handleAPI_typographytext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/typographytext", "text", "typographytext")
}

func handleAPI_pixelglitch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/pixelglitch", "text", "pixelglitch")
}

func handleAPI_neonglitch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/neonglitch", "text", "neonglitch")
}

func handleAPI_flagtext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/flagtext", "text", "flagtext")
}

func handleAPI_flag3dtext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/flag3dtext", "text", "flag3dtext")
}

func handleAPI_deletingtext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/deletingtext", "text", "deletingtext")
}

func handleAPI_blackpinkstyle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/blackpinkstyle", "text", "blackpinkstyle")
}

func handleAPI_glowingtext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/glowingtext", "text", "glowingtext")
}

func handleAPI_underwatertext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/underwatertext", "text", "underwatertext")
}

func handleAPI_logomaker(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/logomaker", "text", "logomaker")
}

func handleAPI_cartoonstyle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/cartoonstyle", "text", "cartoonstyle")
}

func handleAPI_papercutstyle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/papercutstyle", "text", "papercutstyle")
}

func handleAPI_watercolortext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/watercolortext", "text", "watercolortext")
}

func handleAPI_effectclouds(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/effectclouds", "text", "effectclouds")
}

func handleAPI_blackpinklogo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/blackpinklogo", "text", "blackpinklogo")
}

func handleAPI_gradienttext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/gradienttext", "text", "gradienttext")
}

func handleAPI_summerbeach(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/summerbeach", "text", "summerbeach")
}

func handleAPI_luxurygold(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/luxurygold", "text", "luxurygold")
}

func handleAPI_multicoloredneon(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/multicoloredneon", "text", "multicoloredneon")
}

func handleAPI_sandsummer(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/sandsummer", "text", "sandsummer")
}

func handleAPI_galaxywallpaper(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/galaxywallpaper", "text", "galaxywallpaper")
}

func handleAPI_style1917(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/style1917", "text", "style1917")
}

func handleAPI_makingneon(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/makingneon", "text", "makingneon")
}

func handleAPI_royaltext(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/royaltext", "text", "royaltext")
}

func handleAPI_freecreate(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/freecreate", "text", "freecreate")
}

func handleAPI_galaxystyle(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/galaxystyle", "text", "galaxystyle")
}

func handleAPI_lighteffects(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/lighteffects", "text", "lighteffects")
}

func handleAPI_sendemail(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/sendemail", "to", "sendemail")
}

func handleAPI_codeanalyzer(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/codeanalyzer", "code", "codeanalyzer")
}

func handleAPI_codeconverter(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/codeconverter", "code", "codeconverter")
}

func handleAPI_tojavascript(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tojavascript", "code", "tojavascript")
}

func handleAPI_topython(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/topython", "code", "topython")
}

func handleAPI_tojava(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tojava", "code", "tojava")
}

func handleAPI_tocpp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tocpp", "code", "tocpp")
}

func handleAPI_tophp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tophp", "code", "tophp")
}

func handleAPI_compiler(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compiler", "code", "compiler")
}

func handleAPI_compilejs(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compilejs", "code", "compilejs")
}

func handleAPI_compilepython(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compilepython", "code", "compilepython")
}

func handleAPI_compilejava(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compilejava", "code", "compilejava")
}

func handleAPI_compilec(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compilec", "code", "compilec")
}

func handleAPI_compilecpp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compilecpp", "code", "compilecpp")
}

func handleAPI_compilecsharp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/compilecsharp", "code", "compilecsharp")
}

func handleAPI_emojiencrypt(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/emoji-encrypt", "input", "emojiencrypt")
}

func handleAPI_emojidecrypt(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/emoji-decrypt", "input", "emojidecrypt")
}

func handleAPI_htmlecnc(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/htmlecnc", "html", "htmlecnc")
}

func handleAPI_htmlbasic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/htmlbasic", "html", "htmlbasic")
}

func handleAPI_htmlextended(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/htmlextended", "html", "htmlextended")
}

func handleAPI_htmlhigh(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/htmlhigh", "html", "htmlhigh")
}

func handleAPI_htmlmaximum(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/htmlmaximum", "html", "htmlmaximum")
}

func handleAPI_fdroidsearch(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/fdroidsearch", "q", "fdroidsearch")
}

func handleAPI_fdroidpackage(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/fdroidpackage", "url", "fdroidpackage")
}

func handleAPI_fdroidapp(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/fdroidapp", "package", "fdroidapp")
}

func handleAPI_geoip(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/geoip", "ip", "geoip")
}

func handleAPI_myip(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/tools/myip", "query", "myip")
}

func handleAPI_hostcheck(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/hostcheck", "domain", "hostcheck")
}

func handleAPI_hostchecksimple(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/hostchecksimple", "domain", "hostchecksimple")
}

func handleAPI_html2img(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/html2img", "html", "html2img")
}

func handleAPI_html2imgdirect(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/html2imgdirect", "html", "html2imgdirect")
}

func handleAPI_obflow(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/obflow", "code", "obflow")
}

func handleAPI_obfmedium(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/obfmedium", "code", "obfmedium")
}

func handleAPI_obfhigh(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/obfhigh", "code", "obfhigh")
}

func handleAPI_obfextreme(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/obfextreme", "code", "obfextreme")
}

func handleAPI_tiktoktranscript(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tiktoktranscript", "url", "tiktoktranscript")
}

func handleAPI_entoid(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/entoid", "text", "entoid")
}

func handleAPI_idtoen(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/idtoen", "text", "idtoen")
}

func handleAPI_jatoid(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/jatoid", "text", "jatoid")
}

func handleAPI_kotoid(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/kotoid", "text", "kotoid")
}

func handleAPI_zhtoid(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/zhtoid", "text", "zhtoid")
}

func handleAPI_artoid(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/artoid", "text", "artoid")
}

func handleAPI_detectlanguage(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/detectlanguage", "text", "detectlanguage")
}

func handleAPI_languages(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/tools/languages", "query", "languages")
}

func handleAPI_youtubetranscript(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/youtube-transcript", "url", "youtubetranscript")
}

func handleAPI_dagd(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/dagd", "url", "dagd")
}

func handleAPI_vgd(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/vgd", "url", "vgd")
}

func handleAPI_tinube(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tinube", "url", "tinube")
}

func handleAPI_spoome(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/spooMe", "url", "spoome")
}

func handleAPI_spooemoji(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/spooEmoji", "url", "spooemoji")
}

func handleAPI_random(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/random", "url", "random")
}

func handleAPI_allstyles(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/allstyles", "text", "allstyles")
}

func handleAPI_circled(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/circled", "text", "circled")
}

func handleAPI_circledneg(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/circledneg", "text", "circledneg")
}

func handleAPI_fullwidth(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/fullwidth", "text", "fullwidth")
}

func handleAPI_mathbold(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathbold", "text", "mathbold")
}

func handleAPI_mathboldfraktur(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathboldfraktur", "text", "mathboldfraktur")
}

func handleAPI_mathbolditalic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathbolditalic", "text", "mathbolditalic")
}

func handleAPI_mathboldscript(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathboldscript", "text", "mathboldscript")
}

func handleAPI_mathdoublestruck(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathdoublestruck", "text", "mathdoublestruck")
}

func handleAPI_mathmonospace(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathmonospace", "text", "mathmonospace")
}

func handleAPI_mathsans(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathsans", "text", "mathsans")
}

func handleAPI_mathsansbold(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathsansbold", "text", "mathsansbold")
}

func handleAPI_mathsansbolditalic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathsansbolditalic", "text", "mathsansbolditalic")
}

func handleAPI_mathsansitalic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathsansitalic", "text", "mathsansitalic")
}

func handleAPI_parenthesized(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/parenthesized", "text", "parenthesized")
}

func handleAPI_regionalindicator(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/regionalindicator", "text", "regionalindicator")
}

func handleAPI_squared(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/squared", "text", "squared")
}

func handleAPI_squaredneg(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/squaredneg", "text", "squaredneg")
}

func handleAPI_tag(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/tag", "text", "tag")
}

func handleAPI_acute(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/acute", "text", "acute")
}

func handleAPI_cjkthai(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/cjkthai", "text", "cjkthai")
}

func handleAPI_curvy1(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/curvy1", "text", "curvy1")
}

func handleAPI_curvy2(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/curvy2", "text", "curvy2")
}

func handleAPI_curvy3(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/curvy3", "text", "curvy3")
}

func handleAPI_fauxcyrillic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/fauxcyrillic", "text", "fauxcyrillic")
}

func handleAPI_fauxethiopic(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/fauxethiopic", "text", "fauxethiopic")
}

func handleAPI_mathfraktur(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/mathfraktur", "text", "mathfraktur")
}

func handleAPI_rockdots(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/rockdots", "text", "rockdots")
}

func handleAPI_smallcaps(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/smallcaps", "text", "smallcaps")
}

func handleAPI_stroked(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/stroked", "text", "stroked")
}

func handleAPI_subscript(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/subscript", "text", "subscript")
}

func handleAPI_superscript(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/superscript", "text", "superscript")
}

func handleAPI_inverted(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/inverted", "text", "inverted")
}

func handleAPI_reversed(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tools/reversed", "text", "reversed")
}

func handleAPI_ttsen(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-en", "text", "ttsen")
}

func handleAPI_ttsid(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-id", "text", "ttsid")
}

func handleAPI_ttses(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-es", "text", "ttses")
}

func handleAPI_ttsfr(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-fr", "text", "ttsfr")
}

func handleAPI_ttsde(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-de", "text", "ttsde")
}

func handleAPI_ttsit(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-it", "text", "ttsit")
}

func handleAPI_ttspt(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-pt", "text", "ttspt")
}

func handleAPI_ttsnl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-nl", "text", "ttsnl")
}

func handleAPI_ttspl(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-pl", "text", "ttspl")
}

func handleAPI_ttsru(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-ru", "text", "ttsru")
}

func handleAPI_ttsja(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-ja", "text", "ttsja")
}

func handleAPI_ttsko(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-ko", "text", "ttsko")
}

func handleAPI_ttszhcn(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-zhcn", "text", "ttszhcn")
}

func handleAPI_ttszhtw(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-zhtw", "text", "ttszhtw")
}

func handleAPI_ttsar(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-ar", "text", "ttsar")
}

func handleAPI_ttshi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-hi", "text", "ttshi")
}

func handleAPI_ttsth(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-th", "text", "ttsth")
}

func handleAPI_ttsvi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-vi", "text", "ttsvi")
}

func handleAPI_ttstr(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-tr", "text", "ttstr")
}

func handleAPI_ttssv(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-sv", "text", "ttssv")
}

func handleAPI_ttsno(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-no", "text", "ttsno")
}

func handleAPI_ttsda(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-da", "text", "ttsda")
}

func handleAPI_ttsfi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-fi", "text", "ttsfi")
}

func handleAPI_ttsel(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-el", "text", "ttsel")
}

func handleAPI_ttshe(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-he", "text", "ttshe")
}

func handleAPI_ttscs(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-cs", "text", "ttscs")
}

func handleAPI_ttshu(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-hu", "text", "ttshu")
}

func handleAPI_ttsro(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-ro", "text", "ttsro")
}

func handleAPI_ttsuk(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-uk", "text", "ttsuk")
}

func handleAPI_xiaofei(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaofei", "text", "xiaofei")
}

func handleAPI_xiaolei(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaolei", "text", "xiaolei")
}

func handleAPI_xiaojie(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaojie", "text", "xiaojie")
}

func handleAPI_xiaohua(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaohua", "text", "xiaohua")
}

func handleAPI_xiaofeng(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaofeng", "text", "xiaofeng")
}

func handleAPI_xiaoze(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoze", "text", "xiaoze")
}

func handleAPI_xiaoyuan(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoyuan", "text", "xiaoyuan")
}

func handleAPI_xiaozheng(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaozheng", "text", "xiaozheng")
}

func handleAPI_xiaoying(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoying", "text", "xiaoying")
}

func handleAPI_xiaoqing(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoqing", "text", "xiaoqing")
}

func handleAPI_xiaoxiang(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoxiang", "text", "xiaoxiang")
}

func handleAPI_xiaoyan(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoyan", "text", "xiaoyan")
}

func handleAPI_xianran(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xianran", "text", "xianran")
}

func handleAPI_xiaoxue(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoxue", "text", "xiaoxue")
}

func handleAPI_xiaoxuan(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoxuan", "text", "xiaoxuan")
}

func handleAPI_xiaolu(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaolu", "text", "xiaolu")
}

func handleAPI_xiaowei(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaowei", "text", "xiaowei")
}

func handleAPI_xiaozhe(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaozhe", "text", "xiaozhe")
}

func handleAPI_xiaohao(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaohao", "text", "xiaohao")
}

func handleAPI_xiaoyi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoyi", "text", "xiaoyi")
}

func handleAPI_xiaotao(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaotao", "text", "xiaotao")
}

func handleAPI_xiaoming(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xiaoming", "text", "xiaoming")
}

func handleAPI_david(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/david", "text", "david")
}

func handleAPI_layla(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/layla", "text", "layla")
}

func handleAPI_james(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/james", "text", "james")
}

func handleAPI_joey(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/joey", "text", "joey")
}

func handleAPI_jennifer(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/jennifer", "text", "jennifer")
}

func handleAPI_john(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/john", "text", "john")
}

func handleAPI_paul(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/paul", "text", "paul")
}

func handleAPI_xena(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/xena", "text", "xena")
}

func handleAPI_marcus(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/marcus", "text", "marcus")
}

func handleAPI_jacob(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/jacob", "text", "jacob")
}

func handleAPI_sam(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/sam", "text", "sam")
}

func handleAPI_camila(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/camila", "text", "camila")
}

func handleAPI_amy(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/amy", "text", "amy")
}

func handleAPI_quincy(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/quincy", "text", "quincy")
}

func handleAPI_sally(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/sally", "text", "sally")
}

func handleAPI_emma(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/emma", "text", "emma")
}

func handleAPI_ethan(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/ethan", "text", "ethan")
}

func handleAPI_michael(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/michael", "text", "michael")
}

func handleAPI_olivia(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/olivia", "text", "olivia")
}

func handleAPI_mia(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/mia", "text", "mia")
}

func handleAPI_jackson(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/jackson", "text", "jackson")
}

func handleAPI_matthew(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/matthew", "text", "matthew")
}

func handleAPI_sophia(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/sophia", "text", "sophia")
}

func handleAPI_owen(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/owen", "text", "owen")
}

func handleAPI_beatrice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/beatrice", "text", "beatrice")
}

func handleAPI_scott(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/scott", "text", "scott")
}

func handleAPI_ivy(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/ivy", "text", "ivy")
}

func handleAPI_eric(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/eric", "text", "eric")
}

func handleAPI_kevin(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/kevin", "text", "kevin")
}

func handleAPI_hannah(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/hannah", "text", "hannah")
}

func handleAPI_katrina(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/katrina", "text", "katrina")
}

func handleAPI_victor(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/victor", "text", "victor")
}

func handleAPI_justin(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/justin", "text", "justin")
}

func handleAPI_leo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/leo", "text", "leo")
}

func handleAPI_grace(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/grace", "text", "grace")
}

func handleAPI_casey(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/casey", "text", "casey")
}

func handleAPI_dylan(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/dylan", "text", "dylan")
}

func handleAPI_julie(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/julie", "text", "julie")
}

func handleAPI_thomas(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/thomas", "text", "thomas")
}

func handleAPI_freya(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/freya", "text", "freya")
}

func handleAPI_max(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/max", "text", "max")
}

func handleAPI_phoebe(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/phoebe", "text", "phoebe")
}

func handleAPI_noah(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/noah", "text", "noah")
}

func handleAPI_sophie(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/sophie", "text", "sophie")
}

func handleAPI_isla(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/isla", "text", "isla")
}

func handleAPI_theo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/theo", "text", "theo")
}

func handleAPI_ella(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/ella", "text", "ella")
}

func handleAPI_freddie(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/freddie", "text", "freddie")
}

func handleAPI_arthur(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/arthur", "text", "arthur")
}

func handleAPI_isabella(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/isabella", "text", "isabella")
}

func handleAPI_evie(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/evie", "text", "evie")
}

func handleAPI_william(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/william", "text", "william")
}

func handleAPI_henry(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/henry", "text", "henry")
}

func handleAPI_lily(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/lily", "text", "lily")
}

func handleAPI_charlie(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/charlie", "text", "charlie")
}

func handleAPI_ttsvoices(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/tts/tts-voices", "query", "ttsvoices")
}

func handleAPI_ttsadultfemale1americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-female--1-american-english-truvoice", "text", "ttsadultfemale1americanenglishtruvoice")
}

func handleAPI_ttsadultfemale2americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-female--2-american-english-truvoice", "text", "ttsadultfemale2americanenglishtruvoice")
}

func handleAPI_ttsadultmale1americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--1-american-english-truvoice", "text", "ttsadultmale1americanenglishtruvoice")
}

func handleAPI_ttsadultmale2americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--2-american-english-truvoice", "text", "ttsadultmale2americanenglishtruvoice")
}

func handleAPI_ttsadultmale3americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--3-american-english-truvoice", "text", "ttsadultmale3americanenglishtruvoice")
}

func handleAPI_ttsadultmale4americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--4-american-english-truvoice", "text", "ttsadultmale4americanenglishtruvoice")
}

func handleAPI_ttsadultmale5americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--5-american-english-truvoice", "text", "ttsadultmale5americanenglishtruvoice")
}

func handleAPI_ttsadultmale6americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--6-american-english-truvoice", "text", "ttsadultmale6americanenglishtruvoice")
}

func handleAPI_ttsadultmale7americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--7-american-english-truvoice", "text", "ttsadultmale7americanenglishtruvoice")
}

func handleAPI_ttsadultmale8americanenglishtruvoice(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-adult-male--8-american-english-truvoice", "text", "ttsadultmale8americanenglishtruvoice")
}

func handleAPI_ttsfemalewhisper(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-female-whisper", "text", "ttsfemalewhisper")
}

func handleAPI_ttsmalewhisper(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-male-whisper", "text", "ttsmalewhisper")
}

func handleAPI_ttsmary(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mary", "text", "ttsmary")
}

func handleAPI_ttsmaryfortelephone(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mary-for-telephone", "text", "ttsmaryfortelephone")
}

func handleAPI_ttsmaryinhall(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mary-in-hall", "text", "ttsmaryinhall")
}

func handleAPI_ttsmaryinspace(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mary-in-space", "text", "ttsmaryinspace")
}

func handleAPI_ttsmaryinstadium(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mary-in-stadium", "text", "ttsmaryinstadium")
}

func handleAPI_ttsmike(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mike", "text", "ttsmike")
}

func handleAPI_ttsmikefortelephone(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mike-for-telephone", "text", "ttsmikefortelephone")
}

func handleAPI_ttsmikeinhall(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mike-in-hall", "text", "ttsmikeinhall")
}

func handleAPI_ttsmikeinspace(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mike-in-space", "text", "ttsmikeinspace")
}

func handleAPI_ttsmikeinstadium(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-mike-in-stadium", "text", "ttsmikeinstadium")
}

func handleAPI_ttsrobosoftfive(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-robosoft-five", "text", "ttsrobosoftfive")
}

func handleAPI_ttsrobosoftfour(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-robosoft-four", "text", "ttsrobosoftfour")
}

func handleAPI_ttsrobosoftone(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-robosoft-one", "text", "ttsrobosoftone")
}

func handleAPI_ttsrobosoftsix(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-robosoft-six", "text", "ttsrobosoftsix")
}

func handleAPI_ttsrobosoftthree(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-robosoft-three", "text", "ttsrobosoftthree")
}

func handleAPI_ttsrobosofttwo(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-robosoft-two", "text", "ttsrobosofttwo")
}

func handleAPI_ttssam(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-sam", "text", "ttssam")
}

func handleAPI_ttsbonzi(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/tts/tts-bonzi", "text", "ttsbonzi")
}

func handleAPI_sms24countries(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/vnum/sms24-countries", "query", "sms24countries")
}

func handleAPI_sms24numbers(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/vnum/sms24-numbers", "country", "sms24numbers")
}

func handleAPI_sms24messages(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/vnum/sms24-messages", "number", "sms24messages")
}

func handleAPI_veepncountries(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, "test", "/vnum/veepn-countries", "query", "veepncountries")
}

func handleAPI_veepnnumbers(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/vnum/veepn-numbers", "country", "veepnnumbers")
}

func handleAPI_veepnmessages(client *whatsmeow.Client, v *events.Message, args string) {
    react(client, v.Info.Chat, v.Info.ID, "⏳")
    handleGenericTextAPI(client, v, args, "/vnum/veepn-messages", "country", "veepnmessages")
}