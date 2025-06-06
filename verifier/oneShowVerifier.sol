// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CascadingBloomFilter} from "bloom/sol/cascadingBloomFilter.sol";
import "./vrf/VRF.sol";

/// @title OneShowVerifier
/// @notice Verifies credentials by checking an ECDSA signature on a VRF public key hash and looking up a revocation token in a Bloom filter.
contract OneShowVerifier {
    CascadingBloomFilter public bloom;
    address public issuer;

    constructor(address _bloom) {
        bloom = CascadingBloomFilter(_bloom);
        issuer = msg.sender;
    }

    modifier onlyIssuer() {
        require(msg.sender == issuer, "Not issuer");
        _;
    }

    /// @notice Updates the Bloom filter by delegating to the underlying contract.
    /// @param newFilters Packed Bloom filter layers
    /// @param ks Hash function count per layer
    /// @param bitLens Valid bit lengths per layer
    function update(bytes[] calldata newFilters, uint256[] calldata ks, uint256[] calldata bitLens) external onlyIssuer {
        bloom.updateCascade(newFilters, ks, bitLens);
    }

    /// @notice Checks if a credential is valid (not revoked).
    /// @param pubKey The uncompressed public key (64 bytes: x||y or 65 bytes with prefix)
    /// @param signature ECDSA signature on the keccak256 hash of the public key
    /// @param token Revocation token (VRF output)
    /// @return valid Whether the credential is valid
    /// @return errorCode A numeric error code:
    ///         0 = valid,
    ///         1 = invalid signature length,
    ///         2 = signature does not match issuer,
    ///         3 = revoked credential (bloom filter rejected)
    /// @return pubKeyHash The keccak256(pubKey), for debugging purposes
    /// @return issuerAddress The issuer address used for signature verification
    function checkCredential(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata token
    ) external view returns (
        bool valid,
        uint8 errorCode,
        bytes32 pubKeyHash,
        address issuerAddress
    ) {
        pubKeyHash = keccak256(pubKey);
        issuerAddress = issuer;

        if (signature.length != 65) {
            return (false, 1, pubKeyHash, issuerAddress); // Invalid signature length
        }

        address recovered = ecrecover(
            pubKeyHash,
            uint8(signature[64]),
            bytes32(signature[0:32]),
            bytes32(signature[32:64])
        );
        if (recovered != issuer) {
            return (false, 2, pubKeyHash, issuerAddress); // Signature mismatch
        }

        (bool accepted, ) = bloom.testToken(token);
        if (accepted) { // If in bloom filter, we have a revoked credential.
            return (false, 3, pubKeyHash, issuerAddress);
        }

        return (true, 0, pubKeyHash, issuerAddress); // Success
    }
}