package migration

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"os"

	"github.com/nuttapon-first/omma-kebab-server/modules/model"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

func CheckAdminUser(db *gorm.DB) error {
	admin := &model.User{
		UserName:     "admin",
		UserFullName: "admin",
		UserRole:     "superAdmin",
	}
	err := db.Model(&model.User{}).Where("user_name = ?", "admin").First(&model.User{}).Error
	if err != nil {
		if err.Error() != "record not found" {
			return err
		}
		fmt.Println("Start migrate admin user .........")
		tx := db.Begin()
		defer func() {
			if r := recover(); r != nil {
				tx.Rollback()
			}
		}()

		err := tx.Model(&model.User{}).Create(admin).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		//Decrypt the text:
		decrypted, err := decrypt([]byte(os.Getenv("ENCRYPT_KEY")), os.Getenv("ENCRYPT_PASSWORD"))

		//IF the decryption failed:
		if err != nil {
			return err
		}

		//Print re-decrypted text:
		password, err := hex.DecodeString(decrypted)
		if err != nil {
			return err
		}

		// // Hashing the password with the default cost of 10
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		userCredential := &model.UserCredential{
			UserId:     admin.ID,
			Credential: string(hashedPassword),
		}
		err = tx.Model(&model.UserCredential{}).Create(userCredential).Error
		if err != nil {
			tx.Rollback()
			return err
		}

		if err := tx.Commit().Error; err != nil {
			return err
		}
		fmt.Println("Migrate admin user success.")
	}
	return nil
}

// func encrypt(key []byte, message string) (encoded string, err error) {
// 	//Create byte array from the input string
// 	plainText := []byte(message)

// 	//Create a new AES cipher using the key
// 	block, err := aes.NewCipher(key)

// 	//IF NewCipher failed, exit:
// 	if err != nil {
// 		return
// 	}

// 	//Make the cipher text a byte array of size BlockSize + the length of the message
// 	cipherText := make([]byte, aes.BlockSize+len(plainText))

// 	//iv is the ciphertext up to the blocksize (16)
// 	iv := cipherText[:aes.BlockSize]
// 	if _, err = io.ReadFull(rand.Reader, iv); err != nil {
// 		return
// 	}

// 	//Encrypt the data:
// 	stream := cipher.NewCFBEncrypter(block, iv)
// 	stream.XORKeyStream(cipherText[aes.BlockSize:], plainText)

// 	fmt.Println(cipherText)
// 	//Return string encoded in base64
// 	return base64.RawStdEncoding.EncodeToString(cipherText), err
// }

// https://gist.github.com/STGDanny/03acf29a90684c2afc9487152324e832
/*
 *	FUNCTION		: decrypt
 *	DESCRIPTION		:
 *		This function takes a string and a key and uses AES to decrypt the string into plain text
 *
 *	PARAMETERS		:
 *		byte[] key	: Byte array containing the cipher key
 *		string secure	: String containing an encrypted message
 *
 *	RETURNS			:
 *		string decoded	: String containing the decrypted equivalent of secure
 *		error err	: Error message
 */
func decrypt(key []byte, secure string) (decoded string, err error) {
	//Remove base64 encoding:
	cipherText, err := base64.RawStdEncoding.DecodeString(secure)
	if err != nil {
		return "", err
	}

	//Create a new AES cipher with the key and encrypted message
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	iv := cipherText[:aes.BlockSize]
	cipherText = cipherText[aes.BlockSize:]

	//Decrypt the message
	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherText, cipherText)

	return string(cipherText), nil
}
