package sdk_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/mises-id/sdk"
	"github.com/mises-id/sdk/bip39"
	"github.com/mises-id/sdk/misesid"
	"github.com/mises-id/sdk/types"
	"github.com/mises-id/sdk/user"
	tmtypes "github.com/tendermint/tendermint/types"
	"github.com/tyler-smith/assert"
)

const (
	SimplePassWord = "12345678"
)

func CreateRandomUser() types.MUser {
	//create user
	entropy, _ := bip39.NewEntropy(128)

	mnemonics, _ := bip39.NewMnemonic(entropy)

	var ugr user.MisesUserMgr
	pUgr := &ugr
	cuser, _ := pUgr.CreateUser(mnemonics, SimplePassWord)

	return cuser
}
func TestMain(m *testing.M) {
	// call flag.Parse() here if TestMain uses flags
	types.NodeHome = "test.misestm"
	os.Exit(m.Run())
}

func TestSdkNewForUesr(t *testing.T) {
	misesid.DeleteKeyStoreFile()
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	// test NewSdkForUser
	s := sdk.NewSdkForUser(mo)

	ugr := s.UserMgr()

	// test CreateUser
	mnemonics, err := misesid.RandomMnemonics()
	assert.NoError(t, err)
	fmt.Printf("mnemonics is: %s\n", mnemonics)

	newUser, err := ugr.CreateUser(mnemonics, SimplePassWord)
	assert.NoError(t, err)

	// fmt.Printf("keystore version is: %d\n", misesid.Ks.Version)
	// fmt.Printf("keystore id is: %s\n", misesid.Ks.MId)
	// fmt.Printf("keystore address is: %s\n", misesid.Ks.PubKey)
	// fmt.Printf("keystore kdf is: %s\n", misesid.Ks.Crypto.Kdf)
	// fmt.Printf("keystore dklen is: %d\n", misesid.Ks.Crypto.KdfParams.Dklen)
	// fmt.Printf("keystore salt is: %s\n", misesid.Ks.Crypto.KdfParams.Salt)
	// fmt.Printf("keystore n is: %d\n", misesid.Ks.Crypto.KdfParams.N)
	// fmt.Printf("keystore r is: %d\n", misesid.Ks.Crypto.KdfParams.R)
	// fmt.Printf("keystore p is: %d\n", misesid.Ks.Crypto.KdfParams.P)
	// fmt.Printf("keystore cipher is: %s\n", misesid.Ks.Crypto.Cipher)
	// fmt.Printf("keystore ciphertext is: %s\n", misesid.Ks.Crypto.Ciphertext)
	// fmt.Printf("keystore iv is: %s\n", misesid.Ks.Crypto.CipherParams.Iv)
	// fmt.Printf("keystore mac is: %s\n", misesid.Ks.Crypto.Mac)

	u := ugr.ActiveUser()
	assert.True(t, u != nil)
	assert.True(t, newUser.MisesID() == u.MisesID())
	fmt.Printf("user's mid is %s\n", u.MisesID())
	fmt.Printf("user's pubKey is %s\n", u.Signer().PubKey())

	// test Login, sign & verify
	permissions := []string{"user_info_r", "user_info_w"}
	auth, err := s.Login("mises.site", permissions)
	assert.NoError(t, err)

	fmt.Printf("auth string is: %s\n", auth)

	// v, err := url.ParseQuery(auth)
	// assert.NoError(t, err)
	// misesID := v.Get("mises_id")
	// sigStr := v.Get("sig")
	// nonce := v.Get("nonce")

	// err = misesid.Verify(misesID+"&"+nonce, u.Signer().PubKey(), sigStr)
	// assert.NoError(t, err)
	// if err == nil {
	// 	fmt.Printf("signature is verified\n")
	// } else {
	// 	fmt.Printf("Signature verification is failed\n")
	// }

	misesid.DeleteKeyStoreFile()
	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Discover",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	sApp, _ := sdk.NewSdkForApp(mo, appinfo)

	mid, _, err := sApp.VerifyLogin(auth)
	assert.NoError(t, err)
	assert.True(t, mid == u.MisesID())

}

func TestSdkVerifyLogin(t *testing.T) {
	misesid.DeleteKeyStoreFile()
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Discover",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	sApp, _ := sdk.NewSdkForApp(mo, appinfo)
	auth := "mises_id=did:mises:mises1y53kz80x5gm2w0ype8x7a3w6sstztxxg7qkl5n&nonce=0123456789&sig=304402201ada63a9dccc8ace5b3c96b00817311a36096c997e081b57f8b39b2392a51905022041e74283ec05333062a3a7180ba2775b5e203e596c3cefd8b92b775b519b7e06&pubkey=03e78b0e4bddddabd37bca173c9df270096ec55aa97bed2ba82d72c830d400c8e5"

	mid, _, err := sApp.VerifyLogin(auth)
	assert.NoError(t, err)
	assert.True(t, mid == "did:mises:mises1y53kz80x5gm2w0ype8x7a3w6sstztxxg7qkl5n")

}

func TestSdkActiveUesr(t *testing.T) {
	misesid.DeleteKeyStoreFile()
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	// test NewSdkForUser
	s := sdk.NewSdkForUser(mo)

	ugr := s.UserMgr()

	u := ugr.ActiveUser()
	assert.True(t, u == nil)
	mnemonics, err := misesid.RandomMnemonics()
	assert.NoError(t, err)
	fmt.Printf("mnemonics is: %s\n", mnemonics)

	newUser, err := ugr.CreateUser(mnemonics, SimplePassWord)
	assert.NoError(t, err)

	u = ugr.ActiveUser()
	assert.True(t, u.MisesID() == newUser.MisesID())

	s = sdk.NewSdkForUser(mo)
	ugr = s.UserMgr()
	u = ugr.ActiveUser()
	assert.NotNil(t, u)

	mo1 := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: "",
	}

	s = sdk.NewSdkForUser(mo1)
	ugr = s.UserMgr()
	u = ugr.ActiveUser()
	assert.Nil(t, u)

}

func TestSdkRegisterUser(t *testing.T) {
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Discover",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(mo, appinfo)

	newUser := CreateRandomUser()

	err := app.RunSync(app.NewRegisterUserCmd(newUser.MisesID(), newUser.Signer().PubKey(), 100000))
	assert.NoError(t, err)

}

type RegisterUserCallback struct {
	done         chan bool
	successCount int
	failCount    int
	maxCount     int
}

func (cb *RegisterUserCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	fmt.Printf("OnTxGenerated  %s %s\n", cmd.MisesUID(), cmd.TxID())
}
func (cb *RegisterUserCallback) OnSucceed(cmd types.MisesAppCmd) {
	fmt.Printf("OnSucceed %d %s %s\n", cb.successCount, cmd.MisesUID(), cmd.TxID())
	cb.successCount += 1
	if cb.successCount+cb.failCount == cb.maxCount {
		cb.done <- true
	}
}
func (cb *RegisterUserCallback) OnFailed(cmd types.MisesAppCmd, err error) {
	fmt.Printf("OnFailed %s, %s\n", cmd.TxID(), err)
	cb.failCount += 1
	if cb.successCount+cb.failCount == cb.maxCount {
		cb.done <- true
	}
}
func (cb *RegisterUserCallback) wait() {
	<-cb.done
}

func TestHexDuplication(t *testing.T) {
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Discover",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(mo, appinfo)

	callback := &RegisterUserCallback{}
	callback.done = make(chan bool)
	callback.maxCount = 1
	app.SetListener(callback)

	err := app.RunAsync(app.NewRegisterUserCmd("did:mises:mises1jprvqruscc6unxdyan9738s3qcxzgsykhx97me", "032e1cc6a8396a407c0c2c9586c16531f09ee727750db919bf116b28146f1c7454", 1000), false)
	if err != nil {
		fmt.Printf(err.Error())
		return
	}
	callback.wait()

}

func BenchmarkSdkRegisterUserFlooding(t *testing.B) {
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Discover",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(mo, appinfo)

	callback := &RegisterUserCallback{}
	callback.done = make(chan bool)
	callback.maxCount = 10000
	app.SetListener(callback)
	for userIndex := 0; userIndex < callback.maxCount; userIndex++ {
		newUser := CreateRandomUser()

		err := app.RunAsync(app.NewRegisterUserCmd(newUser.MisesID(), newUser.Signer().PubKey(), 1000), false)
		if err != nil {
			fmt.Printf(err.Error())
			return
		}
	}
	callback.wait()

}

type FaucetCallback struct {
}

func (cb *FaucetCallback) OnTxGenerated(cmd types.MisesAppCmd) {
	fmt.Printf("OnTxGenerated\n")
}
func (cb *FaucetCallback) OnSucceed(cmd types.MisesAppCmd) {
	fmt.Printf("OnSucceed\n")
}
func (cb *FaucetCallback) OnFailed(cmd types.MisesAppCmd, err error) {
	fmt.Printf("OnFailed %s\n", err)
}

func TestSdkFaucet(t *testing.T) {
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
	}

	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Faucet",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(mo, appinfo)

	app.SetListener(&FaucetCallback{})
	newUser := CreateRandomUser()

	err := app.RunSync(app.NewFaucetCmd(newUser.MisesID(), newUser.Signer().PubKey(), 100))
	assert.NoError(t, err)

}

type EventStreamingCallback struct {
	done       chan bool
	eventCount int
	maxCount   int
	header     *tmtypes.EventDataNewBlockHeader
	tx         *tmtypes.EventDataTx
}

func (cb *EventStreamingCallback) OnTxEvent(t *tmtypes.EventDataTx) {
	fmt.Printf("OnTxEvent\n")
	cb.eventCount++
	cb.tx = t
	if cb.eventCount > cb.maxCount || (cb.tx != nil && cb.header != nil) {
		cb.done <- true
	}
}
func (cb *EventStreamingCallback) OnNewBlockHeaderEvent(h *tmtypes.EventDataNewBlockHeader) {
	fmt.Printf("OnNewBlockHeaderEvent\n")
	cb.eventCount++
	cb.header = h
	if cb.eventCount > cb.maxCount || (cb.tx != nil && cb.header != nil) {
		cb.done <- true
	}
}
func (cb *EventStreamingCallback) OnEventStreamingTerminated() {
	fmt.Printf("OnEventStreamingTerminated")
	cb.done <- true
}
func (cb *EventStreamingCallback) wait() {
	<-cb.done
}

func TestSdkEventStreaming(t *testing.T) {
	mo := types.MSdkOption{
		ChainID:    "test",
		Debug:      true,
		PassPhrase: SimplePassWord,
		RpcURI:     "http://mises.ihuaj.com:26657",
	}

	appinfo := types.NewMisesAppInfoReadonly(
		"Mises Faucet",
		"https://www.mises.site",
		"https://home.mises.site",
		[]string{"mises.site"},
		"Mises Network",
	)
	_, app := sdk.NewSdkForApp(mo, appinfo)

	callback := &EventStreamingCallback{}
	callback.done = make(chan bool)
	callback.maxCount = 10

	err := app.StartEventStreaming(callback)
	assert.NoError(t, err)
	newUser := CreateRandomUser()

	err = app.RunAsync(app.NewRegisterUserCmd(newUser.MisesID(), newUser.Signer().PubKey(), 1000), false)
	assert.NoError(t, err)
	callback.wait()

	resp, err := app.ParseEvent(callback.header, callback.tx)
	assert.NoError(t, err)
	fmt.Printf("ParseEvent %s", resp.TxHash)

}
