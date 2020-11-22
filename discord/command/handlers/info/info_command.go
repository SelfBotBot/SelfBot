package info

import (
	"fmt"
	"math"
	"runtime"
	"selfbot/discord/command/handlers"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/shirou/gopsutil/cpu"
	"github.com/shirou/gopsutil/host"

	_ "github.com/shirou/w32"
)

const (
	BlankText      = "\u200B"
	HorizontalRule = "──────────────────────────────────────────"
)
const (
	ColourGood = 0x73C62F
	ColourEh   = 0xE2D044
	ColourBad  = 0xE05545
)

const (
	RuntimeInfoTitles = "Version: \nGOARCH: "
	RuntimeInfo       = "%s\n%s"

	CPUInfoTitles = "No GoRoutines: \nNo CGo Calls: "
	CPUInfo       = "%d\n%d"

	GCInfoTitles = "GC Runs: \nGC System: "
	GCInfo       = "%d\n%s"

	HeapInfoTitles = "Heap Usage: \nHeap In Use: \nHeap Objects: "
	HeapInfo       = "%s/%s\n%s\n%d"

	SysInfoTitles = "Total Used: \nTotal GC'd: "
	SysInfo       = "%s\n%s"

	ProcInfoErrTitles = "Total Procs: "
	ProcInfoErr       = "%d"

	ProcInfoTitles = "Total Procs: \nModel Name: \nSpeed"
	ProcInfo       = "%d\n%s\n%.2fMhz"

	HostInfoTitles = "OS: \nPlatform: \nVirtualization: \nRole: \nUptime: "
	HostInfo       = "%s\n%s\n%s\n%s\n%s"
)

var _ handlers.Handler = &Handler{}

type Handler struct{}

func (p *Handler) Handle(s *discordgo.Session, m *discordgo.MessageCreate, args ...string) error {
	var subCommand = "not a letter lol (we need to make it default)"
	if len(args) > 0 {
		subCommand = strings.ToLower(args[0][0:1]) // First letter
	} else {
		subCommand = "a"
	}

	var err error
	switch subCommand {
	case "h":
		err = p.hostStatsCommand(s, m)
	case "b":
		err = p.botStatsCommand(s, m)
	case "a":
		err = p.allStatsCommand(s, m)
	default:
		_, err = s.ChannelMessageSend(m.ChannelID, "Hey, you need to /stats [all|bot|host]")
	}

	if err != nil {
		return fmt.Errorf("info command: %w", err)
	}

	return nil
}

func (p *Handler) ShouldReact() bool {
	return true
}

func (p *Handler) allStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	embedBad := p.buildEmbed(p.populateMemStats(), true)
	embedBad.Author = &discordgo.MessageEmbedAuthor{
		Name: "Bot Stats.",
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embedBad)
	if err != nil {
		return fmt.Errorf("all subcommand: %w", err)
	}
	return nil
}

func (p *Handler) botStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	embedBad := p.buildEmbed(p.populateMemStats(), false)
	embedBad.Author = &discordgo.MessageEmbedAuthor{
		Name: "Bot Runtime Stats.",
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embedBad)
	if err != nil {
		return fmt.Errorf("runtime subcommand: %w", err)
	}
	return nil
}

func (p *Handler) hostStatsCommand(s *discordgo.Session, m *discordgo.MessageCreate) error {
	embed := p.buildEmbed(nil, true)
	embed.Author = &discordgo.MessageEmbedAuthor{
		Name: "Bot Host Stats.",
	}

	_, err := s.ChannelMessageSendEmbed(m.ChannelID, embed)
	if err != nil {
		return fmt.Errorf("host subcommand: %w", err)
	}
	return nil
}

func (p *Handler) populateMemStats() *runtime.MemStats {
	stats := &runtime.MemStats{}
	runtime.ReadMemStats(stats)
	return stats
}

func (p *Handler) buildEmbed(memStats *runtime.MemStats, system bool) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField
	colour := ColourGood

	// Runtime section
	if memStats != nil {

		// Horizontal rule for separation.
		fields = append(fields, RuleField())

		// Runtime
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Runtime Information",
			Inline: true,
			Value:  RuntimeInfoTitles,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   BlankText,
			Inline: true,
			Value:  fmt.Sprintf(RuntimeInfo, runtime.Version(), runtime.GOARCH),
		})

		// Separator field
		fields = append(fields, SeperatorField())

		// CPU
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "CPU Usage Information",
			Inline: true,
			Value:  CPUInfoTitles,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   BlankText,
			Inline: true,
			Value:  fmt.Sprintf(CPUInfo, runtime.NumGoroutine(), runtime.NumCgoCall()),
		})

		// Separator field
		fields = append(fields, SeperatorField())

		// GC
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "GC Information",
			Inline: true,
			Value:  GCInfoTitles,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   BlankText,
			Inline: true,
			Value:  fmt.Sprintf(GCInfo, memStats.NumGC, ByteString(memStats.GCSys)),
		})

		// Separator field
		fields = append(fields, SeperatorField())

		// Heap
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Heap Information",
			Inline: true,
			Value:  HeapInfoTitles,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   BlankText,
			Inline: true,
			Value:  fmt.Sprintf(HeapInfo, ByteString(memStats.HeapAlloc), ByteString(memStats.HeapSys), ByteString(memStats.HeapInuse), memStats.HeapObjects),
		})

		// Separator field
		fields = append(fields, SeperatorField())

		// System
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   "Memory Information",
			Inline: true,
			Value:  SysInfoTitles,
		})
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   BlankText,
			Inline: true,
			Value:  fmt.Sprintf(SysInfo, ByteString(memStats.Sys), ByteString(memStats.TotalAlloc)),
		})

		// Separator field
		fields = append(fields, SeperatorField())

		if runtime.NumGoroutine() > 20 || memStats.HeapAlloc >= uint64(float64(memStats.HeapSys)*0.70) {
			colour = ColourEh
			if memStats.HeapAlloc >= uint64(float64(memStats.HeapSys)*0.80) {
				colour = ColourBad
			}
		}

	}

	// System section!
	if system {

		// Horizontal rule for separation.
		fields = append(fields, RuleField())

		cpuInfo, err := cpu.Info()
		if err != nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Processor Info",
				Inline: true,
				Value:  ProcInfoErrTitles,
			})
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   BlankText,
				Inline: true,
				Value:  fmt.Sprintf(ProcInfoErr, runtime.NumCPU()),
			})
		} else {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Processor Info",
				Inline: true,
				Value:  ProcInfoTitles,
			})
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   BlankText,
				Inline: true,
				Value:  fmt.Sprintf(ProcInfo, runtime.NumCPU(), cpuInfo[0].ModelName, cpuInfo[0].Mhz),
			})
		}

		// Separator field
		fields = append(fields, SeperatorField())

		hostInfo, err := host.Info()
		if err == nil {
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   "Host Info",
				Inline: true,
				Value:  HostInfoTitles,
			})
			fields = append(fields, &discordgo.MessageEmbedField{
				Name:   BlankText,
				Inline: true,
				Value:  fmt.Sprintf(HostInfo, hostInfo.OS, hostInfo.Platform, hostInfo.VirtualizationSystem, hostInfo.VirtualizationRole, time.Duration(int64(hostInfo.Uptime))*time.Second),
			})

			// Separator field
			fields = append(fields, SeperatorField())
		}

	}

	return &discordgo.MessageEmbed{
		Fields: fields,
		Type:   "rich",
		Color:  colour,
	}

}

func RuleField() *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Inline: false,
		Value:  BlankText,
		Name:   HorizontalRule,
	}
}

func SeperatorField() *discordgo.MessageEmbedField {
	return &discordgo.MessageEmbedField{
		Inline: true,
		Value:  BlankText,
		Name:   BlankText,
	}
}

func ByteString(bytes uint64) string {
	if bytes < 1024 {
		return fmt.Sprintf("%d B", bytes)
	}

	floatBytes := float64(bytes)
	exp := uint64(math.Log(floatBytes) / math.Log(1024))
	pre := "KMGTPE"[exp-1]
	return fmt.Sprintf("%.2f %ciB", floatBytes/math.Pow(1024, float64(exp)), pre)

}
