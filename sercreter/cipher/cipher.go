package cipher

import (
	// aes パッケージは、米国連邦情報処理標準出版物 197 で定義されている AES 暗号化 (旧称 Rijndael) を実装します。
	"crypto/aes"
	"crypto/cipher"
	"crypto/md5"
	"crypto/rand"

	"errors"
	"fmt"
	"io"
)

func encryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBEncrypter(block, iv), nil
}

// キーと 平文をうけとり 暗号化された値の 16進数を返す（暗号化する）
// func Encrypt(key, plaintext string) (string, error) {
// 	ciphertext := make([]byte, aes.BlockSize+len(plaintext))
// 	iv := ciphertext[:aes.BlockSize]
// 	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
// 		return "", err
// 	}
// 	stream, err := encryptStream(key, iv)
// 	if err != nil {
// 		return "", err
// 	}
// 	stream.XORKeyStream(ciphertext[aes.BlockSize:], []byte(plaintext))

// 	return fmt.Sprintf("%x", ciphertext), nil
// }

// 暗号化されたデータを 元のライターに書き込む
func EncryptWriter(key string, w io.Writer) (*cipher.StreamWriter, error) {
	iv := make([]byte, aes.BlockSize)
	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return nil, err
	}
	stream, err := encryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	n, err := w.Write(iv)
	if n != len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to write full iv to writer")
	}
	return &cipher.StreamWriter{S: stream, W: w}, nil
}

// decrypt : 復号化   // blockをつかって復号化する
func decryptStream(key string, iv []byte) (cipher.Stream, error) {
	block, err := newCipherBlock(key)
	if err != nil {
		return nil, err
	}
	return cipher.NewCFBDecrypter(block, iv), nil
}

// cipherHex :  暗号16進数
// 暗号の16進数表現を受け取り それを復号
// func Decrypt(key, cipherHex string) (string, error) {
// 	ciphertext, err := hex.DecodeString(cipherHex)
// 	if err != nil {
// 		return "", err
// 	}

// 	if len(ciphertext) < aes.BlockSize {
// 		return "", errors.New("encrypt: cipher too short")
// 	}
// 	iv := ciphertext[:aes.BlockSize]
// 	ciphertext = ciphertext[aes.BlockSize:]

// 	stream, err := decryptStream(key, iv)
// 	if err != nil {
// 		return "", err
// 	}
// 	stream.XORKeyStream(ciphertext, ciphertext)
// 	return string(ciphertext), nil
// }
// decrypt : 復号
// 提供されたリーダからデータを復号化するリーダを返す
func DecryptReader(key string, r io.Reader) (*cipher.StreamReader, error) {
	iv := make([]byte, aes.BlockSize)
	n, err := r.Read(iv)
	if n < len(iv) || err != nil {
		return nil, errors.New("encrypt: unable to read the full iv")
	}
	stream, err := decryptStream(key, iv)
	if err != nil {
		return nil, err
	}
	return &cipher.StreamReader{S: stream, R: r}, nil
}


// AEA暗号  無線LAN などの通信データの暗号化に用いられる暗号化アルゴリズ
// 固定長のキーを使って データを暗号化 復号化する
func newCipherBlock(key string) (cipher.Block, error) {
	hasher := md5.New()
	// hasher に キーを書き込む  fmt.Fprint の中の writeがなかで 第一引数のポインタを使用
	fmt.Fprint(hasher, key)
	cipherKey := hasher.Sum(nil)
	return aes.NewCipher(cipherKey)
}
