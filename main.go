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

type RigStats struct { GPU string `json:"gpu"`; VPN string `json:"vpn"` }
type PlayerStats struct { Cooldown int `json:"cooldown"`; Luck int `json:"luck"`; BuffTime int `json:"buffTime"` }
type ProfileResponse struct { Address string `json:"address"`; Balance string `json:"balance"`; Stats PlayerStats `json:"stats"`; Rig RigStats `json:"rig"` }
type InventoryItem struct { ID int `json:"id"`; Name string `json:"name"`; Type string `json:"type"`; Qty int `json:"qty"` }
type MineResult struct { Name string `json:"name"`; Type string `json:"type"`; Stats string `json:"stats"`; IsEquipment bool `json:"isEquipment"` }
type MineResponse struct { Success bool `json:"success"`; Data *MineResult `json:"data"`; Error string `json:"error"` }
type Mission struct { ID int `json:"id"`; Desc string `json:"desc"`; Progress int `json:"progress"`; Target int `json:"target"`; Reward string `json:"reward"`; Claimed bool `json:"claimed"` }
type ContractsRes struct { Date string `json:"date"`; Missions []Mission `json:"missions"` }
type LeaderboardRes struct { Top10 []struct { Address string `json:"address"`; Balance string `json:"balance"` } `json:"top10"` }

func main() {
	pk, addr := loadWallet()
	for {
		clearScreen()
		printBanner()
		printHUD(addr)

		fmt.Println("\n[SYSTEM MENU]")
		fmt.Println("1. " + bold("HACK_NODE") + "    (Mine - 30s CD)")
		fmt.Println("2. " + bold("CYBERDECK") + "    (Equip Hardware)")
		fmt.Println("3. " + bold("WORKSHOP") + "     (Overclock/Enchant)")
		fmt.Println("4. " + bold("INVENTORY") + "    (Salvage / Use)")
		fmt.Println("5. " + bold("SHADOW_NET") + "   (Missions & Leaderboard)")
		fmt.Println("6. " + bold("EXIT"))
		fmt.Print("\nrunner@terminal:~$ ")

		var choice string; fmt.Scanln(&choice)
		switch choice {
		case "1": performHack(pk, addr)
		case "2": manageRig(pk, addr)
		case "3": openWorkshop(pk, addr)
		case "4": manageInventory(pk, addr)
		case "5": shadowNet(pk, addr)
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

func shadowNet(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(cyan("\n[A] Daily Contracts"))
	fmt.Println(cyan("[B] Global Leaderboard"))
	fmt.Print("Choice: "); var c string; fmt.Scanln(&c)
	if c == "A" || c == "a" {
		resp, _ := http.Get(SERVER_URL + "/contracts/" + addr)
		var r ContractsRes
		json.NewDecoder(resp.Body).Decode(&r)
		fmt.Println(bold("\n=== DAILY CONTRACTS ==="))
		for _, m := range r.Missions {
			st := red("[INCOMPLETE]")
			if m.Claimed { st = green("[CLAIMED]") } else if m.Progress >= m.Target { st = yellow("[READY]") }
			fmt.Printf("[%d] %s (%d/%d) -> %s %s\n", m.ID, m.Desc, m.Progress, m.Target, m.Reward, st)
		}
		fmt.Print("\nClaim Mission ID (0 back): "); var mid string; fmt.Scanln(&mid)
		if mid!="0" {
			midInt, _ := strconv.Atoi(mid)
			sig := sign(pk, "CLAIM_MISSION")
			printRes(postRequest("/contracts/claim", map[string]interface{}{"userAddress":addr, "signature":sig, "missionId":midInt}))
		}
	} else {
		resp, _ := http.Get(SERVER_URL + "/leaderboard")
		var l LeaderboardRes
		json.NewDecoder(resp.Body).Decode(&l)
		fmt.Println(bold("\n=== TOP 10 HACKERS ==="))
		for i, u := range l.Top10 {
			fmt.Printf("#%d %-42s %s\n", i+1, u.Address, yellow(u.Balance))
		}
	}
}

func performHack(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(yellow("\n>> HACKING..."))
	for i := 0; i < 10; i++ { fmt.Print("█"); time.Sleep(50 * time.Millisecond) }
	fmt.Println("")
	sig := sign(pk, "MINT_ACTION")
	res := postRequest("/mine", map[string]string{"userAddress": addr, "signature": sig})
	var m MineResponse
	json.Unmarshal(res, &m)

	if m.Success {
		d := m.Data
		if d.IsEquipment {
			fmt.Println(magenta("\n>>> JACKPOT! HARDWARE! <<<"))
			fmt.Printf("[%s] %s\n", d.Type, bold(d.Name))
		} else if d.Type == "MAT" || d.Type == "CONS" {
			fmt.Println(cyan("\n>>> ITEM FOUND <<<"))
			fmt.Printf("[%s] %s\n", d.Type, bold(d.Name))
		} else {
			fmt.Println(green("\n>> ARTIFACT FOUND <<"))
			fmt.Printf("%s (%s)\n", d.Name, d.Stats)
		}
	} else { fmt.Println(red("FAILED: "), m.Error) }
}

func openWorkshop(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(magenta("\n=== WORKSHOP ==="))
	resp, _ := http.Get(SERVER_URL + "/inventory/" + addr)
	type InvRes struct { Items []InventoryItem `json:"items"` }
	var inv InvRes
	json.NewDecoder(resp.Body).Decode(&inv)
	for _, it := range inv.Items { if it.Type != "ARTIFACT" { fmt.Printf("[%d] %s x%d\n", it.ID, it.Name, it.Qty) } }

	fmt.Print("\nTarget ID: "); var tId string; fmt.Scanln(&tId)
	fmt.Print("Material ID: "); var mId string; fmt.Scanln(&mId)
	sig := sign(pk, "ENCHANT_ACTION")
	res := postRequest("/workshop/enchant", map[string]string{"userAddress": addr, "signature": sig, "targetId": tId, "materialId": mId})
	var r struct { Success bool; Message string; Level int }
	json.Unmarshal(res, &r)
	if r.Success { fmt.Println(green("SUCCESS!"), r.Message) } else { fmt.Println(red("FAIL:"), r.Message) }
}

func manageRig(pk *ecdsa.PrivateKey, addr string) {
	fmt.Print("\nEquip ID (0 cancel): "); var id string; fmt.Scanln(&id)
	if id != "0" {
		sig := sign(pk, "EQUIP_ACTION")
		printRes(postRequest("/equip", map[string]string{"userAddress": addr, "signature": sig, "itemId": id}))
	}
}

func manageInventory(pk *ecdsa.PrivateKey, addr string) {
	resp, _ := http.Get(SERVER_URL + "/inventory/" + addr)
	type InvRes struct { Items []InventoryItem `json:"items"` }
	var inv InvRes
	json.NewDecoder(resp.Body).Decode(&inv)
	for _, it := range inv.Items { fmt.Printf("[%d] %s (%s) x%d\n", it.ID, it.Name, it.Type, it.Qty) }

	fmt.Println("\n[S] Salvage  [U] Use Item")
	var c string; fmt.Scanln(&c)
	if c=="U" || c=="u" {
		fmt.Print("Item ID: "); var id string; fmt.Scanln(&id)
		sig := sign(pk, "USE_ACTION")
		printRes(postRequest("/use", map[string]string{"userAddress": addr, "signature": sig, "itemId": id}))
	} else {
		fmt.Print("ID: "); var id string; fmt.Scanln(&id)
		fmt.Print("Amount: "); var amt string; fmt.Scanln(&amt)
		amtInt, _ := strconv.Atoi(amt)
		sig := sign(pk, "SALVAGE_ACTION")
		printRes(postRequest("/salvage", map[string]interface{}{"userAddress": addr, "signature": sig, "tokenId": id, "amount": amtInt}))
	}
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
func printBanner() { fmt.Println(cyan(` :: RNG MINER v6.0 :: SHADOW NET :: `)) }