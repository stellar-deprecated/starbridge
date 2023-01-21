//lint:file-ignore U1001 Ignore all unused code, this is only used in tests.
package integration

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	soroban_bridge "github.com/stellar/starbridge/soroban-bridge"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime/debug"
	"sync"
	"syscall"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/clients/stellarcore"
	"github.com/stellar/go/keypair"
	proto "github.com/stellar/go/protocols/horizon"
	stellarcoreproto "github.com/stellar/go/protocols/stellarcore"
	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/go/xdr"
	"github.com/stretchr/testify/require"

	"github.com/stellar/starbridge/app"
	"github.com/stellar/starbridge/client"
	"github.com/stellar/starbridge/controllers"
	"github.com/stellar/starbridge/stellar"
)

const (
	StandaloneNetworkPassphrase = "Standalone Network ; February 2017"
	EthereumBridgeAddress       = "0x31995201773dA53F950f15278Ea1538eA37A68A1"
	EthereumXLMTokenAddress     = "0x4Ee50847CD1278DBE87190080DD53055672755F6"
	EthereumRPCURL              = "http://127.0.0.1:8545"
	ethereumSenderPrivateKey    = "c1a4af60400ffd1473ada8425cff9f91b533194d6dd30424a17f356e418ac35b"
)

var (
	dockerHost     = "localhost"
	ethPrivateKeys = []string{
		// 0x89bFfeDAB59580576f7b95DbC500Ac1657EA9119
		"cff41ce3c1708e589b87198c9ee494eef407ca2a765a4353cf162c85ddc81cd9",
		// 0xAe1B35129e5924C3a7313EE579878f829f3e8495
		"51138e68e8a5fa906d38c5b42bc01b805d7adb3fce037743fb406bb10aa83307",
		// 0xCe3535F6f176128A882db28Cca00E2b1FbC38F09
		"0b1037a08795be0955e39e7e279e0530eb89e0ec06d372ff6f122a5a4e1a6f84",
	}
)

type Config struct {
	Servers                int
	EthereumFinalityBuffer uint64
	WithdrawalWindow       time.Duration
}

type Test struct {
	t *testing.T

	composePath string

	client        *http.Client
	horizonClient *horizonclient.Client
	coreClient    *stellarcore.Client

	app           []*app.App
	runningApps   *sync.WaitGroup
	shutdownOnce  sync.Once
	shutdownCalls []func()
	masterKey     *keypair.Full
	passPhrase    string

	mainKey     *keypair.Full
	signerKeys  []*keypair.Full
	clientKey   *keypair.Full
	mainAccount txnbuild.Account

	bridgeClient client.BridgeClient
}

// NewIntegrationTest starts a new environment for integration test.
//
// WARNING: This requires Docker Compose installed.
func NewIntegrationTest(t *testing.T, config Config) *Test {
	if os.Getenv("STARBRIDGE_INTEGRATION_TESTS_ENABLED") == "" {
		t.Skip("skipping integration test: STARBRIDGE_INTEGRATION_TESTS_ENABLED not set")
	}

	if host := os.Getenv("STARBRIDGE_INTEGRATION_TESTS_DOCKER_HOST"); host != "" {
		dockerHost = host
	}

	test := &Test{
		t:           t,
		composePath: findDockerComposePath(t),
		passPhrase:  StandaloneNetworkPassphrase,

		client: &http.Client{},
		horizonClient: &horizonclient.Client{
			HorizonURL: fmt.Sprintf("http://%s:8000", dockerHost),
		},
		coreClient: &stellarcore.Client{URL: fmt.Sprintf("http://%s:11626", dockerHost)},

		runningApps: &sync.WaitGroup{},
	}

	test.runComposeCommand("down", "-v")
	test.runComposeCommand("build")
	test.runComposeCommand("up", "--detach", "--quiet-pull", "--no-color", "quickstart")
	test.runComposeCommand("up", "--detach", "--no-color", "ethereum-node")
	test.runComposeCommand("up", "--no-color", "deploy-ethereum-contract")
	test.prepareShutdownHandlers()
	test.waitForHorizon()
	test.waitForFriendbot()

	if config.Servers == 0 {
		config.Servers = 1
	}

	// Create main account
	keys, accounts := test.CreateAccounts(1)
	test.mainKey, test.mainAccount = keys[0], accounts[0]

	test.signerKeys = make([]*keypair.Full, config.Servers)
	for i := 0; i < config.Servers; i++ {
		test.signerKeys[i] = keypair.MustRandom()
	}

	// Configure main account signers and configure client key
	test.clientKey = keypair.MustParseFull("SBEICGMVMPF2WWIYV34IP7ON2Q6BUOT7F7IGHOTUMYUIG5K4IWIOUQC3")

	threshold := txnbuild.Threshold(config.Servers/2) + 1
	ops := []txnbuild.Operation{
		&txnbuild.CreateAccount{
			Destination: test.clientKey.Address(),
			Amount:      "100",
		},
		&txnbuild.ChangeTrust{
			SourceAccount: test.clientKey.Address(),
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "ETH",
					Issuer: test.mainKey.Address(),
				},
			},
		},
		&txnbuild.SetOptions{
			MasterWeight:    txnbuild.NewThreshold(0),
			LowThreshold:    txnbuild.NewThreshold(threshold),
			MediumThreshold: txnbuild.NewThreshold(threshold),
			HighThreshold:   txnbuild.NewThreshold(txnbuild.Threshold(config.Servers)),
		},
	}

	for i := 0; i < config.Servers; i++ {
		ops = append(ops, &txnbuild.SetOptions{
			Signer: &txnbuild.Signer{
				Address: test.signerKeys[i].Address(),
				Weight:  txnbuild.Threshold(1),
			},
		})
	}

	test.MustSubmitMultiSigOperations(
		test.mainAccount,
		[]*keypair.Full{test.mainKey, test.clientKey},
		ops...,
	)

	stellarBridgeContractID := setupStellarBridge(test)
	test.bridgeClient = client.BridgeClient{
		ValidatorURLs: []string{
			"http://localhost:9000",
			"http://localhost:9001",
			"http://localhost:9002",
		},
		EthereumURL:                 EthereumRPCURL,
		EthereumChainID:             31337,
		HorizonURL:                  test.horizonClient.HorizonURL,
		NetworkPassphrase:           StandaloneNetworkPassphrase,
		EthereumBridgeAddress:       EthereumBridgeAddress,
		StellarBridgeAccount:        test.mainAccount.GetAccountID(),
		EthereumBridgeConfigVersion: 0,
		StellarPrivateKey:           test.clientKey.Seed(),
		EthereumPrivateKey:          ethereumSenderPrivateKey,
		StellarBridgeContractID:     stellarBridgeContractID,
	}

	stellarBridgeContractIDHex := hex.EncodeToString(stellarBridgeContractID[:])
	test.app = make([]*app.App, config.Servers)
	for i := 0; i < config.Servers; i++ {
		if innerErr := test.StartStarbridge(i, config, stellarBridgeContractIDHex); innerErr != nil {
			t.Fatalf("Failed to start Starbridge: %v", innerErr)
		}
	}
	test.waitForStarbridge(config.Servers)

	return test
}

func setupStellarBridge(itest *Test) xdr.Hash {
	// Install the contract
	installContractOp := assembleInstallContractCodeOp(itest.CurrentTest(), itest.Master().Address(), soroban_bridge.SorobanBridgeWasm)
	itest.MustSubmitOperations(itest.MasterAccount(), itest.Master(), installContractOp)

	// Create the contract
	createContractOp, bridgeContractID := assembleCreateContractOp(itest.t, itest.Master().Address(), soroban_bridge.SorobanBridgeWasm, "bridge", itest.passPhrase)

	tx, err := itest.SubmitOperations(itest.MasterAccount(), itest.Master(), createContractOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	initializeOp := initializeBridge(itest, bridgeContractID)
	tx, err = itest.SubmitOperations(itest.MasterAccount(), itest.Master(), initializeOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	ethAsset := xdr.MustNewCreditAsset("ETH", itest.mainKey.Address())
	createContractOp = createSAC(itest, itest.Master().Address(), ethAsset)
	tx, err = itest.SubmitOperations(itest.MasterAccount(), itest.Master(), createContractOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	createContractOp = createSAC(itest, itest.Master().Address(), xdr.MustNewNativeAsset())
	tx, err = itest.SubmitOperations(itest.MasterAccount(), itest.Master(), createContractOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	setAdminOp := setAdminOnSAC(itest, itest.mainAccount.GetAccountID(), ethAsset, bridgeContractID)
	mainAccount := itest.MustGetAccount(itest.mainKey)
	threshold := txnbuild.Threshold(len(itest.signerKeys)/2) + 1
	tx, err = itest.SubmitMultiSigOperations(&mainAccount, itest.signerKeys[:threshold], setAdminOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	clientAccount := itest.MustGetAccount(itest.clientKey)
	incrAllowEthOp := incrAllow(itest, itest.clientKey.Address(), ethAsset, bridgeContractID)
	tx, err = itest.SubmitOperations(&clientAccount, itest.clientKey, incrAllowEthOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	clientAccount = itest.MustGetAccount(itest.clientKey)
	incrAllowXlmOp := incrAllow(itest, itest.clientKey.Address(), xdr.MustNewNativeAsset(), bridgeContractID)
	tx, err = itest.SubmitOperations(&clientAccount, itest.clientKey, incrAllowXlmOp)
	require.NoError(itest.t, err)
	require.True(itest.t, tx.Successful)

	mainAccount = itest.MustGetAccount(itest.mainKey)
	itest.mainAccount = &mainAccount

	observer := stellar.NewObserver(bridgeContractID, itest.horizonClient, itest.coreClient)
	itest.t.Log(observer.GetRequestStatus(context.Background(), [32]byte{}))

	return bridgeContractID
}

func initializeBridge(itest *Test, contractID xdr.Hash) *txnbuild.InvokeHostFunction {
	return addFootprint(itest, &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeArgs: &xdr.ScVec{
				contractIDParam(contractID),
				functionNameParam("init"),
				accountIDEnumParam(itest.mainAccount.GetAccountID()),
			},
		},
		SourceAccount: itest.MasterAccount().GetAccountID(),
	})
}

func accountIDEnumParam(accountID string) xdr.ScVal {
	accountObj := &xdr.ScObject{
		Type:      xdr.ScObjectTypeScoAccountId,
		AccountId: xdr.MustAddressPtr(accountID),
	}
	accountSym := xdr.ScSymbol("Account")
	accountEnum := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoVec,
		Vec: &xdr.ScVec{
			xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  &accountSym,
			},
			xdr.ScVal{
				Type: xdr.ScValTypeScvObject,
				Obj:  &accountObj,
			},
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &accountEnum,
	}
}

func contractIDEnumParam(contractID xdr.Hash) xdr.ScVal {
	contractIDBytes := contractID[:]
	contractIDObj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoBytes,
		Bin:  &contractIDBytes,
	}
	contractSym := xdr.ScSymbol("Contract")
	contractEnum := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoVec,
		Vec: &xdr.ScVec{
			xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  &contractSym,
			},
			xdr.ScVal{
				Type: xdr.ScValTypeScvObject,
				Obj:  &contractIDObj,
			},
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &contractEnum,
	}
}

func functionNameParam(name string) xdr.ScVal {
	contractFnParameterSym := xdr.ScSymbol(name)
	return xdr.ScVal{
		Type: xdr.ScValTypeScvSymbol,
		Sym:  &contractFnParameterSym,
	}
}

func contractIDParam(contractID xdr.Hash) xdr.ScVal {
	contractIdBytes := contractID[:]
	contractIdParameterObj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoBytes,
		Bin:  &contractIdBytes,
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &contractIdParameterObj,
	}
}

func assembleInstallContractCodeOp(t *testing.T, sourceAccount string, contract []byte) *txnbuild.InvokeHostFunction {
	// Assemble the InvokeHostFunction CreateContract operation:
	// CAP-0047 - https://github.com/stellar/stellar-protocol/blob/master/core/cap-0047.md#creating-a-contract-using-invokehostfunctionop
	t.Logf("Contract File Contents: %v", hex.EncodeToString(contract))

	installContractCodeArgs, err := xdr.InstallContractCodeArgs{Code: contract}.MarshalBinary()
	assert.NoError(t, err)
	contractHash := sha256.Sum256(installContractCodeArgs)

	return &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInstallContractCode,
			InstallContractCodeArgs: &xdr.InstallContractCodeArgs{
				Code: contract,
			},
		},
		Footprint: xdr.LedgerFootprint{
			ReadWrite: []xdr.LedgerKey{
				{
					Type: xdr.LedgerEntryTypeContractCode,
					ContractCode: &xdr.LedgerKeyContractCode{
						Hash: contractHash,
					},
				},
			},
		},
		SourceAccount: sourceAccount,
	}
}

func assembleCreateContractOp(t *testing.T, sourceAccount string, contract []byte, contractSalt string, passPhrase string) (*txnbuild.InvokeHostFunction, xdr.Hash) {
	// Assemble the InvokeHostFunction CreateContract operation:
	// CAP-0047 - https://github.com/stellar/stellar-protocol/blob/master/core/cap-0047.md#creating-a-contract-using-invokehostfunctionop

	salt := sha256.Sum256([]byte(contractSalt))
	t.Logf("Salt hash: %v", hex.EncodeToString(salt[:]))

	networkId := xdr.Hash(sha256.Sum256([]byte(passPhrase)))
	preImage := xdr.HashIdPreimage{
		Type: xdr.EnvelopeTypeEnvelopeTypeContractIdFromSourceAccount,
		SourceAccountContractId: &xdr.HashIdPreimageSourceAccountContractId{
			NetworkId: networkId,
			Salt:      salt,
		},
	}
	preImage.SourceAccountContractId.SourceAccount.SetAddress(sourceAccount)
	xdrPreImageBytes, err := preImage.MarshalBinary()
	require.NoError(t, err)
	hashedContractID := sha256.Sum256(xdrPreImageBytes)

	saltParameter := xdr.Uint256(salt)

	installContractCodeArgs, err := xdr.InstallContractCodeArgs{Code: contract}.MarshalBinary()
	assert.NoError(t, err)
	contractHash := xdr.Hash(sha256.Sum256(installContractCodeArgs))

	ledgerKeyContractCodeAddr := xdr.ScStaticScsLedgerKeyContractCode
	ledgerKey := xdr.LedgerKeyContractData{
		ContractId: xdr.Hash(hashedContractID),
		Key: xdr.ScVal{
			Type: xdr.ScValTypeScvStatic,
			Ic:   &ledgerKeyContractCodeAddr,
		},
	}

	return &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeCreateContract,
			CreateContractArgs: &xdr.CreateContractArgs{
				ContractId: xdr.ContractId{
					Type: xdr.ContractIdTypeContractIdFromSourceAccount,
					Salt: &saltParameter,
				},
				Source: xdr.ScContractCode{
					Type:   xdr.ScContractCodeTypeSccontractCodeWasmRef,
					WasmId: &contractHash,
				},
			},
		},
		Footprint: xdr.LedgerFootprint{
			ReadWrite: []xdr.LedgerKey{
				{
					Type:         xdr.LedgerEntryTypeContractData,
					ContractData: &ledgerKey,
				},
			},
			ReadOnly: []xdr.LedgerKey{
				{
					Type: xdr.LedgerEntryTypeContractCode,
					ContractCode: &xdr.LedgerKeyContractCode{
						Hash: contractHash,
					},
				},
			},
		},
		SourceAccount: sourceAccount,
	}, hashedContractID
}

func createSAC(itest *Test, sourceAccount string, asset xdr.Asset) *txnbuild.InvokeHostFunction {
	return addFootprint(itest, &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeCreateContract,
			CreateContractArgs: &xdr.CreateContractArgs{
				ContractId: xdr.ContractId{
					Type:  xdr.ContractIdTypeContractIdFromAsset,
					Asset: &asset,
				},
				Source: xdr.ScContractCode{
					Type: xdr.ScContractCodeTypeSccontractCodeToken,
				},
			},
		},
		SourceAccount: sourceAccount,
	})
}

func stellarAssetContractID(t *testing.T, passPhrase string, asset xdr.Asset) xdr.Hash {
	networkId := xdr.Hash(sha256.Sum256([]byte(passPhrase)))
	preImage := xdr.HashIdPreimage{
		Type: xdr.EnvelopeTypeEnvelopeTypeContractIdFromAsset,
		FromAsset: &xdr.HashIdPreimageFromAsset{
			NetworkId: networkId,
			Asset:     asset,
		},
	}
	xdrPreImageBytes, err := preImage.MarshalBinary()
	require.NoError(t, err)
	return sha256.Sum256(xdrPreImageBytes)
}

func i128Param(hi, lo uint64) xdr.ScVal {
	i128Obj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoI128,
		I128: &xdr.Int128Parts{
			Hi: xdr.Uint64(hi),
			Lo: xdr.Uint64(lo),
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &i128Obj,
	}
}

func setAdminOnSAC(itest *Test, sourceAccount string, asset xdr.Asset, contractID xdr.Hash) *txnbuild.InvokeHostFunction {
	return addFootprint(itest, &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeArgs: &xdr.ScVec{
				contractIDParam(stellarAssetContractID(itest.CurrentTest(), itest.passPhrase, asset)),
				functionNameParam("set_admin"),
				invokerSignatureParam(),
				i128Param(0, 0),
				contractIDEnumParam(contractID),
			},
		},
		SourceAccount: sourceAccount,
	})
}

func incrAllow(itest *Test, sourceAccount string, asset xdr.Asset, spenderContractID xdr.Hash) *txnbuild.InvokeHostFunction {
	return addFootprint(itest, &txnbuild.InvokeHostFunction{
		Function: xdr.HostFunction{
			Type: xdr.HostFunctionTypeHostFunctionTypeInvokeContract,
			InvokeArgs: &xdr.ScVec{
				contractIDParam(stellarAssetContractID(itest.CurrentTest(), itest.passPhrase, asset)),
				functionNameParam("incr_allow"),
				invokerSignatureParam(),
				i128Param(0, 0),
				contractIDEnumParam(spenderContractID),
				i128Param(0, math.MaxInt64),
			},
		},
		SourceAccount: sourceAccount,
	})
}

func invokerSignatureParam() xdr.ScVal {
	invokerSym := xdr.ScSymbol("Invoker")
	obj := &xdr.ScObject{
		Type: xdr.ScObjectTypeScoVec,
		Vec: &xdr.ScVec{
			xdr.ScVal{
				Type: xdr.ScValTypeScvSymbol,
				Sym:  &invokerSym,
			},
		},
	}
	return xdr.ScVal{
		Type: xdr.ScValTypeScvObject,
		Obj:  &obj,
	}
}

func addFootprint(itest *Test, invokeHostFn *txnbuild.InvokeHostFunction) *txnbuild.InvokeHostFunction {
	opXDR, err := invokeHostFn.BuildXDR()
	require.NoError(itest.CurrentTest(), err)

	invokeHostFunctionOp := opXDR.Body.MustInvokeHostFunctionOp()

	// clear footprint so we can verify preflight response
	response, err := itest.coreClient.Preflight(
		context.Background(),
		invokeHostFn.SourceAccount,
		invokeHostFunctionOp,
	)
	require.NoError(itest.CurrentTest(), err)
	require.Equal(itest.CurrentTest(), stellarcoreproto.PreflightStatusOk, response.Status, response.Detail)
	err = xdr.SafeUnmarshalBase64(response.Footprint, &invokeHostFn.Footprint)
	require.NoError(itest.CurrentTest(), err)
	return invokeHostFn
}

// Runs a docker-compose command applied to the above configs
func (i *Test) runComposeCommand(args ...string) {
	integrationYaml := filepath.Join(i.composePath, "docker-compose.integration-tests.yml")

	cmdline := append([]string{"-f", integrationYaml}, args...)
	cmd := exec.Command("docker-compose", cmdline...)
	i.t.Log("Running", cmd.Env, cmd.Args)
	out, innerErr := cmd.Output()
	if exitErr, ok := innerErr.(*exec.ExitError); ok {
		fmt.Printf("stdout:\n%s\n", string(out))
		fmt.Printf("stderr:\n%s\n", string(exitErr.Stderr))
	}

	if innerErr != nil {
		i.t.Fatalf("Compose command failed: %v", innerErr)
	}
}

func (i *Test) prepareShutdownHandlers() {
	i.shutdownCalls = append(i.shutdownCalls,
		func() {
			i.runComposeCommand("down", "-v")
		},
		i.StopStarbridge,
	)

	// Register cleanup handlers (on panic and ctrl+c) so the containers are
	// stopped even if ingestion or testing fails.
	i.t.Cleanup(i.Shutdown)

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		i.Shutdown()
		os.Exit(int(syscall.SIGTERM))
	}()
}

// Shutdown stops the integration tests and destroys all its associated
// resources. It will be implicitly called when the calling test (i.e. the
// `testing.Test` passed to `New()`) is finished if it hasn't been explicitly
// called before.
func (i *Test) Shutdown() {
	i.shutdownOnce.Do(func() {
		// run them in the opposite order in which they were added
		for callI := len(i.shutdownCalls) - 1; callI >= 0; callI-- {
			i.shutdownCalls[callI]()
		}
	})
}

func (i *Test) StartStarbridge(id int, config Config, stellarBridgeContractID string) error {
	i.app[id] = app.NewApp(app.Config{
		Port:                    9000 + uint16(id),
		PostgresDSN:             fmt.Sprintf("postgres://postgres:mysecretpassword@%s:5641/starbridge%d?sslmode=disable", dockerHost, id),
		HorizonURL:              fmt.Sprintf("http://%s:8000/", dockerHost),
		CoreURL:                 fmt.Sprintf("http://%s:11626", dockerHost),
		NetworkPassphrase:       StandaloneNetworkPassphrase,
		StellarBridgeAccount:    i.mainAccount.GetAccountID(),
		StellarBridgeContractID: stellarBridgeContractID,
		StellarPrivateKey:       i.signerKeys[id].Seed(),
		EthereumRPCURL:          EthereumRPCURL,
		EthereumBridgeAddress:   EthereumBridgeAddress,
		EthereumPrivateKey:      ethPrivateKeys[id],
		AssetMapping: []controllers.AssetMappingConfigEntry{
			{
				StellarAsset:      "native",
				EthereumToken:     EthereumXLMTokenAddress,
				StellarToEthereum: "1",
			},
			{
				StellarAsset:      "ETH:" + i.mainAccount.GetAccountID(),
				EthereumToken:     (common.Address{}).String(),
				StellarToEthereum: "100000000000",
			},
		},
		EthereumFinalityBuffer: config.EthereumFinalityBuffer,
		WithdrawalWindow:       config.WithdrawalWindow,
	})

	i.runningApps.Add(1)
	go func() {
		defer i.runningApps.Done()
		i.app[id].RunHTTPServer()
	}()

	return nil
}

func (i *Test) waitForHorizon() uint32 {
	for t := 60; t >= 0; t -= 1 {
		time.Sleep(time.Second)

		i.t.Log("Waiting for ingestion and protocol upgrade...")
		root, err := i.horizonClient.Root()
		if err != nil {
			i.t.Log(err)
			continue
		}

		if root.HorizonSequence < 3 ||
			int(root.HorizonSequence) != int(root.IngestSequence) {
			continue
		}

		if uint32(root.CurrentProtocolVersion) != 0 {
			i.t.Logf("Horizon protocol version upgraded to %d",
				root.CurrentProtocolVersion)
			return root.IngestSequence
		}
	}

	i.t.Fatal("Horizon not ingesting...")
	return 0
}

func (i *Test) waitForFriendbot() {
	for t := 60; t >= 0; t -= 1 {
		time.Sleep(time.Second)

		i.t.Log("Waiting for friendbot...")
		url := fmt.Sprintf("http://%s:8000/friendbot", dockerHost)
		resp, err := http.Get(url)
		if err != nil {
			continue
		}

		if resp.StatusCode == http.StatusBadGateway {
			continue
		}

		return
	}

	i.t.Fatal("Friendbot not working...")
}

func (it *Test) waitForStarbridge(count int) {
	g := new(errgroup.Group)

	for i := 0; i < count; i++ {
		i := i
		g.Go(func() error {
			for t := 60; t >= 0; t -= 1 {
				time.Sleep(time.Second)

				port := 9000 + i
				it.t.Logf("Waiting for Starbridge at port %d...", port)
				url := fmt.Sprintf("http://localhost:%d", port)
				_, err := http.Get(url)
				if err != nil {
					continue
				}

				return nil
			}

			return errors.New("Starbridge not responding...")
		})
	}

	if err := g.Wait(); err != nil {
		it.t.Fatal(err)
	}
}

// HorizonClient returns horizon.Client connected to started Horizon instance.
func (i *Test) HorizonClient() *horizonclient.Client {
	return i.horizonClient
}

// Client returns http.Client connected to started Starbridge instance.
func (i *Test) Client() *http.Client {
	return i.client
}

// StopStarbridge shuts down the running starbridge processes
func (i *Test) StopStarbridge() {
	for _, app := range i.app {
		if app != nil {
			app.Close()
		}
	}
	i.runningApps.Wait()
	i.app = nil
}

// Master returns a keypair of the network masterKey account.
func (i *Test) Master() *keypair.Full {
	if i.masterKey != nil {
		return i.masterKey
	}
	return keypair.Master(i.passPhrase).(*keypair.Full)
}

func (i *Test) MasterAccount() txnbuild.Account {
	account := i.MasterAccountDetails()
	return &account
}

func (i *Test) MasterAccountDetails() proto.Account {
	return i.MustGetAccount(i.Master())
}

func (i *Test) CurrentTest() *testing.T {
	return i.t
}

/* Utility functions for easier test case creation. */

// Creates new accounts via friendbot.
//
// Returns: The slice of created keypairs and account objects.
//
// Note: panics on any errors, since we assume that tests cannot proceed without
// this method succeeding.
func (i *Test) CreateAccounts(count int) ([]*keypair.Full, []txnbuild.Account) {
	pairs := make([]*keypair.Full, count)
	accounts := make([]txnbuild.Account, count)

	g := new(errgroup.Group)

	for j := range pairs {
		j := j
		g.Go(func() error {
			pair, _ := keypair.Random()
			pairs[j] = pair

			_, err := i.horizonClient.Fund(pair.Address())
			if err != nil {
				return err
			}

			i.t.Logf("Funded %s (%s)\n", pair.Seed(), pair.Address())

			request := horizonclient.AccountRequest{AccountID: pair.Address()}
			account, err := i.horizonClient.AccountDetail(request)
			if err != nil {
				return err
			}

			accounts[j] = &account
			return nil
		})
	}

	if err := g.Wait(); err != nil {
		i.t.Fatal(err)
	}

	return pairs, accounts
}

// MustGetAccount panics on any error retrieves an account's details from its
// key. This means it must have previously been funded.
func (i *Test) MustGetAccount(source *keypair.Full) proto.Account {
	client := i.HorizonClient()
	account, err := client.AccountDetail(horizonclient.AccountRequest{AccountID: source.Address()})
	i.panicIf(err)
	return account
}

// MustSubmitOperations submits a signed transaction from an account with
// standard options.
//
// Namely, we set the standard fee, time bounds, etc. to "non-production"
// defaults that work well for tests.
//
// Most transactions only need one signer, so see the more verbose
// `MustSubmitOperationsWithSigners` below for multi-sig transactions.
//
// Note: We assume that transaction will be successful here so we panic in case
// of all errors. To allow failures, use `SubmitOperations`.
func (i *Test) MustSubmitOperations(
	source txnbuild.Account, signer *keypair.Full, ops ...txnbuild.Operation,
) proto.Transaction {
	tx, err := i.SubmitOperations(source, signer, ops...)
	i.panicIf(err)
	return tx
}

func (i *Test) SubmitOperations(
	source txnbuild.Account, signer *keypair.Full, ops ...txnbuild.Operation,
) (proto.Transaction, error) {
	return i.SubmitMultiSigOperations(source, []*keypair.Full{signer}, ops...)
}

func (i *Test) SubmitMultiSigOperations(
	source txnbuild.Account, signers []*keypair.Full, ops ...txnbuild.Operation,
) (proto.Transaction, error) {
	tx, err := i.CreateSignedTransactionFromOps(source, signers, ops...)
	if err != nil {
		return proto.Transaction{}, err
	}
	return i.HorizonClient().SubmitTransaction(tx)
}

func (i *Test) MustSubmitMultiSigOperations(
	source txnbuild.Account, signers []*keypair.Full, ops ...txnbuild.Operation,
) proto.Transaction {
	tx, err := i.SubmitMultiSigOperations(source, signers, ops...)
	i.panicIf(err)
	return tx
}

func (i *Test) MustSubmitTransaction(signer *keypair.Full, txParams txnbuild.TransactionParams,
) proto.Transaction {
	tx, err := i.SubmitTransaction(signer, txParams)
	i.panicIf(err)
	return tx
}

func (i *Test) SubmitTransaction(
	signer *keypair.Full, txParams txnbuild.TransactionParams,
) (proto.Transaction, error) {
	return i.SubmitMultiSigTransaction([]*keypair.Full{signer}, txParams)
}

func (i *Test) SubmitMultiSigTransaction(
	signers []*keypair.Full, txParams txnbuild.TransactionParams,
) (proto.Transaction, error) {
	tx, err := i.CreateSignedTransaction(signers, txParams)
	if err != nil {
		return proto.Transaction{}, err
	}
	return i.HorizonClient().SubmitTransaction(tx)
}

func (i *Test) MustSubmitMultiSigTransaction(
	signers []*keypair.Full, txParams txnbuild.TransactionParams,
) proto.Transaction {
	tx, err := i.SubmitMultiSigTransaction(signers, txParams)
	i.panicIf(err)
	return tx
}

func (i *Test) CreateSignedTransaction(signers []*keypair.Full, txParams txnbuild.TransactionParams,
) (*txnbuild.Transaction, error) {
	tx, err := txnbuild.NewTransaction(txParams)
	if err != nil {
		return nil, err
	}

	for _, signer := range signers {
		tx, err = tx.Sign(i.passPhrase, signer)
		if err != nil {
			return nil, err
		}
	}

	return tx, nil
}

func (i *Test) CreateSignedTransactionFromOps(
	source txnbuild.Account, signers []*keypair.Full, ops ...txnbuild.Operation,
) (*txnbuild.Transaction, error) {
	txParams := txnbuild.TransactionParams{
		SourceAccount:        source,
		Operations:           ops,
		BaseFee:              txnbuild.MinBaseFee,
		Preconditions:        txnbuild.Preconditions{TimeBounds: txnbuild.NewInfiniteTimeout()},
		IncrementSequenceNum: true,
	}

	return i.CreateSignedTransaction(signers, txParams)
}

// LogFailedTx is a convenience function to provide verbose information about a
// failing transaction to the test output log, if it's expected to succeed.
func (i *Test) LogFailedTx(txResponse proto.Transaction, horizonResult error) {
	t := i.CurrentTest()
	assert.NoErrorf(t, horizonResult, "Submitting the transaction failed")
	if prob := horizonclient.GetError(horizonResult); prob != nil {
		t.Logf("  problem: %s\n", prob.Problem.Detail)
		t.Logf("  extras: %s\n", prob.Problem.Extras["result_codes"])
		return
	}

	var txResult xdr.TransactionResult
	err := xdr.SafeUnmarshalBase64(txResponse.ResultXdr, &txResult)
	assert.NoErrorf(t, err, "Unmarshalling transaction failed.")
	assert.Equalf(t, xdr.TransactionResultCodeTxSuccess, txResult.Result.Code,
		"Transaction did not succeed: %d", txResult.Result.Code)
}

// Cluttering code with if err != nil is absolute nonsense.
func (i *Test) panicIf(err error) {
	if err != nil {
		debug.PrintStack()
		i.t.Fatal(err)
	}
}

// findDockerComposePath performs a best-effort attempt to find the project's
// Docker Compose files.
func findDockerComposePath(t *testing.T) string {
	// Lets you check if a particular directory contains a file.
	directoryContainsFilename := func(dir string, filename string) bool {
		files, innerErr := ioutil.ReadDir(dir)
		if innerErr != nil {
			t.Fatal(innerErr)
		}

		for _, file := range files {
			if file.Name() == filename {
				return true
			}
		}

		return false
	}

	current, err := os.Getwd()
	if err != nil {
		t.Fatal(err)
	}

	//
	// We have a primary and backup attempt for finding the necessary docker
	// files: via $GOPATH and via local directory traversal.
	//

	if gopath := os.Getenv("GOPATH"); gopath != "" {
		monorepo := filepath.Join(gopath, "src", "github.com", "stellar", "starbridge")
		if _, err = os.Stat(monorepo); !os.IsNotExist(err) {
			current = monorepo
		}
	}

	// In either case, we try to walk up the tree until we find "go.mod",
	// which we hope is the root directory of the project.
	for !directoryContainsFilename(current, "go.mod") {
		current, err = filepath.Abs(filepath.Join(current, ".."))

		// FIXME: This only works on *nix-like systems.
		if err != nil || filepath.Base(current)[0] == filepath.Separator {
			fmt.Println("Failed to establish project root directory.")
			panic(err)
		}
	}

	// Directly jump down to the folder that should contain the configs
	return filepath.Join(current, "integration")
}
