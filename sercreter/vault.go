package secret

//vault: 金庫室

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"sync"

	"github.com/rensawamo/sercreter/cipher"
)

func File(encodingKey, filepath string) *Vault {
	return &Vault{
		encodingKey: encodingKey,
		filepath:    filepath,
	}
}

type Vault struct {
	encodingKey string
	filepath    string
	// lock unlock 用のmutex
	mutex sync.Mutex
	// set で キーとバリューを保存して get で キーをもちいて valueを
	keyValues map[string]string
}

// ファイルの存在の確認
func (v *Vault) load() error {
	fmt.Println("v", v.filepath)
	f, err := os.Open(v.filepath)
	fmt.Println(f)
	if err != nil {
		fmt.Println("新規")
		v.keyValues = make(map[string]string)
		return nil
	}
	defer f.Close()
	r, err := cipher.DecryptReader(v.encodingKey, f)
	if err != nil {
		return err
	}
	return v.readKeyValues(r)
}

func (v *Vault) readKeyValues(r io.Reader) error {
	dec := json.NewDecoder(r)
	return dec.Decode(&v.keyValues)
}

func (v *Vault) save() error {
	// os.OpenFile(file名, フラグ, ファイルモード)
	f, err := os.OpenFile(v.filepath, os.O_RDWR|os.O_CREATE, 0755)
	if err != nil {
		return err
	}
	defer f.Close()
	// 暗号化の ライターを受け渡す
	w, err := cipher.EncryptWriter(v.encodingKey, f)
	if err != nil {
		return err
	}
	// ライターをエンコードする
	return v.writeKeyValues(w)
}

// encode
func (v *Vault) writeKeyValues(w io.Writer) error {
	enc := json.NewEncoder(w)
	return enc.Encode(v.keyValues)
}

// 金庫の 構造体を定義してそこの methodをつくる
func (v *Vault) Get(key string) (string, error) {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return "", err
	}
	value, ok := v.keyValues[key]
	if !ok {
		return "", errors.New("secret: no value for that key")
	}
	return value, nil
}

func (v *Vault) Set(key, value string) error {
	v.mutex.Lock()
	defer v.mutex.Unlock()
	err := v.load()
	if err != nil {
		return err
	}
	v.keyValues[key] = value
	err = v.save()
	return err
}
