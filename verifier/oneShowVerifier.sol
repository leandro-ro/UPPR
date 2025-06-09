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

    /// @notice Verifies the validity of a credential by checking an ECDSA signature, a VRF proof, and a Bloom filter.
    /// @dev This function is `view` and free to call from off-chain (e.g. via `eth_call`). On-chain calls consume gas.
    /// @param pubKey Compressed VRF public key (33 bytes, SEC1 format: 0x02 or 0x03 prefix + X coordinate)
    /// @param signature ECDSA signature on `keccak256(pubKey)`, signed by the issuer address
    /// @param proof VRF proof as 81-byte concatenation of gammaX (32), gammaY (32), c (16), s (1–32 padded to 32)
    /// @param epoch 64-bit epoch value (used as input seed to the VRF, encoded big-endian)
    /// @return valid `true` if the credential is valid and not revoked, otherwise `false`
    /// @return errorCode A numeric code indicating the failure reason:
    ///         0 = success (valid credential)
    ///         1 = invalid signature length
    ///         2 = signature mismatch (invalid issuer)
    ///         3 = VRF proof verification failed
    ///         4 = token found in Bloom filter (revoked)
    function checkCredential(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch
    ) public view returns (
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

    /// @notice Calls `checkCredential` in a state-changing context to allow gas benchmarking via transaction receipts.
    /// @dev Intended for gas measurement and testing only. Unlike the `view` function, this version consumes gas
    ///      when invoked from a transaction. The internal logic and result are identical to `checkCredential`.
    /// @param pubKey Compressed VRF public key (33 bytes, SEC1 format: 0x02 or 0x03 prefix + X coordinate)
    /// @param signature ECDSA signature on `keccak256(pubKey)`, signed by the issuer address
    /// @param proof VRF proof as 81-byte concatenation of gammaX (32), gammaY (32), c (16), s (1–32 padded to 32)
    /// @param epoch 64-bit epoch value (used as input seed to the VRF, encoded big-endian)
    /// @return valid `true` if the credential is valid and not revoked, otherwise `false`
    /// @return errorCode A numeric code indicating the failure reason:
    ///         0 = success (valid credential)
    ///         1 = invalid signature length
    ///         2 = signature mismatch (invalid issuer)
    ///         3 = VRF proof verification failed
    ///         4 = token found in Bloom filter (revoked)
    function measureCheckCredentialGas(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch
    ) external returns (bool valid, uint8 errorCode) {
        (valid, errorCode) = checkCredential(pubKey, signature, proof, epoch);
        return (valid, errorCode);
    }

/// @notice Checks if a credential is valid (not revoked) using gas-optimized fast VRF verification.
/// @dev Uses `VRF.fastVerify` with precomputed elliptic curve points to reduce gas cost.
/// @param pubKey Compressed VRF public key (33 bytes: 0x02/0x03 prefix + X)
/// @param signature ECDSA signature over keccak256(pubKey), must be 65 bytes
/// @param proof VRF proof as four 256-bit integers: [gammaX, gammaY, c, s]
/// @param epoch Epoch input (used to derive VRF message via 8-byte big-endian encoding)
/// @param uPoint EC point U = s·B − c·Y, required for fast verification
/// @param vComponents Precomputed data to reconstruct V = s·H − c·Gamma: [Hx, Hy, cGammaX, cGammaY]
/// @return valid Whether the credential is valid
/// @return errorCode Status code:
///         0 = success,
///         1 = invalid signature length,
///         2 = signature mismatch,
///         3 = invalid VRF proof,
///         4 = credential is revoked (token in Bloom filter)
    function checkCredentialFast(
        bytes calldata pubKey,
        bytes calldata signature,
        uint256[4] calldata proof,
        uint256 epoch,
        uint256[2] calldata uPoint,
        uint256[4] calldata vComponents
    ) public view returns (
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

        // Encode epoch as 8-byte big-endian
        bytes memory message = new bytes(8);
        uint256 e = epoch;
        for (uint8 i = 0; i < 8; i++) {
            message[7 - i] = bytes1(uint8(e & 0xff));
            e >>= 8;
        }

        bool ok = VRF.fastVerify(pubkeyXY, proof, message, uPoint, vComponents);
        if (!ok) {
            return (false, 3); // Invalid VRF proof
        }

        bytes32 token = VRF.gammaToHash(proof[0], proof[1]);
        (bool accepted, ) = bloom.testToken(abi.encodePacked(token));
        if (accepted) {
            return (false, 4); // Token found → credential revoked
        }

        return (true, 0); // Success
    }
}
