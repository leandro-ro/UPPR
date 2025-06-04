// SPDX-License-Identifier: MIT
pragma solidity ^0.8.0;

/// @notice A “cascading” Bloom filter where each layer is stored as one big `bytes`.
contract CascadingBloomFilter {
    // --- ownable boilerplate (omitted for brevity) ---
    address private _owner;
    modifier onlyOwner() {
        require(msg.sender == _owner, "Not owner");
        _;
    }
    constructor() {
        _owner = msg.sender;
    }

    // --- LAYER STRUCT: one contiguous `bytes` per layer ---
    struct Layer {
        bytes   filter;         // the full bitvector in bytes
        uint256 filterSizeBits; // how many bits of `filter` are valid
        uint256 k;              // number of hash functions
    }
    Layer[] private layers;

    /// @notice Replace all layers in one call.
    /// @param newFilters   An array of `bytes` where `newFilters[i]` is the entire
    ///                     bitvector (packed) for layer i.
    /// @param ks           ks[i] = number of hash functions at layer i
    /// @param bitLens      bitLens[i] = number of bits used in that layer’s filter
    function updateCascade(
        bytes[] calldata newFilters,
        uint256[] calldata ks,
        uint256[] calldata bitLens
    ) external onlyOwner {
        require(newFilters.length == ks.length,     "Need k for each layer");
        require(ks.length == bitLens.length,        "Need bitLen for each layer");
        require(newFilters.length > 0,              "At least one layer");

        // Wipe out old layers
        delete layers;

        for (uint256 i = 0; i < newFilters.length; i++) {
            bytes calldata f = newFilters[i];
            uint256 k_   = ks[i];
            uint256 bits = bitLens[i];

            require(k_ > 0,                         "k must be > 0");
            require(bits > 0,                       "bitLen must be > 0");
            require(bits <= f.length * 8,          "bitLen exceeds f.length*8");

            // Push a new Layer and populate it
            layers.push();
            Layer storage L = layers[layers.length - 1];
            L.filter          = f;      // copy the entire byte-array
            L.filterSizeBits  = bits;
            L.k               = k_;
        }
    }

    /// @notice Number of layers in the cascade
    function layerCount() external view returns (uint256) {
        return layers.length;
    }

    /// @notice Test `token` against every layer. Returns (accepted, layerIndexReached).
    function testToken(bytes calldata token) external view returns (bool, uint256) {
        uint256 n = layers.length;
        require(n > 0, "No layers");

        uint64[4] memory h = extractHashes(token);

        for (uint256 li = 0; li < n; li++) {
            Layer storage L = layers[li];
            bool match_ = testInLayer(L.filter, L.filterSizeBits, L.k, h);

            if (li == n - 1) {
                bool expected = (li % 2 == 0); // even = expected match
                return (match_ == expected, li);
            }

            if (!match_) {
                bool acceptEarly = (li % 2 == 1); // odd layers: accept on mismatch
                return (acceptEarly, li);
            }
        }

        revert("unreachable");
    }

    /// @notice Test a single layer’s `filter` as a Bloom filter.
    /// @param filter         the full bytes‐packed bitvector
    /// @param filterSizeBits how many bits of `filter` are actually used
    /// @param k              how many hash‐probes to do
    /// @param h              the 4×64‐bit “precomputed” hashes
    function testInLayer(
        bytes storage filter,
        uint256 filterSizeBits,
        uint256 k,
        uint64[4] memory h
    ) internal view returns (bool) {
        for (uint256 i = 0; i < k; i++) {
            uint256 bitPos = getLocation(h, i, filterSizeBits);

            uint256 byteIndex = bitPos >> 3;
            require(byteIndex < filter.length, "byteIndex out of range");

            uint8 bitOffset = uint8(bitPos & 7);

            uint8 b = uint8(filter[byteIndex]);

            if (((b >> bitOffset) & 1) == 0) {
                return false;
            }
        }

        return true;
    }

    /// @notice Return metadata and full filter bytes for one layer.
    /// @param i Layer index
    /// @return filterSizeBits_ number of valid bits in filter
    /// @return k_ number of hash functions
    /// @return filter_ the actual bytes of the filter (bitvector)
        function getLayerMetadata(uint256 i)
        external
        view
        returns (
            uint256 filterSizeBits_,
            uint256 k_,
            bytes memory filter_
        )
        {
            require(i < layers.length, "Invalid layer");
            Layer storage L = layers[i];
            return (
                L.filterSizeBits,
                L.k,
                L.filter
            );
        }

    /// @notice Extract four 64-bit values from Keccak256(token)
    function extractHashes(bytes calldata token) internal pure returns (uint64[4] memory h) {
        bytes32 digest = keccak256(token);
        h[0] = uint64(uint256(digest >> 192));
        h[1] = uint64((uint256(digest) >> 128) & 0xFFFFFFFFFFFFFFFF);
        h[2] = uint64((uint256(digest) >> 64)  & 0xFFFFFFFFFFFFFFFF);
        h[3] = uint64(uint256(digest)         & 0xFFFFFFFFFFFFFFFF);
    }

    /// @notice “Go‐style” mixing of the four 64-bit hashes into k locations
    function getLocation(uint64[4] memory h, uint256 i, uint256 mod) internal pure returns (uint256) {
        uint64 ii    = uint64(i);
        uint64 base  = h[ii % 2];
        uint64 idx   = 2 + (((ii + (ii % 2)) % 4) / 2);
        uint64 mult  = h[idx];
        uint64 sum64;
        unchecked {
            sum64 = base + ii * mult;
        }
        return uint256(sum64) % mod;
    }
}