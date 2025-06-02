package zkp

import (
	"context"
	"fmt"
	"math/big"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
)

// Helper to compile Solidity and return paths to .bin and .abi
func compileSolidityContract(solFile string) (string, string, error) {
	cmd := exec.Command(
		"solc",
		"--bin",
		"--overwrite",
		"--abi",
		"--evm-version", "istanbul",
		solFile,
		"-o", "build",
	)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return "", "", fmt.Errorf("solc failed: %s\n%s", err, string(output))
	}

	return filepath.Join("build", "Verifier.bin"), filepath.Join("build", "Verifier.abi"), nil
}

// Test function that deploys the contract and prints deployment gas cost
func TestDeployRevocationTokenVerifier(t *testing.T) {
	privateKey, err := crypto.GenerateKey()
	if err != nil {
		t.Fatalf("key generation failed: %v", err)
	}

	auth, err := bind.NewKeyedTransactorWithChainID(privateKey, big.NewInt(1337))
	if err != nil {
		t.Fatalf("auth setup failed: %v", err)
	}

	address := crypto.PubkeyToAddress(privateKey.PublicKey)
	alloc := map[common.Address]core.GenesisAccount{
		address: {Balance: big.NewInt(1e18)},
	}

	sim := backends.NewSimulatedBackend(alloc, 10_000_000)

	_, _, err = compileSolidityContract("revocationTokenVerifier.sol")
	if err != nil {
		t.Fatalf("contract compilation failed: %v", err)
	}

	contractName := "Verifier"
	binPath := filepath.Join("build", contractName+".bin")
	abiPath := filepath.Join("build", contractName+".abi")

	binData, err := os.ReadFile(binPath)
	if err != nil {
		t.Fatalf("reading bin failed: %v", err)
	}
	abiData, err := os.ReadFile(abiPath)
	if err != nil {
		t.Fatalf("reading abi failed: %v", err)
	}

	parsedABI, err := abi.JSON(strings.NewReader(string(abiData)))
	if err != nil {
		t.Fatalf("ABI parsing failed: %v", err)
	}

	contractAddr, tx, _, err := bind.DeployContract(auth, parsedABI, common.FromHex(string(binData)), sim)
	if err != nil {
		t.Fatalf("contract deployment failed: %v", err)
	}

	sim.Commit()

	receipt, err := sim.TransactionReceipt(context.Background(), tx.Hash())
	if err != nil {
		t.Fatalf("could not get receipt: %v", err)
	}

	t.Logf("Contract deployed at: %s", contractAddr.Hex())
	t.Logf("Deployment gas used: %d", receipt.GasUsed)
}
