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

// replace 0x or 0X with empty String
func hexStr(hexaString string) string {
	numberStr := strings.Replace(hexaString, "0x", "", -1)
	numberStr = strings.Replace(numberStr, "0X", "", -1)
	return numberStr
}

// decrypt strings with the provided values
func decrypt(ciphertext string, sub string, xor string) (cleartext string) {

	// Base64-Decode the string
	decodedCipher, base64Err := base64.StdEncoding.DecodeString(ciphertext)
	check(base64Err)
	// cast the string to a byte slice
	byteSlice := []byte(decodedCipher)

	// convert base16 values to base10 for the sake of it
	subKey, subErr := strconv.ParseUint(hexStr(sub), 16, 64)
	check(subErr)
	xorKey, xorErr := strconv.ParseUint(hexStr(xor), 16, 64)
	check(xorErr)

	counter := 0
	for counter < len(byteSlice) {
		// subtract and xor the slice char by char
		byteSlice[counter] = (byteSlice[counter] - byte(subKey)) ^ byte(xorKey)
		counter++
	}
	// return the decrypted byte slice and cast it back to a string
	return string(byteSlice)
}

func main() {

	// Commandline Arguments/ Options
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

	// read the text file
	encryptedLines, readErr := readInputFile(opts.InFlag)
	check(readErr)
	// set default output filename if none was provided
	if opts.OutFlag == "" {
		opts.OutFlag = "decrypted_strings.txt"
	}

	fmt.Printf("   [!] Visit https://github.com/f0wl/avaddon-strings for more Information\n\n")

	// read first value from stdin
	fmt.Printf("   [>] Please enter a value (e.g. 0x2) for the SUB operation: ")
	var subVal string
	_, numError := fmt.Scanf("%s", &subVal)
	check(numError)

	// read second value from stdin
	fmt.Printf("\n   [>] Please enter a value (e.g. 0x43) for the XOR operation: ")
	var xorVal string
	_, xorError := fmt.Scanf("%s", &xorVal)
	check(xorError)
	print("\n   [!] Decrypted strings: \n\n")

	// create or open output file
	f, createErr := os.Create(opts.OutFlag)
	check(createErr)

	// print encrypted and decrypted strings
	counter := 0
	for counter < len(encryptedLines) {
		enc := encryptedLines[counter]
		dec := decrypt(enc, subVal, xorVal)
		fmt.Println("   " + enc)
		fmt.Printf(Sprintf(Green("   " + dec + "\n\n")))

		// write the string to the output file
		_, writeErr := f.WriteString(dec + "\n")
		check(writeErr)

		counter++
	}

	// close the output file
	closeErr := f.Close()
	check(closeErr)

	print("\n   [!] Wrote decrypted strings to " + opts.OutFlag + "\n\n")

}
