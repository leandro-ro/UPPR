// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CascadingBloomFilter} from "bloom/sol/cascadingBloomFilter.sol";
import {Verifier} from "../zkp/sol/revocationTokenVerifier.sol";

/// @title MultiShowVerifier
/// @notice Verifies revocation status of MultiShow credentials via zkSNARK proof and Bloom filter.
contract MultiShowVerifier {
    CascadingBloomFilter public bloom;
    Verifier public verifier;

    address public issuer;
    uint256 public issuerPubKeyX;
    uint256 public issuerPubKeyY;

    /// @notice Deploys the verifier with a reference to Bloom filter and ZK proof verifier.
    /// @param _bloom Address of the Bloom filter contract.
    /// @param _zkpVerifier Address of the ZKP verifier contract.
    /// @param _x X coordinate of issuer’s secp256k1 public key.
    /// @param _y Y coordinate of issuer’s secp256k1 public key.
    constructor(
        address _bloom,
        address _zkpVerifier,
        uint256 _x,
        uint256 _y
    ) {
        // Validate that (_x, _y) actually corresponds to msg.sender
        bytes memory pubkey = abi.encodePacked(
            hex"04", // Uncompressed SEC1 format
            _x,
            _y
        );
        address derived = address(uint160(uint256(keccak256(pubkey))));
        require(derived == msg.sender, "Public key does not match sender");

        issuer = msg.sender;
        bloom = CascadingBloomFilter(_bloom);
        verifier = Verifier(_zkpVerifier);
        issuerPubKeyX = _x;
        issuerPubKeyY = _y;
    }

    /// @notice Verifies a MultiShow credential using a zkSNARK proof and checks revocation.
    /// @param proof zkSNARK proof.
    /// @param token Revocation token (as input to Bloom filter and zkSNARK).
    /// @param epoch Epoch associated with the credential.
    /// @return valid Whether the credential is valid (proof ok and token not revoked).
    function checkCredential(
        uint256[8] calldata proof,     // Groth16 proof
        uint256 token,
        uint256 epoch
    ) public view returns (bool valid, uint8 errorCode) {
        uint256[4] memory input = [
                    issuerPubKeyX,
                    issuerPubKeyY,
                    token,
                    epoch
            ];

        // Verify zkSNARK
        if (!verifier.verifyProof(proof, input)) {
            return (false, 1);
        }

        // Check Bloom filter
        (bool revoked, ) = bloom.testToken(abi.encodePacked(bytes32(token)));
        return revoked ? (false, 2) : (true, 0);
    }
}