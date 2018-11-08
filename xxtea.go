package tools

import (
	"github.com/xxtea/xxtea-go/xxtea"
	"fmt"
	"encoding/hex"
)

func main()  {
	//str := "Hello World! 你好，中国！"
	//key := "xK*63<qmO@456fhgHRT*"
	//encrypt_data := xxtea.Encrypt([]byte(str), []byte(key))
	//a := hex.EncodeToString(encrypt_data)
	//a = strings.ToUpper(a)
	//fmt.Println("Encrypt：",a)

	b,_ := hex.DecodeString("B5AFD5A546D7E09173E96C4298F54975")
	decrypt_data := string(xxtea.Decrypt(b, []byte(`#%&(!*65#@$^&)@_`)))
	fmt.Println("Decrypt:",decrypt_data)

}