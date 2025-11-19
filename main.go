package main

import (
	"bytes"
	"crypto/ecdsa"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/fatih/color"
)

const (
	SERVER_URL = "http://localhost:3000"
	KEY_FILE   = "session.key"
)

var (
	cyan    = color.New(color.FgCyan).SprintFunc()
	green   = color.New(color.FgGreen).SprintFunc()
	red     = color.New(color.FgRed).SprintFunc()
	yellow  = color.New(color.FgYellow).SprintFunc()
	magenta = color.New(color.FgMagenta).SprintFunc()
	white   = color.New(color.FgWhite).SprintFunc()
	bold    = color.New(color.Bold).SprintFunc()
)

// Structs
type RigStats struct { GPU string `json:"gpu"`; VPN string `json:"vpn"` }
type PlayerStats struct { Cooldown int `json:"cooldown"`; Luck int `json:"luck"`; BuffTime int `json:"buffTime"` }
type ProfileResponse struct { Address string `json:"address"`; Balance string `json:"balance"`; Stats PlayerStats `json:"stats"`; Rig RigStats `json:"rig"` }
type InventoryItem struct { ID int `json:"id"`; Name string `json:"name"`; Type string `json:"type"`; Qty int `json:"qty"` }
type MineResult struct { Name string `json:"name"`; Type string `json:"type"`; Stats string `json:"stats"`; IsEquipment bool `json:"isEquipment"` }
type MineResponse struct { Success bool `json:"success"`; Data *MineResult `json:"data"`; Error string `json:"error"` }

// Struct Leaderboard
type HackerRank struct { Address string `json:"address"`; Balance string `json:"balance"` }
type LeaderboardResponse struct { Top10 []HackerRank `json:"top10"` }

func main() {
	pk, addr := loadWallet()
	for {
		clearScreen()
		printBanner()
		printHUD(addr)

		fmt.Println("\n[SYSTEM MENU]")
		fmt.Println("1. " + bold("HACK_NODE") + "    (Mine - 5s CD)")
		fmt.Println("2. " + bold("CYBERDECK") + "    (Equip Hardware)")
		fmt.Println("3. " + bold("WORKSHOP") + "     (Overclock/Enchant)")
		fmt.Println("4. " + bold("INVENTORY") + "    (Salvage / View)")
		fmt.Println("5. " + bold("NETWORK") + "      (Global Leaderboard)") // <-- NEW
		fmt.Println("6. " + bold("EXIT"))
		fmt.Print("\nrunner@terminal:~$ ")

		var choice string; fmt.Scanln(&choice)
		switch choice {
		case "1": performHack(pk, addr)
		case "2": manageRig(pk, addr)
		case "3": openWorkshop(pk, addr)
		case "4": manageInventory(pk, addr)
		case "5": showLeaderboard(addr)
		case "6": os.Exit(0)
		}
		waitForKey()
	}
}

func printHUD(addr string) {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(SERVER_URL + "/profile/" + addr)
	if err != nil { fmt.Println(red("OFFLINE")); return }
	defer resp.Body.Close()
	var p ProfileResponse
	json.NewDecoder(resp.Body).Decode(&p)

	fmt.Println(cyan("=================================================="))
	fmt.Printf(" OPERATOR : %s\n", green(addr))
	fmt.Printf(" BALANCE  : %s\n", yellow(p.Balance+" $HASH"))
	fmt.Println(cyan("=================================================="))
	fmt.Println(bold(" [ ACTIVE RIG ]"))
	fmt.Printf(" GPU : %s\n", magenta(p.Rig.GPU))
	fmt.Printf(" VPN : %s\n", magenta(p.Rig.VPN))
	fmt.Println(bold("\n [ STATS ]"))
	fmt.Printf(" LUCK : %d\n", p.Stats.Luck)
	fmt.Printf(" CD   : %ds\n", p.Stats.Cooldown)
	if p.Stats.BuffTime > 0 { fmt.Printf(" BUFF : %s (%dm)\n", green("ACTIVE"), p.Stats.BuffTime/60) }
	fmt.Println(cyan("=================================================="))
}

func showLeaderboard(myAddr string) {
	fmt.Println(cyan("\nScanning Global Network..."))
	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Get(SERVER_URL + "/leaderboard")
	if err != nil { fmt.Println(red("Network Timeout.")); return }
	defer resp.Body.Close()
	
	var r LeaderboardResponse
	json.NewDecoder(resp.Body).Decode(&r)

	fmt.Println(bold("\n=== FBI MOST WANTED LIST ==="))
	fmt.Printf("%-4s | %-42s | %s\n", "RANK", "IDENTITY", "BOUNTY ($HASH)")
	fmt.Println("----------------------------------------------------------------")
	for i, hacker := range r.Top10 {
		rankStr := fmt.Sprintf("#%d", i+1)
		addrStr := hacker.Address
		if hacker.Address == myAddr {
			addrStr = green(hacker.Address + " (YOU)")
			rankStr = green(rankStr)
		}
		fmt.Printf("%-4s | %-50s | %s\n", rankStr, addrStr, yellow(hacker.Balance))
	}
	fmt.Println("----------------------------------------------------------------")
}

func performHack(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(yellow("\n>> BRUTE FORCE ATTACK..."))
	for i := 0; i < 10; i++ { fmt.Print("█"); time.Sleep(50 * time.Millisecond) }
	fmt.Println("")

	sig := sign(pk, "MINT_ACTION")
	res := postRequest("/mine", map[string]string{"userAddress": addr, "signature": sig})
	var m MineResponse
	json.Unmarshal(res, &m)

	if m.Success {
		d := m.Data
		if d.IsEquipment {
			fmt.Println(magenta("\n>>> JACKPOT! NEW HARDWARE! <<<"))
			fmt.Printf("[%s] %s (%s)\n", d.Type, bold(d.Name), d.Stats)
		} else if d.Type == "MAT" || d.Type == "SECRET" {
			fmt.Println(cyan("\n>>> MATERIAL FOUND! <<<"))
			fmt.Printf("[%s] %s\n", d.Type, bold(d.Name))
		} else {
			fmt.Println(green("\n>> DATA FOUND <<"))
			fmt.Printf("%s (%s)\n", d.Name, d.Stats)
		}
	} else {
		fmt.Println(red("\n>> FAILED: "), m.Error)
	}
}

func openWorkshop(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(magenta("\n=== HARDWARE WORKSHOP ==="))
	fmt.Println("Combine Equipment + Chip (401) or Core (99)")
	resp, _ := http.Get(SERVER_URL + "/inventory/" + addr)
	defer resp.Body.Close()
	type InvRes struct { Items []InventoryItem `json:"items"` }
	var inv InvRes
	json.NewDecoder(resp.Body).Decode(&inv)
	fmt.Println("\n[INVENTORY]")
	for _, it := range inv.Items { fmt.Printf("[%d] %s (%s) x%d\n", it.ID, white(it.Name), it.Type, it.Qty) }

	fmt.Print("\nTarget Item ID: "); var tId string; fmt.Scanln(&tId)
	fmt.Print("Material ID: "); var mId string; fmt.Scanln(&mId)
	fmt.Println("Overclocking...")
	sig := sign(pk, "ENCHANT_ACTION")
	res := postRequest("/workshop/enchant", map[string]string{"userAddress": addr, "signature": sig, "targetId": tId, "materialId": mId})
	var r struct { Success bool; Message string; Level int }
	json.Unmarshal(res, &r)
	if r.Success { fmt.Println(green("SUCCESS!"), r.Message) } else { fmt.Println(red("RESULT:"), r.Message) }
}

func manageRig(pk *ecdsa.PrivateKey, addr string) {
	fmt.Print("\nEnter ID to EQUIP (0 cancel): "); var id string; fmt.Scanln(&id)
	if id == "0" { return }
	sig := sign(pk, "EQUIP_ACTION")
	printRes(postRequest("/equip", map[string]string{"userAddress": addr, "signature": sig, "itemId": id}))
}

func manageInventory(pk *ecdsa.PrivateKey, addr string) {
	resp, _ := http.Get(SERVER_URL + "/inventory/" + addr)
	defer resp.Body.Close()
	type InvRes struct { Items []InventoryItem `json:"items"` }
	var inv InvRes
	json.NewDecoder(resp.Body).Decode(&inv)
	for _, it := range inv.Items { fmt.Printf("[%d] %s (%s) x%d\n", it.ID, it.Name, it.Type, it.Qty) }
	fmt.Print("\nSalvage ID: "); var id string; fmt.Scanln(&id)
	if id=="0" {return}
	fmt.Print("Amount: "); var amt string; fmt.Scanln(&amt)
	amtInt, _ := strconv.Atoi(amt)
	sig := sign(pk, "SALVAGE_ACTION")
	printRes(postRequest("/salvage", map[string]interface{}{"userAddress": addr, "signature": sig, "tokenId": id, "amount": amtInt}))
}

func loadWallet() (*ecdsa.PrivateKey, string) {
	if _, err := os.Stat(KEY_FILE); os.IsNotExist(err) {
		pk, _ := crypto.GenerateKey()
		ioutil.WriteFile(KEY_FILE, []byte(hex.EncodeToString(crypto.FromECDSA(pk))), 0600)
		return pk, crypto.PubkeyToAddress(pk.PublicKey).Hex()
	}
	kb, _ := ioutil.ReadFile(KEY_FILE)
	pk, _ := crypto.HexToECDSA(string(bytes.TrimSpace(kb)))
	return pk, crypto.PubkeyToAddress(pk.PublicKey).Hex()
}
func sign(pk *ecdsa.PrivateKey, msg string) string {
	d := []byte(msg); h := crypto.Keccak256Hash([]byte(fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(d))), d)
	s, _ := crypto.Sign(h.Bytes(), pk); if s[64]<27 {s[64]+=27}; return hexutil.Encode(s)
}
func postRequest(ep string, d interface{}) []byte {
	b, _ := json.Marshal(d); r, err := http.Post(SERVER_URL+ep, "application/json", bytes.NewBuffer(b))
	if err != nil { return nil }; defer r.Body.Close(); body, _ := ioutil.ReadAll(r.Body); return body
}
func printRes(b []byte) {
	if b==nil { fmt.Println(red("Error")); return }
	var r struct{Success bool; Message string; Error string}; json.Unmarshal(b, &r)
	if r.Success { fmt.Println(green("OK: "), r.Message) } else { fmt.Println(red("FAIL: "), r.Error) }
}
func waitForKey() { fmt.Print(color.HiBlackString("\n[ENTER]")); var s string; fmt.Scanln(&s) }
func clearScreen() { c:=exec.Command("clear"); if runtime.GOOS=="windows"{c=exec.Command("cmd","/c","cls")}; c.Stdout=os.Stdout; c.Run() }
func printBanner() { fmt.Println(cyan(` :: RNG MINER v5.1 :: MOST WANTED :: `)) }