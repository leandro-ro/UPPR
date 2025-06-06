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
    /// @param pubKey Compressed public key (33 bytes, 0x02/0x03 prefix + X)
    /// @param signature ECDSA signature on the keccak256(pubKey)
    /// @param proof VRF proof [gammaX, gammaY, c, s] as 81 bytes
    /// @param epoch Epoch used as VRF input seed
    /// @return valid Whether the credential is valid
    /// @return errorCode See below
    ///
    /// Error codes:
    /// 0 = success
    /// 1 = invalid signature length
    /// 2 = signature mismatch
    /// 3 = invalid VRF proof
    /// 4 = revoked (token found in Bloom filter)
    function checkCredential(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch
    ) external view returns (
        bool valid,
        uint8 errorCode
    ) {
        if (signature.length != 65) {
            return (false, 1); // Invalid signature length
        }

        bytes32 pubKeyHash = keccak256(pubKey);
        address recovered = ecrecover(
            pubKeyHash,
            uint8(signature[64]),
            bytes32(signature[0:32]),
            bytes32(signature[32:64])
        );
        if (recovered != issuer) {
            return (false, 2); // Signature mismatch
        }

        uint256[2] memory pubkeyXY = VRF.decodePoint(pubKey);
        uint256[4] memory decodedProof = VRF.decodeProof(proof);

        // Encode epoch as 8-byte big-endian
        bytes memory message = new bytes(8);
        uint256 e = epoch;
        for (uint8 i = 0; i < 8; i++) {
            message[7 - i] = bytes1(uint8(e & 0xff));
            e >>= 8;
        }

        if (!VRF.verify(pubkeyXY, decodedProof, message)) {
            return (false, 3); // Invalid VRF proof
        }

        bytes32 token = VRF.gammaToHash(decodedProof[0], decodedProof[1]);

        (bool accepted, ) = bloom.testToken(abi.encodePacked(token));
        if (accepted) {
            return (false, 4); // Token found → revoked
        }

        return (true, 0); // Success
    }
}