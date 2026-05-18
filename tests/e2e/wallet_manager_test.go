package e2e_test

import (
	"strings"
	"testing"

	"github.com/2Finance-Labs/go-client-2finance/wallet_manager"
)

const (
	testWalletPassword    = "test-password-123"
	testWalletNewPassword = "test-new-password-456"
)

func TestWalletManager_GenerateEd25519KeyPairHex(t *testing.T) {
	publicKey, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if publicKey == "" {
		t.Fatal("expected public key")
	}

	if privateKey == "" {
		t.Fatal("expected private key")
	}
}

func TestWalletManager_ImportPrivateKey(t *testing.T) {
	manager := setupWalletManager(t)

	publicKey, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	gotOwner := manager.OwnerAddress()
	if gotOwner != publicKey {
		t.Fatalf("expected owner %q, got %q", publicKey, gotOwner)
	}

	if manager.IsUnlocked() {
		t.Fatal("expected wallet to be locked after import")
	}
}

func TestWalletManager_ImportPrivateKey_RequiresPassword(t *testing.T) {
	manager := setupWalletManager(t)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	err = manager.ImportPrivateKey([]byte(privateKey), "")
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "password is required")
}

func TestWalletManager_ImportPrivateKey_RequiresPrivateKey(t *testing.T) {
	manager := setupWalletManager(t)

	err := manager.ImportPrivateKey(nil, testWalletPassword)
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "private key is required")
}

func TestWalletManager_UnlockWithPassword(t *testing.T) {
	manager := setupWalletManager(t)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	if manager.IsUnlocked() {
		t.Fatal("expected wallet to be locked before unlock")
	}

	if err := manager.UnlockWithPassword(testWalletPassword); err != nil {
		t.Fatalf("UnlockWithPassword error: %v", err)
	}

	if !manager.IsUnlocked() {
		t.Fatal("expected wallet to be unlocked")
	}
}

func TestWalletManager_UnlockWithPassword_RequiresPassword(t *testing.T) {
	manager := setupWalletManager(t)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	err = manager.UnlockWithPassword("")
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "password is required")
}

func TestWalletManager_UnlockWithPassword_WrongPassword(t *testing.T) {
	manager := setupWalletManager(t)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	err = manager.UnlockWithPassword("wrong-password")
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "failed to decrypt wallet file")
}

func TestWalletManager_Lock(t *testing.T) {
	manager := setupWalletManager(t)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	if err := manager.UnlockWithPassword(testWalletPassword); err != nil {
		t.Fatalf("UnlockWithPassword error: %v", err)
	}

	if !manager.IsUnlocked() {
		t.Fatal("expected wallet to be unlocked before lock")
	}

	if err := manager.Lock(); err != nil {
		t.Fatalf("Lock error: %v", err)
	}

	if manager.IsUnlocked() {
		t.Fatal("expected wallet to be locked")
	}
}

func TestWalletManager_ChangePassword(t *testing.T) {
	manager := setupWalletManager(t)

	_, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	if err := manager.UnlockWithPassword(testWalletPassword); err != nil {
		t.Fatalf("UnlockWithPassword error: %v", err)
	}

	if !manager.IsUnlocked() {
		t.Fatal("expected wallet to be unlocked before password change")
	}

	if err := manager.ChangePassword(testWalletPassword, testWalletNewPassword); err != nil {
		t.Fatalf("ChangePassword error: %v", err)
	}

	if manager.IsUnlocked() {
		t.Fatal("expected wallet to be locked after password change")
	}

	err = manager.UnlockWithPassword(testWalletPassword)
	if err == nil {
		t.Fatal("expected old password to fail after password change")
	}

	if err := manager.UnlockWithPassword(testWalletNewPassword); err != nil {
		t.Fatalf("expected new password to unlock wallet: %v", err)
	}

	if !manager.IsUnlocked() {
		t.Fatal("expected wallet to be unlocked with new password")
	}
}

func TestWalletManager_ChangePassword_RequiresCurrentPassword(t *testing.T) {
	manager := setupWalletManager(t)

	err := manager.ChangePassword("", testWalletNewPassword)
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "current password is required")
}

func TestWalletManager_ChangePassword_RequiresNewPassword(t *testing.T) {
	manager := setupWalletManager(t)

	err := manager.ChangePassword(testWalletPassword, "")
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "new password is required")
}

func TestWalletManager_ChangePassword_NewPasswordMustBeDifferent(t *testing.T) {
	manager := setupWalletManager(t)

	err := manager.ChangePassword(testWalletPassword, testWalletPassword)
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "new password must be different from current password")
}

func TestWalletManager_ChangePassword_RequiresOwner(t *testing.T) {
	manager := setupWalletManager(t)

	err := manager.ChangePassword(testWalletPassword, testWalletNewPassword)
	if err == nil {
		t.Fatal("expected error")
	}

	assertWalletErrorContains(t, err, "owner is required")
}

func TestWalletManager_OwnerAddress(t *testing.T) {
	manager := setupWalletManager(t)

	publicKey, privateKey, err := wallet_manager.GenerateEd25519KeyPairHex()
	if err != nil {
		t.Fatalf("GenerateEd25519KeyPairHex error: %v", err)
	}

	if err := manager.ImportPrivateKey([]byte(privateKey), testWalletPassword); err != nil {
		t.Fatalf("ImportPrivateKey error: %v", err)
	}

	gotOwner := manager.OwnerAddress()
	if gotOwner != publicKey {
		t.Fatalf("expected owner %q, got %q", publicKey, gotOwner)
	}
}

func assertWalletErrorContains(t *testing.T, err error, expected string) {
	t.Helper()

	if err == nil {
		t.Fatalf("expected error containing %q, got nil", expected)
	}

	if !strings.Contains(err.Error(), expected) {
		t.Fatalf("expected error containing %q, got %q", expected, err.Error())
	}
}