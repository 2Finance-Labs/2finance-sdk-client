package e2e_test

import (
	"strings"
	"testing"

	"github.com/2Finance-Labs/go-client-2finance/wallet_manager"
	"gitlab.com/2finance/2finance-network/blockchain/encryption/keys"
)

const testClientWalletPassword = "test-wallet-password-123"

func TestNetworkClient_SignAndSendTransaction_ReturnsWalletLocked_WhenWalletIsLocked(t *testing.T) {
	walletPath := t.TempDir() + "/wallet.json"

	walletManager := wallet_manager.NewWalletManager(walletPath)

	publicKey, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := walletManager.ImportPrivateKey([]byte(privateKey), testClientWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	if walletManager.IsUnlocked() {
		t.Fatal("expected wallet to be locked after import")
	}

	client := setupClient(t, walletManager)

	output, err := client.SignAndSendTransaction(
		1,
		publicKey,
		publicKey,
		"AddCashback",
		map[string]interface{}{
			"percentage": "10",
		},
		1,
		"018f5f4e-8f6a-7c4a-9a2e-000000000005",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "wallet is locked") {
		t.Fatalf("expected wallet locked error, got: %v", err)
	}

	_ = output
}

func TestNetworkClient_SignAndSendTransaction_ValidatesFromAddress(t *testing.T) {
	walletPath := t.TempDir() + "/wallet.json"

	walletManager := wallet_manager.NewWalletManager(walletPath)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := walletManager.ImportPrivateKey([]byte(privateKey), testClientWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	if err := walletManager.UnlockWithPassword(testClientWalletPassword); err != nil {
		t.Fatalf("UnlockWithPassword error: %v", err)
	}

	client := setupClient(t, walletManager)

	output, err := client.SignAndSendTransaction(
		1,
		"invalid-from-address",
		"invalid-to-address",
		"AddCashback",
		map[string]interface{}{
			"percentage": "10",
		},
		1,
		"018f5f4e-8f6a-7c4a-9a2e-000000000006",
	)

	if err == nil {
		t.Fatal("expected error")
	}

	if !strings.Contains(err.Error(), "invalid from address") {
		t.Fatalf("expected invalid from address error, got: %v", err)
	}

	_ = output
}

func TestNetworkClient_SignAndSendTransaction_CanReachSigning_WhenWalletIsUnlocked(t *testing.T) {
	walletPath := t.TempDir() + "/wallet.json"

	walletManager := wallet_manager.NewWalletManager(walletPath)

	publicKey, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := keys.ValidateEDDSAPublicKeyHex(publicKey); err != nil {
		t.Fatalf("generated public key should be valid: %v", err)
	}

	if err := walletManager.ImportPrivateKey([]byte(privateKey), testClientWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	if err := walletManager.UnlockWithPassword(testClientWalletPassword); err != nil {
		t.Fatalf("UnlockWithPassword error: %v", err)
	}

	if !walletManager.IsUnlocked() {
		t.Fatal("expected wallet to be unlocked")
	}

	// We do not call SignAndSendTransaction here because it would try to send
	// to the network through SendTransaction.
	//
	// This test validates the exact precondition required by your usage:
	//
	// walletManager.UnlockWithPassword(password)
	// client.SetWalletManager(walletManager)
	// client.AddCashback(...)
}