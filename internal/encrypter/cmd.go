package encrypter

import (
	"github.com/spf13/cobra"
)

var CmdEncrypterKeyGenerate = &cobra.Command{
	Use:   "generate",
	Short: "Generate new public/private key",
	Long:  `Generate both a public and private RSA key and dump it to stdout.`,
	Run: func(cmd *cobra.Command, args []string) {
		keysize, _ := cmd.Flags().GetInt("keysize")
		GenerateKey(keysize)
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

		key, err := LoadPrivateKeyFile(privateKeyFile)
		if err != nil {
			panic(err)
		}

		e := New(nil, key)
		if err = e.DecryptFile(srcFile, destFile); err != nil {
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
	CmdEncrypterKey.AddCommand(CmdEncrypterKeyGenerate)

	CmdEncrypterDecrypt.Flags().String("privateKey", "", "Private key file")
	CmdEncrypterDecrypt.Flags().String("srcFile", "", "file to decrypt")
	CmdEncrypterDecrypt.Flags().String("destFile", "", "file to write to")
	CmdEncrypterDecrypt.MarkFlagRequired("privateKey")
	CmdEncrypterDecrypt.MarkFlagRequired("srcFile")
	CmdEncrypterDecrypt.MarkFlagRequired("destFile")
	CmdEncrypter.AddCommand(CmdEncrypterDecrypt)

	CmdEncrypter.AddCommand(CmdEncrypterKey)
}
