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

    /// @notice Checks if a credential is valid (not revoked).
    /// @param pubKey Compressed public key (33 bytes, 0x02/0x03 prefix + X)
    /// @param signature ECDSA signature on the keccak256(pubKey)
    /// @param proof VRF proof [gammaX, gammaY, c, s] as 81 bytes
    /// @param epoch Epoch used as VRF input seed
    /// @return valid Whether the credential is valid
    /// @return errorCode See below
    /// @return pubKeyHash keccak256(pubKey)
    /// @return issuerAddress address used for ECDSA recovery
    ///
    /// Error codes:
    /// 0 = success
    /// 1 = invalid signature length
    /// 2 = signature mismatch
    /// 3 = invalid VRF proof
    /// 4 = revoked (token found in Bloom filter)
    function checkCredentialVrfDebug(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch
    )
    external view
    returns (
        bool valid,
        uint8 errorCode,
        bytes32 pubKeyHash,
        address issuerAddress,
        uint256[2] memory pubkeyXY,
        uint256[4] memory decodedProof,
        bytes32 token
    )
    {
        pubKeyHash = keccak256(pubKey);
        issuerAddress = issuer;

        if (signature.length != 65) {
            return (false, 1, pubKeyHash, issuerAddress, [uint256(0), 0], [uint256(0), 0, 0, 0], bytes32(0));
        }

        address recovered = ecrecover(
            pubKeyHash,
            uint8(signature[64]),
            bytes32(signature[0:32]),
            bytes32(signature[32:64])
        );
        if (recovered != issuer) {
            return (false, 2, pubKeyHash, issuerAddress, [uint256(0), 0], [uint256(0), 0, 0, 0], bytes32(0));
        }

        pubkeyXY = VRF.decodePoint(pubKey);
        decodedProof = VRF.decodeProof(proof);
        // Encode epoch as 8-byte big-endian
        bytes memory message = new bytes(8);
        uint256 e = epoch;
        for (uint8 i = 0; i < 8; i++) {
            message[7 - i] = bytes1(uint8(e & 0xff));
            e >>= 8;
        }

        if (!VRF.verify(pubkeyXY, decodedProof, message)) {
            return (false, 3, pubKeyHash, issuerAddress, pubkeyXY, decodedProof, bytes32(0));
        }

        token = VRF.gammaToHash(decodedProof[0], decodedProof[1]);

        (bool accepted, ) = bloom.testToken(abi.encodePacked(token));
        if (accepted) {
            return (false, 4, pubKeyHash, issuerAddress, pubkeyXY, decodedProof, token);
        }

        return (true, 0, pubKeyHash, issuerAddress, pubkeyXY, decodedProof, token);
    }
}