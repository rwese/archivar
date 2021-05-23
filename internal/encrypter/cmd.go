package encrypter

import (
	"crypto/rsa"
	"log"
	"os"

	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
)

var CmdEncrypterKeyGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generate new public/private key",
	Long:  `Generate both a public and private RSA key and dump it to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		keysize, _ := cmd.Flags().GetInt("keysize")
		getPassphraseInput, _ := cmd.Flags().GetBool("with-passphrase")

		if getPassphraseInput {
			passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}
			GenerateKeyWithPassphrase(keysize, passphrase)
		} else {
			GenerateKey(keysize)
		}
	},
}

var CmdEncrypterKey = &cobra.Command{
	Use:   "key",
	Short: "key operations",
	Long:  `key operations to create, read, dump`,
	Run: func(cmd *cobra.Command, args []string) {

	},
}

var CmdEncrypterDecrypt = &cobra.Command{
	Use:   "decrypt",
	Short: "decrypt file",
	Long:  `decrypt a file with the given privateKey in pem format`,
	Run: func(cmd *cobra.Command, args []string) {
		privateKeyFile, _ := cmd.Flags().GetString("privateKey")
		srcFile, _ := cmd.Flags().GetString("srcFile")
		destFile, _ := cmd.Flags().GetString("destFile")
		getPassphraseInput, _ := cmd.Flags().GetBool("with-passphrase")

		var (
			privateKey *rsa.PrivateKey
			err        error
		)

		if getPassphraseInput {
			passphrase, err := terminal.ReadPassword(int(os.Stdin.Fd()))
			if err != nil {
				log.Fatal(err)
			}

			privateKey, err = LoadEncryptedPrivateKeyFile(privateKeyFile, passphrase)
			if err != nil {
				panic(err)
			}
		} else {
			privateKey, err = LoadPrivateKeyFile(privateKeyFile)
			if err != nil {
				panic(err)
			}
		}

		if err != nil {
			panic(err)
		}

		e := New(nil, privateKey)
		if err = e.DecryptFile(srcFile, destFile); err != nil {
			panic(err)
		}
	},
}
var CmdEncrypterSplit = &cobra.Command{
	Use:   "split",
	Short: "split file into rsa-Key and aes-Body",
	Long: `split will separate the rsa encrypted key from the aes encrypted
	body. This will in theory allow you to separately decrypt the key to then
	decrypt the body. Good luck.`,
	Run: func(cmd *cobra.Command, args []string) {
		srcFile, _ := cmd.Flags().GetString("srcFile")

		e := New(nil, nil)
		if err := e.SplitFile(srcFile); err != nil {
			panic(err)
		}
	},
}

var CmdEncrypter = &cobra.Command{
	Use:   "encrypter",
	Short: "encrypter",
}

func init() {
	CmdEncrypterKeyGenerate.Flags().Int("keysize", 1024, "keysize in bits")
	CmdEncrypterKeyGenerate.Flags().BoolP("with-passphrase", "p", false, "Ask for passphrase")
	CmdEncrypterKey.AddCommand(CmdEncrypterKeyGenerate)
	CmdEncrypter.AddCommand(CmdEncrypterKey)

	CmdEncrypterDecrypt.Flags().BoolP("with-passphrase", "p", false, "Ask for passphrase")
	CmdEncrypterDecrypt.Flags().String("privateKey", "", "Private key file")
	CmdEncrypterDecrypt.Flags().String("srcFile", "", "file to decrypt")
	CmdEncrypterDecrypt.Flags().String("destFile", "", "file to write to")
	CmdEncrypterDecrypt.MarkFlagRequired("privateKey")
	CmdEncrypterDecrypt.MarkFlagRequired("srcFile")
	CmdEncrypterDecrypt.MarkFlagRequired("destFile")
	CmdEncrypter.AddCommand(CmdEncrypterDecrypt)

	CmdEncrypterSplit.Flags().String("srcFile", "", "file to Split")
	CmdEncrypterSplit.MarkFlagRequired("srcFile")

	CmdEncrypter.AddCommand(CmdEncrypterSplit)

}
