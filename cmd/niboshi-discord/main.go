package main

import (
    "fmt"
    "time"
    "github.com/bwmarrin/discordgo"
    "log"
    "strings"
    "os"
    "github.com/cancer/niboshi-discord/pkg/message"
    "github.com/joho/godotenv"
)

var(
    stopBot = make(chan bool)
    vcsession *discordgo.VoiceConnection
    HelloWorld = "!helloworld"
    ChannelInfo = "!channelinfo"
    Members = "!showmembers"
    ChannelVoiceJoin = "!vcjoin"
    ChannelVoiceLeave = "!vcleave"
    Genius = "天才"
    TimerStart = "!timerstart"
    TimerStop = "!timerstop"
    stopTimer = make(chan bool)
)

func main() {
    configure()

    // Discordのセッションを作成
    discord, err := discordgo.New()
    discord.Token = token()

    if err != nil {
        fmt.Println("Error logging in")
        fmt.Println(err)
    }

    // 全てのWSAPIイベントが発生したときのイベントハンドラを追加
    discord.AddHandler(onMessageCreate)
    // websocketを開いてlistening開始
    err = discord.Open()
    if err != nil {
        fmt.Println(err)
    }

    fmt.Println("Listening...")
    <- stopBot
    return
}

func configure() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
}

func token() string {
    return fmt.Sprintf("Bot %s", os.Getenv("BOT_TOKEN"))
}

func botName() string {
    return fmt.Sprintf("<@!%s>", os.Getenv("BOT_ID"))
}

func onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
    // チャンネル取得
    c, err := s.State.Channel(m.ChannelID)
    if err != nil {
        log.Println("Error getting channel: ", err)
        return
    }

    fmt.Printf("%20s %20s %20s > %s\n", m.ChannelID, time.Now().Format(time.Stamp), m.Author.Username, m.Content)

    switch {
        // Bot宛に !helloworld コマンドが実行されたとき
        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", botName(), HelloWorld)):
            sendMessage(s, c, "Hello world!")

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", botName(), ChannelInfo)):
            guildChannels, _ := s.GuildChannels(c.GuildID)
            var sendText string
            for _, a := range guildChannels{
                sendText += fmt.Sprintf("%vチャンネルの%v(IDは%v)\n", a.Type, a.Name, a.ID)
            }
            sendMessage(s, c, sendText)

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", botName(), Members)):
            members, _ := s.GuildMembers(c.GuildID, "", 100)
            var sendText string
            for _, m := range members{
                sendText += fmt.Sprintf("%v\n", m.User.Username)
            }
            sendMessage(s, c, sendText)

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", botName(), TimerStart)):
            startTimer(func() {
                sendMessage(s, c, fmt.Sprintf("にぼしが %s ぐらいをおしらせします", time.Now().Format(time.Stamp)))
            })

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", botName(), TimerStop)):
            stopTimer <- true

        case strings.HasPrefix(m.Content, fmt.Sprintf("%s %s", botName(), Genius)):
            sendMessage(s, c, "ぼくはかしこいので")
    }
}

//メッセージを送信する関数
func sendMessage(s *discordgo.Session, c *discordgo.Channel, msg string) {
    message.Send(s, c, msg)
}

func startTimer(f func()) {
    ticker := time.NewTicker(time.Second * 30)
    defer ticker.Stop()
    for {
        select {
            case <- stopTimer:
                ticker.Stop()
            case <- ticker.C:
                f()
        }
    }
}
