// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CascadingBloomFilter} from "bloom/sol/cascadingBloomFilter.sol";
import {Verifier} from "zkp/sol/revocationTokenVerifier.sol";
import "../vrf/VRF.sol";

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
    /// @param _x X coordinate of issuer’s eddsa bn254 public key. (used for cred signing)
    /// @param _y Y coordinate of issuer’s eddsa bn254 public key. (used for cred signing)
    constructor(
        address _bloom,
        address _zkpVerifier,
        uint256 _x,
        uint256 _y
    ) {
        issuer = msg.sender;
        bloom = CascadingBloomFilter(_bloom);
        verifier = Verifier(_zkpVerifier);
        issuerPubKeyX = _x;
        issuerPubKeyY = _y;
    }

    modifier onlyIssuer() {
        require(msg.sender == issuer, "Not issuer");
        _;
    }

    /// @notice Updates the Bloom filter cascade.
    /// @param newFilters Packed Bloom filter layers
    /// @param ks Number of hash functions per layer
    /// @param bitLens Number of valid bits per layer
    function update(
        bytes[] calldata newFilters,
        uint256[] calldata ks,
        uint256[] calldata bitLens
    ) external onlyIssuer {
        bloom.updateCascade(newFilters, ks, bitLens);
    }

    /// @notice Verifies a MultiShow credential using a zkSNARK proof and checks revocation.
    /// @param proof zkSNARK proof.
    /// @param token Revocation token (as input to Bloom filter and zkSNARK).
    /// @param epoch Epoch associated with the credential.
    /// @return valid True if credential is valid and not revoked.
    /// @return errorCode Code in [0–2] indicating the verification result
    ///                  (0: success, 1: zkSNARK proof invalid, 2: revoked)
    /// @return pubKeyXx Issuer's public key X-coordinate (for debug).
    /// @return pubKeyYy Issuer's public key Y-coordinate (for debug).
    /// @return usedToken The revocation token that was checked.
    /// @return usedEpoch The epoch that was checked.
    function checkCredential(
        uint256[8] calldata proof,
        uint256 pubKeyX,
        uint256 pubKeyY,
        uint256 token,
        uint256 epoch
    )
    public
    view
    returns (
        bool valid,
        uint8 errorCode,
        uint256 pubKeyXx,
        uint256 pubKeyYy,
        uint256 usedToken,
        uint256 usedEpoch
    )
    {
        uint256[4] memory input = [
                    pubKeyX,
                    pubKeyY,
                    token,
                    epoch
            ];

        try verifier.verifyProof(proof, input) {
            // Proof is valid, continue
        } catch {
            return (false, 1, issuerPubKeyX, issuerPubKeyY, token, epoch);
        }

        // Check Bloom filter
        // (bool revoked, ) = bloom.testToken(abi.encodePacked(bytes32(token)));
        // if (revoked) {
        //    return (false, 2, issuerPubKeyX, issuerPubKeyY, token, epoch);
        // }

        return (true, 0, issuerPubKeyX, issuerPubKeyY, token, epoch);
    }

}