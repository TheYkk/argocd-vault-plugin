package vault_test

import (
	"encoding/json"
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"testing"

	"github.com/IBM/argocd-vault-plugin/pkg/helpers"
	"github.com/IBM/argocd-vault-plugin/pkg/vault"
)

func writeToken(token string) error {
	home, err := os.UserHomeDir()
	if err != nil {
		return err
	}
	path := filepath.Join(home, ".avp")
	os.Mkdir(path, 0755)
	data := map[string]interface{}{
		"vault_token": token,
	}
	file, _ := json.MarshalIndent(data, "", " ")
	err = ioutil.WriteFile(filepath.Join(path, "config.json"), file, 0644)
	if err != nil {
		return err
	}
	return nil
}

func removeToken() error {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".avp")
	err := os.RemoveAll(path)
	if err != nil {
		return err
	}
	return nil
}

func readToken() interface{} {
	home, _ := os.UserHomeDir()
	path := filepath.Join(home, ".avp", "config.json")
	dat, _ := ioutil.ReadFile(path)
	var result map[string]interface{}
	json.Unmarshal([]byte(dat), &result)
	return result["vault_token"]
}

func TestSetToken(t *testing.T) {
	cluster, _, _ := helpers.CreateTestAppRoleVault(t)
	defer cluster.Cleanup()

	vc := &vault.Client{
		VaultAPIClient: cluster.Cores[0].Client,
	}

	vault.SetToken(vc, "token")

	err := removeToken()
	if err != nil {
		t.Fatal(err)
	}
}

func TestReadVaultSecret(t *testing.T) {
	cluster, _, _ := helpers.CreateTestAppRoleVault(t)
	defer cluster.Cleanup()

	vc := &vault.Client{
		VaultAPIClient: cluster.Cores[0].Client,
	}

	_, err := vault.ReadVaultSecret(*vc, "kv/data/test", "3")
	expected := "Unsupported kvVersion specified"
	if !reflect.DeepEqual(err.Error(), expected) {
		t.Errorf("expected: %s, got: %s.", expected, err.Error())
	}
}

func TestReadVaultSecretWrongPath(t *testing.T) {
	cluster, _, _ := helpers.CreateTestAppRoleVault(t)
	defer cluster.Cleanup()

	vc := &vault.Client{
		VaultAPIClient: cluster.Cores[0].Client,
	}

	_, err := vault.ReadVaultSecret(*vc, "kv/test", "2")
	expected := "The Vault path: kv/test is empty - did you forget to include /data/ in the Vault path for kv-v2?"
	if err == nil {
		t.Fatalf("Vault path kv/test should be non-existent for kv-v2 Vault")
	}
	if !reflect.DeepEqual(err.Error(), expected) {
		t.Errorf("expected: %s, got: %s.", expected, err.Error())
	}
}
