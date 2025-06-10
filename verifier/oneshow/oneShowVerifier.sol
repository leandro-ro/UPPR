// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

import {CascadingBloomFilter} from "bloom/sol/cascadingBloomFilter.sol";
import "../vrf/VRF.sol";

/// @title OneShowVerifier
/// @notice Verifies credentials by checking an ECDSA signature on a VRF public key, reconstructing the VRF output, and querying a Bloom filter for revocation.
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

    /// @notice Verifies a credential by checking issuer authenticity, VRF validity, and non-revocation.
    /// @dev Off-chain calls are gas-free; on-chain usage incurs cost.
    /// @param pubKey Compressed VRF public key (33 bytes, SEC1 format)
    /// @param signature ECDSA signature over keccak256(pubKey), signed by the issuer
    /// @param proof VRF proof (81 bytes)
    /// @param epoch 64-bit challenge input (big-endian encoded)
    /// @return valid True if credential is valid and not revoked
    /// @return errorCode Code in [0–4] indicating the verification result
    /// (0: success, 1: signature format invalid, 2: signature invalid, 3: VRF verification failed, 4: revoked)
    function checkCredential(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch
    ) public view returns (bool valid, uint8 errorCode) {
        if (signature.length != 65) return (false, 1);

        address recovered = ecrecover(
            keccak256(pubKey),
            uint8(signature[64]),
            bytes32(signature[0:32]),
            bytes32(signature[32:64])
        );
        if (recovered != issuer) return (false, 2);

        uint256[2] memory pubkeyXY = VRF.decodePoint(pubKey);
        uint256[4] memory decodedProof = VRF.decodeProof(proof);

        bytes memory message = new bytes(8);
        for (uint8 i = 0; i < 8; i++) {
            message[7 - i] = bytes1(uint8(epoch >> (i * 8)));
        }

        if (!VRF.verify(pubkeyXY, decodedProof, message)) return (false, 3);

        bytes32 token = VRF.gammaToHash(decodedProof[0], decodedProof[1]);
        (bool revoked, ) = bloom.testToken(abi.encodePacked(token));

        return revoked ? (false, 4) : (true, 0);
    }

    /// @notice Gas-measurable variant of `checkCredential`, intended for benchmarking only.
    function measureCheckCredentialGas(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch
    ) external returns (bool valid, uint8 errorCode) {
        return checkCredential(pubKey, signature, proof, epoch);
    }

    /// @notice Efficient on-chain verification using precomputed elliptic curve data.
    /// @dev Saves gas by avoiding repeated EC operations.
    /// @param pubKey Compressed VRF public key (33 bytes)
    /// @param signature ECDSA signature over keccak256(pubKey)
    /// @param proof VRF proof: [gammaX, gammaY, c, s]
    /// @param epoch 64-bit challenge input (big-endian)
    /// @param uPoint Precomputed U = sB - cY
    /// @param vComponents Precomputed [Hx, Hy, cGammaX, cGammaY] for V = sH - cGamma
    /// @return valid True if credential is valid and not revoked
    /// @return errorCode Code in [0–4] indicating the verification result
    /// (0: success, 1: signature format invalid, 2: signature invalid, 3: VRF verification failed, 4: revoked)
    function checkCredentialFast(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch,
        uint256[2] calldata uPoint,
        uint256[4] calldata vComponents
    ) public view returns (bool valid, uint8 errorCode) {
        if (signature.length != 65) return (false, 1);

        address recovered = ecrecover(
            keccak256(pubKey),
            uint8(signature[64]),
            bytes32(signature[0:32]),
            bytes32(signature[32:64])
        );
        if (recovered != issuer) return (false, 2);

        uint256[2] memory pubkeyXY = VRF.decodePoint(pubKey);
        uint256[4] memory decodedProof = VRF.decodeProof(proof);

        bytes memory message = new bytes(8);
        for (uint8 i = 0; i < 8; i++) {
            message[7 - i] = bytes1(uint8(epoch >> (i * 8)));
        }

        if (!VRF.fastVerify(pubkeyXY, decodedProof, message, uPoint, vComponents)) return (false, 3);

        bytes32 token = VRF.gammaToHash(decodedProof[0], decodedProof[1]);
        (bool revoked, ) = bloom.testToken(abi.encodePacked(token));

        return revoked ? (false, 4) : (true, 0);
    }

    /// @notice Gas-measurable variant of `checkCredentialFast`, intended for benchmarking only.
    function measureCheckCredentialFastGas(
        bytes calldata pubKey,
        bytes calldata signature,
        bytes calldata proof,
        uint256 epoch,
        uint256[2] calldata uPoint,
        uint256[4] calldata vComponents
    ) external returns (bool valid, uint8 errorCode) {
        return checkCredentialFast(pubKey, signature, proof, epoch, uPoint, vComponents);
    }

    /// @notice Computes auxiliary EC points required for fast on-chain verification.
    /// @dev Offloads heavy elliptic curve arithmetic to off-chain clients.
    /// @param pubKey Compressed VRF public key (33 bytes)
    /// @param proof VRF proof (81 bytes)
    /// @param epoch 64-bit epoch input
    /// @return uPoint EC point U = sB - cY
    /// @return vComponents [Hx, Hy, cGammaX, cGammaY] for computing V = sH - cGamma
    function getFastVerifyParams(
        bytes calldata pubKey,
        bytes calldata proof,
        uint256 epoch
    ) external view returns (uint256[2] memory uPoint, uint256[4] memory vComponents) {
        uint256[2] memory pubkeyXY = VRF.decodePoint(pubKey);
        uint256[4] memory decodedProof = VRF.decodeProof(proof);

        bytes memory message = new bytes(8);
        for (uint8 i = 0; i < 8; i++) {
            message[7 - i] = bytes1(uint8(epoch >> (i * 8)));
        }

        return VRF.computeFastVerifyParams(pubkeyXY, decodedProof, message);
    }
}