// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @title  Cascading Bloom Filter (multi-layer, owner can batch-upload)
/// @notice Minimal “ownable” pattern; owner can upload a multi-layer Bloom cascade.
///         Anyone can test a token against all layers. This version avoids “stack too deep.”
contract CascadingBloomFilter {
    // -----------------------------------------------------------------------
    // Simple ownable pattern
    // -----------------------------------------------------------------------
    address private _owner;

    modifier onlyOwner() {
        require(msg.sender == _owner, "Caller is not owner");
        _;
    }

    constructor() {
        _owner = msg.sender;
    }

    function owner() public view returns (address) {
        return _owner;
    }

    function transferOwnership(address newOwner) external onlyOwner {
        require(newOwner != address(0), "New owner is zero address");
        _owner = newOwner;
    }

    // -----------------------------------------------------------------------
    // Layer struct: holds one “off-chain” Bloom filter layer
    // -----------------------------------------------------------------------
    struct Layer {
        bytes[] chunks;
        uint256 chunkSizeBytes;
        uint256 filterSizeBits;
        uint256 k;
    }

    Layer[] private layers;

    /// @notice Replace the entire cascade in one transaction.
    /// @param newChunksByLayer 2D array: newChunksByLayer[i] is the []byte-chunks for layer i.
    /// @param ks               ks[i] is the number of hash functions used off-chain for layer i.
    function updateCascade(bytes[][] calldata newChunksByLayer, uint256[] calldata ks) external onlyOwner {
        require(newChunksByLayer.length == ks.length, "Must supply one k per layer");
        require(newChunksByLayer.length > 0, "Need at least one layer");

        // Clear existing layers
        delete layers;

        for (uint256 layerIdx = 0; layerIdx < newChunksByLayer.length; layerIdx++) {
            bytes[] calldata layerChunks = newChunksByLayer[layerIdx];
            uint256 k = ks[layerIdx];
            require(k > 0, "k must be > 0");
            require(layerChunks.length > 0, "Each layer needs more than 1 chunk");

            uint256 chunkSize = layerChunks[0].length;
            require(chunkSize > 0, "Chunks cannot be zero-length");

            uint256 totalBytes = 0;
            for (uint256 j = 0; j < layerChunks.length; j++) {
                if (j < layerChunks.length - 1) {
                    require(layerChunks[j].length == chunkSize, "Chunks mismatch size");
                }
                totalBytes += layerChunks[j].length;
            }

            layers.push();
            Layer storage L = layers[layers.length - 1];
            L.chunkSizeBytes = chunkSize;
            L.filterSizeBits = totalBytes * 8;
            L.k = k;

            L.chunks = new bytes[](layerChunks.length);
            for (uint256 j = 0; j < layerChunks.length; j++) {
                L.chunks[j] = layerChunks[j];
            }
        }
    }

    function layerCount() external view returns (uint256) {
        return layers.length;
    }

    function chunkCount(uint256 layerIdx) external view returns (uint256) {
        require(layerIdx < layers.length, "Invalid layer index");
        return layers[layerIdx].chunks.length;
    }

    function getChunk(uint256 layerIdx, uint256 chunkIdx) external view returns (bytes memory) {
        require(layerIdx < layers.length, "Invalid layer index");
        Layer storage L = layers[layerIdx];
        require(chunkIdx < L.chunks.length, "Invalid chunk index");
        return L.chunks[chunkIdx];
    }

    /// @notice Test a token against all layers. Returns true only if it passes every layer.
    function testToken(bytes calldata token) external view returns (bool) {
        uint256 numLayers = layers.length;
        require(numLayers > 0, "No layers initialized");

        // Compute base Keccak256 once:
        bytes32 digest = keccak256(token);
        uint256 big = uint256(digest);
        uint64 h1 = uint64(big >> 192);
        uint64 h2 = uint64((big >> 128) & 0xFFFFFFFFFFFFFFFF);

        for (uint256 li = 0; li < numLayers; li++) {
            // Inline to avoid excess stack variables
            uint256 sizeBits  = layers[li].filterSizeBits;
            uint256 cSize     = layers[li].chunkSizeBytes;
            uint256 layerK    = layers[li].k;
            bytes[] storage arr = layers[li].chunks;

            for (uint256 i = 0; i < layerK; i++) {
                uint256 combined = uint256(h1) + i * uint256(h2);
                uint256 bitPos   = combined % sizeBits;

                // Calculate byteIndex, then chunk & offset
                uint256 byteIndex = bitPos >> 3; // divide by 8
                uint256 chunkIndex = byteIndex / cSize;
                uint256 indexInChunk = byteIndex - (chunkIndex * cSize);

                bytes storage chosenChunk = arr[chunkIndex];
                uint8 rawByte = uint8(chosenChunk[indexInChunk]);
                uint8 bitInByte = uint8(bitPos & 7);

                if (((rawByte >> bitInByte) & 0x01) == 0) {
                    return false;
                }
            }
        }
        return true;
    }

    function totalBytesOfLayer(uint256 layerIdx) external view returns (uint256) {
        require(layerIdx < layers.length, "Invalid layer index");
        return layers[layerIdx].filterSizeBits / 8;
    }

    function getLayerMetadata(uint256 layerIdx)
    external
    view
    returns (
        uint256 chunkSizeBytes_,
        uint256 filterSizeBits_,
        uint256 k_
    )
    {
        require(layerIdx < layers.length, "Invalid layer index");
        Layer storage L = layers[layerIdx];
        return (L.chunkSizeBytes, L.filterSizeBits, L.k);
    }
}