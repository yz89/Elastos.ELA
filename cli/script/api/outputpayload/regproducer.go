package outputpayload

import (
	"bytes"
	"fmt"
	"os"

	"github.com/elastos/Elastos.ELA/cli/script/api/client"
	"github.com/elastos/Elastos.ELA/common"
	"github.com/elastos/Elastos.ELA/core/types/outputpayload"
	"github.com/elastos/Elastos.ELA/crypto"

	"github.com/yuin/gopher-lua"
)

const (
	luaRegisterProducerName = "registerproduceroutput"
)

// Registers my person type to given L.
func RegisterRegisterProducerOutputType(L *lua.LState) {
	mt := L.NewTypeMetatable(luaRegisterProducerName)
	L.SetGlobal("registerproduceroutput", mt)
	// static attributes
	L.SetField(mt, "new", L.NewFunction(newRegisterProducer))
	// methods
	L.SetField(mt, "__index", L.SetFuncs(L.NewTable(), registerProducerMethods))
}

// Constructor
func newRegisterProducer(L *lua.LState) int {
	publicKeyStr := L.ToString(1)
	nickName := L.ToString(2)
	url := L.ToString(3)
	location := L.ToInt64(4)
	address := L.ToString(5)
	client := client.CheckClient(L, 6)

	publicKey, err := common.HexStringToBytes(publicKeyStr)
	if err != nil {
		fmt.Println("wrong producer public key")
		os.Exit(1)
	}
	registerProducer := &outputpayload.RegisterProducer{
		PublicKey: []byte(publicKey),
		NickName:  nickName,
		Url:       url,
		Location:  uint64(location),
		Address:   address,
	}

	rpSignBuf := new(bytes.Buffer)
	err = registerProducer.SerializeUnsigned(rpSignBuf)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	acc, err := client.GetDefaultAccount()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	rpSig, err := crypto.Sign(acc.PrivKey(), rpSignBuf.Bytes())
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	registerProducer.Signature = rpSig

	ud := L.NewUserData()
	ud.Value = registerProducer
	L.SetMetatable(ud, L.GetTypeMetatable(luaRegisterProducerName))
	L.Push(ud)

	return 1
}

// Checks whether the first lua argument is a *LUserData with *PayloadRegisterProducer and
// returns this *PayloadRegisterProducer.
func checkRegisterProducer(L *lua.LState, idx int) *outputpayload.RegisterProducer {
	ud := L.CheckUserData(idx)
	if v, ok := ud.Value.(*outputpayload.RegisterProducer); ok {
		return v
	}
	L.ArgError(1, "PayloadRegisterProducer expected")
	return nil
}

var registerProducerMethods = map[string]lua.LGFunction{
	"get": registerProducerGet,
}

// Getter and setter for the Person#Name
func registerProducerGet(L *lua.LState) int {
	p := checkRegisterProducer(L, 1)
	fmt.Println(p)

	return 0
}
