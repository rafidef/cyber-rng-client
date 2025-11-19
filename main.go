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
type MineResponse struct { Success bool `json:"success"`; Data *MineResult `json:"data"`; Error string `json:"error"`; CooldownRemaining int `json:"cooldownRemaining"` }
type Mission struct { ID int `json:"id"`; Desc string `json:"desc"`; Progress int `json:"progress"`; Target int `json:"target"`; Reward string `json:"reward"`; Claimed bool `json:"claimed"` }
type ContractsRes struct { Date string `json:"date"`; Missions []Mission `json:"missions"` }
type LeaderboardRes struct { Top10 []struct { Address string `json:"address"`; Balance string `json:"balance"` } `json:"top10"` }
type StakeInfo struct { ID int `json:"id"`; Name string `json:"name"`; Yield string `json:"yieldRate"`; Staked int `json:"staked"`; Pending string `json:"pending"` }
type StakeRes struct { Stakes []StakeInfo `json:"stakes"` }

func main() {
	pk, addr := loadWallet()
	for {
		clearScreen()
		printBanner()
		printHUD(addr)

		fmt.Println("\n[SYSTEM MENU]")
		fmt.Println("1. " + bold("HACK_NODE") + "    (Mine)")
		fmt.Println("2. " + bold("CYBERDECK") + "    (Equip Hardware)")
		fmt.Println("3. " + bold("WORKSHOP") + "     (Overclock/Enchant)")
		fmt.Println("4. " + bold("INVENTORY") + "    (Salvage / Use)")
		fmt.Println("5. " + bold("SHADOW_NET") + "   (Missions & Leaderboard)")
		fmt.Println("6. " + bold("SERVER_ROOM") + "  (Passive Income)")
		fmt.Println("7. " + bold("BLACK_MARKET") + " (Buy Boosts & Tools)")
		fmt.Println("8. " + bold("EXIT"))
		fmt.Print("\nrunner@terminal:~$ ")

		var choice string; fmt.Scanln(&choice)
		switch choice {
		case "1": performHack(pk, addr)
		case "2": manageRig(pk, addr)
		case "3": openWorkshop(pk, addr)
		case "4": manageInventory(pk, addr)
		case "5": shadowNet(pk, addr)
		case "6": serverRoom(pk, addr)
		case "7": openShop(pk, addr)
		case "8": os.Exit(0)
		}
	}
}

func printHUD(addr string) {
	client := &http.Client{Timeout: 2 * time.Second}
	resp, err := client.Get(SERVER_URL + "/profile/" + addr)
	if err != nil { fmt.Println(red("OFFLINE - CANNOT CONNECT")); return }
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
	if p.Stats.BuffTime > 0 { 
		mins := p.Stats.BuffTime / 60
		fmt.Printf(" BOOST : %s (%dm left)\n", green("ACTIVE"), mins) 
	}
	fmt.Println(cyan("=================================================="))
}

func performHack(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(yellow("\n>> HACKING NODE..."))
	fmt.Print("[")
	for i := 0; i < 15; i++ { fmt.Print(green("=")); time.Sleep(20 * time.Millisecond) }
	fmt.Println("]")

	sig := sign(pk, "MINT_ACTION")
	res := postRequest("/mine", map[string]string{"userAddress": addr, "signature": sig})
	var m MineResponse
	json.Unmarshal(res, &m)

	if m.Success {
		d := m.Data
		fmt.Println(green("\n>> SUCCESS: ACCESS GRANTED <<"))
		if d.IsEquipment {
			fmt.Println(magenta(">>> CRITICAL: NEW HARDWARE FOUND! <<<"))
			fmt.Printf("[%s] %s (%s)\n", d.Type, bold(d.Name), d.Stats)
		} else if d.Type == "MAT" || d.Type == "CONS" {
			fmt.Println(cyan(">>> ITEM FOUND <<<"))
			fmt.Printf("[%s] %s\n", d.Type, bold(d.Name))
		} else {
			fmt.Printf("DATA: %s (%s)\n", white(d.Name), d.Stats)
		}

		// Countdown
		client := &http.Client{Timeout: 2 * time.Second}
		resp, _ := client.Get(SERVER_URL + "/profile/" + addr)
		var p ProfileResponse
		json.NewDecoder(resp.Body).Decode(&p)
		
		fmt.Println(yellow("\nSystem Cooling Down..."))
		for i := p.Stats.Cooldown; i > 0; i-- {
			fmt.Printf("\rrecharge_cycle: [%-30s] %ds ", getBar(i, p.Stats.Cooldown), i)
			time.Sleep(1 * time.Second)
		}
		fmt.Printf("\rrecharge_cycle: [%-30s] READY\n", getBar(0, p.Stats.Cooldown))
		time.Sleep(500 * time.Millisecond)

	} else {
		if m.CooldownRemaining > 0 {
			fmt.Printf("\n%s %d seconds remaining.\n", red(">> SYSTEM OVERHEAT! <<"), m.CooldownRemaining)
		} else if m.Error == "System Overheat" {
			fmt.Println(red("\n>> SYSTEM OVERHEAT! <<"))
		} else {
			fmt.Println(red("\n>> FAILED: "), m.Error)
		}
		waitForKey()
	}
}

func openShop(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(bold("\n=== BLACK MARKET ==="))
	fmt.Println("\nID  | NAME             | EFFECT              | DURATION | COST")
	fmt.Println("-----------------------------------------------------------------")
	fmt.Printf("301 | %-16s | %-19s | %-8s | %s\n", "Script Kiddie", "Luck +1000", "1h", yellow("100 $HASH"))
	fmt.Printf("302 | %-16s | %-19s | %-8s | %s\n", "Black Hat Tool", "Luck +5000", "1h", yellow("500 $HASH"))
	fmt.Printf("303 | %-16s | %-19s | %-8s | %s\n", "State Sponsored", "Luck +20000", "1h", yellow("2000 $HASH"))
	fmt.Println("- - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - - -")
	fmt.Printf("501 | %-16s | %-19s | %-8s | %s\n", "Thermal Paste", "Reset Cooldown", "1x", yellow("150 $HASH"))
	fmt.Printf("502 | %-16s | %-19s | %-8s | %s\n", "Loot Crate", "Random Reward", "1x", yellow("500 $HASH"))
	fmt.Println("-----------------------------------------------------------------")
	
	fmt.Print("\nEnter ID to BUY (0 back): ")
	var idStr string; fmt.Scanln(&idStr)
	if idStr == "0" { return }

	fmt.Println("Transferring funds...")
	sig := sign(pk, "BUY_ACTION")
	res := postRequest("/shop/buy", map[string]string{"userAddress": addr, "signature": sig, "itemId": idStr})
	printRes(res)
	waitForKey()
}

func getBar(current, total int) string {
	if total == 0 { return "||||||||||||||||||||||||||||||" }
	width := 30
	percent := float64(current) / float64(total)
	fill := int(percent * float64(width))
	bar := ""
	for i := 0; i < width; i++ {
		if i < fill { bar += "#" } else { bar += "-" }
	}
	return bar
}

func serverRoom(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(bold("\n=== SERVER ROOM ==="))
	resp, _ := http.Get(SERVER_URL + "/rig/status/" + addr)
	var r StakeRes
	json.NewDecoder(resp.Body).Decode(&r)
	
	fmt.Println("ID  | HARDWARE        | YIELD/s | STAKED | PENDING")
	fmt.Println("-------------------------------------------------------")
	for _, s := range r.Stakes {
		fmt.Printf("%-3d | %-15s | %-7s | %-6d | %s\n", s.ID, s.Name, s.Yield, s.Staked, yellow(s.Pending+" $HASH"))
	}
	fmt.Println("-------------------------------------------------------")

	fmt.Println("\n[A] Install Hardware (Stake)")
	fmt.Println("[B] Remove Hardware (Unstake)")
	fmt.Println("[C] Collect Yields")
	fmt.Print("Choice (0 back): "); var c string; fmt.Scanln(&c)

	if c=="0" { return }

	if c=="A" || c=="a" {
		fmt.Print("Item ID: "); var id string; fmt.Scanln(&id)
		fmt.Print("Amount: "); var amt string; fmt.Scanln(&amt)
		amtInt, _ := strconv.Atoi(amt)
		sig := sign(pk, "STAKE_ACTION")
		printRes(postRequest("/rig/stake", map[string]interface{}{"userAddress":addr, "signature":sig, "itemId":id, "amount":amtInt}))
	} else if c=="B" || c=="b" {
		fmt.Print("Item ID: "); var id string; fmt.Scanln(&id)
		fmt.Print("Amount: "); var amt string; fmt.Scanln(&amt)
		amtInt, _ := strconv.Atoi(amt)
		sig := sign(pk, "UNSTAKE_ACTION")
		printRes(postRequest("/rig/unstake", map[string]interface{}{"userAddress":addr, "signature":sig, "itemId":id, "amount":amtInt}))
	} else if c=="C" || c=="c" {
		fmt.Print("Item ID to Claim: "); var id string; fmt.Scanln(&id)
		sig := sign(pk, "CLAIM_YIELD")
		printRes(postRequest("/rig/claim", map[string]interface{}{"userAddress":addr, "signature":sig, "itemId":id}))
	}
	waitForKey()
}

func shadowNet(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(cyan("\n[A] Daily Contracts"))
	fmt.Println(cyan("[B] Global Leaderboard"))
	fmt.Print("Choice (0 back): "); var c string; fmt.Scanln(&c)
	if c=="0" { return }
	
	if c == "A" || c == "a" {
		resp, _ := http.Get(SERVER_URL + "/contracts/" + addr)
		var r ContractsRes
		json.NewDecoder(resp.Body).Decode(&r)

		fmt.Println(bold("\n=== DAILY CONTRACTS ==="))
		
		now := time.Now().UTC()
		nextMidnight := time.Date(now.Year(), now.Month(), now.Day()+1, 0, 0, 0, 0, time.UTC)
		timeLeft := nextMidnight.Sub(now)
		h := int(timeLeft.Hours())
		m := int(timeLeft.Minutes()) % 60
		s := int(timeLeft.Seconds()) % 60
		fmt.Printf("REFRESH IN: %s\n\n", color.HiYellowString(fmt.Sprintf("%02dh %02dm %02ds", h, m, s)))

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
		waitForKey()
	}
}

func openWorkshop(pk *ecdsa.PrivateKey, addr string) {
	fmt.Println(magenta("\n=== WORKSHOP ==="))
	resp, _ := http.Get(SERVER_URL + "/inventory/" + addr)
	defer resp.Body.Close() 
	type InvRes struct { Items []InventoryItem `json:"items"` }
	var inv InvRes
	json.NewDecoder(resp.Body).Decode(&inv)
	fmt.Println("\n[INVENTORY - HARDWARE & MATERIALS]")
	found := false
	for _, it := range inv.Items {
		if it.ID >= 100 {
			fmt.Printf("[%d] %s (%s) x%d\n", it.ID, white(it.Name), it.Type, it.Qty)
			found = true
		}
	}
	if !found { fmt.Println(yellow("No upgradeable hardware found.")) }

	fmt.Print("\nTarget ID (0 back): "); var tId string; fmt.Scanln(&tId)
	if tId == "0" { return }
	
	fmt.Print("Material ID: "); var mId string; fmt.Scanln(&mId)
	sig := sign(pk, "ENCHANT_ACTION")
	res := postRequest("/workshop/enchant", map[string]string{"userAddress": addr, "signature": sig, "targetId": tId, "materialId": mId})
	printRes(res)
	waitForKey()
}

func manageRig(pk *ecdsa.PrivateKey, addr string) {
	fmt.Print("\nEquip ID (0 cancel): "); var id string; fmt.Scanln(&id)
	if id != "0" {
		sig := sign(pk, "EQUIP_ACTION")
		printRes(postRequest("/equip", map[string]string{"userAddress": addr, "signature": sig, "itemId": id}))
	}
	waitForKey()
}

func manageInventory(pk *ecdsa.PrivateKey, addr string) {
	resp, _ := http.Get(SERVER_URL + "/inventory/" + addr)
	defer resp.Body.Close()
	type InvRes struct { Items []InventoryItem `json:"items"` }
	var inv InvRes
	json.NewDecoder(resp.Body).Decode(&inv)
	
	for _, it := range inv.Items { 
		fmt.Printf("[%d] %s (%s) x%d\n", it.ID, it.Name, it.Type, it.Qty) 
	}

	fmt.Println("\n[S] Salvage  [U] Use Item  [0] Back")
	var c string; fmt.Scanln(&c)
	if c=="0" { return }
	
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
	waitForKey()
}

// --- UTILS ---
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
func printBanner() { fmt.Println(cyan(` :: RNG MINER v7.5 :: FINAL CLI :: `)) }