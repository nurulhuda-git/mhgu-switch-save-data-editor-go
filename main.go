// package main

// import (
// 	"encoding/binary"
// 	"flag"
// 	"fmt"
// 	"io/ioutil"
// 	"os"
// )

// // Player structure based on C# code
// type Player struct {
// 	SaveOffset int

// 	// Basic info
// 	Name          string
// 	PlayTime      int
// 	Funds         int
// 	HunterRank    int
// 	HRPoints      int
// 	AcademyPoints int

// 	// Village points
// 	BhernaPoints int
// 	KokotoPoints int
// 	PokkePoints  int
// 	YukumoPoints int

// 	// Appearance
// 	Voice     byte
// 	EyeColor  byte
// 	Clothing  byte
// 	Gender    byte
// 	HairStyle byte
// 	Face      byte
// 	Features  byte

// 	// Colors
// 	SkinColorRGBA     [4]byte
// 	HairColorRGBA     [4]byte
// 	FeaturesColorRGBA [4]byte
// 	ClothingColorRGBA [4]byte

// 	// Item box
// 	ItemId    []string
// 	ItemCount []string

// 	// Equipment
// 	EquipmentInfo   []byte
// 	EquipmentPalico []byte

// 	// Other data
// 	PalicoData         []byte
// 	GuildCardData      []byte
// 	ArenaData          []byte
// 	MonsterKills       []byte
// 	MonsterCaptures    []byte
// 	MonsterSizes       []byte
// 	ManualShoutouts    []byte
// 	AutomaticShoutouts []byte
// }

// func main() {
// 	inputFile := flag.String("input", "", "Input save file path")
// 	// outputFile := flag.String("output", "system_modified.bin", "Output save file path")
// 	displayAll := flag.Bool("all", false, "Display all information")
// 	displayItems := flag.Bool("items", false, "Display item box summary")
// 	displayEquips := flag.Bool("equips", false, "Display equipment box summary")
// 	displayPalico := flag.Bool("palico", false, "Display palico summary")
// 	slot := flag.Int("slot", 1, "Character slot (1-3)")
// 	debug := flag.Bool("debug", false, "Debug mode")

// 	flag.Parse()

// 	if *inputFile == "" {
// 		fmt.Println("Error: --input flag is required")
// 		os.Exit(1)
// 	}

// 	fmt.Printf("=== MHGU SAVE EDITOR (Go) ===\n")
// 	fmt.Printf("Based on MHXX Save Editor v0.09c by Ukee\n")

// 	// Load save file
// 	saveData, extractedData, isSwitch, err := loadSaveFile(*inputFile)
// 	if err != nil {
// 		fmt.Printf("Error loading save: %v\n", err)
// 		os.Exit(1)
// 	}

// 	fmt.Printf("\nFile: %s\n", *inputFile)
// 	fmt.Printf("Size: %d bytes\n", len(saveData))
// 	if isSwitch {
// 		fmt.Printf("Type: MHGU Switch (removed 36-byte header)\n")
// 	}
// 	fmt.Printf("Extracted size: %d bytes\n", len(extractedData))

// 	// Extract player data
// 	player, err := extractPlayerData(extractedData, *slot)
// 	if err != nil {
// 		fmt.Printf("Error extracting player data: %v\n", err)
// 		os.Exit(1)
// 	}

// 	// Display information based on flags
// 	if *displayAll || !(*displayItems || *displayEquips || *displayPalico) {
// 		// Default: show character info
// 		displayCharacterInfo(player, *slot, *debug)
// 	}

// 	if *displayAll || *displayItems {
// 		displayItemBoxInfo(player, *debug)
// 	}

// 	if *displayAll || *displayEquips {
// 		displayEquipmentInfo(player, *debug)
// 	}

// 	if *displayAll || *displayPalico {
// 		displayPalicoInfo(player, *debug)
// 	}

// 	// Show data section sizes
// 	if *displayAll {
// 		displayDataSections(player)
// 	}
// }

// func loadSaveFile(filename string) ([]byte, []byte, bool, error) {
// 	data, err := ioutil.ReadFile(filename)
// 	if err != nil {
// 		return nil, nil, false, err
// 	}

// 	var extractedData []byte
// 	isSwitch := false

// 	// Check save type
// 	switch len(data) {
// 	case 4726152: // 3DS
// 		extractedData = data
// 	case 4726152 + 36: // Switch
// 		fallthrough
// 	case 4726152 + 432948: // MHGU
// 		if len(data) >= 36 {
// 			extractedData = data[36:]
// 			isSwitch = true
// 		} else {
// 			return nil, nil, false, fmt.Errorf("Switch save too small")
// 		}
// 	default:
// 		return nil, nil, false, fmt.Errorf("unknown save size: %d bytes", len(data))
// 	}

// 	return data, extractedData, isSwitch, nil
// }

// func extractPlayerData(data []byte, slot int) (*Player, error) {
// 	if slot < 1 || slot > 3 {
// 		return nil, fmt.Errorf("invalid slot: %d", slot)
// 	}

// 	// Check slot usage
// 	if len(data) < 0x07 {
// 		return nil, fmt.Errorf("save data too small")
// 	}

// 	var slotUsed bool
// 	var slotOffset int

// 	switch slot {
// 	case 1:
// 		slotUsed = data[0x04] == 1
// 		if len(data) >= 0x14 {
// 			slotOffset = int(binary.LittleEndian.Uint32(data[0x10:]))
// 		}
// 	case 2:
// 		slotUsed = data[0x05] == 1
// 		if len(data) >= 0x18 {
// 			slotOffset = int(binary.LittleEndian.Uint32(data[0x14:]))
// 		}
// 	case 3:
// 		slotUsed = data[0x06] == 1
// 		if len(data) >= 0x1C {
// 			slotOffset = int(binary.LittleEndian.Uint32(data[0x18:]))
// 		}
// 	}

// 	if !slotUsed {
// 		return nil, fmt.Errorf("slot %d is not used", slot)
// 	}

// 	if slotOffset == 0 || slotOffset >= len(data) {
// 		return nil, fmt.Errorf("invalid character offset: 0x%X", slotOffset)
// 	}

// 	player := &Player{
// 		SaveOffset: slotOffset,
// 	}

// 	// Extract basic info
// 	extractBasicInfo(player, data)

// 	// Extract item box
// 	extractItemBox(player, data)

// 	// Extract other data sections
// 	extractOtherData(player, data)

// 	return player, nil
// }

// func extractBasicInfo(player *Player, data []byte) {
// 	offset := player.SaveOffset

// 	// Name
// 	nameOffset := offset + 0x23B7D
// 	if nameOffset+32 <= len(data) {
// 		player.Name = extractNullTerminatedString(data, nameOffset, 32)
// 	}

// 	// Play time
// 	if offset+0x24 <= len(data) {
// 		player.PlayTime = int(binary.LittleEndian.Uint32(data[offset+0x20:]))
// 	}

// 	// Funds
// 	if offset+0x28 <= len(data) {
// 		player.Funds = int(binary.LittleEndian.Uint32(data[offset+0x24:]))
// 	}

// 	// Hunter rank
// 	if offset+0x2A <= len(data) {
// 		player.HunterRank = int(binary.LittleEndian.Uint16(data[offset+0x28:]))
// 	}

// 	// HR points
// 	if offset+0x280F <= len(data) {
// 		player.HRPoints = int(binary.LittleEndian.Uint32(data[offset+0x280B:]))
// 	}

// 	// Academy points
// 	if offset+0x281B <= len(data) {
// 		player.AcademyPoints = int(binary.LittleEndian.Uint32(data[offset+0x2817:]))
// 	}

// 	// Village points
// 	if offset+0x282B <= len(data) {
// 		player.BhernaPoints = int(binary.LittleEndian.Uint32(data[offset+0x281B:]))
// 		player.KokotoPoints = int(binary.LittleEndian.Uint32(data[offset+0x281F:]))
// 		player.PokkePoints = int(binary.LittleEndian.Uint32(data[offset+0x2823:]))
// 		player.YukumoPoints = int(binary.LittleEndian.Uint32(data[offset+0x2827:]))
// 	}

// 	// Appearance
// 	if offset+0x23B50 <= len(data) {
// 		player.Voice = data[offset+0x23B48]
// 		player.EyeColor = data[offset+0x23B49]
// 		player.Clothing = data[offset+0x23B4A]
// 		player.Gender = data[offset+0x23B4B]
// 		player.HairStyle = data[offset+0x23B4D]
// 		player.Face = data[offset+0x23B4E]
// 		player.Features = data[offset+0x23B4F]
// 	}

// 	// Colors
// 	if offset+0x23B77 <= len(data) {
// 		copy(player.SkinColorRGBA[:], data[offset+0x23B67:])
// 		copy(player.HairColorRGBA[:], data[offset+0x23B6B:])
// 		copy(player.FeaturesColorRGBA[:], data[offset+0x23B6F:])
// 		copy(player.ClothingColorRGBA[:], data[offset+0x23B73:])
// 	}
// }

// func extractItemBox(player *Player, data []byte) {
// 	offset := player.SaveOffset
// 	itemBoxOffset := offset + 0x0278

// 	if itemBoxOffset+5463 > len(data) {
// 		return
// 	}

// 	// Extract item box data (simplified - actual extraction needs bit manipulation)
// 	player.ItemId = make([]string, 2300)
// 	player.ItemCount = make([]string, 2300)

// 	// For now, just store the raw bytes
// 	itemBoxData := make([]byte, 5463)
// 	copy(itemBoxData, data[itemBoxOffset:itemBoxOffset+5463])

// 	// Simplified: mark all as empty
// 	for i := 0; i < 2300; i++ {
// 		player.ItemId[i] = "0"
// 		player.ItemCount[i] = "0"
// 	}
// }

// func extractOtherData(player *Player, data []byte) {
// 	offset := player.SaveOffset

// 	// Equipment box
// 	equipOffset := offset + 0x62EE
// 	if equipOffset+72000 <= len(data) {
// 		player.EquipmentInfo = make([]byte, 72000)
// 		copy(player.EquipmentInfo, data[equipOffset:equipOffset+72000])
// 	}

// 	// Palico equipment
// 	palicoEquipOffset := offset + 0x17C2E
// 	if palicoEquipOffset+36000 <= len(data) {
// 		player.EquipmentPalico = make([]byte, 36000)
// 		copy(player.EquipmentPalico, data[palicoEquipOffset:palicoEquipOffset+36000])
// 	}

// 	// Palico data
// 	palicoOffset := offset + 0x23BB6
// 	if palicoOffset+27216 <= len(data) {
// 		player.PalicoData = make([]byte, 27216)
// 		copy(player.PalicoData, data[palicoOffset:palicoOffset+27216])
// 	}

// 	// Guild card
// 	guildCardOffset := offset + 0xC71BD
// 	if guildCardOffset+4986 <= len(data) {
// 		player.GuildCardData = make([]byte, 4986)
// 		copy(player.GuildCardData, data[guildCardOffset:guildCardOffset+4986])
// 	}

// 	// Arena data
// 	arenaOffset := offset + 0xC83E1
// 	if arenaOffset+342 <= len(data) {
// 		player.ArenaData = make([]byte, 342)
// 		copy(player.ArenaData, data[arenaOffset:arenaOffset+342])
// 	}

// 	// Monster data
// 	monsterKillsOffset := offset + 0x5EA6
// 	if monsterKillsOffset+274 <= len(data) {
// 		player.MonsterKills = make([]byte, 274)
// 		copy(player.MonsterKills, data[monsterKillsOffset:monsterKillsOffset+274])
// 	}

// 	monsterCaptureOffset := offset + 0x5FB8
// 	if monsterCaptureOffset+274 <= len(data) {
// 		player.MonsterCaptures = make([]byte, 274)
// 		copy(player.MonsterCaptures, data[monsterCaptureOffset:monsterCaptureOffset+274])
// 	}

// 	monsterSizeOffset := offset + 0x60CA
// 	if monsterSizeOffset+548 <= len(data) {
// 		player.MonsterSizes = make([]byte, 548)
// 		copy(player.MonsterSizes, data[monsterSizeOffset:monsterSizeOffset+548])
// 	}

// 	// Shoutouts
// 	manualShoutOffset := offset + 0x11D629
// 	if manualShoutOffset+2880 <= len(data) {
// 		player.ManualShoutouts = make([]byte, 2880)
// 		copy(player.ManualShoutouts, data[manualShoutOffset:manualShoutOffset+2880])
// 	}

// 	autoShoutOffset := offset + 0x11E169
// 	if autoShoutOffset+1620 <= len(data) {
// 		player.AutomaticShoutouts = make([]byte, 1620)
// 		copy(player.AutomaticShoutouts, data[autoShoutOffset:autoShoutOffset+1620])
// 	}
// }

// func displayCharacterInfo(player *Player, slot int, debug bool) {
// 	fmt.Printf("\n=== CHARACTER SLOT %d ===\n", slot)
// 	fmt.Printf("Save Offset: 0x%08X\n", player.SaveOffset)

// 	// Basic Info
// 	fmt.Printf("\n--- BASIC INFORMATION ---\n")
// 	fmt.Printf("Name:          %s\n", player.Name)
// 	fmt.Printf("Play Time:     %s\n", formatPlayTime(player.PlayTime))
// 	fmt.Printf("Funds:         %dz\n", player.Funds)
// 	fmt.Printf("Hunter Rank:   %d\n", player.HunterRank)
// 	fmt.Printf("HR Points:     %d\n", player.HRPoints)
// 	fmt.Printf("Academy Points:%d\n", player.AcademyPoints)

// 	// Village Points
// 	fmt.Printf("\n--- VILLAGE POINTS ---\n")
// 	fmt.Printf("Bherna:   %d\n", player.BhernaPoints)
// 	fmt.Printf("Kokoto:   %d\n", player.KokotoPoints)
// 	fmt.Printf("Pokke:    %d\n", player.PokkePoints)
// 	fmt.Printf("Yukumo:   %d\n", player.YukumoPoints)

// 	// Appearance
// 	fmt.Printf("\n--- APPEARANCE ---\n")
// 	fmt.Printf("Gender:        %d\n", player.Gender)
// 	fmt.Printf("Voice:         %d\n", player.Voice)
// 	fmt.Printf("Eye Color:     %d\n", player.EyeColor)
// 	fmt.Printf("Clothing:      %d\n", player.Clothing)
// 	fmt.Printf("Hair Style:    %d\n", player.HairStyle)
// 	fmt.Printf("Face:          %d\n", player.Face)
// 	fmt.Printf("Features:      %d\n", player.Features)

// 	// Colors
// 	fmt.Printf("\n--- COLORS (RGBA) ---\n")
// 	fmt.Printf("Skin:      R:%3d G:%3d B:%3d A:%3d\n",
// 		player.SkinColorRGBA[0], player.SkinColorRGBA[1],
// 		player.SkinColorRGBA[2], player.SkinColorRGBA[3])
// 	fmt.Printf("Hair:      R:%3d G:%3d B:%3d A:%3d\n",
// 		player.HairColorRGBA[0], player.HairColorRGBA[1],
// 		player.HairColorRGBA[2], player.HairColorRGBA[3])
// 	fmt.Printf("Features:  R:%3d G:%3d B:%3d A:%3d\n",
// 		player.FeaturesColorRGBA[0], player.FeaturesColorRGBA[1],
// 		player.FeaturesColorRGBA[2], player.FeaturesColorRGBA[3])
// 	fmt.Printf("Clothing:  R:%3d G:%3d B:%3d A:%3d\n",
// 		player.ClothingColorRGBA[0], player.ClothingColorRGBA[1],
// 		player.ClothingColorRGBA[2], player.ClothingColorRGBA[3])

// 	if debug {
// 		fmt.Printf("\n--- DEBUG INFO ---\n")
// 		fmt.Printf("Name offset: 0x%08X + 0x23B7D = 0x%08X\n",
// 			player.SaveOffset, player.SaveOffset+0x23B7D)
// 	}
// }

// func displayItemBoxInfo(player *Player, debug bool) {
// 	fmt.Printf("\n=== ITEM BOX ===\n")
// 	fmt.Printf("Total slots: 2300\n")
// 	fmt.Printf("Data size: %d bytes\n", len(player.ItemId)*19) // 19 bits per item

// 	// Count non-empty items
// 	nonEmpty := 0
// 	for i := 0; i < 2300; i++ {
// 		if player.ItemId[i] != "0" && player.ItemCount[i] != "0" {
// 			nonEmpty++
// 		}
// 	}

// 	fmt.Printf("Non-empty items: %d\n", nonEmpty)
// 	fmt.Printf("Empty slots: %d\n", 2300-nonEmpty)

// 	if debug && nonEmpty > 0 {
// 		fmt.Printf("\nFirst 10 non-empty items:\n")
// 		count := 0
// 		for i := 0; i < 2300 && count < 10; i++ {
// 			if player.ItemId[i] != "0" && player.ItemCount[i] != "0" {
// 				fmt.Printf("  Slot %4d: ID=%s, Count=%s\n",
// 					i+1, player.ItemId[i], player.ItemCount[i])
// 				count++
// 			}
// 		}
// 	}
// }

// func displayEquipmentInfo(player *Player, debug bool) {
// 	fmt.Printf("\n=== EQUIPMENT BOX ===\n")
// 	fmt.Printf("Total slots: 2000\n")
// 	fmt.Printf("Data size: %d bytes (36 bytes per equipment)\n", len(player.EquipmentInfo))

// 	// Count non-empty equipment
// 	nonEmpty := 0
// 	for i := 0; i < 2000; i++ {
// 		// Check if equipment type is not 0 (empty)
// 		if len(player.EquipmentInfo) > i*36 && player.EquipmentInfo[i*36] != 0 {
// 			nonEmpty++
// 		}
// 	}

// 	fmt.Printf("Non-empty equipment: %d\n", nonEmpty)
// 	fmt.Printf("Empty slots: %d\n", 2000-nonEmpty)

// 	if debug && nonEmpty > 0 {
// 		fmt.Printf("\nFirst 5 equipment items:\n")
// 		count := 0
// 		for i := 0; i < 2000 && count < 5; i++ {
// 			if len(player.EquipmentInfo) > i*36+36 {
// 				eqType := player.EquipmentInfo[i*36]
// 				eqID := binary.LittleEndian.Uint16(player.EquipmentInfo[i*36+2:])
// 				if eqType != 0 {
// 					fmt.Printf("  Slot %4d: Type=%d, ID=%d\n", i+1, eqType, eqID)
// 					count++
// 				}
// 			}
// 		}
// 	}

// 	// Palico equipment
// 	fmt.Printf("\n=== PALICO EQUIPMENT ===\n")
// 	fmt.Printf("Total slots: 1000\n")
// 	fmt.Printf("Data size: %d bytes\n", len(player.EquipmentPalico))

// 	palicoNonEmpty := 0
// 	for i := 0; i < 1000; i++ {
// 		if len(player.EquipmentPalico) > i*36 && player.EquipmentPalico[i*36] != 0 {
// 			palicoNonEmpty++
// 		}
// 	}

// 	fmt.Printf("Non-empty palico equipment: %d\n", palicoNonEmpty)
// 	fmt.Printf("Empty slots: %d\n", 1000-palicoNonEmpty)
// }

// func displayPalicoInfo(player *Player, debug bool) {
// 	fmt.Printf("\n=== PALICO DATA ===\n")
// 	fmt.Printf("Total slots: 84\n")
// 	fmt.Printf("Data size: %d bytes (324 bytes per palico)\n", len(player.PalicoData))

// 	// Count palicos with names
// 	palicoCount := 0
// 	for i := 0; i < 84; i++ {
// 		if len(player.PalicoData) > i*324 {
// 			// Check if name is not empty (first byte not 0)
// 			if player.PalicoData[i*324] != 0 {
// 				palicoCount++
// 			}
// 		}
// 	}

// 	fmt.Printf("Palicos with names: %d\n", palicoCount)

// 	if debug && palicoCount > 0 {
// 		fmt.Printf("\nFirst 3 palicos:\n")
// 		count := 0
// 		for i := 0; i < 84 && count < 3; i++ {
// 			if len(player.PalicoData) > i*324+32 && player.PalicoData[i*324] != 0 {
// 				name := extractNullTerminatedString(player.PalicoData, i*324, 32)
// 				if name != "" {
// 					// Get palico type (offset 37)
// 					palicoType := "Unknown"
// 					if len(player.PalicoData) > i*324+37 {
// 						pt := player.PalicoData[i*324+37]
// 						palicoType = fmt.Sprintf("Type %d", pt)
// 					}
// 					fmt.Printf("  Slot %2d: %s (%s)\n", i+1, name, palicoType)
// 					count++
// 				}
// 			}
// 		}
// 	}
// }

// func displayDataSections(player *Player) {
// 	fmt.Printf("\n=== DATA SECTIONS ===\n")
// 	fmt.Printf("Item Box:          %7d bytes\n", 5463)
// 	fmt.Printf("Equipment Box:     %7d bytes\n", 72000)
// 	fmt.Printf("Palico Equipment:  %7d bytes\n", 36000)
// 	fmt.Printf("Palico Data:       %7d bytes\n", 27216)
// 	fmt.Printf("Monster Kills:     %7d bytes\n", 274)
// 	fmt.Printf("Monster Captures:  %7d bytes\n", 274)
// 	fmt.Printf("Monster Sizes:     %7d bytes\n", 548)
// 	fmt.Printf("Guild Card:        %7d bytes\n", 4986)
// 	fmt.Printf("Arena Data:        %7d bytes\n", 342)
// 	fmt.Printf("Manual Shoutouts:  %7d bytes\n", 2880)
// 	fmt.Printf("Auto Shoutouts:    %7d bytes\n", 1620)

// 	// Show actual loaded sizes
// 	fmt.Printf("\n--- ACTUALLY LOADED ---\n")
// 	fmt.Printf("Equipment Info:    %7d bytes\n", len(player.EquipmentInfo))
// 	fmt.Printf("Equipment Palico:  %7d bytes\n", len(player.EquipmentPalico))
// 	fmt.Printf("Palico Data:       %7d bytes\n", len(player.PalicoData))
// 	fmt.Printf("Guild Card:        %7d bytes\n", len(player.GuildCardData))
// 	fmt.Printf("Arena Data:        %7d bytes\n", len(player.ArenaData))
// 	fmt.Printf("Monster Data:      %7d bytes (kills: %d, captures: %d, sizes: %d)\n",
// 		len(player.MonsterKills)+len(player.MonsterCaptures)+len(player.MonsterSizes),
// 		len(player.MonsterKills), len(player.MonsterCaptures), len(player.MonsterSizes))
// 	fmt.Printf("Shoutouts:         %7d bytes (manual: %d, auto: %d)\n",
// 		len(player.ManualShoutouts)+len(player.AutomaticShoutouts),
// 		len(player.ManualShoutouts), len(player.AutomaticShoutouts))
// }

// // Helper functions
// func extractNullTerminatedString(data []byte, offset, maxLen int) string {
// 	if offset < 0 || offset >= len(data) {
// 		return ""
// 	}

// 	end := offset
// 	for end < len(data) && end-offset < maxLen && data[end] != 0 {
// 		end++
// 	}

// 	if end == offset {
// 		return ""
// 	}

// 	return string(data[offset:end])
// }

// func formatPlayTime(seconds int) string {
// 	hours := seconds / 3600
// 	minutes := (seconds % 3600) / 60
// 	secs := seconds % 60
// 	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
// }

// // Search function (from previous code)
// func searchInSave(data []byte, searchStr string) []int {
// 	searchBytes := []byte(searchStr)
// 	var results []int

// 	for i := 0; i <= len(data)-len(searchBytes); i++ {
// 		found := true
// 		for j := 0; j < len(searchBytes); j++ {
// 			if data[i+j] != searchBytes[j] {
// 				found = false
// 				break
// 			}
// 		}
// 		if found {
// 			results = append(results, i)
// 		}
// 	}

// 	return results
// }

package main

import (
	"encoding/binary"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

// Player structure based on C# code
type Player struct {
	SaveOffset int

	// Basic info
	Name          string
	PlayTime      int
	Funds         int
	HunterRank    int
	HRPoints      int
	AcademyPoints int

	// Village points
	BhernaPoints int
	KokotoPoints int
	PokkePoints  int
	YukumoPoints int

	// Appearance
	Voice     byte
	EyeColor  byte
	Clothing  byte
	Gender    byte
	HairStyle byte
	Face      byte
	Features  byte

	// Colors
	SkinColorRGBA     [4]byte
	HairColorRGBA     [4]byte
	FeaturesColorRGBA [4]byte
	ClothingColorRGBA [4]byte

	// Item box
	ItemId    []string
	ItemCount []string

	// Equipment
	EquipmentInfo   []byte
	EquipmentPalico []byte

	// Other data
	PalicoData         []byte
	GuildCardData      []byte
	ArenaData          []byte
	MonsterKills       []byte
	MonsterCaptures    []byte
	MonsterSizes       []byte
	ManualShoutouts    []byte
	AutomaticShoutouts []byte

	// Save data for writing back
	SaveData []byte
}

// Customization options
type Customization struct {
	Name      string
	Gender    int
	Voice     int
	EyeColor  int
	Clothing  int
	HairStyle int
	Face      int
	Features  int

	// Colors
	SkinColor     [4]int
	HairColor     [4]int
	FeaturesColor [4]int
	ClothingColor [4]int

	// Other stats
	Funds         int
	HRPoints      int
	HunterRank    int
	AcademyPoints int
	BhernaPoints  int
	KokotoPoints  int
	PokkePoints   int
	YukumoPoints  int
}

func genTime(base string) string {
	timestamp := time.Now().Unix() // Unix() returns seconds
	return base + "_" + strconv.FormatInt(timestamp, 10)
}

func main() {
	inputFile := flag.String("input", "", "Input save file path")
	outputFile := flag.String("output", genTime("system_modified"), "Output save file path")
	searchStr := flag.String("search", "", "Search for a string in the save file")
	replaceStr := flag.String("replace", "", "Replace found strings with this string")
	replaceAll := flag.Bool("replaceall", false, "Replace all occurrences (default: only at name locations)")
	displayAll := flag.Bool("all", false, "Display all information")
	displayItems := flag.Bool("items", false, "Display item box summary")
	displayEquips := flag.Bool("equips", false, "Display equipment box summary")
	displayPalico := flag.Bool("palico", false, "Display palico summary")
	slot := flag.Int("slot", 1, "Character slot (1-3)")
	debug := flag.Bool("debug", false, "Debug mode")

	// Customization flags
	// setName := flag.String("name", "", "Set character name")
	setGender := flag.Int("gender", -1, "Set gender (0=male, 1=female)")
	setVoice := flag.Int("voice", -1, "Set voice (0-?)")
	setEyeColor := flag.Int("eyecolor", -1, "Set eye color (0-?)")
	setClothing := flag.Int("clothing", -1, "Set clothing (0-?)")
	setHairStyle := flag.Int("hairstyle", -1, "Set hair style (0-?)")
	setFace := flag.Int("face", -1, "Set face (0-?)")
	setFeatures := flag.Int("features", -1, "Set features (0-?)")

	// Color flags
	setSkinR := flag.Int("skinr", -1, "Set skin color red (0-255)")
	setSkinG := flag.Int("sking", -1, "Set skin color green (0-255)")
	setSkinB := flag.Int("skinb", -1, "Set skin color blue (0-255)")
	setSkinA := flag.Int("skina", -1, "Set skin color alpha (0-255)")

	setHairR := flag.Int("hairr", -1, "Set hair color red (0-255)")
	setHairG := flag.Int("hairg", -1, "Set hair color green (0-255)")
	setHairB := flag.Int("hairb", -1, "Set hair color blue (0-255)")
	setHairA := flag.Int("haira", -1, "Set hair color alpha (0-255)")

	// Stats flags
	setFunds := flag.Int("funds", -1, "Set funds")
	setHRPoints := flag.Int("hrpoints", -1, "Set HR points")
	setHunterRank := flag.Int("hr", -1, "Set hunter rank")
	setAcademyPoints := flag.Int("academy", -1, "Set academy points")

	// Village points flags
	setBherna := flag.Int("bherna", -1, "Set Bherna points")
	setKokoto := flag.Int("kokoto", -1, "Set Kokoto points")
	setPokke := flag.Int("pokke", -1, "Set Pokke points")
	setYukumo := flag.Int("yukumo", -1, "Set Yukumo points")

	flag.Parse()

	if *inputFile == "" {
		fmt.Println("Error: --input flag is required")
		fmt.Println("\nUsage examples:")
		fmt.Println("  # Display character info:")
		fmt.Println("  ./mhgu-editor --input system.bin --slot 1")
		fmt.Println("\n  # Change to female character:")
		fmt.Println("  ./mhgu-editor --input system.bin --slot 1 --gender 1")
		fmt.Println("\n  # Customize appearance:")
		fmt.Println("  ./mhgu-editor --input system.bin --slot 1 --gender 1 --voice 15 --hairstyle 20")
		os.Exit(1)
	}

	fmt.Printf("=== MHGU SAVE EDITOR ===\n")
	fmt.Printf("Based on MHXX Save Editor v0.09c by Ukee\n")

	if *searchStr != "" && *replaceStr == "" {
		// Search mode only
		// Load and search
		_, extractedData, isSwitch, err := loadSaveFile(*inputFile)
		if err != nil {
			fmt.Printf("Error loading save: %v\n", err)
			os.Exit(1)
		}

		results := searchInMemory(extractedData, []byte(*searchStr))

		fmt.Printf("\n=== SEARCH RESULTS ===\n")
		fmt.Printf("Searching for: \"%s\"\n", *searchStr)
		fmt.Printf("Found %d occurrence(s)\n", len(results))

		for i, offset := range results {
			fmt.Printf("\n[%d] Offset: 0x%08X (in extracted data)\n", i+1, offset)

			// Show context
			context := extractStringContext(extractedData, offset, len(*searchStr), 32)
			fmt.Printf("Context: \"%s\"\n", context)

			// Show hex dump
			if i < 3 { // Limit output
				fmt.Printf("Hex (16 bytes before and after):\n")
				hexDump(extractedData, offset-16, 48, true)
			}
		}

		// Also show absolute offsets
		fmt.Printf("\n=== ABSOLUTE FILE OFFSETS ===\n")
		for i, offset := range results {
			absoluteOffset := offset
			if isSwitch {
				absoluteOffset += 36
			}
			fmt.Printf("[%d] File offset: 0x%08X\n", i+1, absoluteOffset)
		}

		return
	}

	// Load save file
	saveData, extractedData, isSwitch, err := loadSaveFile(*inputFile)
	if err != nil {
		fmt.Printf("Error loading save: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("\nFile: %s\n", *inputFile)
	fmt.Printf("Size: %d bytes\n", len(saveData))
	if isSwitch {
		fmt.Printf("Type: MHGU Switch\n")
	} else {
		fmt.Printf("Type: 3DS\n")
	}
	fmt.Printf("Extracted size: %d bytes\n", len(extractedData))

	// Check slot usage
	slotUsed := false
	switch *slot {
	case 1:
		slotUsed = extractedData[0x04] == 1
	case 2:
		slotUsed = extractedData[0x05] == 1
	case 3:
		slotUsed = extractedData[0x06] == 1
	}

	if !slotUsed {
		fmt.Printf("\nError: Slot %d is not used\n", *slot)
		os.Exit(1)
	}

	// Extract player data
	player, err := extractPlayerData(extractedData, *slot)
	if err != nil {
		fmt.Printf("Error extracting player data: %v\n", err)
		os.Exit(1)
	}

	// Store the save data for modification
	player.SaveData = make([]byte, len(extractedData))
	copy(player.SaveData, extractedData)

	// Check if we need to customize
	needCustomize := *setGender != -1 || *setVoice != -1 ||
		*setEyeColor != -1 || *setClothing != -1 || *setHairStyle != -1 ||
		*setFace != -1 || *setFeatures != -1 ||
		*setSkinR != -1 || *setSkinG != -1 || *setSkinB != -1 || *setSkinA != -1 ||
		*setHairR != -1 || *setHairG != -1 || *setHairB != -1 || *setHairA != -1 ||
		*setFunds != -1 || *setHRPoints != -1 || *setHunterRank != -1 ||
		*setAcademyPoints != -1 || *setBherna != -1 || *setKokoto != -1 ||
		*setPokke != -1 || *setYukumo != -1

	// Add this after checking for customization but before displaying info
	if *searchStr != "" && *replaceStr != "" {
		// Search and replace mode
		fmt.Printf("\n=== SEARCH AND REPLACE ===\n")
		fmt.Printf("Search: \"%s\"\n", *searchStr)
		fmt.Printf("Replace: \"%s\"\n", *replaceStr)

		searchBytes := []byte(*searchStr)
		replaceBytes := []byte(*replaceStr)

		if len(searchBytes) != len(replaceBytes) {
			fmt.Printf("Warning: Length mismatch (search: %d, replace: %d)\n",
				len(searchBytes), len(replaceBytes))
			fmt.Println("This may cause issues if strings have different null termination.")
		}

		// Find all occurrences
		allResults := searchInMemory(extractedData, searchBytes)
		if len(allResults) == 0 {
			fmt.Printf("String \"%s\" not found\n", *searchStr)
			os.Exit(1)
		}

		fmt.Printf("Found %d total occurrence(s)\n", len(allResults))

		// Filter results if not replacing all
		var replaceResults []int
		if *replaceAll {
			replaceResults = allResults
			fmt.Printf("Will replace ALL %d occurrence(s)\n", len(replaceResults))
		} else {
			// Only replace at valid name locations (null-terminated, reasonable context)
			for _, offset := range allResults {
				if isValidNameLocation(extractedData, offset, len(searchBytes)) {
					replaceResults = append(replaceResults, offset)
				}
			}
			fmt.Printf("Will replace %d occurrence(s) at valid name locations\n", len(replaceResults))
		}

		if len(replaceResults) == 0 {
			fmt.Println("No valid occurrences to replace")
			os.Exit(1)
		}

		// Show what will be replaced
		fmt.Printf("\nOccurrences to replace:\n")
		for i, offset := range replaceResults {
			if i < 5 { // Limit display
				context := extractStringContext(extractedData, offset, len(searchBytes), 16)
				fmt.Printf("[%d] Offset: 0x%08X - \"%s\"\n", i+1, offset, context)
			}
		}
		if len(replaceResults) > 5 {
			fmt.Printf("... and %d more\n", len(replaceResults)-5)
		}

		// Ask for confirmation
		fmt.Printf("\nReplace %d occurrence(s)? (y/n): ", len(replaceResults))
		var response string
		fmt.Scanln(&response)

		if response != "y" && response != "Y" {
			fmt.Println("Replacement cancelled")
			return
		}

		// Perform replacement
		modifiedData := replaceInMemory(extractedData, searchBytes, replaceBytes, replaceResults)

		// Save modified file
		finalData := make([]byte, len(saveData))
		if isSwitch {
			copy(finalData, saveData[:36]) // Keep original header
			copy(finalData[36:], modifiedData)
		} else {
			copy(finalData, modifiedData)
		}

		err = ioutil.WriteFile(*outputFile, finalData, 0644)
		if err != nil {
			fmt.Printf("Error saving file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nSuccessfully replaced %d occurrence(s)\n", len(replaceResults))
		fmt.Printf("Saved to: %s\n", *outputFile)

		// Verify
		verifyResults := searchInMemory(modifiedData, searchBytes)
		verifyReplace := searchInMemory(modifiedData, replaceBytes)
		fmt.Printf("Verification: Original found %d times, Replacement found %d times\n",
			len(verifyResults), len(verifyReplace))

		return
	}
	// End of Add this after checking for customization but before displaying info

	if needCustomize {
		// Create customization struct with player's current values
		custom := Customization{
			Name:      player.Name,
			Gender:    int(player.Gender),
			Voice:     int(player.Voice),
			EyeColor:  int(player.EyeColor),
			Clothing:  int(player.Clothing),
			HairStyle: int(player.HairStyle),
			Face:      int(player.Face),
			Features:  int(player.Features),
			// Convert [4]byte to [4]int
			SkinColor:     byteArrayToIntArray(player.SkinColorRGBA),
			HairColor:     byteArrayToIntArray(player.HairColorRGBA),
			FeaturesColor: byteArrayToIntArray(player.FeaturesColorRGBA),
			ClothingColor: byteArrayToIntArray(player.ClothingColorRGBA),
			Funds:         player.Funds,
			HRPoints:      player.HRPoints,
			HunterRank:    player.HunterRank,
			AcademyPoints: player.AcademyPoints,
			BhernaPoints:  player.BhernaPoints,
			KokotoPoints:  player.KokotoPoints,
			PokkePoints:   player.PokkePoints,
			YukumoPoints:  player.YukumoPoints,
		}

		// Override only if flags are explicitly set (not -1)
		fmt.Println("AAAA", *setGender)
		if *setGender != -1 {
			custom.Gender = *setGender
		}
		if *setVoice != -1 {
			custom.Voice = *setVoice
		}
		if *setEyeColor != -1 {
			custom.EyeColor = *setEyeColor
		}
		if *setClothing != -1 {
			custom.Clothing = *setClothing
		}
		if *setHairStyle != -1 {
			custom.HairStyle = *setHairStyle
		}
		if *setFace != -1 {
			custom.Face = *setFace
		}
		if *setFeatures != -1 {
			custom.Features = *setFeatures
		}

		// Only override color components that are explicitly set
		if *setSkinR != -1 {
			custom.SkinColor[0] = clampColor(*setSkinR)
		}
		if *setSkinG != -1 {
			custom.SkinColor[1] = clampColor(*setSkinG)
		}
		if *setSkinB != -1 {
			custom.SkinColor[2] = clampColor(*setSkinB)
		}
		if *setSkinA != -1 {
			custom.SkinColor[3] = clampColor(*setSkinA)
		}

		if *setHairR != -1 {
			custom.HairColor[0] = clampColor(*setHairR)
		}
		if *setHairG != -1 {
			custom.HairColor[1] = clampColor(*setHairG)
		}
		if *setHairB != -1 {
			custom.HairColor[2] = clampColor(*setHairB)
		}
		if *setHairA != -1 {
			custom.HairColor[3] = clampColor(*setHairA)
		}

		// Set stats only if flags are explicitly set
		if *setFunds != -1 {
			custom.Funds = *setFunds
		}
		if *setHRPoints != -1 {
			custom.HRPoints = *setHRPoints
		}
		if *setHunterRank != -1 {
			custom.HunterRank = *setHunterRank
		}
		if *setAcademyPoints != -1 {
			custom.AcademyPoints = *setAcademyPoints
		}

		// Set village points only if flags are explicitly set
		if *setBherna != -1 {
			custom.BhernaPoints = *setBherna
		}
		if *setKokoto != -1 {
			custom.KokotoPoints = *setKokoto
		}
		if *setPokke != -1 {
			custom.PokkePoints = *setPokke
		}
		if *setYukumo != -1 {
			custom.YukumoPoints = *setYukumo
		}

		// Apply customization
		modifiedData, err := customizePlayer(player, custom)
		if err != nil {
			fmt.Printf("Error customizing player: %v\n", err)
			os.Exit(1)
		}

		// Save modified file
		finalData := make([]byte, len(saveData))
		if isSwitch {
			// Re-add Switch header
			copy(finalData, saveData[:36]) // Keep original header
			copy(finalData[36:], modifiedData)
		} else {
			copy(finalData, modifiedData)
		}

		err = ioutil.WriteFile(*outputFile, finalData, 0644)
		if err != nil {
			fmt.Printf("Error saving file: %v\n", err)
			os.Exit(1)
		}

		fmt.Printf("\nCustomization applied successfully!\n")
		fmt.Printf("Saved to: %s\n", *outputFile)

		// Reload to show changes
		player, _ = extractPlayerData(modifiedData, *slot)
		player.SaveData = modifiedData
	}

	// Display information based on flags
	if *displayAll || !(*displayItems || *displayEquips || *displayPalico) {
		// Default: show character info
		displayCharacterInfo(player, *slot, *debug)
	}

	if *displayAll || *displayItems {
		displayItemBoxInfo(player, *debug)
	}

	if *displayAll || *displayEquips {
		displayEquipmentInfo(player, *debug)
	}

	if *displayAll || *displayPalico {
		displayPalicoInfo(player, *debug)
	}

	// Show data section sizes
	if *displayAll {
		displayDataSections(player)
	}
}

func clampColor(value int) int {
	if value < 0 {
		return 0
	}
	if value > 255 {
		return 255
	}
	return value
}

func loadSaveFile(filename string) ([]byte, []byte, bool, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, nil, false, err
	}

	var extractedData []byte
	isSwitch := false

	// Check save type
	switch len(data) {
	case 4726152: // 3DS
		extractedData = data
	case 4726152 + 36: // Switch
		fallthrough
	case 4726152 + 432948: // MHGU
		if len(data) >= 36 {
			extractedData = data[36:]
			isSwitch = true
		} else {
			return nil, nil, false, fmt.Errorf("Switch save too small")
		}
	default:
		return nil, nil, false, fmt.Errorf("unknown save size: %d bytes", len(data))
	}

	return data, extractedData, isSwitch, nil
}

func extractPlayerData(data []byte, slot int) (*Player, error) {
	if slot < 1 || slot > 3 {
		return nil, fmt.Errorf("invalid slot: %d", slot)
	}

	// Get slot offset
	var slotOffset int
	switch slot {
	case 1:
		if len(data) >= 0x14 {
			slotOffset = int(binary.LittleEndian.Uint32(data[0x10:]))
		}
	case 2:
		if len(data) >= 0x18 {
			slotOffset = int(binary.LittleEndian.Uint32(data[0x14:]))
		}
	case 3:
		if len(data) >= 0x1C {
			slotOffset = int(binary.LittleEndian.Uint32(data[0x18:]))
		}
	}

	if slotOffset == 0 || slotOffset >= len(data) {
		return nil, fmt.Errorf("invalid character offset: 0x%X", slotOffset)
	}

	player := &Player{
		SaveOffset: slotOffset,
	}

	// Extract basic info
	extractBasicInfo(player, data)

	// Extract item box (simplified)
	player.ItemId = make([]string, 2300)
	player.ItemCount = make([]string, 2300)
	for i := 0; i < 2300; i++ {
		player.ItemId[i] = "0"
		player.ItemCount[i] = "0"
	}

	// Extract other data sections
	extractOtherData(player, data)

	return player, nil
}

func extractBasicInfo(player *Player, data []byte) {
	offset := player.SaveOffset

	// Name - using the offset we found earlier (0x18CC78 from character base)
	// First try to find name by searching
	nameOffset := findNameOffset(data, offset)
	if nameOffset > 0 {
		player.Name = extractNullTerminatedString(data, nameOffset, 32)
	}

	// Play time
	if offset+0x24 <= len(data) {
		player.PlayTime = int(binary.LittleEndian.Uint32(data[offset+0x20:]))
	}

	// Funds
	if offset+0x28 <= len(data) {
		player.Funds = int(binary.LittleEndian.Uint32(data[offset+0x24:]))
	}

	// Hunter rank
	if offset+0x2A <= len(data) {
		player.HunterRank = int(binary.LittleEndian.Uint16(data[offset+0x28:]))
	}

	// HR points
	if offset+0x280F <= len(data) {
		player.HRPoints = int(binary.LittleEndian.Uint32(data[offset+0x280B:]))
	}

	// Academy points
	if offset+0x281B <= len(data) {
		player.AcademyPoints = int(binary.LittleEndian.Uint32(data[offset+0x2817:]))
	}

	// Village points
	if offset+0x282B <= len(data) {
		player.BhernaPoints = int(binary.LittleEndian.Uint32(data[offset+0x281B:]))
		player.KokotoPoints = int(binary.LittleEndian.Uint32(data[offset+0x281F:]))
		player.PokkePoints = int(binary.LittleEndian.Uint32(data[offset+0x2823:]))
		player.YukumoPoints = int(binary.LittleEndian.Uint32(data[offset+0x2827:]))
	}

	// Appearance - using C# offsets
	if offset+0x23B50 <= len(data) {
		player.Voice = data[offset+0x23B48]
		player.EyeColor = data[offset+0x23B49]
		player.Clothing = data[offset+0x23B4A]
		player.Gender = data[offset+0x23B4B]
		player.HairStyle = data[offset+0x23B4D]
		player.Face = data[offset+0x23B4E]
		player.Features = data[offset+0x23B4F]
	}

	// Colors
	if offset+0x23B77 <= len(data) {
		copy(player.SkinColorRGBA[:], data[offset+0x23B67:offset+0x23B67+4])
		copy(player.HairColorRGBA[:], data[offset+0x23B6B:offset+0x23B6B+4])
		copy(player.FeaturesColorRGBA[:], data[offset+0x23B6F:offset+0x23B6F+4])
		copy(player.ClothingColorRGBA[:], data[offset+0x23B73:offset+0x23B73+4])
	}
}

func findNameOffset(data []byte, charOffset int) int {
	// Search for printable strings near character offset
	searchStart := charOffset
	searchEnd := charOffset + 0x100000 // Search 1MB forward
	if searchEnd > len(data) {
		searchEnd = len(data)
	}

	// Look for the longest printable string
	bestOffset := -1
	bestLength := 0

	for i := searchStart; i < searchEnd; i++ {
		// Start of string check
		if i > 0 && data[i-1] != 0 {
			continue
		}

		// Check for printable string
		length := 0
		for j := i; j < searchEnd && j-i < 32; j++ {
			b := data[j]
			if b >= 32 && b <= 126 {
				length++
			} else if b == 0 {
				// Found null terminator
				if length > 2 && length > bestLength {
					bestOffset = i
					bestLength = length
				}
				break
			} else {
				// Non-printable
				break
			}
		}
	}

	return bestOffset
}

func extractOtherData(player *Player, data []byte) {
	offset := player.SaveOffset

	// Equipment box
	equipOffset := offset + 0x62EE
	if equipOffset+72000 <= len(data) {
		player.EquipmentInfo = make([]byte, 72000)
		copy(player.EquipmentInfo, data[equipOffset:equipOffset+72000])
	} else {
		player.EquipmentInfo = make([]byte, 0)
	}

	// Palico equipment
	palicoEquipOffset := offset + 0x17C2E
	if palicoEquipOffset+36000 <= len(data) {
		player.EquipmentPalico = make([]byte, 36000)
		copy(player.EquipmentPalico, data[palicoEquipOffset:palicoEquipOffset+36000])
	} else {
		player.EquipmentPalico = make([]byte, 0)
	}

	// Palico data
	palicoOffset := offset + 0x23BB6
	if palicoOffset+27216 <= len(data) {
		player.PalicoData = make([]byte, 27216)
		copy(player.PalicoData, data[palicoOffset:palicoOffset+27216])
	} else {
		player.PalicoData = make([]byte, 0)
	}

	// Guild card
	guildCardOffset := offset + 0xC71BD
	if guildCardOffset+4986 <= len(data) {
		player.GuildCardData = make([]byte, 4986)
		copy(player.GuildCardData, data[guildCardOffset:guildCardOffset+4986])
	} else {
		player.GuildCardData = make([]byte, 0)
	}

	// Initialize other arrays
	player.ArenaData = make([]byte, 0)
	player.MonsterKills = make([]byte, 0)
	player.MonsterCaptures = make([]byte, 0)
	player.MonsterSizes = make([]byte, 0)
	player.ManualShoutouts = make([]byte, 0)
	player.AutomaticShoutouts = make([]byte, 0)
}

func customizePlayer(player *Player, custom Customization) ([]byte, error) {
	// Create a copy of the save data
	modifiedData := make([]byte, len(player.SaveData))
	copy(modifiedData, player.SaveData)

	offset := player.SaveOffset

	// Apply name change
	if custom.Name != "" && len(custom.Name) <= 32 {
		// First find where the name is actually stored
		nameOffset := findNameOffset(modifiedData, offset)
		if nameOffset > 0 && nameOffset+32 <= len(modifiedData) {
			// Clear existing name
			for i := 0; i < 32; i++ {
				modifiedData[nameOffset+i] = 0
			}
			// Write new name
			nameBytes := []byte(custom.Name)
			copy(modifiedData[nameOffset:], nameBytes)
			player.Name = custom.Name
		}
	}

	// Apply appearance changes
	if custom.Gender != -1 {
		if offset+0x23B4C <= len(modifiedData) {
			modifiedData[offset+0x23B4B] = byte(custom.Gender)
			player.Gender = byte(custom.Gender)
		}
	}

	if custom.Voice != -1 {
		if offset+0x23B49 <= len(modifiedData) {
			modifiedData[offset+0x23B48] = byte(custom.Voice)
			player.Voice = byte(custom.Voice)
		}
	}

	if custom.EyeColor != -1 {
		if offset+0x23B4A <= len(modifiedData) {
			modifiedData[offset+0x23B49] = byte(custom.EyeColor)
			player.EyeColor = byte(custom.EyeColor)
		}
	}

	if custom.Clothing != -1 {
		if offset+0x23B4B <= len(modifiedData) {
			modifiedData[offset+0x23B4A] = byte(custom.Clothing)
			player.Clothing = byte(custom.Clothing)
		}
	}

	if custom.HairStyle != -1 {
		if offset+0x23B4E <= len(modifiedData) {
			modifiedData[offset+0x23B4D] = byte(custom.HairStyle)
			player.HairStyle = byte(custom.HairStyle)
		}
	}

	if custom.Face != -1 {
		if offset+0x23B4F <= len(modifiedData) {
			modifiedData[offset+0x23B4E] = byte(custom.Face)
			player.Face = byte(custom.Face)
		}
	}

	if custom.Features != -1 {
		if offset+0x23B50 <= len(modifiedData) {
			modifiedData[offset+0x23B4F] = byte(custom.Features)
			player.Features = byte(custom.Features)
		}
	}

	// Apply color changes
	if custom.SkinColor[0] != -1 || custom.SkinColor[1] != -1 ||
		custom.SkinColor[2] != -1 || custom.SkinColor[3] != -1 {
		skinOffset := offset + 0x23B67
		if skinOffset+4 <= len(modifiedData) {
			for i := 0; i < 4; i++ {
				if custom.SkinColor[i] != -1 {
					modifiedData[skinOffset+i] = byte(custom.SkinColor[i])
					player.SkinColorRGBA[i] = byte(custom.SkinColor[i])
				}
			}
		}
	}

	if custom.HairColor[0] != -1 || custom.HairColor[1] != -1 ||
		custom.HairColor[2] != -1 || custom.HairColor[3] != -1 {
		hairOffset := offset + 0x23B6B
		if hairOffset+4 <= len(modifiedData) {
			for i := 0; i < 4; i++ {
				if custom.HairColor[i] != -1 {
					modifiedData[hairOffset+i] = byte(custom.HairColor[i])
					player.HairColorRGBA[i] = byte(custom.HairColor[i])
				}
			}
		}
	}

	// Apply stat changes
	if custom.Funds != -1 {
		if offset+0x28 <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x24:], uint32(custom.Funds))
			player.Funds = custom.Funds
		}
	}

	if custom.HRPoints != -1 {
		if offset+0x280F <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x280B:], uint32(custom.HRPoints))
			player.HRPoints = custom.HRPoints
		}
	}

	if custom.HunterRank != -1 {
		if offset+0x2A <= len(modifiedData) {
			binary.LittleEndian.PutUint16(modifiedData[offset+0x28:], uint16(custom.HunterRank))
			player.HunterRank = custom.HunterRank
		}
	}

	if custom.AcademyPoints != -1 {
		if offset+0x281B <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x2817:], uint32(custom.AcademyPoints))
			player.AcademyPoints = custom.AcademyPoints
		}
	}

	// Apply village point changes
	if custom.BhernaPoints != -1 {
		if offset+0x281F <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x281B:], uint32(custom.BhernaPoints))
			player.BhernaPoints = custom.BhernaPoints
		}
	}

	if custom.KokotoPoints != -1 {
		if offset+0x2823 <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x281F:], uint32(custom.KokotoPoints))
			player.KokotoPoints = custom.KokotoPoints
		}
	}

	if custom.PokkePoints != -1 {
		if offset+0x2827 <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x2823:], uint32(custom.PokkePoints))
			player.PokkePoints = custom.PokkePoints
		}
	}

	if custom.YukumoPoints != -1 {
		if offset+0x282B <= len(modifiedData) {
			binary.LittleEndian.PutUint32(modifiedData[offset+0x2827:], uint32(custom.YukumoPoints))
			player.YukumoPoints = custom.YukumoPoints
		}
	}

	// Also update guild card appearance
	updateGuildCardAppearance(player, custom, modifiedData)

	return modifiedData, nil
}

func updateGuildCardAppearance(player *Player, custom Customization, data []byte) {
	offset := player.SaveOffset
	guildCardOffset := offset + 0xC71BD

	// Update guild card gender
	if custom.Gender != -1 && guildCardOffset+0xC71D9 <= len(data) {
		data[guildCardOffset+0xC71D9] = byte(custom.Gender)
	}

	// Update guild card voice
	if custom.Voice != -1 && guildCardOffset+0xC71D6 <= len(data) {
		data[guildCardOffset+0xC71D6] = byte(custom.Voice)
	}

	// Update guild card eye color
	if custom.EyeColor != -1 && guildCardOffset+0xC71D7 <= len(data) {
		data[guildCardOffset+0xC71D7] = byte(custom.EyeColor)
	}

	// Update guild card clothing
	if custom.Clothing != -1 && guildCardOffset+0xC71D8 <= len(data) {
		data[guildCardOffset+0xC71D8] = byte(custom.Clothing)
	}

	// Update guild card hair style
	if custom.HairStyle != -1 && guildCardOffset+0xC71DB <= len(data) {
		data[guildCardOffset+0xC71DB] = byte(custom.HairStyle)
	}

	// Update guild card face
	if custom.Face != -1 && guildCardOffset+0xC71DC <= len(data) {
		data[guildCardOffset+0xC71DC] = byte(custom.Face)
	}

	// Update guild card features
	if custom.Features != -1 && guildCardOffset+0xC71DD <= len(data) {
		data[guildCardOffset+0xC71DD] = byte(custom.Features)
	}

	// Update guild card colors
	if custom.SkinColor[0] != -1 || custom.SkinColor[1] != -1 ||
		custom.SkinColor[2] != -1 || custom.SkinColor[3] != -1 {
		skinOffset := guildCardOffset + 0xC71F5
		if skinOffset+4 <= len(data) {
			for i := 0; i < 4; i++ {
				if custom.SkinColor[i] != -1 {
					data[skinOffset+i] = byte(custom.SkinColor[i])
				}
			}
		}
	}

	if custom.HairColor[0] != -1 || custom.HairColor[1] != -1 ||
		custom.HairColor[2] != -1 || custom.HairColor[3] != -1 {
		hairOffset := guildCardOffset + 0xC71F9
		if hairOffset+4 <= len(data) {
			for i := 0; i < 4; i++ {
				if custom.HairColor[i] != -1 {
					data[hairOffset+i] = byte(custom.HairColor[i])
				}
			}
		}
	}
}

func displayCharacterInfo(player *Player, slot int, debug bool) {
	fmt.Printf("\n=== CHARACTER SLOT %d ===\n", slot)
	fmt.Printf("Save Offset: 0x%08X\n", player.SaveOffset)

	// Basic Info
	fmt.Printf("\n--- BASIC INFORMATION ---\n")
	fmt.Printf("Name:          %s\n", player.Name)
	fmt.Printf("Play Time:     %s\n", formatPlayTime(player.PlayTime))
	fmt.Printf("Funds:         %dz\n", player.Funds)
	fmt.Printf("Hunter Rank:   %d\n", player.HunterRank)
	fmt.Printf("HR Points:     %d\n", player.HRPoints)
	fmt.Printf("Academy Points:%d\n", player.AcademyPoints)

	// Village Points
	fmt.Printf("\n--- VILLAGE POINTS ---\n")
	fmt.Printf("Bherna:   %d\n", player.BhernaPoints)
	fmt.Printf("Kokoto:   %d\n", player.KokotoPoints)
	fmt.Printf("Pokke:    %d\n", player.PokkePoints)
	fmt.Printf("Yukumo:   %d\n", player.YukumoPoints)

	// Appearance
	fmt.Printf("\n--- APPEARANCE ---\n")
	fmt.Printf("Gender:        %d (%s)\n", player.Gender, getGenderName(player.Gender))
	fmt.Printf("Voice:         %d\n", player.Voice)
	fmt.Printf("Eye Color:     %d\n", player.EyeColor)
	fmt.Printf("Clothing:      %d\n", player.Clothing)
	fmt.Printf("Hair Style:    %d\n", player.HairStyle)
	fmt.Printf("Face:          %d\n", player.Face)
	fmt.Printf("Features:      %d\n", player.Features)

	// Colors
	fmt.Printf("\n--- COLORS (RGBA) ---\n")
	fmt.Printf("Skin:      R:%3d G:%3d B:%3d A:%3d\n",
		player.SkinColorRGBA[0], player.SkinColorRGBA[1],
		player.SkinColorRGBA[2], player.SkinColorRGBA[3])
	fmt.Printf("Hair:      R:%3d G:%3d B:%3d A:%3d\n",
		player.HairColorRGBA[0], player.HairColorRGBA[1],
		player.HairColorRGBA[2], player.HairColorRGBA[3])
	fmt.Printf("Features:  R:%3d G:%3d B:%3d A:%3d\n",
		player.FeaturesColorRGBA[0], player.FeaturesColorRGBA[1],
		player.FeaturesColorRGBA[2], player.FeaturesColorRGBA[3])
	fmt.Printf("Clothing:  R:%3d G:%3d B:%3d A:%3d\n",
		player.ClothingColorRGBA[0], player.ClothingColorRGBA[1],
		player.ClothingColorRGBA[2], player.ClothingColorRGBA[3])

	if debug {
		fmt.Printf("\n--- DEBUG INFO ---\n")
		fmt.Printf("Name found at offset: 0x%08X\n", findNameOffset(player.SaveData, player.SaveOffset))
	}
}

func getGenderName(gender byte) string {
	if gender == 0 {
		return "Male"
	} else if gender == 1 {
		return "Female"
	}
	return "Unknown"
}

func displayItemBoxInfo(player *Player, debug bool) {
	fmt.Printf("\n=== ITEM BOX ===\n")
	fmt.Printf("Total slots: 2300\n")

	// Count non-empty items
	nonEmpty := 0
	for i := 0; i < 2300; i++ {
		if player.ItemId[i] != "0" && player.ItemCount[i] != "0" {
			nonEmpty++
		}
	}

	fmt.Printf("Non-empty items: %d\n", nonEmpty)
	fmt.Printf("Empty slots: %d\n", 2300-nonEmpty)
}

func displayEquipmentInfo(player *Player, debug bool) {
	fmt.Printf("\n=== EQUIPMENT BOX ===\n")
	fmt.Printf("Total slots: 2000\n")
	fmt.Printf("Data size: %d bytes (36 bytes per equipment)\n", len(player.EquipmentInfo))

	// Count non-empty equipment
	nonEmpty := 0
	for i := 0; i < 2000; i++ {
		if len(player.EquipmentInfo) > i*36 && player.EquipmentInfo[i*36] != 0 {
			nonEmpty++
		}
	}

	fmt.Printf("Non-empty equipment: %d\n", nonEmpty)
	fmt.Printf("Empty slots: %d\n", 2000-nonEmpty)

	// Palico equipment
	fmt.Printf("\n=== PALICO EQUIPMENT ===\n")
	fmt.Printf("Total slots: 1000\n")
	fmt.Printf("Data size: %d bytes\n", len(player.EquipmentPalico))

	palicoNonEmpty := 0
	for i := 0; i < 1000; i++ {
		if len(player.EquipmentPalico) > i*36 && player.EquipmentPalico[i*36] != 0 {
			palicoNonEmpty++
		}
	}

	fmt.Printf("Non-empty palico equipment: %d\n", palicoNonEmpty)
	fmt.Printf("Empty slots: %d\n", 1000-palicoNonEmpty)
}

func displayPalicoInfo(player *Player, debug bool) {
	fmt.Printf("\n=== PALICO DATA ===\n")
	fmt.Printf("Total slots: 84\n")
	fmt.Printf("Data size: %d bytes (324 bytes per palico)\n", len(player.PalicoData))

	// Count palicos with names
	palicoCount := 0
	for i := 0; i < 84; i++ {
		if len(player.PalicoData) > i*324 && player.PalicoData[i*324] != 0 {
			palicoCount++
		}
	}

	fmt.Printf("Palicos with names: %d\n", palicoCount)
}

func displayDataSections(player *Player) {
	fmt.Printf("\n=== DATA SECTIONS ===\n")
	fmt.Printf("Item Box:          %7d bytes\n", 5463)
	fmt.Printf("Equipment Box:     %7d bytes\n", 72000)
	fmt.Printf("Palico Equipment:  %7d bytes\n", 36000)
	fmt.Printf("Palico Data:       %7d bytes\n", 27216)
	fmt.Printf("Monster Kills:     %7d bytes\n", 274)
	fmt.Printf("Monster Captures:  %7d bytes\n", 274)
	fmt.Printf("Monster Sizes:     %7d bytes\n", 548)
	fmt.Printf("Guild Card:        %7d bytes\n", 4986)
	fmt.Printf("Arena Data:        %7d bytes\n", 342)
	fmt.Printf("Manual Shoutouts:  %7d bytes\n", 2880)
	fmt.Printf("Auto Shoutouts:    %7d bytes\n", 1620)

	// Show actual loaded sizes
	fmt.Printf("\n--- ACTUALLY LOADED ---\n")
	fmt.Printf("Equipment Info:    %7d bytes\n", len(player.EquipmentInfo))
	fmt.Printf("Equipment Palico:  %7d bytes\n", len(player.EquipmentPalico))
	fmt.Printf("Palico Data:       %7d bytes\n", len(player.PalicoData))
	fmt.Printf("Guild Card:        %7d bytes\n", len(player.GuildCardData))
}

// Helper functions
func extractNullTerminatedString(data []byte, offset, maxLen int) string {
	if offset < 0 || offset >= len(data) {
		return ""
	}

	end := offset
	for end < len(data) && end-offset < maxLen && data[end] != 0 {
		end++
	}

	if end == offset {
		return ""
	}

	return string(data[offset:end])
}

func formatPlayTime(seconds int) string {
	hours := seconds / 3600
	minutes := (seconds % 3600) / 60
	secs := seconds % 60
	return fmt.Sprintf("%d:%02d:%02d", hours, minutes, secs)
}

// Add these helper functions at the end of the file, before the last closing brace:

// Search for all occurrences of bytes in data
func searchInMemory(data []byte, search []byte) []int {
	var results []int

	if len(search) == 0 || len(data) < len(search) {
		return results
	}

	for i := 0; i <= len(data)-len(search); i++ {
		found := true
		for j := 0; j < len(search); j++ {
			if data[i+j] != search[j] {
				found = false
				break
			}
		}
		if found {
			results = append(results, i)
		}
	}

	return results
}

// Replace occurrences at specific offsets
func replaceInMemory(data []byte, search []byte, replace []byte, offsets []int) []byte {
	// Create a copy of the data
	result := make([]byte, len(data))
	copy(result, data)

	replaced := 0
	for _, offset := range offsets {
		// Check bounds
		if offset < 0 || offset+len(replace) > len(result) {
			continue
		}

		// Verify the search string still exists at this offset
		matches := true
		for j := 0; j < len(search); j++ {
			if result[offset+j] != search[j] {
				matches = false
				break
			}
		}

		if matches {
			// Perform replacement
			for j := 0; j < len(replace); j++ {
				result[offset+j] = replace[j]
			}
			replaced++
		}
	}

	return result
}

// Check if a location looks like a valid name field
func isValidNameLocation(data []byte, offset int, searchLen int) bool {
	if offset < 0 || offset >= len(data) {
		return false
	}

	// Check if it's null-terminated
	if offset+searchLen < len(data) && data[offset+searchLen] == 0 {
		return true
	}

	// Check if it's within a reasonable area (name fields are usually 32 bytes)
	// and has mostly zeros around it
	start := offset - 16
	if start < 0 {
		start = 0
	}
	end := offset + 48
	if end > len(data) {
		end = len(data)
	}

	// Count zeros in the area
	zeroCount := 0
	totalCount := end - start
	for i := start; i < end; i++ {
		if data[i] == 0 {
			zeroCount++
		}
	}

	// If more than 50% zeros, likely a name field
	return float64(zeroCount)/float64(totalCount) > 0.5
}

// Extract string with context
func extractStringContext(data []byte, offset, searchLen, contextSize int) string {
	if offset < 0 || offset >= len(data) {
		return ""
	}

	start := offset - contextSize
	if start < 0 {
		start = 0
	}

	end := offset + searchLen + contextSize
	if end > len(data) {
		end = len(data)
	}

	// Extract and clean up non-printable characters for display
	result := make([]byte, 0, end-start)
	for i := start; i < end; i++ {
		b := data[i]
		if b >= 32 && b <= 126 {
			result = append(result, b)
		} else if b == 0 {
			result = append(result, '.')
		} else {
			result = append(result, '?')
		}
	}

	return string(result)
}

// Hex dump helper
func hexDump(data []byte, start, length int, ascii bool) {
	if start < 0 {
		start = 0
	}

	if start >= len(data) {
		return
	}

	end := start + length
	if end > len(data) {
		end = len(data)
	}

	for i := start; i < end; i += 16 {
		fmt.Printf("  %08X: ", i)

		// Hex bytes
		for j := 0; j < 16; j++ {
			if i+j < end {
				fmt.Printf("%02X ", data[i+j])
			} else {
				fmt.Print("   ")
			}
		}

		// ASCII
		if ascii {
			fmt.Print(" ")
			for j := 0; j < 16; j++ {
				if i+j < end {
					b := data[i+j]
					if b >= 32 && b <= 126 {
						fmt.Printf("%c", b)
					} else {
						fmt.Print(".")
					}
				}
			}
		}
		fmt.Println()
	}
}

// Helper function to convert [4]byte to [4]int
func byteArrayToIntArray(b [4]byte) [4]int {
	return [4]int{int(b[0]), int(b[1]), int(b[2]), int(b[3])}
}
