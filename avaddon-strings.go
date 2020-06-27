package main

import (
	"bufio"
	"encoding/base64"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/jessevdk/go-flags"
	. "github.com/logrusorgru/aurora"
)

// check errors as they occur and panic :o
func check(e error) {
	if e != nil {
		panic(e)
	}
}

// read specified file containing encrypted strings
func readInputFile(path string) ([]string, error) {
	var lines []string
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	// read input file line by line with bufio Scanner
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}

func hexStr(hexaString string) string {
	// replace 0x or 0X with empty String
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

func decrypt(ciphertext string, sub string, xor string) (cleartext string) {

	decodedCipher, base64Err := base64.StdEncoding.DecodeString(ciphertext)
	check(base64Err)
	byteSlice := []byte(decodedCipher)

	subKey, convErr := strconv.ParseUint(hexStr(sub), 16, 64)
	xorKey, convErr := strconv.ParseUint(hexStr(xor), 16, 64)
	check(convErr)

	counter := 0
	for counter < len(byteSlice) {
		byteSlice[counter] = (byteSlice[counter] - byte(subKey)) ^ byte(xorKey)
		counter++
	}

	return string(byteSlice)
}

func main() {

	var opts struct {
		InFlag  string `short:"i" long:"in" description:"path to the file containing the encrypted strings" required:"true"`
		OutFlag string `short:"o" long:"out" description:"optional: path/filename of the output file. If none is specified the file will be named 'decrypted_strings.txt"`
	}

	_, flagErr := flags.Parse(&opts) // parse the command arguments/flags given above
	check(flagErr)

	fmt.Printf(Sprintf(Red("\n       ___                  __    __                 _____ __       _ \n")))
	fmt.Printf(Sprintf(Red("      /   |_   ______ _____/ /___/ /___  ____       / ___// /______(_)___  ____ ______\n")))
	fmt.Printf(Sprintf(Red("     / /| | | / / __ `/ __  / __  / __ \\/ __ \\______\\__ \\/ __/ ___/ / __ \\/ __ `/ ___/\n")))
	fmt.Printf(Sprintf(Red("    / ___ | |/ / /_/ / /_/ / /_/ / /_/ / / / /_____/__/ / /_/ /  / / / / / /_/ (__  ) \n")))
	fmt.Printf(Sprintf(Red("   /_/  |_|___/\\__,_/\\__,_/\\__,_/\\____/_/ /_/     /____/\\__/_/  /_/_/ /_/\\__, /____/  \n")))
	fmt.Printf(Sprintf(Red("                                                                        /____/\n")))
	fmt.Printf(Sprintf(White("\n              Avaddon Ransomware String Decrypter - BASE64 > SUB > XOR\n")))
	fmt.Printf(Sprintf(White("                Marius 'f0wL' Genheimer | https://dissectingmalwa.re\n\n\n")))

	encryptedLines, readErr := readInputFile(opts.InFlag)
	check(readErr)
	if opts.OutFlag == "" {
		opts.OutFlag = "decrypted_strings.txt"
	}

	fmt.Printf("   [>] Please enter a value (e.g. 0x2) for the SUB operation: ")
	var subVal string
	_, numError := fmt.Scanf("%s", &subVal)
	check(numError)

	fmt.Printf("\n   [>] Please enter a value (e.g. 0x43) for the XOR operation: ")
	var xorVal string
	_, xorError := fmt.Scanf("%s", &xorVal)
	check(xorError)
	print("\n   [!] Decrypted strings: \n\n")

	f, createErr := os.Create(opts.OutFlag)
	check(createErr)

	counter := 0
	for counter < len(encryptedLines) {
		enc := encryptedLines[counter]
		dec := decrypt(enc, subVal, xorVal)
		fmt.Println("   " + enc)
		fmt.Printf(Sprintf(Green("   " + dec + "\n\n")))

		_, writeErr := f.WriteString(dec + "\n")
		check(writeErr)

		counter++
	}
	closeErr := f.Close()
	check(closeErr)

	print("\n   [!] Wrote decrypted strings to " + opts.OutFlag + "\n\n")

}