package message

import (
    "github.com/bwmarrin/discordgo"
    "log"
)

func Send(s *discordgo.Session, c *discordgo.Channel, msg string) {
    _, err := s.ChannelMessageSend(c.ID, msg)

    log.Println(">>> " + msg)
    if err != nil {
        log.Println("Error sending message: ", err)
    }
}

