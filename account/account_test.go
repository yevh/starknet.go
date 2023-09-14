package account_test

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"math/big"
	"os"
	"testing"

	"github.com/NethermindEth/juno/core/felt"
	starknetgo "github.com/NethermindEth/starknet.go"
	"github.com/joho/godotenv"

	"github.com/NethermindEth/starknet.go/account"
	"github.com/NethermindEth/starknet.go/mocks"
	"github.com/NethermindEth/starknet.go/rpc"
	"github.com/NethermindEth/starknet.go/test"
	"github.com/NethermindEth/starknet.go/types"
	"github.com/NethermindEth/starknet.go/utils"
	"github.com/golang/mock/gomock"
	"github.com/test-go/testify/require"
)

var (
	// set the environment for the test, default: mock
	testEnv = "mock"
	base    = ""
)

// TestMain is used to trigger the tests and, in that case, check for the environment to use.
func TestMain(m *testing.M) {
	flag.StringVar(&testEnv, "env", "mock", "set the test environment")
	flag.Parse()
	godotenv.Load(fmt.Sprintf(".env.%s", testEnv), ".env")
	base = os.Getenv("INTEGRATION_BASE")
	if base == "" && testEnv != "mock" {
		panic(fmt.Sprint("Failed to set INTEGRATION_BASE for ", testEnv))
	}
	os.Exit(m.Run())
}

func TestTransactionHash(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	// https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8
	t.Run("Transaction hash mock", func(t *testing.T) {
		if testEnv != "mock" {
			t.Skip("Skipping test as it requires a mock environment")
		}
		expectedHash := utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8")
		acntAddress := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")
		privKey := utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
		privKeyBI, ok := new(big.Int).SetString(privKey.String(), 0)
		require.True(t, ok)
		ks := starknetgo.NewMemKeystore()
		ks.Put(acntAddress.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, acntAddress, acntAddress.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		call := rpc.FunctionCall{
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
				utils.TestHexToFelt(t, "0x0"),
				utils.TestHexToFelt(t, "0x3"),
				utils.TestHexToFelt(t, "0x3"),
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x0"),
			},
		}
		txDetails := rpc.TxDetails{
			Nonce:   utils.TestHexToFelt(t, "0x2"),
			MaxFee:  utils.TestHexToFelt(t, "0x574fbde6000"),
			Version: rpc.TransactionV1,
		}
		hash, err := account.TransactionHash2(call.Calldata, txDetails.Nonce, txDetails.MaxFee, account.AccountAddress)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, expectedHash.String(), hash.String(), "transaction hash does not match expected")
	})

	t.Run("Transaction hash testnet", func(t *testing.T) {
		if testEnv != "testnet" {
			t.Skip("Skipping test as it requires a testnet environment")
		}
		expectedHash := utils.TestHexToFelt(t, "0x135c34f53f8b7f59efd450eb689fccd9dd4cfe7f9d9dc4d09954c5653138698")
		address := &felt.Zero

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err, "error returned from account.NewAccount()")

		call := rpc.FunctionCall{
			ContractAddress:    &felt.Zero,
			EntryPointSelector: &felt.Zero,
			Calldata:           []*felt.Felt{&felt.Zero},
		}
		txDetails := rpc.TxDetails{
			Nonce:  &felt.Zero,
			MaxFee: &felt.Zero,
		}
		hash, err := account.TransactionHash2(call.Calldata, txDetails.Nonce, txDetails.MaxFee, account.AccountAddress)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, hash.String(), expectedHash.String(), "transaction hash does not match expected")
	})

	t.Run("Transaction hash mainnet", func(t *testing.T) {
		if testEnv != "mainnet" {
			t.Skip("Skipping test as it requires a mainnet environment")
		}
		expectedHash := utils.TestHexToFelt(t, "0x3476c76a81522fe52616c41e95d062f5c3ea4eeb6c652904ad389fcd9ff4637")
		accountAddress := utils.TestHexToFelt(t, "0x59cd166e363be0a921e42dd5cfca0049aedcf2093a707ef90b5c6e46d4555a8")

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_MAIN", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, accountAddress, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err, "error returned from account.NewAccount()")

		call := rpc.FunctionCall{
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x5dbdedc203e92749e2e746e2d40a768d966bd243df04a6b712e222bc040a9af"),
				utils.TestHexToFelt(t, "0x2f0b3c5710379609eb5495f1ecd348cb28167711b73609fe565a72734550354"),
				utils.TestHexToFelt(t, "0x0"),
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x1"),
				utils.TestHexToFelt(t, "0x52884ee3f"),
			},
		}
		txDetails := rpc.TxDetails{
			Nonce:   utils.TestHexToFelt(t, "0x1"),
			MaxFee:  utils.TestHexToFelt(t, "0x2a173cd36e400"),
			Version: rpc.TransactionV1,
		}
		hash, err := account.TransactionHash2(call.Calldata, txDetails.Nonce, txDetails.MaxFee, account.AccountAddress)
		require.NoError(t, err, "error returned from account.TransactionHash()")
		require.Equal(t, expectedHash.String(), hash.String(), "transaction hash does not match expected")
	})
}

func TestFmtCallData(t *testing.T) {

	t.Run("ChainId mainnet - mock", func(t *testing.T) {

		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("transfer"),
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x1"),
			},
		}
		expectedCallData := []*felt.Felt{
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			utils.TestHexToFelt(t, "0x83afd3f4caedc6eebf44246fe54e38c95e3179a5ec9ea81740eca5b482d12e"),
			utils.TestHexToFelt(t, "0x0"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x3"),
			utils.TestHexToFelt(t, "0x49d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			utils.TestHexToFelt(t, "0x1"),
			utils.TestHexToFelt(t, "0x0"),
		}
		fmt.Println("fnCall.asd", fnCall.EntryPointSelector)
		fmtCallData := account.FmtCalldata2([]rpc.FunctionCall{fnCall})
		fmt.Println("fmtCallData", fmtCallData)
		require.Equal(t, fmtCallData, expectedCallData)
	})
}

func TestChainId(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	t.Run("ChainId mainnet - mock", func(t *testing.T) {
		mainnetID := utils.TestHexToFelt(t, "0x534e5f4d41494e")
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_MAIN", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), mainnetID.String())
	})

	t.Run("ChainId testnet - mock", func(t *testing.T) {
		testnetID := utils.TestHexToFelt(t, "0x534e5f474f45524c49")
		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, &felt.Zero, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err)
		require.Equal(t, account.ChainId.String(), testnetID.String())
	})

	t.Run("ChainId devnet", func(t *testing.T) {
		if testEnv != "devnet" {
			t.Skip("Skipping test as it requires a devnet environment")
		}
		devNetURL := "http://0.0.0.0:5050/rpc"

		fmt.Println("devNetURL", devNetURL)
		client, err := rpc.NewClient(devNetURL)
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		_, err = account.NewAccount(provider, 1, &felt.Zero, "pubkey", starknetgo.NewMemKeystore())
		require.NoError(t, err)
	})
}

func TestSign(t *testing.T) {
	mockCtrl := gomock.NewController(t)
	t.Cleanup(mockCtrl.Finish)
	mockRpcProvider := mocks.NewMockRpcProvider(mockCtrl)

	// Accepted on testnet https://goerli.voyager.online/tx/0x2a7eec54aab835323a810e893354368a496f1a217e8b6ef295476568ef08f0d
	t.Run("Sign testnet - mock", func(t *testing.T) {
		expectedS1 := utils.TestHexToFelt(t, "0x6bf7980d98fa300ed9565b8cd5efcf5582133daa961b5e1d9477bf1bd750727")
		expectedS2 := utils.TestHexToFelt(t, "0x5886b8236b7dc3665c0014876a644ddd0800a167ff1036fb82af1a6f4134c91")

		address := utils.TestHexToFelt(t, "0x476466998f22e0b0177ddc76afcf8e3b5d30164f3eb33031aae7a9cb63c831")
		// pubKey := utils.TestHexToFelt(t, "0xc1e5fc3f93e04dac29878e4efccd81a96547884628d83b40fc9b758b5a349")
		privKey := utils.TestHexToFelt(t, "0x15d0b81e6140f4cce02b47609879a723f9f5b7b9f3ffca346018c73fe81847e")
		privKeyBI, ok := new(big.Int).SetString(privKey.String(), 0)
		require.True(t, ok)
		ks := starknetgo.NewMemKeystore()
		ks.Put(address.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, address.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		msg := utils.TestHexToFelt(t, "0x2a7eec54aab835323a810e893354368a496f1a217e8b6ef295476568ef08f0d")
		sig, err := account.Sign(context.Background(), msg)

		require.NoError(t, err, "error returned from account.Sign()")
		require.Equal(t, expectedS1.String(), sig[0].String(), "s1 does not match expected")
		require.Equal(t, expectedS2.String(), sig[1].String(), "s2 does not match expected")
	})

	// Accepted on testnet https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8
	t.Run("Sign testnet - mock 2", func(t *testing.T) {
		expectedS1 := utils.TestHexToFelt(t, "0x10d405427040655f118bc8b897e2f2f8147858bbcb0e3d6bc6dfbc6d0205e8")
		expectedS2 := utils.TestHexToFelt(t, "0x5cdfe4a3d5b63002e9011ec0ba59ae2b75a43cb2a3bc1699b35aa64cb9ca3cf")

		address := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")
		privKey := utils.TestHexToFelt(t, "0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
		privKeyBI, ok := new(big.Int).SetString(privKey.String(), 0)
		require.True(t, ok)
		ks := starknetgo.NewMemKeystore()
		ks.Put(address.String(), privKeyBI)

		mockRpcProvider.EXPECT().ChainID(context.Background()).Return("SN_GOERLI", nil)
		account, err := account.NewAccount(mockRpcProvider, 1, address, address.String(), ks)
		require.NoError(t, err, "error returned from account.NewAccount()")

		msg := utils.TestHexToFelt(t, "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8")
		sig, err := account.Sign(context.Background(), msg)

		require.NoError(t, err, "error returned from account.Sign()")
		require.Equal(t, expectedS1.String(), sig[0].String(), "s1 does not match expected")
		require.Equal(t, expectedS2.String(), sig[1].String(), "s2 does not match expected")
	})
}

func TestAddInvoke(t *testing.T) {

	// https://goerli.voyager.online/tx/0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8#overview
	t.Run("Test AddInvokeTransction testnet", func(t *testing.T) {
		if testEnv != "testnet" {
			t.Skip("Skipping test as it requires a testnet environment")
		}
		// Why does deploying an account work (which reqs a sig), but this doesnt?
		// Should use another library and force it to print data out to see what's going wrong here
		// New Client
		fmt.Println("base", base)

		fmt.Println("base", new(felt.Felt).SetBytes([]byte(account.TRANSACTION_PREFIX)))
		fmt.Println("base", types.UTF8StrToBig(starknetgo.TRANSACTION_PREFIX))
		bigg := types.UTF8StrToBig(starknetgo.TRANSACTION_PREFIX)
		fellt := new(felt.Felt).SetBytes(bigg.Bytes())
		fmt.Println("base", fellt)
		client, err := rpc.NewClient(base + "/rpc")
		require.NoError(t, err, "Error in rpc.NewClient")
		provider := rpc.NewProvider(client)

		// account address
		accountAddress := utils.TestHexToFelt(t, "0x043784df59268c02b716e20bf77797bd96c68c2f100b2a634e448c35e3ad363e")

		// Set up ks
		ks := starknetgo.NewMemKeystore()
		fakePubKey, _ := new(felt.Felt).SetString("0x049f060d2dffd3bf6f2c103b710baf519530df44529045f92c3903097e8d861f")
		fakePrivKey, _ := new(big.Int).SetString("0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa", 0)
		fakePrivKeyFelt, _ := new(felt.Felt).SetString("0x043b7fe9d91942c98cd5fd37579bd99ec74f879c4c79d886633eecae9dad35fa")
		ks.Put(fakePubKey.String(), fakePrivKey)

		// Get account
		acnt, err := account.NewAccount(provider, 1, accountAddress, fakePubKey.String(), ks)
		require.NoError(t, err)

		// Now build the trasaction
		nonce, _ := new(felt.Felt).SetString("0x3") // should be 0x2 for the requries
		mxfee, _ := new(felt.Felt).SetString("0x574fbde6000")
		invokeTx := rpc.BroadcastedInvokeV1Transaction{
			BroadcastedTxnCommonProperties: rpc.BroadcastedTxnCommonProperties{
				Nonce:   nonce,
				MaxFee:  mxfee,
				Version: rpc.TransactionV1,
				Type:    rpc.TransactionType_Invoke,
			},
			SenderAddress: acnt.AccountAddress,
		}
		fmt.Println(" ====== ======invokeTx", invokeTx)
		fnCall := rpc.FunctionCall{
			ContractAddress:    utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
			EntryPointSelector: types.GetSelectorFromNameFelt("transfer"),
			Calldata: []*felt.Felt{
				utils.TestHexToFelt(t, "0x049d36570d4e46f48e99674bd3fcc84644ddd6b96f7c741b1562b82f9e004dc7"),
				utils.TestHexToFelt(t, "0x1"),
			},
		}
		// fmt.Println(fnCall.EntryPointSelector)
		// err = account.BuildInvokeTx(context.Background(), &invokeTx, &[]rpc.FunctionCall{fnCall})
		// require.NoError(t, err, "Error in BuildInvokeTx")

		fmtdCallData := account.FmtCalldata2([]rpc.FunctionCall{fnCall})
		invokeTx.Calldata = fmtdCallData

		/// NEED TO FORMAT THE CALL DATA
		fmt.Println("+++++++ pre TransactionHash2", invokeTx.Nonce, invokeTx.MaxFee)
		txHash, err := acnt.TransactionHash2(invokeTx.Calldata, invokeTx.Nonce, invokeTx.MaxFee, acnt.AccountAddress)
		// require.Equal(t, txHash.String(), "0x73cf79c4bfa0c7a41f473c07e1be5ac25faa7c2fdf9edcbd12c1438f40f13d8")
		x, y, err := starknetgo.Curve.SignFelt(txHash, fakePrivKeyFelt)
		if err != nil {
			panic(err)
		}
		// require.Equal(t, x.String(), "0x10d405427040655f118bc8b897e2f2f8147858bbcb0e3d6bc6dfbc6d0205e8")
		// require.Equal(t, y.String(), "0x5cdfe4a3d5b63002e9011ec0ba59ae2b75a43cb2a3bc1699b35aa64cb9ca3cf")

		fmt.Println("+++++++ post TransactionHash2", invokeTx.Nonce, invokeTx.MaxFee)
		invokeTx.Signature = []*felt.Felt{x, y}
		fmt.Println(" ====== invokeTx", invokeTx)
		// fmt.Println("sig", invokeTx.Signature)
		// fmt.Println("invokeTx.Calldata", invokeTx.Calldata)
		// fmt.Println("txHash", txHash)
		// fmt.Println(x, y)
		fmt.Println(" ====== invokeTx", invokeTx)
		qwe, _ := json.MarshalIndent(invokeTx, "", "")
		fmt.Println("MarshalIndent", string(qwe))
		fmt.Println("+++++++ post MarshalIndent", invokeTx.Nonce, invokeTx.MaxFee)
		fmt.Println(" ====== invokeTx", invokeTx)

		///
		// pub, _ := new(big.Int).SetString("2090221843434510384432085791482977629840322403554658343615172301617258923551", 0)
		// hash, _ := new(big.Int).SetString("2122438891878094235855424351251599998644607787805056148746903357402852875679", 0)
		// fmt.Println("")
		// fmt.Println("txHash", txHash) // It seems the txhash is incorrect. Signature is correct.
		// fmt.Println("acntadr", accountAddress, accountAddress.BigInt(new(big.Int)))
		// fmt.Println("pub", pub, new(felt.Felt).SetBytes(pub.Bytes()))
		// fmt.Println("hash", hash, new(felt.Felt).SetBytes(hash.Bytes()))
		// fmt.Println("")
		// fmt.Println(invokeTx.Signature[0].BigInt(new(big.Int)))
		// fmt.Println(invokeTx.Signature[1].BigInt(new(big.Int)))
		// fmt.Println(invokeTx.Signature[0].BigInt(new(big.Int)))
		// fmt.Println(invokeTx.Signature[1].BigInt(new(big.Int)))
		// // Send tx
		resp, err := acnt.AddInvokeTransaction(context.Background(), &invokeTx)
		fmt.Println("resp", resp, err)
		require.NoError(t, err)
	})
}

func newDevnet(t *testing.T, url string) ([]test.TestAccount, error) {
	// url := SetupLocalStarknetNode(t)
	devnet := test.NewDevNet(url)
	acnts, err := devnet.Accounts()
	return acnts, err
}